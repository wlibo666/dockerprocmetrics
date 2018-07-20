package container

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	log "github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/dockerprocmetrics/utils"
	procpid "github.com/wlibo666/procinfo/pid"
)

type ContainerInfo struct {
	Labels    map[string]string
	LablesStr string
}

var (
	// map[containerId]*ContainerInfo
	ContainerInfos *sync.Map = &sync.Map{}
	// map[containerId]int
	ContainerPidNum *sync.Map = &sync.Map{}
	// map[pid]containerId
	PidsMonitor *sync.Map = &sync.Map{}
)

func DelContainerInfoById(id string) {
	_, ok := ContainerInfos.Load(id)
	if ok {
		ContainerInfos.Delete(id)
	}
}

func handleContainerError(id string, err error) {
	if strings.Contains(err.Error(), "is not running") {
		ContainerInfos.Delete(id)
		ContainerPidNum.Delete(id)
	}
}

func StoreContainerInfoById(id string) error {
	// 如果已经保存,则不需要更改信息
	_, ok := ContainerInfos.Load(id)
	if ok {
		return nil
	}
	cinfo := &ContainerInfo{
		Labels: make(map[string]string),
	}
	// 获取单个容器信息
	client := NewDockerClient()
	defer client.Close()
	info, err := client.ContainerInspect(utils.NewContextWithTimeout(3), id)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"containerId": id,
			"position":    utils.GetFileAndLine(),
			"error":       err.Error(),
		}).Warn("ContainerInspect failed")
		handleContainerError(id, err)
		return err
	}

	// 从容器信息里获取必要字段
	cinfo.Labels["id"] = info.ID
	cinfo.Labels["created"] = info.Created
	cinfo.Labels["image"] = info.Image
	cinfo.Labels["name"] = info.Name
	if info.HostConfig != nil {
		if info.HostConfig.Privileged {
			cinfo.Labels["privileged"] = "1"
		} else {
			cinfo.Labels["privileged"] = "0"
		}
	}
	if info.Config != nil {
		cinfo.Labels["k8s_hostname"] = info.Config.Hostname
		if info.Config.Labels != nil {
			for k, v := range info.Config.Labels {
				switch k {
				case "annotation.io.kubernetes.container.restartCount":
					cinfo.Labels["container_restart_cnt"] = v
				case "io.kubernetes.container.logpath":
					cinfo.Labels["container_logpath"] = v
				case "io.kubernetes.container.name":
					cinfo.Labels["container_name"] = v
				case "io.kubernetes.pod.name":
					cinfo.Labels["pod_name"] = v
				case "io.kubernetes.pod.namespace":
					cinfo.Labels["pod_ns"] = v
				case "io.kubernetes.pod.uid":
					cinfo.Labels["pod_uid"] = v
				}
			}
		}
	}
	var labels []string
	for k, v := range cinfo.Labels {
		labels = append(labels, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	cinfo.LablesStr = strings.Join(labels, ",")
	ContainerInfos.Store(id, cinfo)
	log.DefFileLogger.WithFields(logrus.Fields{
		"containerId": id,
		"position":    utils.GetFileAndLine(),
		"labels":      cinfo.LablesStr,
	}).Info("add container info by id ok")
	return nil
}

func StorePidsByContainerId(id string) error {
	client := NewDockerClient()
	defer client.Close()
	body, err := client.ContainerTop(utils.NewContextWithTimeout(3), id, []string{})
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			"containerId": id,
			"position":    utils.GetFileAndLine(),
			"error":       err.Error(),
		}).Warn("ContainerTop failed")
		handleContainerError(id, err)
		return err
	}
	// 记录容器内进程数量
	pidNum := len(body.Processes)
	ContainerPidNum.Store(id, pidNum)

	for _, pids := range body.Processes {
		pid := pids[1]
		log.DefFileLogger.WithFields(logrus.Fields{
			"containerId": id,
			"pid":         pid,
			"position":    utils.GetFileAndLine(),
		}).Debug("will add pid from container")
		v, ok := PidsMonitor.Load(pid)
		if ok {
			if v.(string) != id {
				PidsMonitor.Store(pid, id)
				log.DefFileLogger.WithFields(logrus.Fields{
					"containerId": id,
					"pid":         pid,
					"position":    utils.GetFileAndLine(),
				}).Info("update pid ok")
				procpid.MoniPidCpu(pid)
			}
		} else {
			PidsMonitor.Store(pid, id)
			log.DefFileLogger.WithFields(logrus.Fields{
				"containerId": id,
				"pid":         pid,
				"position":    utils.GetFileAndLine(),
			}).Info("add pid ok")
			procpid.MoniPidCpu(pid)
		}
	}
	return nil
}

// 每秒钟检测一次pid是否还存在
func checkPids() {
	go func() {
		for {
			PidsMonitor.Range(func(key, value interface{}) bool {
				// 检测PID是否存在
				dir := fmt.Sprintf(procpid.PROC_DIR, key.(string))
				_, err := os.Stat(dir)
				if err == os.ErrNotExist {
					PidsMonitor.Delete(key)
				}
				return true
			})
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for {
			ContainerInfos.Range(func(key, value interface{}) bool {
				// 查看容器内是否有新进程启动
				StorePidsByContainerId(key.(string))
				return true
			})
			time.Sleep(3 * time.Second)
		}
	}()
}

func StoreContainerInfoAndPid() error {
	// 定时检测PID是否还存活
	checkPids()

	for {
		// 获取所有容器列表
		client := NewDockerClient()
		containers, err := client.ContainerList(utils.NewContextWithTimeout(3), types.ContainerListOptions{})
		if err != nil {
			client.Close()
			log.DefFileLogger.WithFields(logrus.Fields{
				"position": utils.GetFileAndLine(),
				"error":    err.Error(),
			}).Warn("get container list failed")
			time.Sleep(5 * time.Second)
			continue
		}
		client.Close()
		// 获取容器信息及PID列表
		for _, container := range containers {
			StoreContainerInfoById(container.ID)
			StorePidsByContainerId(container.ID)
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}
