package middleware

import (
	"strconv"
	"time"

	"github.com/Verano-20/go-crud/internal/telemetry"
	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := telemetry.GetMetrics()
		if metrics == nil {
			c.Next()
			return
		}

		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		ctx := c.Request.Context()
		metrics.RecordActiveHTTPRequest(ctx, 1)

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		metrics.RecordHTTPRequest(ctx, method, path, status, duration)
		metrics.RecordActiveHTTPRequest(ctx, -1)
	}
}
