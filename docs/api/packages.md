# Packages

## Public Endpoints

### GET /packages

List active publicly available packages. **Auth:** None | **Rate limit:** 30 req/min per IP

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "770e8400-...", "name": "Free", "code": "free", "price": 0, "duration_days": 7, "guest_limit": 50, "template_group": "standard", "features": { "music.upload": false, "music.preset": true, "video.youtube": false, "custom_domain": false, "watermark.removed": false, "whatsapp.bulk": false, "guestbook.qr": false, "rsvp.export": false, "gallery.photos": 10, "gallery.video": false } },
    { "id": "770e8400-...", "name": "Basic", "code": "basic", "price": 50000, "duration_days": 30, "guest_limit": 200, "template_group": "standard", "features": { "music.upload": true, "music.preset": true, "video.youtube": false, "custom_domain": false, "watermark.removed": false, "whatsapp.bulk": false, "guestbook.qr": false, "rsvp.export": true, "gallery.photos": 20, "gallery.video": false } },
    { "id": "770e8400-...", "name": "Premium", "code": "premium", "price": 150000, "duration_days": 90, "guest_limit": 500, "template_group": "premium", "features": { "music.upload": true, "music.preset": true, "video.youtube": true, "custom_domain": false, "watermark.removed": true, "whatsapp.bulk": true, "guestbook.qr": true, "rsvp.export": true, "gallery.photos": 50, "gallery.video": true } },
    { "id": "770e8400-...", "name": "Pro", "code": "pro", "price": 350000, "duration_days": null, "guest_limit": 1000, "template_group": "all", "features": { "music.upload": true, "music.preset": true, "video.youtube": true, "custom_domain": true, "watermark.removed": true, "whatsapp.bulk": true, "guestbook.qr": true, "rsvp.export": true, "gallery.photos": 100, "gallery.video": true } }
  ],
  "meta": { "request_id": "req-abc123" }
}
```
**Business rules:** Only `is_active = true`. Price in IDR. `duration_days = null` = lifetime (Pro). `template_group` limits accessible templates: standard/premium/all.

---

## Admin Endpoints

### GET /admin/packages

List all packages including inactive. **Auth:** Admin required. Response same structure plus `is_active`, `created_at`, `updated_at`.

---

### POST /admin/packages

Create a package. **Auth:** Admin required

**Request:**
```json
{ "name": "Ultimate", "code": "ultimate", "price": 500000, "duration_days": 365, "guest_limit": 2000, "template_group": "all", "features": { "music.upload": true, "video.youtube": true, "custom_domain": true }, "is_active": true }
```
**Response (201):** Full package object. **Errors:** VALIDATION_ERROR 422, CONFLICT 409 (duplicate name/code).

---

### PUT /admin/packages/:id

Update a package. **Auth:** Admin required

**Request:** `{ "price": 450000, "features": { "whatsapp.bulk": false } }`

**Response (200):** Updated package. **Errors:** VALIDATION_ERROR 422, NOT_FOUND 404, CONFLICT 409.

**Business rules:** Price changes do NOT affect existing active subscriptions. Deactivating a package (`is_active = false`) blocks new purchases but existing subscriptions continue. Features merged at JSONB level (deep merge of provided keys).
