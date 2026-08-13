# Frontend Routing

## App Router Structure

This project uses Next.js App Router with TypeScript. Routes are organized into three groups: public routes, admin dashboard, and user dashboard. There are no route groups (parenthesized folders) in the current structure — all routes are flat within `app/`.

## Route Map

```
/                           Landing page (public)
/pricing                    Pricing plans (public)
/[slug]                     Public invitation page (dynamic, SEO-friendly)
/admin                      Admin dashboard overview
/admin/users                User management
/admin/transactions         Transaction monitoring
/admin/packages             Package/subscription plan management
/admin/templates            Template management
/user                       User dashboard overview
/user/editor                Invitation content editor
/user/guests                Guest list management
/user/rsvp                  RSVP recap & guestbook
/user/billing               Subscription & payment history
```

## Dynamic Route: `/[slug]`

The `[slug]` route is a **dynamic segment** that renders the public-facing wedding invitation page. The slug corresponds to the event's unique URL identifier (e.g., `/andi-sinta`). This route:

- Is a **server component** by default for optimal SEO.
- Calls `generateMetadata({ params })` to produce dynamic `<title>`, `<meta>`, OpenGraph, and Twitter card tags based on the couple's names and wedding date.
- Fetches invitation data from `GET /api/v1/public/events/:slug` at request time.
- Renders section-based invitation components (hero, couple, event details, gallery, RSVP, guestbook, digital gifts, music).

404 handling: if the slug does not exist or the event is unpublished, the server component returns `notFound()`.

## Layout Nesting

```
Root Layout (app/layout.tsx)
├── Public pages (/, /pricing, /[slug])
├── Admin Layout (app/admin/layout.tsx)
│   ├── /admin
│   ├── /admin/users
│   ├── /admin/transactions
│   ├── /admin/packages
│   └── /admin/templates
└── User Layout (app/user/layout.tsx)
    ├── /user
    ├── /user/editor
    ├── /user/guests
    ├── /user/rsvp
    └── /user/billing
```

## Layout Behavior

| Route | Layout | Auth Required | Sidebar |
|-------|--------|---------------|---------|
| `/` | Root | No | No |
| `/pricing` | Root | No | No |
| `/[slug]` | Root | No | No |
| `/admin/*` | Admin | JWT + admin role | Yes (admin sidebar) |
| `/user/*` | User | JWT + user role | Yes (user sidebar) |

## Not Found & Error Pages

- `app/not-found.tsx` — Custom 404 page for all unmatched routes.
- `app/error.tsx` — Global error boundary with retry button.

## Route Protection Strategy

Protected routes (`/admin/*` and `/user/*`) rely on a **middleware** (or layout-level check) that reads the JWT from an httpOnly cookie, verifies it, and redirects to `/login` if invalid. Admin routes additionally check `role === 'admin'` and redirect to `/user` if the role is `user`. If no token exists, the user is redirected to `/login?redirect=<original_path>` so they return after logging in.

## Future Considerations

- Route groups (`(public)`, `(dashboard)`) may be introduced to scope specific layouts without affecting the URL.
- Loading states: `loading.tsx` files at each segment level for Suspense fallbacks.
- Parallel routes for modals in the editor (e.g., `/user/editor/@preview`).
- API routes (`app/api/`) if a Next.js API proxy layer is needed for server-side cookie forwarding.
