package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func PrometheusMiddleware() gin.HandlerFunc {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

		c.Next()

		statusCode := c.Writer.Status()
		totalRequests.WithLabelValues(path).Inc()
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		timer.ObserveDuration()
	}
}
