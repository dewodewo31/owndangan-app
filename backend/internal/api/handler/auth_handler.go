package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	user, accessToken, refreshToken, expiresIn, err := h.authSvc.Register(r.Context(), req.Name, req.Email, req.Password, req.Phone)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, toAuthResponse(user, accessToken, refreshToken, expiresIn), r)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	user, accessToken, refreshToken, expiresIn, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, toAuthResponse(user, accessToken, refreshToken, expiresIn), r)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	accessToken, refreshToken, expiresIn, err := h.authSvc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, r)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	err := h.authSvc.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"}, r)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", r)
		return
	}
	var req dto.ChangePasswordRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}

	err := h.authSvc.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"}, r)
}

func (h *AuthHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	r.With(authRequired).Post("/change-password", h.ChangePassword)
}

func (h *AuthHandler) RefreshWithAuth(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", r)
		return
	}
	var req dto.RefreshRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	accessToken, refreshToken, expiresIn, err := h.authSvc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, r)
}

func toAuthResponse(user *model.User, accessToken, refreshToken string, expiresIn int64) *dto.AuthResponse {
	return &dto.AuthResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		Role:         user.Role,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}
