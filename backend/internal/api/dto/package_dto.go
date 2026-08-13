package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PackageResponse struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Code          string         `json:"code"`
	Price         int64          `json:"price"`
	DurationDays  *int           `json:"duration_days"`
	GuestLimit    *int           `json:"guest_limit"`
	TemplateGroup string         `json:"template_group"`
	Features      datatypes.JSON `json:"features"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at,omitempty"`
}

type CreatePackageRequest struct {
	Name          string         `json:"name" validate:"required,max=100"`
	Code          string         `json:"code" validate:"required,max=50"`
	Price         int64          `json:"price" validate:"required,gte=0"`
	DurationDays  *int           `json:"duration_days"`
	GuestLimit    *int           `json:"guest_limit"`
	TemplateGroup string         `json:"template_group" validate:"required,oneof=standard premium all"`
	Features      datatypes.JSON `json:"features"`
	IsActive      bool           `json:"is_active"`
}

type UpdatePackageRequest struct {
	Name          *string        `json:"name,omitempty"`
	Code          *string        `json:"code,omitempty"`
	Price         *int64         `json:"price,omitempty"`
	DurationDays  *int           `json:"duration_days"`
	GuestLimit    *int           `json:"guest_limit"`
	TemplateGroup *string        `json:"template_group,omitempty"`
	Features      datatypes.JSON `json:"features,omitempty"`
	IsActive      *bool          `json:"is_active,omitempty"`
}
