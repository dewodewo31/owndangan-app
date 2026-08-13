package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"gorm.io/gorm"
)

type GuestRepository interface {
	Create(ctx context.Context, guest *model.Guest) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Guest, error)
	GetByToken(ctx context.Context, token string) (*model.Guest, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID, page, perPage int) ([]model.Guest, int64, error)
	Update(ctx context.Context, guest *model.Guest) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	CountByEvent(ctx context.Context, eventID uuid.UUID) (int64, error)
	BulkCreate(ctx context.Context, guests []model.Guest) error
	IsTokenTaken(ctx context.Context, token string) (bool, error)
	WithTx(tx *gorm.DB) GuestRepository
}

type guestRepo struct {
	db *gorm.DB
}

func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepo{db: db}
}

func (r *guestRepo) WithTx(tx *gorm.DB) GuestRepository {
	return &guestRepo{db: tx}
}

func (r *guestRepo) Create(ctx context.Context, guest *model.Guest) error {
	return r.db.WithContext(ctx).Create(guest).Error
}

func (r *guestRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Guest, error) {
	var guest model.Guest
	err := r.db.WithContext(ctx).
		Preload("RSVP").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&guest).Error
	if err != nil {
		return nil, err
	}
	return &guest, nil
}

func (r *guestRepo) GetByToken(ctx context.Context, token string) (*model.Guest, error) {
	var guest model.Guest
	err := r.db.WithContext(ctx).
		Preload("RSVP").
		Where("token = ? AND deleted_at IS NULL", token).
		First(&guest).Error
	if err != nil {
		return nil, err
	}
	return &guest, nil
}

func (r *guestRepo) ListByEvent(ctx context.Context, eventID uuid.UUID, page, perPage int) ([]model.Guest, int64, error) {
	var guests []model.Guest
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Guest{}).
		Where("event_id = ? AND deleted_at IS NULL", eventID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := query.Preload("RSVP").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&guests).Error
	return guests, total, err
}

func (r *guestRepo) Update(ctx context.Context, guest *model.Guest) error {
	return r.db.WithContext(ctx).Save(guest).Error
}

func (r *guestRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Guest{}, id).Error
}

func (r *guestRepo) CountByEvent(ctx context.Context, eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Guest{}).
		Where("event_id = ? AND deleted_at IS NULL", eventID).
		Count(&count).Error
	return count, err
}

func (r *guestRepo) BulkCreate(ctx context.Context, guests []model.Guest) error {
	return r.db.WithContext(ctx).Create(&guests).Error
}

func (r *guestRepo) IsTokenTaken(ctx context.Context, token string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Guest{}).
		Where("token = ? AND deleted_at IS NULL", token).
		Count(&count).Error
	return count > 0, err
}

type RSVPRepository interface {
	Create(ctx context.Context, rsvp *model.RSVP) error
	Update(ctx context.Context, rsvp *model.RSVP) error
	GetByGuestID(ctx context.Context, guestID uuid.UUID) (*model.RSVP, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.RSVP, error)
	CountRespondedByEvent(ctx context.Context, eventID uuid.UUID) (int64, error)
	CountByAttendance(ctx context.Context, eventID uuid.UUID, attendance string) (int64, error)
	SumGuestCountByAttendance(ctx context.Context, eventID uuid.UUID, attendance string) (int64, error)
	WithTx(tx *gorm.DB) RSVPRepository
}

type rsvpRepo struct {
	db *gorm.DB
}

func NewRSVPRepository(db *gorm.DB) RSVPRepository {
	return &rsvpRepo{db: db}
}

func (r *rsvpRepo) WithTx(tx *gorm.DB) RSVPRepository {
	return &rsvpRepo{db: tx}
}

func (r *rsvpRepo) Create(ctx context.Context, rsvp *model.RSVP) error {
	return r.db.WithContext(ctx).Create(rsvp).Error
}

func (r *rsvpRepo) Update(ctx context.Context, rsvp *model.RSVP) error {
	return r.db.WithContext(ctx).Save(rsvp).Error
}

func (r *rsvpRepo) GetByGuestID(ctx context.Context, guestID uuid.UUID) (*model.RSVP, error) {
	var rsvp model.RSVP
	err := r.db.WithContext(ctx).Where("guest_id = ?", guestID).First(&rsvp).Error
	if err != nil {
		return nil, err
	}
	return &rsvp, nil
}

func (r *rsvpRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.RSVP, error) {
	var rsvps []model.RSVP
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&rsvps).Error
	return rsvps, err
}

func (r *rsvpRepo) CountRespondedByEvent(ctx context.Context, eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RSVP{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}

func (r *rsvpRepo) CountByAttendance(ctx context.Context, eventID uuid.UUID, attendance string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RSVP{}).
		Where("event_id = ? AND attendance = ?", eventID, attendance).
		Count(&count).Error
	return count, err
}

func (r *rsvpRepo) SumGuestCountByAttendance(ctx context.Context, eventID uuid.UUID, attendance string) (int64, error) {
	var sum int64
	err := r.db.WithContext(ctx).Model(&model.RSVP{}).
		Where("event_id = ? AND attendance = ?", eventID, attendance).
		Select("COALESCE(SUM(guest_count), 0)").
		Row().Scan(&sum)
	return sum, err
}

type GuestbookRepository interface {
	Create(ctx context.Context, msg *model.GuestbookMessage) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.GuestbookMessage, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID, approved bool) ([]model.GuestbookMessage, error)
	Approve(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	WithTx(tx *gorm.DB) GuestbookRepository
}

type guestbookRepo struct {
	db *gorm.DB
}

