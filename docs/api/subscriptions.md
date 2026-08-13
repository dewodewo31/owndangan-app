# Subscriptions

## Endpoints

### GET /users/me/subscription

Get the authenticated user's current subscription.

Also documented in [users.md](./users.md). This is the authoritative reference for the subscription-specific fields.

**Auth:** Required

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
        "whatsapp.bulk": true,
        "gallery.photos": 50
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
| NOT_FOUND | 404 | User has no subscription record |

**Business rules:**
- Returns the user's active subscription with the latest `start_at`.
- If all subscriptions have expired, returns the most recent one with `status: "expired"`.
- Free users: system auto-creates a Free subscription with 7-day duration upon registration.

---

### POST /subscriptions/upgrade

Initiate a subscription upgrade. Creates a pending transaction and returns a Midtrans Snap token for payment.

**Auth:** Required

**Request:**
```json
{
  "package_id": "770e8400-e29b-41d4-a716-446655440004",
  "redirect_url": "https://app.example.com/payment/finish"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "transaction_id": "660e8400-e29b-41d4-a716-446655440001",
    "order_id": "INV-20250215-XXXX-ABC123",
    "snap_token": "3c5c8f52-9c0e-4a1a-8b5a-1b2c3d4e5f6a",
    "snap_redirect_url": "https://app.sandbox.midtrans.com/snap/v3/redirection/3c5c8f52-...",
    "gross_amount": 150000,
    "expires_at": "2025-02-16T08:30:00Z"
  },
  "meta": { "request_id": "req-abc123" }
}
```

**Error cases:**
| Code | Status | Condition |
|------|--------|-----------|
| VALIDATION_ERROR | 422 | Missing package_id or redirect_url |
| NOT_FOUND | 404 | Package does not exist or is inactive |
| PAYMENT_REQUIRED | 402 | (Not used here — payment is the purpose) |
| LIMIT_EXCEEDED | 422 | User already has an active subscription equal to or higher than the selected package |

**Business rules:**
- The user must NOT already have an active subscription for the same or higher-tier package. If they do, return LIMIT_EXCEEDED.
- If the user has an active lower-tier subscription, the upgrade creates a new subscription that starts after the current one expires (proration: TBD in future).
- `order_id` format: `INV-YYYYMMDD-XXXX-{random8}` where XXXX is a short hash of user_id.
- The `snap_token` is obtained from Midtrans Snap API and is valid for 1 hour.
- After successful payment (via webhook), the new subscription is activated. The old subscription (if any) is marked as `cancelled`.
- Free subscriptions are overridden immediately on payment settlement (no need to wait for expiry).
- The `redirect_url` is passed to Midtrans so the user is redirected there after payment.

**Subscription status lifecycle:**
```
pending (payment initiated)
    │
    ├── active (payment settled) ──→ expired (duration passed)
    │                                  │
    │                                  └── cancelled (user upgraded)
    │
    └── cancelled (user cancelled before payment)

(When upgrading: current active → cancelled; new subscription → active)
```
