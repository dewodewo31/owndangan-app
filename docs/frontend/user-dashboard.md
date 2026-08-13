# User Dashboard

## Overview

The user dashboard (`/user/*`) is the primary workspace for the wedding couple (pengantin). From here, they create and manage their digital invitation, add guests, track RSVPs, and handle billing. The user layout provides a sidebar with navigation links and requires JWT authentication with `role === 'user'`.

## Dashboard Overview (`/user`)

The main user page shows a summary of the user's events and subscription status.

**When the user has no events:**
- A welcome card with a "Buat Undangan Baru" (Create Invitation) button.
- Brief instructions for getting started.
- Link to the pricing page if they haven't subscribed.

**When the user has events:**
- Event cards showing: couple name, slug, status (draft/published), view count, and a "Buka" (Open) link.
- Quick stats: total guests, RSVP received, guestbook messages.
- Subscription info: current plan, expiry date, days remaining.
- Low-urgency alerts (e.g., "Undangan Anda belum dipublikasikan" or "Masa aktif akan berakhir dalam 7 hari").

## Editor (`/user/editor`)

The editor is the core feature of the user dashboard. It is a multi-step form that covers all 9 invitation sections. The editor is a client component that fetches the current event data and allows the user to modify each section.

**Sections:**

| # | Section | Fields |
|---|---------|--------|
| 1 | Cover | Couple names, wedding date, "Buka Undangan" button text |
| 2 | Opening | Greeting text, Quran verse, optional |
| 3 | Couple Profile | Groom name & parents, bride name & parents |
| 4 | Event Details | Akad date/time/location/map, Resepsi date/time/location/map |
| 5 | Gallery | Upload photos, reorder, captions, max count depends on plan |
| 6 | RSVP | Toggle on/off, message options |
| 7 | Guestbook | Toggle on/off, moderation setting |
| 8 | Digital Gifts | Bank accounts, e-wallet, QRIS upload, customizable message |
| 9 | Music | Select preset or upload MP3, toggle on/off |

**Editor UX:**
- Tab-based navigation at the top (or sidebar) to switch between sections.
- "Save" button per section (auto-save on blur is considered for future).
- Unsaved changes prompt: "Anda memiliki perubahan yang belum disimpan."
- Plan-based feature gating: if a section is not available in the user's plan (e.g., custom music upload), the UI shows a locked state with an "Upgrade" link.

**Publish:**
- A "Publikasikan" button in the editor header.
- On publish: `POST /api/v1/events/:id/publish` → slug becomes publicly accessible.
- On unpublish: `POST /api/v1/events/:id/unpublish` → `/[slug]` returns 404.

## Guest Management (`/user/guests`)

Guest list page with full CRUD and WhatsApp integration.

**Features:**
- Table: name, phone, category (keluarga/teman/rekan kerja), RSVP status (pending/hadir/tidak hadir), actions.
- Add guest form (modal or inline): name (required), phone (optional, `+62` format), category (optional).
- CSV import: upload CSV with columns `name, phone, category`. Backend validates and returns import summary (X added, Y duplicates skipped).
- Bulk actions: select multiple guests, send WhatsApp invitations.
- WhatsApp link: per guest, a "Kirim WA" button that opens `https://wa.me/{phone}?text={encoded_template}`. The template message includes the invitation URL and couple names.
- Search and filter by name, category, RSVP status.
- Pagination (20 per page).

## RSVP Recap (`/user/rsvp`)

RSVP page shows the attendance summary for the wedding.

**Summary cards:**
- Total invited guests count.
- RSVP received count (with percentage).
- Will attend (count).
- Will not attend (count).
- No response yet (count).

**Guestbook moderation:**
- Pending approval messages list.
- Approve / reject / delete actions.
- Approved messages appear on the public invitation page.

**Export:**
- "Export Excel" button: downloads guest list with RSVP status.
- Backend generates the XLSX file; the frontend triggers a download via the API response.

## Billing (`/user/billing`)

Billing page shows the user's current subscription and available plans to purchase or upgrade.

**Current subscription:**
- Plan name (Basic / Premium / Pro).
- Status (active / expired / pending).
- Expiry date.
- Days remaining.
- Feature list with checkmarks.

**Plan comparison:**
- Three plan cards (Basic, Premium, Pro) in a side-by-side layout.
- Each card shows: price, duration, feature highlights.
- Current plan is marked with a "Paket Saat Ini" badge.
- "Upgrade" button for plans above the current tier.

**Upgrade flow:**
1. User clicks "Upgrade" or "Beli Paket".
2. `POST /api/v1/payments/snap` is called with the selected plan.
3. Backend returns a Snap token.
4. Midtrans Snap popup opens (via `window.snap.pay(token)`).
5. On payment success (Snap callback), the frontend shows a success message and refreshes subscription data.
6. The backend webhook handles the actual subscription activation.

**Transaction history:**
- Table of past transactions: order ID, plan, amount, payment method, status, date.
- Status badges: settlement (green), pending (yellow), expire (red), cancel (gray).

## Related API Documentation

See `docs/api/events.md`, `docs/api/guests.md`, `docs/api/rsvp.md`, `docs/api/payments.md`, and `docs/api/subscriptions.md`.