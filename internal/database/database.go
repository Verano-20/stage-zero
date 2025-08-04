package database

import (
	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	config := config.Get()
	log := logger.Get()

	log.Info("Connecting to database...",
		zap.String("host", config.Database.Host),
		zap.String("database", config.Database.Name),
		zap.String("port", config.Database.Port))

	dsn := config.GetDBConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database",
			zap.Error(err),
			zap.String("dsn", dsn))
	}

	log.Info("Database connection established")
	return db
}
