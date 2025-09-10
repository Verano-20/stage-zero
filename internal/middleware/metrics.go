package middleware

import (
	"strconv"
	"time"

	"github.com/Verano-20/go-crud/internal/telemetry"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

		attrs := []attribute.KeyValue{
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.String("status", status),
		}

		metrics.HTTPRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
		metrics.HTTPRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))

		metrics.RecordActiveHTTPRequest(ctx, -1)
	}
}
