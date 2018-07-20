package consul

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

// 一直执行注册直到成功
func registerServiceAlways(consulAddrs, dataCenter string, id, name, addr string, port int, checkUrl string) {
	for {
		err := RegisterService(consulAddrs, dataCenter, id, name, addr, port, checkUrl)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func RegisterService(consulAddrs, dataCenter string, id, name, addr string, port int, checkUrl string) error {
	var regErr error
	for _, consulAddr := range strings.Split(consulAddrs, ",") {
		config := api.DefaultConfig()
		config.Address = consulAddr
		config.Datacenter = dataCenter

		client, err := api.NewClient(config)
		if err != nil {
			continue
		}
		asr := &api.AgentServiceRegistration{
			ID:      id,
			Name:    name,
			Address: addr,
			Port:    port,
			Tags:    []string{name},
			Check: &api.AgentServiceCheck{
				HTTP:     fmt.Sprintf("http://%s:%d%s", addr, port, checkUrl),
				Method:   "GET",
				Status:   api.HealthPassing,
				Interval: "5s",
			},
		}
		agent := client.Agent()
		regErr = agent.ServiceRegister(asr)
		if regErr == nil {
			return nil
		}
	}
	go registerServiceAlways(consulAddrs, dataCenter, id, name, addr, port, checkUrl)
	return regErr
}

func DeRegisterService(consulAddrs, dataCenter, id string) error {
	var deRegErr error
	for _, consulAddr := range strings.Split(consulAddrs, ",") {
		config := api.DefaultConfig()
		config.Address = consulAddr
		config.Datacenter = dataCenter

		client, err := api.NewClient(config)
		if err != nil {
			continue
		}
		deRegErr = client.Agent().ServiceDeregister(id)
		if deRegErr == nil {
			return nil
		}
	}
	return deRegErr
}

func GetServiceEndpoints(consulAddrs, dataCenter, name string, passingOnly bool) ([]string, error) {
	var addrs []string
	var consulErr error
	for _, consulAddr := range strings.Split(consulAddrs, ",") {
		config := api.DefaultConfig()
		config.Address = consulAddr
		config.Datacenter = dataCenter

		client, err := api.NewClient(config)
		if err != nil {
			consulErr = err
			continue
		}
		services, _, err := client.Health().Service(name, name, passingOnly, &api.QueryOptions{AllowStale: true})
		if err != nil {
			consulErr = err
			continue
		}
		for _, v := range services {
			if v.Service != nil {
				addr := fmt.Sprintf("%s:%d", v.Service.Address, v.Service.Port)
				addrs = append(addrs, addr)
			}
		}
		return addrs, nil
	}
	if consulErr == nil {
		consulErr = fmt.Errorf("not found endpoints by name:%s", name)
	}
	return []string{}, consulErr
}
