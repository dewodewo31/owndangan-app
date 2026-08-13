package database

import (
	"fmt"
	"time"

	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/model"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func Connect(cfg config.DatabaseConfig, log zerolog.Logger) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	var gormLogLevel logger.LogLevel = logger.Silent
	if cfg.LogSQL {
		gormLogLevel = logger.Info
	}

	gormCfg := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(gormLogLevel),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Info().Msg("database connected successfully")
	return &DB{db}, nil
}

func (d *DB) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Package{},
		&model.Transaction{},
		&model.Subscription{},
		&model.Template{},
		&model.Music{},
		&model.Event{},
		&model.EventSection{},
		&model.Guest{},
		&model.RSVP{},
		&model.GuestbookMessage{},
		&model.DigitalGift{},
		&model.GalleryPhoto{},
		&model.AuditLog{},
		&model.AnalyticsEvent{},
		&model.WebhookIdempotency{},
	)
}

func Close(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func Migrate(cfg config.DatabaseConfig, log zerolog.Logger) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	gormCfg := &gorm.Config{
		SkipDefaultTransaction: true,
	}
	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Package{},
		&model.Transaction{},
		&model.Subscription{},
		&model.Template{},
		&model.Music{},
		&model.Event{},
		&model.EventSection{},
		&model.Guest{},
		&model.RSVP{},
		&model.GuestbookMessage{},
		&model.DigitalGift{},
		&model.GalleryPhoto{},
		&model.AuditLog{},
		&model.AnalyticsEvent{},
		&model.WebhookIdempotency{},
	); err != nil {
		return err
	}

	log.Info().Msg("database migration completed")
	return nil
}
