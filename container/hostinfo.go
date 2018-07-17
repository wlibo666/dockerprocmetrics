package container

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	log "github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/procinfo"
)

var (
	HostLabels string
)

func InitHostLabels() error {
	procinfo.InitHostInfo()

	HostLabels = fmt.Sprintf("host_name=\"%s\",host_cpu=\"%d\",host_memory=\"%d\",mem_unit=\"kB\"",
		procinfo.Host.HostName, procinfo.Host.CpuCnt, procinfo.Host.Memory)
	log.DefFileLogger.WithFields(logrus.Fields{
		"HostLabels": HostLabels,
	}).Info("InitHostLabels ok")
	return nil
}
