package container

import (
	"encoding/json"
	"testing"
)

func TestInitHostLabels(t *testing.T) {
	hostinfo := struct {
		HostName string `json:"hostname"`
		CpuCnt   string `json:"cpucnt"`
		Memory   string `json:"memory"`
	}{
		HostName: "k8s-129-master",
		CpuCnt:   "2",
		Memory:   "9024520",
	}

	content, _ := json.Marshal(hostinfo)
	t.Logf("content:%s\n", string(content))
	t.Logf("str:%s\n", string(content)[1:len(content)-1])
}
