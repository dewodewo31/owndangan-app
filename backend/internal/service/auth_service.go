package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/email"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepository
	refreshRepo repository.RefreshTokenRepository
	pkgRepo     repository.PackageRepository
	subRepo     repository.SubscriptionRepository
	jwtService  *jwt.Service
	auditRepo   repository.AuditLogRepository
	emailSvc    EmailSender
	loginURL    string
}

func NewAuthService(userRepo repository.UserRepository, refreshRepo repository.RefreshTokenRepository, pkgRepo repository.PackageRepository, subRepo repository.SubscriptionRepository, jwtSvc *jwt.Service, auditRepo repository.AuditLogRepository, emailSvc EmailSender, loginURL string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		pkgRepo:     pkgRepo,
		subRepo:     subRepo,
		jwtService:  jwtSvc,
		auditRepo:   auditRepo,
		emailSvc:    emailSvc,
		loginURL:    loginURL,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password, phone string) (*model.User, string, string, int64, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", "", 0, fmt.Errorf("%w: email already registered", errors.ErrConflict)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Phone:        phone,
		Role:         "user",
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", "", 0, fmt.Errorf("create user: %w", err)
	}

	s.sendWelcomeEmail(ctx, user)

	freePkg, err := s.pkgRepo.GetByCode(ctx, "free")
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("get free package: %w", err)
	}

	start := time.Now()
	duration := 7 * 24 * time.Hour
	expiresAt := start.Add(duration)
	freeSub := &model.Subscription{
		UserID:    user.ID,
		PackageID: freePkg.ID,
		Status:    "active",
		StartAt:   start,
		ExpiresAt: &expiresAt,
	}
	if err := s.subRepo.Create(ctx, freeSub); err != nil {
		_ = s.userRepo.SoftDelete(ctx, user.ID)
		return nil, "", "", 0, fmt.Errorf("create free subscription: %w", err)
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, _, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("generate refresh token: %w", err)
	}

	tokenHash := hashToken(refreshToken)
	tokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: tokenExpiry,
	}
	if err := s.refreshRepo.Create(ctx, rt); err != nil {
		return nil, "", "", 0, fmt.Errorf("store refresh token: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &user.ID,
		Action:     "user.registered",
		EntityType: "user",
		EntityID:   &user.ID,
		Metadata:   datatypesJSON(map[string]interface{}{"email": email}),
	})

	return user, accessToken, refreshToken, expiresIn, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, string, int64, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, "", "", 0, errors.ErrUnauthorized
	}

	if user.Status == "suspended" {
		return nil, "", "", 0, errors.ErrForbidden
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", 0, errors.ErrUnauthorized
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, _, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("generate refresh token: %w", err)
	}

	tokenHash := hashToken(refreshToken)
	tokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: tokenExpiry,
	}
	if err := s.refreshRepo.Create(ctx, rt); err != nil {
		return nil, "", "", 0, fmt.Errorf("store refresh token: %w", err)
	}

	_ = s.refreshRepo.DeleteExpired(ctx)

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &user.ID,
		Action:     "user.login",
		EntityType: "user",
		EntityID:   &user.ID,
	})

	return user, accessToken, refreshToken, expiresIn, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, int64, error) {
	tokenHash := hashToken(refreshToken)
	rt, err := s.refreshRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil || rt == nil {
		return "", "", 0, errors.ErrUnauthorized
	}

	if err := s.refreshRepo.Revoke(ctx, rt.ID); err != nil {
		return "", "", 0, fmt.Errorf("revoke old token: %w", err)
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(rt.UserID, "")
	if err != nil {
		return "", "", 0, fmt.Errorf("generate access token: %w", err)
	}
	newRefreshToken, _, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, fmt.Errorf("generate refresh token: %w", err)
	}

	newHash := hashToken(newRefreshToken)
	newExpiry := time.Now().Add(7 * 24 * time.Hour)
	newRT := &model.RefreshToken{
		UserID:    rt.UserID,
		TokenHash: newHash,
		ExpiresAt: newExpiry,
	}
	if err := s.refreshRepo.Create(ctx, newRT); err != nil {
		return "", "", 0, fmt.Errorf("store new refresh token: %w", err)
	}

	return accessToken, newRefreshToken, expiresIn, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	rt, err := s.refreshRepo.GetByTokenHash(ctx, tokenHash)
	if err == nil && rt != nil {
		_ = s.refreshRepo.Revoke(ctx, rt.ID)
	}
	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.ErrForbidden
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user.PasswordHash = string(hashed)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (s *AuthService) VerifyPassword(ctx context.Context, userID uuid.UUID, password string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.ErrUnauthorized
	}
	return user, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (s *AuthService) sendWelcomeEmail(ctx context.Context, user *model.User) {
	if s.emailSvc == nil {
		return
	}
	html, err := email.RenderWelcome(email.WelcomeData{Name: user.Name, LoginURL: s.loginURL})
	if err != nil {
		return
	}
	s.emailSvc.SendAsync(user.Email, "Selamat datang di Owndangan", html)
}
