package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
}

type SubscriptionResponse struct {
	ID        uuid.UUID    `json:"id"`
	Package   PackageBrief `json:"package"`
	Status    string       `json:"status"`
	StartAt   time.Time    `json:"start_at"`
	ExpiresAt *time.Time   `json:"expires_at,omitempty"`
}

type PackageBrief struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Code          string         `json:"code"`
	Price         int64          `json:"price"`
	GuestLimit    *int           `json:"guest_limit"`
	TemplateGroup string         `json:"template_group"`
	Features      datatypes.JSON `json:"features,omitempty"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
