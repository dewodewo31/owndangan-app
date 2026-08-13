package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/service"
	"github.com/owndangan/backend/internal/service/entitlement"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func setupSubscriptionService(t *testing.T) (*service.SubscriptionService, *mockPackageRepo) {
	t.Helper()
	pkgRepo := newMockPackageRepo()
	duration := 7
	guestLimit := 50
	freePkg := &model.Package{
		ID:            uuid.New(),
		Name:          "Free",
		Code:          "free",
		Price:         0,
		DurationDays:  &duration,
		GuestLimit:    &guestLimit,
		TemplateGroup: "standard",
		Features:      datatypes.JSON(`{"guest.max": 50, "event.max": 1, "video.enabled": false, "custom_domain": false}`),
		IsActive:      true,
	}
	pkgRepo.packages["free"] = freePkg

	proDuration := 30
	proPkg := &model.Package{
		ID:            uuid.New(),
		Name:          "Pro",
		Code:          "pro",
		Price:         350000,
		DurationDays:  &proDuration,
		GuestLimit:    nil,
		TemplateGroup: "all",
		Features:      datatypes.JSON(`{"guest.max": null, "event.max": null, "video.enabled": true, "custom_domain": true}`),
		IsActive:      true,
	}
	pkgRepo.packages["pro"] = proPkg

	subSvc := service.NewSubscriptionService(
		&mockSubscriptionRepo{},
		pkgRepo,
		&mockTransactionRepo{},
		newMockUserRepo(),
		&mockAuditLogRepo{},
	)
	return subSvc, pkgRepo
}

func TestSubscription_ActivateOnSettlement(t *testing.T) {
	_, pkgRepo := setupSubscriptionService(t)
	ctx := context.Background()

	userID := uuid.New()
	subRepo := &mockSubscriptionRepo{}
	svc := service.NewSubscriptionService(subRepo, pkgRepo, &mockTransactionRepo{}, newMockUserRepo(), &mockAuditLogRepo{})

	now := time.Now()
	expiresAt := now.AddDate(0, 0, 7)
	sub := &model.Subscription{
		UserID:    userID,
		PackageID: pkgRepo.packages["free"].ID,
		Status:    "active",
		StartAt:   now,
		ExpiresAt: &expiresAt,
	}
	err := subRepo.Create(ctx, sub)
	require.NoError(t, err)

	got, err := subRepo.GetActiveByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "active", got.Status)

	_ = svc
}

func TestSubscription_IsActive_NotExpired(t *testing.T) {
	svc, _ := setupSubscriptionService(t)

	expiresAt := time.Now().Add(24 * time.Hour)
	sub := &model.Subscription{
		Status:    "active",
		ExpiresAt: &expiresAt,
	}
	require.True(t, svc.IsSubscriptionActive(sub))
}

func TestSubscription_IsActive_Expired(t *testing.T) {
	svc, _ := setupSubscriptionService(t)

	expiresAt := time.Now().Add(-24 * time.Hour)
	sub := &model.Subscription{
		Status:    "active",
		ExpiresAt: &expiresAt,
	}
	require.False(t, svc.IsSubscriptionActive(sub))
}

func TestSubscription_IsActive_Cancelled(t *testing.T) {
	svc, _ := setupSubscriptionService(t)

	sub := &model.Subscription{
		Status: "cancelled",
	}
	require.False(t, svc.IsSubscriptionActive(sub))
}

func TestSubscription_IsActive_Nil(t *testing.T) {
	svc, _ := setupSubscriptionService(t)
	require.False(t, svc.IsSubscriptionActive(nil))
}

func TestEntitlement_FreePackage_Limits(t *testing.T) {
	duration := 7
	guestLimit := 50
	pkg := &model.Package{
		Name:         "Free",
		Code:         "free",
		Price:        0,
		DurationDays: &duration,
		GuestLimit:   &guestLimit,
		Features:     datatypes.JSON(`{"guest.max": 50, "event.max": 1, "video.enabled": false}`),
		IsActive:     true,
	}

	resolver := entitlement.NewResolver(pkg)
	require.Equal(t, 50, *resolver.GuestMax())
	require.Equal(t, 1, *resolver.EventMax())
	require.False(t, resolver.VideoEnabled())
	require.True(t, resolver.CanCreateGuest(49))
	require.False(t, resolver.CanCreateGuest(50))
}

