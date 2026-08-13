# Backend Deployment

## Build Artifact

The backend is compiled into a single static binary with embedded migration files.

```bash
cd backend

# Build for production (Linux amd64)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server ./cmd/server

# Verify the binary
file bin/server
# Output: ELF 64-bit LSB executable, x86-64, statically linked
```

The binary includes:
- All Go source code compiled.
- Embedded SQL migration files (via `embed` package).
- No external runtime dependencies (no Go installation needed).

## Docker Container

```dockerfile
# Dockerfile (backend)
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /server /server
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["/server"]
```

### Build and Run

```bash
docker build -t owndangan-backend:latest -f Dockerfile .
docker run -d --name owndangan-api \
  -p 8080:8080 \
  --env-file .env.production \
  owndangan-backend:latest
```

## Environment Configuration

Environment variables are injected at runtime (not build time):

```bash
docker run -e DATABASE_URL="postgres://..." -e JWT_SECRET="..." owndangan-backend:latest
```

For production:
- Use Docker secrets or a secret manager (e.g., AWS Secrets Manager, Vault).
- Never bake secrets into the Docker image.

## Health Check

The backend exposes `GET /health` which returns:

```json
{
  "status": "ok",
  "timestamp": "2025-01-15T10:00:00Z",
  "version": "1.2.3",
  "database": "connected",
  "uptime_seconds": 3600
}
```

- Used by the container orchestrator for liveness/readiness probes.
- Returns `503` if the database is unreachable.
- Returns `200` with no auth required.

## Deployment Options

### Option 1: VPS / Bare Metal

```bash
# Using systemd service
[Unit]
Description=owndangan-backend
After=network.target

[Service]
ExecStart=/opt/owndangan/server
WorkingDirectory=/opt/owndangan
EnvironmentFile=/opt/owndangan/.env
Restart=always
User=owndangan

[Install]
WantedBy=multi-user.target
```

### Option 2: Container Orchestrator (Kubernetes, Nomad)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: owndangan-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: owndangan-backend
  template:
    spec:
      containers:
      - name: api
        image: owndangan-backend:latest
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: backend-env
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
```

### Option 3: Platform as a Service (Railway, Fly.io)

```bash
# Railway CLI
railway up

# Fly.io
flyctl deploy
```

## Post-Deployment

1. Verify health endpoint returns 200.
2. Run smoke tests against the deployed endpoint.
3. Monitor error rates for 15 minutes.
4. Check database connection pool metrics.