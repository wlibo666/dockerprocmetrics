package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wlibo666/dockerprocmetrics/metrics"
)

func Metrics(c *gin.Context) {
	args := make(map[string]string)
	args[metrics.CONTAINER_STATE_ARG_REMOTEADDR] = strings.Split(c.Request.RemoteAddr, ":")[0]
	metrics.WriteAllMetricsData(c.Writer, args)
}
