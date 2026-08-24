package media

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

type ValidationConfig struct {
	MaxImageSize      int64
	MaxAudioSize      int64
	MaxVideoSize      int64
	AllowedImageTypes []string
	AllowedAudioTypes []string
	AllowedVideoTypes []string
}

func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		MaxImageSize: 5 * 1024 * 1024,
		MaxAudioSize: 10 * 1024 * 1024,
		MaxVideoSize: 50 * 1024 * 1024,
		AllowedImageTypes: []string{
			"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml",
		},
		AllowedAudioTypes: []string{
			"audio/mpeg", "audio/mp3", "audio/wav", "audio/ogg", "audio/aac", "audio/flac",
		},
		AllowedVideoTypes: []string{
			"video/mp4", "video/webm", "video/ogg", "video/quicktime",
		},
	}
}

type MediaType string

const (
	ImageMedia MediaType = "image"
	AudioMedia MediaType = "audio"
	VideoMedia MediaType = "video"
)

type Validator struct {
	config ValidationConfig
}

func NewValidator(config ValidationConfig) *Validator {
	return &Validator{config: config}
}

func (v *Validator) ValidateFile(filename string, size int64, contentType string) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	if ext == "" {
		return fmt.Errorf("file extension is required")
	}

	mediaType := v.detectMediaType(contentType)
	if mediaType == "" {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	if err := v.validateSize(size, mediaType); err != nil {
		return err
	}

	if err := v.validateExtension(ext, mediaType); err != nil {
		return err
	}

	if err := v.validateMIME(contentType, ext); err != nil {
		return err
	}

	return nil
}

func (v *Validator) detectMediaType(contentType string) MediaType {
	for _, t := range v.config.AllowedImageTypes {
		if t == contentType {
			return ImageMedia
		}
	}
	for _, t := range v.config.AllowedAudioTypes {
		if t == contentType {
			return AudioMedia
		}
	}
	for _, t := range v.config.AllowedVideoTypes {
		if t == contentType {
			return VideoMedia
		}
	}
	return ""
}

func (v *Validator) validateSize(size int64, mediaType MediaType) error {
	var maxSize int64
	switch mediaType {
	case ImageMedia:
		maxSize = v.config.MaxImageSize
	case AudioMedia:
		maxSize = v.config.MaxAudioSize
	case VideoMedia:
		maxSize = v.config.MaxVideoSize
	}
	if size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d for %s", size, maxSize, mediaType)
	}
	return nil
}

func (v *Validator) validateExtension(ext string, mediaType MediaType) error {
	var allowed []string
	switch mediaType {
	case ImageMedia:
		allowed = []string{"jpg", "jpeg", "png", "gif", "webp", "svg"}
	case AudioMedia:
		allowed = []string{"mp3", "wav", "ogg", "aac", "flac"}
	case VideoMedia:
		allowed = []string{"mp4", "webm", "ogg", "mov"}
	}
	for _, a := range allowed {
		if a == ext {
			return nil
		}
	}
	return fmt.Errorf("file extension .%s is not allowed for %s", ext, mediaType)
}

func (v *Validator) validateMIME(contentType, ext string) error {
	expectedTypes, err := mime.ExtensionsByType(contentType)
	if err == nil && len(expectedTypes) > 0 {
		for _, expected := range expectedTypes {
			if strings.TrimPrefix(expected, ".") == ext {
				return nil
			}
		}
	}

	mediaType := v.detectMediaType(contentType)
	switch mediaType {
	case ImageMedia:
		if contentType == "image/jpeg" && (ext == "jpg" || ext == "jpeg") {
			return nil
		}
		if contentType == "image/png" && ext == "png" {
			return nil
		}
		if contentType == "image/gif" && ext == "gif" {
			return nil
		}
		if contentType == "image/webp" && ext == "webp" {
			return nil
		}
	case AudioMedia:
		if (contentType == "audio/mpeg" || contentType == "audio/mp3") && ext == "mp3" {
			return nil
		}
		if contentType == "audio/wav" && ext == "wav" {
			return nil
		}
		if contentType == "audio/ogg" && ext == "ogg" {
			return nil
		}
	case VideoMedia:
		if contentType == "video/mp4" && ext == "mp4" {
			return nil
		}
		if contentType == "video/webm" && ext == "webm" {
			return nil
		}
	}

	return fmt.Errorf("content type %s does not match extension .%s", contentType, ext)
}

func SanitizeFilename(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	base = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, base)
	if base == "" {
		base = "file"
	}
	if len(base) > 100 {
		base = base[:100]
	}
	return base + strings.ToLower(ext)
}
