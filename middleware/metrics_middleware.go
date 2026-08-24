package middleware

import (
	"strconv"
	"time"
	"url-shortener/metrics"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		metrics.HTTPRequestsInFlight.Inc()

		defer metrics.HTTPRequestsInFlight.Dec()

		c.Next()

		duration := time.Since(start).Seconds()

		status := c.Writer.Status()

		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(status),
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
