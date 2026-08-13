package gallery

import (
	"context"
	"io"
	"fmt"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
)

type Service struct {
	galleryRepo repository.GalleryPhotoRepository
	eventRepo   repository.EventRepository
	storage     Storage
	validator   MediaValidator
}

type Storage interface {
	Upload(ctx context.Context, key string, data []byte, opts UploadOpts) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) string
}

type UploadOpts struct {
	ContentType string
	Extension   string
	MaxSize     int64
}

type UploadResult struct {
	URL string
	Key string
}

type MediaValidator interface {
	ValidateFile(filename string, size int64, contentType string) error
	SanitizeFilename(name string) string
}

func NewService(galleryRepo repository.GalleryPhotoRepository, eventRepo repository.EventRepository,
	storage Storage, validator MediaValidator) *Service {
	return &Service{
		galleryRepo: galleryRepo,
		eventRepo:   eventRepo,
		storage:     storage,
		validator:   validator,
	}
}

func (s *Service) AddPhoto(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req AddPhotoRequest) (*model.GalleryPhoto, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if err := s.validator.ValidateFile(req.Filename, req.FileSize, req.ContentType); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidInput, err)
	}

	safeName := s.validator.SanitizeFilename(req.Filename)
	key := fmt.Sprintf("gallery/%s/%s", eventID.String(), safeName)

	data, _ := io.ReadAll(req.Data)
	result, err := s.storage.Upload(ctx, key, data, UploadOpts{
		ContentType: req.ContentType,
		Extension:   req.Extension,
		MaxSize:     5 * 1024 * 1024,
	})
	if err != nil {
		return nil, fmt.Errorf("upload photo: %w", err)
	}

	count, _ := s.galleryRepo.CountByEvent(ctx, eventID)

	photo := &model.GalleryPhoto{
		EventID:   eventID,
		ImageURL:  result.URL,
		Caption:   req.Caption,
		SortOrder: int(count),
	}

	if err := s.galleryRepo.Create(ctx, photo); err != nil {
		_ = s.storage.Delete(ctx, result.Key)
		return nil, fmt.Errorf("create photo record: %w", err)
	}

	return photo, nil
}

func (s *Service) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.GalleryPhoto, error) {
	return s.galleryRepo.ListByEvent(ctx, eventID)
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error {
	photo, err := s.galleryRepo.GetByID(ctx, photoID)
	if err != nil || photo == nil {
		return errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, photo.EventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.galleryRepo.Delete(ctx, photoID); err != nil {
		return fmt.Errorf("delete photo: %w", err)
	}

	return nil
}
