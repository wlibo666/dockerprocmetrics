package metrics

import (
	"bytes"
	"fmt"
	"io"

	"github.com/wlibo666/dockerprocmetrics/container"
	"github.com/wlibo666/procinfo"
	"github.com/wlibo666/procinfo/pid"
)

var (
	metricsItem []MetericItemFunc
)

const (
	METRICS_PREFIX = "dpm_"

	NODE_SYS_CPU_USAGE_RATE = METRICS_PREFIX + "node_sys_cpu_usage_rate"
	NODE_SYS_MEM_USAGE      = METRICS_PREFIX + "node_sys_mem_usage"
	CONTAINER_PID_CPU_UAGE  = METRICS_PREFIX + "container_pid_cpu_usage"
	CONTAINER_PID_MEM_USAGE = METRICS_PREFIX + "container_pid_mem_usage"
	CONTAINER_PID_COUNT     = METRICS_PREFIX + "container_pid_count"
	CONTAINER_STATE         = METRICS_PREFIX + "container_state"

	// metrics args
	CONTAINER_STATE_ARG_REMOTEADDR = "remoteaddr"
)

var (
	node_sys_cpu_usage_rate_comment = fmt.Sprintf("# HELP %s Cpu usage rate on the machine.\n", NODE_SYS_CPU_USAGE_RATE) +
		fmt.Sprintf("# TYPE %s gauge\n", NODE_SYS_CPU_USAGE_RATE)
	node_sys_mem_usage_comment = fmt.Sprintf("# HELP %s memory usage on the machine.\n", NODE_SYS_MEM_USAGE) +
		fmt.Sprintf("# TYPE %s gauge\n", NODE_SYS_MEM_USAGE)
	container_pid_cpu_uage = fmt.Sprintf("# HELP %s cpu usage of process.\n", CONTAINER_PID_CPU_UAGE) +
		fmt.Sprintf("# TYPE %s gauge\n", CONTAINER_PID_CPU_UAGE)
	container_pid_mem_usage = fmt.Sprintf("# HELP %s memory usage of process.\n", CONTAINER_PID_MEM_USAGE) +
		fmt.Sprintf("# TYPE %s gauge\n", CONTAINER_PID_MEM_USAGE)
	container_pid_count = fmt.Sprintf("# HELP %s process count of container.\n", CONTAINER_PID_COUNT) +
		fmt.Sprintf("# TYPE %s gauge\n", CONTAINER_PID_COUNT)
	container_pid_state = fmt.Sprintf("# HELP %s state of container.\n", CONTAINER_STATE) +
		fmt.Sprintf("# TYPE %s gauge\n", CONTAINER_STATE)
)

func AddMetricItem(name, comment string, metricfun func(map[string]string) (string, error)) {
	item := MetericItemFunc{
		ItemName: name,
		Comment:  comment,
		Func:     metricfun,
	}
	metricsItem = append(metricsItem, item)
}

type MetericItemFunc struct {
	ItemName string
	Comment  string
	Func     func(args map[string]string) (string, error)
}

func SysCpuUsageRateMetric(args map[string]string) (string, error) {
	rate := procinfo.GetCpuUsageRate()
	itemData := fmt.Sprintf("%s{%s} %.2f\n", NODE_SYS_CPU_USAGE_RATE, container.HostLabels, 1.0-rate.Id)
	return itemData, nil
}

func SysMemUsageMetric(args map[string]string) (string, error) {
	var data []byte
	buff := bytes.NewBuffer(data)

	minfo := procinfo.GetMemInfo()
	itemData := fmt.Sprintf("%s{%s,free=\"%d\",avaliable=\"%d\"} %d\n",
		NODE_SYS_MEM_USAGE, container.HostLabels, minfo.Free, minfo.Available, minfo.Total-minfo.Available)
	buff.Write([]byte(itemData))
	return buff.String(), nil
}

