package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Verano-20/stage-zero/internal/config"
	"github.com/Verano-20/stage-zero/internal/container"
	"github.com/Verano-20/stage-zero/internal/database"
	"github.com/Verano-20/stage-zero/internal/logger"
	"github.com/Verano-20/stage-zero/internal/router"
	"github.com/Verano-20/stage-zero/internal/telemetry"
	"go.uber.org/zap"

	_ "github.com/Verano-20/stage-zero/docs"
)

func init() {
	config.InitConfig()
	logger.InitLogger()
}

// @title           Stage-Zero API
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
	defer telemetry.Shutdown()

	db := database.InitDatabase()
	defer database.Shutdown(db)

	container := container.NewContainerWithDB(db)
	ginRouter := router.InitRouter(container)

	server := &http.Server{
		Addr:    ":" + config.ServicePort,
		Handler: ginRouter,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Server starting", zap.String("address", server.Addr))
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", zap.Error(err))
		}

	case sig := <-shutdown:
		log.Info("Shutdown signal received", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Info("Starting graceful shutdown...")

		if err := server.Shutdown(ctx); err != nil {
			log.Error("Server shutdown failed", zap.Error(err))
			if err := server.Close(); err != nil {
				log.Error("Server force close failed", zap.Error(err))
			}
		} else {
			log.Info("HTTP server shutdown completed")
		}
	}

	log.Info("Application shutdown completed")
}
