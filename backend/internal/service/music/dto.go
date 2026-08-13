package music

import (
	"bytes"

	"github.com/owndangan/backend/internal/model"
)

type SetMusicRequest struct {
	Title       string
	Filename    string
	ContentType string
	Extension   string
	FileSize    int64
	Data        *bytes.Reader
	IsPreset    bool
	PresetName  string
}

type MusicResponse struct {
	ID       string `json:"id"`
	EventID  string `json:"event_id,omitempty"`
	Title    string `json:"title"`
	FileURL  string `json:"file_url,omitempty"`
	Preset   string `json:"preset,omitempty"`
	IsPreset bool   `json:"is_preset"`
}

func ToMusicResponse(m *model.Music) MusicResponse {
	return MusicResponse{
		ID:       m.ID.String(),
		EventID:  m.EventID.String(),
		Title:    m.Title,
		FileURL:  m.FileURL,
		Preset:   m.Preset,
		IsPreset: m.IsPreset,
	}
}
