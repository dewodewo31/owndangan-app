package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PublishResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	PublishedAt string    `json:"published_at"`
	PublicURL   string    `json:"public_url"`
}

type CreateEventRequest struct {
	Title            string `json:"title" validate:"required,max=255"`
	Slug             string `json:"slug,omitempty" validate:"omitempty,slug,max=100"`
	CoupleName       string `json:"couple_name,omitempty" validate:"max=255"`
	GroomName        string `json:"groom_name,omitempty" validate:"max=255"`
	BrideName        string `json:"bride_name,omitempty" validate:"max=255"`
	GroomParents     string `json:"groom_parents,omitempty" validate:"max=255"`
	BrideParents     string `json:"bride_parents,omitempty" validate:"max=255"`
	WeddingDate      string `json:"wedding_date,omitempty" validate:"omitempty,date"`
	WeddingTime      string `json:"wedding_time,omitempty" validate:"omitempty,max=20"`
	CeremonyVenue    string `json:"ceremony_venue,omitempty" validate:"max=1000"`
	CeremonyAddress  string `json:"ceremony_address,omitempty" validate:"max=2000"`
	CeremonyMapURL   string `json:"ceremony_map_url,omitempty" validate:"omitempty,url"`
	ReceptionVenue   string `json:"reception_venue,omitempty" validate:"max=1000"`
	ReceptionAddress string `json:"reception_address,omitempty" validate:"max=2000"`
	ReceptionMapURL  string `json:"reception_map_url,omitempty" validate:"omitempty,url"`
}

type UpdateEventRequest struct {
	Title            *string    `json:"title,omitempty" validate:"omitempty,max=255"`
	CoupleName       *string    `json:"couple_name,omitempty" validate:"omitempty,max=255"`
	GroomName        *string    `json:"groom_name,omitempty" validate:"omitempty,max=255"`
	BrideName        *string    `json:"bride_name,omitempty" validate:"omitempty,max=255"`
	GroomParents     *string    `json:"groom_parents,omitempty" validate:"omitempty,max=255"`
	BrideParents     *string    `json:"bride_parents,omitempty" validate:"omitempty,max=255"`
	WeddingDate      *time.Time `json:"wedding_date,omitempty"`
	WeddingTime      *string    `json:"wedding_time,omitempty" validate:"omitempty,max=20"`
	CeremonyVenue    *string    `json:"ceremony_venue,omitempty" validate:"omitempty,max=1000"`
	CeremonyAddress  *string    `json:"ceremony_address,omitempty" validate:"omitempty,max=2000"`
	CeremonyMapURL   *string    `json:"ceremony_map_url,omitempty" validate:"omitempty,url"`
	ReceptionVenue   *string    `json:"reception_venue,omitempty" validate:"omitempty,max=1000"`
	ReceptionAddress *string    `json:"reception_address,omitempty" validate:"omitempty,max=2000"`
	ReceptionMapURL  *string    `json:"reception_map_url,omitempty" validate:"omitempty,url"`
	TemplateID       *uuid.UUID `json:"template_id,omitempty"`
}

