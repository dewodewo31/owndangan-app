package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service/rsvp"
)

type RSVPHandler struct {
	rsvpSvc *rsvp.Service
}

func NewRSVPHandler(rsvpSvc *rsvp.Service) *RSVPHandler {
	return &RSVPHandler{rsvpSvc: rsvpSvc}
}

func (h *RSVPHandler) Submit(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	var req rsvp.SubmitRSVPRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	result, err := h.rsvpSvc.Submit(r.Context(), eventID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, rsvp.ToRSVPResponse(result), r)
}

func (h *RSVPHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	rsvps, err := h.rsvpSvc.GetByEvent(r.Context(), userID, eventID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	result := make([]rsvp.RSVPResponse, len(rsvps))
	for i, r := range rsvps {
		result[i] = rsvp.ToRSVPResponse(&r)
	}
	response.JSON(w, http.StatusOK, result, r)
}

func (h *RSVPHandler) Recap(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	recap, err := h.rsvpSvc.GetRecap(r.Context(), userID, eventID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, recap, r)
}

func (h *RSVPHandler) Export(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	rows, err := h.rsvpSvc.ListForExport(r.Context(), userID, eventID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=rsvp-%s.csv", eventID.String()[:8]))
	w.WriteHeader(http.StatusOK)

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"nama", "telepon", "kehadiran", "jumlah_tamu", "pesan", "waktu_kirim"})
	for _, row := range rows {
		attendance := map[string]string{
			"attending":     "hadir",
			"not_attending": "tidak_hadir",
			"maybe":         "ragu",
		}[row.Attendance]
		if attendance == "" {
			attendance = row.Attendance
		}
		_ = cw.Write([]string{
			row.GuestName,
			row.GuestPhone,
			attendance,
			strconv.Itoa(row.GuestCount),
			row.Message,
			row.SubmittedAt.Format(time.RFC3339),
		})
	}
	cw.Flush()
}

func (h *RSVPHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler, publicHandler func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(publicHandler)
		r.Post("/{eventID}/submit", h.Submit)
	})
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Get("/{eventID}", h.List)
		r.Get("/{eventID}/recap", h.Recap)
		r.Get("/{eventID}/export", h.Export)
	})
}

// ExportByEvent streams the RSVP export for an event, supporting ?format=csv
// (all tiers) and ?format=xlsx (Pro tier only).
func (h *RSVPHandler) ExportByEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	result, err := h.rsvpSvc.Export(r.Context(), userID, eventID, r.URL.Query().Get("format"))
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", result.Filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result.Data)
}
