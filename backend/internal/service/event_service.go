package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/pkg/pagination"
	"github.com/owndangan/backend/internal/pkg/slug"
	"github.com/owndangan/backend/internal/pkg/storage"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/entitlement"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EventService struct {
	db               *gorm.DB
	eventRepo        repository.EventRepository
	sectionRepo      repository.EventSectionRepository
	digitalRepo      repository.DigitalGiftRepository
	subRepo          repository.SubscriptionRepository
	pkgRepo          repository.PackageRepository
	guestRepo        repository.GuestRepository
	rsvpRepo         repository.RSVPRepository
	guestbookRepo    repository.GuestbookRepository
	loveStoryRepo    repository.LoveStoryRepository
	templateRepo     repository.TemplateRepository
	musicRepo        repository.MusicRepository
	galleryPhotoRepo repository.GalleryPhotoRepository
	auditRepo        repository.AuditLogRepository
	analyticsRepo    repository.AnalyticsEventRepository
	storage          Storage
}

type Storage interface {
	Upload(ctx context.Context, key string, data io.Reader, opts storage.UploadOptions) (*storage.UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) string
}

func NewEventService(db *gorm.DB, eventRepo repository.EventRepository, sectionRepo repository.EventSectionRepository,
	digitalRepo repository.DigitalGiftRepository, subRepo repository.SubscriptionRepository,
	pkgRepo repository.PackageRepository, guestRepo repository.GuestRepository,
	rsvpRepo repository.RSVPRepository, guestbookRepo repository.GuestbookRepository,
	loveStoryRepo repository.LoveStoryRepository,
	templateRepo repository.TemplateRepository, musicRepo repository.MusicRepository,
	galleryPhotoRepo repository.GalleryPhotoRepository,
	auditRepo repository.AuditLogRepository, analyticsRepo repository.AnalyticsEventRepository,
	storage Storage) *EventService {
	return &EventService{
		db:               db,
		eventRepo:        eventRepo,
		sectionRepo:      sectionRepo,
		digitalRepo:      digitalRepo,
		subRepo:          subRepo,
		pkgRepo:          pkgRepo,
		guestRepo:        guestRepo,
		rsvpRepo:         rsvpRepo,
		guestbookRepo:    guestbookRepo,
		loveStoryRepo:    loveStoryRepo,
		templateRepo:     templateRepo,
		musicRepo:        musicRepo,
		galleryPhotoRepo: galleryPhotoRepo,
		auditRepo:        auditRepo,
		analyticsRepo:    analyticsRepo,
		storage:          storage,
	}
}

func (s *EventService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	if err := s.checkEventLimit(ctx, userID); err != nil {
		return nil, err
	}

	slugStr := req.Slug
	if slugStr == "" {
		slugStr = slug.Generate(req.CoupleName)
	}
	if slugStr == "" {
		slugStr = slug.Generate(req.Title)
	}

	if err := slug.Validate(slugStr); err != nil {
		slugStr = slug.Generate(req.Title + "-" + uuid.New().String()[:4])
	}

	slugStr = s.ensureSlugUnique(ctx, slugStr)

	var weddingDate *time.Time
	if req.WeddingDate != "" {
		t, err := time.Parse("2006-01-02", req.WeddingDate)
		if err == nil {
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			weddingDate = &t
		}
	}

	event := &model.Event{
		UserID:           userID,
		Title:            req.Title,
		Slug:             slugStr,
		CoupleName:       req.CoupleName,
		GroomName:        req.GroomName,
		BrideName:        req.BrideName,
		GroomParents:     req.GroomParents,
		BrideParents:     req.BrideParents,
		WeddingDate:      weddingDate,
		WeddingTime:      req.WeddingTime,
		CeremonyVenue:    req.CeremonyVenue,
		CeremonyAddress:  req.CeremonyAddress,
		CeremonyMapURL:   req.CeremonyMapURL,
		ReceptionVenue:   req.ReceptionVenue,
		ReceptionAddress: req.ReceptionAddress,
		ReceptionMapURL:  req.ReceptionMapURL,
		Status:           "draft",
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		eventRepo := s.eventRepo.WithTx(tx)
		sectionRepo := s.sectionRepo.WithTx(tx)
		digitalRepo := s.digitalRepo.WithTx(tx)

		if err := eventRepo.Create(ctx, event); err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				return fmt.Errorf("%w: slug already taken", errors.ErrConflict)
			}
			return err
		}

		section := &model.EventSection{
			EventID:             event.ID,
			HeroEnabled:         true,
			CoupleEnabled:       true,
			EventDetailsEnabled: true,
			GalleryEnabled:      true,
			RSVPEnabled:         true,
			GuestbookEnabled:    true,
			DigitalGiftsEnabled: false,
		}
		if err := sectionRepo.Create(ctx, section); err != nil {
			return err
		}

		gift := &model.DigitalGift{
			EventID: event.ID,
		}
		if err := digitalRepo.Create(ctx, gift); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &event.UserID,
		Action:     "event.created",
		EntityType: "event",
		EntityID:   &event.ID,
	})

	return s.toResponse(ctx, event), nil
}

