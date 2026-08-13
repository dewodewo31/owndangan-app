# Layouts

## Layout Hierarchy

Three levels of layout nesting exist: root, admin, and user. No route groups are used in the current structure; layouts are applied via folder nesting.

```
app/layout.tsx              Root layout — fonts, metadata, providers
├── Pages without layout override (/, /pricing, /[slug])
├── app/admin/layout.tsx    Admin layout — sidebar, header, auth guard
└── app/user/layout.tsx     User layout — sidebar, header, auth guard
```

## Root Layout (`app/layout.tsx`)

The root layout wraps every page. It is the only layout rendered for public-facing routes.

**Responsibilities:**

- Loads Google Fonts (Playfair Display for headings, Inter for body) via `next/font`.
- Sets global `<html lang="id">` and default metadata (title template, description).
- Wraps children in `AuthProvider` (React context for auth state).
- Wraps children in `ThemeProvider` if client-side theme toggling is needed.
- Renders `<Toaster />` for toast notifications.
- Does NOT include any header, footer, or navigation bar.

**File:**
```tsx
// app/layout.tsx (server component by default)
export default function RootLayout({ children }) {
  return (
    <html lang="id">
      <body className={`${playfair.variable} ${inter.variable} font-sans`}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
```

## Admin Layout (`app/admin/layout.tsx`)

The admin layout applies to all routes under `/admin/*`.

**Responsibilities:**

- Auth guard: checks for JWT in httpOnly cookie. Redirects to `/login` if absent.
- Role guard: checks that `role === 'admin'`. Redirects to `/user` if role is `user`.
- Renders a persistent sidebar with navigation links:
  - Dashboard (`/admin`)
  - Users (`/admin/users`)
  - Transactions (`/admin/transactions`)
  - Packages (`/admin/packages`)
  - Templates (`/admin/templates`)
- Renders a top header bar showing admin name, notification bell, and logout button.
- Sidebar is collapsible on mobile; uses a hamburger toggle.

**Sidebar behavior:**
- Active link is highlighted based on `pathname` from `usePathname()`.
- Links use `next/link` for client-side navigation.
- Sidebar collapses to icons-only at `md` breakpoint, hidden behind toggle at `sm`.

## User Layout (`app/user/layout.tsx`)

The user layout applies to all routes under `/user/*`.

**Responsibilities:**

- Auth guard: same JWT check as admin layout.
- Role guard: redirects admins to `/admin`.
- Renders a user-sidebar with navigation:
  - Dashboard (`/user`)
  - Editor (`/user/editor`)
  - Guests (`/user/guests`)
  - RSVP / Guestbook (`/user/rsvp`)
  - Billing (`/user/billing`)

**Sidebar items differ from admin.**
- "Billing" shows the user's current subscription badge (Basic / Premium / Pro) next to the link.
- Sidebar footer shows the user's name and email.

## Auth Guard Implementation

```tsx
// Shared auth guard function used by admin/layout.tsx and user/layout.tsx
async function getAuthSession(): Promise<Session | null> {
  // Read JWT from cookie (server-side via cookies())
  // Verify token (decode + validate expiry)
  // Fetch user from GET /api/v1/users/me
  // Return session or null
}
```

If `getAuthSession()` returns `null`, the layout calls `redirect('/login')`. If the role does not match the layout's required role, the layout calls `redirect('/' + correctDashboard)`.

## Public Invitation Layout

The `[slug]` route does NOT have its own layout. It inherits the root layout directly. This ensures:

- No sidebar or dashboard chrome appears on the public invitation page.
- The invitation renders full-screen for an immersive experience.
- The page appears as a standalone website when shared on social media.

## Future Considerations

- Loading skeletons: each layout can define `loading.tsx` at its segment level.
- Error boundaries: each layout can define `error.tsx` for segment-scoped error handling.
- Parallel routes: the editor layout may use `@preview` slot for live preview pane.
