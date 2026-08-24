package handler_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api"
	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/database"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testServer *api.Server
var testEmailSender service.EmailSender

func setupAuthTestServer(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan_handler_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Skipf("test db not available: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	_ = db.AutoMigrate(
		&model.User{}, &model.RefreshToken{}, &model.Package{}, &model.Transaction{},
		&model.Subscription{}, &model.Event{}, &model.EventSection{}, &model.Guest{},
		&model.RSVP{}, &model.GuestbookMessage{}, &model.DigitalGift{},
		&model.GalleryPhoto{}, &model.AuditLog{}, &model.AnalyticsEvent{},
		&model.WebhookIdempotency{}, &model.Template{}, &model.Music{},
		&model.LoveStory{},
	)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:             "test-secret-key",
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
		LoveStoryRepo:          repository.NewLoveStoryRepository(db),
		DigitalGiftRepo:        repository.NewDigitalGiftRepository(db),
		GalleryPhotoRepo:       repository.NewGalleryPhotoRepository(db),
		AnalyticsRepo:          repository.NewAnalyticsEventRepository(db),
		AuditLogRepo:           repository.NewAuditLogRepository(db),
		WebhookIdempotencyRepo: repository.NewWebhookIdempotencyRepository(db),
		JWTService:             jwtSvc,
		EmailSender:            testEmailSender,
	}

	if err := database.SeedPackages(db); err != nil {
		t.Fatalf("seed packages: %v", err)
	}

	testServer = api.NewServer(cfg, deps, db, log)
	return db
}

func doAuthRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
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
	testServer.Handler().(*chi.Mux).ServeHTTP(w, req)
	return w
}

func generateExpiredTestToken(t *testing.T) string {
	t.Helper()
	secret := []byte("test-secret-key")
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	now := time.Now()
	exp := now.Add(-1 * time.Hour).Unix()
	iat := now.Add(-2 * time.Hour).Unix()
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":"00000000-0000-0000-0000-000000000000","role":"user","exp":%d,"iat":%d}`, exp, iat)))
	signingInput := header + "." + payload
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + signature
}

func TestAuth_Register(t *testing.T) {
	setupAuthTestServer(t)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Auth Test User",
		"email":    "auth_test@example.com",
		"password": "securepassword123",
		"phone":    "6281234567890",
	}, "")
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.True(t, resp["success"].(bool))

	data := resp["data"].(map[string]interface{})
	require.NotEmpty(t, data["id"])
	require.Equal(t, "auth_test@example.com", data["email"].(string))
	require.NotEmpty(t, data["access_token"].(string))
	require.NotEmpty(t, data["refresh_token"].(string))
	require.Greater(t, data["expires_in"].(float64), float64(0))

	bodyStr := w.Body.String()
	require.NotContains(t, bodyStr, "password_hash", "response should not contain password hash")
	require.NotContains(t, bodyStr, `"password"`, "response should not contain password")
}

func TestAuth_Register_DuplicateEmail(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "User One",
		"email":    "dup_test@example.com",
		"password": "securepassword123",
	}, "")

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "User Two",
		"email":    "dup_test@example.com",
		"password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestAuth_Register_ValidationError(t *testing.T) {
	setupAuthTestServer(t)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "",
		"email":    "invalid-email",
		"password": "short",
	}, "")
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuth_Login(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Login Test",
		"email":    "login_test@example.com",
		"password": "mypassword123",
	}, "")

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "login_test@example.com",
		"password": "mypassword123",
	}, "")
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	require.NotEmpty(t, data["access_token"].(string))
	require.NotEmpty(t, data["refresh_token"].(string))

	bodyStr := w.Body.String()
	require.NotContains(t, bodyStr, "password_hash")
	require.NotContains(t, bodyStr, `"password"`)
}

func TestAuth_Login_InvalidCredentials(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Invalid Login",
		"email":    "invalid_login@example.com",
		"password": "correctpassword123",
	}, "")

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "invalid_login@example.com",
		"password": "wrongpassword123",
	}, "")
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.False(t, resp["success"].(bool))

	err := resp["error"].(map[string]interface{})
	require.Contains(t, strings.ToLower(err["message"].(string)), "authentication")
}

func TestAuth_Login_NonExistentUser(t *testing.T) {
	setupAuthTestServer(t)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "nonexistent_user@example.com",
		"password": "anypassword123",
	}, "")
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	err := resp["error"].(map[string]interface{})
	require.Contains(t, strings.ToLower(err["message"].(string)), "authentication")
}

