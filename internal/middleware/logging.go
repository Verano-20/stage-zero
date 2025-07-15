package middleware

import (
	"time"

	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()

		log := logger.Get().With(
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		)

		ctx.Set("logger", log)

		log.Info("Incoming HTTP request")

		ctx.Next()

		latency := time.Since(start)
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
			zap.Duration("latency", latency),
		)
	}
}
