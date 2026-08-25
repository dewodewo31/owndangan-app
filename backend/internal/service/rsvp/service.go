package rsvp

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/email"
	"github.com/owndangan/backend/internal/service/entitlement"
	"github.com/xuri/excelize/v2"
)

type EmailSender interface {
	SendAsync(to, subject, htmlBody string)
}

type Service struct {
	rsvpRepo  repository.RSVPRepository
	guestRepo repository.GuestRepository
	eventRepo repository.EventRepository
	userRepo  repository.UserRepository
	subRepo   repository.SubscriptionRepository
	pkgRepo   repository.PackageRepository
	auditRepo repository.AuditLogRepository
	emailSvc  EmailSender
}

func NewService(rsvpRepo repository.RSVPRepository, guestRepo repository.GuestRepository,
	eventRepo repository.EventRepository, userRepo repository.UserRepository,
	subRepo repository.SubscriptionRepository, pkgRepo repository.PackageRepository,
	auditRepo repository.AuditLogRepository, emailSvc EmailSender) *Service {
	return &Service{
		rsvpRepo:  rsvpRepo,
		guestRepo: guestRepo,
		eventRepo: eventRepo,
		userRepo:  userRepo,
		subRepo:   subRepo,
		pkgRepo:   pkgRepo,
		auditRepo: auditRepo,
		emailSvc:  emailSvc,
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
		s.notifyOwnerRSVP(ctx, event, guest, existingRSVP)
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

	s.notifyOwnerRSVP(ctx, event, guest, rsvp)
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
		TotalResponded:  int(attending + notAttending + maybe),
		Attending:       int(attending),
		NotAttending:    int(notAttending),
		Maybe:           int(maybe),
		TotalGuestCount: int(attendingGuests),
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

// ExportResult is the rendered export payload ready to stream to the client.
type ExportResult struct {
	ContentType string
	Filename    string
	Data        []byte
}

// Export renders the RSVP export for an event. CSV is available on all tiers.
// XLSX requires a Pro-tier subscription and returns ErrForbidden otherwise.
func (s *Service) Export(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, format string) (*ExportResult, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	if format == "" {
		format = "csv"
	}
	format = strings.ToLower(format)

	rows, err := s.rsvpRepo.ListExportRows(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if format == "xlsx" {
		resolver, err := s.resolveEventTier(ctx, eventID)
		if err != nil {
			return nil, err
		}
		if !resolver.IsProTier() {
			return nil, errors.ErrForbidden
		}
		data, err := buildXLSX(rows)
		if err != nil {
			return nil, fmt.Errorf("build xlsx: %w", err)
		}
		return &ExportResult{
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Filename:    fmt.Sprintf("rsvp-%s.xlsx", eventID.String()[:8]),
			Data:        data,
		}, nil
	}

	if format != "csv" {
		return nil, fmt.Errorf("%w: unsupported format %q", errors.ErrInvalidInput, format)
	}

	var buf bytes.Buffer
	cw := csv.NewWriter(&buf)
	_ = cw.Write([]string{"nama", "telepon", "status_rsvp", "kehadiran", "jumlah_tamu", "waktu_kirim"})
	for _, row := range rows {
		_ = cw.Write([]string{
			row.GuestName,
			row.GuestPhone,
			"responded",
			attendanceLabel(row.Attendance),
			strconv.Itoa(row.GuestCount),
			row.SubmittedAt.Format(time.RFC3339),
		})
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return nil, fmt.Errorf("write csv: %w", err)
	}
	return &ExportResult{
		ContentType: "text/csv; charset=utf-8",
		Filename:    fmt.Sprintf("rsvp-%s.csv", eventID.String()[:8]),
		Data:        buf.Bytes(),
	}, nil
}

func buildXLSX(rows []repository.RSVPExportRow) ([]byte, error) {
	f := excelize.NewFile()
	const sheet = "RSVP"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"Nama", "Telepon", "Status RSVP", "Kehadiran", "Jumlah Tamu", "Waktu Kirim"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for i, row := range rows {
		r := i + 2
		values := []interface{}{
			row.GuestName,
			row.GuestPhone,
			"responded",
			attendanceLabel(row.Attendance),
			row.GuestCount,
			row.SubmittedAt.Format(time.RFC3339),
		}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, r)
			f.SetCellValue(sheet, cell, v)
		}
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Service) resolveEventTier(ctx context.Context, eventID uuid.UUID) (*entitlement.Resolver, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, errors.ErrNotFound
	}
	return entitlement.ResolveForUser(ctx, s.subRepo, s.pkgRepo, event.UserID), nil
}

func attendanceLabel(attendance string) string {
	switch attendance {
	case "attending":
		return "Hadir"
	case "not_attending":
		return "Tidak Hadir"
	case "maybe":
		return "Ragu-ragu"
	default:
		return attendance
	}
}

func (s *Service) notifyOwnerRSVP(ctx context.Context, event *model.Event, guest *model.Guest, rsvp *model.RSVP) {
	if s.emailSvc == nil {
		return
	}
	owner, err := s.userRepo.GetByID(ctx, event.UserID)
	if err != nil || owner == nil {
		return
	}
	html, err := email.RenderRSVP(email.RSVPData{
		OwnerName:   owner.Name,
		GuestName:   guest.Name,
		Invitation:  event.Title,
		Attendance:  attendanceLabel(rsvp.Attendance),
		GuestCount:  rsvp.GuestCount,
		SubmittedAt: rsvp.SubmittedAt.Format("2 January 2006 15:04"),
	})
	if err != nil {
		return
	}
	s.emailSvc.SendAsync(owner.Email, "Ada konfirmasi kehadiran baru", html)
}
