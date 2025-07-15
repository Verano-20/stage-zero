package main

import (
	"github.com/Verano-20/go-crud/internal/initializer"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/router"
	"go.uber.org/zap"

	_ "github.com/Verano-20/go-crud/docs"
)

const (
	port = "8080"
)

// @title           Go-CRUD API
// @version         1.0
// @description     This is a CRUD API for a simple application.

// @host      localhost:8080
// @BasePath  /
func main() {
	logger.Init()
	defer logger.Sync()
	log := logger.Get()

	log.Info("Starting Go-CRUD API", zap.String("version", "1.0"), zap.String("port", port))

	log.Info("Initializing database connection...")
	db := initializer.InitializeDatabase()
	log.Info("Database connection established")

	log.Info("Initializing routes...")
	router := router.InitializeRouter(db)
	log.Info("Routes configured")

	log.Info("Server starting", zap.String("address", ":"+port))
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
