package utils

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	EXIT_WAIT_DURATION = 3
)

func ExitWait(duration, code int) {
	time.Sleep(time.Duration(duration) * time.Second)
	os.Exit(code)
}

func ExitWaitDef(code int) {
	time.Sleep(time.Duration(EXIT_WAIT_DURATION) * time.Second)
	os.Exit(code)
}

func GetWanAddr(name string) (string, error) {
	dev, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addrs, err := dev.Addrs()
	if err != nil {
		return "", err
	}
	return strings.Split(addrs[0].String(), "/")[0], nil
}

func GenListenAddr(addr string, port int) string {
	if addr == "" {
		return fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("%s:%d", addr, port)
}

func AddNanoSufix(key string) string {
	return fmt.Sprintf("%s-%x", key, time.Now().UnixNano())
}

func NewContextWithTimeout(timeout int) context.Context {
	timeoutContext, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return timeoutContext
}

func GetFileAndLine() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", file, line)
}