func (s *EventService) GetByID(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}
	return s.toResponse(ctx, event), nil
}

func (s *EventService) Update(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.CoupleName != nil {
		event.CoupleName = *req.CoupleName
	}
	if req.GroomName != nil {
		event.GroomName = *req.GroomName
	}
	if req.BrideName != nil {
		event.BrideName = *req.BrideName
	}
	if req.GroomParents != nil {
		event.GroomParents = *req.GroomParents
	}
	if req.BrideParents != nil {
		event.BrideParents = *req.BrideParents
	}
	if req.WeddingDate != nil {
		event.WeddingDate = req.WeddingDate
	}
	if req.WeddingTime != nil {
		event.WeddingTime = *req.WeddingTime
	}
	if req.CeremonyVenue != nil {
		event.CeremonyVenue = *req.CeremonyVenue
	}
	if req.CeremonyAddress != nil {
		event.CeremonyAddress = *req.CeremonyAddress
	}
	if req.CeremonyMapURL != nil {
		event.CeremonyMapURL = *req.CeremonyMapURL
	}
	if req.ReceptionVenue != nil {
		event.ReceptionVenue = *req.ReceptionVenue
	}
	if req.ReceptionAddress != nil {
		event.ReceptionAddress = *req.ReceptionAddress
	}
	if req.ReceptionMapURL != nil {
		event.ReceptionMapURL = *req.ReceptionMapURL
	}
	if req.VideoURL != nil {
		event.VideoURL = *req.VideoURL
	}
	if req.TemplateID != nil {
		if err := s.ensureTemplateAssignable(ctx, userID, *req.TemplateID); err != nil {
			return nil, err
		}
		event.TemplateID = req.TemplateID
	}

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	return s.toResponse(ctx, event), nil
}

func (s *EventService) Delete(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	if event.Status == "published" {
		event.Status = "unpublished"
		event.PublishedAt = nil
		if err := s.eventRepo.Update(ctx, event); err != nil {
			return fmt.Errorf("update event: %w", err)
		}
	}

	if err := s.eventRepo.SoftDelete(ctx, eventID); err != nil {
		return err
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &event.UserID,
		Action:     "event.deleted",
		EntityType: "event",
		EntityID:   &event.ID,
	})

	return nil
}

func (s *EventService) Publish(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.PublishResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if err := s.checkEventLimit(ctx, userID); err != nil {
		return nil, err
	}

	if event.Status == "published" {
		return nil, fmt.Errorf("%w: event is already published", errors.ErrConflict)
	}

	if event.GroomName == "" || event.BrideName == "" || event.WeddingDate == nil {
		return nil, fmt.Errorf("%w: complete groom name, bride name, and wedding date before publishing", errors.ErrInvalidInput)
	}

	event.Status = "published"
	now := time.Now()
	event.PublishedAt = &now

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &event.UserID,
		Action:     "event.published",
		EntityType: "event",
		EntityID:   &event.ID,
	})

	return &dto.PublishResponse{
		ID:          event.ID,
		Status:      event.Status,
		PublishedAt: event.PublishedAt.Format(time.RFC3339),
		PublicURL:   fmt.Sprintf("/%s", event.Slug),
	}, nil
}

func (s *EventService) Unpublish(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	if event.Status != "published" {
		return fmt.Errorf("%w: event is not published", errors.ErrConflict)
	}

	event.Status = "unpublished"
	event.PublishedAt = nil

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	return nil
}

func (s *EventService) List(ctx context.Context, userID uuid.UUID, params pagination.Params, status string) ([]dto.EventResponse, int64, error) {
	events, total, err := s.eventRepo.ListByUser(ctx, userID, params.Page, params.PerPage, status)
	if err != nil {
		return nil, 0, err
	}
	result := make([]dto.EventResponse, len(events))
	for i, event := range events {
		result[i] = *s.toResponse(ctx, &event)
	}
	return result, total, nil
}

