package logger

import (
	"github.com/Verano-20/go-crud/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init() {
	var err error

	env := config.GetEnvironment()

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)

		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
		config.EncoderConfig.CallerKey = "caller"

		Logger, err = config.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	} else {
		gin.SetMode(gin.DebugMode)

		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		Logger, err = config.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	zap.ReplaceGlobals(Logger)
}

// Get returns the global logger instance
func Get() *zap.Logger {
	return Logger
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
	if Logger != nil {
		Logger.Sync()
	}
}
