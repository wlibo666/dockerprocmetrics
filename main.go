package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	log "github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/dockerprocmetrics/config"
	"github.com/wlibo666/dockerprocmetrics/controller"
	"github.com/wlibo666/dockerprocmetrics/metrics"
	"github.com/wlibo666/dockerprocmetrics/register/consul"
	"github.com/wlibo666/dockerprocmetrics/utils"
	"github.com/wlibo666/procinfo"
)

var (
	logFile  = flag.String("log.file", "/var/log/dpm.log", "-log.file /var/log/dpm.log")
	logCnt   = flag.Int("log.count", 3, "-log.count 3")
	confFile = flag.String("config.file", "/etc/dpm/metrics.json", "-config.file ../conf/metrics.json")
)

var (
	serviceId string
)

func prepare() {
	flag.Parse()
	// 初始化日志
	err := log.NewDefaultFileLogger(*logFile, *logCnt, log.LOG_FORMAT_JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new log file failed,err:%s,will exit with code 1\n", err.Error())
		utils.ExitWaitDef(1)
	}
	log.SetDefaultLoggerLevel(logrus.InfoLevel)
	// 加载配置文件
	err = config.LoadConfig(*confFile)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"config.file": *confFile,
			"position":    utils.GetFileAndLine(),
			"error":       err.Error(),
		}).Error("LoadConfig failed,will exit with code 2")
		utils.ExitWaitDef(2)
	}
	if config.GMertricConfig.Docker.ProcDir != "" {
		procinfo.SetProcBaseDir(config.GMertricConfig.Docker.ProcDir)
	}
}

func sigAction() {
	go func() {
		// 注册信号处理
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		log.DefFileLogger.WithFields(logrus.Fields{
			"position": utils.GetFileAndLine(),
			"signal":   s.String(),
		}).Info("receive signal,program will exit with 0")
		// 退出前取消服务注册
		err := consul.DeRegisterService(config.GMertricConfig.ConsulReg.Addr, config.GMertricConfig.ConsulReg.Dc, serviceId)
		if err != nil {
			log.DefFileLogger.WithFields(logrus.Fields{
				"serviceId": serviceId,
				"position":  utils.GetFileAndLine(),
				"error":     err.Error(),
			}).Warn("DeRegisterService service failed")
		} else {
			log.DefFileLogger.WithFields(logrus.Fields{
				"position":  utils.GetFileAndLine(),
				"serviceId": serviceId,
			}).Info("DeRegisterService service success")
		}
		os.Exit(0)
	}()
}

func register() error {
	var moniAddr string
	// 优先使用网卡名称获取监听地址
	addr, err := utils.GetWanAddr(config.GMertricConfig.Listen.WanName)
	if err == nil {
		serviceId = config.GMertricConfig.ConsulReg.ServiceName + "-" + addr
	} else {
		serviceId = config.GMertricConfig.ConsulReg.ServiceName
	}
	if addr != "" {
		moniAddr = addr
	} else {
		moniAddr = config.GMertricConfig.Listen.Addr
	}
	// 在consul注册服务
	err = consul.RegisterService(config.GMertricConfig.ConsulReg.Addr, config.GMertricConfig.ConsulReg.Dc,
		serviceId, config.GMertricConfig.ConsulReg.ServiceName, moniAddr, config.GMertricConfig.Listen.Port,
		controller.CTL_URL_HEALTHZ)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"regType":     "consul",
			"serviceName": config.GMertricConfig.ConsulReg.ServiceName,
			"serviceID":   serviceId,
			"consulAddr":  config.GMertricConfig.ConsulReg.Addr,
			"consulDc":    config.GMertricConfig.ConsulReg.Dc,
			"listenAddr":  moniAddr,
			"listenPort":  config.GMertricConfig.Listen.Port,
			"position":    utils.GetFileAndLine(),
			"error":       err.Error(),
		}).Warn("RegisterService failed")
		return err
	}
	log.DefFileLogger.WithFields(logrus.Fields{
		"regType":     "consul",
		"serviceName": config.GMertricConfig.ConsulReg.ServiceName,
		"serviceID":   serviceId,
		"consulAddr":  config.GMertricConfig.ConsulReg.Addr,
		"consulDc":    config.GMertricConfig.ConsulReg.Dc,
		"listenAddr":  moniAddr,
		"listenPort":  config.GMertricConfig.Listen.Port,
		"position":    utils.GetFileAndLine(),
	}).Info("RegisterService success")
	return nil
}

func RunGinServer() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// 记录web访问日志
	engine.Use(log.MiddleAccessLog)
	engine.GET(controller.CTL_URL_MTERICS, controller.Metrics)
	engine.GET(controller.CTL_URL_HEALTHZ, controller.Health)
	register()
	return engine.Run(utils.GenListenAddr(config.GMertricConfig.Listen.Addr, config.GMertricConfig.Listen.Port))
}

func main() {
	// 加载配置文件并初始化日志文件
	prepare()
	// 添加信号处理
	sigAction()
	// 启动metric协程
	err := metrics.Run()
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"position": utils.GetFileAndLine(),
			"error":    err.Error(),
		}).Error("metrics Run failed,will exit with code 4")
		utils.ExitWaitDef(3)
	}
	// 启动web服务
	err = RunGinServer()
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"position": utils.GetFileAndLine(),
			"error":    err.Error(),
		}).Error("RunGinServer failed,will exit with code 5")
		utils.ExitWaitDef(5)
	}

	log.DefFileLogger.WithFields(logrus.Fields{
		"position": utils.GetFileAndLine(),
		"error":    "",
	}).Info("program exit with code 0")
}
