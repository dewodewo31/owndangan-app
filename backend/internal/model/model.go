package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Guest category constants. Unknown/empty values normalize to GuestCategoryLainnya.
const (
	GuestCategoryKeluarga   = "keluarga"
	GuestCategoryTeman      = "teman"
	GuestCategoryRekanKerja = "rekan_kerja"
	GuestCategoryTetangga   = "tetangga"
	GuestCategoryLainnya    = "lainnya"
)

// NormalizeGuestCategory maps a raw category value to a known value, defaulting
// unknown or empty input to GuestCategoryLainnya.
func NormalizeGuestCategory(c string) string {
	switch strings.ToLower(strings.TrimSpace(c)) {
	case GuestCategoryKeluarga, GuestCategoryTeman, GuestCategoryRekanKerja, GuestCategoryTetangga, GuestCategoryLainnya:
		return strings.ToLower(strings.TrimSpace(c))
	default:
		return GuestCategoryLainnya
	}
}

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Email        string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Phone        string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Role         string         `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	Status       string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	AvatarURL    string         `gorm:"type:text" json:"avatar_url,omitempty"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Subscriptions []Subscription `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	Events        []Event        `gorm:"foreignKey:UserID" json:"events,omitempty"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	IsRevoked bool      `gorm:"not null;default:false" json:"is_revoked"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Package struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	Code          string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Price         int64          `gorm:"type:bigint;not null" json:"price"`
	DurationDays  *int           `gorm:"type:integer" json:"duration_days"`
	GuestLimit    *int           `gorm:"type:integer" json:"guest_limit"`
	TemplateGroup string         `gorm:"type:varchar(50);not null;default:'standard'" json:"template_group"`
	Features      datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"features"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`

	Subscriptions []Subscription `gorm:"foreignKey:PackageID" json:"subscriptions,omitempty"`
	Transactions  []Transaction  `gorm:"foreignKey:PackageID" json:"transactions,omitempty"`
}

type Subscription struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	PackageID     uuid.UUID      `gorm:"type:uuid;not null" json:"package_id"`
	TransactionID *uuid.UUID     `gorm:"type:uuid" json:"transaction_id,omitempty"`
	Status        string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	StartAt       time.Time      `gorm:"type:timestamptz;not null" json:"start_at"`
	ExpiresAt     *time.Time     `gorm:"type:timestamptz" json:"expires_at,omitempty"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Package     Package      `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	Transaction *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
}

type Transaction struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	PackageID        uuid.UUID      `gorm:"type:uuid;not null" json:"package_id"`
	OrderID          string         `gorm:"type:varchar(100);unique;not null" json:"order_id"`
	SnapToken        string         `gorm:"type:varchar(255)" json:"snap_token,omitempty"`
	GrossAmount      int64          `gorm:"type:bigint;not null" json:"gross_amount"`
	Status           string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	PaymentType      string         `gorm:"type:varchar(50)" json:"payment_type,omitempty"`
	TransactionTime  *time.Time     `gorm:"type:timestamptz" json:"transaction_time,omitempty"`
	SettlementTime   *time.Time     `gorm:"type:timestamptz" json:"settlement_time,omitempty"`
	MidtransResponse datatypes.JSON `gorm:"type:jsonb" json:"-"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Package Package `gorm:"foreignKey:PackageID" json:"package,omitempty"`
}

type Event struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TemplateID       *uuid.UUID     `gorm:"type:uuid" json:"template_id,omitempty"`
	Title            string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug             string         `gorm:"type:varchar(100);unique;not null" json:"slug"`
	CoupleName       string         `gorm:"type:varchar(255)" json:"couple_name,omitempty"`
	GroomName        string         `gorm:"type:varchar(255)" json:"groom_name,omitempty"`
	BrideName        string         `gorm:"type:varchar(255)" json:"bride_name,omitempty"`
	GroomParents     string         `gorm:"type:varchar(255)" json:"groom_parents,omitempty"`
	BrideParents     string         `gorm:"type:varchar(255)" json:"bride_parents,omitempty"`
	WeddingDate      *time.Time     `gorm:"type:date" json:"wedding_date,omitempty"`
	WeddingTime      string         `gorm:"type:varchar(20)" json:"wedding_time,omitempty"`
	CeremonyVenue    string         `gorm:"type:text" json:"ceremony_venue,omitempty"`
	CeremonyAddress  string         `gorm:"type:text" json:"ceremony_address,omitempty"`
	CeremonyMapURL   string         `gorm:"type:text" json:"ceremony_map_url,omitempty"`
	ReceptionVenue   string         `gorm:"type:text" json:"reception_venue,omitempty"`
	ReceptionAddress string         `gorm:"type:text" json:"reception_address,omitempty"`
	ReceptionMapURL  string         `gorm:"type:text" json:"reception_map_url,omitempty"`
	MusicURL         string         `gorm:"type:varchar(255)" json:"music_url,omitempty"`
	VideoURL         string         `gorm:"type:text" json:"video_url,omitempty"`
	Status           string         `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	PublishedAt      *time.Time     `gorm:"type:timestamptz" json:"published_at,omitempty"`
	ViewCount        int64          `gorm:"not null;default:0" json:"view_count"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Sections      *EventSection  `gorm:"foreignKey:EventID" json:"sections,omitempty"`
	Guests        []Guest        `gorm:"foreignKey:EventID" json:"guests,omitempty"`
	GalleryPhotos []GalleryPhoto `gorm:"foreignKey:EventID" json:"gallery_photos,omitempty"`
	DigitalGift   *DigitalGift   `gorm:"foreignKey:EventID" json:"digital_gift,omitempty"`
	LoveStories   []LoveStory    `gorm:"foreignKey:EventID" json:"love_stories,omitempty"`
}

type EventSection struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID             uuid.UUID  `gorm:"type:uuid;unique;not null" json:"event_id"`
	HeroEnabled         bool       `gorm:"not null;default:true" json:"hero_enabled"`
	CoupleEnabled       bool       `gorm:"not null;default:true" json:"couple_enabled"`
	EventDetailsEnabled bool       `gorm:"not null;default:true" json:"event_details_enabled"`
	GalleryEnabled      bool       `gorm:"not null;default:true" json:"gallery_enabled"`
	VideoEnabled        bool       `gorm:"not null;default:false" json:"video_enabled"`
	MusicID             *uuid.UUID `gorm:"type:uuid" json:"music_id,omitempty"`
	RSVPEnabled         bool       `gorm:"not null;default:true" json:"rsvp_enabled"`
	GuestbookEnabled    bool       `gorm:"not null;default:true" json:"guestbook_enabled"`
	LoveStoryEnabled    bool       `gorm:"not null;default:false" json:"love_story_enabled"`
	DigitalGiftsEnabled bool       `gorm:"not null;default:false" json:"digital_gifts_enabled"`
	DressCode           string     `gorm:"type:varchar(500);default:''" json:"dress_code,omitempty"`
	ClosingMessage      string     `gorm:"type:text;default:''" json:"closing_message,omitempty"`
	OpeningMessage      string     `gorm:"type:text;default:''" json:"opening_message,omitempty"`
	VerseEnabled        bool       `gorm:"not null;default:false" json:"verse_enabled"`
	VerseReligion       string     `gorm:"type:varchar(20);default:'quran'" json:"verse_religion"`
	VerseText           string     `gorm:"type:text;default:''" json:"verse_text,omitempty"`
	VerseSource         string     `gorm:"type:varchar(255);default:''" json:"verse_source,omitempty"`
	CreatedAt           time.Time  `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
}

