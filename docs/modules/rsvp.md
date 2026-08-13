# Module: RSVP

## Purpose

Manage guest attendance responses to the wedding invitation. Handles RSVP submission from guests (anonymous, token-based), status updates, guest count tracking, and recap aggregation for the couple.

## Responsibilities

- Accept RSVP submissions from guests (public endpoint).
- Store attendance status: `yes`, `no`, `maybe`.
- Track guest count (number of people attending).
- Allow guests to update their existing RSVP (one response per guest).
- Provide RSVP recap/statistics for the couple (total invited, attending, declining, pending).
- Provide RSVP list for the couple's dashboard.
- Store optional guest message with the RSVP.

## Non-Responsibilities

- Guest list management (handled by Guests module).
- Guestbook messages (handled by Guestbook module — RSVP message is separate and brief).
- WhatsApp reminders (handled by WhatsApp module).
- Export of RSVP data (handled by Export module).

## Actors

- **Guest (invitee):** Submits or updates their RSVP via the public invitation page.
- **User (couple):** Views RSVP recap and individual responses.
- **Admin:** Read-only access.

## Business Rules

- Each guest can submit at most one RSVP (1:1, enforced by unique guest_id).
- Subsequent submissions from the same guest UPDATE the existing RSVP, not create a new one.
- RSVP submission is anonymous — no login required; identity comes from the guest token.
- The guest token must correspond to a valid guest record of a published event.
- Attendance values: `yes`, `no`, `maybe`.
- `guest_count` is the number of people attending (default 1); relevant only when attendance = `yes`.
- RSVP is only accepted for published events.
- Couple can view recap: totals per status, attendance percentage.
- RSVP recap is visible in the user dashboard and can be exported (Export module, Pro feature).
- RSVP submission creates an analytics event (`rsvp_submitted`).
- RSVP records are immutable snapshots once updated — keep latest state (single row per guest).
- RSVP message is optional and distinct from guestbook messages.

## Entities

- **RSVP:** `{ id, guest_id, event_id, attendance, guest_count, message, submitted_at, updated_at }`

## Database

- Table: `rsvps`
- Unique constraint on `guest_id` (1:1 per guest).
- Index on `event_id`.
- No soft delete (immutable snapshot, kept for history).
- `attendance` values: `yes`, `no`, `maybe`.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/public/rsvps` | Public (token) | Submit or update RSVP |
| GET | `/api/v1/events/:id/rsvps` | JWT | List RSVPs for event (couple) |
| GET | `/api/v1/events/:id/rsvps/stats` | JWT | RSVP recap statistics |
| GET | `/api/v1/events/:id/rsvps/export` | JWT | Export RSVP to Excel (Pro) |

## Request Flow

```
POST /public/rsvps
  → Handler: parse body (token, attendance, guest_count, message)
  → Service: look up guest by token
  → Service: verify guest's event is published (else 404)
  → Service: validate attendance value
  → Service: if RSVP exists for guest → update; else → create (upsert)
  → Service: record analytics event rsvp_submitted
  → Handler: return confirmation
```

```
GET /events/:id/rsvps/stats
  → Handler: parse event ID
  → Service: verify ownership
  → Service: aggregate rsvps by attendance status + guest_count sums
  → Handler: return { total_guests, yes, no, maybe, pending, attendance_rate }
```

## Validation

- token: required, must match an existing guest token.
- attendance: one of `yes`, `no`, `maybe`.
- guest_count: integer ≥ 1, max 20 (reasonable party size limit).
- message: optional, max 1000 chars.
- Event must be published.

## Authorization

- Public submission: no JWT — authenticated by guest token only.
- Recap/list: JWT + event ownership.
- Export: JWT + ownership + `rsvp.export` capability (Pro).

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT (couple endpoints) |
| 403 | FORBIDDEN | Not event owner |
| 404 | NOT_FOUND | Guest token invalid or event not published |
| 422 | VALIDATION_ERROR | Invalid attendance, guest_count, message |

## Security Considerations

- Guest tokens are bearer credentials — they must be unguessable and handled with care in URLs.
- Rate-limit the public RSVP endpoint to prevent token brute-forcing.
- Do not leak whether a token exists (use generic "invalid or expired link" message).
- RSVP data is couple-private — the public endpoint must not expose other guests' responses.
- `guest_count` validation prevents abuse (huge party sizes).
- Log RSVP submissions at info level with event ID, never the full token.

## Testing Requirements

- Unit tests for upsert logic (first submit creates, second updates).
- Integration tests for public submission flow with valid/invalid tokens.
- Test submission to unpublished event rejected.
- Test recap aggregation correctness.
- Test guest_count limits.
- Test rate limiting on public endpoint.
- Test that other guests' data is not leaked in responses.

## Dependencies

- Guests module — guest records and tokens.
- Events module — published status.
- Analytics module — rsvp_submitted events.
- Export module — Excel export (Pro).

## Related Modules

- **Guests** — Source of guest identity.
- **Events** — Parent entity and published gate.
- **Export** — RSVP recap export.
- **Analytics** — Submission tracking.

## Known Limitations

- No RSVP deadline / closing date per event.
- No RSVP per event schedule (akad vs resepsi separately).
- No validation of guest_count when attendance = `no` (should force count 0 conceptually).
- No reminder automation for non-responders.
- No WebSocket real-time updates for the couple's recap (polling only).

## TODO

- [ ] Add RSVP deadline configuration per event.
- [ ] Add separate attendance for akad and resepsi.
- [ ] Enforce guest_count = 0 when attendance = `no`.
- [ ] Add reminder automation for pending guests.
- [ ] Add real-time recap updates (WebSocket, future).
- [ ] Add per-guest RSVP history tracking.