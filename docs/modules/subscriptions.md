# Module: Subscriptions

## Purpose

Manage the subscription lifecycle for users: activation, expiry, extension, cancellation, and feature entitlement. This module is the source of truth for what a user can and cannot do within the platform.

## Responsibilities

- Activate a subscription after a successful (verified) Midtrans settlement.
- Track subscription status (`active`, `expired`, `cancelled`).
- Calculate and store `start_at` and `expires_at` based on package `duration_days`.
- Auto-activate the Free 7-day subscription on new user registration.
- Detect and process subscription expiry.
- Provide entitlement checks (feature capability queries for other modules).
- Allow admin overrides (extend, terminate, create).
- Handle subscription upgrades (replace current subscription with a new one).

## Non-Responsibilities

- Payment processing (handled by Payments module).
- Package definition/pricing (handled by Packages module).
- Sending expiry reminder notifications (planned but not implemented).
- Renewal billing (no auto-renewal in current scope).

## Actors

- **User (couple):** Owns a subscription, benefits from entitlement checks.
- **System (webhook):** Triggers subscription activation on Midtrans settlement.
- **Admin:** Can create, extend, or terminate any subscription.

## Business Rules

- A user can have at most one active subscription at a time.
- Free subscription: 7 days from registration, same limits as Basic plan.
- Subscription starts only after Midtrans webhook confirms settlement (never from frontend callback).
- Subscription activation is idempotent per transaction — a settled transaction cannot activate twice.
- Basic: 90 days, Premium: 365 days, Pro: lifetime (no expiry, `expires_at` = NULL).
- Expired subscriptions revert features to Free tier.
- Expired subscription: invitation stays published but premium features are locked.
- Guest access to invitations is NOT affected by subscription expiry.
- Editing (event/section updates) is disabled after expiry until renewed.
- On settlement, any existing Free subscription is replaced/overridden by the new paid subscription.
- On settlement, if a paid subscription already exists and is active, it is extended/upgraded (new expiry = max(current expiry, now + duration)).
- 7 days before expiry, a warning should be sent (TODO: implement notification).
- Pro subscription (lifetime) never expires; `expires_at` is NULL.

## Entities

- **Subscription:** `{ id, user_id, package_id, transaction_id, status, start_at, expires_at, created_at, updated_at, deleted_at }`

## Database

- Table: `subscriptions`
- Soft delete enabled.
- Indexes: `user_id`, `status`, `expires_at` (WHERE status = 'active').
- Status values: `active`, `expired`, `cancelled`.
- Partial unique index concept: one active subscription per user enforced in service layer.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/users/me/subscription` | JWT | Get current user's subscription and entitlements |
| POST | `/api/v1/subscriptions/upgrade` | JWT | Create payment intent for upgrading package |

## Request Flow

```
Activation (triggered by Payments webhook)
  → Service: verify transaction status = settlement, verify signature (done in payments)
  → Service: check transaction not already used (idempotency)
  → Service: determine new subscription state (replace Free, extend paid, or create)
  → Service: compute expires_at from package duration_days
  → Service: deactivate existing subscription(s), insert new active subscription
  → Repository: persist in transaction
  → Audit Log: record subscription.activated
```

```
Entitlement check (called by other services)
  → Service: load active subscription for user
  → Service: if none or expired → return Free package capabilities
  → Service: map subscription.package_id → package.features JSONB
  → Service: return capability map (e.g., guest_limit, music.upload)
```

## Validation

- Package must exist and be `is_active = true`.
- Transaction must exist and have status `settlement`.
- Admin extension requires a positive number of days.
- Subscription can only be terminated if it is currently active.

## Authorization

- Users can only view their own subscription.
- Admin can view and modify any subscription.
- Entitlement checks are internal service calls (no HTTP exposure).

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 402 | PAYMENT_REQUIRED | No active subscription when required |
| 404 | NOT_FOUND | Subscription/package not found |
| 409 | CONFLICT | Activation attempted for already-settled transaction |
| 422 | VALIDATION_ERROR | Invalid package or duration |

## Security Considerations

- Never activate a subscription based on frontend response — only the Midtrans webhook is authoritative.
- Subscription activation must be idempotent to prevent double-activation from webhook retries.
- Admin override actions must be logged in audit log with admin identity.
- Entitlement checks must read from the database (or cached DB data), never from client-supplied plan info.
- Expiry detection must be server-side; do not trust client timestamps.
- A user must not be able to downgrade their subscription manually.

## Testing Requirements

- Unit tests for activation logic (new, replace Free, extend paid).
- Unit tests for expiry calculation (90 days, 365 days, lifetime).
- Integration tests for webhook-triggered activation (idempotency on duplicate webhook).
- Test entitlement checks for each package tier.
- Test expired subscription reverts to Free capabilities.
- Test admin extend/terminate operations.
- Test concurrent activation race (two webhook calls at once).

## Dependencies

- Payments module — triggers activation on settlement.
- Packages module — package definition, duration, features.
- Users module — user existence check.
- Audit Log module — action logging.
- A scheduler/job runner for expiry detection (future).

## Related Modules

- **Payments** — Source of subscription activation events.
- **Packages** — Defines duration and capabilities.
- **Events** — Enforces subscription limits on creation/publishing.
- **Guests** — Enforces guest limit from subscription.
- **Admin** — Admin overrides.

## Known Limitations

- No auto-renewal (users must manually repurchase).
- No prorated refunds or partial refunds.
- No grace period after expiry.
- Expiry reminder notifications not yet implemented.
- No subscription history UI beyond the current active subscription.
- No plan downgrade flow (Premium → Basic).
- Expiry is detected lazily (on access) rather than by scheduled job; no background process yet.

## TODO

- [ ] Implement expiry reminder notification (7 days before expiry, email/WhatsApp).
- [ ] Implement background job for marking expired subscriptions.
- [ ] Implement auto-renewal with saved payment method (future phase).
- [ ] Add subscription history view.
- [ ] Implement downgrade flow.
- [ ] Add grace period policy decision.