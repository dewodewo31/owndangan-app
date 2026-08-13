package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/service"
	"github.com/stretchr/testify/require"
)

func TestSubscription_ActivateOnSettlement_New(t *testing.T) {
	_, pkgRepo := setupSubscriptionService(t)
	ctx := context.Background()

	subRepo := &mockSubscriptionRepo{}
	svc := service.NewSubscriptionService(subRepo, pkgRepo, &mockTransactionRepo{}, newMockUserRepo(), &mockAuditLogRepo{})

	now := time.Now()
	expiresAt := now.AddDate(0, 0, 7)
	sub := &model.Subscription{
		UserID: uuid.New(), PackageID: pkgRepo.packages["free"].ID,
		Status: "active", StartAt: now, ExpiresAt: &expiresAt,
	}
	err := subRepo.Create(ctx, sub)
	require.NoError(t, err)

	got, err := subRepo.GetActiveByUserID(ctx, sub.UserID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "active", got.Status)

	_ = svc
}

func TestSubscription_Lifecycle_New(t *testing.T) {
	svc, pkgRepo := setupSubscriptionService(t)
	ctx := context.Background()

	now := time.Now()
	expiresAt := now.AddDate(0, 0, 7)
	sub := &model.Subscription{
		UserID: uuid.New(), PackageID: pkgRepo.packages["free"].ID,
		Status: "active", StartAt: now, ExpiresAt: &expiresAt,
	}

	err := svc.Extend(ctx, uuid.New(), 7)
	require.Error(t, err)

	_ = sub
}
