package db

import (
	"headless-cms/internal/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	var level logger.LogLevel
	switch cfg.Env {
	case "prod":
		level = logger.Warn
	default:
		level = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             750 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
		},
	)

	return gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
}

