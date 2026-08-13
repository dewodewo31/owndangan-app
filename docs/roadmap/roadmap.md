# Roadmap

## Status

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1: Backend + DB | [ ] Planned | Go REST API, auth, CRUD, Midtrans integration |
| Phase 2: Frontend | [ ] Planned | Next.js setup, dashboards, editor, public pages |
| Phase 3: Templates & Tiers | [ ] Planned | Template system, feature entitlement, export |
| Phase 4: Production Readiness | [ ] Planned | Production Midtrans, SEO, load testing, deployment |

## Phase 1: Backend + Database

**Goal:** Complete Go REST API with database schema, authentication, and core CRUD.

- [ ] Project scaffolding (Go modules, project structure)
- [ ] Database schema and migrations
- [ ] User registration and login (JWT)
- [ ] Package CRUD (admin)
- [ ] Subscription management
- [ ] Event CRUD
- [ ] Guest CRUD
- [ ] RSVP submission (public)
- [ ] Guestbook (public + moderation)
- [ ] Digital gift configuration
- [ ] Gallery and music upload
- [ ] Midtrans Snap integration
- [ ] Midtrans webhook handler
- [ ] Admin endpoints
- [ ] Audit logging
- [ ] Basic analytics

## Phase 2: Next.js Frontend

**Goal:** Complete frontend with user dashboard, admin dashboard, and public invitation pages.

- [ ] Next.js project setup
- [ ] Authentication pages (login, register)
- [ ] Landing page
- [ ] Pricing page
- [ ] User dashboard layout
- [ ] User dashboard overview
- [ ] Invitation editor (9 sections)
- [ ] Guest management page
- [ ] RSVP recap page
- [ ] Billing page (Midtrans Snap)
- [ ] Admin dashboard layout
- [ ] Admin dashboard overview
- [ ] Admin user management
- [ ] Admin transaction monitoring
- [ ] Admin package management
- [ ] Admin template management
- [ ] Public invitation page (`/[slug]`)
- [ ] SEO metadata and OpenGraph

## Phase 3: Templates & Tiers

**Goal:** Complete template system, feature entitlement enforcement, and advanced features.

- [ ] Template system (upload, assign, render)
- [ ] Feature entitlement enforcement (backend)
- [ ] Feature tier gating (frontend)
- [ ] WhatsApp link generator
- [ ] RSVP Excel export (Pro)
- [ ] Guest list CSV import
- [ ] QR guestbook (Pro)
- [ ] Custom domain (Pro)
- [ ] Watermark removal (Premium+)

## Phase 4: Production Readiness

**Goal:** Production deployment, performance optimization, and security hardening.

- [ ] Midtrans production switch
- [ ] Load testing (backend)
- [ ] Performance optimization
- [ ] SEO finalization (sitemap, robots.txt, JSON-LD)
- [ ] Security audit
- [ ] Rate limiting enforcement
- [ ] Error monitoring (Sentry)
- [ ] CI/CD pipeline
- [ ] Production deployment
- [ ] Monitoring and alerting setup
- [ ] Documentation finalization
- [ ] Backup and disaster recovery