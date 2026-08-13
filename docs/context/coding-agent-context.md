# Coding Agent Context

This document provides concise instructions for AI coding agents (e.g., OpenCode) when working on this project. Read this first before making any code changes.

## Project Identity

- **Platform:** SaaS Undangan Pernikahan Digital & Cetak (Indonesia Market)
- **Languages:** Go (backend), TypeScript (frontend)
- **Database:** PostgreSQL with GORM
- **Auth:** JWT
- **Payment:** Midtrans Snap

## Repository Structure

```
/
├── backend/
│   ├── cmd/           # Main entry points
│   ├── internal/
│   │   ├── api/       # Handlers, middleware
│   │   ├── service/   # Business logic
│   │   ├── repository/ # Database access
│   │   ├── model/     # Database models
│   │   ├── dto/       # Request/response DTOs
│   │   ├── config/    # Configuration
│   │   └── pkg/       # Shared utilities
│   ├── migrations/    # Database migrations
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── app/       # Next.js App Router pages
│   │   ├── components/ # React components
│   │   ├── lib/       # Utilities, API client
│   │   ├── hooks/     # Custom hooks
│   │   └── types/     # TypeScript types
│   └── package.json
├── docs/              # Full documentation
└── AGENTS.md          # Agent instructions
```

## Coding Rules

1. **Handlers are thin.** Parse request, validate basic format, call service, return response. No business logic.
2. **Services contain business logic.** One service per module/domain.
3. **Repositories handle database.** One repository per model. No business logic.
4. **Validate all external input.** Request validation at handler level. Business rule validation in service layer.
5. **Enforce authorization server-side.** Never trust frontend roles/permissions.
6. **Never trust frontend payment status.** Always verify via Midtrans webhook.
7. **No hardcoded secrets.** Use environment variables or config.
8. **No fake API responses in production.**
9. **No silent schema changes.**
10. **No large unrelated refactors.**

## When Modifying Code

1. Read relevant docs in `/docs`
2. Inspect existing implementation
3. Identify all affected modules
4. Plan minimal change
5. Implement
6. Update/add tests
7. Run `go test ./...` or `npm test`
8. Update documentation if contracts changed
9. Report: changed files, tests, TODOs

## Database Rules

- Use GORM auto-migration for development only.
- Use versioned SQL migrations for production.
- Never drop columns without confirmation.
- Soft-delete where appropriate (deleted_at).
- Timestamps: created_at, updated_at, deleted_at.

## API Rules

- RESTful, base path `/api/v1`
- Standard error format: `{ "error": { "code": "...", "message": "..." } }`
- Standard pagination: `{ "data": [...], "pagination": { "page": 1, "per_page": 20, "total": 100 } }`
- Throttle unauthenticated public routes.
- Webhook endpoints must verify source (Midtrans signature).

## Payments

- Midtrans Snap (not Core API).
- Midnotif is outdated. Use Snap with webhook.
- Frontend callback is UI hint only.
- Subscription activates ONLY after Midtrans webhook confirms settlement.
- Transaction status: `pending`, `settlement`, `expire`, `deny`, `cancel`, `refund`.

## Key Business Rules

- Free subscription: 7 days, limited features.
- One active subscription per user at a time.
- Invitation slug must be unique.
- Guest can RSVP only once.
- Digital gift info is static content, not real-time payment.
- Public invitation is viewable without auth.
- Admin can override subscription status.

## Frontend Rules

- Use Next.js App Router.
- Server components by default.
- Client components when interactivity needed.
- API calls via lib/api-client.ts.
- Form validation on client-side + server-side validation.
- Optimistic updates for CRUD where appropriate.
- Do not store payment data on frontend.

## Documentation Priority

When updating docs:
1. API contracts
2. Database schema
3. Module docs
4. Architecture docs
5. Development docs

## Agent Communication

When reporting results, include:
- Files created/modified
- Tests written/updated
- Any deviations from docs
- Unresolved TODOs
- Trade-offs made
