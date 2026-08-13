# Architecture Overview

## Principles

1. **Decoupled architecture.** Go backend and Next.js frontend are independent deployments.
2. **Thin handlers.** Handlers parse request, call service, return response. No business logic.
3. **Service layer.** Business logic resides in services. Services are testable without HTTP.
4. **Repository layer.** Database access isolated in repositories.
5. **Dependency direction.** Handler → Service → Repository (inward only).
6. **Stateless backend.** JWT contains session information. No server-side session store.
7. **Webhook as source of truth.** Payment status is authoritative only from Midtrans webhook.
8. **Server-side validation.** All validation happens server-side. Frontend validation is for UX only.

## Key Architecture Decisions

### Why Decoupled?
- Independent scaling of frontend and backend.
- Frontend can be served via CDN (Vercel/Cloudflare).
- Backend can be deployed as container or binary.
- Each layer can be tested independently.
- Team can work on frontend and backend in parallel.

### Why Go?
- Strong performance for API workload.
- Simple concurrency model.
- Single binary deployment.
- Good fit for REST API + PostgreSQL workload.

### Why Next.js App Router?
- Server Components for SEO-friendly public pages.
- App Router for nested layouts.
- TypeScript for type safety.
- Large ecosystem.

### Why Midtrans Snap?
- Snap (not Core API) provides hosted payment page.
- Reduces PCI compliance scope.
- Supports all Indonesian payment methods.
- Webhook-based notification.

### Why GORM?
- Auto-migration for development.
- Good PostgreSQL support.
- Popular in Go ecosystem.
- Avoids raw SQL boilerplate.

## Architecture Diagram

See `docs/context/architecture-context.md` for detailed diagram.

## Route Architecture

### Frontend Route Groups

```
(app)             → landing, pricing, static pages
(slug)/[slug]     → public invitation page
admin/            → admin dashboard (prefixed /admin)
user/             → user dashboard (prefixed /user)
api/              → Next.js API routes (proxy to backend, optional)
```

### Backend API Versioning

- All routes under `/api/v1/`.
- Version prefix allows backwards-compatible changes.
- Breaking changes → `/api/v2/`.

## Data Flow

### Invitation Creation Flow
```
User → POST /api/v1/events → Handler → Service (validate subscription, limits) → 
Repository (save to PostgreSQL) → Return event ID → User configures in editor
```

### Payment Flow
```
User → POST /api/v1/payments/snap → Handler → Service (create Midtrans transaction) → 
Repository (save transaction record) → Return Snap token → Frontend launches Snap.js → 
Midtrans handles payment → Midtrans sends webhook → 
POST /api/v1/webhook/midtrans → Handler (verify signature) → 
Service (update transaction, activate subscription) → Repository (save)
```

### Public Invitation Flow
```
Guest → GET /api/v1/public/events/:slug → Handler → Service (check event published) → 
Repository (load event + sections) → Return JSON → Frontend renders invitation
```

### Guest RSVP Flow
```
Guest → GET /[slug] (public invitation) → Fill RSVP form → 
POST /api/v1/public/rsvps → Handler → Service (validate guest token, check duplicates) → 
Repository (save/update RSVP) → Return confirmation
```

## Scalability Considerations

- Backend is stateless → horizontal scaling via load balancer.
- Database connection pool.
- Read replicas for read-heavy workloads (future).
- Redis cache for public invitation pages (future).
- CDN for static assets.
- Object storage for user-uploaded media.

## Failure Modes

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Midtrans down | Can't create new payments | Graceful error, retry logic |
| Database down | All services unavailable | Connection pool, monitoring |
| Webhook delivery delay | Subscription activation delayed | Idempotency, manual admin override |
| File storage down | Gallery/music unavailable | Graceful degradation (placeholder) |
| Frontend down | Users can't access UI | CDN caching of static assets |

## Future Considerations

- Redis caching layer.
- Database read replicas.
- Background job queue (RabbitMQ / Redis Queue).
- WebSocket for real-time RSVP updates.
- Mobile app (Flutter/React Native).
- Multi-language support.
- Printed invitation order system integration.