func NewGuestbookRepository(db *gorm.DB) GuestbookRepository {
	return &guestbookRepo{db: db}
}

func (r *guestbookRepo) WithTx(tx *gorm.DB) GuestbookRepository {
	return &guestbookRepo{db: tx}
}

func (r *guestbookRepo) Create(ctx context.Context, msg *model.GuestbookMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *guestbookRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.GuestbookMessage, error) {
	var msg model.GuestbookMessage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *guestbookRepo) ListByEvent(ctx context.Context, eventID uuid.UUID, approved bool) ([]model.GuestbookMessage, error) {
	var msgs []model.GuestbookMessage
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND is_approved = ?", eventID, approved).
		Order("created_at DESC").
		Find(&msgs).Error
	return msgs, err
}

func (r *guestbookRepo) Approve(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.GuestbookMessage{}).
		Where("id = ?", id).
		Update("is_approved", true).Error
}

func (r *guestbookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.GuestbookMessage{}, id).Error
}

type DigitalGiftRepository interface {
	Create(ctx context.Context, gift *model.DigitalGift) error
	GetByEventID(ctx context.Context, eventID uuid.UUID) (*model.DigitalGift, error)
	Update(ctx context.Context, gift *model.DigitalGift) error
	WithTx(tx *gorm.DB) DigitalGiftRepository
}

type digitalGiftRepo struct {
	db *gorm.DB
}

func NewDigitalGiftRepository(db *gorm.DB) DigitalGiftRepository {
	return &digitalGiftRepo{db: db}
}

func (r *digitalGiftRepo) WithTx(tx *gorm.DB) DigitalGiftRepository {
	return &digitalGiftRepo{db: tx}
}

func (r *digitalGiftRepo) Create(ctx context.Context, gift *model.DigitalGift) error {
	return r.db.WithContext(ctx).Create(gift).Error
}

func (r *digitalGiftRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) (*model.DigitalGift, error) {
	var gift model.DigitalGift
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(&gift).Error
	if err != nil {
		return nil, err
	}
	return &gift, nil
}

func (r *digitalGiftRepo) Update(ctx context.Context, gift *model.DigitalGift) error {
	return r.db.WithContext(ctx).Save(gift).Error
}

type GalleryPhotoRepository interface {
	Create(ctx context.Context, photo *model.GalleryPhoto) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.GalleryPhoto, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.GalleryPhoto, error)
	CountByEvent(ctx context.Context, eventID uuid.UUID) (int64, error)
	UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error
	Delete(ctx context.Context, id uuid.UUID) error
	WithTx(tx *gorm.DB) GalleryPhotoRepository
}

type galleryPhotoRepo struct {
	db *gorm.DB
}

func NewGalleryPhotoRepository(db *gorm.DB) GalleryPhotoRepository {
	return &galleryPhotoRepo{db: db}
}

func (r *galleryPhotoRepo) WithTx(tx *gorm.DB) GalleryPhotoRepository {
	return &galleryPhotoRepo{db: tx}
}

func (r *galleryPhotoRepo) Create(ctx context.Context, photo *model.GalleryPhoto) error {
	return r.db.WithContext(ctx).Create(photo).Error
}

func (r *galleryPhotoRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.GalleryPhoto, error) {
	var photo model.GalleryPhoto
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&photo).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *galleryPhotoRepo) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.GalleryPhoto, error) {
	var photos []model.GalleryPhoto
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("sort_order ASC, created_at ASC").
		Find(&photos).Error
	return photos, err
}

func (r *galleryPhotoRepo) UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error {
	return r.db.WithContext(ctx).Model(&model.GalleryPhoto{}).
		Where("id = ?", id).Update("sort_order", sortOrder).Error
}

func (r *galleryPhotoRepo) CountByEvent(ctx context.Context, eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.GalleryPhoto{}).
		Where("event_id = ?", eventID).
		Count(&count).Error
	return count, err
}

func (r *galleryPhotoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.GalleryPhoto{}, id).Error
}

type AnalyticsEventRepository interface {
	Create(ctx context.Context, event *model.AnalyticsEvent) error
	CountByType(ctx context.Context, eventType string, since time.Time) (int64, error)
	SumRevenueLast30Days(ctx context.Context) (int64, error)
	CountEventsByTypeLast30Days(ctx context.Context) (map[string]int64, error)
}

type analyticsEventRepo struct {
	db *gorm.DB
}

func NewAnalyticsEventRepository(db *gorm.DB) AnalyticsEventRepository {
	return &analyticsEventRepo{db: db}
}

func (r *analyticsEventRepo) Create(ctx context.Context, event *model.AnalyticsEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *analyticsEventRepo) CountByType(ctx context.Context, eventType string, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AnalyticsEvent{}).
		Where("event_type = ? AND created_at > ?", eventType, since).
		Count(&count).Error
	return count, err
}

func (r *analyticsEventRepo) SumRevenueLast30Days(ctx context.Context) (int64, error) {
	var sum int64
	err := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("status = 'settlement' AND created_at > ?", time.Now().AddDate(0, 0, -30)).
		Select("COALESCE(SUM(gross_amount), 0)").
		Row().Scan(&sum)
	return sum, err
}

func (r *analyticsEventRepo) CountEventsByTypeLast30Days(ctx context.Context) (map[string]int64, error) {
	type result struct {
		EventType string
		Count     int64
	}
	var results []result
	err := r.db.WithContext(ctx).Model(&model.AnalyticsEvent{}).
		Select("event_type, COUNT(*) as count").
		Where("created_at > ?", time.Now().AddDate(0, 0, -30)).
		Group("event_type").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64)
	for _, r := range results {
		m[r.EventType] = r.Count
	}
	return m, nil
}
