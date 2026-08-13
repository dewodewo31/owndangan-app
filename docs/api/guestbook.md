# Guestbook

Guestbook messages for wedding invitations. Public endpoints for submission and reading. Authenticated endpoints for moderation.

## Endpoints

### POST /public/guestbook

Submit a guestbook message. **Auth:** None | **Rate limit:** 10 req/min per IP

**Request:**
```json
{ "slug": "andi-sinta-2025", "name": "Budi Santoso", "message": "Selamat menempuh hidup baru! Semoga langgeng selalu." }
```

**Response (201):**
```json
{
  "success": true, "data": {
    "id": "bb0e8400-...", "name": "Budi Santoso",
    "message": "Selamat menempuh hidup baru! Semoga langgeng selalu.",
    "is_approved": false, "created_at": "2025-03-01T15:00:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Errors:** VALIDATION_ERROR 422 (missing slug/name/message), NOT_FOUND 404 (slug not found or event not published).

**Business rules:**
- `slug` is the public event slug (not UUID).
- Name required (max 255), message required (max 1000 characters).
- New messages = `is_approved: false` (requires moderation).
- Event must be published; otherwise return 404.
- Duplicate prevention: same name + event limited to 1 submission per hour.

---

### GET /public/guestbook/:slug

List approved guestbook messages. **Auth:** None | **Rate limit:** 30 req/min per IP

**Query:** `?page=1&per_page=20`

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "bb0e8400-...", "name": "Budi Santoso", "message": "Selamat menempuh hidup baru!", "created_at": "2025-03-01T15:00:00Z" }
  ],
  "meta": { "pagination": { "page": 1, "per_page": 20, "total": 15, "total_pages": 1 }, "request_id": "req-abc123" }
}
```
**Errors:** NOT_FOUND 404 (slug not found or not published).

**Business rules:**
- Only returns messages where `is_approved = true`.
- Sorted by `created_at` descending (newest first).
- `is_approved` field is never exposed to the public.

---

### GET /events/:id/guestbook

List messages with moderation status (event owner view). **Auth:** Required

**Query:** `?page=1&per_page=20&is_approved=true&sort=created_at&order=desc`

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "bb0e8400-...", "name": "Budi Santoso", "message": "Selamat menempuh hidup baru!", "is_approved": false, "created_at": "2025-03-01T15:00:00Z" },
    { "id": "bb0e8400-...", "name": "Sinta Dewi", "message": "Happy wedding!", "is_approved": true, "created_at": "2025-03-01T14:00:00Z" }
  ],
  "meta": { "pagination": { "page": 1, "per_page": 20, "total": 15, "total_pages": 1 }, "request_id": "req-abc123" }
}
```
**Errors:** UNAUTHORIZED 401, NOT_FOUND 404.

**Business rules:**
- Event owner sees all messages including unapproved ones.
- `is_approved` is exposed here (unlike public endpoint).
- Filter by `is_approved=true/false` to see only approved or pending messages.
- Moderation actions (approve/reject) use a future `PUT /events/:id/guestbook/:gid` endpoint.
- Delete via future `DELETE /events/:id/guestbook/:gid`.
