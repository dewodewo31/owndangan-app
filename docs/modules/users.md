# Module: Users

## Purpose

Manage user account information, profile settings, and account lifecycle. Provides the identity layer for both regular users (couples) and administrators.

## Responsibilities

- Retrieve and update user profile (name, email, phone, avatar).
- Change password.
- Manage account status (active, suspended).
- Provide user data to other modules (e.g., subscription, events).
- Enforce data ownership boundaries (users can only access their own data).

## Non-Responsibilities

- Authentication (login, register, token management) — handled by Authentication module.
- Role assignment — handled by Admin module for admin-created accounts.
- User deletion — soft delete is handled within this module; hard delete is not supported.
- Email verification — not implemented (future feature).

## Actors

- **User (couple):** Can view and update their own profile.
- **Admin:** Can view any user profile, update user status, change user role.

## Business Rules

- A user can only access their own profile data (ownership enforced at service layer).
- Email changes require uniqueness check (case-insensitive).
- Phone number must be a valid Indonesia format (`62xxxxxxxxxx`).
- Password change requires current password verification.
- Suspended users cannot log in but their public invitations remain accessible (publihed content is not affected by user status).
- Soft delete preserves data but prevents login.
- Admin can view any user's profile but cannot change the password.
- Each user has exactly one role: `user` or `admin`.

## Entities

- **User:** `{ id, name, email, password_hash, phone, role, status, avatar_url, created_at, updated_at, deleted_at }`

## Database

- Table: `users`
- Soft delete enabled (via `deleted_at`).
- Unique indexes: `email` (WHERE deleted_at IS NULL), `phone` (WHERE deleted_at IS NULL).
- Status values: `active`, `suspended`.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/users/me` | JWT | Get current user profile |
| PUT | `/api/v1/users/me` | JWT | Update profile (name, email, phone, avatar) |
| PUT | `/api/v1/users/me/password` | JWT | Change password |

## Request Flow

```
GET /users/me
  → Handler: extract user ID from JWT context
  → Service: look up user by ID from repository
  → Handler: return user profile (exclude password_hash)
```

```
PUT /users/me
  → Handler: parse JSON body, validate fields
  → Service: verify email uniqueness if changed
  → Service: update user fields in repository
  → Handler: return updated user profile
```

## Validation

- Name: required, max 150 characters.
- Email: valid format, max 255 characters, unique across active users.
- Phone: optional, must match `62xxxxxxxxxx` format.
- Avatar URL: optional, must be a valid URL (max 500 characters).
- Password change: current_password required, new_password min 8 chars.

## Authorization

- JWT required for all user endpoints.
- Users can only access their own profile (user_id from JWT must match).
- Admin can access any user profile via admin endpoints.
- Service layer enforces ownership before returning data.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing or invalid JWT |
| 404 | NOT_FOUND | User not found |
| 409 | CONFLICT | Email or phone already taken |
| 422 | VALIDATION_ERROR | Invalid field values |
| 403 | FORBIDDEN | Wrong password on change |

## Security Considerations

- `password_hash` must never be returned in any API response.
- Current password must be verified before allowing password change.
- Email/phone uniqueness must be enforced at the database level, not just application level.
- User profile updates are logged in audit log (except password changes — log only that a change occurred, not the new password).
- Soft delete preserves referential integrity for events, subscriptions, and transactions.

## Testing Requirements

- Unit tests for profile update, password change, email uniqueness.
- Integration tests for ownership enforcement (user A cannot access user B's profile).
- Test soft delete behavior (user deleted but events remain).
- Test suspended user behavior.
- Test concurrent email change race condition.

## Dependencies

- Authentication module for JWT middleware.
- Database: `users` table.
- Audit log service for logging profile changes.

## Related Modules

- **Authentication** — Login, registration, token management.
- **Admin** — Admin user management (list, suspend, activate).
- **Subscriptions** — User subscription lookup.
- **Events** — User's invitation ownership.

## Known Limitations

- No email verification on registration or email change.
- No phone verification.
- No profile picture upload endpoint (avatar_url is set externally via storage).
- No bulk user operations (admin must manage users individually).
- No account deletion endpoint (users must contact admin).

## TODO

- [ ] Implement email verification flow.
- [ ] Add phone verification (OTP) for WhatsApp integration.
- [ ] Add profile picture upload endpoint.
- [ ] Add account deletion endpoint (with data export).
- [ ] Implement email change notification.