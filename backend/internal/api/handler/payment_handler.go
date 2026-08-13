package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service"
)

type PaymentHandler struct {
	paySvc *service.PaymentService
}

func NewPaymentHandler(paySvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paySvc: paySvc}
}

func (h *PaymentHandler) CreateSnap(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req dto.CreateSnapRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.paySvc.CreateSnapTransaction(r.Context(), userID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, resp, r)
}

func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload dto.MidtransWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload", r)
		return
	}
	if err := h.paySvc.HandleWebhook(r.Context(), payload); err != nil {
		response.FromError(w, err, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PaymentHandler) ListUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	_, page, perPage := parsePagination(r)
	resp, total, err := h.paySvc.ListUserTransactions(r.Context(), userID, page, perPage)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSONPaginated(w, http.StatusOK, resp, page, perPage, int(total), r)
}

func (h *PaymentHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Post("/snap", h.CreateSnap)
		r.Get("/transactions", h.ListUserTransactions)
	})
	r.Post("/webhook", h.HandleWebhook)
}
