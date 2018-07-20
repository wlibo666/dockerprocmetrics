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
	err := container.NewDefaultDockerClient(config.GMertricConfig.Docker.DaemonSock, config.GMertricConfig.Docker.ApiVersion)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"apiversion": config.GMertricConfig.Docker.ApiVersion,
			"host":       config.GMertricConfig.Docker.DaemonSock,
			"position":   utils.GetFileAndLine(),
			"error":      err.Error(),
		}).Error("NewDefaultDockerClient failed, will exit with code 3")
		utils.ExitWaitDef(10)
	}
	// 检测daemon是否存活
	go container.DefaultDockerClientCheck(config.GMertricConfig.Docker.DaemonSock, config.GMertricConfig.Docker.ApiVersion)
	// 从docker daemon获取容器信息和容器内PID
	go container.StoreContainerInfoAndPid()
	// 处理容器事件
	go container.HandleDockerEvent()
	return nil
}
