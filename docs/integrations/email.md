# Email Integration

## Overview

The platform sends transactional emails for user onboarding, payment confirmations, and subscription lifecycle events. Email is **not** used for marketing.

## Mail Service Provider

### Current: SMTP (Direct)

| Variable | Description | Example |
|---|---|---|
| `SMTP_HOST` | SMTP server hostname | `smtp.sendgrid.net` |
| `SMTP_PORT` | SMTP port (TLS) | `587` |
| `SMTP_USERNAME` | SMTP username | `apikey` |
| `SMTP_PASSWORD` | SMTP password or API key | `SG.abc123...` |
| `SMTP_FROM` | Sender email address | `noreply@owndangan.com` |
| `SMTP_FROM_NAME` | Sender display name | `Owndangan` |

### Future: Dedicated Email Service Provider

When migrating to a dedicated service (SendGrid, Mailgun, Amazon SES):

| Provider | Pros | Cons |
|---|---|---|
| SendGrid | Good deliverability, templates, analytics | Cost at scale |
| Mailgun | API-first, good for developers | Fewer templates |
| Amazon SES | Cheapest at scale | Higher setup complexity |
| Resend | Modern API, good DX | Newer, fewer features |

## Transactional Emails

### Welcome Email

**Trigger**: User registration.

**Content**:
- Subject: `Selamat datang di Owndangan, {name}!`
- Body: Account created, quick start guide, link to create first event.
- Priority: High (send immediately).

### Payment Confirmation

**Trigger**: Midtrans webhook with `settlement` status.

**Content**:
- Subject: `Pembayaran berhasil — {plan_name} aktif!`
- Body: Plan name, amount paid, activation date, expiry date, invoice link.
- Priority: High (send immediately).
- Attachment: Invoice PDF (future).

### Payment Failed

**Trigger**: Midtrans webhook with `deny`, `cancel`, or `expire` status.

**Content**:
- Subject: `Pembayaran {plan_name} gagal`
- Body: Reason, retry link, support contact.
- Priority: High (send immediately).

### Subscription Expiry Reminder

**Trigger**: 7 days before subscription expiry.

**Content**:
- Subject: `Langganan Owndangan akan berakhir`
- Body: Current plan, expiry date, renewal link.
- Priority: Medium.

### Subscription Expired

**Trigger**: Subscription expiry date reached.

**Content**:
- Subject: `Langganan Owndangan telah berakhir`
- Body: What features are now locked, renewal link.
- Priority: Medium.

### Event Invitation

**Trigger**: User sends invitation to a guest (currently via WhatsApp, email planned).

**Content**:
- Subject: `Undangan Pernikahan {couple_names}`
- Body: Event details, RSVP link, personalized greeting.
- Priority: Low (batch send).

## Implementation

### Sending Emails

```go
func (s *EmailService) Send(ctx context.Context, to string, subject string, body string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", s.fromAddress)
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/html", body)

    dialer := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPass)
    return dialer.DialAndSend(msg)
}
```

### Queue Implementation

Emails are sent asynchronously via a work queue to avoid blocking API responses.

```go
// Publisher (in handler/service)
func (s *EmailService) QueueWelcomeEmail(ctx context.Context, user *User) {
    s.queue.Publish(ctx, "email:welcome", EmailPayload{
        To:      user.Email,
        Name:    user.Name,
        UserID:  user.ID,
    })
}

// Worker (background goroutine)
func (w *EmailWorker) ProcessWelcomeEmail(ctx context.Context, payload EmailPayload) error {
    subject := fmt.Sprintf("Selamat datang di Owndangan, %s!", payload.Name)
    body := w.renderWelcomeTemplate(payload.Name)
    return w.send(ctx, payload.To, subject, body)
}
```

### Retry Logic
- First attempt: immediate.
- Retry 1: 1 minute.
- Retry 2: 5 minutes.
- Retry 3: 30 minutes.
- After 3 failures: log alert, move to dead-letter queue, notify admin.

## Email Templates (Future)

### Planned Architecture
- Templates stored in the database or file system.
- Rendered server-side with Go `html/template`.
- Plain text alternative generated from HTML.
- Template preview available in admin dashboard.

### Template Variables

| Variable | Available In |
|---|---|
| `{{.Name}}` | All emails |
| `{{.PlanName}}` | Payment, subscription emails |
| `{{.Amount}}` | Payment emails |
| `{{.ExpiryDate}}` | Subscription emails |
| `{{.InvoiceURL}}` | Payment confirmation |
| `{{.LoginURL}}` | Welcome email |
| `{{.SupportEmail}}` | All emails |

### Template Example

```html
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif;">
  <h1>Halo {{.Name}}!</h1>
  <p>Pembayaran Anda untuk paket <strong>{{.PlanName}}</strong> sebesar <strong>Rp{{.Amount}}</strong> telah berhasil.</p>
  <p>Langganan Anda aktif hingga {{.ExpiryDate}}.</p>
  <a href="{{.InvoiceURL}}" style="...">Lihat Invoice</a>
</body>
</html>
```

## Best Practices

1. **Always include unsubscribe link** (even for transactional emails, per Indonesian regulation).
2. **Use plain text alternative** for better deliverability and accessibility.
3. **Avoid image-heavy emails** — many email clients block images by default.
4. **SPF, DKIM, DMARC** — Configure all three for the sending domain to prevent spoofing.
5. **Bounce handling** — Monitor hard bounces and mark email as invalid.
6. **Rate limiting** — Max 10 emails per second (SMTP connection limit). Queue handles this.
7. **No sensitive data in subject lines** — Never include passwords, reset tokens, or payment amounts in subject.
8. **Test emails** — Send test emails to multiple clients (Gmail, Outlook, Yahoo) before production deployment.