package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	err := LoadConfig("../conf/metrics.json")
	if err != nil {
		t.Fatalf("load config failed,err:%s", err.Error())
	}
	t.Logf("config:%v\n", GMertricConfig)
}
