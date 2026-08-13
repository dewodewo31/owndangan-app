package admin_test

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
	"github.com/owndangan/backend/internal/repository"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testAdminServer *api.Server

func setupAdminTestServer(t *testing.T) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan_admin_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Skipf("test db not available: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")

	db.AutoMigrate(
		&model.User{}, &model.RefreshToken{}, &model.Package{}, &model.Transaction{},
		&model.Subscription{}, &model.Template{}, &model.Music{}, &model.Event{},
		&model.EventSection{}, &model.Guest{}, &model.RSVP{}, &model.GuestbookMessage{},
		&model.DigitalGift{}, &model.GalleryPhoto{}, &model.AuditLog{},
		&model.AnalyticsEvent{}, &model.WebhookIdempotency{},
	)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:             "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:3000"}},
	}

	log := zerolog.Nop()
	jwtSvc := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)

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
	}

	database.SeedPackages(db)
	testAdminServer = api.NewServer(cfg, deps, db, log)
}

func doAdminRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
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
	testAdminServer.Handler().(*chi.Mux).ServeHTTP(w, req)
	return w
}

func getAdminToken(t *testing.T) string {
	doAdminRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name": "Admin", "email": "admin@test.com", "password": "password123",
	}, "")
	
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan_admin_test sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	sqlDB, _ := db.DB()
	result, _ := sqlDB.Exec("UPDATE users SET role = 'admin' WHERE email = 'admin@test.com'")
	rows, _ := result.RowsAffected()
	t.Logf("Updated %d rows", rows)
	
	var role string
	sqlDB.QueryRow("SELECT role FROM users WHERE email = 'admin@test.com'").Scan(&role)
	t.Logf("User role after update: %s", role)
	sqlDB.Close()
	
	loginResp := doAdminRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "admin@test.com", "password": "password123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	t.Logf("Login role: %s", loginData["data"].(map[string]interface{})["role"])
	return loginData["data"].(map[string]interface{})["access_token"].(string)
}

func TestAdmin_Analytics(t *testing.T) {
	setupAdminTestServer(t)
	token := getAdminToken(t)
	t.Logf("Token: %s", token[:20]+"...")

	w := doAdminRequest(t, http.MethodGet, "/api/v1/admin/analytics", nil, token)
	t.Logf("Analytics status: %d", w.Code)
	t.Logf("Analytics body: %s", w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.True(t, resp["success"].(bool))
}

func TestAdmin_ListUsers(t *testing.T) {
	setupAdminTestServer(t)
	token := getAdminToken(t)

	w := doAdminRequest(t, http.MethodGet, "/api/v1/admin/users", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAdmin_ListPackages(t *testing.T) {
	setupAdminTestServer(t)
	token := getAdminToken(t)

	w := doAdminRequest(t, http.MethodGet, "/api/v1/admin/packages", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAdmin_ListTemplates(t *testing.T) {
	setupAdminTestServer(t)
	token := getAdminToken(t)

	w := doAdminRequest(t, http.MethodGet, "/api/v1/admin/templates", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAdmin_ListTransactions(t *testing.T) {
	setupAdminTestServer(t)
	token := getAdminToken(t)

	w := doAdminRequest(t, http.MethodGet, "/api/v1/admin/transactions", nil, token)
	require.Equal(t, http.StatusOK, w.Code)
}

