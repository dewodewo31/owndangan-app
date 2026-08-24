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
		Category: model.NormalizeGuestCategory(req.Category),
		Note:     strings.TrimSpace(req.Note),
		Token:    token,
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
		guest.Category = model.NormalizeGuestCategory(*req.Category)
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

// PreviewImport parses a CSV guest file and returns a preview with per-row
// validation and duplicate detection (within file and against existing guests).
// It does NOT persist anything.
func (s *Service) PreviewImport(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, reader io.Reader, mapping *ImportMapping) (*ImportPreview, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: invalid CSV format", errors.ErrInvalidInput)
	}

	nameIdx := findColumn(header, candidates([]string{"name", "nama"}, mappingVal(mapping, "name")))
	emailIdx := findColumn(header, candidates([]string{"email", "e-mail", "alamat email"}, mappingVal(mapping, "email")))
	phoneIdx := findColumn(header, candidates([]string{"phone", "telepon", "hp", "no hp"}, mappingVal(mapping, "phone")))
	categoryIdx := findColumn(header, candidates([]string{"category", "kategori", "group"}, mappingVal(mapping, "category")))

	if nameIdx == -1 {
		return nil, fmt.Errorf("%w: CSV must have a 'name' column", errors.ErrInvalidInput)
	}

	existing, err := s.guestRepo.FindExistingForEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("load existing guests: %w", err)
	}
	existingPhones := make(map[string]bool, len(existing))
	existingNames := make(map[string]bool, len(existing))
	for _, g := range existing {
		if p := strings.TrimSpace(strings.ToLower(g.Phone)); p != "" {
			existingPhones[p] = true
		}
		if n := strings.TrimSpace(strings.ToLower(g.Name)); n != "" {
			existingNames[n] = true
		}
	}

	preview := &ImportPreview{Columns: header}
	seen := make(map[string]bool)
	index := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		index++
		if err != nil {
			preview.Rows = append(preview.Rows, ImportPreviewRow{
				Index:  index,
				Status: "invalid",
				Errors: []string{fmt.Sprintf("could not parse row: %v", err)},
			})
			preview.Summary.Total++
			preview.Summary.Invalid++
			continue
		}

		row := buildPreviewRow(record, index, nameIdx, emailIdx, phoneIdx, categoryIdx)
		key := contactKey(row.Email, row.Phone, row.Name)

		switch {
		case row.Status == "invalid":
			preview.Summary.Invalid++
		case seen[key]:
			row.Status = "duplicate"
			row.Errors = append(row.Errors, "duplicate within file")
			preview.Summary.Duplicate++
		case (row.Phone != "" && existingPhones[strings.ToLower(row.Phone)]) ||
			(row.Name != "" && existingNames[strings.ToLower(row.Name)]):
			row.Status = "duplicate"
			row.Errors = append(row.Errors, "already exists for this event")
			preview.Summary.Duplicate++
		default:
			seen[key] = true
			preview.Summary.Valid++
		}
		preview.Rows = append(preview.Rows, row)
		preview.Summary.Total++
	}

	return preview, nil
}

// ConfirmImport re-validates the selected rows and inserts them. It de-duplicates
// again against the database and within the request batch.
func (s *Service) ConfirmImport(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req ImportConfirmRequest) (*ImportConfirmResult, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if err := s.checkGuestLimit(ctx, eventID, userID, len(req.Rows)); err != nil {
		return nil, err
	}

	existing, err := s.guestRepo.FindExistingForEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("load existing guests: %w", err)
	}
	existingPhones := make(map[string]bool, len(existing))
	existingNames := make(map[string]bool, len(existing))
	for _, g := range existing {
		if p := strings.TrimSpace(strings.ToLower(g.Phone)); p != "" {
			existingPhones[p] = true
		}
		if n := strings.TrimSpace(strings.ToLower(g.Name)); n != "" {
			existingNames[n] = true
		}
	}

	result := &ImportConfirmResult{Total: len(req.Rows)}
	seen := make(map[string]bool)
	for i, r := range req.Rows {
		idx := i + 1
		name := strings.TrimSpace(r.Name)
		phone := strings.TrimSpace(r.Phone)
		email := strings.TrimSpace(r.Email)
		category := model.NormalizeGuestCategory(r.Category)

		if name == "" {
			result.Errors = append(result.Errors, ImportConfirmError{Index: idx, Errors: []string{"name is required"}})
			continue
		}

		key := contactKey(email, phone, name)
		if seen[key] ||
			(phone != "" && existingPhones[strings.ToLower(phone)]) ||
			(name != "" && existingNames[strings.ToLower(name)]) {
			result.Duplicates++
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
			Name:     name,
			Phone:    phone,
			Category: category,
			Token:    token,
		}
		if err := s.guestRepo.Create(ctx, guest); err != nil {
			result.Errors = append(result.Errors, ImportConfirmError{Index: idx, Errors: []string{err.Error()}})
			continue
		}
		seen[key] = true
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

func candidates(def []string, override string) []string {
	if override != "" {
		return []string{override}
	}
	return def
}

func mappingVal(m *ImportMapping, field string) string {
	if m == nil {
		return ""
	}
	switch field {
	case "name":
		return m.Name
	case "email":
		return m.Email
	case "phone":
		return m.Phone
	case "category":
		return m.Category
	}
	return ""
}

func findColumn(header []string, names []string) int {
	for _, n := range names {
		for i, h := range header {
			if strings.EqualFold(strings.TrimSpace(h), n) {
				return i
			}
		}
	}
	return -1
}

func contactKey(email, phone, name string) string {
	e := strings.TrimSpace(strings.ToLower(email))
	p := strings.TrimSpace(strings.ToLower(phone))
	n := strings.TrimSpace(strings.ToLower(name))
	if e != "" {
		return "e:" + e
	}
	if p != "" {
		return "p:" + p
	}
	return "n:" + n
}

func buildPreviewRow(record []string, index, nameIdx, emailIdx, phoneIdx, categoryIdx int) ImportPreviewRow {
	get := func(i int) string {
		if i >= 0 && i < len(record) {
			return strings.TrimSpace(record[i])
		}
		return ""
	}
	row := ImportPreviewRow{
		Index:    index,
		Name:     get(nameIdx),
		Email:    get(emailIdx),
		Phone:    get(phoneIdx),
		Category: model.NormalizeGuestCategory(get(categoryIdx)),
	}
	if row.Name == "" {
		row.Status = "invalid"
		row.Errors = append(row.Errors, "name is required")
	} else {
		row.Status = "valid"
	}
	return row
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

func generateToken() string {
	return uuid.New().String()[:8]
}
