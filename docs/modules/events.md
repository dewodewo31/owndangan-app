# Module: Events

## Purpose

Manage the core entity of the platform — the wedding invitation (event). Handles CRUD operations, slug management, publish/unpublish workflow, and status transitions. This module is the backbone that all other invitation-content modules attach to.

## Responsibilities

- Create, read, update, and delete (soft) wedding invitation events.
- Auto-generate and manage unique slugs.
- Manage event lifecycle: `draft` → `published` → `unpublished`.
- Enforce subscription-based event limits.
- Validate minimal data before publishing.
- Track view counts for published events.
- Provide event data for the public invitation page.

## Non-Responsibilities

- Section content configuration (handled by Invitation Editor module).
- Template assignment (handled by Templates module).
- Guest list management (handled by Guests module).
- SEO metadata management (handled by SEO module).
- Theming/visual layout (handled by Templates module).

## Actors

- **User (couple):** Creates and manages their own events.
- **Guest (public):** Views published events via `/[slug]`.
- **Admin:** Can view any event (read-only access).

## Business Rules

- Each event belongs to exactly one user.
- Event starts in `draft` status.
- Event limits by plan: Free=1, Basic=3, Premium=10, Pro=unlimited.
- Creating an event beyond the plan limit returns LIMIT_EXCEEDED (422).
- Slug: 3–100 chars, lowercase letters, numbers, hyphens; globally unique.
- Slug is auto-generated from `title`/`couple_name` on creation (can be customized).
- Slug can be changed while remaining unique; changing it invalidates the old URL (no redirect).
- Publish requires: active subscription + minimal data (title, couple_name, wedding_date).
- Publishing sets `status = 'published'`, `published_at = now()`.
- Publishing an already-published event returns CONFLICT (409).
- Unpublish sets `status = 'unpublished'` and clears `published_at`; published invitations return 404 publicly.
- Soft delete sets `deleted_at`; if published, unpublish first.
- Expired subscriptions: event stays published but premium features are locked; editing disabled.
- Guest access is not affected by subscription expiry.
- `event_sections` and `digital_gifts` records are auto-created when the event is created.
- Template can be assigned per event (not per subscription).

## Entities

- **Event:** `{ id, user_id, template_id, title, slug, couple_name, groom_name, bride_name, groom_parents, bride_parents, wedding_date, wedding_time, ceremony_venue, ceremony_address, ceremony_map_url, reception_venue, reception_address, reception_map_url, status, published_at, view_count, created_at, updated_at, deleted_at }`
- **EventSection:** 1:1 with event (see Invitation Editor module).

## Database

- Table: `events`
- Soft delete enabled.
- Unique index on `slug` (WHERE deleted_at IS NULL).
- Index on `user_id`, `status`.
- Status values: `draft`, `published`, `unpublished`.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events` | JWT | List user's events (paginated, filter by status) |
| POST | `/api/v1/events` | JWT | Create event |
| GET | `/api/v1/events/:id` | JWT | Get event details |
| PUT | `/api/v1/events/:id` | JWT | Update event (slug not changed via PUT) |
| DELETE | `/api/v1/events/:id` | JWT | Soft delete event |
| POST | `/api/v1/events/:id/publish` | JWT | Publish invitation |
| POST | `/api/v1/events/:id/unpublish` | JWT | Unpublish invitation |
| GET | `/api/v1/public/events/:slug` | Public | Get public invitation |

## Request Flow

```
POST /events
  → Handler: parse JSON body, validate fields
  → Service: get user ID from JWT context
  → Service: check subscription is active (else 402)
  → Service: check event limit for plan (else LIMIT_EXCEEDED)
  → Service: generate unique slug from title/couple_name
  → Service: create event + event_sections + digital_gifts in a DB transaction
  → Handler: return created event (201)
```

```
POST /events/:id/publish
  → Handler: parse event ID
  → Service: verify event ownership
  → Service: check active subscription (else 402)
  → Service: validate minimal data (title, couple_name, wedding_date)
  → Service: check status is draft/unpublished (else 409)
  → Service: set status=published, published_at=now()
  → Handler: return public_url
```

## Validation

- title: required, max 255 chars.
- couple_name: optional, max 255 chars.
- slug (when provided): 3–100 chars, `^[a-z0-9-]+$`.
- wedding_date: valid date.
- groom/bride names and parents: optional, max 255 chars.
- Venue/address: optional, max 1000 chars.
- map_url: valid URL (http/https).

## Authorization

- JWT required for all event management endpoints.
- Ownership enforced: users can only access their own events (403 otherwise).
- Admin read-only access to any event.
- Public endpoint requires the event to be `published` and not soft-deleted.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 402 | PAYMENT_REQUIRED | No active subscription |
| 403 | FORBIDDEN | Not the owner |
| 404 | NOT_FOUND | Event not found |
| 409 | CONFLICT | Slug taken, already published/unpublished |
| 422 | VALIDATION_ERROR | Invalid fields, LIMIT_EXCEEDED |

## Security Considerations

- Ownership checks must happen at service layer, never trust client-supplied user IDs.
- Slug uniqueness enforced at DB level (unique index) plus application check.
- Published events are publicly readable; ensure no private fields leak (e.g., internal notes).
- Deleting an event must cascade (soft) or handle dependent records (guests, sections, gallery).
- Publish state changes logged in audit log.

## Testing Requirements

- Unit tests for slug generation and uniqueness.
- Unit tests for event limit enforcement per plan.
- Integration tests for CRUD with ownership.
- Test publish/unpublish workflow.
- Test publish without subscription (402).
- Test publish without minimal data.
- Test public endpoint returns 404 for unpublished/deleted events.
- Test soft delete behavior with dependencies.

## Dependencies

- Subscriptions module — active subscription and event limit checks.
- Packages module — event.max capability.
- Templates module — template assignment.
- Audit Log module — publish/unpublish/delete logging.
- Analytics module — view count increments.

## Related Modules

- **Invitation Editor** — Section content for the event.
- **Templates** — Theme applied to the event.
- **Guests** — Guest list per event.
- **RSVP**, **Guestbook**, **Digital Gifts**, **Gallery**, **Music** — Event-bound content modules.
- **SEO** — Slug/metadata for public pages.
- **Analytics** — View counting.

## Known Limitations

- Slug change does not redirect old URLs.
- No event duplication/copy feature.
- No scheduled publishing (draft date).
- No event-level collaboration (only the owner can edit).
- No multiple template preview comparison (see Editor module).
- view_count is a simple counter, not detailed analytics (see Analytics module).

## TODO

- [ ] Implement slug change redirect (301) or 404 handling decision.
- [ ] Add event duplication/copy feature.
- [ ] Add scheduled publish feature.
- [ ] Add co-editor collaboration (future).
- [ ] Add event-level settings (timezone, language).