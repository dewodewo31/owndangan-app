package rsvp

import "github.com/owndangan/backend/internal/model"

type SubmitRSVPRequest struct {
	Token      string `json:"token" validate:"required,len=8"`
	Attendance string `json:"attendance" validate:"required,oneof=attending not_attending maybe"`
	GuestCount int    `json:"guest_count" validate:"min=1,max=10"`
	Message    string `json:"message,omitempty" validate:"max=1000"`
}

type RSVPRecap struct {
	TotalResponded  int `json:"total_responded"`
	Attending       int `json:"attending"`
	NotAttending    int `json:"not_attending"`
	Maybe           int `json:"maybe"`
	TotalGuestCount int `json:"total_guest_count"`
}

type RSVPResponse struct {
	ID           string `json:"id"`
	GuestID      string `json:"guest_id"`
	EventID      string `json:"event_id"`
	Attendance   string `json:"attendance"`
	GuestCount   int    `json:"guest_count"`
	Message      string `json:"message,omitempty"`
	SubmittedAt  string `json:"submitted_at"`
}

func ToRSVPResponse(rsvp *model.RSVP) RSVPResponse {
	return RSVPResponse{
		ID:          rsvp.ID.String(),
		GuestID:     rsvp.GuestID.String(),
		EventID:     rsvp.EventID.String(),
		Attendance:  rsvp.Attendance,
		GuestCount:  rsvp.GuestCount,
		Message:     rsvp.Message,
		SubmittedAt: rsvp.SubmittedAt.Format("2006-01-02T15:04:05Z"),
	}
}
