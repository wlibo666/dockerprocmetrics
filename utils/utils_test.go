package utils

import (
	"testing"
)

func TestGetWanAddr(t *testing.T) {
	addr, err := GetWanAddr("lo")
	if err != nil {
		t.Fatalf("get wan addr failed,err:%s", err.Error())
	}
	t.Logf("addr:%s", addr)
	if addr != "127.0.0.1" {
		t.Fatalf("lo addr is not 127.0.0.1")
	}
}

func TestAddNanoSufix(t *testing.T) {
	tmp := AddNanoSufix("key")
	t.Logf("tmp key:%s", tmp)
}
