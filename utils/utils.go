package utils

import (
	"os"
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
