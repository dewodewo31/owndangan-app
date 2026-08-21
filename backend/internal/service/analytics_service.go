package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/entitlement"
)

var validClickTypes = map[string]bool{
	"whatsapp_click": true,
	"map_click":      true,
	"phone_click":    true,
}

type EventAnalytics struct {
	Views          int64 `json:"views"`
	UniqueViews    int64 `json:"unique_views"`
	WhatsappClicks int64 `json:"whatsapp_clicks"`
	MapClicks      int64 `json:"map_clicks"`
	PhoneClicks    int64 `json:"phone_clicks"`
	RSVPCount      int64 `json:"rsvp_count"`
}

type AnalyticsService struct {
	eventRepo    repository.EventRepository
	analyticsRepo repository.AnalyticsEventRepository
	rsvpRepo     repository.RSVPRepository
	subRepo      repository.SubscriptionRepository
	pkgRepo      repository.PackageRepository
}

func NewAnalyticsService(eventRepo repository.EventRepository, analyticsRepo repository.AnalyticsEventRepository,
	rsvpRepo repository.RSVPRepository, subRepo repository.SubscriptionRepository,
	pkgRepo repository.PackageRepository) *AnalyticsService {
	return &AnalyticsService{
		eventRepo:     eventRepo,
		analyticsRepo: analyticsRepo,
		rsvpRepo:      rsvpRepo,
		subRepo:       subRepo,
		pkgRepo:       pkgRepo,
	}
}

func (s *AnalyticsService) TrackEvent(ctx context.Context, eventID uuid.UUID, eventType, ip, userAgent string) error {
	if !validClickTypes[eventType] {
		return fmt.Errorf("%w: invalid event type", errors.ErrInvalidInput)
	}

	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.Status != "published" {
		return errors.ErrNotFound
	}

	return s.analyticsRepo.Create(ctx, &model.AnalyticsEvent{
		EventID:   &eventID,
		EventType: eventType,
		IPAddress: ip,
		UserAgent: userAgent,
	})
}

func (s *AnalyticsService) GetEventAnalytics(ctx context.Context, userID, eventID uuid.UUID) (*EventAnalytics, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if !s.resolveEntitlement(ctx, userID).AnalyticsEnabled() {
		return nil, errors.ErrForbidden
	}

	views, err := s.analyticsRepo.CountByTypeForEvent(ctx, eventID, "page_view")
	if err != nil {
		return nil, err
	}
	uniqueViews, err := s.analyticsRepo.CountUniqueByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	whatsappClicks, err := s.analyticsRepo.CountByTypeForEvent(ctx, eventID, "whatsapp_click")
	if err != nil {
		return nil, err
	}
	mapClicks, err := s.analyticsRepo.CountByTypeForEvent(ctx, eventID, "map_click")
	if err != nil {
		return nil, err
	}
	phoneClicks, err := s.analyticsRepo.CountByTypeForEvent(ctx, eventID, "phone_click")
	if err != nil {
		return nil, err
	}
	rsvpCount, err := s.rsvpRepo.CountRespondedByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return &EventAnalytics{
		Views:          views,
		UniqueViews:    uniqueViews,
		WhatsappClicks: whatsappClicks,
		MapClicks:      mapClicks,
		PhoneClicks:    phoneClicks,
		RSVPCount:      rsvpCount,
	}, nil
}

func (s *AnalyticsService) resolveEntitlement(ctx context.Context, userID uuid.UUID) *entitlement.Resolver {
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