func (s *EventService) GetPublicBySlug(ctx context.Context, slug string) (*dto.PublicEventResponse, error) {
	event, err := s.eventRepo.GetBySlug(ctx, slug)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.Status != "published" {
		return nil, errors.ErrNotFound
	}

	_ = s.eventRepo.IncrementViewCount(ctx, slug)

	_ = s.analyticsRepo.Create(ctx, &model.AnalyticsEvent{
		EventID:   &event.ID,
		EventType: "page_view",
	})

	section := event.Sections
	if section == nil {
		section, _ = s.sectionRepo.GetByEventID(ctx, event.ID)
	}

	var musicDTO *dto.MusicDTO
	if section != nil && section.MusicID != nil {
		music, err := s.musicRepo.GetByID(ctx, *section.MusicID)
		if err == nil && music != nil {
			musicDTO = &dto.MusicDTO{
				Title:    music.Title,
				FileURL:  music.FileURL,
				Preset:   music.Preset,
				IsPreset: music.IsPreset,
			}
		}
	}

	var gallery []dto.GalleryPhotoDTO
	photos := event.GalleryPhotos
	sort.Slice(photos, func(i, j int) bool {
		return photos[i].SortOrder < photos[j].SortOrder
	})
	for _, p := range photos {
		gallery = append(gallery, dto.GalleryPhotoDTO{
			ImageURL:  p.ImageURL,
			Caption:   p.Caption,
			SortOrder: p.SortOrder,
		})
	}

	var guestbook []dto.GuestbookPublicDTO
	if section != nil && section.GuestbookEnabled {
		gb, _, _ := s.guestbookRepo.ListByEventPaged(ctx, event.ID, true, 50, 0)
		for _, msg := range gb {
			guestbook = append(guestbook, dto.GuestbookPublicDTO{
				Name:      msg.Name,
				Message:   msg.Message,
				CreatedAt: msg.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	var loveStories []dto.LoveStoryPublicDTO
	if section != nil && section.LoveStoryEnabled {
		stories := event.LoveStories
		sort.Slice(stories, func(i, j int) bool {
			return stories[i].SortOrder < stories[j].SortOrder
		})
		for _, st := range stories {
			loveStories = append(loveStories, dto.LoveStoryPublicDTO{
				ID:        st.ID,
				Title:     st.Title,
				Story:     st.Story,
				Year:      st.Year,
				Date:      st.Date,
				ImageURL:  st.ImageURL,
				SortOrder: st.SortOrder,
			})
		}
	}

	var digitalGift *dto.DigitalGiftPublicDTO
	if section != nil && section.DigitalGiftsEnabled && event.DigitalGift != nil {
		digitalGift = &dto.DigitalGiftPublicDTO{
			BankAccounts: decodeJSON(event.DigitalGift.BankAccounts),
			EWallet:      decodeJSONMap(event.DigitalGift.EWallet),
			QRISImageURL: event.DigitalGift.QRISImageURL,
			GiftMessage:  event.DigitalGift.GiftMessage,
		}
	}

	var templatePreview *dto.TemplatePreview
	if event.TemplateID != nil {
		t, err := s.templateRepo.GetByID(ctx, *event.TemplateID)
		if err == nil && t != nil {
			templatePreview = &dto.TemplatePreview{
				Name:         t.Name,
				GroupName:    t.GroupName,
				CSSConfig:    t.CSSConfig,
				LayoutConfig: t.LayoutConfig,
			}
		}
	}

	return &dto.PublicEventResponse{
		Event: dto.PublicEventDetail{
			ID:               event.ID,
			Title:            event.Title,
			CoupleName:       event.CoupleName,
			GroomName:        event.GroomName,
			BrideName:        event.BrideName,
			GroomParents:     event.GroomParents,
			BrideParents:     event.BrideParents,
			WeddingDate:      formatDate(event.WeddingDate),
			WeddingTime:      event.WeddingTime,
			CeremonyVenue:    event.CeremonyVenue,
			CeremonyAddress:  event.CeremonyAddress,
			CeremonyMapURL:   event.CeremonyMapURL,
			ReceptionVenue:   event.ReceptionVenue,
			ReceptionAddress: event.ReceptionAddress,
			ReceptionMapURL:  event.ReceptionMapURL,
			VideoURL:         event.VideoURL,
			ViewCount:        event.ViewCount,
		},
		Template:    templatePreview,
		Sections:    toSectionsDTO(section, musicDTO),
		Gallery:     gallery,
		Guestbook:   guestbook,
		DigitalGift: digitalGift,
		LoveStories: loveStories,
	}, nil
}

func (s *EventService) checkEventLimit(ctx context.Context, userID uuid.UUID) error {
	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil || sub == nil {
		return errors.ErrPaymentRequired
	}

	pkg, err := s.pkgRepo.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return errors.ErrNotFound
	}

	resolver := entitlement.NewResolver(pkg)
	count, _ := s.eventRepo.CountByUser(ctx, userID)
	if !resolver.CanCreateEvent(int(count)) {
		return fmt.Errorf("%w: event limit reached for your plan", errors.ErrLimitExceeded)
	}
	return nil
}

func (s *EventService) resolveEntitlement(ctx context.Context, userID uuid.UUID) *entitlement.Resolver {
	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil || sub == nil {
		freePkg, err := s.pkgRepo.GetByCode(ctx, "free")
		if err != nil || freePkg == nil {
			return entitlement.NewResolver(nil)
		}
		return entitlement.NewResolver(freePkg)
	}
	pkg, err := s.pkgRepo.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return entitlement.NewResolver(nil)
	}
	return entitlement.NewResolver(pkg)
}

func (s *EventService) ensureOwner(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*model.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}
	return event, nil
}

func (s *EventService) ensureTemplateAssignable(ctx context.Context, userID uuid.UUID, templateID uuid.UUID) error {
	t, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil || t == nil {
		return errors.ErrNotFound
	}
	if !t.IsActive {
		return fmt.Errorf("%w: template is not active", errors.ErrConflict)
	}
	resolver := s.resolveEntitlement(ctx, userID)
	switch t.GroupName {
	case "standard":
	case "premium":
		if !resolver.CanAccessPremiumTemplates() {
			return errors.ErrForbidden
		}
	case "all":
		if !resolver.CanAccessAllTemplates() {
			return errors.ErrForbidden
		}
	default:
		return fmt.Errorf("%w: unsupported template group", errors.ErrConflict)
	}
	return nil
}

func (s *EventService) ListTemplates(ctx context.Context, userID uuid.UUID) ([]dto.TemplateSummary, error) {
	resolver := s.resolveEntitlement(ctx, userID)
	groups := []string{"standard"}
	if resolver.CanAccessPremiumTemplates() {
		groups = append(groups, "premium")
	}
	if resolver.CanAccessAllTemplates() {
		groups = append(groups, "all")
	}

	templates, err := s.templateRepo.ListByGroups(ctx, groups)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TemplateSummary, len(templates))
	for i, t := range templates {
		result[i] = dto.TemplateSummary{
			ID:           t.ID,
			Name:         t.Name,
			GroupName:    t.GroupName,
			ThumbnailURL: t.ThumbnailURL,
			CSSConfig:    t.CSSConfig,
			LayoutConfig: t.LayoutConfig,
		}
	}
	return result, nil
}

