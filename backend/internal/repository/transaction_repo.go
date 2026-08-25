package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, txn *model.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetByOrderID(ctx context.Context, orderID string) (*model.Transaction, error)
	GetPendingByUserAndPackage(ctx context.Context, userID uuid.UUID, packageID uuid.UUID) (*model.Transaction, error)
	Update(ctx context.Context, txn *model.Transaction) error
	ListByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]model.Transaction, int64, error)
	ListAll(ctx context.Context, page, perPage int, status, packageID string) ([]model.Transaction, int64, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, txn *model.Transaction) error {
	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *transactionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	var txn model.Transaction
	err := r.db.WithContext(ctx).
		Preload("Package").
		Where("id = ?", id).First(&txn).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *transactionRepo) GetByOrderID(ctx context.Context, orderID string) (*model.Transaction, error) {
	var txn model.Transaction
	err := r.db.WithContext(ctx).
		Preload("Package").
		Where("order_id = ?", orderID).First(&txn).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *transactionRepo) Update(ctx context.Context, txn *model.Transaction) error {
	return r.db.WithContext(ctx).Save(txn).Error
}

func (r *transactionRepo) GetPendingByUserAndPackage(ctx context.Context, userID uuid.UUID, packageID uuid.UUID) (*model.Transaction, error) {
	var txn model.Transaction
	err := r.db.WithContext(ctx).
		Preload("Package").
		Where("user_id = ? AND package_id = ? AND status = ?", userID, packageID, "pending").
		Order("created_at DESC").
		First(&txn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &txn, nil
}

func (r *transactionRepo) ListByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]model.Transaction, int64, error) {
	var txns []model.Transaction
	var total int64
	offset := (page - 1) * perPage
	query := r.db.WithContext(ctx).Preload("Package").Model(&model.Transaction{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&txns).Error
	return txns, total, err
}

func (r *transactionRepo) ListAll(ctx context.Context, page, perPage int, status, packageID string) ([]model.Transaction, int64, error) {
	var txns []model.Transaction
	var total int64
	offset := (page - 1) * perPage
	query := r.db.WithContext(ctx).Preload("Package").Preload("User").Model(&model.Transaction{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if packageID != "" {
		query = query.Where("package_id = ?", packageID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&txns).Error
	return txns, total, err
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error)
	GetByTransactionID(ctx context.Context, txnID uuid.UUID) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	CountActive(ctx context.Context) (int64, error)
	CountExpired(ctx context.Context) (int64, error)
	DeactivateActive(ctx context.Context, userID uuid.UUID) error
	ListExpiringBetween(ctx context.Context, from, to time.Time) ([]model.Subscription, error)
	ListExpiredActive(ctx context.Context, now time.Time) ([]model.Subscription, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *model.Subscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.WithContext(ctx).
		Preload("Package").
		Preload("User").
		Where("id = ?", id).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepo) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.WithContext(ctx).
		Preload("Package").
		Preload("User").
		Where("user_id = ? AND status = ? AND (expires_at IS NULL OR expires_at > ?)", userID, "active", time.Now()).
		First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepo) GetByTransactionID(ctx context.Context, txnID uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.WithContext(ctx).
		Preload("Package").
		Preload("User").
		Where("transaction_id = ?", txnID).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, sub *model.Subscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *subscriptionRepo) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Subscription{}).
		Where("status = ?", "active").
		Count(&count).Error
	return count, err
}

func (r *subscriptionRepo) CountExpired(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Subscription{}).
		Where("expires_at < ? AND status = ?", time.Now(), "active").
		Count(&count).Error
	return count, err
}

func (r *subscriptionRepo) DeactivateActive(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Subscription{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "cancelled").Error
}

func (r *subscriptionRepo) ListExpiringBetween(ctx context.Context, from, to time.Time) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at >= ? AND expires_at <= ?", "active", from, to).
		Find(&subs).Error
	return subs, err
}

func (r *subscriptionRepo) ListExpiredActive(ctx context.Context, now time.Time) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at < ?", "active", now).
		Find(&subs).Error
	return subs, err
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	ExistsSince(ctx context.Context, action, entityType string, entityID uuid.UUID, since time.Time) (bool, error)
}

type auditLogRepo struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepo) ExistsSince(ctx context.Context, action, entityType string, entityID uuid.UUID, since time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AuditLog{}).
		Where("action = ? AND entity_type = ? AND entity_id = ? AND created_at >= ?", action, entityType, entityID, since).
		Count(&count).Error
	return count > 0, err
}

type WebhookIdempotencyRepository interface {
	IsProcessed(ctx context.Context, requestID string) (bool, error)
	MarkProcessed(ctx context.Context, requestID, orderID, status string) error
}

type webhookIdempotencyRepo struct {
	db *gorm.DB
}

func NewWebhookIdempotencyRepository(db *gorm.DB) WebhookIdempotencyRepository {
	return &webhookIdempotencyRepo{db: db}
}

func (r *webhookIdempotencyRepo) IsProcessed(ctx context.Context, requestID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.WebhookIdempotency{}).
		Where("transaction_id = ?", requestID).Count(&count).Error
	return count > 0, err
}

func (r *webhookIdempotencyRepo) MarkProcessed(ctx context.Context, requestID, orderID, status string) error {
	return r.db.WithContext(ctx).Create(&model.WebhookIdempotency{
		TransactionID: requestID,
		OrderID:       orderID,
		Status:        status,
	}).Error
}
