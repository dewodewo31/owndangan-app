# Users

Auth required on all endpoints in this document. Standard: `Authorization: Bearer <jwt_token>`.

## Endpoints

### GET /users/me

Get the authenticated user's profile.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Andi Pratama",
    "email": "andi@example.com",
    "phone": "6281234567890",
    "avatar_url": null,
    "role": "user",
    "status": "active",
    "created_at": "2025-01-15T08:30:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| UNAUTHORIZED | 401 | Missing or invalid token |

---

### PUT /users/me

Update the authenticated user's profile.

**Request:**
```json
{
  "name": "Andi Pratama Putra",
  "phone": "6289876543210",
  "avatar_url": "https://storage.example.com/avatars/user-123.jpg"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Andi Pratama Putra",
    "email": "andi@example.com",
    "phone": "6289876543210",
    "avatar_url": "https://storage.example.com/avatars/user-123.jpg",
    "role": "user",
    "status": "active",
    "updated_at": "2025-02-01T10:00:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| VALIDATION_ERROR | 422 | Invalid name (empty, too long), invalid phone format |
| UNAUTHORIZED | 401 | Missing or invalid token |

**Business rules:**
- All fields are optional. Only provided fields are updated (partial/PATCH semantics via PUT).
- Name: max 255 characters, must not be empty if provided.
- Phone: must be valid Indonesian format (62xxx, 10-15 digits) if provided.
- Email cannot be changed through this endpoint (use a separate email change flow).
- avatar_url, if provided, should point to a URL already uploaded to object storage.

---

### GET /users/me/subscription

Get the authenticated user's current active subscription.

See [subscriptions.md](./subscriptions.md) for full documentation.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "package": {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "Premium",
      "code": "premium",
      "price": 150000,
      "guest_limit": 500,
      "template_group": "premium",
      "features": {
        "music.upload": true,
        "video.youtube": true,
        "custom_domain": false,
        "watermark.removed": true,
        "whatsapp.bulk": false
      }
    },
    "status": "active",
    "start_at": "2025-01-20T00:00:00Z",
    "expires_at": "2025-04-20T00:00:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| UNAUTHORIZED | 401 | Missing or invalid token |
| NOT_FOUND | 404 | User has no subscription record yet |

**Business rules:**
- Returns the active subscription. If there are multiple active subscriptions (edge case), returns the one with the latest start_at.
- If the subscription has expired, status will be "expired" even if within the grace period.
- Free tier users see their default Free subscription with 7-day duration.
- If no subscription record exists (very new user), return NOT_FOUND. Frontend should handle this by treating the user as Free tier.