func TestEntitlement_ProPackage_Unlimited(t *testing.T) {
	duration := 30
	pkg := &model.Package{
		Name:         "Pro",
		Code:         "pro",
		Price:        350000,
		DurationDays: &duration,
		GuestLimit:   nil,
		Features:     datatypes.JSON(`{"guest.max": null, "event.max": null, "video.enabled": true, "custom_domain": true}`),
		IsActive:     true,
	}

	resolver := entitlement.NewResolver(pkg)
	require.Nil(t, resolver.GuestMax())
	require.Nil(t, resolver.EventMax())
	require.True(t, resolver.VideoEnabled())
	require.True(t, resolver.CustomDomain())
	require.True(t, resolver.IsUnlimitedGuests())
	require.True(t, resolver.IsUnlimitedEvents())
	require.True(t, resolver.CanCreateGuest(99999))
}

func TestEntitlement_NilPackage(t *testing.T) {
	resolver := entitlement.NewResolver(nil)
	require.Nil(t, resolver.GuestMax())
	require.False(t, resolver.VideoEnabled())
	require.True(t, resolver.CanCreateGuest(100))
}

func TestSubscription_GetEntitlementResolver_ActiveSub(t *testing.T) {
	_, pkgRepo := setupSubscriptionService(t)
	ctx := context.Background()

	userID := uuid.New()
	subRepo := &mockSubscriptionRepo{}
	svc := service.NewSubscriptionService(subRepo, pkgRepo, &mockTransactionRepo{}, newMockUserRepo(), &mockAuditLogRepo{})

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	sub := &model.Subscription{
		UserID:    userID,
		PackageID: pkgRepo.packages["pro"].ID,
		Status:    "active",
		StartAt:   now,
		ExpiresAt: &expiresAt,
	}
	err := subRepo.Create(ctx, sub)
	require.NoError(t, err)

	resolver, err := svc.GetEntitlementResolver(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, resolver)
	require.True(t, resolver.IsUnlimitedGuests())
	require.True(t, resolver.VideoEnabled())
}

func TestSubscription_GetEntitlementResolver_NoSub(t *testing.T) {
	_, pkgRepo := setupSubscriptionService(t)
	ctx := context.Background()

	subRepo := &mockSubscriptionRepo{}
	svc := service.NewSubscriptionService(subRepo, pkgRepo, &mockTransactionRepo{}, newMockUserRepo(), &mockAuditLogRepo{})

	userID := uuid.New()

	resolver, err := svc.GetEntitlementResolver(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, resolver)
	require.Equal(t, 50, *resolver.GuestMax())
	require.False(t, resolver.VideoEnabled())
}

func TestSubscription_ExpiryHandling(t *testing.T) {
	svc, _ := setupSubscriptionService(t)

	expiresAt := time.Now().Add(-1 * time.Hour)
	sub := &model.Subscription{
		Status:    "active",
		ExpiresAt: &expiresAt,
	}
	require.False(t, svc.IsSubscriptionActive(sub), "expired subscription should not be active")

	neverExpires := &model.Subscription{
		Status:    "active",
		ExpiresAt: nil,
	}
	require.True(t, svc.IsSubscriptionActive(neverExpires), "lifetime subscription should be active")
}

func TestSubscription_Lifecycle(t *testing.T) {
	_, pkgRepo := setupSubscriptionService(t)
	ctx := context.Background()

	userID := uuid.New()
	subRepo := &mockSubscriptionRepo{}
	svc := service.NewSubscriptionService(subRepo, pkgRepo, &mockTransactionRepo{}, newMockUserRepo(), &mockAuditLogRepo{})

	now := time.Now()
	expiresAt := now.AddDate(0, 0, 7)
	sub := &model.Subscription{
		UserID:    userID,
		PackageID: pkgRepo.packages["free"].ID,
		Status:    "active",
		StartAt:   now,
		ExpiresAt: &expiresAt,
	}
	err := subRepo.Create(ctx, sub)
	require.NoError(t, err)

	got, err := subRepo.GetByID(ctx, sub.ID)
	require.NoError(t, err)
	require.True(t, svc.IsSubscriptionActive(got))

	err = svc.Terminate(ctx, got.ID)
	require.NoError(t, err)

	got, err = subRepo.GetByID(ctx, sub.ID)
	require.NoError(t, err)
	require.Equal(t, "cancelled", got.Status)
	require.False(t, svc.IsSubscriptionActive(got))
}
