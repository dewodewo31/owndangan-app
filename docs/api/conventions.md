# API Conventions

## Base URL

```
/api/v1
```

All endpoints are prefixed with `/api/v1`.

## HTTP Methods

| Method | Operation |
|--------|-----------|
| GET | Read/List resources |
| POST | Create resource or action |
| PUT | Full update |
| PATCH | Partial update |
| DELETE | Delete resource |

## Authentication

- Authenticated endpoints: `Authorization: Bearer <jwt_token>`
- Public endpoints: no auth header required
- Webhook endpoint: verified by signature (not JWT)

## Standard Response Envelope

### Success

```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "request_id": "req-abc123"
  }
}
```

### List with Pagination

```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 100,
      "total_pages": 5
    },
    "request_id": "req-abc123"
  }
}
```

### Error

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable description"
  },
  "meta": {
    "request_id": "req-abc123"
  }
}
```

### Validation Error

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "Email is required",
      "password": "Password must be at least 8 characters"
    }
  },
  "meta": {
    "request_id": "req-abc123"
  }
}
```

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 422 | Request validation failed |
| UNAUTHORIZED | 401 | Missing or invalid token |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| CONFLICT | 409 | Resource conflict (duplicate) |
| RATE_LIMITED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |
| PAYMENT_REQUIRED | 402 | Subscription required |
| LIMIT_EXCEEDED | 422 | Plan limit exceeded |

## Pagination

Query parameters: `?page=1&per_page=20`

- Default: page=1, per_page=20
- Max per_page: 100
- Response includes pagination metadata.

## Filtering

Query parameters: `?field=value`

- Exact match: `?status=active`
- Nested: `?search=keyword` (implementation-specific)

## Sorting

Query parameter: `?sort=field&order=asc|desc`

- Default: `?sort=created_at&order=desc`
- Valid fields vary by endpoint.

## Searching

Query parameter: `?q=keyword`

Search across relevant text fields (name, email, title, etc.). Implementation varies by endpoint.

## Request ID

- Each request receives a unique `X-Request-ID` header.
- If client sends `X-Request-ID`, it's preserved.
- Included in response meta as `request_id`.

## Standard Headers

### Request
```
Content-Type: application/json
Authorization: Bearer <token>  (if authenticated)
X-Request-ID: <uuid>  (optional, client-provided)
```

### Response
```
Content-Type: application/json
X-Request-ID: <uuid>
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1628640000
```

## Rate Limiting

> TODO: Define rate limit strategy.

- Public endpoints: 30 requests/minute per IP
- Authenticated endpoints: 100 requests/minute per user
- Admin endpoints: 200 requests/minute per admin
- Webhook endpoint: no rate limit (IP-whitelisted)

## Idempotency

> TODO: Define idempotency strategy.

- POST endpoints should support `Idempotency-Key` header for payment operations.
- Duplicate webhook notifications are handled by order_id uniqueness.

## Webhook Handling

- Webhook endpoint: `POST /api/v1/webhook/midtrans`
- No JWT auth; verified via Midtrans signature.
- Response: 200 OK (acknowledgment) — even if processing fails.
- Actual processing may be async (future: queue system).

## List of Endpoints by Category

### Public (No Auth)
- `GET /api/v1/packages` — List packages
- `GET /api/v1/public/events/:slug` — Get public invitation
- `POST /api/v1/public/rsvps` — Submit RSVP
- `POST /api/v1/public/guestbook` — Submit guestbook message
- `GET /api/v1/public/guestbook/:slug` — List approved guestbook messages
- `GET /api/v1/health` — Health check

### Authentication
- `POST /api/v1/auth/register` — Register
- `POST /api/v1/auth/login` — Login
- `POST /api/v1/auth/refresh` — Refresh token
- `POST /api/v1/auth/logout` — Logout

### User (JWT required)
- `GET /api/v1/users/me` — Get profile
- `PUT /api/v1/users/me` — Update profile
- `GET /api/v1/users/me/subscription` — Get current subscription
- `GET /api/v1/events` — List user's events
- `POST /api/v1/events` — Create event
- `GET /api/v1/events/:id` — Get event details
- `PUT /api/v1/events/:id` — Update event
- `DELETE /api/v1/events/:id` — Delete event
- `POST /api/v1/events/:id/publish` — Publish event
- `POST /api/v1/events/:id/unpublish` — Unpublish event
- `GET /api/v1/events/:id/guests` — List guests
- `POST /api/v1/events/:id/guests` — Add guest
- `POST /api/v1/events/:id/guests/import` — CSV import
- `PUT /api/v1/events/:id/guests/:gid` — Update guest
- `DELETE /api/v1/events/:id/guests/:gid` — Delete guest
- `GET /api/v1/events/:id/rsvps` — RSVP recap
- `GET /api/v1/events/:id/guestbook` — List guestbook (with moderation)
- `GET /api/v1/events/:id/gallery` — List gallery
- `POST /api/v1/events/:id/gallery` — Upload photo
- `GET /api/v1/events/:id/music` — List music
- `POST /api/v1/events/:id/music` — Upload music
- `GET /api/v1/events/:id/digital-gifts` — Get gift config
- `PUT /api/v1/events/:id/digital-gifts` — Update gift config
- `GET /api/v1/events/:id/sections` — Get section config
- `PUT /api/v1/events/:id/sections` — Update section config
- `POST /api/v1/payments/snap` — Create Snap transaction
- `POST /api/v1/subscriptions/upgrade` — Upgrade subscription

### Admin (JWT + admin role)
- `GET /api/v1/admin/dashboard` — Dashboard stats
- `GET /api/v1/admin/users` — List users
- `GET /api/v1/admin/users/:id` — Get user details
- `PUT /api/v1/admin/users/:id/status` — Suspend/activate user
- `GET /api/v1/admin/transactions` — List transactions
- `GET /api/v1/admin/packages` — List packages
- `POST /api/v1/admin/packages` — Create package
- `PUT /api/v1/admin/packages/:id` — Update package
- `GET /api/v1/admin/templates` — List templates
- `POST /api/v1/admin/templates` — Upload template
- `PUT /api/v1/admin/templates/:id` — Update template
- `DELETE /api/v1/admin/templates/:id` — Delete template

### Webhook
- `POST /api/v1/webhook/midtrans` — Midtrans notification

## Versioning

- URL-based: `/api/v1/...`
- Backward-compatible changes: same version.
- Breaking changes: new version (`/api/v2/`).
- Deprecated endpoints: include `Deprecated` header.