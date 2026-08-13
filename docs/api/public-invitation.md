# Public Invitation

Returns full invitation data for rendering the public invitation page. No authentication required.

## Endpoints

### GET /public/events/:slug

Get full public invitation data by slug. This is the main endpoint for rendering the invitation website.

**Auth:** None | **Rate limit:** 60 req/min per IP

**Response (200):**
```json
{
  "success": true,
  "data": {
    "event": {
      "title": "Andi & Sinta",
      "couple_name": "Andi & Sinta",
      "groom_name": "Andi Pratama",
      "bride_name": "Sinta Dewi",
      "groom_parents": "Mr. & Mrs. Budi",
      "bride_parents": "Mr. & Mrs. Agus",
      "wedding_date": "2025-06-15",
      "wedding_time": "09:00:00",
      "ceremony_venue": "Gedung Serbaguna Jakarta",
      "ceremony_address": "Jl. Merdeka No. 1, Jakarta",
      "ceremony_map_url": "https://maps.google.com/?q=...",
      "reception_venue": "Hotel Indonesia",
      "reception_address": "Jl. Thamrin No. 2, Jakarta",
      "reception_map_url": "https://maps.google.com/?q=...",
      "view_count": 142
    },
    "template": {
      "name": "Elegant",
      "group_name": "premium",
      "css_config": {
        "primary_color": "#d4af37",
        "secondary_color": "#ffffff",
        "font_family": "Playfair Display",
        "background_pattern": "floral"
      },
      "layout_config": {
        "hero_style": "fullscreen",
        "gallery_layout": "masonry",
        "rsvp_style": "inline"
      }
    },
    "sections": {
      "hero_enabled": true,
      "couple_enabled": true,
      "event_details_enabled": true,
      "gallery_enabled": true,
      "video_enabled": false,
      "rsvp_enabled": true,
      "guestbook_enabled": true,
      "digital_gifts_enabled": true,
      "music": {
        "title": "Perfect - Ed Sheeran",
        "file_url": null,
        "preset": "perfect-ed-sheeran",
        "is_preset": true
      }
    },
    "gallery": [
      {
        "image_url": "https://storage.example.com/gallery/event-123/photo-1.jpg",
        "caption": "Prewedding",
        "sort_order": 0
      },
      {
        "image_url": "https://storage.example.com/gallery/event-123/photo-2.jpg",
        "caption": "Engagement",
        "sort_order": 1
      }
    ],
    "guestbook": [
      {
        "name": "Budi Santoso",
        "message": "Selamat! Semoga langgeng.",
        "created_at": "2025-03-01T15:00:00Z"
      }
    ],
    "digital_gifts": {
      "enabled": true,
      "data": {
        "bank_accounts": [
          {
            "bank_name": "BCA",
            "account_number": "1234567890",
            "account_holder": "Andi Pratama"
          }
        ],
        "ewallet": {
          "gopay_number": "6281234567890",
          "gopay_holder": "Andi Pratama"
        },
        "qris_image_url": "https://storage.example.com/qris/event-123-qris.png",
        "gift_message": "Terima kasih atas doa dan hadiahnya!"
      }
    }
  },
  "meta": { "request_id": "req-abc123" }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| NOT_FOUND | 404 | Slug does not exist OR event is not published OR event is expired |
| GONE | 410 | Event was published but wedding date has passed (configurable: return 404 or 410) |

**Business rules:**

**Visibility rules (most important):**
- Only events with `status = 'published'` are accessible. Draft or unpublished events return 404.
- The slug is the unique public identifier. The event UUID is never exposed.
- Internal fields like `user_id`, `id`, `deleted_at`, `status` are never included in the public response.

**Expiry check:**
- If the wedding date has passed (more than 1 day ago) AND the subscription has expired, the event may be considered expired. For now, the event remains accessible indefinitely after publishing (future: expiry based on subscription).

**Section-gated data:**
- If `sections.guestbook_enabled` is `false`, the `guestbook` array should be empty or omitted.
- If `sections.digital_gifts_enabled` is `false`, `digital_gifts.data` should be `null` (but `digital_gifts.enabled` is `false`).
- If `sections.gallery_enabled` is `false`, the `gallery` array should be empty.
- The `music` block is included only if `sections.music_id` is set.

**Guestbook data:**
- Only returns approved messages (`is_approved = true`).

**View counting:**
- Each request to this endpoint increments `view_count` on the event by 1.
- Rate-limited view counting: the same IP within 5 minutes does not count as a duplicate view (simple memory-based dedup).

**Frontend usage:**
- The frontend uses this single response to render the entire invitation page.
- The `sections` object determines which UI components to show/hide.
- The `template.css_config` values are applied as CSS custom properties on the root element.
- The `template.name` is used to load the correct template component.
- All URLs (images, maps) are absolute and directly usable.
