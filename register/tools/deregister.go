package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/wlibo666/dockerprocmetrics/register/consul"
)

var (
	consulAddr = flag.String("consul", "http://172.16.13.129:8500,http://172.16.13.130:8500,http://172.16.13.131:8500", "-consul http://xxxx:8500")
	dc         = flag.String("dc", "dc1", "-dc dc1")
	service    = flag.String("service", "", "-service service1,service2")
)

func main() {
	flag.Parse()
	if *service == "" {
		fmt.Printf("lost service name\n")
		return
	}
	for _, ser := range strings.Split(*service, ",") {
		err := consul.DeRegisterService(*consulAddr, *dc, ser)
		if err != nil {
			fmt.Printf("DeRegisterService for %s failed,err:%s\n", ser, err.Error())
		} else {
			fmt.Printf("DeRegisterService for %s success\n", ser)
		}
	}
}
