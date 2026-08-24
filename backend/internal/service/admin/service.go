package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"gorm.io/datatypes"
)

type Service struct {
	userRepo     repository.UserRepository
	pkgRepo      repository.PackageRepository
	txnRepo      repository.TransactionRepository
	templateRepo repository.TemplateRepository
	subRepo      repository.SubscriptionRepository
	eventRepo    repository.EventRepository
	auditRepo    repository.AuditLogRepository
}

func NewService(userRepo repository.UserRepository, pkgRepo repository.PackageRepository,
	txnRepo repository.TransactionRepository, templateRepo repository.TemplateRepository,
	subRepo repository.SubscriptionRepository, eventRepo repository.EventRepository,
	auditRepo repository.AuditLogRepository) *Service {
	return &Service{
		userRepo:     userRepo,
		pkgRepo:      pkgRepo,
		txnRepo:      txnRepo,
		templateRepo: templateRepo,
		subRepo:      subRepo,
		eventRepo:    eventRepo,
		auditRepo:    auditRepo,
	}
}

func (s *Service) GetAnalytics(ctx context.Context) (*Analytics, error) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	totalUsers, _ := s.userRepo.Count(ctx)
	activeUsers, _ := s.userRepo.CountByStatus(ctx, "active")

	_, totalTxns, _ := s.txnRepo.ListAll(ctx, 1, 1, "", "")
	settlementTxns, _, _ := s.txnRepo.ListAll(ctx, 1, 10000, "settlement", "")
	var totalRevenue int64
	for _, t := range settlementTxns {
		totalRevenue += t.GrossAmount
	}

	pkgs, _ := s.pkgRepo.GetAllActive(ctx)
	templates, _ := s.templateRepo.ListAll(ctx)

	activeSubs, _ := s.subRepo.CountActive(ctx)
	totalInv, _ := s.eventRepo.Count(ctx)

	recentTxns, _, _ := s.txnRepo.ListAll(ctx, 1, 5, "", "")
	recentUsers, _, _ := s.userRepo.List(ctx, 1, 5, "", "")

	return &Analytics{
		TotalUsers:          totalUsers,
		ActiveUsers:         activeUsers,
		TotalTransactions:   totalTxns,
		TotalRevenue:        totalRevenue,
		SettlementCount:     len(settlementTxns),
		ActivePackages:      len(pkgs),
		ActiveTemplates:     len(templates),
		ActiveSubscriptions: activeSubs,
		TotalInvitations:    totalInv,
		RecentTransactions:  toRecentTransactions(recentTxns),
		RecentUsers:         toRecentUsers(recentUsers),
		PeriodStart:         thirtyDaysAgo,
		PeriodEnd:           now,
	}, nil
}

func toRecentTransactions(txns []model.Transaction) []RecentTransaction {
	out := make([]RecentTransaction, 0, len(txns))
	for _, t := range txns {
		ts := t.CreatedAt
		if t.TransactionTime != nil {
			ts = *t.TransactionTime
		}
		name, email := "", ""
		if t.User.ID != uuid.Nil {
			name = t.User.Name
			email = t.User.Email
		}
		out = append(out, RecentTransaction{
			ID:        t.ID,
			UserName:  name,
			UserEmail: email,
			Amount:    t.GrossAmount,
			Status:    t.Status,
			Timestamp: ts,
		})
	}
	return out
}

func toRecentUsers(users []model.User) []RecentUser {
	out := make([]RecentUser, 0, len(users))
	for _, u := range users {
		out = append(out, RecentUser{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		})
	}
	return out
}

func (s *Service) ListUsers(ctx context.Context, page, perPage int, search, status string) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, page, perPage, search, status)
}

func (s *Service) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error {
	if status != "active" && status != "suspended" {
		return fmt.Errorf("%w: invalid status", errors.ErrInvalidInput)
	}

	err := s.userRepo.UpdateStatus(ctx, userID, status)
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "admin.user_status_updated",
		EntityType: "user",
		EntityID:   &userID,
	})

	return nil
}

func (s *Service) GetPackages(ctx context.Context) ([]model.Package, error) {
	return s.pkgRepo.GetAllWithInactive(ctx)
}

func (s *Service) CreatePackage(ctx context.Context, req CreatePackageRequest) (*model.Package, error) {
	pkg := &model.Package{
		Name:          req.Name,
		Code:          req.Code,
		Price:         req.Price,
		DurationDays:  req.DurationDays,
		GuestLimit:    req.GuestLimit,
		TemplateGroup: req.TemplateGroup,
		Features:      req.Features,
		IsActive:      req.IsActive,
	}

	if err := s.pkgRepo.Create(ctx, pkg); err != nil {
		return nil, fmt.Errorf("create package: %w", err)
	}

	return pkg, nil
}

