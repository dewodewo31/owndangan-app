package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/entitlement"
)

type SubscriptionService struct {
	subRepo   repository.SubscriptionRepository
	pkgRepo   repository.PackageRepository
	txnRepo   repository.TransactionRepository
	userRepo  repository.UserRepository
	auditRepo repository.AuditLogRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, pkgRepo repository.PackageRepository, txnRepo repository.TransactionRepository, userRepo repository.UserRepository, auditRepo repository.AuditLogRepository) *SubscriptionService {
	return &SubscriptionService{
		subRepo:   subRepo,
		pkgRepo:   pkgRepo,
		txnRepo:   txnRepo,
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (s *SubscriptionService) GetCurrentUserSubscription(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error) {
	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if sub == nil {
		return nil, errors.ErrNotFound
	}
	pkg, _ := s.pkgRepo.GetByID(ctx, sub.PackageID)
	return toSubscriptionResponse(sub, pkg), nil
}

func (s *SubscriptionService) GetUserSubscriptionOrDefault(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error) {
	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil || sub == nil {
		freePkg, err := s.pkgRepo.GetByCode(ctx, "free")
		if err != nil {
			return nil, errors.ErrNotFound
		}
		return &dto.SubscriptionResponse{
			ID:      uuid.New(),
			Package: toPackageBrief(freePkg),
			Status:  "active",
			StartAt: time.Now(),
			ExpiresAt: func() *time.Time {
				t := time.Now().Add(7 * 24 * time.Hour)
				return &t
			}(),
		}, nil
	}
	pkg, _ := s.pkgRepo.GetByID(ctx, sub.PackageID)
	return toSubscriptionResponse(sub, pkg), nil
}

func (s *SubscriptionService) ActivateOnSettlement(ctx context.Context, transactionID uuid.UUID) (*model.Subscription, error) {
	txn, err := s.txnRepo.GetByID(ctx, transactionID)
	if err != nil || txn == nil {
		return nil, fmt.Errorf("transaction not found: %w", errors.ErrNotFound)
	}
	if txn.Status != "settlement" {
		return nil, fmt.Errorf("transaction not settled: %w", errors.ErrConflict)
	}

	pkg, err := s.pkgRepo.GetByID(ctx, txn.PackageID)
	if err != nil || pkg == nil {
		return nil, fmt.Errorf("package not found: %w", errors.ErrNotFound)
	}

	now := time.Now()
	expiresAt := computeExpiry(pkg, now)

	if err := s.subRepo.DeactivateActive(ctx, txn.UserID); err != nil {
		return nil, fmt.Errorf("deactivate existing subscription: %w", err)
	}

	newSub := &model.Subscription{
		UserID:        txn.UserID,
		PackageID:     txn.PackageID,
		TransactionID: &txn.ID,
		Status:        "active",
		StartAt:       now,
		ExpiresAt:     expiresAt,
	}
	if err := s.subRepo.Create(ctx, newSub); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &txn.UserID,
		Action:     "subscription.activated",
		EntityType: "subscription",
		EntityID:   &newSub.ID,
		Metadata:   datatypesJSON(map[string]interface{}{"transaction_id": txn.ID, "package_id": txn.PackageID}),
	})

	return newSub, nil
}

func (s *SubscriptionService) Extend(ctx context.Context, subID uuid.UUID, days int) error {
	sub, err := s.subRepo.GetByID(ctx, subID)
	if err != nil || sub == nil {
		return errors.ErrNotFound
	}
	if days <= 0 {
		return fmt.Errorf("%w: days must be positive", errors.ErrInvalidInput)
	}
	if sub.ExpiresAt == nil {
		return fmt.Errorf("%w: cannot extend lifetime subscription", errors.ErrConflict)
	}
	newExpiry := sub.ExpiresAt.AddDate(0, 0, days)
	sub.ExpiresAt = &newExpiry
	return s.subRepo.Update(ctx, sub)
}

func (s *SubscriptionService) Terminate(ctx context.Context, subID uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, subID)
	if err != nil || sub == nil {
		return errors.ErrNotFound
	}
	if sub.Status != "active" {
		return fmt.Errorf("%w: subscription is not active", errors.ErrConflict)
	}
	sub.Status = "cancelled"
	return s.subRepo.Update(ctx, sub)
}

func computeExpiry(pkg *model.Package, now time.Time) *time.Time {
	if pkg.DurationDays == nil {
		return nil
	}
	t := now.AddDate(0, 0, *pkg.DurationDays)
	return &t
}

func toSubscriptionResponse(sub *model.Subscription, pkg *model.Package) *dto.SubscriptionResponse {
	return &dto.SubscriptionResponse{
		ID:        sub.ID,
		Package:   toPackageBrief(pkg),
		Status:    sub.Status,
		StartAt:   sub.StartAt,
		ExpiresAt: sub.ExpiresAt,
	}
}

func toPackageBrief(pkg *model.Package) dto.PackageBrief {
	if pkg == nil {
		return dto.PackageBrief{}
	}
	return dto.PackageBrief{
		ID:            pkg.ID,
		Name:          pkg.Name,
		Code:          pkg.Code,
		Price:         pkg.Price,
		GuestLimit:    pkg.GuestLimit,
		TemplateGroup: pkg.TemplateGroup,
		Features:      pkg.Features,
	}
}

func (s *SubscriptionService) GetEntitlementResolver(ctx context.Context, userID uuid.UUID) (*entitlement.Resolver, error) {
	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil || sub == nil {
		freePkg, err := s.pkgRepo.GetByCode(ctx, "free")
		if err != nil || freePkg == nil {
			return entitlement.NewResolver(nil), nil
		}
		return entitlement.NewResolver(freePkg), nil
	}
	pkg, err := s.pkgRepo.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return entitlement.NewResolver(nil), nil
	}
	return entitlement.NewResolver(pkg), nil
}

func (s *SubscriptionService) IsSubscriptionActive(sub *model.Subscription) bool {
	if sub == nil {
		return false
	}
	if sub.Status != "active" {
		return false
	}
	if sub.ExpiresAt != nil && sub.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

func (s *SubscriptionService) CheckExpiration(ctx context.Context) error {
	expiredCount, err := s.subRepo.CountExpired(ctx)
	if err != nil {
		return fmt.Errorf("count expired subscriptions: %w", err)
	}
	if expiredCount == 0 {
		return nil
	}
	return nil
}

func (s *SubscriptionService) Upgrade(ctx context.Context, userID uuid.UUID, newPackageID uuid.UUID) error {
	_, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: no active subscription", errors.ErrNotFound)
	}
	newPkg, err := s.pkgRepo.GetByID(ctx, newPackageID)
	if err != nil || newPkg == nil {
		return fmt.Errorf("%w: package not found", errors.ErrNotFound)
	}
	now := time.Now()
	expiresAt := computeExpiry(newPkg, now)
	if err := s.subRepo.DeactivateActive(ctx, userID); err != nil {
		return fmt.Errorf("deactivate current subscription: %w", err)
	}
	newSub := &model.Subscription{
		UserID:    userID,
		PackageID: newPackageID,
		Status:    "active",
		StartAt:   now,
		ExpiresAt: expiresAt,
	}
	if err := s.subRepo.Create(ctx, newSub); err != nil {
		return fmt.Errorf("create new subscription: %w", err)
	}
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "subscription.upgraded",
		EntityType: "subscription",
		EntityID:   &newSub.ID,
		Metadata:   datatypesJSON(map[string]interface{}{"package_id": newPackageID}),
	})
	return nil
}
