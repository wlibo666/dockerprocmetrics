package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wlibo666/dockerprocmetrics/metrics"
)

func Metrics(c *gin.Context) {
	args := make(map[string]string)
	args[metrics.CONTAINER_STATE_ARG_REMOTEADDR] = strings.Split(c.Request.RemoteAddr, ":")[0]
	c.String(http.StatusOK, metrics.GetAllMetricsData(args))
}
