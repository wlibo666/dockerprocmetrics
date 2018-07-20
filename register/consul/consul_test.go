package consul

import (
	"testing"
)

func TestRegisterService(t *testing.T) {
	err := RegisterService("http://172.16.13.129:8500", "dc1",
		"dpm-001", "dpm", "172.16.13.129", 3280, "/healthz")
	if err != nil {
		t.Logf("RegisterService failed,err:%s", err.Error())
	}
}

func TestDeRegisterService(t *testing.T) {
	err := DeRegisterService("http://172.16.13.129:8500", "dc1", "dpm-172.16.13.129-1542c22f9f535088")
	if err != nil {
		t.Logf("DeRegisterService failed,err:%s", err.Error())
	}
	err = DeRegisterService("http://172.16.13.129:8500", "dc1", "dpm-001")
	if err != nil {
		t.Logf("DeRegisterService failed,err:%s", err.Error())
	}
}

func TestGetServiceEndpoints(t *testing.T) {
	addrs, err := GetServiceEndpoints("http://172.16.13.129:8500", "dc1", "dpm", true)
	if err != nil {
		t.Logf("GetServiceEndpoints err:%s", err.Error())
	}

	t.Logf("addrs :%v\n", addrs)
}
