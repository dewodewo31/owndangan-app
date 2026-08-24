package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	GetBySlug(ctx context.Context, slug string) (*model.Event, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, perPage int, status string) ([]model.Event, int64, error)
	Update(ctx context.Context, event *model.Event) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	IncrementViewCount(ctx context.Context, slug string) error
	WithTx(tx *gorm.DB) EventRepository
}

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) WithTx(tx *gorm.DB) EventRepository {
	return &eventRepo{db: tx}
}

func (r *eventRepo) Create(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).
		Preload("Sections").
		Preload("DigitalGift").
		Preload("LoveStories").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepo) GetBySlug(ctx context.Context, slug string) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).
		Preload("Sections").
		Preload("DigitalGift").
		Preload("GalleryPhotos").
		Preload("LoveStories").
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepo) ListByUser(ctx context.Context, userID uuid.UUID, page, perPage int, status string) ([]model.Event, int64, error) {
	var events []model.Event
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Event{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&events).Error
	return events, total, err
}

func (r *eventRepo) Update(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *eventRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Event{}, id).Error
}

func (r *eventRepo) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Event{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count).Error
	return count, err
}

func (r *eventRepo) IncrementViewCount(ctx context.Context, slug string) error {
	return r.db.WithContext(ctx).Model(&model.Event{}).Where("slug = ? AND deleted_at IS NULL", slug).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

type EventSectionRepository interface {
	Create(ctx context.Context, section *model.EventSection) error
	GetByEventID(ctx context.Context, eventID uuid.UUID) (*model.EventSection, error)
	Update(ctx context.Context, section *model.EventSection) error
	WithTx(tx *gorm.DB) EventSectionRepository
}

type eventSectionRepo struct {
	db *gorm.DB
}

func NewEventSectionRepository(db *gorm.DB) EventSectionRepository {
	return &eventSectionRepo{db: db}
}

func (r *eventSectionRepo) WithTx(tx *gorm.DB) EventSectionRepository {
	return &eventSectionRepo{db: tx}
}

func (r *eventSectionRepo) Create(ctx context.Context, section *model.EventSection) error {
	return r.db.WithContext(ctx).Create(section).Error
}

func (r *eventSectionRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) (*model.EventSection, error) {
	var section model.EventSection
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *eventSectionRepo) Update(ctx context.Context, section *model.EventSection) error {
	return r.db.WithContext(ctx).Save(section).Error
}

type TemplateRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Template, error)
	ListByGroups(ctx context.Context, groups []string) ([]model.Template, error)
	ListAll(ctx context.Context) ([]model.Template, error)
	Create(ctx context.Context, t *model.Template) error
	Update(ctx context.Context, t *model.Template) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type templateRepo struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepo{db: db}
}

func (r *templateRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Template, error) {
	var t model.Template
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *templateRepo) ListByGroups(ctx context.Context, groups []string) ([]model.Template, error) {
	var templates []model.Template
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND group_name IN ?", true, groups).
		Find(&templates).Error
	return templates, err
}

func (r *templateRepo) ListAll(ctx context.Context) ([]model.Template, error) {
	var templates []model.Template
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *templateRepo) Create(ctx context.Context, t *model.Template) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *templateRepo) Update(ctx context.Context, t *model.Template) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *templateRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Template{}).Where("id = ?", id).Update("is_active", false).Error
}

type MusicRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Music, error)
	ListPresets(ctx context.Context) ([]model.Music, error)
	Create(ctx context.Context, m *model.Music) error
	GetByEvent(ctx context.Context, eventID uuid.UUID) (*model.Music, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type LoveStoryRepository interface {
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.LoveStory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.LoveStory, error)
	Create(ctx context.Context, story *model.LoveStory) error
	Update(ctx context.Context, story *model.LoveStory) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByEvent(ctx context.Context, eventID uuid.UUID) (int64, error)
	WithTx(tx *gorm.DB) LoveStoryRepository
}

type loveStoryRepo struct {
	db *gorm.DB
}

func NewLoveStoryRepository(db *gorm.DB) LoveStoryRepository {
	return &loveStoryRepo{db: db}
}

func (r *loveStoryRepo) WithTx(tx *gorm.DB) LoveStoryRepository {
	return &loveStoryRepo{db: tx}
}

func (r *loveStoryRepo) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.LoveStory, error) {
	var stories []model.LoveStory
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("sort_order ASC, created_at ASC").
		Find(&stories).Error
	return stories, err
}

func (r *loveStoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.LoveStory, error) {
	var story model.LoveStory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&story).Error
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func (r *loveStoryRepo) Create(ctx context.Context, story *model.LoveStory) error {
	return r.db.WithContext(ctx).Create(story).Error
}

func (r *loveStoryRepo) Update(ctx context.Context, story *model.LoveStory) error {
	return r.db.WithContext(ctx).Save(story).Error
}

func (r *loveStoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.LoveStory{}, id).Error
}

func (r *loveStoryRepo) CountByEvent(ctx context.Context, eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LoveStory{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}

type musicRepo struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) MusicRepository {
	return &musicRepo{db: db}
}

func (r *musicRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Music, error) {
	var m model.Music
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *musicRepo) ListPresets(ctx context.Context) ([]model.Music, error) {
	var music []model.Music
	err := r.db.WithContext(ctx).Where("is_preset = ? AND event_id IS NULL", true).Find(&music).Error
	return music, err
}

func (r *musicRepo) Create(ctx context.Context, m *model.Music) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *musicRepo) GetByEvent(ctx context.Context, eventID uuid.UUID) (*model.Music, error) {
	var m model.Music
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *musicRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Music{}).Error
}
