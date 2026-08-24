package gallery

import (
	"bytes"

	"github.com/owndangan/backend/internal/model"
)

type AddPhotoRequest struct {
	Filename    string
	ContentType string
	Extension   string
	FileSize    int64
	Data        *bytes.Reader
	Caption     string
}

type PhotoResponse struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	ImageURL  string `json:"image_url"`
	Caption   string `json:"caption,omitempty"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

func ToPhotoResponse(photo *model.GalleryPhoto) PhotoResponse {
	return PhotoResponse{
		ID:        photo.ID.String(),
		EventID:   photo.EventID.String(),
		ImageURL:  photo.ImageURL,
		Caption:   photo.Caption,
		SortOrder: photo.SortOrder,
		CreatedAt: photo.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