func (s *EventService) AssignTemplate(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, templateID uuid.UUID) (*dto.EventResponse, error) {
	if err := s.ensureTemplateAssignable(ctx, userID, templateID); err != nil {
		return nil, err
	}
	event, err := s.ensureOwner(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	event.TemplateID = &templateID
	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}
	return s.toResponse(ctx, event), nil
}

func (s *EventService) GetSections(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.SectionsResponse, error) {
	event, err := s.ensureOwner(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	if event.Sections == nil {
		sec, err := s.sectionRepo.GetByEventID(ctx, eventID)
		if err != nil {
			return nil, err
		}
		event.Sections = sec
	}
	sec := event.Sections
	if sec == nil {
		sec = &model.EventSection{EventID: event.ID}
	}
	return &dto.SectionsResponse{
		ID:                  sec.ID,
		EventID:             event.ID,
		HeroEnabled:         sec.HeroEnabled,
		CoupleEnabled:       sec.CoupleEnabled,
		EventDetailsEnabled: sec.EventDetailsEnabled,
		GalleryEnabled:      sec.GalleryEnabled,
		VideoEnabled:        sec.VideoEnabled,
		MusicID:             sec.MusicID,
		RSVPEnabled:         sec.RSVPEnabled,
		GuestbookEnabled:    sec.GuestbookEnabled,
		LoveStoryEnabled:    sec.LoveStoryEnabled,
		DigitalGiftsEnabled: sec.DigitalGiftsEnabled,
		DressCode:           sec.DressCode,
		ClosingMessage:      sec.ClosingMessage,
		OpeningMessage:      sec.OpeningMessage,
		VerseEnabled:        sec.VerseEnabled,
		VerseReligion:       sec.VerseReligion,
		VerseText:           sec.VerseText,
		VerseSource:         sec.VerseSource,
	}, nil
}

func (s *EventService) UpdateSections(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req dto.UpdateSectionsRequest) (*dto.EventResponse, error) {
	event, err := s.ensureOwner(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSectionsEntitlement(ctx, userID, &req); err != nil {
		return nil, err
	}

	sec := event.Sections
	if sec == nil {
		sec, err = s.sectionRepo.GetByEventID(ctx, eventID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sec = &model.EventSection{EventID: eventID}
			} else {
				return nil, err
			}
		}
		event.Sections = sec
	}

	if req.HeroEnabled != nil {
		sec.HeroEnabled = *req.HeroEnabled
	}
	if req.CoupleEnabled != nil {
		sec.CoupleEnabled = *req.CoupleEnabled
	}
	if req.EventDetailsEnabled != nil {
		sec.EventDetailsEnabled = *req.EventDetailsEnabled
	}
	if req.GalleryEnabled != nil {
		sec.GalleryEnabled = *req.GalleryEnabled
	}
	if req.VideoEnabled != nil {
		sec.VideoEnabled = *req.VideoEnabled
	}
	if req.MusicID != nil {
		sec.MusicID = req.MusicID
	}
	if req.RSVPEnabled != nil {
		sec.RSVPEnabled = *req.RSVPEnabled
	}
	if req.GuestbookEnabled != nil {
		sec.GuestbookEnabled = *req.GuestbookEnabled
	}
	if req.LoveStoryEnabled != nil {
		sec.LoveStoryEnabled = *req.LoveStoryEnabled
	}
	if req.DigitalGiftsEnabled != nil {
		sec.DigitalGiftsEnabled = *req.DigitalGiftsEnabled
	}
	if req.DressCode != nil {
		sec.DressCode = *req.DressCode
	}
	if req.ClosingMessage != nil {
		sec.ClosingMessage = *req.ClosingMessage
	}
	if req.OpeningMessage != nil {
		sec.OpeningMessage = *req.OpeningMessage
	}
	if req.VerseEnabled != nil {
		sec.VerseEnabled = *req.VerseEnabled
	}
	if req.VerseReligion != nil {
		sec.VerseReligion = *req.VerseReligion
	}
	if req.VerseText != nil {
		sec.VerseText = *req.VerseText
	}
	if req.VerseSource != nil {
		sec.VerseSource = *req.VerseSource
	}

	if err := s.sectionRepo.Update(ctx, sec); err != nil {
		return nil, err
	}
	return s.toResponse(ctx, event), nil
}

func (s *EventService) ensureSectionsEntitlement(ctx context.Context, userID uuid.UUID, req *dto.UpdateSectionsRequest) error {
	resolver := s.resolveEntitlement(ctx, userID)
	if req.GalleryEnabled != nil && *req.GalleryEnabled && resolver.GalleryMax() != nil {
		// gallery.gallery.max is not a defined feature; gallery count is enforced on upload
	}
	if req.VideoEnabled != nil && *req.VideoEnabled {
		_ = resolver
	}
	if req.MusicID != nil && !resolver.MusicUpload() && !resolver.MusicPreset() {
		return fmt.Errorf("%w: music feature is not available on your plan", errors.ErrLimitExceeded)
	}
	return nil
}

func (s *EventService) GetDigitalGift(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.DigitalGiftResponse, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	gift, err := s.digitalRepo.GetByEventID(ctx, eventID)
	if err != nil || gift == nil {
		return nil, errors.ErrNotFound
	}
	return s.digitalGiftToResponse(gift), nil
}

func (s *EventService) UpdateDigitalGift(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req dto.UpdateDigitalGiftRequest) (*dto.DigitalGiftResponse, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	resolver := s.resolveEntitlement(ctx, userID)
	if req.QRISImageURL != nil && *req.QRISImageURL != "" && !resolver.DigitalGiftQRIS() {
		return nil, fmt.Errorf("%w: QRIS feature is not available on your plan", errors.ErrLimitExceeded)
	}
	gift, err := s.digitalRepo.GetByEventID(ctx, eventID)
	if err != nil || gift == nil {
		return nil, errors.ErrNotFound
	}

	if req.BankAccounts != nil {
		gift.BankAccounts = encodeJSONSlice(req.BankAccounts)
	}
	if req.EWallet != nil {
		gift.EWallet = encodeJSONMap(req.EWallet)
	}
	if req.QRISImageURL != nil {
		gift.QRISImageURL = *req.QRISImageURL
	}
	if req.GiftMessage != nil {
		gift.GiftMessage = *req.GiftMessage
	}

	if err := s.digitalRepo.Update(ctx, gift); err != nil {
		return nil, err
	}
	return s.digitalGiftToResponse(gift), nil
}

func (s *EventService) digitalGiftToResponse(gift *model.DigitalGift) *dto.DigitalGiftResponse {
	return &dto.DigitalGiftResponse{
		ID:           gift.ID,
		EventID:      gift.EventID,
		BankAccounts: decodeJSON(gift.BankAccounts),
		EWallet:      decodeJSONMap(gift.EWallet),
		QRISImageURL: gift.QRISImageURL,
		GiftMessage:  gift.GiftMessage,
	}
}

func (s *EventService) ListGallery(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) ([]dto.GalleryPhotoResponse, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	photos, err := s.galleryPhotoRepo.ListByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	sort.Slice(photos, func(i, j int) bool {
		return photos[i].SortOrder < photos[j].SortOrder
	})
	result := make([]dto.GalleryPhotoResponse, len(photos))
	for i, p := range photos {
		result[i] = dto.GalleryPhotoResponse{
			ID:        p.ID,
			ImageURL:  p.ImageURL,
			Caption:   p.Caption,
			SortOrder: p.SortOrder,
		}
	}
	return result, nil
}

func (s *EventService) UploadGallery(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, file io.Reader, filename, caption string) (*dto.GalleryPhotoResponse, error) {
	event, err := s.ensureOwner(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	resolver := s.resolveEntitlement(ctx, userID)
	if max := resolver.GalleryMax(); max != nil {
		count, _ := s.galleryPhotoRepo.CountByEvent(ctx, eventID)
		if int(count) >= *max {
			return nil, fmt.Errorf("%w: gallery limit reached for your plan", errors.ErrLimitExceeded)
		}
	}

	ext := filepath.Ext(filename)
	key := storage.GenerateKey("gallery/"+event.ID.String(), ext)
	upload, err := s.storage.Upload(ctx, key, file, storage.UploadOptions{
		ContentType: "image/" + strings.TrimPrefix(filepath.Ext(filename), "."),
		Extension:   ext,
	})
	if err != nil {
		return nil, fmt.Errorf("upload gallery: %w", err)
	}

	count, _ := s.galleryPhotoRepo.CountByEvent(ctx, eventID)
	photo := &model.GalleryPhoto{
		EventID:   event.ID,
		ImageURL:  upload.URL,
		Caption:   caption,
		SortOrder: int(count),
	}
	if err := s.galleryPhotoRepo.Create(ctx, photo); err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, err
	}
	return &dto.GalleryPhotoResponse{
		ID:        photo.ID,
		ImageURL:  photo.ImageURL,
		Caption:   photo.Caption,
		SortOrder: photo.SortOrder,
	}, nil
}

func (s *EventService) DeleteGallery(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, photoID uuid.UUID) error {
	event, err := s.ensureOwner(ctx, userID, eventID)
	if err != nil {
		return err
	}
	photo, err := s.galleryPhotoRepo.GetByID(ctx, photoID)
	if err != nil || photo == nil {
		return errors.ErrNotFound
	}
	if photo.EventID != event.ID {
		return errors.ErrForbidden
	}
	if err := s.galleryPhotoRepo.Delete(ctx, photoID); err != nil {
		return err
	}
	_ = s.storage.Delete(ctx, filepath.Base(photo.ImageURL))
	return nil
}

func (s *EventService) ReorderGallery(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req dto.ReorderGalleryRequest) error {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, p := range req.Photos {
			if err := s.galleryPhotoRepo.WithTx(tx).UpdateSortOrder(ctx, p.ID, p.SortOrder); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *EventService) GetMusic(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.MusicResponse, error) {
	event, err := s.ensureOwner(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	if event.Sections == nil {
		sec, _ := s.sectionRepo.GetByEventID(ctx, eventID)
		event.Sections = sec
	}
	if event.Sections == nil || event.Sections.MusicID == nil {
		// check for an event-owned uploaded music as a fallback
		music, err := s.musicRepo.GetByEvent(ctx, eventID)
		if err == nil && music != nil {
			return &dto.MusicResponse{
				ID:       music.ID,
				EventID:  &eventID,
				Title:    music.Title,
				FileURL:  music.FileURL,
				Preset:   music.Preset,
				IsPreset: music.IsPreset,
			}, nil
		}
		// No music configured yet — return an empty configuration (200 + null),
		// not a 404, so newly created events don't produce an unexplained error.
		return nil, nil
	}
	music, err := s.musicRepo.GetByID(ctx, *event.Sections.MusicID)
	if err != nil || music == nil {
		// Linked music record is missing — treat as no music configured.
		return nil, nil
	}
	return &dto.MusicResponse{
		ID:       music.ID,
		EventID:  music.EventID,
		Title:    music.Title,
		FileURL:  music.FileURL,
		Preset:   music.Preset,
		IsPreset: music.IsPreset,
	}, nil
}

func (s *EventService) UploadMusic(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, file io.Reader, filename, title string) (*dto.MusicResponse, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	resolver := s.resolveEntitlement(ctx, userID)
	if !resolver.MusicUpload() {
		return nil, fmt.Errorf("%w: music upload is not available on your plan", errors.ErrLimitExceeded)
	}
	if title == "" {
		title = "Uploaded music"
	}
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".mp3"
	}
	key := storage.GenerateKey("music/"+eventID.String(), ext)
	upload, err := s.storage.Upload(ctx, key, file, storage.UploadOptions{
		ContentType: "audio/mpeg",
		Extension:   ext,
	})
	if err != nil {
		return nil, fmt.Errorf("upload music: %w", err)
	}
	music := &model.Music{
		EventID: &eventID,
		Title:   title,
		FileURL: upload.URL,
	}
	if err := s.musicRepo.Create(ctx, music); err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, err
	}
	if err := s.assignMusic(ctx, eventID, music.ID); err != nil {
		return nil, err
	}
	return s.musicToResponse(music), nil
}

func (s *EventService) ListMusicPresets(ctx context.Context, userID uuid.UUID) ([]dto.MusicResponse, error) {
	resolver := s.resolveEntitlement(ctx, userID)
	if !resolver.MusicPreset() && !resolver.MusicUpload() {
		return nil, fmt.Errorf("%w: music feature is not available on your plan", errors.ErrLimitExceeded)
	}
	music, err := s.musicRepo.ListPresets(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.MusicResponse, len(music))
	for i, m := range music {
		result[i] = *s.musicToResponse(&m)
	}
	return result, nil
}

func (s *EventService) AssignPresetMusic(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, presetID uuid.UUID) (*dto.MusicResponse, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	resolver := s.resolveEntitlement(ctx, userID)
	if !resolver.MusicPreset() {
		return nil, fmt.Errorf("%w: music presets are not available on your plan", errors.ErrLimitExceeded)
	}
	music, err := s.musicRepo.GetByID(ctx, presetID)
	if err != nil || music == nil {
		return nil, errors.ErrNotFound
	}
	if !music.IsPreset {
		return nil, fmt.Errorf("%w: not a preset", errors.ErrInvalidInput)
	}
	if err := s.assignMusic(ctx, eventID, music.ID); err != nil {
		return nil, err
	}
	return s.musicToResponse(music), nil
}

// RemoveMusic detaches and deletes the currently selected music for an event,
// effectively disabling background music. Safe to call when no music is set.
func (s *EventService) RemoveMusic(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return err
	}
	sec, err := s.sectionRepo.GetByEventID(ctx, eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	if sec.MusicID == nil {
		return nil
	}
	musicID := *sec.MusicID
	if err := s.musicRepo.Delete(ctx, musicID); err != nil {
		return err
	}
	sec.MusicID = nil
	return s.sectionRepo.Update(ctx, sec)
}

func (s *EventService) assignMusic(ctx context.Context, eventID, musicID uuid.UUID) error {
	sec, err := s.sectionRepo.GetByEventID(ctx, eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			sec = &model.EventSection{EventID: eventID, MusicID: &musicID}
			return s.sectionRepo.Create(ctx, sec)
		}
		return err
	}
	sec.MusicID = &musicID
	return s.sectionRepo.Update(ctx, sec)
}

