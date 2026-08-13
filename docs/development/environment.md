# Environment Configuration

## Overview

The application uses environment variables for all configuration. No secrets or environment-specific values are hardcoded.

## Backend Environment Variables

### `backend/.env`

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8080` | HTTP server port |
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `DATABASE_MAX_OPEN` | No | `25` | Max open DB connections |
| `DATABASE_MAX_IDLE` | No | `5` | Max idle DB connections |
| `DATABASE_CONN_MAX_LIFETIME` | No | `5m` | Max connection lifetime |
| `JWT_SECRET` | Yes | — | HMAC key for signing JWT tokens |
| `JWT_ACCESS_TTL` | No | `15m` | Access token expiry |
| `JWT_REFRESH_TTL` | No | `7d` | Refresh token expiry |
| `MIDTRANS_SERVER_KEY` | Yes | — | Midtrans server key (sandbox or production) |
| `MIDTRANS_CLIENT_KEY` | Yes | — | Midtrans client key (sandbox or production) |
| `MIDTRANS_IS_PRODUCTION` | No | `false` | Toggle sandbox/production mode |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000` | Comma-separated allowed origins |
| `LOG_LEVEL` | No | `info` | Log level: debug, info, warn, error |
| `SENTRY_DSN` | No | — | Sentry error tracking DSN |

### Environment-Specific Overrides

Create `backend/.env.development`, `backend/.env.staging`, `backend/.env.production` as needed. The deployment pipeline loads the appropriate file.

## Frontend Environment Variables

### `frontend/.env.local`

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | Yes | `http://localhost:8080/api` | Backend API base URL |
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
DATABASE_URL=postgres://user:password@localhost:5432/owndangan
JWT_SECRET=change-me-to-a-random-64-char-string
MIDTRANS_SERVER_KEY=SB-Mid-server-xxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxx
```