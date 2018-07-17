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

type MetricConfig struct {
	DockerDaemonSock string   `json:"dockerDaemonSock"`
	DockerApiVersion string   `json:"dockerApiVersion"`
	MonitorItems     []string `json:"monitorItems"`
	Frequency        int      `json:"frequency"`
	ListenAddr       string   `json:"listenAddr"`
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
	if config.DockerDaemonSock == "" {
		config.DockerDaemonSock = DEFAULT_DOCKER_DAEMON_SOCK
	}
	if config.Frequency <= 0 {
		config.Frequency = DEFAULT_FREQUENCY
	}
	if len(config.MonitorItems) == 0 {
		config.MonitorItems = append(config.MonitorItems, MONITOR_CPU)
		config.MonitorItems = append(config.MonitorItems, MONITOR_MEMORY)
	}
	if config.ListenAddr == "" {
		config.ListenAddr = DEFAULT_LISTEN_ADDR
	}
	return nil
}
