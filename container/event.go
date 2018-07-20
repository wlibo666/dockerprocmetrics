package container

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	log "github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/dockerprocmetrics/utils"
)

type EventData struct {
	// 事件label
	Data string
	// 事件发生事件
	Time int64
	// 该事件被哪些IP采集过 map[string]bool
	Accessd *sync.Map
}

const (
	CONTAINER_STATE_START = 1
	CONTAINER_STATE_DIE   = 0

	MAX_EVENT_ALIVE = 300
)

var (
	// map[int64]*EventData
	EventDatas *sync.Map = &sync.Map{}
)

func cleanEventData() {
	now := time.Now().Unix()
	EventDatas.Range(func(key, value interface{}) bool {
		if now-value.(*EventData).Time >= MAX_EVENT_ALIVE {
			EventDatas.Delete(key)
		}
		return true
	})
}

func handleContainerStartMsg(msg events.Message) error {
	log.DefFileLogger.WithFields(logrus.Fields{
		"position": utils.GetFileAndLine(),
		"msg":      msg,
	}).Info("receive start msg from docker daemon")

	StoreContainerInfoById(msg.ID)
	// 生成容器启动事件
	ev := &EventData{
		Data: fmt.Sprintf("{status=\"start\",container_id=\"%s\",image_id=\"%s\",type=\"container\"} %d",
			msg.ID, msg.From, CONTAINER_STATE_START),
		Time:    time.Now().Unix(),
		Accessd: &sync.Map{},
	}
	EventDatas.Store(msg.TimeNano, ev)

	StorePidsByContainerId(msg.ID)
	return nil
}

func handleContainerDieMsg(msg events.Message) error {
	log.DefFileLogger.WithFields(logrus.Fields{
		"position": utils.GetFileAndLine(),
		"msg":      msg,
	}).Info("receive die msg from docker daemon")
	DelContainerInfoById(msg.ID)

	// 生成容器结束事件
	ev := &EventData{
		Data: fmt.Sprintf("{status=\"die\",container_id=\"%s\",image_id=\"%s\",type=\"container\"} %s",
			msg.ID, msg.From, msg.Actor.Attributes["exitCode"]),
		Time:    time.Now().Unix(),
		Accessd: &sync.Map{},
	}
	EventDatas.Store(msg.TimeNano, ev)
	return nil
}

func HandleDockerEvent() {
	filter := filters.NewArgs()
	filter.Add("type", "container")

	msgChan := make(<-chan events.Message, 100)
	errChan := make(<-chan error, 100)
	client := NewDockerClient()
	msgChan, errChan = client.Events(context.Background(), types.EventsOptions{Filters: filter})
	for {
		select {
		case msg := <-msgChan:
			switch msg.Status {
			case "start":
				handleContainerStartMsg(msg)
			case "die":
				handleContainerDieMsg(msg)
			}

		case err := <-errChan:
			if err != nil {
				log.DefFileLogger.WithFields(logrus.Fields{
					"position": utils.GetFileAndLine(),
					"error":    err,
				}).Info("receive EOF error from docker daemon,will reregister")
				client.Close()
				client = NewDockerClient()
				msgChan, errChan = client.Events(context.Background(), types.EventsOptions{Filters: filter})
				time.Sleep(time.Second)
			}
		}
	}
}

func init() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			cleanEventData()
		}
	}()
}