func PidCpuUsageMetric(args map[string]string) (string, error) {
	var data []byte
	buff := bytes.NewBuffer(data)
	// key=pid value=containerId
	container.PidsMonitor.Range(func(key, value interface{}) bool {
		// 根据容器ID获取容器信息
		containerInfo, ok := container.ContainerInfos.Load(value)
		if !ok {
			return false
		}
		rate, err := pid.GetProcCpuUsageRate(key.(string))
		if err == nil {
			cmdLine, _ := pid.GetPidCmdline(key.(string))
			itemData := fmt.Sprintf("%s{%s,%s,pid=\"%s\",cmdline=\"%s\"} %.2f\n",
				CONTAINER_PID_CPU_UAGE, container.HostLabels, containerInfo.(*container.ContainerInfo).LablesStr,
				key.(string), cmdLine, rate.Rate)
			buff.Write([]byte(itemData))
		}
		return true
	})
	return buff.String(), nil
}

func PidMemUsageMetric(args map[string]string) (string, error) {
	var data []byte
	buff := bytes.NewBuffer(data)
	// key=pid value=containerId
	container.PidsMonitor.Range(func(key, value interface{}) bool {
		// 根据容器ID获取容器信息
		containerInfo, ok := container.ContainerInfos.Load(value)
		if !ok {
			return false
		}
		minfo, err := pid.GetProcMemInfo(key.(string))
		if err == nil {
			cmdLine, _ := pid.GetPidCmdline(key.(string))
			itemData := fmt.Sprintf("%s{%s,%s,pid=\"%s\",cmdline=\"%s\"} %d\n",
				CONTAINER_PID_MEM_USAGE, container.HostLabels, containerInfo.(*container.ContainerInfo).LablesStr,
				key.(string), cmdLine, minfo.VmRss)
			buff.Write([]byte(itemData))
		}
		return true
	})
	return buff.String(), nil
}

func ContainerPidCntMetric(args map[string]string) (string, error) {
	var data []byte
	buff := bytes.NewBuffer(data)

	container.ContainerPidNum.Range(func(key, value interface{}) bool {
		// 根据容器ID获取容器信息
		containerInfo, ok := container.ContainerInfos.Load(key)
		if !ok {
			return false
		}

		itemData := fmt.Sprintf("%s{%s,%s} %d\n",
			CONTAINER_PID_COUNT, container.HostLabels, containerInfo.(*container.ContainerInfo).LablesStr,
			value.(int))
		buff.Write([]byte(itemData))
		return true
	})
	return buff.String(), nil
}

func ContainerStateMetric(args map[string]string) (string, error) {
	var data []byte
	buff := bytes.NewBuffer(data)

	container.EventDatas.Range(func(key, value interface{}) bool {
		// 查询采集数据的server的地址
		if args != nil {
			addr, ok := args[CONTAINER_STATE_ARG_REMOTEADDR]
			if ok {
				// 如果该节点已采集过,忽略
				_, accessed := value.(*container.EventData).Accessd.Load(addr)
				if accessed {
					return false
				}
			}

			// 如果没采集过,则采集数据并标为已采集
			itemData := fmt.Sprintf("%s%s\n", CONTAINER_STATE, value.(*container.EventData).Data)
			value.(*container.EventData).Accessd.Store(addr, true)
			buff.Write([]byte(itemData))
			return true
		}
		evData := fmt.Sprintf("%s%s\n", CONTAINER_STATE, value.(*container.EventData).Data)
		buff.Write([]byte(evData))
		return true
	})
	return buff.String(), nil
}

func WriteAllMetricsData(writer io.Writer, args map[string]string) error {
	for _, item := range metricsItem {
		data, err := item.Func(args)
		if err != nil {
			continue
		}
		if len(data) > 0 {
			writer.Write([]byte(item.Comment))
			writer.Write([]byte(data))
		}
	}
	return nil
}

func init() {
	AddMetricItem(NODE_SYS_CPU_USAGE_RATE, node_sys_cpu_usage_rate_comment, SysCpuUsageRateMetric)
	AddMetricItem(NODE_SYS_MEM_USAGE, node_sys_mem_usage_comment, SysMemUsageMetric)
	AddMetricItem(CONTAINER_PID_CPU_UAGE, container_pid_cpu_uage, PidCpuUsageMetric)
	AddMetricItem(CONTAINER_PID_MEM_USAGE, container_pid_mem_usage, PidMemUsageMetric)
	AddMetricItem(CONTAINER_PID_COUNT, container_pid_count, ContainerPidCntMetric)
	AddMetricItem(CONTAINER_STATE, container_pid_state, ContainerStateMetric)
}
