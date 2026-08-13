# Module: WhatsApp

## Purpose

Facilitate sharing invitations and communicating with guests via WhatsApp. Generates personalized invitation links, provides message templates, and (in the future) supports bulk sending.

## Responsibilities

- Generate per-guest WhatsApp share links using guest phone + access token.
- Provide message templates for invitation announcements and reminders.
- Build the `wa.me` URL with pre-filled message text.
- (Future) Bulk message sending via WhatsApp Business API.
- (Future) Broadcast tracking.

## Non-Responsibilities

- Storing WhatsApp chat history.
- Sending messages server-side (current scope: link generation only, user clicks to send).
- Phone number verification (future).
- WhatsApp Business API integration (future phase).

## Actors

- **User (couple):** Generates links and shares with guests via WhatsApp.
- **Guest (invitee):** Receives the message, opens the invitation, RSVPs.
- **System:** Builds URLs and message text (no direct sending in current scope).

## Business Rules

- WhatsApp link format: `https://wa.me/{phone}?text={url_encoded_message}`.
- The message contains the invitation URL: `https://undangan.example.com/{slug}?token={guest_token}`.
- The guest token personalizes the link so the guest's RSVP is pre-identified.
- If a guest has no phone number, link generation is skipped (user notified).
- Phone numbers are normalized to the `62` format before building `wa.me` links.
- Message templates are configurable per event (system defaults provided).
- Manual mode (current): user clicks the link to open WhatsApp with the pre-filled message.
- Auto/bulk mode (Premium/Pro future): server-side message queue via WhatsApp Business API.
- Bulk sending requires the `whatsapp.bulk` capability (Premium/Pro).
- Message length should stay under WhatsApp limits (plain text, no attachments in current scope).

## Entities

- No dedicated database table in current scope. Uses `guests` (phone, token), `events` (slug, couple_name), and configurable message templates.
- Future: `message_templates`, `broadcast_jobs` tables.

## Database

- Current scope: reads `guests` (phone, token) and `events` (slug, couple_name, title).
- No writes.
- Future tables: `whatsapp_message_templates`, `whatsapp_broadcasts`, `whatsapp_messages`.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/events/:id/whatsapp/message` | JWT | Get the invitation message template |
| PUT | `/api/v1/events/:id/whatsapp/message` | JWT | Update message template |
| GET | `/api/v1/events/:id/guests/:gid/whatsapp-link` | JWT | Get personalized wa.me link for a guest |
| GET | `/api/v1/events/:id/whatsapp/links` | JWT | Get all personalized links (list) |

## Request Flow

```
GET /events/:id/guests/:gid/whatsapp-link
  → Handler: parse event ID + guest ID
  → Service: verify event ownership
  → Service: load guest (phone, token), load event (slug, couple_name)
  → Service: normalize phone to 62 format
  → Service: build message from template (replace {name}, {link} placeholders)
  → Service: build wa.me URL with url-encoded message
  → Handler: return { whatsapp_link, message_text }
```

## Validation

- Guest must have a valid phone number.
- Event must be published or at least have a valid slug (links to draft events are allowed for preview but flagged).
- Message template length: max 1000 characters.
- Placeholders: only known placeholders allowed (`{name}`, `{link}`, `{couple_name}`, `{date}`).

## Authorization

- JWT required; event ownership enforced.
- WhatsApp links are only accessible to the event owner (they are sent to guests by the couple).

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not event owner |
| 404 | NOT_FOUND | Event or guest not found |
| 422 | VALIDATION_ERROR | Guest has no phone, invalid template |

## Security Considerations

- Guest tokens embedded in links are sensitive — the link grants RSVP identity. Do not log full links.
- URL encoding prevents message injection into the wa.me link.
- Message templates are user content — sanitize to prevent link/URL injection.
- Phone numbers are PII — minimize exposure in API responses (masking in list endpoints).
- Bulk sending (future) requires opt-in consent handling per Indonesia regulations (UU PDP).
- Rate-limit link generation to prevent scraping of guest phone numbers.

## Testing Requirements

- Unit tests for phone normalization (08xxx → 62xxx).
- Unit tests for URL building and encoding.
- Unit tests for template placeholder replacement.
- Integration tests for link generation with ownership.
- Test guest without phone returns clear error.
- Test token not exposed in logs.

## Dependencies

- Guests module — phone and token data.
- Events module — slug, couple_name.
- Subscriptions module — `whatsapp.bulk` capability check (future).
- Configuration for the platform URL base (`APP_BASE_URL`).

## Related Modules

- **Guests** — Source of phone/token data.
- **Events** — Invitation URL target.
- **RSVP** — Guests use the link to submit RSVP with token.

## Known Limitations

- Manual mode only — link generation, no actual message delivery.
- No WhatsApp Business API integration yet.
- No delivery/sent tracking.
- No reminder scheduling.
- Message template is a single per-event template (no per-guest customization beyond placeholders).
- No multimedia message support (photos/videos in WhatsApp messages).

## TODO

- [ ] Implement WhatsApp Business API integration for bulk sending (Premium/Pro).
- [ ] Add broadcast job table and queue.
- [ ] Add per-guest message history.
- [ ] Add reminder scheduling (e.g., 1 day before event).
- [ ] Add opt-out handling.
- [ ] Add analytics for link clicks (via Analytics module).