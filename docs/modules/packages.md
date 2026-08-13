# Module: Packages

## Purpose

Define the SaaS pricing plans (Free, Basic, Premium, Pro), their feature flags, pricing, durations, and limits. Packages are the configuration source of truth for what each tier entitles users to.

## Responsibilities

- Define package plans: code, name, price, duration, guest limit, template group.
- Store feature capabilities as a JSONB feature flag set.
- Activate/deactivate packages (is_active flag).
- Provide package listing for the public pricing page.
- Provide package lookup for subscription activation and entitlement checks.
- Allow admin CRUD of packages.

## Non-Responsibilities

- Subscription lifecycle management (handled by Subscriptions module).
- Payment processing (handled by Payments module).
- Enforcing limits (enforced by consuming modules, e.g., Guests, Gallery).
- Currency conversion (pricing is fixed in IDR).

## Actors

- **Guest/public:** Views available packages on the pricing page.
- **User (couple):** Purchases a package; capabilities apply to their account.
- **Admin:** Creates, updates, and deactivates packages.

## Business Rules

- Package price is stored in IDR (BIGINT, no decimals).
- Duration: `duration_days` — Free=7, Basic=90, Premium=365, Pro=NULL (lifetime).
- Guest limit: Free=100, Basic=100, Premium=500, Pro=NULL (unlimited).
- Template group: Free/Basic='standard', Premium='premium', Pro='all'.
- Only `is_active = true` packages can be purchased.
- Feature capabilities are stored as a JSONB map of capability key → boolean/value (capability-based entitlement, not hardcoded plan names).
- Capability keys follow the `module.capability` convention, e.g., `gallery.photo.max`, `music.upload`, `rsvp.export`.
- Deleting a package is not allowed if it has associated subscriptions; deactivate instead.
- Package code is unique and used for internal logic; name is unique and user-facing.
- Free package cannot be deleted or purchased — it is auto-granted on registration.

## Entities

- **Package:** `{ id, name, code, price, duration_days, guest_limit, template_group, features, is_active, created_at, updated_at }`

## Database

- Table: `packages`
- No soft delete — packages are deactivated via `is_active = false`.
- `features` column: JSONB.
- Unique constraints: `name`, `code`.

### Feature flag structure

```json
{
  "guest.max": 100,
  "theme.max": 5,
  "gallery.photo.max": 5,
  "gallery.video.max": 0,
  "video.enabled": false,
  "music.upload": false,
  "custom_domain": false,
  "watermark.removed": false,
  "whatsapp.bulk": false,
  "guestbook.qr": false,
  "rsvp.export": false,
  "digital_gift.qris": false,
  "template.custom_request": false,
  "event.max": 1
}
```

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/packages` | Public | List active packages (pricing page) |
| GET | `/api/v1/admin/packages` | Admin | List all packages (incl. inactive) |
| POST | `/api/v1/admin/packages` | Admin | Create package |
| PUT | `/api/v1/admin/packages/:id` | Admin | Update package |
| DELETE | `/api/v1/admin/packages/:id` | Admin | Deactivate package (soft) |

## Request Flow

```
GET /packages (public)
  → Handler: no auth
  → Service: load is_active = true packages
  → Service: exclude Free package from public pricing (internal only)
  → Handler: return package list with price and feature summary
```

```
PUT /admin/packages/:id
  → Handler: parse body, validate fields
  → Service: verify admin role, verify package exists
  → Service: validate feature JSONB structure
  → Service: update package, log in audit log
  → Handler: return updated package
```

## Validation

- name: required, unique, max 100 chars.
- code: required, unique, lowercase snake_case, max 50 chars.
- price: positive integer (IDR).
- duration_days: positive integer or null (lifetime).
- guest_limit: positive integer or null (unlimited).
- template_group: one of `standard`, `premium`, `all`.
- features: valid JSON object; known keys validated against schema.

## Authorization

- Public listing: no auth.
- Admin CRUD: JWT with `role = 'admin'`.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT (admin routes) |
| 403 | FORBIDDEN | Not admin |
| 404 | NOT_FOUND | Package not found |
| 409 | CONFLICT | Duplicate name/code, delete with active subscriptions |
| 422 | VALIDATION_ERROR | Invalid price/duration/feature JSON |

## Security Considerations

- Prices and limits are server-authoritative; never trust client-supplied feature flags.
- Free package must always exist and be active (system bootstraps it).
- Do not expose internal codes or feature JSONB to the public pricing endpoint — return a curated summary.
- Admin package changes affect active entitlements immediately; changes must be intentional and audited.

## Testing Requirements

- Unit tests for feature flag parsing and capability lookup.
- Integration tests for public package listing (Free excluded, only active shown).
- Admin CRUD tests (create, update, deactivate).
- Test that deactivated packages cannot be purchased.
- Test duplicate name/code rejection.
- Test that Free package cannot be deleted.

## Dependencies

- Database: `packages` table.
- Audit Log module for admin changes.
- No external services.

## Related Modules

- **Subscriptions** — Uses package definition for activation and entitlement.
- **Payments** — Uses package price for transactions.
- **Admin** — Package management UI.
- **Events**, **Guests**, **Gallery**, **Music** — Enforce capability limits from package features.

## Known Limitations

- No promotional/discount pricing per package (flat price only).
- No trial extension beyond the Free 7-day package.
- Feature JSONB schema is not formally versioned; changes require care.
- No package localization (Indonesian-only descriptions).
- Pricing in IDR only.

## TODO

- [ ] Add discount/promo price fields.
- [ ] Formalize feature JSONB schema versioning.
- [ ] Add package description fields (public marketing copy).
- [ ] Add feature comparison data for pricing page.
- [ ] Consider time-based promotional pricing.