func (s *EventService) musicToResponse(music *model.Music) *dto.MusicResponse {
	return &dto.MusicResponse{
		ID:       music.ID,
		EventID:  music.EventID,
		Title:    music.Title,
		FileURL:  music.FileURL,
		Preset:   music.Preset,
		IsPreset: music.IsPreset,
	}
}

func (s *EventService) loveStoryToResponse(story *model.LoveStory) dto.LoveStoryDTO {
	return dto.LoveStoryDTO{
		ID:        story.ID,
		Title:     story.Title,
		Story:     story.Story,
		Year:      story.Year,
		Date:      story.Date,
		ImageURL:  story.ImageURL,
		SortOrder: story.SortOrder,
		CreatedAt: story.CreatedAt,
		UpdatedAt: story.UpdatedAt,
	}
}

func (s *EventService) ListLoveStories(ctx context.Context, userID, eventID uuid.UUID) ([]dto.LoveStoryDTO, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	stories, err := s.loveStoryRepo.ListByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.LoveStoryDTO, 0, len(stories))
	for _, st := range stories {
		result = append(result, s.loveStoryToResponse(&st))
	}
	return result, nil
}

func (s *EventService) CreateLoveStory(ctx context.Context, userID, eventID uuid.UUID, req dto.CreateLoveStoryRequest) (*dto.LoveStoryDTO, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	if req.SortOrder == 0 {
		count, _ := s.loveStoryRepo.CountByEvent(ctx, eventID)
		req.SortOrder = int(count) + 1
	}
	story := &model.LoveStory{
		EventID:   eventID,
		Title:     req.Title,
		Story:     req.Story,
		Year:      req.Year,
		Date:      req.Date,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
	}
	if err := s.loveStoryRepo.Create(ctx, story); err != nil {
		return nil, err
	}
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "lovestory.created",
		EntityType: "love_story",
		EntityID:   &story.ID,
	})
	resp := s.loveStoryToResponse(story)
	return &resp, nil
}

