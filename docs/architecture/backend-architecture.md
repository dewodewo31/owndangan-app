# Backend Architecture

## Layer Structure

```
cmd/
  server/
    main.go           ← Entry point, wire everything together

internal/
  api/
    router.go          ← Route registration
    handler/           ← HTTP handlers per module
    middleware/        ← Auth, logging, CORS, rate limiter
    dto/               ← Request/response structs

  service/             ← Business logic per module

  repository/          ← Database access per model

  model/               ← GORM models

  config/              ← Configuration loading

  pkg/                 ← Shared utilities
    errors/            ← Custom error types
    validator/         ← Validation helpers
    jwt/               ← JWT utilities
    pagination/        ← Pagination helpers
    response/          ← Standard response helpers

migrations/            ← SQL migration files (versioned)

go.mod
go.sum
Makefile
Dockerfile
```

## Handler Responsibilities

- Parse HTTP request (path params, query params, body).
- Validate transport-level input (required fields, format).
- Call service layer.
- Map service result to HTTP response.
- Handle errors and return appropriate HTTP status.

```go
// Example handler shape
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateEventRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
        return
    }
    if err := h.validator.Struct(req); err != nil {
        response.ValidationError(w, err)
        return
    }
    event, err := h.service.Create(r.Context(), req)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create event")
        return
    }
    response.JSON(w, http.StatusCreated, event)
}
```

## Service Responsibilities

- Business logic and business validation.
- Orchestrating multiple repository calls.
- Transaction management across repositories.
- Feature entitlement checks.
- Authorization enforcement (ownership, role).

```go
// Example service shape
func (s *EventService) Create(ctx context.Context, req dto.CreateEventRequest) (*model.Event, error) {
    userID := auth.GetUserID(ctx)
    // Check subscription limits
    sub, err := s.subRepo.GetActiveByUser(ctx, userID)
    if err != nil { return nil, err }
    if sub == nil {
        return nil, errors.New("no active subscription")
    }
    // Check event limit for the plan
    // ...
    event := &model.Event{
        UserID: userID,
        Slug:   req.Slug,
        // ...
    }
    if err := s.eventRepo.Create(ctx, event); err != nil {
        return nil, err
    }
    return event, nil
}
```

## Repository Responsibilities

- Database CRUD operations.
- Query building with GORM.
- No business logic.
- Returns model structs.

```go
// Example repository shape
type EventRepository interface {
    Create(ctx context.Context, event *model.Event) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Event, error)
    GetBySlug(ctx context.Context, slug string) (*model.Event, error)
    Update(ctx context.Context, event *model.Event) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

## Middleware Stack

```
1. RequestID      ← Add unique request ID
2. Logger         ← Log request method, path, duration
3. CORS           ← Allow frontend origin
4. Rate Limiter   ← Throttle public endpoints
5. Auth           ← Parse and validate JWT
6. Authorizer     ← Check role/permissions (admin routes)
7. Handler        ← Business logic
8. Recovery       ← Recover from panics
```

## Dependency Injection

> TODO: Define DI approach (manual / wire / dig).

## Routing

> TODO: Define router library (chi / gorilla-mux / stdlib).

## Configuration

Loaded from environment variables at startup.

```go
type Config struct {
    Port     string
    Database DatabaseConfig
    JWT      JWTConfig
    Midtrans MidtransConfig
    Storage  StorageConfig
    Env      string
}
```

## Health Check

`GET /api/v1/health` — Returns 200 with database connection status.

## Related Documentation

- `docs/backend/handlers.md`
- `docs/backend/services.md`
- `docs/backend/repositories.md`
- `docs/backend/models.md`
- `docs/backend/middleware.md`
