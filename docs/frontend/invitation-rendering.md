# Invitation Rendering

## Overview

The public invitation page (`/[slug]`) is the most important page for SEO and user experience. It is a **server component** that fetches all invitation data at request time, generates dynamic metadata, and renders section-based content. No client-side data fetching is needed for the initial render.

## Data Flow

```
Guest visits /andi-sinta
  → app/[slug]/page.tsx (server component)
  → generateMetadata({ params }) runs first
  → Page component calls fetch(BASE_URL + /public/events/andi-sinta)
  → API returns event + template + sections + gallery + guestbook + gifts
  → Page renders sections conditionally based on sections flags
  → Template CSS config is applied as CSS custom properties
```

## Server Component

```tsx
// app/[slug]/page.tsx
export default async function InvitationPage({ params }: { params: { slug: string } }) {
  const data = await getInvitation(params.slug);
  return (
    <div
      className="invitation-root"
      style={{
        '--primary-color': data.template.css_config.primary_color,
        '--secondary-color': data.template.css_config.secondary_color,
        '--font-family': data.template.css_config.font_family,
      } as React.CSSProperties}
    >
      <MusicPlayer music={data.sections.music} />
      {data.sections.hero_enabled && <HeroSection event={data.event} />}
      {data.sections.couple_enabled && <CoupleSection event={data.event} />}
      <EventDetailsSection event={data.event} />
      <GallerySection images={data.gallery} />
      <RsvpSection eventId={data.event.id} />
      <GuestbookSection messages={data.guestbook} eventId={data.event.id} />
      <DigitalGiftsSection gifts={data.digital_gifts} />
    </div>
  );
}
```

## Section-based Rendering

The `sections` object from the API controls which components are rendered. Each boolean flag maps to a section component. This allows the user to toggle sections on/off without changing the code.

| API Flag | Component | Condition |
|----------|-----------|-----------|
| `hero_enabled` | `HeroSection` | Always shown if enabled |
| `couple_enabled` | `CoupleSection` | Groom & bride profiles |
| `event_details_enabled` | `EventDetailsSection` | Akad & Resepsi info |
| `gallery_enabled` | `GallerySection` | Photo grid (empty if no photos) |
| `video_enabled` | `VideoSection` | YouTube embed |
| `rsvp_enabled` | `RsvpSection` | Interactive form (client component) |
| `guestbook_enabled` | `GuestbookSection` | Messages + form (client component) |
| `digital_gifts_enabled` | `DigitalGiftsSection` | Bank/QRIS info |

## Theme-based Styling

The `template.css_config` object controls the visual appearance of the invitation. Values are applied as CSS custom properties on the root `<div>`:

```css
.invitation-root {
  --primary: #d4af37;
  --secondary: #ffffff;
  --font-heading: 'Playfair Display', serif;
  --font-body: 'Inter', sans-serif;
  --bg-pattern: url('/patterns/floral.png');
}
```

Each template component uses these CSS custom properties for consistent theming. The `template.name` determines which component set is loaded, but the CSS variables provide fine-grained styling control.

## Interactive Components (Client)

Three sections require client-side interactivity and are wrapped in `'use client'`:

1. **RsvpSection** — form to submit attendance. After submission, shows a confirmation message. Uses `POST /api/v1/public/rsvps`.
2. **GuestbookSection** — message list with a form to add new messages. Uses `POST /api/v1/public/guestbook`.
3. **MusicPlayer** — floating play/pause button for background music. Auto-plays on user interaction (first click on page).

These components are isolated islands of interactivity within the server-rendered page. They do not require any dashboard context or auth.

## Loading & Error States

The server component handles loading via Suspense boundaries and error via `notFound()`:

- **Loading:** The `[slug]` segment can have a `loading.tsx` file showing a skeleton of the invitation layout (couple names, date, placeholder images).
- **Not found:** If the slug does not exist or the event is unpublished, `notFound()` is called, which renders `app/not-found.tsx`.
- **Error:** If the API is unreachable, the page throws an error and `app/error.tsx` is rendered with a "Gagal memuat undangan. Silakan coba lagi." message and a retry button.

## Related API Documentation

See `docs/api/public-invitation.md` for the full API contract.