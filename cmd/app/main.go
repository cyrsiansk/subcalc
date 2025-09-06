package main

// @title Subscriptions API
// @version 1.0
// @description Simple service to track user subscriptions (monthly prices).

import (
	"log"
	"subcalc/internal/config"
	"subcalc/internal/infrastructure/db"
	loggerpkg "subcalc/internal/logger"

	"go.uber.org/zap"
)

func main() {
	// CFG
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// LOGGER

	rawLog, err := loggerpkg.New(cfg.LogLevel, false)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func() {
		_ = rawLog.Sync()
	}()

	sugar := rawLog.Sugar()
	sugar.Infof("starting subscriptions service on port %s", cfg.AppPort)

	// DB

	database, err := db.NewPostgres(cfg, sugar)
	if err != nil {
		sugar.Fatalf("failed to connect to db: %v", err)
	}

	if cfg.AutoMigrate {
		if err := db.AutoMigrate(database); err != nil {
			sugar.Fatalf("auto-migrate failed: %v", err)
		}
		sugar.Info("auto-migrate finished")
	} else {
		sugar.Info("auto-migrate disabled")
	}

	// APP

	sugar.Info("config loaded", zap.Any("config", cfg))
	sugar.Info("service started")
}
