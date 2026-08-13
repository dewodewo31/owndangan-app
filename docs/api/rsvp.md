# RSVP

RSVP (attendance response) endpoints. Public submission is by guest token; authenticated recap is by event ownership.

## Endpoints

### POST /public/rsvps

Submit an RSVP response. No authentication required. The guest is identified by their unique token.

**Auth:** None | **Rate limit:** 10 req/min per IP

**Request:**
```json
{
  "token": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3",
  "attendance": "yes",
  "guest_count": 3,
  "message": "Selamat menempuh hidup baru! Kami sangat antusias."
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "aa0e8400-e29b-41d4-a716-446655440008",
    "guest_id": "990e8400-e29b-41d4-a716-446655440007",
    "attendance": "yes",
    "guest_count": 3,
    "message": "Selamat menempuh hidup baru! Kami sangat antusias.",
    "submitted_at": "2025-03-01T14:00:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| VALIDATION_ERROR | 422 | Missing token, invalid attendance value |
| NOT_FOUND | 404 | Token does not match any guest |
| CONFLICT | 409 | Guest has already submitted an RSVP |

**Business rules:**
- `attendance` must be one of: `'yes'`, `'no'`, `'maybe'`.
- `guest_count` is the total number of people attending (including the guest). Must be at least 1.
- `message` is optional (max 500 characters).
- The RSVP is linked to both the `guest_id` (from token lookup) and `event_id`.
- **Duplicate handling:** Each guest can only submit one RSVP. If the guest has already submitted, return CONFLICT (409). The frontend should fetch the existing RSVP status via the guest token (future endpoint) or allow the guest to edit via a separate `PUT /public/rsvps` endpoint (future).
- The token is looked up across all events. The guest does not need to specify which event.
- RSVP is immutable once submitted — no update endpoint exists yet. A future `PUT /public/rsvps` may allow updating attendance.
- `submitted_at` is set server-side to the current timestamp.

---

### GET /events/:id/rsvps

Get RSVP recap for an event. Shows aggregated statistics and individual responses.

**Auth:** Required (event owner)

**Query parameters:**
| Param | Type | Description |
|-------|------|-------------|
| page | int | Default 1 |
| per_page | int | Default 20, max 100 |
| attendance | string | Filter: 'yes', 'no', 'maybe' |
| sort | string | name, attendance, submitted_at |
| order | string | asc, desc |

**Response (200):**
```json
{
  "success": true,
  "data": {
    "summary": {
      "total_invited": 48,
      "responded": 35,
      "pending": 13,
      "attending_yes": 25,
      "attending_no": 5,
      "attending_maybe": 5,
      "total_guests_yes": 82,
      "total_guests_no": 5,
      "total_guests_maybe": 12
    },
    "responses": [
      {
        "guest_id": "990e8400-e29b-41d4-a716-446655440007",
        "guest_name": "Budi Santoso",
        "guest_category": "family",
        "attendance": "yes",
        "guest_count": 3,
        "message": "Akan hadir!",
        "submitted_at": "2025-03-01T14:00:00Z"
      }
    ]
  },
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 35,
      "total_pages": 2
    },
    "request_id": "req-abc123"
  }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| UNAUTHORIZED | 401 | Missing or invalid token |
| NOT_FOUND | 404 | Event not found or not owned by user |

**Business rules:**
- `summary` provides aggregate counts:
  - `total_invited`: total guests for the event (including soft-deleted? no, only active).
  - `responded`: number of guests who submitted RSVP.
  - `pending`: `total_invited - responded`.
  - `attending_yes/no/maybe`: count of responses by attendance type.
  - `total_guests_yes`: sum of `guest_count` for all 'yes' responses (total people attending).
  - `total_guests_no` and `total_guests_maybe`: sum of `guest_count` for each category (for 'no', guest_count is typically 1).
- Pagination applies to the `responses` array only, not the summary.
- Export feature (CSV) for RSVP data is planned but not yet implemented.
