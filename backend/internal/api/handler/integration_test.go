package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api"
	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/database"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/owndangan/backend/internal/pkg/storage"
	"github.com/owndangan/backend/internal/repository"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) (*api.Server, *gorm.DB) {
	t.Helper()
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan_handler_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	if err := db.AutoMigrate(
		&model.User{}, &model.RefreshToken{}, &model.Package{}, &model.Transaction{},
		&model.Subscription{}, &model.Event{}, &model.EventSection{}, &model.Guest{},
		&model.RSVP{}, &model.GuestbookMessage{}, &model.DigitalGift{},
		&model.GalleryPhoto{}, &model.AuditLog{}, &model.AnalyticsEvent{},
		&model.WebhookIdempotency{}, &model.Template{}, &model.Music{},
	); err != nil {
		t.Logf("auto migrate note: %v", err)
	}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:             "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		Storage: config.StorageConfig{
			Provider:  "local",
			LocalPath: t.TempDir(),
			PublicURL: "/uploads",
		},
	}

	log := zerolog.Nop()
	jwtSvc := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
	store := storage.NewLocalStorage(cfg.Storage.LocalPath, cfg.Storage.PublicURL)

	deps := &api.Dependencies{
		UserRepo:               repository.NewUserRepository(db),
		RefreshTokenRepo:       repository.NewRefreshTokenRepository(db),
		PackageRepo:            repository.NewPackageRepository(db),
		TransactionRepo:        repository.NewTransactionRepository(db),
		SubscriptionRepo:       repository.NewSubscriptionRepository(db),
		EventRepo:              repository.NewEventRepository(db),
		EventSectionRepo:       repository.NewEventSectionRepository(db),
		TemplateRepo:           repository.NewTemplateRepository(db),
		MusicRepo:              repository.NewMusicRepository(db),
		GuestRepo:              repository.NewGuestRepository(db),
		RSVPRepo:               repository.NewRSVPRepository(db),
		GuestbookRepo:          repository.NewGuestbookRepository(db),
		DigitalGiftRepo:        repository.NewDigitalGiftRepository(db),
		GalleryPhotoRepo:       repository.NewGalleryPhotoRepository(db),
		AnalyticsRepo:          repository.NewAnalyticsEventRepository(db),
		AuditLogRepo:           repository.NewAuditLogRepository(db),
		WebhookIdempotencyRepo: repository.NewWebhookIdempotencyRepository(db),
		JWTService:             jwtSvc,
		Storage:                store,
	}

	if err := database.SeedPackages(db); err != nil {
		t.Fatalf("seed packages: %v", err)
	}

	server := api.NewServer(cfg, deps, db, log)
	return server, db
}

func doRequest(t *testing.T, s *chi.Mux, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	var req *http.Request
	if reader != nil {
		req = httptest.NewRequest(method, path, reader)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{}))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	return w
}

func TestAPI_RegisterAndLogin(t *testing.T) {
	server, _ := setupTestServer(t)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Integration User",
		"email":    "integration@test.com",
		"password": "password123",
		"phone":    "6281234567890",
	}, "")

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert := require.New(t)
	assert.True(resp["success"].(bool))

	data := resp["data"].(map[string]interface{})
	assert.Equal("integration@test.com", data["email"].(string))
	assert.NotEmpty(data["access_token"].(string))
	assert.NotEmpty(data["refresh_token"].(string))

	// Login
	w2 := doRequest(t, server.Handler().(*chi.Mux), http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "integration@test.com",
		"password": "password123",
	}, "")

	require.Equal(t, http.StatusOK, w2.Code)
}

func TestAPI_PackagesList(t *testing.T) {
	server, _ := setupTestServer(t)

	w := doRequest(t, server.Handler().(*chi.Mux), http.MethodGet, "/api/v1/packages", nil, "")

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert := require.New(t)
	assert.True(resp["success"].(bool))

	data := resp["data"].([]interface{})
	assert.Len(data, 3)
}
