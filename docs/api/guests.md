# Guests

Guest management for an event. All endpoints require JWT auth + event ownership.

## Endpoints

### GET /events/:id/guests

List guests with pagination, search, category filter. **Query:** `?page=1&per_page=20&category=family&q=keyword&sort=name&order=asc`

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "990e8400-...", "event_id": "880e8400-...", "name": "Budi Santoso", "phone": "628123456789", "category": "family", "note": "Kakak dari mempelai pria", "token": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3", "rsvp": { "attendance": "yes", "guest_count": 3, "message": "Akan hadir!", "submitted_at": "2025-03-01T14:00:00Z" }, "created_at": "2025-02-10T09:00:00Z" }
  ],
  "meta": { "pagination": { "page": 1, "per_page": 20, "total": 48, "total_pages": 3 }, "request_id": "req-abc123" }
}
```
**Errors:** UNAUTHORIZED 401, NOT_FOUND 404. **Rules:** RSVP embedded if submitted, else null. Token is 40-char hex.

---

### POST /events/:id/guests

Add a guest. **Request:** `{ "name": "Budi Santoso", "phone": "628123456789", "category": "family", "note": "Kakak dari mempelai pria" }`

**Response (201):** Full guest object with generated token.

**Errors:** VALIDATION_ERROR 422, LIMIT_EXCEEDED 422 (subscription limit), UNAUTHORIZED 401, NOT_FOUND 404.

**Rules:** Name required (max 255). Category: family/friend/colleague/other (default: family). Phone optional, must be valid 62xxx. Token = 40-byte random hex. Check subscription guest limit before adding.

---

### PUT /events/:id/guests/:gid

Update guest. **Request:** `{ "name": "Budi Santoso Putra", "category": "friend" }`

**Response (200):** Updated guest. **Errors:** VALIDATION_ERROR 422, UNAUTHORIZED 401, NOT_FOUND 404. Partial update; token never changes.

---

### DELETE /events/:id/guests/:gid

Soft delete guest. **Response:** `{ "success": true, "data": null, "meta": { "request_id": "req-abc123" } }`

**Errors:** UNAUTHORIZED 401, NOT_FOUND 404. RSVP records preserved. Token freed for reuse.

---

### POST /events/:id/guests/import

CSV import. **Accept:** `multipart/form-data`. Field: `file`.

CSV format:
```
name,phone,category,note
Budi Santoso,628123456789,family,Kakak mempelai pria
Sinta Dewi,628987654321,friend,Teman SMP
```

**Response (200):**
```json
{ "success": true, "data": { "imported": 45, "skipped": 3, "errors": [ { "row": 12, "reason": "Duplicate name" }, { "row": 23, "reason": "Invalid phone format" } ] }, "meta": { "request_id": "req-abc123" } }
```
**Errors:** VALIDATION_ERROR 422, LIMIT_EXCEEDED 422, UNAUTHORIZED 401, NOT_FOUND 404.

**Rules:** Required columns: `name`. Optional: `phone`, `category`, `note`. UTF-8 CSV. Max 2 MB. Duplicate names skipped (case-insensitive). Transactional: critical errors roll back all; non-critical collected in errors array.
