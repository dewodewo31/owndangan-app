package guestbook

import "github.com/owndangan/backend/internal/model"

type SubmitMessageRequest struct {
	Name    string `json:"name" validate:"required,max=255"`
	Message string `json:"message" validate:"required,max=2000"`
}

type GuestbookResponse struct {
	ID         string `json:"id"`
	EventID    string `json:"event_id"`
	Name       string `json:"name"`
	Message    string `json:"message"`
	IsApproved bool   `json:"is_approved"`
	CreatedAt  string `json:"created_at"`
}

func ToGuestbookResponse(msg *model.GuestbookMessage) GuestbookResponse {
	return GuestbookResponse{
		ID:         msg.ID.String(),
		EventID:    msg.EventID.String(),
		Name:       msg.Name,
		Message:    msg.Message,
		IsApproved: msg.IsApproved,
		CreatedAt:  msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
