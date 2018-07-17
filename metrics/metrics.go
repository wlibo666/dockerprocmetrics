package metrics

import (
	"github.com/Sirupsen/logrus"
	log "github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/dockerprocmetrics/config"
	"github.com/wlibo666/dockerprocmetrics/container"
	"github.com/wlibo666/dockerprocmetrics/utils"
)

func Run() error {
	// 初始化主机固有信息(主机名称/CPU总数/内存总量)
	container.InitHostLabels()

	// 初始化docker客户端
	err := container.NewDefaultDockerClient(config.GMertricConfig.DockerDaemonSock, config.GMertricConfig.DockerApiVersion)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"apiversion": config.GMertricConfig.DockerApiVersion,
			"host":       config.GMertricConfig.DockerDaemonSock,
			"error":      err.Error(),
		}).Error("NewDefaultDockerClient failed, will exit with code 3")
		utils.ExitWaitDef(10)
	}
	// 检测daemon是否存活
	go container.DefaultDockerClientCheck(config.GMertricConfig.DockerDaemonSock, config.GMertricConfig.DockerApiVersion)
	// 从docker daemon获取容器信息和容器内PID
	err = container.StoreContainerInfoAndPid()
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("StoreContainerInfoAndPid failed")
		return err
	}
	go container.HandleDockerEvent()
	return nil
}
