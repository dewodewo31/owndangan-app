# Module Index

## Purpose

This directory contains documentation for every domain module in the Platform Undangan Pernikahan Digital & Cetak. Each module doc describes the module's purpose, responsibilities, business rules, API surface, dependencies, and security considerations. These documents serve as the primary reference for AI coding agents and developers implementing or modifying functionality.

## Module List

| Module | File | Purpose |
|--------|------|---------|
| **Authentication** | [authentication.md](./authentication.md) | Register, login, JWT access tokens, refresh token rotation, session management. |
| **Users** | [users.md](./users.md) | User profile CRUD, account settings, account status management. |
| **Admin** | [admin.md](./admin.md) | Platform administration: user management, package management, template management, dashboard analytics. |
| **Subscriptions** | [subscriptions.md](./subscriptions.md) | Subscription lifecycle: activation, expiry, extension, cancellation, feature entitlement checks. |
| **Packages** | [packages.md](./packages.md) | Package plan definitions, feature flags, pricing tiers, duration configuration. |
| **Payments** | [payments.md](./payments.md) | Midtrans Snap integration, transaction lifecycle, webhook handling, signature verification. |
| **Events** | [events.md](./events.md) | Wedding invitation core entity: CRUD, publish/unpublish, slug management, status workflow. |
| **Invitation Editor** | [invitation-editor.md](./invitation-editor.md) | Nine-section invitation content editor, section toggles, content configuration, preview. |
| **Templates** | [templates.md](./templates.md) | Invitation themes, template groups (standard/premium/all), versioning, asset management. |
| **Guests** | [guests.md](./guests.md) | Guest CRUD, category management, CSV import, guest limit enforcement, access token generation. |
| **WhatsApp** | [whatsapp.md](./whatsapp.md) | WhatsApp link generation, message templates, individual and bulk send (future). |
| **RSVP** | [rsvp.md](./rsvp.md) | Attendance submission, status management, guest count tracking, recap aggregation. |
| **Guestbook** | [guestbook.md](./guestbook.md) | Guest messages, admin/couple moderation, approval workflow, public display. |
| **Digital Gifts** | [digital-gifts.md](./digital-gifts.md) | Bank account info, e-wallet configuration, QRIS image upload, gift message display. |
| **Gallery** | [gallery.md](./gallery.md) | Photo upload, video embed (YouTube), album sorting, storage integration. |
| **Music** | [music.md](./music.md) | Background music, preset tracks, MP3 upload, autoplay configuration. |
| **SEO** | [seo.md](./seo.md) | Metadata management, OpenGraph tags, slug validation, sitemap generation. |
| **Analytics** | [analytics.md](./analytics.md) | View counting, RSVP statistics, page visit tracking, platform-wide metrics. |
| **Audit Log** | [audit-log.md](./audit-log.md) | Immutable action logging, admin activity tracking, entity change history. |
| **Export** | [export.md](./export.md) | RSVP Excel export, guest list export, CSV download. |

## Module Dependency Graph

```
Authentication  ──────┐
                       │
Users  ────────────────┤
                       │
Subscriptions  ────────┤
                       │
Packages  ─────────────┤
                       │
Payments  ─────────────┤
                       │
Events  ───────────────┼──►  Shared services (JWT, DB, Config)
                       │
Templates  ────────────┤
                       │
Guests  ───────────────┤
                       │
RSVP  ─────────────────┤
                       │
Guestbook  ────────────┤
                       │
Digital Gifts  ────────┤
                       │
Gallery  ──────────────┤
                       │
Music  ────────────────┤
                       │
WhatsApp  ─────────────┘
```

## Layer Mapping

Each module in the Go backend follows the standard three-layer architecture:

- **Handler** (`internal/api/handler/`) — HTTP request parsing, transport validation, response formatting
- **Service** (`internal/service/`) — Business logic, authorization, orchestration
- **Repository** (`internal/repository/`) — GORM database operations

## Reading Order

1. Start with **Authentication** and **Users** to understand user identity.
2. Read **Packages** and **Subscriptions** to understand the pricing model.
3. Read **Payments** to understand the revenue flow.
4. Read **Events** and **Invitation Editor** to understand the core product.
5. Read **Templates** to understand theming.
6. Read **Guests**, **RSVP**, **Guestbook**, **Digital Gifts**, **Gallery**, **Music** for feature modules.
7. Read **WhatsApp** for sharing integration.
8. Read **SEO**, **Analytics**, **Audit Log**, **Export**, **Admin** for cross-cutting concerns.