func (s *EventService) UpdateLoveStory(ctx context.Context, userID, eventID, storyID uuid.UUID, req dto.UpdateLoveStoryRequest) (*dto.LoveStoryDTO, error) {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return nil, err
	}
	story, err := s.loveStoryRepo.GetByID(ctx, storyID)
	if err != nil || story == nil || story.EventID != eventID {
		return nil, errors.ErrNotFound
	}
	if req.Title != nil {
		story.Title = *req.Title
	}
	if req.Story != nil {
		story.Story = *req.Story
	}
	if req.Year != nil {
		story.Year = *req.Year
	}
	if req.Date != nil {
		story.Date = *req.Date
	}
	if req.ImageURL != nil {
		story.ImageURL = *req.ImageURL
	}
	if req.SortOrder != nil {
		story.SortOrder = *req.SortOrder
	}
	if err := s.loveStoryRepo.Update(ctx, story); err != nil {
		return nil, err
	}
	resp := s.loveStoryToResponse(story)
	return &resp, nil
}

func (s *EventService) DeleteLoveStory(ctx context.Context, userID, eventID, storyID uuid.UUID) error {
	if _, err := s.ensureOwner(ctx, userID, eventID); err != nil {
		return err
	}
	story, err := s.loveStoryRepo.GetByID(ctx, storyID)
	if err != nil || story == nil || story.EventID != eventID {
		return errors.ErrNotFound
	}
	if err := s.loveStoryRepo.Delete(ctx, storyID); err != nil {
		return err
	}
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "lovestory.deleted",
		EntityType: "love_story",
		EntityID:   &storyID,
	})
	return nil
}

