package guestbook

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/email"
)

type EmailSender interface {
	SendAsync(to, subject, htmlBody string)
}

type Service struct {
	guestbookRepo repository.GuestbookRepository
	eventRepo     repository.EventRepository
	userRepo      repository.UserRepository
	auditRepo     repository.AuditLogRepository
	emailSvc      EmailSender
}

func NewService(guestbookRepo repository.GuestbookRepository, eventRepo repository.EventRepository,
	userRepo repository.UserRepository, auditRepo repository.AuditLogRepository, emailSvc EmailSender) *Service {
	return &Service{
		guestbookRepo: guestbookRepo,
		eventRepo:     eventRepo,
		userRepo:      userRepo,
		auditRepo:     auditRepo,
		emailSvc:      emailSvc,
	}
}

func (s *Service) Submit(ctx context.Context, eventID uuid.UUID, req SubmitMessageRequest) (*model.GuestbookMessage, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.Status != "published" {
		return nil, fmt.Errorf("%w: event is not published", errors.ErrConflict)
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" || req.Message == "" {
		return nil, fmt.Errorf("%w: name and message are required", errors.ErrInvalidInput)
	}

	msg := &model.GuestbookMessage{
		EventID:    eventID,
		Name:       req.Name,
		Message:    req.Message,
		IsApproved: false,
	}

	if err := s.guestbookRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("create guestbook message: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		Action:     "guestbook.submitted",
		EntityType: "guestbook_message",
		EntityID:   &msg.ID,
	})

	s.notifyOwnerGuestbook(ctx, event, msg)
	return msg, nil
}

func (s *Service) ListPublic(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.GuestbookMessage, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.guestbookRepo.ListByEventPaged(ctx, eventID, true, limit, offset)
}

func (s *Service) ListAll(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, limit, offset int) ([]model.GuestbookMessage, int64, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, 0, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, 0, errors.ErrForbidden
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.guestbookRepo.ListByEventPaged(ctx, eventID, false, limit, offset)
}

func (s *Service) Approve(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) error {
	msg, err := s.guestbookRepo.GetByID(ctx, messageID)
	if err != nil || msg == nil {
		return errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, msg.EventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.guestbookRepo.Approve(ctx, messageID); err != nil {
		return fmt.Errorf("approve message: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "guestbook.approved",
		EntityType: "guestbook_message",
		EntityID:   &messageID,
	})

	return nil
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) error {
	msg, err := s.guestbookRepo.GetByID(ctx, messageID)
	if err != nil || msg == nil {
		return errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, msg.EventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	return s.guestbookRepo.Delete(ctx, messageID)
}

func (s *Service) notifyOwnerGuestbook(ctx context.Context, event *model.Event, msg *model.GuestbookMessage) {
	if s.emailSvc == nil {
		return
	}
	owner, err := s.userRepo.GetByID(ctx, event.UserID)
	if err != nil || owner == nil {
		return
	}
	html, err := email.RenderGuestbook(email.GuestbookData{
		OwnerName:   owner.Name,
		GuestName:   msg.Name,
		Message:     msg.Message,
		Invitation:  event.Title,
		SubmittedAt: msg.CreatedAt.Format("2 January 2006 15:04"),
	})
	if err != nil {
		return
	}
	s.emailSvc.SendAsync(owner.Email, "Ada pesan baru di buku tamu", html)
}
