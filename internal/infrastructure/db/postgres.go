package db

import (
	"fmt"
	"subcalc/internal/config"
	"subcalc/internal/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg *config.Config, logger *zap.SugaredLogger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	logger.Infof("connected to postgres: %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.Subscription{})
}