func encodeJSONSlice(data []map[string]interface{}) datatypes.JSON {
	if data == nil {
		return datatypes.JSON(nil)
	}
	bytes, _ := json.Marshal(data)
	return bytes
}

func encodeJSONMap(data map[string]interface{}) datatypes.JSON {
	if data == nil {
		return datatypes.JSON(nil)
	}
	bytes, _ := json.Marshal(data)
	return bytes
}

func (s *EventService) toResponse(ctx context.Context, event *model.Event) *dto.EventResponse {
	count, _ := s.guestRepo.CountByEvent(ctx, event.ID)
	responded, _ := s.rsvpRepo.CountRespondedByEvent(ctx, event.ID)

	return &dto.EventResponse{
		ID:               event.ID,
		UserID:           event.UserID,
		TemplateID:       event.TemplateID,
		Title:            event.Title,
		Slug:             event.Slug,
		CoupleName:       event.CoupleName,
		GroomName:        event.GroomName,
		BrideName:        event.BrideName,
		GroomParents:     event.GroomParents,
		BrideParents:     event.BrideParents,
		WeddingDate:      event.WeddingDate,
		WeddingTime:      event.WeddingTime,
		CeremonyVenue:    event.CeremonyVenue,
		CeremonyAddress:  event.CeremonyAddress,
		CeremonyMapURL:   event.CeremonyMapURL,
		ReceptionVenue:   event.ReceptionVenue,
		ReceptionAddress: event.ReceptionAddress,
		ReceptionMapURL:  event.ReceptionMapURL,
		MusicURL:         event.MusicURL,
		VideoURL:         event.VideoURL,
		Status:           event.Status,
		PublishedAt:      event.PublishedAt,
		ViewCount:        event.ViewCount,
		GuestCount:       count,
		RsvpCount:        responded,
		CreatedAt:        event.CreatedAt,
		UpdatedAt:        event.UpdatedAt,
	}
}

