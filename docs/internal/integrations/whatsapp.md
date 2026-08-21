# WhatsApp Integration

## Current Implementation

The platform uses **wa.me links** for WhatsApp messaging. This is a deep link that opens WhatsApp with a pre-filled message. No API key or third-party service is required for the current implementation.

## wa.me Link Generation

### Format

```
https://wa.me/{phone_number}?text={url_encoded_message}
```

### Implementation

```go
func GenerateWhatsAppLink(phone string, message string) string {
    // Normalize phone: remove +, spaces, dashes
    phone = strings.TrimPrefix(phone, "+")
    phone = strings.ReplaceAll(phone, " ", "")
    phone = strings.ReplaceAll(phone, "-", "")

    // Encode message
    encoded := url.QueryEscape(message)

    return fmt.Sprintf("https://wa.me/%s?text=%s", phone, encoded)
}
```

### Phone Number Format
- Must be in international format (E.164).
- Examples: `6281234567890` (Indonesia), `+6281234567890` (input format).
- The code strips the leading `+` for wa.me links.
- Do not add `+` or `00` prefix in the final URL.

## Message Templates

### Template Variables

| Variable | Description | Example |
|---|---|---|
| `{guest_name}` | Guest's full name | `Budi` |
| `{couple_names}` | Bride and groom names | `Ani & Budi` |
| {event_date}` | Event date (formatted) | `15 Januari 2025` |
| `{event_time}` | Event time | `10:00 WIB` |
| `{event_location}` | Event venue | `Hotel Indonesia` |
| `{rsvp_link}` | RSVP URL (shortened) | `https://undangan.io/rsvp/abc` |
| `{invitation_link}` | Invitation page URL | `https://undangan.io/inv/xyz` |

### Default Templates

**Invitation Message:**
```
Halo {guest_name}!
Kami mengundang Anda untuk hadir di acara pernikahan kami:
{event_date} pukul {event_time}
di {event_location}
Konfirmasi kehadiran: {rsvp_link}
— {couple_names}
```

**RSVP Reminder:**
```
Halo {guest_name}!
Kami masih menunggu konfirmasi kehadiran Anda untuk acara pernikahan {couple_names}.
Mohon konfirmasi sebelum {rsvp_deadline} di:
{rsvp_link}
Terima kasih!
```

### Customization
- Admins can edit message templates from the dashboard.
- Templates use Go `text/template` syntax (HTML escaping is not needed for plain text).
- Template preview is available before sending.

## Personalized Guest Messages

### Sending Messages

```go
func (s *WhatsAppService) SendInvitation(ctx context.Context, guest *Guest, event *Event) (string, error) {
    // 1. Build personalized message
    message := s.renderTemplate(event.InviteTemplate, map[string]string{
        "guest_name":       guest.Name,
        "couple_names":     fmt.Sprintf("%s & %s", event.GroomName, event.BrideName),
        "event_date":       event.Date.Format("2 Januari 2006"),
        "event_time":       event.Time,
        "event_location":   event.Location,
        "rsvp_link":        fmt.Sprintf("https://undangan.io/rsvp/%s", guest.ID),
        "invitation_link":  fmt.Sprintf("https://undangan.io/inv/%s", event.ID),
    })

    // 2. Generate wa.me link
    link := GenerateWhatsAppLink(guest.Phone, message)

    // 3. Log the send attempt
    s.logger.Info("whatsapp link generated",
        "guest_id", guest.ID,
        "event_id", event.ID,
    )

    // 4. Return link for frontend to open
    return link, nil
}
```

### Delivery Method
- The backend generates the `wa.me` link.
- The frontend opens the link in a new tab.
- The user must click "Send" in WhatsApp manually.
- No server-side delivery (no WhatsApp Business API — yet).

## Rate Limits

### wa.me Links (Current)
- No rate limits from the API (no API used).
- Recommend: limit link generation to 50 requests per minute per user to prevent abuse.
- Reasonable use: generating 1 link per guest is normal; bulk generation should be throttled.

### WhatsApp Business API (Future)
If WhatsApp Business API is integrated:
- **Free tier**: 1,000 messages per day, 10 messages per second.
- **Paid tier**: Higher limits based on business account quality.
- Template messages must be pre-approved by Meta.
- Opt-out handling is mandatory.

## Bulk Sending (Planned)

### Future Architecture
```
Admin uploads CSV / selects guests
  → Backend validates phone numbers
  → Backend queues messages (Redis/Kafka)
  → Worker processes queue (rate-limited)
  → For each guest: send via WhatsApp Business API
  → Status tracked per recipient
```

### When Not to Use wa.me Links
- More than 50 recipients.
- Scheduled/delayed sending.
- Delivery tracking required.
- Opt-out management needed.

## Best Practices

1. **Always include opt-out instructions** in promotional messages (not applicable for transactional invitations).
2. **Respect Indonesian hours**: 08:00–20:00 WIB only (avoid sending late at night).
3. **Phone validation**: Validate phone numbers with Indonesia country code before generating links.
4. **Logging**: Log every generated link (guest ID, timestamp, IP) for abuse monitoring.
5. **No HTML in messages**: WhatsApp messages are plain text only.
6. **URL shortening**: Use a URL shortener for RSVP links to avoid broken messages (URLs in WhatsApp are auto-previewed).