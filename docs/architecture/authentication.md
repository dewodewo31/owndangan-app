# Authentication Architecture

## Overview

Authentication uses JWT (JSON Web Tokens) with access and refresh token pattern. Backend issues tokens, frontend stores them and includes them in API requests.

## Authentication Flow

```
User → Login Form → POST /api/v1/auth/login → Validate credentials →
Generate JWT (access token) + Refresh token → Return tokens →
Frontend stores tokens → Frontend includes Authorization: Bearer header →
Backend middleware validates JWT on every request
```

## Token Design

### Access Token

- Format: JWT (HS256)
- Payload: `{ user_id, email, role, iat, exp }`
- Expiry: 15 minutes (configurable)
- Storage: httpOnly cookie or in-memory (frontend)
- Sent in: `Authorization: Bearer <token>` header

### Refresh Token

- Format: Opaque string (not JWT)
- Storage: httpOnly cookie (preferred) or localStorage
- Expiry: 7 days
- Rotation: old refresh token is invalidated when new one is issued
- Usage: `POST /api/v1/auth/refresh`

## Registration

```
POST /api/v1/auth/register
Body: { name, email, password, phone (optional) }
Validation: email uniqueness, password strength (min 8 chars), valid email format
Response: { user, access_token, refresh_token }
Side effects: Creates user record, assigns Free subscription
```

## Login

```
POST /api/v1/auth/login
Body: { email, password }
Validation: email exists, password matches hash
Response: { user, access_token, refresh_token }
Rate limiting: 5 attempts per minute per email
```

## Token Refresh

```
POST /api/v1/auth/refresh
Body: { refresh_token }
Validation: refresh token exists and not expired
Response: { access_token, refresh_token (rotated) }
```

## Logout

```
POST /api/v1/auth/logout
Auth: Valid access token
Action: Invalidates refresh token
```

## Password Security

- Algorithm: bcrypt
- Cost factor: 12
- Password never returned in API responses
- Password change requires current password verification

## Middleware

The JWT middleware:
1. Extracts `Authorization: Bearer <token>` header
2. Validates JWT signature and expiry
3. Extracts user_id, email, role from payload
4. Sets user context on request context
5. On failure: returns 401 Unauthorized

## Admin Middleware

Extends JWT middleware:
1. Validates JWT (same as above)
2. Checks `role == "admin"` in token payload
3. On failure: returns 403 Forbidden

## Related Documentation

- `docs/api/authentication.md`
- `docs/modules/authentication.md`
- `docs/security/authentication-security.md`
- `docs/frontend/authentication.md`
- `docs/backend/middleware.md`