func decodeJSON(j datatypes.JSON) []map[string]interface{} {
	if len(j) == 0 {
		return nil
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(j, &result); err != nil {
		return nil
	}
	return result
}

func decodeJSONMap(j datatypes.JSON) map[string]interface{} {
	if len(j) == 0 {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(j, &result); err != nil {
		return nil
	}
	return result
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func (s *EventService) ensureSlugUnique(ctx context.Context, slugStr string) string {
	existing, _ := s.eventRepo.GetBySlug(ctx, slugStr)
	if existing == nil {
		return slugStr
	}
	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d", slugStr, i)
		existing, _ := s.eventRepo.GetBySlug(ctx, candidate)
		if existing == nil {
			return candidate
		}
	}
	return slugStr + "-" + uuid.New().String()[:8]
}

func toSectionsDTO(section *model.EventSection, music *dto.MusicDTO) *dto.EventSectionsDTO {
	if section == nil {
		return nil
	}
	return &dto.EventSectionsDTO{
		HeroEnabled:         section.HeroEnabled,
		CoupleEnabled:       section.CoupleEnabled,
		EventDetailsEnabled: section.EventDetailsEnabled,
		GalleryEnabled:      section.GalleryEnabled,
		VideoEnabled:        section.VideoEnabled,
		Music:               music,
		RSVPEnabled:         section.RSVPEnabled,
		GuestbookEnabled:    section.GuestbookEnabled,
		LoveStoryEnabled:    section.LoveStoryEnabled,
		DigitalGiftsEnabled: section.DigitalGiftsEnabled,
		DressCode:           section.DressCode,
		ClosingMessage:      section.ClosingMessage,
		OpeningMessage:      section.OpeningMessage,
		VerseEnabled:        section.VerseEnabled,
		VerseReligion:       section.VerseReligion,
		VerseText:           section.VerseText,
		VerseSource:         section.VerseSource,
	}
}
