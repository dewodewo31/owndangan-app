# Module: Export

## Purpose

Provide data export functionality for the couple's dashboard. Enables downloading RSVP data and guest lists in spreadsheet formats (Excel/CSV) for offline record-keeping and planning.

## Responsibilities

- Export RSVP data to Excel (.xlsx) format.
- Export guest list to Excel (.xlsx) or CSV format.
- Include all relevant fields in exports (name, phone, category, RSVP status, guest count, message).
- Respect plan entitlements (RSVP export is a Pro feature).
- Generate downloadable files with proper headers and encoding.
- Support large exports (paginated or streamed).

## Non-Responsibilities

- Exporting platform-wide analytics (admin dashboard has its own metrics).
- Exporting invitation content or design data.
- Scheduled/automated exports (manual download only).
- PDF generation (invitation as PDF is not in scope).
- Email delivery of exports (user downloads directly).

## Actors

- **User (couple):** Downloads RSVP recap and guest list for their event.
- **Admin:** Can trigger exports for any event (read-only).

## Business Rules

- RSVP export is available on Pro plan (`rsvp.export` capability).
- Guest list export is available on all plans (basic guest management).
- Exports include only non-deleted guests.
- RSVP export includes: guest name, phone, category, attendance status, guest count, message, submitted_at.
- Guest list export includes: guest name, phone, category, note, RSVP status (if available), created_at.
- File format: `.xlsx` (Excel) with proper column headers (Indonesian labels).
- File encoding: UTF-8 with BOM for Excel compatibility.
- File generation is synchronous for reasonable sizes (< 5000 rows); above that, consider async (future).
- Maximum export rows: 10,000 (capped at DB query level).
- Rate limit: 5 exports per hour per event.

## Entities

- No dedicated entities. Exports read from:
  - `guests` table (name, phone, category, note, created_at).
  - `rsvps` table (attendance, guest_count, message, submitted_at).

## Database

- Read-only: `guests` (with `deleted_at IS NULL`), `rsvps` (JOIN on guest_id).
- No writes.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events/:id/export/rsvp` | JWT | Download RSVP as Excel |
| GET | `/api/v1/events/:id/export/guests` | JWT | Download guest list as Excel |

## Request Flow

```
GET /events/:id/export/rsvp
  → Handler: parse event ID, query params (format: xlsx/csv)
  → Service: verify event ownership
  → Service: check rsvp.export capability (Pro)
  → Service: load event + guests + rsvps (JOIN query)
  → Service: build Excel/CSV file in memory
  → Handler: return file with Content-Disposition: attachment
```

```
GET /events/:id/export/guests
  → Handler: parse event ID
  → Service: verify event ownership
  → Service: load non-deleted guests
  → Service: build Excel/CSV file
  → Handler: return file download
```

## Validation

- event_id: must be valid UUID and owned by the user.
- format: `xlsx` (default) or `csv` (optional query param).
- Rate limit: check export count per hour.

## Authorization

- JWT required; event ownership enforced.
- RSVP export additionally requires `rsvp.export` capability (Pro).
- Guest list export: available to all plans (no additional check).
- Admin: can export any event.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not owner, or RSVP export not allowed by plan |
| 404 | NOT_FOUND | Event not found |
| 429 | RATE_LIMITED | Too many exports |
| 422 | LIMIT_EXCEEDED | Too many rows for synchronous export |

## Security Considerations

- Exports contain PII (guest names, phone numbers, RSVP data) — must be restricted to the event owner.
- Generated files must be served as downloads (Content-Disposition: attachment), not rendered inline.
- Temporary export files should not be stored on disk (generate in memory, stream to response).
- If buffering to disk, ensure files are cleaned up after response.
- Phone numbers in exports are PII — consider the data protection implications of bulk download.
- CSV files must not contain formula injection (start with `=`, `+`, `-`, `@`). Use Excel-compatible escaping or prefix with tab.
- Rate limit exports to prevent data scraping.

## Testing Requirements

- Unit tests for Excel/CSV generation with correct headers and data.
- Unit tests for formula injection prevention in CSV.
- Integration tests for RSVP export with entitlement check.
- Integration tests for guest list export.
- Test rate limiting.
- Test large export behavior (5000+ rows).
- Test that deleted guests are excluded.

## Dependencies

- Excel library (e.g., `excelize` for Go).
- Events module — event ownership.
- Guests module — guest data.
- RSVP module — RSVP data.
- Subscriptions module — `rsvp.export` capability.

## Related Modules

- **RSVP** — RSVP data export.
- **Guests** — Guest list export.
- **Subscriptions** — Export entitlement.
- **Events** — Parent entity.

## Known Limitations

- No async export for large datasets (may timeout for very large guest lists).
- No PDF export (invitation as PDF).
- No export of guestbook messages.
- No export of digital gift configurations.
- No scheduled/automated exports.
- No export of invitation analytics.
- CSV export may have issues with special characters in Excel (BOM helps but not perfect).

## TODO

- [ ] Add async export for large datasets (background job, download link by email).
- [ ] Add PDF export of invitation (future).
- [ ] Add guestbook message export.
- [ ] Add analytics export.
- [ ] Add scheduled export (daily/weekly, sent to email).
- [ ] Add column selection for exports (choose which fields to include).