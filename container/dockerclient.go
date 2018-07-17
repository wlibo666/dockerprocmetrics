package container

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	log "github.com/wlibo666/common-lib/logrus"
)

var (
	DefaultDockerClient *client.Client
)

func NewDockerClient(addr, apiversion string) (*client.Client, error) {
	return client.NewClientWithOpts(client.WithVersion(apiversion), client.WithHost(addr))
}

func DefaultDockerClientCheck(addr, apiversion string) {
	for {
		_, err := DefaultDockerClient.Ping(context.Background())
		if err != nil {
			log.DefFileLogger.WithFields(logrus.Fields{
				"apiversion": apiversion,
				"host":       addr,
				"error":      err.Error(),
			}).Warn("ping docker daemon failed")
			NewDefaultDockerClient(addr, apiversion)
		}
		time.Sleep(time.Second)
	}
}

func NewDefaultDockerClient(addr, apiversion string) error {
	client, err := NewDockerClient(addr, apiversion)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"apiversion": apiversion,
			"host":       addr,
			"error":      err.Error(),
		}).Error("NewClientWithOpts failed")
		return err
	}
	DefaultDockerClient = client

	return nil
}
