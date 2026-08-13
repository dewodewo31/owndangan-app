# Handlers

## Responsibilities

Handlers are the HTTP transport layer. Their sole job is to convert incoming HTTP requests into service calls and map responses back to JSON.

**Handlers must:**
1. Parse the HTTP request (JSON body, path params, query params).
2. Validate transport-level input using `go-playground/validator`.
3. Extract authentication context (user ID, role) from the request context.
4. Call the service layer.
5. Map the service result to a standardised JSON response.
6. Handle errors and return the appropriate HTTP status code.

**Handlers must NOT:**
- Contain business logic or validation.
- Access the database directly.
- Make decisions based on user roles (delegate to service/middleware).

## Request Parsing Pattern

All request bodies use explicit DTO structs with validation tags:

```go
type CreateEventRequest struct {
    Title        string `json:"title"        validate:"required,max=255"`
    Slug         string `json:"slug"         validate:"required,slug,max=100"`
    GroomName    string `json:"groom_name"   validate:"required,max=255"`
    BrideName    string `json:"bride_name"   validate:"required,max=255"`
    WeddingDate  string `json:"wedding_date" validate:"required,date"`
}
```

## Handler Implementation

```go
package handler

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/yourorg/app-owndangan/internal/api/dto"
    "github.com/yourorg/app-owndangan/internal/service"
    "github.com/yourorg/app-owndangan/internal/pkg/response"
)

type EventHandler struct {
    service   *service.EventService
    validator *validator.Validate
}

func NewEventHandler(svc *service.EventService, v *validator.Validate) *EventHandler {
    return &EventHandler{service: svc, validator: v}
}

// POST /api/v1/events
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateEventRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "INVALID_JSON",
            "Request body is not valid JSON")
        return
    }
    if err := h.validator.Struct(req); err != nil {
        response.ValidationError(w, err)
        return
    }
    event, err := h.service.Create(r.Context(), req)
    if err != nil {
        response.FromError(w, err)
        return
    }
    response.JSON(w, http.StatusCreated, event)
}

// GET /api/v1/events/{id}
func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        response.Error(w, http.StatusBadRequest, "MISSING_ID", "Event ID is required")
        return
    }
    event, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        response.FromError(w, err)
        return
    }
    response.JSON(w, http.StatusOK, event)
}
```

## Error Handling

Handlers use `response.FromError()` which maps typed errors to HTTP responses. Never catch a service error and return a generic 500 — always let the error itself dictate the status code.

```go
// Good: let the error type determine the response
event, err := h.service.GetByID(r.Context(), id)
if err != nil {
    response.FromError(w, err)  // maps ErrNotFound → 404, ErrForbidden → 403, etc.
    return
}

// Bad: manual error handling in the handler
if errors.Is(err, service.ErrNotFound) {
    response.Error(w, 404, "NOT_FOUND", "Event not found")
    return
}
```

## Path and Query Parameters

Use `chi.URLParam` for path params and a pagination helper for query params:

```go
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
    page, perPage := pagination.FromRequest(r)
    events, total, err := h.service.List(r.Context(), auth.GetUserID(r.Context()), page, perPage)
    if err != nil {
        response.FromError(w, err)
        return
    }
    response.Paginated(w, http.StatusOK, events, page, perPage, total)
}
```
