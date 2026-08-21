package handler_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/email"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newExpiryWorker(db *gorm.DB, sender *recEmailSender, fixed time.Time) *email.ExpiryWorker {
	worker := email.NewExpiryWorker(
		repository.NewSubscriptionRepository(db),
		repository.NewUserRepository(db),
		repository.NewPackageRepository(db),
		repository.NewAuditLogRepository(db),
		sender,
		"https://app.test/dashboard/billing",
		zerolog.Nop(),
	)
	worker.Clock = func() time.Time { return fixed }
	return worker
}

func seedActiveSub(t *testing.T, db *gorm.DB, fixed time.Time, expiresAt time.Time) uuid.UUID {
	userID := uuid.New()
	require.NoError(t, db.Create(&model.User{
		ID:           userID,
		Name:         "Exp User",
		Email:        fmt.Sprintf("exp_%s@test.com", userID.String()[:8]),
		PasswordHash: "x",
		Role:         "user",
		Status:       "active",
	}).Error)
	var pkg model.Package
	require.NoError(t, db.Where("code = ?", "starter").First(&pkg).Error)
	subID := uuid.New()
	require.NoError(t, db.Create(&model.Subscription{
		ID:        subID,
		UserID:    userID,
		PackageID: pkg.ID,
		Status:    "active",
		StartAt:   fixed.AddDate(0, 0, -30),
		ExpiresAt: &expiresAt,
	}).Error)
	return subID
}

func TestExpiryWorker_Reminder_Dedupe(t *testing.T) {
	db := setupAuthTestServer(t)
	rec := &recEmailSender{}
	fixed := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	subID := seedActiveSub(t, db, fixed, fixed.AddDate(0, 0, 7))

	worker := newExpiryWorker(db, rec, fixed)
	ctx := context.Background()

	require.NoError(t, worker.RunOnce(ctx))
	require.Equal(t, 1, rec.count(), "subscription expiring in 7d gets exactly one reminder")

	require.NoError(t, worker.RunOnce(ctx))
	require.Equal(t, 1, rec.count(), "same-day rerun must not double-send")

	var got model.Subscription
	require.NoError(t, db.First(&got, "id = ?", subID).Error)
	require.Equal(t, "active", got.Status, "reminder must not change subscription status")
}

func TestExpiryWorker_Expired_FlipsStatus(t *testing.T) {
	db := setupAuthTestServer(t)
	rec := &recEmailSender{}
	fixed := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	subID := seedActiveSub(t, db, fixed, fixed.Add(-time.Hour))

	worker := newExpiryWorker(db, rec, fixed)
	ctx := context.Background()

	require.NoError(t, worker.RunOnce(ctx))
	require.Equal(t, 1, rec.count(), "expired active subscription gets the expired notice")

	var got model.Subscription
	require.NoError(t, db.First(&got, "id = ?", subID).Error)
	require.Equal(t, "expired", got.Status)

	require.NoError(t, worker.RunOnce(ctx))
	require.Equal(t, 1, rec.count(), "already-flipped subscription must not re-notify")
}
