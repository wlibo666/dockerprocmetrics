package container

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
)

func TestStorePidsByContainerId(t *testing.T) {
	client, err := client.NewClientWithOpts(client.WithVersion("1.37"), client.WithHost("unix:///var/run/docker.sock"))
	if err != nil {
		t.Fatalf("new client failed,err:%s", err.Error())
	}
	body, err := client.ContainerTop(context.Background(), "c5945343c01783451fda72e12335b911830c5b9beb192c752a8025b9388ae218", []string{})
	if err != nil {
		t.Fatalf("top failed,err:%s", err.Error())
	}
	t.Logf("body:%v\n", body)
	pidnum := len(body.Processes)
	t.Logf("pidnum:%d\n", pidnum)
	for index, pids := range body.Processes {
		t.Logf("index:%d,pids:%v, pid:%s\n", index, pids, pids[1])
	}
}
