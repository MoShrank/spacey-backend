package middleware

// copied and adapted from this repository:
//	https://github.com/toorop/gin-logrus

import (
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/logger"
)

var timeFormat = "02/Jan/2006:15:04:05 -0700"

func Logger(logger logger.LoggerInterface, notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		fields := map[string]interface{}{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency,
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
			"time":       time.Now().Format(timeFormat),
		}

		if len(c.Errors) > 0 {
			logger.Error(fields)
		} else {
			if statusCode >= http.StatusInternalServerError {
				logger.Error(fields)
			} else if statusCode >= http.StatusBadRequest {
				logger.Warn(fields)
			} else {
				logger.Info(fields)
			}
		}
	}
}
