# Module: Digital Gifts

## Purpose

Configure and display digital gift information (amplop digital) on the invitation. Allows couples to share bank account details, e-wallet information, and QRIS codes for guests who wish to send monetary gifts. Note: this module handles information display only — the platform does NOT process digital gift payments.

## Responsibilities

- Store bank account information (array of bank name, account number, account holder).
- Store e-wallet information (provider, ID/phone).
- Store QRIS image URL for static QR code.
- Provide gift/thank you message from the couple.
- Display gift information on the public invitation page.
- Enforce feature entitlements (bank only, e-wallet, QRIS) based on plan.

## Non-Responsibilities

- Processing payments or transfers (no transaction handling).
- Confirming receipt of gifts (couple manages this externally).
- Refunding gifts (not applicable — platform does not handle money).
- Verifying bank account numbers (couple is responsible for accuracy).

## Actors

- **User (couple):** Configures gift information and preferences.
- **Guest (invitee):** Views gift information on the invitation (no interaction).
- **Admin:** Read-only access.

## Business Rules

- Digital gift info is static content displayed on the invitation — it is NOT a payment system.
- The platform does not process, track, or confirm any monetary transfers.
- Basic plan: bank accounts only (up to 2 accounts).
- Premium plan: bank accounts + e-wallet (up to 3 accounts + 2 e-wallets).
- Pro plan: bank accounts + e-wallet + QRIS image (unlimited accounts).
- QRIS is only available for Pro plans (`digital_gift.qris` capability).
- Bank accounts stored as JSONB array of `{ bank_name, account_number, account_holder }`.
- e-Wallet stored as JSONB object: `{ provider, id, name }`.
- QRIS image is uploaded as an image; stored as a URL.
- Gift message is optional text displayed alongside the gift section.
- One digital gift configuration per event (1:1).
- Auto-created when the event is created (defaults: empty, disabled until configured).

## Entities

- **DigitalGift:** `{ id, event_id, bank_accounts, ewallet, qris_image_url, gift_message, created_at, updated_at }`

## Database

- Table: `digital_gifts`
- 1:1 with event (unique `event_id`).
- `bank_accounts`: JSONB array.
- `ewallet`: JSONB.
- `qris_image_url`: TEXT (nullable).
- `gift_message`: TEXT (nullable).

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events/:id/digital-gifts` | JWT | Get gift configuration |
| PUT | `/api/v1/events/:id/digital-gifts` | JWT | Update gift configuration |
| GET | `/api/v1/public/events/:slug/digital-gifts` | Public | Get public gift info |

## Request Flow

```
PUT /events/:id/digital-gifts
  → Handler: parse JSON body (bank_accounts, ewallet, qris_image_url, gift_message)
  → Service: verify event ownership
  → Service: check subscription entitlement (bank/e-wallet/QRIS per plan)
  → Service: validate bank account structure
  → Service: upsert digital_gifts record
  → Handler: return updated configuration
```

```
GET /public/events/:slug/digital-gifts
  → Handler: parse slug
  → Service: resolve event, verify published
  → Service: load digital_gifts record
  → Service: check if section is enabled (event_sections.digital_gifts_enabled)
  → Handler: return gift info (or 404 if disabled)
```

## Validation

- bank_accounts: array of objects, each with `bank_name` (required), `account_number` (required), `account_holder` (required).
- bank_name: max 50 chars.
- account_number: number/string, max 30 chars.
- account_holder: max 255 chars.
- ewallet: object with `provider` (required, max 50), `id` (required, max 100), `name` (optional, max 255).
- qris_image_url: valid URL, must point to approved storage.
- gift_message: optional, max 2000 chars.
- Max bank accounts: 2 (Basic), 3 (Premium), unlimited (Pro).

## Authorization

- Configuration: JWT + event ownership.
- Public reading: no auth (published event only).
- Admin: read-only.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not owner, or feature not allowed by plan |
| 404 | NOT_FOUND | Event not found, or gift section disabled |
| 422 | VALIDATION_ERROR | Invalid bank account structure, exceeds limits |

## Security Considerations

- Bank account numbers are sensitive — minimize exposure in logs and API responses.
- Bank account numbers should be encrypted at rest or stored with restricted access (future).
- QRIS image upload must be validated (type, size) and stored in approved storage.
- The gift message is user text — sanitize HTML to prevent XSS.
- Public endpoint returns bank account numbers in full (guests need them to transfer). This is by design but should be flagged in security review.
- No payment processing means no PCI compliance scope for this module.

## Testing Requirements

- Unit tests for bank account structure validation.
- Unit tests for entitlement enforcement (Basic blocks e-wallet, Pro allows QRIS).
- Integration tests for configuration CRUD.
- Test public endpoint returns 404 when gift section is disabled.
- Test account limit enforcement per plan.
- Test QRIS image URL validation.

## Dependencies

- Events module — event ownership, slug resolution.
- Subscriptions module — feature entitlement (bank, e-wallet, qris).
- Storage module — QRIS image upload.
- Invitation Editor module — digital_gifts_enabled toggle.

## Related Modules

- **Events** — Parent entity.
- **Invitation Editor** — Gift section toggle.
- **Payments** — Note: payment processing is separate; digital gifts do not involve Midtrans.

## Known Limitations

- No guest-to-couple gift confirmation flow.
- No gift amount tracking (couple must track manually).
- No PayLater/transfer code integration.
- QRIS is static image only (no dynamic QRIS).
- Bank account validation is format-only (no IBAN/bank account number verification).
- No support for international transfers.

## TODO

- [ ] Add bank account number encryption at rest.
- [ ] Add gift confirmation message from guests (e.g., "I've sent a gift").
- [ ] Add gift amount tracking (optional, user-input).
- [ ] Add dynamic QRIS generation (future).
- [ ] Consider platform escrow for gift handling (future, high complexity).