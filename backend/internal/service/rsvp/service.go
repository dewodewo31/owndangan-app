package rsvp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
)

type Service struct {
	rsvpRepo  repository.RSVPRepository
	guestRepo repository.GuestRepository
	eventRepo repository.EventRepository
	auditRepo repository.AuditLogRepository
}

func NewService(rsvpRepo repository.RSVPRepository, guestRepo repository.GuestRepository,
	eventRepo repository.EventRepository, auditRepo repository.AuditLogRepository) *Service {
	return &Service{
		rsvpRepo:  rsvpRepo,
		guestRepo: guestRepo,
		eventRepo: eventRepo,
		auditRepo: auditRepo,
	}
}

func (s *Service) Submit(ctx context.Context, eventID uuid.UUID, req SubmitRSVPRequest) (*model.RSVP, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.Status != "published" {
		return nil, fmt.Errorf("%w: event is not published", errors.ErrConflict)
	}

	guest, err := s.guestRepo.GetByToken(ctx, req.Token)
	if err != nil || guest == nil {
		return nil, fmt.Errorf("%w: invalid invitation token", errors.ErrNotFound)
	}
	if guest.EventID != eventID {
		return nil, fmt.Errorf("%w: token does not match event", errors.ErrForbidden)
	}

	existingRSVP, _ := s.rsvpRepo.GetByGuestID(ctx, guest.ID)
	if existingRSVP != nil {
		existingRSVP.Attendance = req.Attendance
		existingRSVP.GuestCount = req.GuestCount
		existingRSVP.Message = strings.TrimSpace(req.Message)
		if err := s.rsvpRepo.Update(ctx, existingRSVP); err != nil {
			return nil, fmt.Errorf("update rsvp: %w", err)
		}
		return existingRSVP, nil
	}

	if req.Attendance != "attending" && req.Attendance != "not_attending" && req.Attendance != "maybe" {
		return nil, fmt.Errorf("%w: attendance must be 'attending', 'not_attending', or 'maybe'", errors.ErrInvalidInput)
	}
	if req.GuestCount < 1 {
		req.GuestCount = 1
	}

	rsvp := &model.RSVP{
		GuestID:     guest.ID,
		EventID:     eventID,
		Attendance:  req.Attendance,
		GuestCount:  req.GuestCount,
		Message:     strings.TrimSpace(req.Message),
		SubmittedAt: time.Now(),
	}

	if err := s.rsvpRepo.Create(ctx, rsvp); err != nil {
		return nil, fmt.Errorf("create rsvp: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		Action:     "rsvp.submitted",
		EntityType: "rsvp",
		EntityID:   &rsvp.ID,
	})

	return rsvp, nil
}

func (s *Service) GetByEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) ([]model.RSVP, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return s.rsvpRepo.GetByEventID(ctx, eventID)
}

func (s *Service) GetRecap(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*RSVPRecap, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	attending, _ := s.rsvpRepo.CountByAttendance(ctx, eventID, "attending")
	notAttending, _ := s.rsvpRepo.CountByAttendance(ctx, eventID, "not_attending")
	maybe, _ := s.rsvpRepo.CountByAttendance(ctx, eventID, "maybe")
	attendingGuests, _ := s.rsvpRepo.SumGuestCountByAttendance(ctx, eventID, "attending")

	return &RSVPRecap{
		TotalResponded:   int(attending + notAttending + maybe),
		Attending:        int(attending),
		NotAttending:     int(notAttending),
		Maybe:            int(maybe),
		TotalGuestCount:  int(attendingGuests),
	}, nil
}

func (s *Service) ListForExport(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) ([]repository.RSVPExportRow, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return s.rsvpRepo.ListExportRows(ctx, eventID)
}
