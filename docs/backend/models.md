# Models

## Responsibilities

Model files define GORM structs that map to database tables. They contain struct tags for the ORM and JSON serialisation, but no business logic.

## Naming Conventions

- Struct name: singular PascalCase (`Event`, `GuestBookMessage`).
- Table name: GORM auto-derives `snake_case` plural from the struct name (`events`, `guest_book_messages`). Override only when necessary.
- Field names: PascalCase, matching Go conventions.
- JSON tags: `snake_case`.
- GORM tags: specify column name, constraints, indexes.

## Example Model Definitions

```go
package model

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/datatypes"
)

type Event struct {
    ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID           uuid.UUID      `gorm:"type:uuid;not null;index"                        json:"user_id"`
    TemplateID       *uuid.UUID     `gorm:"type:uuid"                                       json:"template_id"`
    Title            string         `gorm:"type:varchar(255);not null"                      json:"title"`
    Slug             string         `gorm:"type:varchar(100);uniqueIndex;not null"          json:"slug"`
    GroomName        string         `gorm:"type:varchar(255)"                               json:"groom_name"`
    BrideName        string         `gorm:"type:varchar(255)"                               json:"bride_name"`
    WeddingDate      *time.Time     `gorm:"type:date"                                       json:"wedding_date"`
    Status           string         `gorm:"type:varchar(20);not null;default:draft"         json:"status"`
    PublishedAt      *time.Time     `gorm:"type:timestamptz"                                json:"published_at"`
    ViewCount        int64          `gorm:"not null;default:0"                              json:"view_count"`
    CreatedAt        time.Time      `gorm:"type:timestamptz;not null;autoCreateTime"        json:"created_at"`
    UpdatedAt        time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime"        json:"updated_at"`
    DeletedAt        gorm.DeletedAt `gorm:"index"                                           json:"deleted_at"`

    // Relations
    Sections         *EventSection  `gorm:"foreignKey:EventID"                              json:"sections,omitempty"`
    Guests           []Guest        `gorm:"foreignKey:EventID"                              json:"guests,omitempty"`
    GalleryPhotos    []GalleryPhoto `gorm:"foreignKey:EventID"                              json:"gallery_photos,omitempty"`
}

type EventSection struct {
    ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    EventID               uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"                 json:"event_id"`
    HeroEnabled           bool      `gorm:"not null;default:true"                          json:"hero_enabled"`
    CoupleEnabled         bool      `gorm:"not null;default:true"                          json:"couple_enabled"`
    EventDetailsEnabled   bool      `gorm:"not null;default:true"                          json:"event_details_enabled"`
    GalleryEnabled        bool      `gorm:"not null;default:true"                          json:"gallery_enabled"`
    RSVPEnabled           bool      `gorm:"not null;default:true"                          json:"rsvp_enabled"`
    GuestbookEnabled      bool      `gorm:"not null;default:true"                          json:"guestbook_enabled"`
    DigitalGiftsEnabled   bool      `gorm:"not null;default:false"                         json:"digital_gifts_enabled"`
    CreatedAt             time.Time `gorm:"autoCreateTime"                                  json:"created_at"`
    UpdatedAt             time.Time `gorm:"autoUpdateTime"                                  json:"updated_at"`
}
```

## JSONB Columns

For flexible data like package features or digital gift bank accounts:

```go
type Package struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name      string         `gorm:"type:varchar(100);uniqueIndex;not null"`
    Code      string         `gorm:"type:varchar(50);uniqueIndex;not null"`
    Price     int64          `gorm:"type:bigint;not null"` // IDR in smallest unit
    Features  datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
    IsActive  bool           `gorm:"not null;default:true"`
}
```

## GORM Hooks

Use hooks sparingly — only for cross-cutting data concerns that must never be missed:

```go
// BeforeCreate sets a UUID if not already provided
func (e *Event) BeforeCreate(tx *gorm.DB) error {
    if e.ID == uuid.Nil {
        e.ID = uuid.New()
    }
    return nil
}

// BeforeSave slugifies the slug field
func (e *Event) BeforeSave(tx *gorm.DB) error {
    if e.Slug != "" {
        e.Slug = strings.TrimSpace(strings.ToLower(e.Slug))
    }
    return nil
}
```

Hooks should NOT contain business logic, send HTTP requests, or call external services.

## Soft Delete Models

Models that support soft delete embed `gorm.DeletedAt`:

```go
type User struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name         string         `gorm:"type:varchar(255);not null"`
    Email        string         `gorm:"type:varchar(255);uniqueIndex:idx_users_email,where:deleted_at IS NULL;not null"`
    PasswordHash string         `gorm:"type:varchar(255);not null"`
    Role         string         `gorm:"type:varchar(20);not null;default:user"`
    Status       string         `gorm:"type:varchar(20);not null;default:active"`
    CreatedAt    time.Time      `gorm:"autoCreateTime"`
    UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

For soft-deleted tables with unique constraints on email/slug, use partial unique indexes (`WHERE deleted_at IS NULL`) in the migration, not in GORM tags.
