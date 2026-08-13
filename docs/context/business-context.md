# Business Context

## Core Business Flow

```
User
 ↓
Register / Login
 ↓
Choose Package
 ↓
Create Payment (Midtrans Snap)
 ↓
Payment Settlement (Webhook)
 ↓
Subscription Activated
 ↓
Create Invitation (Event)
 ↓
Configure 9 Sections (Editor)
 ↓
Publish Invitation
 ↓
Share via WhatsApp / Link
 ↓
Guest Opens Invitation
 ↓
Guest RSVP
 ↓
Guest Sends Messages (Guestbook)
 ↓
Couple Views RSVP Recap
```

## Package Rules

| Feature | Basic | Premium | Pro |
|---------|-------|---------|-----|
| Price (Rp) | 49K–99K | 149K–199K | 299K–499K |
| Duration | 3 months | 1 year | Lifetime / Custom |
| Themes | 5 standard | 20+ premium | All + custom request |
| Guest Limit | 100 | 500 | Unlimited |
| Custom Slug | Yes | Yes | Yes + Custom Domain |
| Background Music | Preset only | Upload MP3 | Upload MP3 |
| Photos | Max 5 | Max 20 | Unlimited |
| Video | No | YouTube embed | YouTube + HD |
| RSVP | Yes | Yes | Yes + Excel export |
| Digital Gifts | Bank only | Bank + e-Wallet | QRIS + Direct Transfer |
| WhatsApp | Manual link | Auto broadcast | API bulk sender |
| QR Guestbook | No | No | Yes |
| Watermark Removal | No | Yes | Yes |

## Subscription Rules

- A user can have at most one active subscription at a time.
- Subscription starts after Midtrans settlement webhook is received and verified.
- Subscription expires based on plan_type duration:
  - Basic: 90 days from activation
  - Premium: 365 days from activation
  - Pro: lifetime (no expiry)
- Expired subscriptions revert features to free tier.
- Pro subscription has no expiry (stored as NULL expiry date).
- Admin can manually create, extend, or terminate any subscription.

## Expiration Rules

- Subscription expires at `expires_at` timestamp.
- 7 days before expiry, system sends warning (TODO: implement notification).
- On expiry date, subscription status changes to `expired`.
- Expired subscription: invitation stays published but premium features are locked.
- Guest access to invitation is NOT affected by subscription expiry.
- Editing disabled after expiry until subscription is renewed.

## Guest Limits

- Guests count is the total number of guest records associated with an event.
- If guest limit is reached, adding new guests is blocked.
- Deleting guests frees up quota.
- Admin bypass for Pro plan (unlimited limit).
- CSV import respects limit (reject if import would exceed limit).

## Feature Entitlement

Use capability-based entitlement, not hardcoded plan names.

```
guest.max
theme.max
gallery.photo.max
video.enabled
music.upload.enabled
custom_domain.enabled
watermark.removed
whatsapp.bulk.enabled
guestbook.qr.enabled
rsvp.export.enabled
digital_gift.qris.enabled
template.custom_request.enabled
```

Each package maps to a set of capabilities stored in the packages table.

## Template Entitlement

- Basic: access to Standard theme group (5 themes).
- Premium: access to Premium theme group (20+ themes).
- Pro: access to all themes + option to request custom.
- Template assignment is per event, not per subscription.

## RSVP Rules

- Guest records have an associated RSVP.
- RSVP has attendance status: `yes`, `no`, `maybe`.
- Guest can submit RSVP only once. Subsequent submissions update the existing RSVP.
- RSVP includes guest_count (number of people attending).
- RSVP submission is anonymous (no auth required), identified by guest token or ID.
- Couple can view RSVP recap in user dashboard.

## Gift Rules

- Digital gift info is static content displayed on invitation.
- Bank accounts, e-wallet info, and QRIS are configured by the couple.
- Platform does NOT process digital gift payments.
- Gift messages from guests are stored as guestbook messages.

## Invitation Publishing Rules

- Invitation starts as `draft`.
- User must complete required sections before publishing.
- Published invitations are publicly accessible via `/[slug]`.
- Unpublishing an invitation makes it inaccessible (returns 404).
- Publish/unpublish is logged in audit log.

## Slug Rules

- Slug must be 3–100 characters.
- Allowed characters: lowercase letters, numbers, hyphens.
- Slug must be globally unique across all events.
- Slug can be changed as long as it remains unique.
- Changing slug invalidates the old URL (no automatic redirect).
- Pro users can use custom domain instead of slug.

## Admin Override Rules

- Admin can create subscription for any user.
- Admin can extend/terminate subscriptions.
- Admin can suspend/activate user accounts.
- Admin can change user role.
- Admin can view any invitation (read-only).
- Admin can moderate guestbook messages.
- All admin actions are logged in audit log.

## Business Invariants

- A user cannot exceed the guest limit of their active plan.
- A subscription must be active before premium features can be used.
- A slug must be globally unique across all events.
- Payment status must never be activated solely based on frontend response.
- Midtrans webhook is authoritative for payment settlement.
- Users cannot access another user's events, guests, or RSVP data.
- Guest RSVP is limited to one response per guest record.
- Subscription activation is irreversible for the same transaction (idempotent).
- Guestbook messages can be moderated before public display.
- Free tier (no subscription) has the same limits as the Basic plan.