func TestAuth_ExpiredToken(t *testing.T) {
	setupAuthTestServer(t)

	expiredToken := generateExpiredTestToken(t)
	w := doAuthRequest(t, http.MethodGet, "/api/v1/users/me", nil, expiredToken)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_UnauthorizedAccess(t *testing.T) {
	setupAuthTestServer(t)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/users/me", nil, "")
	require.Equal(t, http.StatusUnauthorized, w.Code)

	w2 := doAuthRequest(t, http.MethodGet, "/api/v1/users/me", nil, "Bearer invalidtoken123")
	require.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestAuth_ResourceOwnership(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "User A",
		"email":    "user_a_ownership@example.com",
		"password": "securepassword123",
	}, "")

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "User B",
		"email":    "user_b_ownership@example.com",
		"password": "securepassword123",
	}, "")

	respA := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "user_a_ownership@example.com",
		"password": "securepassword123",
	}, "")
	var loginA map[string]interface{}
	json.Unmarshal(respA.Body.Bytes(), &loginA)
	require.True(t, loginA["success"].(bool))
	tokenA := loginA["data"].(map[string]interface{})["access_token"].(string)

	respB := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "user_b_ownership@example.com",
		"password": "securepassword123",
	}, "")
	var loginB map[string]interface{}
	json.Unmarshal(respB.Body.Bytes(), &loginB)
	require.True(t, loginB["success"].(bool))
	tokenB := loginB["data"].(map[string]interface{})["access_token"].(string)

	profileA := doAuthRequest(t, http.MethodGet, "/api/v1/users/me", nil, tokenA)
	require.Equal(t, http.StatusOK, profileA.Code)
	var profileAData map[string]interface{}
	json.Unmarshal(profileA.Body.Bytes(), &profileAData)
	userIDA := profileAData["data"].(map[string]interface{})["id"].(string)

	profileB := doAuthRequest(t, http.MethodGet, "/api/v1/users/me", nil, tokenB)
	require.Equal(t, http.StatusOK, profileB.Code)
	var profileBData map[string]interface{}
	json.Unmarshal(profileB.Body.Bytes(), &profileBData)
	userIDB := profileBData["data"].(map[string]interface{})["id"].(string)

	require.NotEqual(t, userIDA, userIDB, "User A and User B should have different IDs")
	require.Equal(t, "user_a_ownership@example.com", profileAData["data"].(map[string]interface{})["email"].(string))
	require.Equal(t, "user_b_ownership@example.com", profileBData["data"].(map[string]interface{})["email"].(string))
}

func TestAuth_AdminEndpoint_Protected(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Regular User",
		"email":    "regular_user@example.com",
		"password": "securepassword123",
	}, "")

	loginResp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "regular_user@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodGet, "/api/v1/packages/all", nil, token)
	require.Equal(t, http.StatusForbidden, w.Code, "regular user should not access admin endpoint")

	w2 := doAuthRequest(t, http.MethodGet, "/api/v1/packages/all", nil, "")
	require.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestAuth_Logout(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Logout User",
		"email":    "logout_user@example.com",
		"password": "securepassword123",
	}, "")

	resp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "logout_user@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &loginData)
	refreshToken := loginData["data"].(map[string]interface{})["refresh_token"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/logout", map[string]string{
		"refresh_token": refreshToken,
	}, "")
	require.Equal(t, http.StatusOK, w.Code)

	w2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/refresh", map[string]string{
		"refresh_token": refreshToken,
	}, "")
	require.Equal(t, http.StatusUnauthorized, w2.Code, "revoked refresh token should not work")
}

func TestAuth_ChangePassword(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Change PW User",
		"email":    "change_pw@example.com",
		"password": "oldpassword123",
	}, "")

	resp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "change_pw@example.com",
		"password": "oldpassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &loginData)
	require.True(t, loginData["success"].(bool))
	token := loginData["data"].(map[string]interface{})["access_token"].(string)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/change-password", map[string]string{
		"current_password": "wrongpassword",
		"new_password":     "newpassword123",
	}, token)
	require.Equal(t, http.StatusForbidden, w.Code)

	w2 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/change-password", map[string]string{
		"current_password": "oldpassword123",
		"new_password":     "newpassword123",
	}, token)
	require.Equal(t, http.StatusOK, w2.Code)

	w3 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "change_pw@example.com",
		"password": "newpassword123",
	}, "")
	require.Equal(t, http.StatusOK, w3.Code)

	w4 := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "change_pw@example.com",
		"password": "oldpassword123",
	}, "")
	require.Equal(t, http.StatusUnauthorized, w4.Code, "old password should not work after change")
}

func TestAuth_JWT_TokenStructure(t *testing.T) {
	setupAuthTestServer(t)

	doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "JWT Test",
		"email":    "jwt_test@example.com",
		"password": "securepassword123",
	}, "")

	resp := doAuthRequest(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "jwt_test@example.com",
		"password": "securepassword123",
	}, "")
	var loginData map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &loginData)

	data := loginData["data"].(map[string]interface{})
	token := data["access_token"].(string)
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3, "JWT should have 3 parts (header.payload.signature)")

	expiresIn := data["expires_in"].(float64)
	require.LessOrEqual(t, expiresIn, float64(900), "access token expiry should be <= 15 minutes (900 seconds)")
}

func TestAuth_Register_SendsWelcomeEmail(t *testing.T) {
	rec := &recEmailSender{}
	testEmailSender = rec
	defer func() { testEmailSender = nil }()
	setupAuthTestServer(t)

	w := doAuthRequest(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Welcome User",
		"email":    "welcome@example.com",
		"password": "securepassword123",
	}, "")
	require.Equal(t, http.StatusCreated, w.Code)

	require.Equal(t, 1, rec.count(), "registration must enqueue exactly one welcome email")
	rec.mu.Lock()
	defer rec.mu.Unlock()
	require.Equal(t, "welcome@example.com", rec.emails[0].To)
	require.Contains(t, rec.emails[0].HTML, "Welcome User")
}
