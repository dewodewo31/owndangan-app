# Module: Authentication

## Purpose

Handle user identity across the platform: registration, login, token issuance, and session management. Provides the authentication foundation for both the user dashboard and admin dashboard.

## Responsibilities

- Register new user accounts with email/password.
- Authenticate users via email/password login.
- Issue JWT access tokens (short-lived, 15 minutes).
- Issue and rotate refresh tokens (opaque, 7-day expiry).
- Invalidate sessions on logout.
- Detect token theft via refresh token rotation.
- Rate-limit login attempts to prevent brute-force attacks.
- Provide middleware for JWT validation (used by all protected routes).

## Non-Responsibilities

- User profile management (handled by Users module).
- Role/permission enforcement beyond basic JWT parsing (handled by middleware).
- Password reset flow (future feature, not yet implemented).
- OAuth/social login (future feature, not yet implemented).

## Actors

- **Guest (unauthenticated user):** Can register and login.
- **User (authenticated):** Uses access token to call protected endpoints.
- **Admin (authenticated):** Uses same auth flow; role checked by admin middleware.

## Business Rules

- Password must be at least 8 characters.
- Passwords are hashed with bcrypt (cost factor 12).
- Email uniqueness is enforced case-insensitively (stored as-is, but compared with LOWER()).
- On 5 failed login attempts within 15 minutes, the account is rate-limited for 1 hour.
- Suspended users receive FORBIDDEN (403) on login, not UNAUTHORIZED.
- Login error messages are intentionally generic — do not reveal whether email or password was wrong.
- JWT payload contains: `sub` (user UUID), `role` (user/admin), `iat`, `exp`.
- Access tokens expire after 15 minutes (configurable via `JWT_EXPIRATION_HOURS`).
- Refresh tokens are 64-byte cryptographically random hex strings, stored in DB.
- Refresh token rotation: each use invalidates the old token and issues a new pair.
- If a revoked/rotated refresh token is presented, invalidate ALL sessions for that user (theft detection).
- New users automatically receive a Free 7-day subscription upon registration.

## Entities

- **User:** Core identity record (see Users module).
- **RefreshToken:** `{ id, user_id, token_hash, expires_at, created_at, is_revoked }`

## Database

- `users` table stores credentials (password_hash).
- `refresh_tokens` table stores hashed refresh tokens.
- No session table; JWT is self-contained.

## API

| Method | Endpoint | Auth | Rate Limit |
|--------|----------|------|------------|
| POST | `/api/v1/auth/register` | None | 5 req/min/IP |
| POST | `/api/v1/auth/login` | None | 10 req/min/IP |
| POST | `/api/v1/auth/refresh` | None | 20 req/min/IP |
| POST | `/api/v1/auth/logout` | JWT | 30 req/min |

## Request Flow

```
POST /auth/login
  → Handler: parse JSON body, validate required fields
  → Service: look up user by email, verify bcrypt hash, check account status
  → Service: generate JWT access token, create refresh token in DB
  → Handler: return { user, access_token, refresh_token }
```

```
POST /auth/refresh
  → Handler: parse refresh_token from body
  → Service: verify token hash exists in DB, not revoked, not expired
  → Service: if token already revoked (replay), revoke ALL user sessions
  → Service: generate new access token + new refresh token, store new, revoke old
  → Handler: return { access_token, refresh_token }
```

## Validation

- Email: valid format, max 255 characters, must be unique.
- Password: min 8 characters, max 72 characters (bcrypt limit).
- Phone (optional): must match `62xxxxxxxxxx` format (Indonesia).
- Name: required, max 150 characters, trimmed.
- Refresh token: must be exactly 64 hex characters.

## Authorization

- Registration and login are public (no auth required).
- Refresh endpoint is public but requires a valid refresh token (not JWT).
- Logout requires a valid JWT access token.
- No role-based restrictions on auth endpoints.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 422 | VALIDATION_ERROR | Missing/invalid fields, password too short |
| 409 | CONFLICT | Email already registered |
| 401 | UNAUTHORIZED | Invalid credentials, invalid/expired token |
| 403 | FORBIDDEN | Account suspended |
| 429 | RATE_LIMITED | Too many requests |

## Security Considerations

- Passwords must never be returned in API responses.
- JWT secret must be a strong random value, configured via environment variable `JWT_SECRET`.
- JWT secret rotation plan must be documented.
- Refresh token rotation is mandatory for theft detection.
- Rate limiting on login prevents brute-force attacks.
- bcrypt cost 12 provides adequate hashing strength.
- On registration, the password must be hashed before any response is sent (no plaintext exposure).
- All auth endpoints must use HTTPS in production.

## Testing Requirements

- Unit tests for token generation, validation, and refresh rotation.
- Unit tests for bcrypt hashing and comparison.
- Integration tests for registration, login, refresh, logout flows.
- Test rate limiting behavior.
- Test suspended user login rejection.
- Test theft detection (reused revoked token invalidates all sessions).
- Test concurrent refresh token usage (race condition).

## Dependencies

- JWT library: `golang-jwt/jwt/v5`
- Bcrypt: `golang.org/x/crypto/bcrypt`
- Config: `JWT_SECRET`, `JWT_EXPIRATION_HOURS` environment variables
- Database: `users` and `refresh_tokens` tables

## Related Modules

- **Users** — User profile management, account status.
- **Admin** — Admin user management (suspend, activate).
- **Audit Log** — Login events, failed login attempts.

## Known Limitations

- No OAuth/social login (Google, Facebook, etc.).
- No password reset flow (users must contact admin).
- No MFA/2FA support.
- JWT revocation requires short expiry (15 min) — no server-side blacklist.
- No refresh token cleanup mechanism for expired tokens (DB will accumulate stale rows).

## TODO

- [ ] Implement password reset flow (email-based).
- [ ] Add refresh token cleanup job (cron: delete expired revoked tokens).
- [ ] Consider OAuth 2.0 / social login for future phases.
- [ ] Add login activity notification (email when new device login detected).
- [ ] Document JWT secret rotation procedure.