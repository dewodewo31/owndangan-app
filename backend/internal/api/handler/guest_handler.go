package handler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service/guest"
)

type GuestHandler struct {
	guestSvc *guest.Service
}

func NewGuestHandler(guestSvc *guest.Service) *GuestHandler {
	return &GuestHandler{guestSvc: guestSvc}
}

func (h *GuestHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	var req guest.CreateGuestRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	g, err := h.guestSvc.Create(r.Context(), userID, eventID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusCreated, guest.ToGuestResponse(g), r)
}

func (h *GuestHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	_, page, perPage := parsePagination(r)
	guests, total, err := h.guestSvc.ListByEvent(r.Context(), userID, eventID, page, perPage)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	result := make([]guest.GuestResponse, len(guests))
	for i, g := range guests {
		result[i] = guest.ToGuestResponse(&g)
	}
	response.JSONPaginated(w, http.StatusOK, result, page, perPage, int(total), r)
}

func (h *GuestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	guestID, err := parseUUID(r, "guestID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid guest ID", r)
		return
	}

	g, err := h.guestSvc.GetByID(r.Context(), userID, guestID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, guest.ToGuestResponse(g), r)
}

func (h *GuestHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	guestID, err := parseUUID(r, "guestID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid guest ID", r)
		return
	}

	var req guest.UpdateGuestRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	g, err := h.guestSvc.Update(r.Context(), userID, guestID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, guest.ToGuestResponse(g), r)
}

func (h *GuestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	guestID, err := parseUUID(r, "guestID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid guest ID", r)
		return
	}

	if err := h.guestSvc.Delete(r.Context(), userID, guestID); err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Guest deleted"}, r)
}

func (h *GuestHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FILE", "CSV file is required", r)
		return
	}
	defer file.Close()

	result, err := h.guestSvc.ImportCSV(r.Context(), userID, eventID, io.Reader(file))
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, result, r)
}

func (h *GuestHandler) Restore(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	guestID, err := parseUUID(r, "guestID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid guest ID", r)
		return
	}

	if err := h.guestSvc.Restore(r.Context(), eventID, guestID, userID); err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Guest restored"}, r)
}

func (h *GuestHandler) ListDeleted(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID, err := parseUUID(r, "eventID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}

	guests, err := h.guestSvc.ListDeleted(r.Context(), userID, eventID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	result := make([]guest.GuestResponse, len(guests))
	for i, g := range guests {
		result[i] = guest.ToGuestResponse(&g)
	}
	response.JSON(w, http.StatusOK, result, r)
}

func (h *GuestHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Post("/import", h.ImportCSV)
		r.Get("/deleted", h.ListDeleted)
		r.Get("/{guestID}", h.GetByID)
		r.Put("/{guestID}", h.Update)
		r.Delete("/{guestID}", h.Delete)
		r.Post("/{guestID}/restore", h.Restore)
	})
}
