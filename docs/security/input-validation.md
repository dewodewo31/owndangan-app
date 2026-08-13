# Input Validation

## Two-Layer Validation Architecture

### Layer 1: Transport Validation (Handler)
Validates the request format before any business logic runs. Catches malformed input early.

```go
type CreateEventRequest struct {
    Title       string `json:"title"       validate:"required,min=3,max=200"`
    Date        string `json:"date"        validate:"required,datetime=2006-01-02"`
    Location    string `json:"location"    validate:"required,max=500"`
    GuestCount  int    `json:"guest_count" validate:"omitempty,min=1,max=10000"`
    Description string `json:"description" validate:"omitempty,max=5000"`
}
```

- Use `github.com/go-playground/validator/v10` for struct tags.
- Return 400 with field-level error messages (no stack traces).
- Validate on every mutation endpoint (POST, PUT, PATCH, DELETE).

### Layer 2: Business Validation (Service)
Validates semantics and business rules.

```go
func (s *EventService) Create(ctx context.Context, req CreateEventRequest, userID string) (*Event, error) {
    // Business validation
    if req.Date < time.Now().Format("2006-01-02") {
        return nil, ErrEventDateInPast
    }
    if s.eventRepo.CountByUser(ctx, userID) >= s.planRepo.GetMaxEvents(userID) {
        return nil, ErrPlanLimitExceeded
    }
    // ... create event
}
```

## SQL Injection Prevention

### Rule: Never Use String Concatenation for SQL

```go
// DANGEROUS — never do this
db.Raw("SELECT * FROM events WHERE title = '" + input + "'")

// SAFE — GORM parameterized query
db.Where("title = ?", input).Find(&events)

// SAFE — Named parameters
db.Where("title = @title", sql.Named("title", input)).Find(&events)
```

### Rule: Avoid Raw SQL Unless Necessary
- If you must use `db.Raw()`, always use `?` placeholders.
- GORM's `Find`, `Where`, `Create`, `Update` methods are parameterized by default.
- `db.Exec()` with `?` placeholders is safe.
- `db.Exec()` with string interpolation is **not** safe.

### Rule: Validate Sort/Order Parameters
If user input controls column names in ORDER BY, validate against a whitelist:

```go
var allowedSorts = map[string]bool{
    "created_at": true,
    "title":      true,
    "date":       true,
}

func validateSortField(field string) bool {
    return allowedSorts[field]
}
```

## XSS Prevention

### Go Templates
- Go's `html/template` package auto-escapes HTML. Use it, not `text/template`.
- Never use `template.HTML` with user-supplied data.
- If rich text is needed, use a sanitizer like `github.com/microcosm-cc/bluemonday`.

### API Responses
- Set `Content-Type: application/json` on all API responses (browsers won't interpret as HTML).
- If returning HTML fragments, use `html/template` or sanitize thoroughly.

### Frontend (Next.js)
- JSX auto-escapes by default. Avoid `dangerouslySetInnerHTML`.
- If rendering user content (e.g., guest names on a wedding website), use `textContent` not `innerHTML`.
- For rich text, use a sanitized Markdown renderer.

## File Upload Validation

### Allowed File Types

```go
var allowedMimeTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/webp": true,
    "image/gif":  true,
    "application/pdf": true,
}
```

### Validation Steps (in order)
1. **Check file size** before reading content (max 5MB for images, 10MB for PDFs).
2. **Check MIME type** from `Content-Type` header AND from file content (magic bytes).
3. **Check extension** against whitelist: `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`, `.pdf`.
4. **Re-encode images** (if applicable) to strip EXIF data and potential payloads.

### Implementation

```go
func validateUpload(fileHeader *multipart.FileHeader) error {
    // 1. Size check
    if fileHeader.Size > 5*1024*1024 {
        return ErrFileTooLarge
    }

    // 2. Extension check
    ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
    if !allowedExtensions[ext] {
        return ErrInvalidExtension
    }

    // 3. MIME type check (from header)
    contentType := fileHeader.Header.Get("Content-Type")
    if !allowedMimeTypes[contentType] {
        return ErrInvalidMimeType
    }

    // 4. Magic bytes check
    file, _ := fileHeader.Open()
    defer file.Close()
    buf := make([]byte, 512)
    file.Read(buf)
    detected := http.DetectContentType(buf)
    if detected != contentType {
        return ErrMimeMismatch
    }

    return nil
}
```

## Validation Rules Summary

| Field Type | Rules | Error Response |
|---|---|---|
| Email | `required,email,max=255` | Invalid email format |
| Phone | `required,e164` (international format) | Invalid phone number |
| Name | `required,min=2,max=100` | Name must be 2-100 characters |
| URL | `omitempty,url` | Invalid URL format |
| UUID | `omitempty,uuid4` | Invalid UUID |
| Numeric ID | `required,numeric,min=1` | Invalid ID |
| Date | `required,datetime=2006-01-02` | Invalid date format (YYYY-MM-DD) |
| File | size ≤ 5MB, allowed MIME types | File too large or invalid type |
| Pagination | `min=1,max=100` (page, per_page) | Invalid pagination params |

## Reject Early, Reject Clearly

- All validation errors return HTTP 400.
- Field-level error messages (no generic "bad request").
- Never reveal internal structure (no SQL errors, no stack traces, no file paths).
- Log the full validation error server-side for debugging.