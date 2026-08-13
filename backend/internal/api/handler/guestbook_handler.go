package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service/guestbook"
)

type GuestbookHandler struct {
	guestbookSvc *guestbook.Service
}

func NewGuestbookHandler(guestbookSvc *guestbook.Service) *GuestbookHandler {
	return &GuestbookHandler{guestbookSvc: guestbookSvc}
}

func (h *GuestbookHandler) Submit(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	var req guestbook.SubmitMessageRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	msg, err := h.guestbookSvc.Submit(r.Context(), eventID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusCreated, guestbook.ToGuestbookResponse(msg), r)
}

func (h *GuestbookHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	msgs, err := h.guestbookSvc.ListPublic(r.Context(), eventID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	result := make([]guestbook.GuestbookResponse, len(msgs))
	for i, m := range msgs {
		result[i] = guestbook.ToGuestbookResponse(&m)
	}
	response.JSON(w, http.StatusOK, result, r)
}

func (h *GuestbookHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	msgs, err := h.guestbookSvc.ListAll(r.Context(), userID, eventID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	result := make([]guestbook.GuestbookResponse, len(msgs))
	for i, m := range msgs {
		result[i] = guestbook.ToGuestbookResponse(&m)
	}
	response.JSON(w, http.StatusOK, result, r)
}

func (h *GuestbookHandler) Approve(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	messageID, err := parseUUID(r, "messageID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid message ID", r)
		return
	}

	if err := h.guestbookSvc.Approve(r.Context(), userID, messageID); err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Message approved"}, r)
}

func (h *GuestbookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	messageID, err := parseUUID(r, "messageID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid message ID", r)
		return
	}

	if err := h.guestbookSvc.Delete(r.Context(), userID, messageID); err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Message deleted"}, r)
}

func (h *GuestbookHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler, publicHandler func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(publicHandler)
		r.Post("/{eventID}/submit", h.Submit)
		r.Get("/{eventID}", h.ListPublic)
	})
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Get("/{eventID}/all", h.ListAll)
		r.Post("/{messageID}/approve", h.Approve)
		r.Delete("/{messageID}", h.Delete)
	})
}
