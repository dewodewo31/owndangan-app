# Environment Configuration

## Overview

The Owndangan platform supports three environments:

| Environment | Purpose | URL |
|-------------|---------|-----|
| Development | Local development | localhost |
| Staging | Pre-production testing | staging.owndangan.com |
| Production | Live environment | owndangan.com |

## Environment Variables

### Backend

| Variable | Development | Staging | Production |
|----------|-------------|---------|------------|
| ENV | development | staging | production |
| PORT | 8080 | 8080 | 8080 |
| DB_HOST | localhost | postgres | postgres |
| DB_PORT | 5433 | 5432 | 5432 |
| DB_SSLMODE | disable | require | require |
| JWT_SECRET | dev-secret | *staging_secret* | *production_secret* |
| MIDTRANS_IS_PRODUCTION | false | false | true |

### Frontend

| Variable | Development | Staging | Production |
|----------|-------------|---------|------------|
| NEXT_PUBLIC_API_URL | http://localhost:8080/api/v1 | https://staging.owndangan.com/api/v1 | https://owndangan.com/api/v1 |

## Secrets Management

### GitHub Secrets (CI/CD)

| Secret | Description |
|--------|-------------|
| GITHUB_TOKEN | Container registry authentication |
| STAGING_DB_PASSWORD | Staging database password |
| PRODUCTION_DB_PASSWORD | Production database password |
| STAGING_JWT_SECRET | Staging JWT signing secret |
| PRODUCTION_JWT_SECRET | Production JWT signing secret |
| MIDTRANS_SERVER_KEY | Midtrans API server key |
| MIDTRANS_CLIENT_KEY | Midtrans API client key |

### Environment Secrets (GitHub Environments)

Configure in GitHub → Settings → Environments:

**staging:**
- `DB_HOST`: Staging database host
- `DB_PASSWORD`: Staging database password
- `JWT_SECRET`: Staging JWT secret

**production:**
- `DB_HOST`: Production database host
- `DB_PASSWORD`: Production database password
- `JWT_SECRET`: Production JWT secret

## Deployment Process

### Staging Deployment

1. Push to `main` branch
2. CI runs tests and builds
3. CD automatically deploys to staging
4. Smoke tests verify deployment

### Production Deployment

1. Create a release tag
2. Manually trigger CD workflow
3. Select `production` environment
4. CD deploys to production
5. Smoke tests verify deployment

### Rollback

If deployment fails:

```bash
# Rollback to previous version
./scripts/rollback.sh production

# Rollback to specific version
./scripts/rollback.sh production v1.2.3
```

## Infrastructure

### Docker Compose (Production)

```
nginx (reverse proxy + SSL)
├── frontend (Next.js)
└── backend (Go API)
    └── postgres (database)
```

### Kubernetes (Alternative)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: owndangan-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: ghcr.io/owndangan/backend:latest
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: backend-secrets
```

## Monitoring

- Health check: `/health`
- Metrics: Prometheus (optional)
- Logs: Structured JSON logging
- Alerts: Configure for 5xx errors

## Security

- All secrets stored in GitHub Secrets
- Environment-specific secrets isolated
- Production requires manual approval
- Database SSL required in staging/production
- JWT secrets are environment-specific
