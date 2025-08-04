package main

import (
	"fmt"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/database"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/router"
	"github.com/Verano-20/go-crud/internal/telemetry"
	"go.uber.org/zap"

	_ "github.com/Verano-20/go-crud/docs"
)

func init() {
	config.InitConfig()
	logger.InitLogger()
}

// @title           Go-CRUD API
// @version         1.0
// @description     This is a CRUD API for a simple application.

// @host      localhost:8080
// @BasePath  /
func main() {
	config := config.Get()
	log := logger.Get()
	defer logger.Sync()

	log.Info(fmt.Sprintf("Starting %s", config.ServiceName),
		zap.String("service_name", config.ServiceName),
		zap.String("service_version", config.ServiceVersion),
		zap.String("port", config.ServicePort),
	)

	telemetry.InitTelemetry()
	db := database.InitDatabase()
	router := router.InitRouter(db)

	address := ":" + config.ServicePort
	log.Info("Server starting", zap.String("address", address))
	if err := router.Run(address); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
