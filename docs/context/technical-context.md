# Technical Context

## Technology Stack

### Backend

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | > TODO: Define exact version |
| Framework | Standard library + chi/gorilla-mux | > TODO: Define router |
| ORM | GORM | > TODO: Define version |
| Auth | JWT (golang-jwt) | > TODO: Define version |
| Validation | go-playground/validator | > TODO: Define version |
| Config | viper / envconfig | > TODO: Define |
| Logging | zerolog / zap | > TODO: Define |
| Testing | Go testing + testify | > TODO: Define |
| Database | PostgreSQL | 16+ |

### Frontend

| Component | Technology | Version |
|-----------|-----------|---------|
| Framework | Next.js | > TODO: Define version |
| Router | App Router | Built-in |
| Language | TypeScript | > TODO: Define version |
| Styling | Tailwind CSS | > TODO: Define version |
| State | React context / SWR / TanStack Query | > TODO: Define |
| HTTP Client | fetch / axios | > TODO: Define |
| Form | react-hook-form | > TODO: Define |
| Validation | zod | > TODO: Define |

### Infrastructure

| Component | Technology | Version |
|-----------|-----------|---------|
| Database | PostgreSQL | 16+ |
| Cache | Redis (future) | > TODO |
| Object Storage | S3-compatible (MinIO / AWS S3) | > TODO |
| Reverse Proxy | Nginx / Caddy / Cloudflare | > TODO |
| Container | Docker | > TODO |

## Dependency Policy

- Pin major dependency versions in `go.mod` and `package.json`.
- Use Go modules for Go dependency management.
- Use npm/pnpm for frontend dependency management.
- Avoid unnecessary dependencies. If a library is >50KB and does one small thing, consider implementing it.
- Document major library choices in ADR.

## Environment Configuration

Configuration via environment variables (`.env` files for local development).

```
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/wedding_invitation

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24

# Midtrans
MIDTRANS_SERVER_KEY=your-server-key
MIDTRANS_CLIENT_KEY=your-client-key
MIDTRANS_IS_PRODUCTION=false

# App
APP_PORT=8080
APP_ENV=development

# Storage (future)
STORAGE_PROVIDER=local
STORAGE_PATH=./uploads

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=your-client-key
```

> TODO: Finalize env variable names and defaults.

## Local Development

> TODO: Define after initial repo setup.

Prerequisites:
- Go 1.x
- Node.js 18+
- PostgreSQL 16+
- Docker (optional, for PostgreSQL)

## Production Architecture

```
User → CDN (Cloudflare) → Reverse Proxy → Frontend (Next.js)
                                       → Backend (Go API)
                                       → PostgreSQL
                                       → Midtrans API
                                       → Object Storage (S3)
```

## API Communication

- Frontend ↔ Backend: REST JSON over HTTPS.
- Backend ↔ Midtrans: HTTPS with server key authentication.
- Midtrans → Backend: Webhook (HTTP POST) with signature verification.
- Frontend → Midtrans: Snap.js (client-side, using client key).

## Authentication Flow

```
User credentials → POST /api/v1/auth/login → JWT token (access token)
                                                       ↓
                                        Token included in Authorization: Bearer header
                                                       ↓
                                        Middleware validates JWT → user context
```

- JWT payload: `{ user_id, email, role }`
- Token expiry: configurable (default 24 hours)
- Password hashing: bcrypt (cost 12)

## Error Handling

- API errors follow standard envelope: `{ "error": { "code": "ERROR_CODE", "message": "..." } }`
- HTTP status codes follow REST conventions.
- Unexpected errors are logged server-side and return generic 500.
- Validation errors include details field.

## Logging

- Structured JSON logging in production.
- Development: human-readable format.
- Include: request ID, timestamp, user ID, action, duration.
- Never log: password hashes, tokens, payment full card data (Midtrans handles this).

## Testing

- Backend: unit tests (service + repository + handler), integration tests with test database.
- Frontend: component tests (Storybook/testing-library), API integration tests.
- E2E: Playwright/Cypress (future).
- Test database: separate PostgreSQL database for test runs.

> TODO: Define test runner, database test helpers, and CI pipeline.
