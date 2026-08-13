# API Client

## Overview

The API client is a thin fetch wrapper located at `src/lib/api-client.ts`. It provides a single interface for all HTTP calls to the Go backend, handling authentication header injection, error parsing, and token refresh transparently.

## Implementation

```tsx
// src/lib/api-client.ts

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

interface ApiResponse<T> {
  success: boolean;
  data: T;
  meta?: {
    pagination?: PaginationMeta;
    request_id: string;
  };
  error?: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}
```

## Methods

The client exposes four HTTP methods matching the backend API conventions:

```tsx
apiClient.get<T>(path: string, params?: Record<string, string>): Promise<T>
apiClient.post<T>(path: string, body?: unknown): Promise<T>
apiClient.put<T>(path: string, body?: unknown): Promise<T>
apiClient.delete<T>(path: string): Promise<T>
```

Each method:

1. Constructs the full URL: `${BASE_URL}${path}`.
2. Appends query params for `GET` requests.
3. Sets `Content-Type: application/json`.
4. Injects `Authorization: Bearer <token>` if a token is available (read from cookie or localStorage).
5. Calls `fetch()`.
6. On non-2xx response, parses the error body and throws an `ApiError` with `code`, `message`, and `status`.
7. On 401 response, attempts token refresh via `POST /api/v1/auth/refresh`, then retries the original request once.
8. On success, unwraps `response.data` and returns it.

## Auth Header Injection

```tsx
function getAuthHeaders(): Record<string, string> {
  const token = getAccessToken(); // reads from cookie or localStorage
  if (!token) return {};
  return { Authorization: `Bearer ${token}` };
}
```

The client reads the token from the same location where the auth provider stores it. This must be consistent — either httpOnly cookie (recommended, more secure) or localStorage (simpler, but XSS-vulnerable). The chosen approach uses **httpOnly cookies** for the access token, so the client does not need to manage token storage manually. The cookie is set by the backend on login/refresh and is automatically included in cross-origin requests if `credentials: 'include'` is set.

If httpOnly cookies are not feasible during development, the client falls back to `Authorization` header from localStorage.

## Error Handling

```tsx
class ApiError extends Error {
  code: string;
  status: number;
  details?: Record<string, string>;

  constructor(code: string, message: string, status: number, details?: Record<string, string>) {
    super(message);
    this.code = code;
    this.status = status;
    this.details = details;
  }
}
```

Components handle errors by catching `ApiError`:

```tsx
try {
  const events = await apiClient.get('/events');
} catch (error) {
  if (error instanceof ApiError && error.code === 'LIMIT_EXCEEDED') {
    toast.error('Batas undangan Anda sudah tercapai. Upgrade paket untuk menambah undangan.');
  }
}
```

## Usage in Components

**Server components** fetch directly:

```tsx
// app/[slug]/page.tsx (server component)
async function getInvitation(slug: string) {
  const res = await fetch(`${BASE_URL}/public/events/${slug}`, { next: { revalidate: 60 } });
  if (!res.ok) notFound();
  return res.json();
}
```

**Client components** use the wrapper via hooks:

```tsx
// In a hook
import useSWR from 'swr';
import { apiClient } from '@/lib/api-client';

export function useGuests(eventId: string) {
  return useSWR(['/events', eventId, 'guests'], () =>
    apiClient.get(`/events/${eventId}/guests`)
  );
}
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080/api/v1` | Backend base URL |
| `NEXT_PUBLIC_MIDTRANS_CLIENT_KEY` | — | Midtrans Snap.js client key |
| `API_INTERNAL_URL` | — | Internal URL for server-side fetch (SSR) |

The `NEXT_PUBLIC_API_URL` is used for client-side calls. For server-side fetches (SSR, SSG, ISR), `API_INTERNAL_URL` is preferred to avoid hitting the public load balancer.

## File Uploads

File uploads use a separate `apiClient.upload()` method that sends `multipart/form-data`:

```tsx
apiClient.upload('/events/:id/gallery', formData);
```

The upload method includes the same auth headers but omits `Content-Type` (let the browser set the boundary).