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
	"github.com/wlibo666/dockerprocmetrics/utils"
)

var (
	logFile  = flag.String("log.file", "/var/log/dockerprocmetrics.log", "-log.file /var/log/dockerprocmetrics.log")
	logCnt   = flag.Int("log.count", 3, "-log.count 3")
	confFile = flag.String("config.file", "./conf/metrics.json", "-config.file ../conf/metrics.json")
)

func prepare() {
	flag.Parse()
	err := log.NewDefaultFileLogger(*logFile, *logCnt, log.LOG_FORMAT_JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new log file failed,err:%s,will exit with code 1\n", err.Error())
		utils.ExitWaitDef(1)
	}
	log.SetDefaultLoggerLevel(logrus.InfoLevel)

	err = config.LoadConfig(*confFile)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"config.file": *confFile,
			"error":       err.Error(),
		}).Error("LoadConfig failed,will exit with code 2")
		utils.ExitWaitDef(2)
	}
}

func sigAction() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		log.DefFileLogger.WithFields(logrus.Fields{
			"signal": s.String(),
		}).Info("receive signal,program will exit with 0")
		os.Exit(0)
	}()
}

func RunGinServer() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// 记录web访问日志
	engine.Use(log.MiddleAccessLog)
	engine.GET("/metrics", controller.Metrics)
	engine.Run(config.GMertricConfig.ListenAddr)
	return nil
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
			"error": err.Error(),
		}).Error("metrics Run failed,will exit with code 4")
		utils.ExitWaitDef(3)
	}
	// 启动web服务
	RunGinServer()

	log.DefFileLogger.WithFields(logrus.Fields{
		"error": "",
	}).Info("program exit with code 0")
}
