package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateGuestRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,phone_id"`
	Category string `json:"category,omitempty" validate:"omitempty,oneof=family friend colleague other"`
	Note     string `json:"note,omitempty" validate:"max=1000"`
}

type UpdateGuestRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,phone_id"`
	Category *string `json:"category,omitempty" validate:"omitempty,oneof=family friend colleague other"`
	Note     *string `json:"note,omitempty" validate:"omitempty,max=1000"`
}

type GuestResponse struct {
	ID        uuid.UUID  `json:"id"`
	EventID   uuid.UUID  `json:"event_id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone,omitempty"`
	Category  string     `json:"category"`
	Note      string     `json:"note,omitempty"`
	RSVP      *GuestRSVP `json:"rsvp,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type GuestRSVP struct {
	Attendance  string `json:"attendance"`
	GuestCount  int    `json:"guest_count"`
	Message     string `json:"message,omitempty"`
	SubmittedAt string `json:"submitted_at"`
}

type ImportSummary struct {
	Imported int           `json:"imported"`
	Skipped  int           `json:"skipped"`
	Errors   []ImportError `json:"errors"`
}

type ImportError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}
