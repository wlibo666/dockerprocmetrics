package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wlibo666/dockerprocmetrics/metrics"
)

const (
	CTL_URL_MTERICS = "/metrics"
	CTL_URL_HEALTHZ = "/healthz"
)

func Metrics(c *gin.Context) {
	args := make(map[string]string)
	args[metrics.CONTAINER_STATE_ARG_REMOTEADDR] = strings.Split(c.Request.RemoteAddr, ":")[0]

	metrics.WriteAllMetricsData(c.Writer, args)
}

func Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
