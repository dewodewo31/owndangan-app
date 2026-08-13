# Module: SEO

## Purpose

Optimize the public invitation pages for search engine discoverability and social media sharing. Manages metadata, OpenGraph tags, canonical URLs, slug validation, and sitemap generation.

## Responsibilities

- Generate page title and description for each invitation.
- Set OpenGraph (OG) tags for social media preview (Facebook, WhatsApp, Twitter, etc.).
- Set Twitter Card tags for Twitter previews.
- Generate canonical URL for each invitation page.
- Validate slug format and uniqueness.
- Generate sitemap.xml for search engines.
- Generate robots.txt.
- Ensure metadata is available server-side (SSR/SSG) for SEO.
- Manage 404 handling for unpublished/deleted invitations.

## Non-Responsibilities

- Search engine ranking (SEO is about optimization, not ranking guarantees).
- Analytics tracking (handled by Analytics module).
- Content generation (content comes from the Invitation Editor module).
- URL shortening (not in scope).

## Actors

- **User (couple):** Sets invitation title and couple name (used for metadata).
- **Guest (public):** Consumes the SEO-optimized page.
- **Search Engine (Google, etc.):** Crawls and indexes invitation pages.
- **Social Media (Facebook, WhatsApp, Twitter):** Renders OG/Twitter card previews.

## Business Rules

- Page title format: `{couple_name} - Wedding Invitation | Platform Name`.
- Page description: auto-generated from event details (couple names, date, venue) or custom.
- OG image: the cover photo of the invitation (first gallery photo or template default).
- OG URL: canonical URL of the invitation (`https://undangan.example.com/{slug}`).
- OG type: `website` (or `article` for future).
- Twitter card: `summary_large_image`.
- Slug must be: 3–100 characters, lowercase letters, numbers, hyphens only.
- Slug must be globally unique (enforced by Events module).
- Slug changes invalidate old URLs with no automatic redirect (see Events module).
- Sitemap includes all published invitations with `lastmod` date.
- 404 returned for unpublished, deleted, or non-existent slugs.
- Metadata is rendered server-side by Next.js (Server Components).
- Template-level metadata (e.g., theme name) is added to page metadata.

## Entities

- No dedicated entities. SEO metadata is derived from:
  - `events` table: slug, couple_name, title, wedding_date.
  - `gallery_photos` table: first photo for OG image.
  - `templates` table: template name, default OG image.
  - Platform configuration: base URL, platform name.

## Database

- Reads from: `events`, `gallery_photos`, `templates`.
- No writes specific to SEO (metadata is computed from existing data).
- Future: `seo_metadata` table for custom OG title/description overrides.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/public/events/:slug/metadata` | Public | Get SEO metadata for an invitation |
| GET | `/sitemap.xml` | Public | XML sitemap of all published invitations |
| GET | `/robots.txt` | Public | Robots configuration |

## Request Flow

```
GET /public/events/:slug/metadata
  → Handler: parse slug
  → Service: resolve event, verify published
  → Service: load couple_name, title, wedding_date, cover photo
  → Service: build OG tags, Twitter card, title, description
  → Handler: return metadata object
```

```
GET /sitemap.xml
  → Handler: no auth
  → Service: load all published events (slug, updated_at)
  → Service: build XML sitemap with URLs
  → Handler: return XML (Content-Type: application/xml)
```

## Validation

- Slug format: `^[a-z0-9-]{3,100}$` — validated on event creation (Events module).
- Metadata description length: max 160 characters (OG standard).
- OG image URL: must be from approved storage or a valid URL.
- Sitemap URLs: max 50,000 URLs per sitemap (standard limit).

## Authorization

- All SEO endpoints are public (no auth).
- Sitemap and robots.txt are fully public.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 404 | NOT_FOUND | Event slug not found or not published |
| 422 | VALIDATION_ERROR | Invalid slug format (should not happen from valid data) |

## Security Considerations

- Metadata must not leak private information (e.g., internal notes, guest counts).
- OG image must be from trusted storage to prevent hotlinking abuse.
- Slug enumeration protection: 404 response for unpublished slugs is identical to non-existent slugs.
- Sitemap only includes published events — no unpublished or deleted events.
- Metadata API is rate-limited (public endpoint).
- No user-editable metadata fields in current scope (future: custom OG title/description — must sanitize).

## Testing Requirements

- Unit tests for metadata generation (title, description, OG tags).
- Unit tests for slug format validation.
- Integration tests for metadata endpoint.
- Integration tests for sitemap generation.
- Test that unpublished events return 404 on metadata endpoint.
- Test that deleted events are excluded from sitemap.

## Dependencies

- Events module — slug, couple_name, title, status.
- Gallery module — cover photo for OG image.
- Templates module — default OG image fallback.
- Frontend (Next.js) — renders metadata client-side (Server Components).

## Related Modules

- **Events** — Source of slug, title, couple_name.
- **Gallery** — Cover photo for OG image.
- **Templates** — Default OG image, theme metadata.
- **Invitation Editor** — Custom meta description (future).

## Known Limitations

- No per-event custom OG title/description override (computed from event data).
- No structured data (JSON-LD) for rich search results.
- No hreflang tags (Indonesia-only, Indonesian language).
- No automatic URL redirect on slug change (404 for old URL).
- Sitemap is dynamically generated (no caching yet — consider Redis for high traffic).
- No breadcrumb structured data.

## TODO

- [ ] Add per-event custom OG title/description override.
- [ ] Add JSON-LD structured data (Event, Person, WeddingEvent schema).
- [ ] Implement sitemap caching.
- [ ] Add 301 redirect support for slug changes.
- [ ] Add breadcrumb structured data.
- [ ] Add canonical URL enforcement in frontend.