func (s *Service) UpdatePackage(ctx context.Context, pkgID uuid.UUID, req UpdatePackageRequest) (*model.Package, error) {
	pkg, err := s.pkgRepo.GetByID(ctx, pkgID)
	if err != nil || pkg == nil {
		return nil, errors.ErrNotFound
	}

	if req.Name != nil {
		pkg.Name = *req.Name
	}
	if req.Price != nil {
		pkg.Price = *req.Price
	}
	if req.DurationDays != nil {
		pkg.DurationDays = req.DurationDays
	}
	if req.GuestLimit != nil {
		pkg.GuestLimit = req.GuestLimit
	}
	if req.TemplateGroup != nil {
		pkg.TemplateGroup = *req.TemplateGroup
	}
	if req.Features != nil {
		pkg.Features = req.Features
	}
	if req.IsActive != nil {
		pkg.IsActive = *req.IsActive
	}

	if err := s.pkgRepo.Update(ctx, pkg); err != nil {
		return nil, fmt.Errorf("update package: %w", err)
	}

	return pkg, nil
}

func (s *Service) GetTemplates(ctx context.Context) ([]model.Template, error) {
	return s.templateRepo.ListAll(ctx)
}

func (s *Service) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*model.Template, error) {
	t := &model.Template{
		Name:         req.Name,
		GroupName:    req.GroupName,
		ThumbnailURL: req.ThumbnailURL,
		CSSConfig:    req.CSSConfig,
		LayoutConfig: req.LayoutConfig,
		IsActive:     req.IsActive,
	}

	if err := s.templateRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create template: %w", err)
	}

	return t, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, templateID uuid.UUID, req UpdateTemplateRequest) (*model.Template, error) {
	t, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil || t == nil {
		return nil, errors.ErrNotFound
	}

	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.GroupName != nil {
		t.GroupName = *req.GroupName
	}
	if req.ThumbnailURL != nil {
		t.ThumbnailURL = *req.ThumbnailURL
	}
	if req.CSSConfig != nil {
		t.CSSConfig = req.CSSConfig
	}
	if req.LayoutConfig != nil {
		t.LayoutConfig = req.LayoutConfig
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}

	if err := s.templateRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("update template: %w", err)
	}

	return t, nil
}

func (s *Service) GetTransactions(ctx context.Context, page, perPage int, status, packageID string) ([]model.Transaction, int64, error) {
	return s.txnRepo.ListAll(ctx, page, perPage, status, packageID)
}

func (s *Service) GetTransactionDetail(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error) {
	return s.txnRepo.GetByID(ctx, txnID)
}

type Analytics struct {
	TotalUsers          int64               `json:"total_users"`
	ActiveUsers         int64               `json:"active_users"`
	TotalTransactions   int64               `json:"total_transactions"`
	TotalRevenue        int64               `json:"total_revenue"`
	SettlementCount     int                 `json:"settlement_count"`
	ActivePackages      int                 `json:"active_packages"`
	ActiveTemplates     int                 `json:"active_templates"`
	ActiveSubscriptions int64               `json:"active_subscriptions"`
	TotalInvitations    int64               `json:"total_invitations"`
	RecentTransactions  []RecentTransaction `json:"recent_transactions"`
	RecentUsers         []RecentUser        `json:"recent_users"`
	PeriodStart         time.Time           `json:"period_start"`
	PeriodEnd           time.Time           `json:"period_end"`
}

type RecentTransaction struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type RecentUser struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePackageRequest struct {
	Name          string         `json:"name" validate:"required,max=100"`
	Code          string         `json:"code" validate:"required,max=50"`
	Price         int64          `json:"price" validate:"required,gte=0"`
	DurationDays  *int           `json:"duration_days"`
	GuestLimit    *int           `json:"guest_limit"`
	TemplateGroup string         `json:"template_group" validate:"required,oneof=standard premium all"`
	Features      datatypes.JSON `json:"features"`
	IsActive      bool           `json:"is_active"`
}

type UpdatePackageRequest struct {
	Name          *string        `json:"name,omitempty"`
	Price         *int64         `json:"price,omitempty"`
	DurationDays  *int           `json:"duration_days"`
	GuestLimit    *int           `json:"guest_limit"`
	TemplateGroup *string        `json:"template_group,omitempty"`
	Features      datatypes.JSON `json:"features,omitempty"`
	IsActive      *bool          `json:"is_active,omitempty"`
}

type CreateTemplateRequest struct {
	Name         string         `json:"name" validate:"required,max=100"`
	GroupName    string         `json:"group_name" validate:"required"`
	ThumbnailURL string         `json:"thumbnail_url,omitempty"`
	CSSConfig    datatypes.JSON `json:"css_config,omitempty"`
	LayoutConfig datatypes.JSON `json:"layout_config,omitempty"`
	IsActive     bool           `json:"is_active"`
}

type UpdateTemplateRequest struct {
	Name         *string        `json:"name,omitempty"`
	GroupName    *string        `json:"group_name,omitempty"`
	ThumbnailURL *string        `json:"thumbnail_url,omitempty"`
	CSSConfig    datatypes.JSON `json:"css_config,omitempty"`
	LayoutConfig datatypes.JSON `json:"layout_config,omitempty"`
	IsActive     *bool          `json:"is_active,omitempty"`
}
