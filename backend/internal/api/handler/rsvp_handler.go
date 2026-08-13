package handler

import (
	"net/http"

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

func (h *RSVPHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler, publicHandler func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(publicHandler)
		r.Post("/{eventID}/submit", h.Submit)
	})
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Get("/{eventID}", h.List)
		r.Get("/{eventID}/recap", h.Recap)
	})
}
