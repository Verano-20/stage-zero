package logger

import (
	"github.com/Verano-20/stage-zero/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func InitLogger() {
	var err error

	config := config.Get()
	env := config.Environment

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)

		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "time"
		zapConfig.EncoderConfig.MessageKey = "msg"
		zapConfig.EncoderConfig.LevelKey = "level"
		zapConfig.EncoderConfig.CallerKey = "caller"
		zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

		// Add service info to all logs
		globalLogger, err = zapConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.Fields(
				zap.String("service", config.ServiceName),
				zap.String("version", config.ServiceVersion),
				zap.String("environment", config.Environment),
			),
		)
	} else {
		gin.SetMode(gin.DebugMode)

		zapConfig := zap.NewDevelopmentConfig()
		// Use JSON encoding for better Loki parsing even in dev
		zapConfig.Encoding = "json"
		zapConfig.EncoderConfig.TimeKey = "time"
		zapConfig.EncoderConfig.MessageKey = "msg"
		zapConfig.EncoderConfig.LevelKey = "level"
		zapConfig.EncoderConfig.CallerKey = "caller"
		zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

		globalLogger, err = zapConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.Fields(
				zap.String("service", config.ServiceName),
				zap.String("version", config.ServiceVersion),
				zap.String("environment", config.Environment),
			),
		)
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	zap.ReplaceGlobals(globalLogger)
}

func Get() *zap.Logger {
	if globalLogger == nil {
		panic("Logger not initialized")
	}
	return globalLogger
}

func GetFromContext(ctx *gin.Context) *zap.Logger {
	if logger, exists := ctx.Get("logger"); exists {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			return zapLogger
		}
	}
	return zap.L() // fallback to global logger
}

// Sync flushes any buffered log entries
func Sync() {
	if globalLogger != nil {
		globalLogger.Sync()
	}
}
