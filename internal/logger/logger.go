package logger

import (
	"github.com/Verano-20/go-crud/internal/config"
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
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.MessageKey = "message"
		zapConfig.EncoderConfig.LevelKey = "level"
		zapConfig.EncoderConfig.CallerKey = "caller"

		globalLogger, err = zapConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	} else {
		gin.SetMode(gin.DebugMode)

		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		globalLogger, err = zapConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
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
