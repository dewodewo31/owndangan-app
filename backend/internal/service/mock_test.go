package service_test

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/repository"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	users    map[uuid.UUID]*model.User
	emailIdx map[string]*model.User
	createFn func(ctx context.Context, user *model.User) error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:    make(map[uuid.UUID]*model.User),
		emailIdx: make(map[string]*model.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	if _, exists := m.emailIdx[user.Email]; exists {
		return gorm.ErrDuplicatedKey
	}
	user.ID = uuid.New()
	m.users[user.ID] = user
	m.emailIdx[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return &model.User{ID: id, Name: "Test", Email: "test@example.com", Role: "user", Status: "active"}, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if u, ok := m.emailIdx[email]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Update(ctx context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(m.users)), nil
}

func (m *mockUserRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
	return 0, nil
}

func (m *mockUserRepo) CountByRole(ctx context.Context, role string) (int64, error) {
	return 0, nil
}

func (m *mockUserRepo) List(ctx context.Context, page, perPage int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return nil
}

func (m *mockUserRepo) WithTx(tx *gorm.DB) repository.UserRepository {
	return m
}

type mockRefreshTokenRepo struct {
	tokens []*model.RefreshToken
}

func (m *mockRefreshTokenRepo) Create(ctx context.Context, token *model.RefreshToken) error {
	m.tokens = append(m.tokens, token)
	return nil
}

func (m *mockRefreshTokenRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	return nil, errors.New("not found")
}

func (m *mockRefreshTokenRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockRefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (m *mockRefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	return nil
}

type mockPackageRepo struct {
	packages map[string]*model.Package
}

func newMockPackageRepo() *mockPackageRepo {
	return &mockPackageRepo{packages: make(map[string]*model.Package)}
}

func (m *mockPackageRepo) Create(ctx context.Context, pkg *model.Package) error {
	m.packages[pkg.Code] = pkg
	return nil
}

func (m *mockPackageRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	for _, p := range m.packages {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockPackageRepo) GetByCode(ctx context.Context, code string) (*model.Package, error) {
	if p, ok := m.packages[code]; ok {
		return p, nil
	}
	return nil, errors.New("not found")
}

func (m *mockPackageRepo) GetAllActive(ctx context.Context) ([]model.Package, error) {
	var result []model.Package
	for _, p := range m.packages {
		if p.IsActive {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPackageRepo) GetAllWithInactive(ctx context.Context) ([]model.Package, error) {
	return nil, nil
}

func (m *mockPackageRepo) Update(ctx context.Context, pkg *model.Package) error {
	m.packages[pkg.Code] = pkg
	return nil
}

func (m *mockPackageRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockPackageRepo) WithTx(tx *gorm.DB) repository.PackageRepository {
	return m
}

type mockSubscriptionRepo struct {
	subs []*model.Subscription
}

func (m *mockSubscriptionRepo) Create(ctx context.Context, sub *model.Subscription) error {
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}
	m.subs = append(m.subs, sub)
	return nil
}

func (m *mockSubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	for _, s := range m.subs {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockSubscriptionRepo) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	for _, s := range m.subs {
		if s.UserID == userID && s.Status == "active" {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockSubscriptionRepo) GetByTransactionID(ctx context.Context, txnID uuid.UUID) (*model.Subscription, error) {
	return nil, errors.New("not found")
}

func (m *mockSubscriptionRepo) Update(ctx context.Context, sub *model.Subscription) error {
	return nil
}

func (m *mockSubscriptionRepo) CountExpired(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockSubscriptionRepo) DeactivateActive(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (m *mockSubscriptionRepo) ListExpiringBetween(ctx context.Context, from, to time.Time) ([]model.Subscription, error) {
	return nil, nil
}

func (m *mockSubscriptionRepo) ListExpiredActive(ctx context.Context, now time.Time) ([]model.Subscription, error) {
	return nil, nil
}

type mockAuditLogRepo struct{}

func (m *mockAuditLogRepo) Create(ctx context.Context, log *model.AuditLog) error {
	return nil
}

func (m *mockAuditLogRepo) ExistsSince(ctx context.Context, action, entityType string, entityID uuid.UUID, since time.Time) (bool, error) {
	return false, nil
}

type mockTransactionRepo struct {
	txns []*model.Transaction
}

func (m *mockTransactionRepo) Create(ctx context.Context, txn *model.Transaction) error {
	if txn.ID == uuid.Nil {
		txn.ID = uuid.New()
	}
	m.txns = append(m.txns, txn)
	return nil
}

func (m *mockTransactionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	for _, t := range m.txns {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTransactionRepo) GetByOrderID(ctx context.Context, orderID string) (*model.Transaction, error) {
	for _, t := range m.txns {
		if t.OrderID == orderID {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTransactionRepo) Update(ctx context.Context, txn *model.Transaction) error {
	for i, t := range m.txns {
		if t.ID == txn.ID {
			m.txns[i] = txn
			return nil
		}
	}
	return nil
}

func (m *mockTransactionRepo) ListByUserID(ctx context.Context, userID uuid.UUID, page, perPage int) ([]model.Transaction, int64, error) {
	var result []model.Transaction
	for _, t := range m.txns {
		result = append(result, *t)
	}
	return result, int64(len(result)), nil
}

func (m *mockTransactionRepo) ListAll(ctx context.Context, page, perPage int, status string) ([]model.Transaction, int64, error) {
	var result []model.Transaction
	for _, t := range m.txns {
		if status == "" || t.Status == status {
			result = append(result, *t)
		}
	}
	return result, int64(len(result)), nil
}

type mockWebhookIdempotencyRepo struct {
	processed map[string]bool
}

func (m *mockWebhookIdempotencyRepo) IsProcessed(ctx context.Context, requestID string) (bool, error) {
	if m.processed == nil {
		m.processed = make(map[string]bool)
	}
	return m.processed[requestID], nil
}

func (m *mockWebhookIdempotencyRepo) MarkProcessed(ctx context.Context, requestID, orderID, status string) error {
	if m.processed == nil {
		m.processed = make(map[string]bool)
	}
	m.processed[requestID] = true
	return nil
}


