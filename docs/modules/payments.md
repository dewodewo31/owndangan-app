# Module: Payments

## Purpose

Integrate with the Midtrans payment gateway (Snap) to process package purchases. Manages the full transaction lifecycle from Snap token creation through webhook-driven status updates, and triggers subscription activation on settlement.

## Responsibilities

- Create Midtrans Snap transactions for package purchases.
- Generate unique order IDs (`INV-YYYYMMDD-{short_user_hash}-{random8}`).
- Store Snap token and redirect URL for the frontend.
- Receive and verify Midtrans webhook notifications (signature check).
- Map Midtrans statuses to internal transaction statuses.
- Update transaction status in the database.
- Trigger subscription activation on settlement (idempotent).
- Provide transaction history for users and admins.
- Refund handling (cancel subscription on refund).

## Non-Responsibilities

- Hosting the payment page (Midtrans Snap.js runs on the frontend).
- Storing or processing card data (Midtrans handles PCI compliance).
- Frontend payment callback processing (UI hint only, never authoritative).
- Generating invoices or receipts (future).
- Managing bank accounts / QRIS content shown on invitations (Digital Gifts module).

## Actors

- **User (couple):** Creates a payment for a package.
- **Midtrans:** Sends webhook notifications with payment status.
- **Admin:** Views all transactions, monitors webhook processing.

## Business Rules

- Only `is_active = true` packages can be purchased.
- Transaction is created with status `pending` before calling Midtrans Snap.
- Snap token has a 1-hour expiry (Midtrans side).
- Order ID format: `INV-YYYYMMDD-{short_user_hash}-{random8}`, unique globally.
- Webhook signature: `sha512(order_id + status_code + gross_amount + server_key)` must match `signature_key`; reject with 401 on mismatch.
- Always return HTTP 200 to Midtrans to acknowledge receipt, even if processing fails internally.
- If `order_id` is not found in webhook, return 200 and log a warning.
- If internal status is already `settlement`, ignore the duplicate webhook (idempotency).
- On settlement: update transaction, activate/replace subscription, compute `expires_at` from package `duration_days`.
- On settlement, existing Free subscription is overridden; existing paid subscription is extended.
- Transaction statuses: `pending`, `settlement`, `expire`, `deny`, `cancel`, `refund`.
- Do NOT trust frontend payment status — only the webhook activates subscriptions.
- Refund/partial_refund of a settled transaction cancels the associated subscription.

## Entities

- **Transaction:** `{ id, user_id, package_id, order_id, snap_token, gross_amount, status, payment_type, transaction_time, settlement_time, midtrans_response, created_at, updated_at }`

## Database

- Table: `transactions`
- Hard delete not allowed (kept for audit).
- Unique index on `order_id`.
- Index on `user_id`, `status`.
- `midtrans_response` stores the full webhook payload (JSONB) for debugging and audit.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/payments/snap` | JWT | Create Snap transaction |
| POST | `/api/v1/webhook/midtrans` | None (signature) | Midtrans notification |
| GET | `/api/v1/payments/transactions` | JWT | User's transaction history |
| GET | `/api/v1/admin/transactions` | Admin | All transactions |

## Request Flow

```
POST /payments/snap
  → Handler: parse package_id, redirect_url
  → Service: verify user, verify package active
  → Service: generate order_id, create transaction (status=pending)
  → Service: call Midtrans Snap API with server key
  → Service: store snap_token, snap_redirect_url in transaction
  → Handler: return { transaction_id, order_id, snap_token, snap_redirect_url, gross_amount }
```

```
POST /webhook/midtrans
  → Handler: parse raw body, verify signature (sha512)
  → Service: look up transaction by order_id
  → Service: if not found → log warning, return 200
  → Service: if status already settlement → return 200 (idempotent)
  → Service: map Midtrans status → internal status, update transaction
  → Service: if settlement → call Subscriptions.activate(...) in same DB transaction
  → Handler: return 200 { success: true }
```

## Validation

- package_id must be a valid UUID and exist with `is_active = true`.
- Webhook signature must match exactly.
- gross_amount in webhook must match transaction gross_amount (integrity check).
- Midtrans status must be a known value.

## Authorization

- `/payments/snap`: JWT required.
- Webhook: no JWT; verified via signature key (server-to-server).
- Admin transactions: JWT with admin role.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT (snap) |
| 401 | INVALID_SIGNATURE | Webhook signature mismatch |
| 404 | NOT_FOUND | Package not found |
| 422 | VALIDATION_ERROR | Invalid package_id, amount mismatch |
| 500 | INTERNAL_ERROR | Midtrans API failure |
| 402 | PAYMENT_REQUIRED | User suspended / cannot purchase |

## Security Considerations

- Server key is server-side only; client key is safe for the frontend.
- Never log full card numbers or raw webhook payloads to logs (store payload in DB, log only order_id and status).
- Signature verification must use a constant-time comparison.
- Webhook endpoint should be rate-limited and optionally IP-whitelisted.
- Idempotency is critical — duplicate webhook deliveries must not double-activate.
- A settled transaction cannot activate a subscription more than once.
- Do not expose `midtrans_response` raw payload to API consumers — return a sanitized transaction object.
- Midtrans keys configured via environment: `MIDTRANS_SERVER_KEY`, `MIDTRANS_CLIENT_KEY`, `MIDTRANS_IS_PRODUCTION`.

## Testing Requirements

- Unit tests for signature verification (valid and tampered payloads).
- Unit tests for status mapping (settlement, pending, deny, expire, refund).
- Unit tests for idempotency (duplicate settlement webhook).
- Integration tests for full payment flow: snap creation → webhook settlement → subscription activated.
- Test webhook with unknown order_id (returns 200, logs warning).
- Test gross amount mismatch rejection.
- Test refund cancels subscription.
- Test Midtrans API failure returns graceful error.

## Dependencies

- Midtrans Go SDK.
- Subscriptions module — activation on settlement.
- Packages module — price and duration lookup.
- Audit Log module — payment events.
- Environment config: `MIDTRANS_SERVER_KEY`, `MIDTRANS_CLIENT_KEY`, `MIDTRANS_IS_PRODUCTION`.

## Related Modules

- **Subscriptions** — Activated by settlement webhook.
- **Packages** — Purchase target.
- **Users** — Transaction ownership.
- **Admin** — Transaction monitoring.
- **Audit Log** — Payment event logging.

## Known Limitations

- No automatic retry for failed Midtrans API calls (manual retry only).
- No support for partial refunds in subscription logic (full refund only).
- No payment method preferences saved.
- No webhook replay tool for debugging (planned).
- No transaction expiry cleanup job (pending transactions accumulate).
- Sandbox/production switch is a single env flag — no per-transaction override.

## TODO

- [ ] Implement pending transaction cleanup job (expire stale pending transactions).
- [ ] Add webhook replay/debugging tool for admins.
- [ ] Add payment method preference persistence.
- [ ] Implement invoice/receipt generation.
- [ ] Add partial refund handling.
- [ ] Add webhook IP whitelist.