# ADR-001: Initial Architecture Decision

- **Status**: Accepted
- **Date**: 2025-01-15
- **Author**: Architecture Team

## Context

We are building a SaaS wedding invitation platform for the Indonesia market. The platform needs to handle user registration, invitation creation with customizable templates, payment processing via Indonesian payment gateways, RSVP management, and guest list management. The system must scale to handle thousands of concurrent users during peak wedding season.

Key requirements:
- Decoupled frontend and backend for independent scaling.
- Rapid development with familiar technologies.
- Support for Midtrans as the primary payment gateway.
- Secure authentication and authorization.
- Indonesian market localization (Bahasa Indonesia, timezone Asia/Jakarta).

## Options

### Option 1: Monolith (Rails/Django)

A single application serving both API and views.

**Pros**: Simple deployment, single codebase, fast initial development.
**Cons**: Tight coupling, cannot scale frontend/backend independently, harder to maintain as the team grows.

### Option 2: Decoupled Go + Next.js (Selected)

Backend API in Go, frontend in Next.js, communicating via REST.

**Pros**:
- Independent scaling of API and frontend.
- Go provides excellent performance and low resource usage.
- Next.js offers SSR for SEO, ISR for static pages, and a rich ecosystem.
- Clear separation of concerns (handler → service → repository).
- TypeScript on frontend provides type safety.
- Go compiles to a single static binary — easy to deploy.

**Cons**:
- Two codebases to maintain.
- API contract must be versioned and documented.
- More complex local development setup.

### Option 3: BFF (Backend for Frontend) with GraphQL

A single BFF layer in Go/Node.js, frontend talks to BFF via GraphQL.

**Pros**: Flexible queries, reduces over-fetching.
**Cons**: GraphQL complexity, caching challenges, overkill for a CRUD-heavy app.

## Decision

We chose **Option 2: Decoupled Go + Next.js** with the following technology stack:

- **Backend**: Go 1.22+ with Chi router, pgx for PostgreSQL, testify for tests.
- **Frontend**: Next.js 14+ with TypeScript, Tailwind CSS, React Query.
- **Database**: PostgreSQL 16+ with goose migrations.
- **Payment**: Midtrans Snap (REST API + Snap.js).
- **Auth**: JWT (access + refresh tokens).
- **Deployment**: Docker containers, Vercel for frontend, Railway/VPS for backend.

## Consequences

### Positive

- Backend and frontend teams can work independently.
- Go backend handles high concurrency efficiently.
- Next.js provides SEO-friendly pages for invitation public pages.
- TypeScript reduces runtime errors in the frontend.
- Containerization makes deployment consistent across environments.

### Negative / Trade-offs

- API contract maintenance: Every endpoint must be documented and versioned.
- Local setup requires both Go and Node.js toolchains.
- Two CI/CD pipelines to maintain.
- Payment webhook handling requires careful idempotency logic.

### Compliance

- Architecture follows the rules defined in `AGENTS.md`.
- All business logic lives in the backend (never trust the frontend).
- Payment status is determined server-side via webhook notifications.
- Secrets are never exposed to the client bundle.