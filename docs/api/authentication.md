# Authentication

## Token Format

- **Access token:** JWT (HS256), 15-minute expiry
- **Refresh token:** Opaque 64-byte hex string, 7-day expiry
- **Access token payload:** `{ "sub": "user-uuid", "role": "user", "iat": 1628000000, "exp": 1628000900 }`

## Endpoints

### POST /auth/register

Create a new user account. **Auth:** None | **Rate limit:** 5 req/min per IP

**Request:**
```json
{ "name": "Andi Pratama", "email": "andi@example.com", "password": "secret1234", "phone": "6281234567890" }
```
**Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Andi Pratama",
    "email": "andi@example.com",
    "phone": "6281234567890",
    "role": "user",
    "status": "active",
    "created_at": "2025-01-15T08:30:00Z",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "a1b2c3d4e5f6...64-byte-hex",
    "expires_in": 900
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Error cases:** VALIDATION_ERROR 422 (missing/invalid fields, password < 8 chars), CONFLICT 409 (duplicate email).

**Business rules:** Password min 8 chars, bcrypt hashed. Email uniqueness enforced case-insensitively. Phone optional, must be valid 62xxx format. New users auto-get Free 7-day subscription.

---

### POST /auth/login

Authenticate and receive tokens. **Auth:** None | **Rate limit:** 10 req/min per IP

**Request:** `{ "email": "andi@example.com", "password": "secret1234" }`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Andi Pratama",
    "email": "andi@example.com",
    "role": "user",
    "status": "active",
    "created_at": "2025-01-15T08:30:00Z",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "a1b2c3d4e5f6...64-byte-hex",
    "expires_in": 900
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Error cases:** VALIDATION_ERROR 422 (missing fields), UNAUTHORIZED 401 (invalid credentials), FORBIDDEN 403 (user suspended).

**Business rules:** On 5 failed attempts within 15 min, rate-limit account for 1 hour. Suspended users get FORBIDDEN (403). Use generic "invalid credentials" — do not reveal whether email or password is wrong.

---

### POST /auth/refresh

Exchange a refresh token for a new access token. **Auth:** None (uses refresh token in body)

**Request:** `{ "refresh_token": "a1b2c3d4e5f6...64-byte-hex" }`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "b2c3d4e5f6a7...new-64-byte-hex",
    "expires_in": 900
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Error cases:** VALIDATION_ERROR 422 (missing refresh_token), UNAUTHORIZED 401 (invalid/expired token).

**Business rules:** Rotation-based: each use invalidates old token, issues new pair. If a revoked token is presented, invalidate ALL sessions (theft detection).

---

### POST /auth/logout

Invalidate the current session. **Auth:** Required | **Rate limit:** 30 req/min

**Request:** `{ "refresh_token": "a1b2c3d4e5f6...64-byte-hex" }`

**Response (200):** `{ "success": true, "data": { "message": "Logged out successfully" }, "meta": { "request_id": "req-abc123" } }`

**Error cases:** UNAUTHORIZED 401 (missing/invalid access token).

**Business rules:** JWT access token lives until its 15-min natural expiry. Refresh token is deleted from DB. Idempotent: calling logout with an already-invalidated token still returns 200.
