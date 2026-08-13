# Frontend Architecture

## App Router Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── layout.tsx              ← Root layout (fonts, metadata, providers)
│   │   ├── page.tsx                ← Landing page
│   │   ├── pricing/
│   │   │   └── page.tsx            ← Pricing page
│   │   ├── [slug]/
│   │   │   └── page.tsx            ← Public invitation page (dynamic)
│   │   ├── admin/
│   │   │   ├── layout.tsx          ← Admin dashboard layout (sidebar, header)
│   │   │   ├── page.tsx            ← Admin dashboard overview
│   │   │   ├── users/
│   │   │   ├── transactions/
│   │   │   ├── packages/
│   │   │   └── templates/
│   │   ├── user/
│   │   │   ├── layout.tsx          ← User dashboard layout (sidebar, header)
│   │   │   ├── page.tsx            ← User dashboard overview
│   │   │   ├── editor/
│   │   │   ├── guests/
│   │   │   ├── rsvp/
│   │   │   └── billing/
│   │   └── not-found.tsx           ← 404 page
│   │   └── error.tsx               ← Global error boundary
│   ├── components/
│   │   ├── ui/                     ← Shared UI components (Button, Input, Modal...)
│   │   ├── forms/                  ← Form components
│   │   ├── layout/                 ← Header, Footer, Sidebar, Nav
│   │   ├── invitation/            ← Invitation section components
│   │   └── admin/                 ← Admin-specific components
│   ├── lib/
│   │   ├── api-client.ts           ← API client (fetch wrapper)
│   │   ├── auth.ts                 ← Auth utilities
│   │   └── utils.ts               ← General helpers
│   ├── hooks/                      ← Custom React hooks
│   ├── types/                      ← TypeScript types/interfaces
│   └── providers/                  ← React context providers
├── public/
│   ├── images/
│   └── favicon.ico
├── package.json
├── next.config.ts
├── tsconfig.json
└── tailwind.config.ts
```

## Route Groups

### `(public)/` — Public routes (no auth)
- Landing page, pricing, about
- Public invitation `/[slug]`

### `admin/` — Admin dashboard
- All routes prefixed with `/admin`
- Protected by JWT (admin role)
- Nested layout with admin sidebar

### `user/` — User dashboard
- All routes prefixed with `/user`
- Protected by JWT (user role)
- Nested layout with user sidebar

### `[slug]` — Public invitation
- Dynamic route for public invitation pages
- SEO-optimized with server components
- No auth required

## Component Categories

### Server Components (default)
- Pages
- Layouts
- Static content
- Public invitation sections (read-only)

### Client Components (`'use client'`)
- Forms
- Interactive editors
- Midtrans Snap button
- Real-time RSVP form
- Guestbook form
- Dashboard charts
- Modals and dialogs

## Data Fetching

- Server components: fetch data directly from backend API.
- Client components: use hooks (SWR/TanStack Query) for data fetching and caching.
- Public invitation page: server-side fetch for SEO.
- Dashboard pages: client-side fetch for real-time updates.

## Authentication

- JWT stored in httpOnly cookie or localStorage.
- Auth context provider wraps the app.
- Protected routes redirect to login if unauthenticated.
- Admin routes check for admin role.

## Related Documentation

- `docs/frontend/routing.md`
- `docs/frontend/layouts.md`
- `docs/frontend/components.md`
- `docs/frontend/api-client.md`
- `docs/frontend/authentication.md`
