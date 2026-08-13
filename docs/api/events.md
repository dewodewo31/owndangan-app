# Events

The core entity — a wedding invitation event. All endpoints require JWT auth.

## Endpoints

### GET /events

List user's events. **Query:** `?page=1&per_page=20&status=published&sort=created_at&order=desc`

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "880e8400-...", "title": "Andi & Sinta", "slug": "andi-sinta-2025", "couple_name": "Andi & Sinta", "status": "published", "wedding_date": "2025-06-15", "wedding_time": "09:00:00", "view_count": 142, "guest_count": 48, "rsvp_count": 35, "created_at": "2025-01-20T10:00:00Z" }
  ],
  "meta": { "pagination": { "page": 1, "per_page": 20, "total": 5, "total_pages": 1 }, "request_id": "req-abc123" }
}
```
**Errors:** UNAUTHORIZED 401, PAYMENT_REQUIRED 402.

**Rules:** Only non-deleted events owned by user. `guest_count` = total guests; `rsvp_count` = guests with RSVP.

---

### POST /events

Create an event.

**Request:**
```json
{
  "title": "Andi & Sinta", "couple_name": "Andi & Sinta", "groom_name": "Andi Pratama", "bride_name": "Sinta Dewi",
  "groom_parents": "Mr. & Mrs. Budi", "bride_parents": "Mr. & Mrs. Agus",
  "wedding_date": "2025-06-15", "wedding_time": "09:00:00",
  "ceremony_venue": "Gedung Serbaguna Jakarta", "ceremony_address": "Jl. Merdeka No. 1, Jakarta", "ceremony_map_url": "https://maps.google.com/?q=...",
  "reception_venue": "Hotel Indonesia", "reception_address": "Jl. Thamrin No. 2, Jakarta", "reception_map_url": "https://maps.google.com/?q=..."
}
```
**Response (201):**
```json
{ "success": true, "data": { "id": "880e8400-...", "user_id": "550e8400-...", "title": "Andi & Sinta", "slug": "andi-sinta-2025", "couple_name": "Andi & Sinta", "status": "draft", "wedding_date": "2025-06-15", "created_at": "2025-01-20T10:00:00Z" }, "meta": { "request_id": "req-abc123" } }
```
**Errors:** VALIDATION_ERROR 422, LIMIT_EXCEEDED 422 (max events per tier), UNAUTHORIZED 401.

**Rules:** Slug auto-generated. Initial status = `draft`. Event limits: Free=1, Basic=3, Premium=10, Pro=unlimited. `event_sections` + `digital_gifts` auto-created.

---

### GET /events/:id

Full event details. **Errors:** UNAUTHORIZED 401, NOT_FOUND 404, FORBIDDEN 403 (wrong owner).

---

### PUT /events/:id

Partial update. All fields optional. **Errors:** 422, 401, 404. Slug NOT changed by PUT.

---

### DELETE /events/:id

Soft delete. Sets `deleted_at`. If published, unpublished first. **Errors:** 401, 404. **Response:** `{ "success": true, "data": null, "meta": { "request_id": "req-abc123" } }`

---

### POST /events/:id/publish

Publish invitation. **Requires:** active subscription + minimal data (title, couple_name, wedding_date).

**Response (200):**
```json
{ "success": true, "data": { "id": "880e8400-...", "status": "published", "published_at": "2025-02-01T12:00:00Z", "public_url": "https://undangan.example.com/andi-sinta-2025" }, "meta": { "request_id": "req-abc123" } }
```
**Errors:** UNAUTHORIZED 401, NOT_FOUND 404, PAYMENT_REQUIRED 402 (no subscription), CONFLICT 409 (already published). Sets `status='published'`, `published_at=now()`.

---

### POST /events/:id/unpublish

Remove public access. Only published events. Sets `status='unpublished'`, clears `published_at`. **Errors:** 401, 404, CONFLICT 409 (not published).

---

### PUT /events/:id (template & video)

`PUT /events/:id` additionally accepts `template_id` (assign a template by id, gated by plan) and `video_url` (string). Assigning a `premium`/`all` group template requires the matching entitlement or returns FORBIDDEN 403.

---

### GET /templates

List templates available to the caller's plan. Free → `standard`; Premium → `standard` + `premium`; Pro → `standard` + `premium` + `all`.

**Response (200):**
```json
{ "success": true, "data": [
  { "id": "uuid", "name": "Classic", "group_name": "standard",
    "thumbnail_url": "https://.../classic.png",
    "css_config": {"primary_color":"#b22234"}, "layout_config": {"hero":"cover"} }
], "meta": { "request_id": "req-abc" } }
```

---

### PUT /events/:id/template

Assign a template to the event. **Request:** `{ "template_id": "uuid" }`. **Errors:** 401, 404 (template not found / event), 403 (plan does not permit template group).

---

### GET /events/:id/sections

Return the event's section configuration. **Response (200):**
```json
{ "success": true, "data": {
  "id": "uuid", "event_id": "uuid",
  "hero_enabled": true, "couple_enabled": true, "event_details_enabled": true,
  "gallery_enabled": true, "video_enabled": false, "music_id": null,
  "rsvp_enabled": true, "guestbook_enabled": true, "digital_gifts_enabled": false,
  "dress_code": "", "closing_message": "", "opening_message": ""
}}
```

---

### PUT /events/:id/sections

Upsert section configuration. All fields optional (bools toggle a section). `music_id` links an uploaded/preset track. Entitlement-gated: `gallery_enabled`, `video_enabled`, and `music_id` require the corresponding plan feature (LIMIT_EXCEEDED 422 otherwise). **Errors:** 401, 404, 403.

---

### GET /events/:id/digital-gifts

Return the event's digital gift info. **Response (200):**
```json
{ "success": true, "data": {
  "id": "uuid", "event_id": "uuid",
  "bank_accounts": [{"bank":"Mandiri","account":"1234","name":"Budi"}],
  "ewallet": {"ovo":"0812"}, "qris_image_url": "", "gift_message": "Terima kasih"
}}
```

### PUT /events/:id/digital-gifts

Update digital gift fields. All optional. QRIS (`qris_image_url`) requires the `digital_gift.qris` feature (LIMIT_EXCEEDED 422). **Errors:** 401, 404, 403.

---

### Gallery

`GET /events/:id/gallery` — list photos sorted by `sort_order`.
`POST /events/:id/gallery/upload` — multipart `file` (+ optional `caption`). Gated by `gallery.max` (LIMIT_EXCEEDED 422). **Response (201):** `{ "success": true, "data": { "id", "image_url", "caption", "sort_order" } }`.
`DELETE /events/:id/gallery/{photoID}` — deletes a photo.
`PUT /events/:id/gallery/reorder` — body `{ "photos": [{ "id": "uuid", "sort_order": <int> }, ...] }`.

Uploaded files are stored under `STORAGE_LOCAL_PATH` (local) and served statically at `/uploads/*`.

---

### Music

`GET /events/:id/music` — currently configured track. Returns `200` with `"data": null` when no music is configured (e.g. a newly created event), **not** a 404, so the editor can render an empty/default state.
`GET /events/:id/music/presets` — list preset tracks (gated by `music.preset`).
`POST /events/:id/music/upload` — multipart `file` (+ optional `title`), requires `music.upload` entitlement. **Response (201):** `{ "success": true, "data": { "id","event_id","title","file_url",... } }`.
`POST /events/:id/music/presets` — body `{ "preset_id": "uuid" }`, assigns a preset; requires `music.preset` entitlement.
`DELETE /events/:id/music` — detach and delete the currently selected track, disabling background music. Safe no-op when none is set. **Response (200):** `{ "success": true, "data": { "message": "Music removed" } }`.

---

### Public: GET /e/{slug}

Unauthenticated public invitation payload combining event detail, resolved template (css/layout), section flags, ordered gallery, guestbook messages, and digital gift info.
