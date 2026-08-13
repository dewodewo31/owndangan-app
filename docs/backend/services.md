# Services

## Responsibilities

The service layer contains all business logic. It sits between handlers and repositories and enforces the rules of the application.

**Services must:**
- Implement business rules and domain validation.
- Orchestrate multiple repository calls within a single operation.
- Manage database transactions across repositories.
- Enforce authorization (ownership checks, feature entitlement).
- Return domain errors that the handler layer maps to HTTP responses.

**Services must NOT:**
- Parse HTTP requests or write HTTP responses.
- Access `*http.Request` or `http.ResponseWriter` directly.
- Know about request DTOs — accept primitive types or repository models.

## Transaction Management

Use GORM's transaction API for operations that span multiple repositories. Transactions should be managed at the service level, never at the handler or repository level.

```go
type EventService struct {
    db          *gorm.DB
    eventRepo   *repository.EventRepository
    sectionRepo *repository.EventSectionRepository
    subRepo     *repository.SubscriptionRepository
}

func (s *EventService) Create(ctx context.Context, req dto.CreateEventRequest) (*model.Event, error) {
    userID := auth.GetUserID(ctx)

    // Step 1: authorise — does the user have permission?
    if err := s.canCreateEvent(ctx, userID); err != nil {
        return nil, err
    }

    // Step 2: build the domain model from request data
    event := &model.Event{
        UserID:    userID,
        Title:     req.Title,
        Slug:      req.Slug,
        GroomName: req.GroomName,
        BrideName: req.BrideName,
        Status:    "draft",
    }

    // Step 3: persist within a transaction
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // Use transaction-scoped repositories
        eventRepo := s.eventRepo.WithTx(tx)
        sectionRepo := s.sectionRepo.WithTx(tx)

        if err := eventRepo.Create(ctx, event); err != nil {
            return err
        }
        // Create default sections for the new event
        section := &model.EventSection{
            EventID: event.ID,
            HeroEnabled:            true,
            CoupleEnabled:          true,
            EventDetailsEnabled:    true,
            GalleryEnabled:         true,
            RSVPEnabled:            true,
            GuestbookEnabled:       true,
            DigitalGiftsEnabled:    false,
        }
        return sectionRepo.Create(ctx, section)
    })
    if err != nil {
        return nil, err
    }
    return event, nil
}

func (s *EventService) canCreateEvent(ctx context.Context, userID uuid.UUID) error {
    // Check active subscription
    sub, err := s.subRepo.GetActiveByUser(ctx, userID)
    if err != nil {
        return fmt.Errorf("check subscription: %w", err)
    }
    if sub == nil {
        return errors.ErrPaymentRequired // requires subscription
    }
    // Check event count limit
    count, _ := s.eventRepo.CountByUser(ctx, userID)
    if count >= sub.Package.GuestLimit {
        return errors.ErrLimitExceeded
    }
    return nil
}
```

## Business Validation

Validation that goes beyond simple field checks (e.g., uniqueness, state transitions) belongs in services:

```go
func (s *EventService) Publish(ctx context.Context, eventID uuid.UUID) (*model.Event, error) {
    event, err := s.eventRepo.GetByID(ctx, eventID)
    if err != nil {
        return nil, err
    }
    // Business rule: cannot publish a draft if required fields are empty
    if event.GroomName == "" || event.BrideName == "" || event.WeddingDate == nil {
        return nil, errors.ErrValidationFailed("complete groom, bride, and wedding date first")
    }
    // Business rule: only draft events can be published
    if event.Status != "draft" {
        return nil, errors.ErrConflict("event is already " + event.Status)
    }
    event.Status = "published"
    now := time.Now()
    event.PublishedAt = &now
    if err := s.eventRepo.Update(ctx, event); err != nil {
        return nil, err
    }
    return event, nil
}
```

## Orchestration Pattern

When a single endpoint needs data from multiple sources, the service composes the result:

```go
func (s *EventService) GetDashboard(ctx context.Context) (*dto.DashboardResponse, error) {
    userID := auth.GetUserID(ctx)
    eventCount, _ := s.eventRepo.CountByUser(ctx, userID)
    guestCount, _ := s.guestRepo.CountByUser(ctx, userID)
    rsvpSummary, _ := s.rsvpRepo.SummaryByUser(ctx, userID)
    sub, _ := s.subRepo.GetActiveByUser(ctx, userID)
    return &dto.DashboardResponse{
        EventCount:  eventCount,
        GuestCount:  guestCount,
        RSVPSummary: rsvpSummary,
        Subscription: sub,
    }, nil
}
```
