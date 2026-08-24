package music

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
)

type Service struct {
	musicRepo repository.MusicRepository
	eventRepo repository.EventRepository
	storage   Storage
	validator MediaValidator
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

func NewService(musicRepo repository.MusicRepository, eventRepo repository.EventRepository,
	storage Storage, validator MediaValidator) *Service {
	return &Service{
		musicRepo: musicRepo,
		eventRepo: eventRepo,
		storage:   storage,
		validator: validator,
	}
}

func (s *Service) SetEventMusic(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req SetMusicRequest) (*model.Music, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if req.IsPreset {
		return s.setPresetMusic(ctx, eventID, req.PresetName)
	}

	if err := s.validator.ValidateFile(req.Filename, req.FileSize, req.ContentType); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidInput, err)
	}

	safeName := s.validator.SanitizeFilename(req.Filename)
	key := fmt.Sprintf("music/%s/%s", eventID.String(), safeName)

	data, _ := io.ReadAll(req.Data)
	result, err := s.storage.Upload(ctx, key, data, UploadOpts{
		ContentType: req.ContentType,
		Extension:   req.Extension,
		MaxSize:     10 * 1024 * 1024,
	})
	if err != nil {
		return nil, fmt.Errorf("upload music: %w", err)
	}

	music := &model.Music{
		EventID:  &eventID,
		Title:    req.Title,
		FileURL:  result.URL,
		IsPreset: false,
	}

	if err := s.musicRepo.Create(ctx, music); err != nil {
		_ = s.storage.Delete(ctx, result.Key)
		return nil, fmt.Errorf("create music record: %w", err)
	}

	return music, nil
}

func (s *Service) setPresetMusic(ctx context.Context, eventID uuid.UUID, presetName string) (*model.Music, error) {
	presets, err := s.musicRepo.ListPresets(ctx)
	if err != nil {
		return nil, fmt.Errorf("list presets: %w", err)
	}

	for _, p := range presets {
		if p.Preset == presetName || p.Title == presetName {
			music := &model.Music{
				EventID:  &eventID,
				Title:    p.Title,
				FileURL:  p.FileURL,
				Preset:   p.Preset,
				IsPreset: true,
			}
			if err := s.musicRepo.Create(ctx, music); err != nil {
				return nil, fmt.Errorf("create music record: %w", err)
			}
			return music, nil
		}
	}

	return nil, fmt.Errorf("%w: preset not found", errors.ErrNotFound)
}

func (s *Service) GetEventMusic(ctx context.Context, eventID uuid.UUID) (*model.Music, error) {
	return s.musicRepo.GetByEvent(ctx, eventID)
}

func (s *Service) ListPresets(ctx context.Context) ([]model.Music, error) {
	return s.musicRepo.ListPresets(ctx)
}
