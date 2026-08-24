package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountByRole(ctx context.Context, role string) (int64, error)
	List(ctx context.Context, page, perPage int, search, status string) ([]model.User, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	WithTx(tx *gorm.DB) UserRepository
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) WithTx(tx *gorm.DB) UserRepository {
	return &userRepo{db: tx}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = LOWER(?) AND deleted_at IS NULL", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}

func (r *userRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ? AND deleted_at IS NULL", status).Count(&count).Error
	return count, err
}

func (r *userRepo) List(ctx context.Context, page, perPage int, search, status string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.WithContext(ctx).Model(&model.User{}).Where("deleted_at IS NULL")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&users).Error
	return users, total, err
}

func (r *userRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *userRepo) CountByRole(ctx context.Context, role string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ? AND deleted_at IS NULL", role).Count(&count).Error
	return count, err
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND is_revoked = false AND expires_at > CURRENT_TIMESTAMP", tokenHash).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).Where("id = ?", id).Update("is_revoked", true).Error
}

func (r *refreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).Where("user_id = ? AND is_revoked = false", userID).Update("is_revoked", true).Error
}

func (r *refreshTokenRepo) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < CURRENT_TIMESTAMP OR is_revoked = true").Delete(&model.RefreshToken{}).Error
}

type PackageRepository interface {
	Create(ctx context.Context, pkg *model.Package) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Package, error)
	GetByCode(ctx context.Context, code string) (*model.Package, error)
	GetAllActive(ctx context.Context) ([]model.Package, error)
	GetAllWithInactive(ctx context.Context) ([]model.Package, error)
	Update(ctx context.Context, pkg *model.Package) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	WithTx(tx *gorm.DB) PackageRepository
}

type packageRepo struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepo{db: db}
}

func (r *packageRepo) WithTx(tx *gorm.DB) PackageRepository {
	return &packageRepo{db: tx}
}

func (r *packageRepo) Create(ctx context.Context, pkg *model.Package) error {
	return r.db.WithContext(ctx).Create(pkg).Error
}

func (r *packageRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	var pkg model.Package
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepo) GetByCode(ctx context.Context, code string) (*model.Package, error) {
	var pkg model.Package
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepo) GetAllActive(ctx context.Context) ([]model.Package, error) {
	var packages []model.Package
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&packages).Error
	return packages, err
}

func (r *packageRepo) GetAllWithInactive(ctx context.Context) ([]model.Package, error) {
	var packages []model.Package
	err := r.db.WithContext(ctx).Find(&packages).Error
	return packages, err
}

func (r *packageRepo) Update(ctx context.Context, pkg *model.Package) error {
	return r.db.WithContext(ctx).Save(pkg).Error
}

func (r *packageRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Package{}).Where("id = ?", id).Update("is_active", false).Error
}