type Template struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	GroupName    string         `gorm:"type:varchar(50);not null" json:"group_name"`
	ThumbnailURL string         `gorm:"type:text" json:"thumbnail_url,omitempty"`
	CSSConfig    datatypes.JSON `gorm:"type:jsonb" json:"css_config,omitempty"`
	LayoutConfig datatypes.JSON `gorm:"type:jsonb" json:"layout_config,omitempty"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
}

type Guest struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"event_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Category  string         `gorm:"type:varchar(32);not null;default:'lainnya'" json:"category"`
	Note      string         `gorm:"type:text" json:"note,omitempty"`
	Token     string         `gorm:"type:varchar(100);unique;not null" json:"-"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RSVP *RSVP `gorm:"foreignKey:GuestID" json:"rsvp,omitempty"`
}

type RSVP struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	GuestID     uuid.UUID `gorm:"type:uuid;unique;not null" json:"guest_id"`
	EventID     uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	Attendance  string    `gorm:"type:varchar(20);not null" json:"attendance"`
	GuestCount  int       `gorm:"not null;default:1" json:"guest_count"`
	Message     string    `gorm:"type:text" json:"message,omitempty"`
	SubmittedAt time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"submitted_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
}

type GuestbookMessage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID    uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	IsApproved bool      `gorm:"not null;default:false" json:"is_approved"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
}

type DigitalGift struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID      uuid.UUID      `gorm:"type:uuid;unique;not null" json:"event_id"`
	BankAccounts datatypes.JSON `gorm:"type:jsonb" json:"bank_accounts,omitempty"`
	EWallet      datatypes.JSON `gorm:"type:jsonb" json:"ewallet,omitempty"`
	QRISImageURL string         `gorm:"type:text" json:"qris_image_url,omitempty"`
	GiftMessage  string         `gorm:"type:text" json:"gift_message,omitempty"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
}

type GalleryPhoto struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	ImageURL  string    `gorm:"type:text;not null" json:"image_url"`
	Caption   string    `gorm:"type:varchar(255)" json:"caption,omitempty"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

type LoveStory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Story     string    `gorm:"type:text;not null" json:"story"`
	Year      string    `gorm:"type:varchar(10)" json:"year,omitempty"`
	Date      string    `gorm:"type:varchar(64)" json:"date,omitempty"`
	ImageURL  string    `gorm:"type:text" json:"image_url,omitempty"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
}

type Music struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID   *uuid.UUID `gorm:"type:uuid" json:"event_id,omitempty"`
	Title     string     `gorm:"type:varchar(255);not null" json:"title"`
	FileURL   string     `gorm:"type:text" json:"file_url,omitempty"`
	Preset    string     `gorm:"type:varchar(100)" json:"preset,omitempty"`
	IsPreset  bool       `gorm:"not null;default:false" json:"is_preset"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

type AuditLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Action     string         `gorm:"type:varchar(100);not null" json:"action"`
	EntityType string         `gorm:"type:varchar(50)" json:"entity_type,omitempty"`
	EntityID   *uuid.UUID     `gorm:"type:uuid" json:"entity_id,omitempty"`
	Metadata   datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	IPAddress  string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

type AnalyticsEvent struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EventID   *uuid.UUID     `gorm:"type:uuid;index" json:"event_id,omitempty"`
	EventType string         `gorm:"type:varchar(50);not null" json:"event_type"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	IPAddress string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent string         `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

type WebhookIdempotency struct {
	TransactionID string    `gorm:"type:varchar(100);primaryKey" json:"transaction_id"`
	OrderID       string    `gorm:"type:varchar(100);not null" json:"order_id"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status"`
	ProcessedAt   time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"processed_at"`
}
