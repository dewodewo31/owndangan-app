# Deployment Architecture

## Environment Strategy

| Environment | Purpose | Backend | Frontend | Database |
|-------------|---------|---------|----------|----------|
| Development | Local coding | `go run` / air | `next dev` | Local PostgreSQL |
| Staging | Pre-production testing | Docker container | Vercel / Docker | Staging PostgreSQL |
| Production | Live | Docker container / K8s | Vercel / Node server | Production PostgreSQL |

## Production Architecture

```
Internet
    │
    ▼
Cloudflare (CDN + DNS + DDoS protection)
    │
    ├── → Next.js Frontend (Vercel / Node server)
    │       └── Static assets served via CDN
    │
    └── → Go Backend (Docker container / K8s)
            │
            ├── PostgreSQL (Managed DB / RDS)
            ├── Object Storage (S3-compatible)
            └── Midtrans API (external)
```

## Backend Deployment

### Option 1: Docker on VPS

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY .env.example .env
EXPOSE 8080
CMD ["./server"]
```

### Option 2: Kubernetes

> TODO: Define K8s manifests.

## Frontend Deployment

### Option 1: Vercel

- Connect GitHub repo → Vercel auto-deploys
- Environment variables configured in Vercel dashboard
- Serverless functions for API routes (if any)

### Option 2: Docker

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./
EXPOSE 3000
CMD ["npm", "start"]
```

## Environment Variables

### Backend

```
APP_PORT=8080
APP_ENV=production
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
JWT_SECRET=...
JWT_EXPIRATION_HOURS=24
MIDTRANS_SERVER_KEY=...
MIDTRANS_CLIENT_KEY=...
MIDTRANS_IS_PRODUCTION=true
STORAGE_PROVIDER=s3
STORAGE_ENDPOINT=https://s3.amazonaws.com
STORAGE_BUCKET=wedding-invitation-assets
```

### Frontend

```
NEXT_PUBLIC_API_URL=https://api.example.com/api/v1
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=...
```

## Database Deployment

- Managed PostgreSQL (AWS RDS / DigitalOcean Managed DB / Supabase)
- Automated backups (daily with 7-day retention)
- Connection pooling via PgBouncer if needed
- SSL/TLS connection required

## CI/CD Pipeline

> TODO: Define CI/CD approach (GitHub Actions / GitLab CI).

## Monitoring

> TODO: Define monitoring tools (Sentry, Grafana, Prometheus).

## Related Documentation

- `docs/deployment/environments.md`
- `docs/deployment/backend-deployment.md`
- `docs/deployment/frontend-deployment.md`
- `docs/deployment/database-deployment.md`
- `docs/deployment/monitoring.md`