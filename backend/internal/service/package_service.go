package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"gorm.io/datatypes"
)

type PackageService struct {
	pkgRepo   repository.PackageRepository
	auditRepo repository.AuditLogRepository
}

func NewPackageService(pkgRepo repository.PackageRepository, auditRepo repository.AuditLogRepository) *PackageService {
	return &PackageService{pkgRepo: pkgRepo, auditRepo: auditRepo}
}

func (s *PackageService) ListActive(ctx context.Context) ([]dto.PackageResponse, error) {
	packages, err := s.pkgRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PackageResponse, len(packages))
	for i, pkg := range packages {
		result[i] = toPackageResponse(&pkg)
	}
	return result, nil
}

func (s *PackageService) ListAllForAdmin(ctx context.Context) ([]dto.PackageResponse, error) {
	packages, err := s.pkgRepo.GetAllWithInactive(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PackageResponse, len(packages))
	for i, pkg := range packages {
		result[i] = toPackageResponse(&pkg)
	}
	return result, nil
}

func (s *PackageService) Create(ctx context.Context, req dto.CreatePackageRequest) (*dto.PackageResponse, error) {
	_, err := s.pkgRepo.GetByCode(ctx, req.Code)
	if err == nil {
		return nil, fmt.Errorf("%w: package code already exists", errors.ErrConflict)
	}
	pkg := &model.Package{
		Name:          req.Name,
		Code:          req.Code,
		Price:         req.Price,
		DurationDays:  req.DurationDays,
		GuestLimit:    req.GuestLimit,
		TemplateGroup: req.TemplateGroup,
		Features:      datatypes.JSON(req.Features),
		IsActive:      req.IsActive,
	}
	if err := s.pkgRepo.Create(ctx, pkg); err != nil {
		return nil, fmt.Errorf("create package: %w", err)
	}
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		Action:     "package.created",
		EntityType: "package",
		EntityID:   &pkg.ID,
	})
	resp := toPackageResponse(pkg)
	return &resp, nil
}

func (s *PackageService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PackageResponse, error) {
	pkg, err := s.pkgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	resp := toPackageResponse(pkg)
	return &resp, nil
}

func (s *PackageService) Update(ctx context.Context, id uuid.UUID, req dto.UpdatePackageRequest) (*dto.PackageResponse, error) {
	pkg, err := s.pkgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if req.Name != nil {
		pkg.Name = *req.Name
	}
	if req.Code != nil {
		pkg.Code = *req.Code
	}
	if req.Price != nil {
		pkg.Price = *req.Price
	}
	pkg.DurationDays = req.DurationDays
	pkg.GuestLimit = req.GuestLimit
	if req.TemplateGroup != nil {
		pkg.TemplateGroup = *req.TemplateGroup
	}
	if req.Features != nil {
		if len(pkg.Features) == 0 {
			pkg.Features = datatypes.JSON{}
		}
	}
	if req.IsActive != nil {
		pkg.IsActive = *req.IsActive
	}
	if err := s.pkgRepo.Update(ctx, pkg); err != nil {
		return nil, fmt.Errorf("update package: %w", err)
	}
	resp := toPackageResponse(pkg)
	return &resp, nil
}

func (s *PackageService) Deactivate(ctx context.Context, id uuid.UUID) error {
	pkg, err := s.pkgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrNotFound
	}
	if pkg.Code == "free" {
		return fmt.Errorf("%w: free package cannot be deactivated", errors.ErrConflict)
	}
	return s.pkgRepo.Deactivate(ctx, id)
}

func (s *PackageService) GetFeatures(ctx context.Context, id uuid.UUID) (map[string]interface{}, error) {
	pkg, err := s.pkgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if pkg.Features == nil {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	if err := pkg.Features.UnmarshalJSON(pkg.Features); err == nil && m != nil {
		return m, nil
	}
	return map[string]interface{}{}, nil
}

func toPackageResponse(pkg *model.Package) dto.PackageResponse {
	return dto.PackageResponse{
		ID:            pkg.ID,
		Name:          pkg.Name,
		Code:          pkg.Code,
		Price:         pkg.Price,
		DurationDays:  pkg.DurationDays,
		GuestLimit:    pkg.GuestLimit,
		TemplateGroup: pkg.TemplateGroup,
		Features:      pkg.Features,
		IsActive:      pkg.IsActive,
		CreatedAt:     pkg.CreatedAt,
		UpdatedAt:     pkg.UpdatedAt,
	}
}
