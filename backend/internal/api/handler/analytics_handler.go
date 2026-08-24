package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/middleware"
	apperrors "github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service"
)

type AnalyticsHandler struct {
	analyticsSvc *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsSvc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsSvc: analyticsSvc}
}

type trackEventRequest struct {
	EventID uuid.UUID `json:"event_id" validate:"required,uuid"`
	Type    string    `json:"type" validate:"required"`
}

func (h *AnalyticsHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	var req trackEventRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	err := h.analyticsSvc.TrackEvent(r.Context(), req.EventID, req.Type, r.RemoteAddr, r.UserAgent())
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidInput) {
			response.Error(w, http.StatusBadRequest, "INVALID_EVENT_TYPE", "Invalid event type", r)
			return
		}
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]string{"message": "Event tracked"}, r)
}

func (h *AnalyticsHandler) GetEventAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	resp, err := h.analyticsSvc.GetEventAnalytics(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}
