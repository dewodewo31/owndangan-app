# State Management

## Approach

This project uses a **three-tier state management** approach. Each tier addresses a specific concern without overlap:

| Concern | Solution | Scope |
|---------|----------|-------|
| Server / API data | SWR or TanStack Query | Remote state, caching, revalidation |
| Auth session | React Context | Global auth state (user, token) |
| Form / UI state | Local `useState` / `useReducer` | Ephemeral, non-shared state |

## Server State (SWR / TanStack Query)

All data fetched from the API is managed by a caching library, not by manual fetch + setState. This provides:

- **Automatic caching** — repeated requests within the stale time return cached data.
- **Background revalidation** — data refreshes when the user refocuses the tab or navigates.
- **Optimistic updates** — mutations update the cache immediately, then sync with the server.
- **Pagination** — infinite scroll and page-based lists are handled natively.

**Pattern:**

```tsx
// hooks/useEvents.ts
import useSWR from 'swr';
import { apiClient } from '@/lib/api-client';

export function useEvents() {
  return useSWR('/events', (url) => apiClient.get(url));
}

// In a component
const { data, error, isLoading, mutate } = useEvents();
```

**Key configuration:**

- `dedupingInterval: 2000` — deduplicate requests within 2 seconds.
- `revalidateOnFocus: true` — re-fetch when the user returns to the tab.
- `errorRetryCount: 3` — retry failed requests up to 3 times.
- `fallbackData` — used for SSR hydration on dashboard pages.

**Module-specific hooks:**

- `useEvents()` — list user's invitations.
- `useEvent(id)` — single event details.
- `useGuests(eventId)` — guest list for an event.
- `useRsvps(eventId)` — RSVP recap.
- `useTransactions()` — user's payment history (billing page).
- `useAdminDashboard()` — admin overview stats.
- `useAdminUsers()` — admin user list.
- `usePackages()` — package definitions.
- `useTemplates()` — template list.
- `useAnalytics()` — admin analytics data.

## Auth State (React Context)

The `AuthProvider` wraps the root layout and provides auth state to all components.

```tsx
// providers/auth-provider.tsx
interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
}
```

**Implementation details:**

- On mount, `AuthProvider` attempts to read the JWT from an httpOnly cookie (via `GET /api/v1/auth/me`).
- If the call succeeds, `user` is set. If it fails (401), `user` is `null` and no redirect occurs — the guard in the layout handles it.
- `login()` calls `POST /api/v1/auth/login`, receives tokens, stores them, and sets `user`.
- `logout()` calls `POST /api/v1/auth/logout` and clears user state.
- `refreshToken()` is called by the API client interceptor when a 401 is received, transparently retrying the original request.

**Usage:**

```tsx
const { user, logout } = useAuth();
// user.name, user.email, user.role, user.status
```

## Local State (Forms)

Form state is managed locally using `useState` or `useReducer`. No global state library is used for forms.

**Simple forms** (login, register, add guest): `useState` with an object.

```tsx
const [form, setForm] = useState({ email: '', password: '' });
```

**Complex forms** (invitation editor with 9 sections): `useReducer` with a structured reducer.

```tsx
const [sections, dispatch] = useReducer(sectionReducer, initialState);
// dispatch({ type: 'UPDATE_SECTION', section: 'couple', data: { groomName: 'Andi' } })
```

**Form validation:** helper functions that return error objects. No schema library dependency (keep it lean). Validation messages are in Indonesian.

## What We Do NOT Use

- Redux / Zustand / Jotai — unnecessary for this project's complexity.
- URL state management — Next.js search params are sufficient for filter/sort state.
- WebSocket state — RSVP updates are polled on interval (SWR revalidate) rather than pushed.

## Cross-module Communication

When one module's data depends on another (e.g., guest list page needs event ID from the selected event), the event ID is passed as a prop or query parameter. No global event bus or pub/sub pattern is used.