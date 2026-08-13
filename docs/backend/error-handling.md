# Error Handling

## Error Types

The application uses a custom error type that carries an HTTP status code, an error code string, and an internal error message. This allows handlers to map errors to responses without understanding business logic.

```go
package errors

import "net/http"

type AppError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    HTTPStatus int    `json:"-"`
    Err        error  `json:"-"` // wrapped internal error
}

func (e *AppError) Error() string {
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}
```

## Predefined Error Sentinels

```go
var (
    ErrNotFound       = &AppError{Code: "NOT_FOUND",       HTTPStatus: http.StatusNotFound,       Message: "Resource not found"}
    ErrUnauthorized   = &AppError{Code: "UNAUTHORIZED",    HTTPStatus: http.StatusUnauthorized,   Message: "Authentication required"}
    ErrForbidden      = &AppError{Code: "FORBIDDEN",       HTTPStatus: http.StatusForbidden,      Message: "Insufficient permissions"}
    ErrConflict       = &AppError{Code: "CONFLICT",        HTTPStatus: http.StatusConflict,       Message: "Resource already exists"}
    ErrRateLimited    = &AppError{Code: "RATE_LIMITED",    HTTPStatus: http.StatusTooManyRequests, Message: "Too many requests"}
    ErrPaymentRequired = &AppError{Code: "PAYMENT_REQUIRED", HTTPStatus: http.StatusPaymentRequired, Message: "Active subscription required"}
    ErrLimitExceeded  = &AppError{Code: "LIMIT_EXCEEDED",  HTTPStatus: http.StatusUnprocessableEntity, Message: "Plan limit exceeded"}
)
```

## Creating Errors with Context

```go
// Wrap an internal error with context
func NotFound(entity string, id string) *AppError {
    return &AppError{
        Code:       "NOT_FOUND",
        HTTPStatus: http.StatusNotFound,
        Message:    fmt.Sprintf("%s with ID %s not found", entity, id),
    }
}

// Validation error with field-level details
func ValidationFailed(details map[string]string) *AppError {
    return &AppError{
        Code:       "VALIDATION_ERROR",
        HTTPStatus: http.StatusUnprocessableEntity,
        Message:    "Validation failed",
    }
}
```

## Error Wrapping in Services

Services wrap errors to add context. The original error type is preserved via `fmt.Errorf("... %w", err)`:

```go
func (s *EventService) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
    event, err := s.eventRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get event %s: %w", id, err)
    }
    if event == nil {
        return nil, fmt.Errorf("event %s: %w", id, errors.ErrNotFound)
    }
    return event, nil
}
```

## HTTP Status Mapping (response.FromError)

The response package provides `FromError` which uses `errors.As` to find an `AppError` in the chain:

```go
package response

func FromError(w http.ResponseWriter, err error) {
    var appErr *errors.AppError
    if errors.As(err, &appErr) {
        JSON(w, appErr.HTTPStatus, ErrorPayload{
            Success: false,
            Error: ErrorDetail{
                Code:    appErr.Code,
                Message: appErr.Message,
            },
        })
        return
    }
    // Fallback for unrecognised errors (should not happen in production)
    JSON(w, http.StatusInternalServerError, ErrorPayload{
        Success: false,
        Error: ErrorDetail{
            Code:    "INTERNAL_ERROR",
            Message: "An unexpected error occurred",
        },
    })
}
```

## Error Handler Registration

The chi router can register a custom error handler for unmatched routes:

```go
r.NotFound(func(w http.ResponseWriter, r *http.Request) {
    response.Error(w, http.StatusNotFound, "NOT_FOUND", "Route not found")
})

r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
    response.Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
        fmt.Sprintf("Method %s not allowed on %s", r.Method, r.URL.Path))
})
```

## Error Response Format

All errors return the standard envelope:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "Must be a valid email address"
    }
  },
  "meta": {
    "request_id": "req-abc123"
  }
}
```

Internal error details (stack traces, internal IDs, full error messages) must never be leaked in production API responses. They are logged separately via the structured logger.
