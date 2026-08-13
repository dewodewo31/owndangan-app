# Validation

## Approach

Input validation happens at two layers:

1. **Handler (transport validation):** Struct tag validation using `go-playground/validator`. Checks field existence, format, length, and type constraints.
2. **Service (business validation):** Explicit Go code that checks uniqueness, state transitions, referential integrity, and subscription limits.

This document covers transport validation only. See `services.md` for business validation.

## Setup

```go
package validator

import (
    "github.com/go-playground/validator/v10"
    "regexp"
    "strings"
    "time"
)

func New() *validator.Validate {
    v := validator.New()

    // Use JSON tag names in error output instead of field names
    v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := fld.Tag.Get("json")
        if name == "" {
            return fld.Name
        }
        return strings.Split(name, ",")[0]
    })

    // Register custom validators
    _ = v.RegisterValidation("slug", slugValidator)
    _ = v.RegisterValidation("date", dateValidator)
    _ = v.RegisterValidation("phone_id", phoneIDValidator)
    _ = v.RegisterValidation("uuid", uuidValidator)

    return v
}
```

## Built-in Validators Used

| Tag | Usage | Description |
|-----|-------|-------------|
| `required` | All required fields | Fails on zero-value |
| `max=255` | Strings | Max character length |
| `min=8` | Password, phone | Min character length |
| `email` | Email fields | RFC 5322 email format |
| `oneof=active suspended` | Status enums | Must match allowed values |
| `gte=0` | Numeric fields | Greater than or equal to zero |
| `uuid` | ID fields | Valid UUID v4 format |

## Custom Validators

```go
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func slugValidator(fl validator.FieldLevel) bool {
    return slugPattern.MatchString(fl.Field().String())
}

func dateValidator(fl validator.FieldLevel) bool {
    _, err := time.Parse("2006-01-02", fl.Field().String())
    return err == nil
}

func phoneIDValidator(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^62\d{8,15}$`, phone)
    return matched
}
```

## Request DTO Example

```go
type RegisterRequest struct {
    Name     string `json:"name"      validate:"required,max=150"`
    Email    string `json:"email"     validate:"required,email,max=255"`
    Password string `json:"password"  validate:"required,min=8,max=72"`
    Phone    string `json:"phone"     validate:"omitempty,phone_id"`
}

type CreateEventRequest struct {
    Title       string `json:"title"        validate:"required,max=255"`
    Slug        string `json:"slug"         validate:"required,slug,max=100"`
    GroomName   string `json:"groom_name"   validate:"required,max=255"`
    BrideName   string `json:"bride_name"   validate:"required,max=255"`
    WeddingDate string `json:"wedding_date" validate:"required,date"`
}
```

## Validation Error Formatting

Errors are formatted into a map of field → message for the API response:

```go
func FormatValidationError(err error) map[string]string {
    details := make(map[string]string)
    var verr validator.ValidationErrors
    if errors.As(err, &verr) {
        for _, fe := range verr {
            details[fe.Field()] = validationMessage(fe)
        }
    }
    return details
}

func validationMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Must be a valid email address"
    case "min":
        return fmt.Sprintf("Must be at least %s characters", fe.Param())
    case "max":
        return fmt.Sprintf("Must be at most %s characters", fe.Param())
    case "slug":
        return "Must be a valid slug (lowercase letters, numbers, hyphens)"
    case "date":
        return "Must be a valid date in YYYY-MM-DD format"
    case "oneof":
        return fmt.Sprintf("Must be one of: %s", fe.Param())
    case "phone_id":
        return "Must be a valid Indonesian phone number starting with 62"
    default:
        return "Invalid value"
    }
}
```

## Usage in Handlers

```go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req dto.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
        return
    }
    if err := h.validator.Struct(req); err != nil {
        response.ValidationError(w, err) // calls FormatValidationError internally
        return
    }
    // ... call service
}
```

The response body for validation errors follows the [API convention](https://github.com/yourorg/app-owndangan/blob/main/docs/api/conventions.md#validation-error):

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "Must be a valid email address",
      "password": "Must be at least 8 characters"
    }
  }
}
```
