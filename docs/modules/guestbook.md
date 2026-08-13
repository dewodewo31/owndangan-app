# Module: Guestbook

## Purpose

Enable guests to leave messages, wishes, and prayers on the couple's invitation page. Provides moderation capabilities for the couple to approve or delete messages before public display.

## Responsibilities

- Accept guestbook message submissions from guests (public endpoint).
- Store messages with sender name, message text, and timestamp.
- Provide moderation workflow: messages are hidden by default until approved.
- Allow the couple to approve, reject, or delete messages.
- Display approved messages on the public invitation page.
- Provide guestbook management for the couple's dashboard.

## Non-Responsibilities

- RSVP submission (handled by RSVP module — the RSVP message is separate from guestbook).
- Digital gift messages (handled by Digital Gifts module).
- WhatsApp notifications for new messages (future).
- Spam detection (future).

## Actors

- **Guest (invitee):** Submits a guestbook message via the public invitation page.
- **User (couple):** Moderates messages (approve, reject, delete).
- **Admin:** Can moderate across all events.
- **Guest (public):** Reads approved messages on the invitation page.

## Business Rules

- Messages are submitted with a name and message text (no auth required).
- Submitted messages are hidden by default (`is_approved = false`).
- Only approved messages are displayed on the public invitation page.
- The couple can approve or delete messages from their dashboard.
- The couple can choose to auto-approve messages (configurable per event via event_sections).
- Admin can moderate (approve/delete) messages for any event.
- Deleted messages are hard-deleted (not recoverable).
- There is no edit capability for messages (submit once, submit again for a new message).
- A single guest can submit multiple messages (no uniqueness constraint).
- Message names are displayed as-is; the couple can moderate offensive names.
- Spam or abusive messages can be reported but no automated spam detection in current scope.

## Entities

- **GuestbookMessage:** `{ id, event_id, name, message, is_approved, created_at, updated_at }`

## Database

- Table: `guestbook_messages`
- Hard delete (no `deleted_at`).
- Index on `event_id`, `is_approved`.
- No soft delete — moderation is via `is_approved` flag.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/public/guestbook` | Public | Submit guestbook message |
| GET | `/api/v1/public/guestbook/:slug` | Public | List approved messages |
| GET | `/api/v1/events/:id/guestbook` | JWT | List all messages (with moderation) |
| PUT | `/api/v1/events/:id/guestbook/:mid/approve` | JWT | Approve message |
| DELETE | `/api/v1/events/:id/guestbook/:mid` | JWT | Delete message |

## Request Flow

```
POST /public/guestbook
  → Handler: parse body (slug, name, message)
  → Service: resolve event by slug, verify published
  → Service: sanitize name and message (strip HTML)
  → Service: check auto-approve setting (event_sections.guestbook_auto_approve)
  → Service: create message with is_approved based on setting
  → Handler: return created message (201)
```

```
GET /public/guestbook/:slug
  → Handler: parse slug
  → Service: resolve event, verify published
  → Service: load is_approved = true messages, order by created_at desc
  → Handler: return message list (name, message, created_at, limited to 50)
```

## Validation

- name: required, max 255 chars, trimmed.
- message: required, max 2000 chars, trimmed.
- slug: must match a published event.
- HTML tags in name/message are stripped (not allowed).
- Rate limit: 5 submissions per IP per minute.

## Authorization

- Public submission: no auth (rate-limited).
- Public reading: no auth.
- Moderation: JWT + event ownership (couple) or admin role.
- Admin: can moderate any event's guestbook.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 404 | NOT_FOUND | Event slug not found or not published |
| 422 | VALIDATION_ERROR | Name/message empty or too long |
| 429 | RATE_LIMITED | Too many submissions |

## Security Considerations

- HTML stripping is essential to prevent stored XSS — guestbook messages are displayed on the public page.
- Name and message must be sanitized on both save and render (defense in depth).
- Rate-limit the public submission endpoint to prevent spam.
- Auto-approve is a convenience feature but increases spam risk — recommended off by default.
- Moderated (deleted) messages are hard-deleted — no recovery.
- No personal data beyond name/message is stored (no IP in message, but IP may be logged for rate limiting).
- Message content should be scanned for malicious links (future).

## Testing Requirements

- Unit tests for message sanitization (HTML stripping).
- Unit tests for auto-approve logic.
- Integration tests for public submission and listing.
- Integration tests for moderation (approve, delete).
- Test rate limiting behavior.
- Test that unapproved messages are not returned in public listing.
- Test HTML injection attempts are stripped.

## Dependencies

- Events module — event slug and published status.
- Invitation Editor module — auto-approve setting in event_sections.
- Audit Log module — moderation actions.

## Related Modules

- **Events** — Parent entity, slug resolution.
- **Invitation Editor** — Auto-approve toggle.
- **RSVP** — Separate communication channel; guestbook is distinct from RSVP.
- **Admin** — Cross-event moderation.

## Known Limitations

- No spam detection or CAPTCHA (future).
- No edit of messages (delete and resubmit).
- No notify-couple-on-new-message feature.
- No pagination on public listing (limited to 50).
- No profanity filter.
- No message reactions (like/heart).

## TODO

- [ ] Add CAPTCHA or honeypot for spam prevention.
- [ ] Add email/WhatsApp notification to couple on new message.
- [ ] Add pagination for public guestbook listing.
- [ ] Add profanity filter (optional).
- [ ] Add message reactions (like/heart emoji).
- [ ] Add auto-delete after event date (configurable retention).