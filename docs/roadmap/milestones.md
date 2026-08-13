# Milestones

## M1: Backend Foundation

**Target:** Week 1-2

**Deliverables:**
- Go project structure with Handler → Service → Repository layers
- PostgreSQL database with all tables
- JWT authentication (register, login, refresh)
- User CRUD
- Package CRUD + seed data
- Subscription management
- Event CRUD (create, read, update, delete, publish, unpublish)
- Guest CRUD

**Definition of Done:**
- All API endpoints return correct responses
- Tests for service layer
- Test coverage > 60%

## M2: Payment Integration

**Target:** Week 2-3

**Deliverables:**
- Midtrans Snap token creation
- Snap.js frontend integration
- Midtrans webhook handler
- Signature verification
- Transaction status management
- Subscription activation via webhook
- Idempotency handling

**Definition of Done:**
- Successful sandbox payment flow
- Webhook signature verification works
- Subscription activates after settlement
- All transaction statuses handled correctly

## M3: Guest & Interaction Features

**Target:** Week 3-4

**Deliverables:**
- RSVP submission (public)
- Guestbook messages (public + moderation)
- Digital gift configuration
- Gallery upload
- Music upload/presets
- WhatsApp link generator
- Admin dashboard endpoints

**Definition of Done:**
- Public invitation page renders all data
- RSVP submission works end-to-end
- Guestbook moderation works
- Admin can manage users and transactions

## M4: Frontend Dashboard

**Target:** Week 4-5

**Deliverables:**
- Next.js project with App Router
- Auth pages (login, register)
- User dashboard (overview, editor, guests, RSVP, billing)
- Admin dashboard (overview, users, transactions, packages, templates)
- Public invitation page (`/[slug]`)
- API client integration

**Definition of Done:**
- User can create and publish invitation
- Guest can RSVP and send messages
- Admin can view dashboard stats
- All pages are responsive

## M5: Template & Entitlement

**Target:** Week 5-6

**Deliverables:**
- Template system (upload, assign, render)
- Feature entitlement enforcement
- Plan-specific feature gating
- RSVP Excel export
- Guest CSV import
- Watermark removal

**Definition of Done:**
- Different plans see different features
- Template assignment works
- Export generates correct file
- Import handles duplicates and limits

## M6: Production Release

**Target:** Week 7-8

**Deliverables:**
- Midtrans production mode
- SEO optimization (sitemap, metadata, OpenGraph)
- Load testing and optimization
- CI/CD pipeline
- Production deployment
- Monitoring setup
- Security audit

**Definition of Done:**
- Production deployment live
- Load test passes (100 concurrent users)
- All monitoring alerts configured
- Security audit passes