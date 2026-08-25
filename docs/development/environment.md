# Environment Configuration

## Overview

The application uses environment variables for all configuration. No secrets or environment-specific values are hardcoded.

## Backend Environment Variables

### `backend/.env`

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_ENV` | No | `development` | Runtime environment |
| `PORT` | No | `8080` | HTTP server port |
| `DB_HOST` | Yes | `localhost` | PostgreSQL host |
| `DB_PORT` | Yes | `5433` | PostgreSQL port (repo convention: 5433) |
| `DB_USER` | Yes | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | Yes | `password` | PostgreSQL password |
| `DB_NAME` | Yes | `owndangan` | Database name |
| `DB_SSLMODE` | No | `disable` | SSL mode |
| `DB_MAX_OPEN_CONNS` | No | `25` | Max open DB connections |
| `DB_MAX_IDLE_CONNS` | No | `5` | Max idle DB connections |
| `DB_LOG_SQL` | No | `false` | Log SQL queries |
| `JWT_SECRET` | Yes | — | HMAC key for signing JWT tokens |
| `JWT_ACCESS_EXPIRY` | No | `15m` | Access token expiry |
| `JWT_REFRESH_EXPIRY` | No | `168h` | Refresh token expiry |
| `MIDTRANS_SERVER_KEY` | No* | — | Midtrans server key (sandbox or production) |
| `MIDTRANS_CLIENT_KEY` | No* | — | Midtrans client key (sandbox or production) |
| `MIDTRANS_IS_PRODUCTION` | No | `false` | Toggle sandbox/production mode |
| `STORAGE_PROVIDER` | No | `local` | Object storage provider (`local` / s3) |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000,http://localhost:8080` | Comma-separated allowed origins |

\* Required only when using payments; the app boots without them.

The backend reads individual `DB_*` variables — there is **no** `DATABASE_URL` connection string. Migrations and package seeds run automatically on server start, so no `goose`/migration tool is needed.

### Environment-Specific Overrides

Create `backend/.env.development`, `backend/.env.staging`, `backend/.env.production` as needed. The deployment pipeline loads the appropriate file.

## Frontend Environment Variables

### `frontend/.env.local`

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | Yes | `http://localhost:8080/api/v1` | Backend API base URL |
| `NEXT_PUBLIC_MIDTRANS_CLIENT_KEY` | Yes | — | Midtrans client key (public) |
| `NEXT_PUBLIC_SENTRY_DSN` | No | — | Sentry DSN for frontend error tracking |
| `NEXT_PUBLIC_GA_ID` | No | — | Google Analytics ID |

### Variable Prefix Convention

- `NEXT_PUBLIC_*` — Exposed to the browser. Safe to include in client bundle.
- All other variables — Server-side only (available in API routes, `getServerSideProps`).

## Development vs Production

| Aspect | Development | Production |
|--------|-------------|------------|
| Midtrans | Sandbox | Production |
| Log level | `debug` | `warn` |
| JWT TTL | `15m` / `7d` | `15m` / `7d` |
| CORS | `http://localhost:3000` | Production domain |
| DB SSL | `disable` | `require` |
| Hot reload | Yes | No |

## Secrets Management

- **Never commit** `.env` files (except `.env.example`).
- **Production secrets**: Stored in environment variables on the deployment platform (Vercel, Railway, Docker secrets).
- **Rotation**: JWT secret and Midtrans keys are rotated every 90 days or immediately if compromised.
- **Access**: Only senior developers and DevOps have access to production secrets.

## `.env.example` Files

Each project contains a `.env.example` with placeholder values:

```bash
# backend/.env.example
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=owndangan
JWT_SECRET=change-me-to-a-random-64-char-string
MIDTRANS_SERVER_KEY=SB-Mid-server-xxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxx
```