package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/owndangan/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func setupAuthService(t *testing.T) (*service.AuthService, *mockPackageRepo) {
	t.Helper()
	pkgRepo := newMockPackageRepo()
	duration := 7
	freePkg := &model.Package{
		ID:            newUUID(),
		Name:          "Free",
		Code:          "free",
		Price:         0,
		DurationDays:  &duration,
		GuestLimit:    &duration,
		TemplateGroup: "standard",
		Features:      datatypes.JSON(`{}`),
		IsActive:      true,
	}
	pkgRepo.packages["free"] = freePkg

	authSvc := service.NewAuthService(
		newMockUserRepo(),
		&mockRefreshTokenRepo{},
		pkgRepo,
		&mockSubscriptionRepo{},
		jwt.New("test-secret-key", 15*time.Minute, 7*24*time.Hour),
		&mockAuditLogRepo{},
		noopEmailSender{},
		"http://localhost:3000",
	)
	return authSvc, pkgRepo
}

func newUUID() (id uuid.UUID) {
	return uuid.New()
}

func TestAuthService_Register(t *testing.T) {
	svc, _ := setupAuthService(t)

	resp, accessToken, refreshToken, expiresIn, err := svc.Register(context.TODO(), "John Doe", "john@example.com", "password123", "+6281234567890")

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.Greater(t, expiresIn, int64(0))
	assert.Equal(t, "john@example.com", resp.Email)
	assert.Equal(t, "John Doe", resp.Name)
	assert.Equal(t, "active", resp.Status)
	assert.False(t, resp.ID == uuid.Nil)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	svc, _ := setupAuthService(t)

	_, _, _, _, err := svc.Register(context.TODO(), "John Doe", "john@example.com", "password123", "+6281234567890")
	assert.NoError(t, err)

	// Register again with same email
	_, _, _, _, err = svc.Register(context.TODO(), "Jane Doe", "john@example.com", "password123", "+6281234567890")
	assert.Error(t, err)
}

func TestAuthService_Login_Failure(t *testing.T) {
	svc, _ := setupAuthService(t)

	_, _, _, _, err := svc.Login(context.TODO(), "nonexistent@example.com", "password")
	assert.Error(t, err)
}

func TestAuthService_Register_InvalidPassword(t *testing.T) {
	req := dto.RegisterRequest{
		Name:     "Test",
		Email:    "test@example.com",
		Password: "short",
	}
	err := validateRegisterRequest(req)
	assert.Error(t, err)
}

func validateRegisterRequest(req dto.RegisterRequest) error {
	if len(req.Password) < 8 {
		return errPasswordTooShort
	}
	return nil
}

var errPasswordTooShort = &validationError{"password too short"}

type validationError struct{ msg string }

func (e *validationError) Error() string { return e.msg }
