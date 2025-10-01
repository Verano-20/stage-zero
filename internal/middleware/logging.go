package middleware

import (
	"time"

	"github.com/Verano-20/stage-zero/internal/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()

		logFields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		}

		span := trace.SpanFromContext(ctx.Request.Context())
		spanContext := span.SpanContext()

		if spanContext.IsValid() {
			logFields = append(logFields,
				zap.String("trace_id", spanContext.TraceID().String()),
				zap.String("span_id", spanContext.SpanID().String()),
			)
		}

		log := logger.Get().With(logFields...)

		ctx.Set("logger", log)

		log.Info("Incoming HTTP request")

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()

		logLevel := zap.InfoLevel
		message := "Request completed"

		if status >= 400 && status < 500 {
			logLevel = zap.WarnLevel
			message = "Request completed with client error"
		} else if status >= 500 {
			logLevel = zap.ErrorLevel
			message = "Request completed with server error"
		}

		log.Log(logLevel, message,
			zap.Int("status", status),
			zap.Duration("duration", duration),
		)
	}
}
