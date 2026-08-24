package guest

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/entitlement"
)

type Service struct {
	guestRepo repository.GuestRepository
	eventRepo repository.EventRepository
	subRepo   repository.SubscriptionRepository
	pkgRepo   repository.PackageRepository
	auditRepo repository.AuditLogRepository
}

func NewService(guestRepo repository.GuestRepository, eventRepo repository.EventRepository,
	subRepo repository.SubscriptionRepository, pkgRepo repository.PackageRepository,
	auditRepo repository.AuditLogRepository) *Service {
	return &Service{
		guestRepo: guestRepo,
		eventRepo: eventRepo,
		subRepo:   subRepo,
		pkgRepo:   pkgRepo,
		auditRepo: auditRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req CreateGuestRequest) (*model.Guest, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if err := s.checkGuestLimit(ctx, eventID, event.UserID, 1); err != nil {
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("%w: guest name is required", errors.ErrInvalidInput)
	}

	token := generateToken()
	for {
		taken, _ := s.guestRepo.IsTokenTaken(ctx, token)
		if !taken {
			break
		}
		token = generateToken()
	}

	guest := &model.Guest{
		EventID:  eventID,
		Name:     req.Name,
		Phone:    strings.TrimSpace(req.Phone),
		Category: req.Category,
		Note:     strings.TrimSpace(req.Note),
		Token:    token,
	}

	if guest.Category == "" {
		guest.Category = "family"
	}

	if err := s.guestRepo.Create(ctx, guest); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, fmt.Errorf("%w: guest already exists", errors.ErrConflict)
		}
		return nil, fmt.Errorf("create guest: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "guest.created",
		EntityType: "guest",
		EntityID:   &guest.ID,
	})

	return guest, nil
}

func (s *Service) GetByID(ctx context.Context, userID uuid.UUID, guestID uuid.UUID) (*model.Guest, error) {
	guest, err := s.guestRepo.GetByID(ctx, guestID)
	if err != nil || guest == nil {
		return nil, errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, guest.EventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return guest, nil
}

func (s *Service) ListByEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, page, perPage int) ([]model.Guest, int64, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, 0, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, 0, errors.ErrForbidden
	}

	return s.guestRepo.ListByEvent(ctx, eventID, page, perPage)
}

func (s *Service) ListDeleted(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) ([]model.Guest, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return s.guestRepo.ListDeleted(ctx, eventID)
}

