package admin

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/pkg/pagination"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service/admin"
)

type Handler struct {
	adminSvc *admin.Service
}

func NewHandler(adminSvc *admin.Service) *Handler {
	return &Handler{adminSvc: adminSvc}
}

func parseUUID(r *http.Request, key string) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, key))
}

func parsePaginationAdmin(r *http.Request) (pagination.Params, int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	return pagination.Params{Page: page, PerPage: perPage}, page, perPage
}

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.adminSvc.GetAnalytics(r.Context())
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, analytics, r)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	_, page, perPage := parsePaginationAdmin(r)
	status := r.URL.Query().Get("status")

	users, total, err := h.adminSvc.ListUsers(r.Context(), page, perPage, status)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSONPaginated(w, http.StatusOK, users, page, perPage, int(total), r)
}

func (h *Handler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(r, "userID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", r)
		return
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=active suspended"`
	}
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	if err := h.adminSvc.UpdateUserStatus(r.Context(), userID, req.Status); err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "User status updated"}, r)
}

func (h *Handler) ListPackages(w http.ResponseWriter, r *http.Request) {
	pkgs, err := h.adminSvc.GetPackages(r.Context())
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, pkgs, r)
}

func (h *Handler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	var req admin.CreatePackageRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	pkg, err := h.adminSvc.CreatePackage(r.Context(), req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusCreated, pkg, r)
}

func (h *Handler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	pkgID, err := parseUUID(r, "packageID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID", r)
		return
	}

	var req admin.UpdatePackageRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	pkg, err := h.adminSvc.UpdatePackage(r.Context(), pkgID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, pkg, r)
}

func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.adminSvc.GetTemplates(r.Context())
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, templates, r)
}

func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req admin.CreateTemplateRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	t, err := h.adminSvc.CreateTemplate(r.Context(), req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusCreated, t, r)
}

func (h *Handler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateID, err := parseUUID(r, "templateID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid template ID", r)
		return
	}

	var req admin.UpdateTemplateRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	t, err := h.adminSvc.UpdateTemplate(r.Context(), templateID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, t, r)
}

func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	_, page, perPage := parsePaginationAdmin(r)
	status := r.URL.Query().Get("status")

	txns, total, err := h.adminSvc.GetTransactions(r.Context(), page, perPage, status)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSONPaginated(w, http.StatusOK, txns, page, perPage, int(total), r)
}

func (h *Handler) GetTransactionDetail(w http.ResponseWriter, r *http.Request) {
	txnID, err := parseUUID(r, "transactionID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid transaction ID", r)
		return
	}

	txn, err := h.adminSvc.GetTransactionDetail(r.Context(), txnID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}

	response.JSON(w, http.StatusOK, txn, r)
}

func (h *Handler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler, adminRequired func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authRequired, adminRequired)

		r.Get("/analytics", h.GetAnalytics)
		r.Get("/users", h.ListUsers)
		r.Put("/users/{userID}/status", h.UpdateUserStatus)

		r.Get("/packages", h.ListPackages)
		r.Post("/packages", h.CreatePackage)
		r.Put("/packages/{packageID}", h.UpdatePackage)

		r.Get("/templates", h.ListTemplates)
		r.Post("/templates", h.CreateTemplate)
		r.Put("/templates/{templateID}", h.UpdateTemplate)

		r.Get("/transactions", h.ListTransactions)
		r.Get("/transactions/{transactionID}", h.GetTransactionDetail)
	})
}
