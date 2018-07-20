package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	DEFAULT_DOCKER_DAEMON_SOCK = "unix:///var/run/docker.sock"
	DEFAULT_FREQUENCY          = 1
	DEFAULT_LISTEN_ADDR        = ":3280"
	MONITOR_CPU                = "cpu"
	MONITOR_MEMORY             = "memory"
)

var (
	GMertricConfig *MetricConfig
)

type DockerConfig struct {
	DaemonSock string `json:"daemonSock"`
	ApiVersion string `json:"apiVersion"`
}

type MonitorConfig struct {
	Items     []string `json:"items"`
	Frequency int      `json:"frequency"`
}

type ListenConfig struct {
	Addr    string `json:"addr"`
	Port    int    `json:"port"`
	WanName string `json:"wanName"`
}

type ConsulRegisterConfig struct {
	Addr        string `json:"addr"`
	Dc          string `json:"dc"`
	ServiceName string `json:"serviceName"`
}

type MetricConfig struct {
	Docker    DockerConfig         `json:"docker"`
	Monitor   MonitorConfig        `json:"monitor"`
	Listen    ListenConfig         `json:"listen"`
	ConsulReg ConsulRegisterConfig `json:"consulRegister"`
}

func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &GMertricConfig)
	if err != nil {
		return err
	}
	return GMertricConfig.check()
}

func (config *MetricConfig) check() error {
	if config.Docker.DaemonSock == "" {
		config.Docker.DaemonSock = DEFAULT_DOCKER_DAEMON_SOCK
	}
	if config.Monitor.Frequency <= 0 {
		config.Monitor.Frequency = DEFAULT_FREQUENCY
	}
	if len(config.Monitor.Items) == 0 {
		config.Monitor.Items = append(config.Monitor.Items, MONITOR_CPU)
		config.Monitor.Items = append(config.Monitor.Items, MONITOR_MEMORY)
	}

	return nil
}
