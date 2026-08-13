# Module: Invitation Editor

## Purpose

Provide the configuration surface for the wedding invitation content. Manages the nine invitation sections, their enable/disable state, and content configuration. The editor drives what guests see on the public invitation page.

## Responsibilities

- Manage the nine invitation sections: Cover, Opening (Pembuka), Couple (Profil), Events (Acara), Gallery (Galeri), RSVP, Gifts (Amplop), Dress Code, Closing (Penutup).
- Toggle section visibility (enabled/disabled).
- Persist section content configuration.
- Provide preview data for the editor UI.
- Enforce feature entitlements per section (e.g., video only for Premium/Pro).
- Auto-create default section configuration when an event is created.

## Non-Responsibilities

- Template/theme selection (handled by Templates module).
- Guest RSVP submission (handled by RSVP module).
- Public rendering of sections (handled by the frontend `[slug]` page).
- Media uploads (handled by Gallery and Music modules).

## Actors

- **User (couple):** Configures sections in the `/user/editor` interface.
- **Guest (public):** Views the configured sections (read-only).

## Business Rules

- Every event has exactly one section configuration record (`event_sections`, 1:1 with event).
- The nine sections are: Cover (hero), Opening, Couple, Events, Gallery, RSVP, Gifts, Dress Code, Closing.
- Sections can be individually enabled/disabled by the couple.
- Default state: hero, couple, events, gallery, rsvp, guestbook enabled; video and digital_gifts disabled.
- Feature entitlements gate certain sections:
  - Gallery photo count limited by plan (`gallery.photo.max`: Free/Basic=5, Premium=20, Pro=unlimited).
  - Video embed requires Premium or Pro.
  - Music upload requires Premium or Pro (Basic = preset only).
  - QRIS gift section requires Pro.
  - Watermark removal is a Pro/Premium feature (template rendering concern).
- Editing is disabled after subscription expiry until renewal.
- Content is saved as structured JSON per section, not free-form HTML (preventing XSS).
- The editor can save draft content; changes do not require re-publishing (published page reflects latest saved content).
- Preview uses the assigned template with draft content.

## Entities

- **EventSection:** `{ id, event_id, hero_enabled, couple_enabled, event_details_enabled, gallery_enabled, video_enabled, music_id, rsvp_enabled, guestbook_enabled, digital_gifts_enabled, created_at, updated_at }`
- **SectionContent** (per section, stored as JSONB in `event_sections` or dedicated content columns): structured fields per section.

## Database

- Table: `event_sections` (1:1 with events, unique event_id).
- Section content stored as JSONB columns or a content table keyed by `(event_id, section_key)`.
- `music_id` references `music` table (selected background track).
- Auto-created on event creation (see Events module).

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events/:id/sections` | JWT | Get all section config + content |
| PUT | `/api/v1/events/:id/sections` | JWT | Update section toggles |
| PUT | `/api/v1/events/:id/sections/:section` | JWT | Update a single section's content |
| GET | `/api/v1/events/:id/preview` | JWT | Preview invitation with current template |

## Request Flow

```
GET /events/:id/sections
  → Handler: parse event ID
  → Service: verify ownership
  → Service: load event_sections record + all section content
  → Service: annotate entitlements (what's locked by plan)
  → Handler: return sections payload
```

```
PUT /events/:id/sections/:section
  → Handler: parse section key + content JSON
  → Service: verify ownership + subscription status (editing allowed?)
  → Service: validate section content against section schema
  → Service: enforce entitlement (e.g., video requires premium)
  → Service: upsert section content in DB
  → Handler: return updated section
```

## Validation

- Section key must be one of the nine known sections.
- Content must validate against the per-section JSON schema.
- Text fields: max lengths per section (e.g., message < 2000 chars).
- URL fields: valid http/https or `wa.me` links.
- Date fields: valid date/time formats.
- Enforced at service layer: content type and structure.

## Authorization

- JWT required; ownership enforced (user must own the event).
- Subscription status checked before allowing edits (expired = read-only).
- Admin read-only access.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not owner, editing disabled (expired) |
| 404 | NOT_FOUND | Event or section not found |
| 422 | VALIDATION_ERROR | Invalid section key or content |
| 402 | PAYMENT_REQUIRED | Editing locked due to no active subscription |

## Security Considerations

- Section content is rendered as data, not raw HTML — sanitize all user-provided content before rendering.
- Never allow script injection via invitation content (escape on render).
- Validate content schemas server-side; do not trust the editor client.
- File URLs in sections (gallery, music) must reference approved storage locations.
- Stored XSS is the primary risk — content must be sanitized both on save and on render.

## Testing Requirements

- Unit tests for section schema validation.
- Integration tests for section CRUD and toggles.
- Test entitlement gating (e.g., video on Basic rejected).
- Test editing locked after subscription expiry.
- Test content sanitization (malicious HTML/JS rejected or escaped).
- Test preview endpoint returns merged template + content.

## Dependencies

- Events module — event ownership and lifecycle.
- Subscriptions module — editing entitlement.
- Packages module — capability limits (gallery count, video, music upload).
- Templates module — template for preview rendering.
- Gallery/Music modules — media referenced by sections.

## Related Modules

- **Events** — Parent entity.
- **Templates** — Visual theme applied to section content.
- **Gallery**, **Music** — Media for gallery/music sections.
- **Digital Gifts** — Gifts section configuration.
- **SEO** — Metadata for public page.

## Known Limitations

- No drag-and-drop section ordering (fixed order).
- No section versioning/history (undo not available).
- No A/B template comparison within editor.
- No collaborative editing.
- Draft vs published content is not versioned — edits apply immediately to published page.
- No visual preview on mobile breakpoints in the editor (frontend concern).

## TODO

- [ ] Add section content versioning/undo history.
- [ ] Add drag-and-drop section reordering.
- [ ] Add template preview switcher.
- [ ] Implement draft/published content separation.
- [ ] Add content validation previews (e.g., character counters).