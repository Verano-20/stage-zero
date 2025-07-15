package initializer

import (
	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	log := logger.Get()

	dsn := config.GetDBConnectionString()

	log.Info("Connecting to database",
		zap.String("host", config.NewConfig().DBHost),
		zap.String("database", config.NewConfig().DBName),
		zap.String("port", config.NewConfig().DBPort))

	// Initialize database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database",
			zap.Error(err),
			zap.String("dsn", dsn))
	}

	return db
}
