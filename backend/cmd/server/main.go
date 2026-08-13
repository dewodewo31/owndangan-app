package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/owndangan/backend/internal/api"
	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/database"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/owndangan/backend/internal/pkg/logger"
	"github.com/owndangan/backend/internal/pkg/storage"
	"github.com/owndangan/backend/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)

	db, err := database.Connect(cfg.Database, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatal().Err(err).Msg("failed to run auto-migration")
	}

	if err := database.SeedPackages(db.DB); err != nil {
		log.Error().Err(err).Msg("failed to seed packages")
	}

	log.Info().Msg("database ready")

	jwtSvc := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)

	deps := &api.Dependencies{
		UserRepo:               repository.NewUserRepository(db.DB),
		RefreshTokenRepo:       repository.NewRefreshTokenRepository(db.DB),
		PackageRepo:            repository.NewPackageRepository(db.DB),
		TransactionRepo:        repository.NewTransactionRepository(db.DB),
		SubscriptionRepo:       repository.NewSubscriptionRepository(db.DB),
		EventRepo:              repository.NewEventRepository(db.DB),
		EventSectionRepo:       repository.NewEventSectionRepository(db.DB),
		TemplateRepo:           repository.NewTemplateRepository(db.DB),
		MusicRepo:              repository.NewMusicRepository(db.DB),
		GuestRepo:              repository.NewGuestRepository(db.DB),
		RSVPRepo:               repository.NewRSVPRepository(db.DB),
		GuestbookRepo:          repository.NewGuestbookRepository(db.DB),
		DigitalGiftRepo:        repository.NewDigitalGiftRepository(db.DB),
		GalleryPhotoRepo:       repository.NewGalleryPhotoRepository(db.DB),
		AnalyticsRepo:          repository.NewAnalyticsEventRepository(db.DB),
		AuditLogRepo:           repository.NewAuditLogRepository(db.DB),
		WebhookIdempotencyRepo: repository.NewWebhookIdempotencyRepository(db.DB),
		JWTService:             jwtSvc,
		Storage:                newStorage(cfg),
	}

	server := api.NewServer(cfg, deps, db.DB, log)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
	database.Close(db.DB)
	log.Info().Msg("server stopped")
}

func newStorage(cfg *config.Config) storage.Storage {
	switch cfg.Storage.Provider {
	case "s3":
		return storage.NewS3Storage(cfg.Storage.Bucket, cfg.Storage.Region, cfg.Storage.AccessKey,
			cfg.Storage.SecretKey, cfg.Storage.Endpoint, cfg.Storage.PublicURL)
	default:
		return storage.NewLocalStorage(cfg.Storage.LocalPath, "/uploads")
	}
}
