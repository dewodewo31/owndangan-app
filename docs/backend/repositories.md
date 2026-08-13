# Repositories

## Responsibilities

The repository layer encapsulates all database access. It provides a clean interface for the service layer and hides GORM implementation details.

**Repositories must:**
- Perform CRUD operations using GORM.
- Build complex queries using GORM's query builder (Scopes, Preload, Joins).
- Return model structs or slices.
- Accept transactions for atomic operations.

**Repositories must NOT:**
- Contain business logic or domain validation.
- Make decisions about what to save or not save.
- Call other repositories.
- Access HTTP context or authentication data.

## Interface Pattern

Every repository exposes an interface that services depend on, enabling unit testing with mocks:

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "github.com/yourorg/app-owndangan/internal/model"
)

type EventRepository interface {
    Create(ctx context.Context, event *model.Event) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
    GetBySlug(ctx context.Context, slug string) (*model.Event, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Event, error)
    ListByUser(ctx context.Context, userID uuid.UUID, page, perPage int) ([]model.Event, int64, error)
    Update(ctx context.Context, event *model.Event) error
    Delete(ctx context.Context, id uuid.UUID) error
    CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}
```

## Concrete Implementation

```go
type eventRepository struct {
    db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
    return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
    return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
    var event model.Event
    err := r.db.WithContext(ctx).
        Preload("Sections").
        Preload("Guests", func(db *gorm.DB) *gorm.DB {
            return db.Where("deleted_at IS NULL")
        }).
        First(&event, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // nil, nil — let service decide if that's an error
        }
        return nil, err
    }
    return &event, nil
}

func (r *eventRepository) ListByUser(ctx context.Context, userID uuid.UUID, page, perPage int) ([]model.Event, int64, error) {
    var events []model.Event
    var total int64

    query := r.db.WithContext(ctx).Model(&model.Event{}).
        Where("user_id = ? AND deleted_at IS NULL", userID)

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    offset := (page - 1) * perPage
    err := query.
        Order("created_at DESC").
        Offset(offset).
        Limit(perPage).
        Find(&events).Error
    return events, total, err
}
```

## Transaction Support

Repositories that need transaction support expose a `WithTx` method:

```go
func (r *eventRepository) WithTx(tx *gorm.DB) EventRepository {
    return &eventRepository{db: tx}
}
```

Services use this when composing multi-repository operations:

```go
err := s.db.Transaction(func(tx *gorm.DB) error {
    eventRepo := s.eventRepo.WithTx(tx)
    sectionRepo := s.sectionRepo.WithTx(tx)
    // ... all operations within the same transaction
    return nil
})
```

## Query Scopes

Extract reusable query filters as GORM scopes:

```go
func ActiveEvent(db *gorm.DB) *gorm.DB {
    return db.Where("status = ? AND deleted_at IS NULL", "published")
}

func BelongsToUser(userID uuid.UUID) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("user_id = ?", userID)
    }
}

// Usage
func (r *eventRepository) FindPublishedBySlug(ctx context.Context, slug string) (*model.Event, error) {
    var event model.Event
    err := r.db.WithContext(ctx).
        Scopes(ActiveEvent).
        Where("slug = ?", slug).
        First(&event).Error
    return &event, err
}
```