type EventResponse struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	TemplateID       *uuid.UUID `json:"template_id,omitempty"`
	Title            string     `json:"title"`
	Slug             string     `json:"slug"`
	CoupleName       string     `json:"couple_name,omitempty"`
	GroomName        string     `json:"groom_name,omitempty"`
	BrideName        string     `json:"bride_name,omitempty"`
	GroomParents     string     `json:"groom_parents,omitempty"`
	BrideParents     string     `json:"bride_parents,omitempty"`
	WeddingDate      *time.Time `json:"wedding_date,omitempty"`
	WeddingTime      string     `json:"wedding_time,omitempty"`
	CeremonyVenue    string     `json:"ceremony_venue,omitempty"`
	CeremonyAddress  string     `json:"ceremony_address,omitempty"`
	CeremonyMapURL   string     `json:"ceremony_map_url,omitempty"`
	ReceptionVenue   string     `json:"reception_venue,omitempty"`
	ReceptionAddress string     `json:"reception_address,omitempty"`
	ReceptionMapURL  string     `json:"reception_map_url,omitempty"`
	MusicURL         string     `json:"music_url,omitempty"`
	Status           string     `json:"status"`
	PublishedAt      *time.Time `json:"published_at,omitempty"`
	ViewCount        int64      `json:"view_count"`
	GuestCount       int64      `json:"guest_count"`
	RsvpCount        int64      `json:"rsvp_count"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type PublicEventResponse struct {
	Event       PublicEventDetail     `json:"event"`
	Template    *TemplatePreview      `json:"template,omitempty"`
	Sections    *EventSectionsDTO     `json:"sections,omitempty"`
	Gallery     []GalleryPhotoDTO     `json:"gallery,omitempty"`
	Guestbook   []GuestbookPublicDTO  `json:"guestbook,omitempty"`
	DigitalGift *DigitalGiftPublicDTO `json:"digital_gift,omitempty"`
}

type PublicEventDetail struct {
	Title            string `json:"title"`
	CoupleName       string `json:"couple_name,omitempty"`
	GroomName        string `json:"groom_name,omitempty"`
	BrideName        string `json:"bride_name,omitempty"`
	GroomParents     string `json:"groom_parents,omitempty"`
	BrideParents     string `json:"bride_parents,omitempty"`
	WeddingDate      string `json:"wedding_date,omitempty"`
	WeddingTime      string `json:"wedding_time,omitempty"`
	CeremonyVenue    string `json:"ceremony_venue,omitempty"`
	CeremonyAddress  string `json:"ceremony_address,omitempty"`
	CeremonyMapURL   string `json:"ceremony_map_url,omitempty"`
	ReceptionVenue   string `json:"reception_venue,omitempty"`
	ReceptionAddress string `json:"reception_address,omitempty"`
	ReceptionMapURL  string `json:"reception_map_url,omitempty"`
	ViewCount        int64  `json:"view_count"`
}

type TemplatePreview struct {
	Name         string         `json:"name"`
	GroupName    string         `json:"group_name"`
	CSSConfig    datatypes.JSON `json:"css_config,omitempty"`
	LayoutConfig datatypes.JSON `json:"layout_config,omitempty"`
}

type EventSectionsDTO struct {
	HeroEnabled         bool      `json:"hero_enabled"`
	CoupleEnabled       bool      `json:"couple_enabled"`
	EventDetailsEnabled bool      `json:"event_details_enabled"`
	GalleryEnabled      bool      `json:"gallery_enabled"`
	VideoEnabled        bool      `json:"video_enabled"`
	Music               *MusicDTO `json:"music,omitempty"`
	RSVPEnabled         bool      `json:"rsvp_enabled"`
	GuestbookEnabled    bool      `json:"guestbook_enabled"`
	DigitalGiftsEnabled bool      `json:"digital_gifts_enabled"`
	DressCode           string    `json:"dress_code,omitempty"`
	ClosingMessage      string    `json:"closing_message,omitempty"`
	OpeningMessage      string    `json:"opening_message,omitempty"`
	VerseEnabled        bool      `json:"verse_enabled"`
	VerseReligion       string    `json:"verse_religion"`
	VerseText           string    `json:"verse_text,omitempty"`
	VerseSource         string    `json:"verse_source,omitempty"`
}

type MusicDTO struct {
	Title    string `json:"title"`
	FileURL  string `json:"file_url,omitempty"`
	Preset   string `json:"preset,omitempty"`
	IsPreset bool   `json:"is_preset"`
}

type GalleryPhotoDTO struct {
	ImageURL  string `json:"image_url"`
	Caption   string `json:"caption,omitempty"`
	SortOrder int    `json:"sort_order"`
}

type GuestbookPublicDTO struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type DigitalGiftPublicDTO struct {
	BankAccounts []map[string]interface{} `json:"bank_accounts,omitempty"`
	EWallet      map[string]interface{}   `json:"ewallet,omitempty"`
	QRISImageURL string                   `json:"qris_image_url,omitempty"`
	GiftMessage  string                   `json:"gift_message,omitempty"`
}

type UpdateSectionsRequest struct {
	HeroEnabled         *bool      `json:"hero_enabled,omitempty"`
	CoupleEnabled       *bool      `json:"couple_enabled,omitempty"`
	EventDetailsEnabled *bool      `json:"event_details_enabled,omitempty"`
	GalleryEnabled      *bool      `json:"gallery_enabled,omitempty"`
	VideoEnabled        *bool      `json:"video_enabled,omitempty"`
	MusicID             *uuid.UUID `json:"music_id,omitempty"`
	RSVPEnabled         *bool      `json:"rsvp_enabled,omitempty"`
	GuestbookEnabled    *bool      `json:"guestbook_enabled,omitempty"`
	DigitalGiftsEnabled *bool      `json:"digital_gifts_enabled,omitempty"`
	DressCode           *string    `json:"dress_code,omitempty"`
	ClosingMessage      *string    `json:"closing_message,omitempty"`
	OpeningMessage      *string    `json:"opening_message,omitempty"`
	VerseEnabled        *bool      `json:"verse_enabled,omitempty"`
	VerseReligion       *string    `json:"verse_religion,omitempty"`
	VerseText           *string    `json:"verse_text,omitempty"`
	VerseSource         *string    `json:"verse_source,omitempty"`
}

type SectionsResponse struct {
	ID                  uuid.UUID  `json:"id"`
	EventID             uuid.UUID  `json:"event_id"`
	HeroEnabled         bool       `json:"hero_enabled"`
	CoupleEnabled       bool       `json:"couple_enabled"`
	EventDetailsEnabled bool       `json:"event_details_enabled"`
	GalleryEnabled      bool       `json:"gallery_enabled"`
	VideoEnabled        bool       `json:"video_enabled"`
	MusicID             *uuid.UUID `json:"music_id,omitempty"`
	RSVPEnabled         bool       `json:"rsvp_enabled"`
	GuestbookEnabled    bool       `json:"guestbook_enabled"`
	DigitalGiftsEnabled bool       `json:"digital_gifts_enabled"`
	DressCode           string     `json:"dress_code,omitempty"`
	ClosingMessage      string     `json:"closing_message,omitempty"`
	OpeningMessage      string     `json:"opening_message,omitempty"`
	VerseEnabled        bool       `json:"verse_enabled"`
	VerseReligion       string     `json:"verse_religion"`
	VerseText           string     `json:"verse_text,omitempty"`
	VerseSource         string     `json:"verse_source,omitempty"`
}

type UpdateDigitalGiftRequest struct {
	BankAccounts []map[string]interface{} `json:"bank_accounts,omitempty"`
	EWallet      map[string]interface{}   `json:"ewallet,omitempty"`
	QRISImageURL *string                  `json:"qris_image_url,omitempty"`
	GiftMessage  *string                  `json:"gift_message,omitempty"`
}

type DigitalGiftResponse struct {
	ID           uuid.UUID                `json:"id"`
	EventID      uuid.UUID                `json:"event_id"`
	BankAccounts []map[string]interface{} `json:"bank_accounts,omitempty"`
	EWallet      map[string]interface{}   `json:"ewallet,omitempty"`
	QRISImageURL string                   `json:"qris_image_url,omitempty"`
	GiftMessage  string                   `json:"gift_message,omitempty"`
}

type ReorderGalleryRequest struct {
	Photos []GalleryPhotoOrder `json:"photos"`
}

type GalleryPhotoOrder struct {
	ID        uuid.UUID `json:"id"`
	SortOrder int       `json:"sort_order"`
}

type GalleryPhotoResponse struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	Caption   string    `json:"caption,omitempty"`
	SortOrder int       `json:"sort_order"`
}

type MusicResponse struct {
	ID       uuid.UUID  `json:"id"`
	EventID  *uuid.UUID `json:"event_id,omitempty"`
	Title    string     `json:"title"`
	FileURL  string     `json:"file_url,omitempty"`
	Preset   string     `json:"preset,omitempty"`
	IsPreset bool       `json:"is_preset"`
}

type TemplateSummary struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	GroupName    string         `json:"group_name"`
	ThumbnailURL string         `json:"thumbnail_url,omitempty"`
	CSSConfig    datatypes.JSON `json:"css_config,omitempty"`
	LayoutConfig datatypes.JSON `json:"layout_config,omitempty"`
}
