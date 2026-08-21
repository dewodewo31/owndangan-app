package email

import (
	"context"
	"fmt"
	"time"

	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/repository"
	"github.com/rs/zerolog"
)

type Sender interface {
	SendWithRetry(to, subject, htmlBody string) error
}

type ExpiryWorker struct {
	subs     repository.SubscriptionRepository
	users    repository.UserRepository
	pkgs     repository.PackageRepository
	audit    repository.AuditLogRepository
	sender   Sender
	renewURL string
	log      zerolog.Logger
	Clock    func() time.Time
}

func NewExpiryWorker(subs repository.SubscriptionRepository, users repository.UserRepository,
	pkgs repository.PackageRepository, audit repository.AuditLogRepository,
	sender Sender, renewURL string, log zerolog.Logger) *ExpiryWorker {
	return &ExpiryWorker{
		subs:     subs,
		users:    users,
		pkgs:     pkgs,
		audit:    audit,
		sender:   sender,
		renewURL: renewURL,
		log:      log,
		Clock:    time.Now,
	}
}

const (
	actionExpiryReminder = "email.expiry_reminder"
	subjectExpiryReminder = "Langganan Owndangan akan berakhir"
	subjectExpired        = "Langganan Owndangan telah berakhir"
)

// RunOnce scans twice: active subscriptions expiring within [now+6d, now+8d]
// get one reminder per day (audit-log dedupe); still-active subscriptions past
// their expiry get a final notice and are flipped to status "expired".
func (w *ExpiryWorker) RunOnce(ctx context.Context) error {
	now := w.Clock()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	reminders, err := w.subs.ListExpiringBetween(ctx, now.AddDate(0, 0, 6), now.AddDate(0, 0, 8))
	if err != nil {
		return fmt.Errorf("list expiring subscriptions: %w", err)
	}
	for i := range reminders {
		sub := &reminders[i]
		sent, err := w.audit.ExistsSince(ctx, actionExpiryReminder, "subscription", sub.ID, dayStart)
		if err != nil {
			w.log.Warn().Err(err).Str("subscription_id", sub.ID.String()).Msg("reminder dedupe check failed")
			continue
		}
		if sent {
			continue
		}
		if w.remind(ctx, sub) {
			w.audit.Create(ctx, &model.AuditLog{
				UserID:     &sub.UserID,
				Action:     actionExpiryReminder,
				EntityType: "subscription",
				EntityID:   &sub.ID,
			})
		}
	}

	expired, err := w.subs.ListExpiredActive(ctx, now)
	if err != nil {
		return fmt.Errorf("list expired subscriptions: %w", err)
	}
	for i := range expired {
		sub := &expired[i]
		w.notifyExpired(ctx, sub)
		sub.Status = "expired"
		if err := w.subs.Update(ctx, sub); err != nil {
			w.log.Error().Err(err).Str("subscription_id", sub.ID.String()).Msg("failed to mark subscription expired")
		}
	}
	return nil
}

func (w *ExpiryWorker) remind(ctx context.Context, sub *model.Subscription) bool {
	user, pkg, ok := w.loadUserAndPkg(ctx, sub)
	if !ok || sub.ExpiresAt == nil {
		return false
	}
	html, err := RenderExpiryReminder(ExpiryReminderData{
		Name:       user.Name,
		PlanName:   pkg.Name,
		ExpiryDate: sub.ExpiresAt.Format("2 January 2006"),
		RenewURL:   w.renewURL,
	})
	if err != nil {
		w.log.Error().Err(err).Msg("render expiry reminder failed")
		return false
	}
	if err := w.sender.SendWithRetry(user.Email, subjectExpiryReminder, html); err != nil {
		w.log.Error().Err(err).Str("to", user.Email).Msg("expiry reminder send failed")
		return false
	}
	return true
}

func (w *ExpiryWorker) notifyExpired(ctx context.Context, sub *model.Subscription) {
	user, pkg, ok := w.loadUserAndPkg(ctx, sub)
	if !ok {
		return
	}
	html, err := RenderExpired(ExpiredData{Name: user.Name, PlanName: pkg.Name, RenewURL: w.renewURL})
	if err != nil {
		w.log.Error().Err(err).Msg("render expired notification failed")
		return
	}
	if err := w.sender.SendWithRetry(user.Email, subjectExpired, html); err != nil {
		w.log.Error().Err(err).Str("to", user.Email).Msg("expired notification send failed")
	}
}

func (w *ExpiryWorker) loadUserAndPkg(ctx context.Context, sub *model.Subscription) (*model.User, *model.Package, bool) {
	user, err := w.users.GetByID(ctx, sub.UserID)
	if err != nil || user == nil {
		return nil, nil, false
	}
	pkg, err := w.pkgs.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return nil, nil, false
	}
	return user, pkg, true
}
