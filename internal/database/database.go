package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/logger"
	"github.com/Verano-20/go-crud/internal/telemetry"
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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying database connection", zap.Error(err))
	}

	maxOpenConns := 25
	maxIdleConns := 5
	connMaxLifetime := 5 * time.Minute
	connMaxIdleTime := 1 * time.Minute

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	log.Info("Database connection established with connection pool configured",
		zap.Int("max_open_conns", maxOpenConns),
		zap.Int("max_idle_conns", maxIdleConns),
		zap.Duration("conn_max_lifetime", connMaxLifetime),
		zap.Duration("conn_max_idle_time", connMaxIdleTime))

	// Start connection metrics collection
	go startConnectionMetricsCollector(sqlDB)

	return db
}

var shutdownMetricsCollector = make(chan struct{})

func startConnectionMetricsCollector(sqlDB *sql.DB) {
	log := logger.Get()
	log.Info("Starting database connection metrics collector")

	// Wait for telemetry to be fully initialized
	time.Sleep(2 * time.Second)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-shutdownMetricsCollector:
			log.Info("Stopping database connection metrics collector")
			return
		case <-ticker.C:
			collectConnectionMetrics(sqlDB)
		}
	}
}

func collectConnectionMetrics(sqlDB *sql.DB) {
	log := logger.Get()

	defer func() {
		if r := recover(); r != nil {
			log.Warn("Failed to collect database connection metrics - telemetry may not be ready",
				zap.Any("error", r))
			return
		}
	}()

	ctx := context.Background()
	metrics := telemetry.GetMetrics()
	stats := sqlDB.Stats()

	metrics.RecordDBConnectionStats(ctx, stats)

	log.Debug("Database connection metrics collected")
}

func Shutdown(db *gorm.DB) {
	log := logger.Get()
	log.Info("Closing database connection...")

	close(shutdownMetricsCollector)

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Error("Failed to close database connection", zap.Error(err))
		}
	} else {
		log.Error("Failed to get underlying database connection", zap.Error(err))
	}
}
