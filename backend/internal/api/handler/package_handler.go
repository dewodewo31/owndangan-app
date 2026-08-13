package handler

import (
	"net/http"

	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service"
)

type PackageHandler struct {
	pkgSvc *service.PackageService
}

func NewPackageHandler(pkgSvc *service.PackageService) *PackageHandler {
	return &PackageHandler{pkgSvc: pkgSvc}
}

func (h *PackageHandler) ListActive(w http.ResponseWriter, r *http.Request) {
	packages, err := h.pkgSvc.ListActive(r.Context())
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, packages, r)
}

func (h *PackageHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	packages, err := h.pkgSvc.ListAllForAdmin(r.Context())
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, packages, r)
}

func (h *PackageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID", r)
		return
	}
	resp, err := h.pkgSvc.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *PackageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePackageRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.pkgSvc.Create(r.Context(), req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, resp, r)
}

func (h *PackageHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID", r)
		return
	}
	var req dto.UpdatePackageRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.pkgSvc.Update(r.Context(), id, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *PackageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID", r)
		return
	}
	err = h.pkgSvc.Deactivate(r.Context(), id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Package deactivated"}, r)
}
