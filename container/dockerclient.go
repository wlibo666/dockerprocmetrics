package container

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	log "github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/dockerprocmetrics/utils"
)

var (
	DefaultDockerClient *client.Client

	dockerAddr string
	apiVersion string
)

func newDockerClient(addr, apiversion string) (*client.Client, error) {
	return client.NewClientWithOpts(client.WithVersion(apiversion), client.WithHost(addr))
}

func NewDockerClient() *client.Client {
	client, err := newDockerClient(dockerAddr, apiVersion)
	if err != nil {
		return DefaultDockerClient
	}
	return client
}

func DefaultDockerClientCheck(addr, apiversion string) {
	for {
		_, err := DefaultDockerClient.Ping(utils.NewContextWithTimeout(3))
		if err != nil {
			log.DefFileLogger.WithFields(logrus.Fields{
				"apiversion": apiversion,
				"host":       addr,
				"position":   utils.GetFileAndLine(),
				"error":      err.Error(),
			}).Warn("ping docker daemon failed")
			DefaultDockerClient.Close()
			NewDefaultDockerClient(addr, apiversion)
		}
		time.Sleep(1 * time.Second)
	}
}

func NewDefaultDockerClient(addr, apiversion string) error {
	client, err := newDockerClient(addr, apiversion)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"apiversion": apiversion,
			"host":       addr,
			"position":   utils.GetFileAndLine(),
			"error":      err.Error(),
		}).Error("NewClientWithOpts failed")
		return err
	}
	DefaultDockerClient = client
	dockerAddr = addr
	apiVersion = apiversion

	return nil
}
