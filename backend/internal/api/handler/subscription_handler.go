package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/service"
)

type SubscriptionHandler struct {
	subSvc *service.SubscriptionService
}

func NewSubscriptionHandler(subSvc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subSvc: subSvc}
}

func (h *SubscriptionHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	resp, err := h.subSvc.GetCurrentUserSubscription(r.Context(), userID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *SubscriptionHandler) GetUserOrDefault(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	resp, err := h.subSvc.GetUserSubscriptionOrDefault(r.Context(), userID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *SubscriptionHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Get("/current", h.GetCurrentUser)
		r.Get("/default", h.GetUserOrDefault)
	})
}
