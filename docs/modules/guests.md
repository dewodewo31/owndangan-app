# Module: Guests

## Purpose

Manage the guest list for each event. Handles guest CRUD, categorization, CSV import, access token generation, and guest limit enforcement.

## Responsibilities

- Create, read, update, and delete guest records.
- Categorize guests (`family`, `friend`, `colleague`, `other`).
- Generate unique access tokens for each guest (used for RSVP).
- Enforce guest limit per subscription plan.
- Import guests via CSV.
- Export guest list (see Export module).
- Link guests to RSVP records.

## Non-Responsibilities

- RSVP submission (handled by RSVP module).
- WhatsApp message sending (handled by WhatsApp module).
- Guestbook messages (handled by Guestbook module).
- RSVP recap aggregation (handled by RSVP module).

## Actors

- **User (couple):** Manages the guest list for their event.
- **Guest (invitee):** Receives invitation link with their token; submits RSVP.
- **Admin:** Read-only access to guest lists.

## Business Rules

- Guests belong to exactly one event.
- Guest count = total non-deleted guest records for an event.
- Guest limit per plan: Free/Basic=100, Premium=500, Pro=unlimited.
- Adding a guest beyond the limit is blocked (LIMIT_EXCEEDED).
- Deleting a guest frees up quota.
- CSV import respects the limit — the whole import is rejected if it would exceed the limit.
- Each guest gets a unique access token (random, unguessable) used to identify their RSVP.
- Token is the only identity mechanism for anonymous guests (no auth).
- A guest can have at most one RSVP record (1:1).
- Phone number is optional but recommended for WhatsApp integration.
- Guest category defaults to `family`.
- `note` is an internal field (couple-only), never exposed publicly.
- Duplicate guests (same name + phone) within an event are allowed unless explicitly deduplicated on import.
- Guest records are soft-deleted.

## Entities

- **Guest:** `{ id, event_id, name, phone, category, note, token, created_at, updated_at, deleted_at }`

## Database

- Table: `guests`
- Soft delete enabled.
- Unique index on `token`.
- Index on `event_id`.
- 1:1 relationship with `rsvps` (guest_id unique in rsvps).

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events/:id/guests` | JWT | List guests (paginated, filter by category/search) |
| POST | `/api/v1/events/:id/guests` | JWT | Add guest |
| PUT | `/api/v1/events/:id/guests/:gid` | JWT | Update guest |
| DELETE | `/api/v1/events/:id/guests/:gid` | JWT | Delete guest |
| POST | `/api/v1/events/:id/guests/import` | JWT | CSV import |

## Request Flow

```
POST /events/:id/guests
  → Handler: parse body, validate fields
  → Service: verify event ownership
  → Service: count current guests, check against plan guest_limit
  → Service: generate unique token
  → Service: create guest record
  → Handler: return created guest (201)
```

```
POST /events/:id/guests/import
  → Handler: parse multipart CSV file
  → Service: verify event ownership
  → Service: parse CSV (validate rows, name/phone format)
  → Service: count existing guests + import rows, reject if exceeds limit
  → Service: generate tokens, bulk insert in transaction
  → Handler: return import summary (imported, skipped, errors)
```

## Validation

- name: required, max 255 chars.
- phone: optional, valid phone format (prefer 62xxx).
- category: one of `family`, `friend`, `colleague`, `other`.
- note: optional, max 1000 chars (internal).
- CSV: headers must include `name`; optional `phone`, `category`, `note`.
- CSV rows with invalid data are skipped with error reporting (partial import allowed, but limit check is on the full batch).

## Authorization

- JWT required; event ownership enforced for all guest operations.
- Guest tokens are not used for authentication of this module's APIs.
- Admin read-only access.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not event owner |
| 404 | NOT_FOUND | Event or guest not found |
| 422 | VALIDATION_ERROR | Invalid guest fields, invalid CSV |
| 422 | LIMIT_EXCEEDED | Guest limit reached or import would exceed |

## Security Considerations

- Guest tokens are secrets — generate with crypto/rand, store hashed or plaintext as needed (token is used for lookup; hash recommended).
- Tokens must never be returned in bulk list responses (return masked/truncated form); expose full token only on creation or via WhatsApp link.
- `note` field is internal — must not appear in public API responses.
- Phone numbers are PII — restrict access, log access carefully, and consider encryption at rest.
- CSV upload must be size-limited and validated to prevent malicious content (formula injection in spreadsheets).
- Prevent enumeration: guest list endpoints require ownership; public endpoints never expose the list.

## Testing Requirements

- Unit tests for token generation uniqueness and format.
- Unit tests for limit enforcement (add, delete frees quota).
- Integration tests for guest CRUD with ownership.
- CSV import tests: valid, malformed rows, partial failures, limit rejection.
- Test that `note` is never exposed publicly.
- Test token masked in list responses.

## Dependencies

- Events module — event ownership and existence.
- Subscriptions module — guest_limit entitlement.
- Audit Log module — guest operations (bulk import events).
- Export module — guest list export.

## Related Modules

- **RSVP** — 1:1 guest RSVP.
- **WhatsApp** — Uses guest phone + token to generate invite links.
- **Events** — Parent entity.
- **Export** — Guest list CSV/Excel export.

## Known Limitations

- No duplicate detection on manual add (dedupe only planned on import).
- No guest grouping/tables (e.g., bride-side vs groom-side).
- No guest status tracking beyond RSVP (e.g., invitation sent/not sent).
- CSV import is all-or-nothing on limit check but partial on row errors — may confuse users.
- No bulk edit or bulk delete.
- No search by token.

## TODO

- [ ] Add duplicate detection on import (name+phone match).
- [ ] Add guest grouping (bride/groom side) field.
- [ ] Add invitation sent/not-sent status tracking.
- [ ] Add bulk delete and bulk category update.
- [ ] Add CSV import template download.
- [ ] Consider phone normalization and deduplication on import.