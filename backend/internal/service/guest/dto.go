package guest

import "github.com/owndangan/backend/internal/model"

type CreateGuestRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Category string `json:"category,omitempty" validate:"omitempty,max=50"`
	Note     string `json:"note,omitempty" validate:"omitempty,max=1000"`
}

type UpdateGuestRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Category *string `json:"category,omitempty" validate:"omitempty,max=50"`
	Note     *string `json:"note,omitempty" validate:"omitempty,max=1000"`
}

type ImportResult struct {
	Total    int      `json:"total"`
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

type GuestResponse struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone,omitempty"`
	Category  string `json:"category"`
	Note      string `json:"note,omitempty"`
	Token     string `json:"token,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToGuestResponse(guest *model.Guest) GuestResponse {
	return GuestResponse{
		ID:        guest.ID.String(),
		EventID:   guest.EventID.String(),
		Name:      guest.Name,
		Phone:     guest.Phone,
		Category:  guest.Category,
		Note:      guest.Note,
		Token:     guest.Token,
		CreatedAt: guest.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: guest.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
