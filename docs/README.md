# Documentation

## Project

Platform Undangan Pernikahan Digital & Cetak — SaaS wedding invitation platform for the Indonesia market. Decoupled architecture with Go backend and Next.js frontend.

## Source of Truth

Priority order when resolving conflicts:
1. Actual production requirements
2. Approved PRD (`/prd.md`)
3. Architecture Decision Records (`docs/decisions/`)
4. API contracts (`docs/api/`)
5. Database schema (`docs/database/`)
6. Module documentation (`docs/modules/`)
7. Implementation/code
8. TODO/planned documentation

## Architecture

```
Frontend (Next.js)  →  API (Go REST)  →  Database (PostgreSQL)
                          ↕
                    Midtrans Payment Gateway
```

See `docs/architecture/overview.md`.

## Modules

| Module | Doc | Purpose |
|--------|-----|---------|
| Authentication | `docs/modules/authentication.md` | Login, register, JWT tokens |
| Users | `docs/modules/users.md` | User profile & account |
| Admin | `docs/modules/admin.md` | Platform management |
| Subscriptions | `docs/modules/subscriptions.md` | Subscription lifecycle |
| Packages | `docs/modules/packages.md` | Package definitions & limits |
| Payments | `docs/modules/payments.md` | Midtrans transaction lifecycle |
| Events | `docs/modules/events.md` | Wedding invitation core entity |
| Invitation Editor | `docs/modules/invitation-editor.md` | Invitation content config |
| Templates | `docs/modules/templates.md` | Invitation themes/templates |
| Guests | `docs/modules/guests.md` | Guest management |
| WhatsApp | `docs/modules/whatsapp.md` | WhatsApp link & broadcast |
| RSVP | `docs/modules/rsvp.md` | Attendance response |
| Guestbook | `docs/modules/guestbook.md` | Guest messages |
| Digital Gifts | `docs/modules/digital-gifts.md` | Bank/QRIS/gift info |
| Gallery | `docs/modules/gallery.md` | Photo & video gallery |
| Music | `docs/modules/music.md` | Background music |
| SEO | `docs/modules/seo.md` | Metadata, slug, discoverability |
| Analytics | `docs/modules/analytics.md` | Platform analytics |
| Audit Log | `docs/modules/audit-log.md` | Activity logging |
| Export | `docs/modules/export.md` | Data export (RSVP, guests) |

## API

Base path: `/api/v1`

See `docs/api/conventions.md` for standard formats, authentication, pagination, error handling.

## Database

PostgreSQL with GORM ORM.

See `docs/database/schema.md` and `docs/database/relationships.md`.

## Frontend

Next.js App Router with TypeScript.

Route groups:
- `(public)/` — Landing, pricing, public invitation `[slug]`
- `admin/` — Admin dashboard
- `user/` — User/pengantin dashboard

See `docs/frontend/routing.md`.

## Backend

Go REST API with layered architecture: Handler → Service → Repository.

See `docs/backend/project-structure.md`.

## Integrations

> Split rule: `docs/` root and topical folders hold product/architecture overview documentation; implementation-level details live under `docs/internal/`.

- Midtrans — `docs/internal/integrations/midtrans.md`
- WhatsApp — `docs/internal/integrations/whatsapp.md`
- Object Storage — `docs/internal/integrations/storage.md`
- Email — `docs/internal/integrations/email.md`

## Security

See `docs/security/security-overview.md`.

## Testing

See `docs/testing/strategy.md`.

## Development

See `docs/development/setup.md`.

## Deployment

See `docs/deployment/environments.md`.

## Roadmap

See `docs/roadmap/roadmap.md`.

## AI Agent Instructions

Before any code change, the AI agent MUST:
1. Read `docs/context/coding-agent-context.md`
2. Read relevant module documentation
3. Read API/database contracts
4. Inspect existing implementation
5. Identify discrepancies
6. Plan minimal changes
7. Implement
8. Write/update tests
9. Run validation
10. Update documentation
11. Report changed files, tests, and unresolved TODOs
