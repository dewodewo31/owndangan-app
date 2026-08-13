package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   repository.UserRepository
	subRepo    repository.SubscriptionRepository
	pkgRepo    repository.PackageRepository
	auditRepo  repository.AuditLogRepository
}

func NewUserService(userRepo repository.UserRepository, subRepo repository.SubscriptionRepository, pkgRepo repository.PackageRepository, auditRepo repository.AuditLogRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		subRepo:   subRepo,
		pkgRepo:   pkgRepo,
		auditRepo: auditRepo,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound
	}
	return toUserResponse(user), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound
	}

	tx := func() {
		if req.Name != nil {
			user.Name = strings.TrimSpace(*req.Name)
		}
		if req.Email != nil {
			existing, _ := s.userRepo.GetByEmail(ctx, *req.Email)
			if existing != nil && existing.ID != userID {
				return
			}
			user.Email = strings.ToLower(strings.TrimSpace(*req.Email))
		}
		if req.Phone != nil {
			user.Phone = *req.Phone
		}
		if req.AvatarURL != nil {
			user.AvatarURL = *req.AvatarURL
		}
	}
	tx()

	if req.Email != nil {
		existing, _ := s.userRepo.GetByEmail(ctx, *req.Email)
		if existing != nil && existing.ID != userID {
			return nil, fmt.Errorf("%w: email already taken", errors.ErrConflict)
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:    &user.ID,
		Action:    "user.profile_updated",
		EntityType: "user",
		EntityID:  &user.ID,
	})

	return toUserResponse(user), nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.ErrForbidden
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user.PasswordHash = string(hashed)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Status:    user.Status,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *UserService) GetSubscription(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error) {
	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil || sub == nil {
		return nil, errors.ErrNotFound
	}
	pkg, _ := s.pkgRepo.GetByID(ctx, sub.PackageID)
	return toSubscriptionResponse(sub, pkg), nil
}

func (s *UserService) EnsureFreeSubscription(ctx context.Context, userID uuid.UUID) error {
	sub, _ := s.subRepo.GetActiveByUserID(ctx, userID)
	if sub != nil {
		return nil
	}
	freePkg, err := s.pkgRepo.GetByCode(ctx, "free")
	if err != nil {
		return fmt.Errorf("get free package: %w", err)
	}
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour)
	freeSub := &model.Subscription{
		UserID:    userID,
		PackageID: freePkg.ID,
		Status:    "active",
		StartAt:   now,
		ExpiresAt: &expiresAt,
	}
	return s.subRepo.Create(ctx, freeSub)
}

func (s *UserService) GetUserWithSubscription(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrNotFound
	}
	return toUserResponse(user), nil
}
