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
	userRepo      repository.UserRepository
	pkgRepo       repository.PackageRepository
	txnRepo       repository.TransactionRepository
	templateRepo  repository.TemplateRepository
	auditRepo     repository.AuditLogRepository
}

func NewService(userRepo repository.UserRepository, pkgRepo repository.PackageRepository,
	txnRepo repository.TransactionRepository, templateRepo repository.TemplateRepository,
	auditRepo repository.AuditLogRepository) *Service {
	return &Service{
		userRepo:     userRepo,
		pkgRepo:      pkgRepo,
		txnRepo:      txnRepo,
		templateRepo: templateRepo,
		auditRepo:    auditRepo,
	}
}

func (s *Service) GetAnalytics(ctx context.Context) (*Analytics, error) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	users, _, _ := s.userRepo.List(ctx, 1, 10000)
	activeUsers := int64(0)
	for _, u := range users {
		if u.Status == "active" {
			activeUsers++
		}
	}

	txns, totalTxns, _ := s.txnRepo.ListAll(ctx, 1, 10000, "")
	var totalRevenue int64
	settlementCount := 0
	for _, t := range txns {
		if t.Status == "settlement" {
			totalRevenue += t.GrossAmount
			settlementCount++
		}
	}

	pkgs, _ := s.pkgRepo.GetAllActive(ctx)
	templates, _ := s.templateRepo.ListAll(ctx)

	return &Analytics{
		TotalUsers:       int64(len(users)),
		ActiveUsers:      activeUsers,
		TotalTransactions: totalTxns,
		TotalRevenue:     totalRevenue,
		SettlementCount:  settlementCount,
		ActivePackages:   len(pkgs),
		ActiveTemplates:  len(templates),
		PeriodStart:      thirtyDaysAgo,
		PeriodEnd:        now,
	}, nil
}

func (s *Service) ListUsers(ctx context.Context, page, perPage int, status string) ([]model.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	if status != "" {
		var filtered []model.User
		for _, u := range users {
			if u.Status == status {
				filtered = append(filtered, u)
			}
		}
		return filtered, int64(len(filtered)), nil
	}

	return users, total, nil
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

func (s *Service) GetTransactions(ctx context.Context, page, perPage int, status string) ([]model.Transaction, int64, error) {
	return s.txnRepo.ListAll(ctx, page, perPage, status)
}

func (s *Service) GetTransactionDetail(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error) {
	return s.txnRepo.GetByID(ctx, txnID)
}

type Analytics struct {
	TotalUsers       int64     `json:"total_users"`
	ActiveUsers      int64     `json:"active_users"`
	TotalTransactions int64    `json:"total_transactions"`
	TotalRevenue     int64     `json:"total_revenue"`
	SettlementCount  int       `json:"settlement_count"`
	ActivePackages   int       `json:"active_packages"`
	ActiveTemplates  int       `json:"active_templates"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
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
