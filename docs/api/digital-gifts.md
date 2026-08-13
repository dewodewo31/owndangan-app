# Digital Gifts

Each event has one gift configuration (1:1). Bank accounts, e-wallet, QRIS.

## Endpoints

### GET /events/:id/digital-gifts

Get gift config (owner). **Auth:** Required (event owner)

**Response (200):**
```json
{
  "success": true, "data": {
    "id": "cc0e8400-...", "event_id": "880e8400-...",
    "bank_accounts": [ { "bank_name": "BCA", "account_number": "1234567890", "account_holder": "Andi Pratama" }, { "bank_name": "Mandiri", "account_number": "9876543210", "account_holder": "Sinta Dewi" } ],
    "ewallet": { "gopay_number": "6281234567890", "gopay_holder": "Andi Pratama" },
    "qris_image_url": "https://storage.example.com/qris/event-123-qris.png",
    "gift_message": "Terima kasih atas doa dan hadiahnya!",
    "created_at": "2025-01-20T10:00:00Z", "updated_at": "2025-02-01T12:00:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Errors:** UNAUTHORIZED 401, NOT_FOUND 404. Record auto-created with event (empty defaults).

---

### PUT /events/:id/digital-gifts

Update gift config. **Auth:** Required (event owner)

**Request:**
```json
{
  "bank_accounts": [ { "bank_name": "BCA", "account_number": "1234567890", "account_holder": "Andi Pratama" } ],
  "ewallet": { "gopay_number": "6281234567890", "gopay_holder": "Andi Pratama" },
  "qris_image_url": "https://storage.example.com/qris/event-123-qris-v2.png",
  "gift_message": "Updated message."
}
```
**Response (200):** Updated config. **Errors:** VALIDATION_ERROR 422 (invalid bank structure), UNAUTHORIZED 401, NOT_FOUND 404.

**Rules:** Partial top-level update. Send `null` to clear a field. `bank_accounts` must be array of `{bank_name, account_number, account_holder}`. `gift_message` max 500 chars. `qris_image_url` should be pre-uploaded URL.

---

### GET /public/events/:slug/digital-gifts

Public gift info for invitation page. **Auth:** None | **Rate limit:** 30 req/min per IP

**Response (200):**
```json
{ "success": true, "data": { "bank_accounts": [ { "bank_name": "BCA", "account_number": "1234567890", "account_holder": "Andi Pratama" } ], "ewallet": { "gopay_number": "6281234567890", "gopay_holder": "Andi Pratama" }, "qris_image_url": "https://storage.example.com/qris/event-123-qris-v2.png", "gift_message": "Terima kasih atas doa dan hadiahnya!" }, "meta": { "request_id": "req-abc123" } }
```
**Errors:** NOT_FOUND 404 (slug not found, not published, or digital gifts disabled).

**Rules:** Only returns if event is published AND `event_sections.digital_gifts_enabled = true`. No UUIDs or internal IDs exposed. If digital gifts disabled, return 404 (do not reveal existence).
