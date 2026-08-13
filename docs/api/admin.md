# Admin

All endpoints require JWT + `admin` role. Unauthorized non-admin gets:
```json
{ "success": false, "error": { "code": "FORBIDDEN", "message": "Admin access required" }, "meta": { "request_id": "req-abc123" } }
```

## Endpoints

### GET /admin/dashboard

Platform-wide statistics.

**Response (200):**
```json
{
  "success": true, "data": {
    "stats": { "total_users": 1240, "active_users": 980, "suspended_users": 12, "total_events": 2100, "published_events": 1850, "active_subscriptions": 520, "total_revenue": 78500000, "revenue_this_month": 12500000 },
    "recent_users": [ { "id": "550e8400-...", "name": "Sinta Dewi", "email": "sinta@example.com", "status": "active", "created_at": "2025-02-10T14:00:00Z" } ],
    "recent_transactions": [ { "id": "660e8400-...", "user_name": "Sinta Dewi", "package_name": "Premium", "gross_amount": 150000, "status": "settlement", "created_at": "2025-02-10T14:05:00Z" } ]
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Error cases:** UNAUTHORIZED 401, FORBIDDEN 403.

**Business rules:** total_revenue = lifetime settled gross. revenue_this_month = current calendar month. recent_users = last 5. recent_transactions = last 5 settled.

---

### GET /admin/users

List all users with filtering/pagination.

**Query:** `?page=1&per_page=20&status=active&role=user&q=keyword&sort=created_at&order=desc`

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "550e8400-...", "name": "Andi Pratama", "email": "andi@example.com", "phone": "6281234567890", "role": "user", "status": "active", "subscription": { "package_name": "Premium", "status": "active", "expires_at": "2025-04-20T00:00:00Z" }, "events_count": 3, "created_at": "2025-01-15T08:30:00Z" }
  ],
  "meta": { "pagination": { "page": 1, "per_page": 20, "total": 1240, "total_pages": 62 }, "request_id": "req-abc123" }
}
```
**Error cases:** UNAUTHORIZED 401, FORBIDDEN 403.

**Business rules:** `q` matches partial name/email case-insensitively. events_count excludes soft-deleted events.

---

### PUT /admin/users/:id/status

Suspend or activate a user. **Auth:** Admin only

**Request:** `{ "status": "suspended" }`

**Response (200):**
```json
{ "success": true, "data": { "id": "550e8400-...", "name": "Andi Pratama", "email": "andi@example.com", "status": "suspended", "updated_at": "2025-02-11T09:00:00Z" }, "meta": { "request_id": "req-abc123" } }
```
**Error cases:** VALIDATION_ERROR 422 (invalid status value — must be 'active' or 'suspended'), NOT_FOUND 404, UNAUTHORIZED 401, FORBIDDEN 403.

**Business rules:** Valid statuses: 'active', 'suspended'. Cannot suspend yourself. Cannot suspend other admins. Suspended users cannot log in or manage events; published events remain publicly accessible. Reactivating restores all access.
