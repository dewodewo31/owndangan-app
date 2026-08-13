# Authentication

## Auth Flow Overview

Authentication uses JWT access tokens with refresh token rotation. The backend (Go) issues tokens, and the frontend manages them transparently via the `AuthProvider` and API client interceptor.

## Token Storage

The project supports two storage strategies. The recommended approach is httpOnly cookies.

**Strategy A: httpOnly cookie (recommended, production)**

- The backend sets the access token as an httpOnly, Secure, SameSite=Strict cookie on login/refresh.
- The frontend never reads the token in JavaScript. It is sent automatically on every request via `credentials: 'include'`.
- Token refresh is triggered by the backend returning a 401, which the API client interceptor catches and calls `POST /api/v1/auth/refresh`. The backend sets a new cookie.
- Advantages: immune to XSS. No token in JS memory.
- Trade-off: requires the frontend and backend to share the same root domain (or use a proxy).

**Strategy B: localStorage (simpler, development)**

- The access token is stored in `localStorage` under the key `access_token`.
- The refresh token is stored in `localStorage` under `refresh_token`.
- The API client reads the access token from localStorage and injects it as `Authorization: Bearer <token>`.
- On 401, the client calls `POST /api/v1/auth/refresh` with the stored refresh token, receives a new pair, and updates localStorage.
- Advantages: simple to implement, works across any domain.
- Trade-off: vulnerable to XSS; tokens are accessible to any JavaScript on the page.

## Login Flow

```
User enters email + password on /login
  → POST /api/v1/auth/login
  → Backend validates credentials, returns { user, access_token, refresh_token }
  → Frontend stores tokens (cookie or localStorage)
  → AuthProvider sets user state
  → Redirect to /user or /admin (based on role)
```

The login page reads `redirect` from the query string (e.g., `/login?redirect=/user/editor`). If present, the user is redirected there after login. If absent, the redirect target is determined by role: `/user` for user role, `/admin` for admin role.

## Registration Flow

```
User fills registration form on /register
  → POST /api/v1/auth/register
  → Backend creates user, auto-assigns free 7-day subscription, returns tokens
  → Frontend stores tokens, sets user, redirects to /user
```

## Protected Routes

Protected routes are enforced at the **layout level**, not in individual pages. This ensures consistent behavior across all pages in a segment.

**Admin layout** (`app/admin/layout.tsx`):
```
1. Check if user is authenticated (JWT valid)
2. If not → redirect to /login?redirect=/admin
3. Check if user.role === 'admin'
4. If not → redirect to /user
5. Render admin sidebar and children
```

**User layout** (`app/user/layout.tsx`):
```
1. Check if user is authenticated (JWT valid)
2. If not → redirect to /login?redirect=/user
3. Check if user.role === 'user'
4. If not → redirect to /admin
5. Render user sidebar and children
```

The auth check in layouts uses a server-side read of the httpOnly cookie (via `cookies()` from `next/headers`). If the token is valid, the layout decodes the user from the JWT payload and passes it to the auth context. This avoids a flash of unauthenticated content before the client-side check runs.

## Token Refresh

The API client handles token refresh transparently:

1. Any API call receives a 401 response.
2. The client checks if the error is a genuine 401 (not a network error).
3. The client calls `POST /api/v1/auth/refresh` with the stored refresh token.
4. If the refresh succeeds, the new tokens are stored and the original request is retried.
5. If the refresh fails (refresh token expired or revoked), the user is logged out and redirected to `/login`.

The refresh endpoint uses rotation-based tokens: each refresh invalidates the previous refresh token and issues a new pair. If a revoked refresh token is presented, all sessions are invalidated (theft detection).

## Logout

```
User clicks Logout
  → POST /api/v1/auth/logout (with refresh token in body)
  → Backend invalidates refresh token
  → Frontend clears stored tokens
  → AuthProvider sets user to null
  → Redirect to /
```

The logout button is present in both admin and user sidebars.

## Suspended User Handling

If a user account is suspended by an admin, any API call returns `FORBIDDEN (403)` with code `USER_SUSPENDED`. The API client interceptor catches this and:

1. Logs the user out locally.
2. Shows a toast: "Akun Anda telah dinonaktifkan. Hubungi admin untuk informasi lebih lanjut."
3. Redirects to `/login`.

## Admin-Only Routes

The admin layout enforces the `role === 'admin'` check. Additionally, the backend enforces authorization on every admin endpoint. Frontend route hiding is a UX convenience, not a security measure.

## Related API Documentation

See `docs/api/authentication.md` for the full API contract (register, login, refresh, logout).