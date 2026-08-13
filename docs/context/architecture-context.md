# Architecture Context

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Users                                │
│  (Couples, Guests, Admins via browser/mobile)               │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTTPS
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    CDN / Reverse Proxy                       │
│                   (Cloudflare / Nginx)                       │
└──────────┬────────────────────────────────────┬─────────────┘
           │                                    │
           ▼                                    ▼
┌─────────────────────┐          ┌──────────────────────────────┐
│   Next.js Frontend  │          │   Go Backend (REST API)      │
│   (App Router)      │◄────────►│   Handler → Service → Repo   │
│                     │   JSON   │                              │
├─────────────────────┤          ├──────────────────────────────┤
│ (public)/ Landing   │          │ /api/v1/auth/*               │
│ admin/ Dashboard    │          │ /api/v1/users/*              │
│ user/ Dashboard     │          │ /api/v1/events/*             │
│ [slug]/ Invitation  │          │ /api/v1/guests/*             │
└─────────────────────┘          │ /api/v1/payments/*           │
           │                     │ /api/v1/admin/*              │
           │                     │ /api/v1/webhook/midtrans     │
           │                     └──────────────┬───────────────┘
           │                                    │
           │                    Snap.js          │ GORM
           │                    Client Key       │
           ▼                                    ▼
┌─────────────────────┐          ┌──────────────────────────────┐
│   Midtrans Gateway  │          │   PostgreSQL Database         │
│   (Snap API)        │◄────────►│                              │
│                     │ Webhook  │ users, subscriptions,         │
│                     │          │ transactions, events,         │
│                     │          │ guests, rsvps, guestbook...   │
└─────────────────────┘          └──────────────────────────────┘
```

## Decoupled Architecture

Frontend and backend are completely separate:
- No server-side rendering coupling between Go and Next.js.
- Next.js can be deployed independently (Vercel or Node server).
- Go backend can be deployed independently (binary or container).
- Communication exclusively via REST JSON API.

## Route Design

### Frontend Routes

| Route Group | Path | Purpose |
|-------------|------|---------|
| Public | `/` | Landing page |
| Public | `/pricing` | Pricing page |
| Public | `/[slug]` | Public invitation (dynamic) |
| User | `/user` | User dashboard overview |
| User | `/user/editor` | Invitation editor |
| User | `/user/guests` | Guest management |
| User | `/user/rsvp` | RSVP recap & guestbook |
| User | `/user/billing` | Package purchase & history |
| Admin | `/admin` | Admin dashboard overview |
| Admin | `/admin/users` | User management |
| Admin | `/admin/transactions` | Transaction management |
| Admin | `/admin/packages` | Package management |
| Admin | `/admin/templates` | Template management |

### API Routes

| Prefix | Purpose | Auth |
|--------|---------|------|
| `/api/v1/auth/*` | Authentication | None / Public |
| `/api/v1/users/*` | User profile | JWT |
| `/api/v1/events/*` | Invitation management | JWT |
| `/api/v1/guests/*` | Guest management | JWT |
| `/api/v1/rsvps/*` | RSVP submission | Public or JWT |
| `/api/v1/guestbook/*` | Guestbook messages | Public or JWT |
| `/api/v1/digital-gifts/*` | Digital gift info | JWT |
| `/api/v1/subscriptions/*` | Subscription management | JWT |
| `/api/v1/packages/*` | Package listing | Public + JWT |
| `/api/v1/payments/*` | Payment transactions | JWT |
| `/api/v1/webhook/midtrans` | Midtrans notification | Webhook Auth |
| `/api/v1/admin/*` | Admin operations | JWT (admin) |
| `/api/v1/public/*` | Public invitation data | None |

## Module Boundaries

```
Authentication
    └── Login / Register / Token / Session

Users
    └── User profile / Account settings

Subscriptions
    └── User subscription lifecycle (activate, extend, expire)

Packages
    └── Package definitions / Feature limits / Pricing

Payments
    └── Midtrans transaction lifecycle / Snap token / Webhook

Events
    └── Wedding invitation core entity / Slug / Status / Sections

Invitation Editor
    └── 9-section invitation content configuration

Templates
    └── Invitation themes / Template versions / Assets

Guests
    └── Guest CRUD / CSV import / Category / Limit check

WhatsApp
    └── WhatsApp link generation / Message templates / Bulk sending (future)

RSVP
    └── Attendance response / Count / Status updates

Guestbook
    └── Guest messages / Moderation / Display

Digital Gifts
    └── Bank accounts / e-Wallet / QRIS / Gift messages

Gallery
    └── Photo upload / Video embed / Album management

Music
    └── Background music / Upload / Preset tracks

SEO
    └── Metadata / OpenGraph / Slug validation / Sitemap

Analytics
    └── Platform metrics / Invitation views / RSVP stats

Audit Log
    └── Admin and user action logging

Export
    └── RSVP Excel export / Guest list export
```

## Layer Separation (Backend)

```
HTTP Request
    │
    ▼
┌──────────────┐
│  Middleware   │  ← Auth, logger, CORS, rate limiter
└──────┬───────┘
       ▼
┌──────────────┐
│   Handler    │  ← Parse request, validate transport, call service, respond
└──────┬───────┘
       ▼
┌──────────────┐
│   Service    │  ← Business logic, business validation, orchestration
└──────┬───────┘
       ▼
┌──────────────┐
│  Repository  │  ← Database access, queries, GORM operations
└──────┬───────┘
       ▼
┌──────────────┐
│  PostgreSQL  │
└──────────────┘
```

## State Management

- Backend is stateless (JWT contains session info).
- Database is the single source of truth.
- Frontend can cache API responses (SWR/TanStack Query).
- Midtrans state is managed by Midtrans; backend mirrors relevant state.

## Security Boundaries

- Public API: rate-limited, non-sensitive data only.
- User API: JWT required, ownership verified.
- Admin API: JWT required, admin role verified.
- Webhook: source IP verified + signature verified.
- File upload: validated type/size, served via signed URLs.
