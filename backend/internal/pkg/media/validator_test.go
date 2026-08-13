package media_test

import (
	"testing"

	"github.com/owndangan/backend/internal/pkg/media"
	"github.com/stretchr/testify/require"
)

func TestValidateFile_ValidImage(t *testing.T) {
	v := media.NewValidator(media.DefaultValidationConfig())
	err := v.ValidateFile("photo.jpg", 1024*1024, "image/jpeg")
	require.NoError(t, err)
}

func TestValidateFile_InvalidExtension(t *testing.T) {
	v := media.NewValidator(media.DefaultValidationConfig())
	err := v.ValidateFile("file.exe", 1024, "application/octet-stream")
	require.Error(t, err)
}

func TestValidateFile_TooLarge(t *testing.T) {
	v := media.NewValidator(media.DefaultValidationConfig())
	err := v.ValidateFile("photo.jpg", 100*1024*1024, "image/jpeg")
	require.Error(t, err)
}

func TestValidateFile_ValidAudio(t *testing.T) {
	v := media.NewValidator(media.DefaultValidationConfig())
	err := v.ValidateFile("song.mp3", 5*1024*1024, "audio/mpeg")
	require.NoError(t, err)
}

func TestValidateFile_ValidVideo(t *testing.T) {
	v := media.NewValidator(media.DefaultValidationConfig())
	err := v.ValidateFile("video.mp4", 20*1024*1024, "video/mp4")
	require.NoError(t, err)
}

func TestValidateFile_InvalidMIME(t *testing.T) {
	v := media.NewValidator(media.DefaultValidationConfig())
	err := v.ValidateFile("photo.jpg", 1024, "application/pdf")
	require.Error(t, err)
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"photo.jpg", "photo.jpg"},
		{"my photo.jpg", "my_photo.jpg"},
		{"../../../etc/passwd", "_________etc_passwd"},
		{"file@name#.png", "file_name_.png"},
		{"", ""},
	}

	for _, tt := range tests {
		result := media.SanitizeFilename(tt.input)
		if tt.expected == "" {
			require.NotEmpty(t, result)
		} else {
			require.Equal(t, tt.expected, result)
		}
	}
}