func (s *Service) Update(ctx context.Context, userID uuid.UUID, guestID uuid.UUID, req UpdateGuestRequest) (*model.Guest, error) {
	guest, err := s.guestRepo.GetByID(ctx, guestID)
	if err != nil || guest == nil {
		return nil, errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, guest.EventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
		if *req.Name == "" {
			return nil, fmt.Errorf("%w: guest name cannot be empty", errors.ErrInvalidInput)
		}
		guest.Name = *req.Name
	}
	if req.Phone != nil {
		guest.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Category != nil {
		guest.Category = *req.Category
	}
	if req.Note != nil {
		guest.Note = strings.TrimSpace(*req.Note)
	}

	if err := s.guestRepo.Update(ctx, guest); err != nil {
		return nil, fmt.Errorf("update guest: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "guest.updated",
		EntityType: "guest",
		EntityID:   &guest.ID,
	})

	return guest, nil
}

func (s *Service) Restore(ctx context.Context, eventID uuid.UUID, guestID uuid.UUID, userID uuid.UUID) error {
	guest, err := s.guestRepo.GetByIdUnscoped(ctx, guestID)
	if err != nil || guest == nil {
		return errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, guest.EventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.guestRepo.Restore(ctx, guestID, eventID); err != nil {
		return fmt.Errorf("restore guest: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "guest.restored",
		EntityType: "guest",
		EntityID:   &guest.ID,
	})

	return nil
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID, guestID uuid.UUID) error {
	guest, err := s.guestRepo.GetByID(ctx, guestID)
	if err != nil || guest == nil {
		return errors.ErrNotFound
	}

	event, err := s.eventRepo.GetByID(ctx, guest.EventID)
	if err != nil || event == nil {
		return errors.ErrNotFound
	}
	if event.UserID != userID {
		return errors.ErrForbidden
	}

	if err := s.guestRepo.SoftDelete(ctx, guestID); err != nil {
		return fmt.Errorf("delete guest: %w", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "guest.deleted",
		EntityType: "guest",
		EntityID:   &guest.ID,
	})

	return nil
}

func (s *Service) ImportCSV(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, reader io.Reader) (*ImportResult, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	currentCount, _ := s.guestRepo.CountByEvent(ctx, eventID)
	if err := s.checkGuestLimit(ctx, eventID, userID, int(currentCount)); err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: invalid CSV format", errors.ErrInvalidInput)
	}

	nameIdx := -1
	phoneIdx := -1
	categoryIdx := -1
	noteIdx := -1

	for i, col := range header {
		switch strings.ToLower(strings.TrimSpace(col)) {
		case "name", "nama":
			nameIdx = i
		case "phone", "telepon", "hp":
			phoneIdx = i
		case "category", "kategori", "group":
			categoryIdx = i
		case "note", "catatan":
			noteIdx = i
		}
	}

	if nameIdx == -1 {
		return nil, fmt.Errorf("%w: CSV must have a 'name' column", errors.ErrInvalidInput)
	}

	result := &ImportResult{}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", result.Total+1, err))
			continue
		}

		result.Total++
		if nameIdx >= len(record) || strings.TrimSpace(record[nameIdx]) == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: name is required", result.Total))
			continue
		}

		guestName := strings.TrimSpace(record[nameIdx])

		var phone, category, note string
		if phoneIdx >= 0 && phoneIdx < len(record) {
			phone = strings.TrimSpace(record[phoneIdx])
		}
		if categoryIdx >= 0 && categoryIdx < len(record) {
			category = strings.TrimSpace(record[categoryIdx])
		}
		if noteIdx >= 0 && noteIdx < len(record) {
			note = strings.TrimSpace(record[noteIdx])
		}

		if s.isDuplicate(ctx, eventID, guestName, phone) {
			result.Skipped++
			continue
		}

		token := generateToken()
		for {
			taken, _ := s.guestRepo.IsTokenTaken(ctx, token)
			if !taken {
				break
			}
			token = generateToken()
		}

		guest := &model.Guest{
			EventID:  eventID,
			Name:     guestName,
			Phone:    phone,
			Category: category,
			Note:     note,
			Token:    token,
		}
		if guest.Category == "" {
			guest.Category = "family"
		}

		if err := s.guestRepo.Create(ctx, guest); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", result.Total, err))
			continue
		}

		result.Imported++
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &userID,
		Action:     "guest.imported",
		EntityType: "guest",
		EntityID:   &eventID,
	})

	return result, nil
}

func (s *Service) checkGuestLimit(ctx context.Context, eventID uuid.UUID, userID uuid.UUID, newGuests int) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil
	}

	sub, err := s.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil || sub == nil {
		return fmt.Errorf("%w: active subscription required to add guests", errors.ErrPaymentRequired)
	}

	pkg, err := s.pkgRepo.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return errors.ErrNotFound
	}

	resolver := entitlement.NewResolver(pkg)
	currentCount, _ := s.guestRepo.CountByEvent(ctx, eventID)
	if !resolver.CanCreateGuest(int(currentCount) + newGuests) {
		return fmt.Errorf("%w: guest limit reached for your plan", errors.ErrLimitExceeded)
	}

	return nil
}

func (s *Service) isDuplicate(ctx context.Context, eventID uuid.UUID, name string, phone string) bool {
	guests, _, _ := s.guestRepo.ListByEvent(ctx, eventID, 1, 10000)
	nameLower := strings.ToLower(strings.TrimSpace(name))
	for _, g := range guests {
		if strings.ToLower(strings.TrimSpace(g.Name)) == nameLower {
			if phone == "" || strings.TrimSpace(g.Phone) == "" {
				return true
			}
			if strings.TrimSpace(g.Phone) == strings.TrimSpace(phone) {
				return true
			}
		}
	}
	return false
}

func generateToken() string {
	return uuid.New().String()[:8]
}
