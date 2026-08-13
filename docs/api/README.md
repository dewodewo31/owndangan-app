# API Documentation

Base path: `/api/v1`

## Standard Conventions

All endpoints follow the response conventions defined in [conventions.md](./conventions.md).

- **Success envelope:** `{ "success": true, "data": {}, "meta": {} }`
- **Error envelope:** `{ "success": false, "error": { "code": "...", "message": "..." } }`
- **Auth header:** `Authorization: Bearer <jwt_token>` (for authenticated endpoints)
- **Content-Type:** `application/json` (request and response)

## Endpoint Index

### Public (No Auth)

| Method | Endpoint | Description | Doc |
|--------|----------|-------------|-----|
| GET | `/api/v1/health` | Health check | — |
| GET | `/api/v1/packages` | List available packages | [packages.md](./packages.md) |
| GET | `/api/v1/e/:slug` | Get full public invitation | [public-invitation.md](./public-invitation.md) |
| POST | `/api/v1/payments/webhook` | Midtrans payment notification (signature verified) | [payments.md](./payments.md) |
| POST | `/api/v1/auth/register` | Create new account | [authentication.md](./authentication.md) |
| POST | `/api/v1/auth/login` | Sign in | [authentication.md](./authentication.md) |
| POST | `/api/v1/auth/refresh` | Renew access token | [authentication.md](./authentication.md) |
| POST | `/api/v1/auth/logout` | Invalidate session | [authentication.md](./authentication.md) |

### Protected (JWT Required)

| Method | Endpoint | Description | Doc |
|--------|----------|-------------|-----|
| POST | `/api/v1/auth/change-password` | Change password | [authentication.md](./authentication.md) |
| GET | `/api/v1/users/me` | Get profile | [users.md](./users.md) |
| PUT | `/api/v1/users/me` | Update profile | [users.md](./users.md) |
| GET | `/api/v1/users/me/subscription` | Get current subscription | [subscriptions.md](./subscriptions.md) |
| GET | `/api/v1/packages/all` | List all packages (admin) | [packages.md](./packages.md) |
| POST | `/api/v1/packages` | Create package (admin) | [packages.md](./packages.md) |
| PUT | `/api/v1/packages/:id` | Update package (admin) | [packages.md](./packages.md) |
| DELETE | `/api/v1/packages/:id` | Deactivate package (admin) | [packages.md](./packages.md) |
| GET | `/api/v1/events` | List user's events | [events.md](./events.md) |
| POST | `/api/v1/events` | Create event | [events.md](./events.md) |
| GET | `/api/v1/events/:id` | Get event details | [events.md](./events.md) |
| PUT | `/api/v1/events/:id` | Update event | [events.md](./events.md) |
| DELETE | `/api/v1/events/:id` | Delete event | [events.md](./events.md) |
| POST | `/api/v1/events/:id/publish` | Publish invitation | [events.md](./events.md) |
| POST | `/api/v1/events/:id/unpublish` | Unpublish invitation | [events.md](./events.md) |
| POST | `/api/v1/payments/snap` | Create Midtrans Snap transaction | [payments.md](./payments.md) |
| GET | `/api/v1/payments/transactions` | List user transactions | [payments.md](./payments.md) |
| GET | `/api/v1/subscriptions/current` | Get active subscription | [subscriptions.md](./subscriptions.md) |
| GET | `/api/v1/subscriptions/default` | Get or default subscription | [subscriptions.md](./subscriptions.md) |

## Related Documents

- [API Conventions](./conventions.md) — Standard formats, pagination, error codes
- [Database Schema](/run/media/oweedd/New Volume/app-owndangan/docs/database/schema.md) — Data model reference
- [Midtrans Integration](/run/media/oweedd/New Volume/app-owndangan/docs/integrations/midtrans.md) — Payment gateway specifics
