# AGENTS.md

## Project Context

This repository contains a SaaS wedding invitation platform (Indonesia market) with decoupled Go backend and Next.js frontend.

## Before Coding

1. Read the relevant documentation in `/docs`.
2. Inspect existing implementation.
3. Identify affected modules.
4. Check API/database contracts.
5. Avoid unnecessary refactoring.

## Architecture Rules

- Keep handlers thin (parse request, call service, return response).
- Business logic belongs in services.
- Database access belongs in repositories.
- Validate external input.
- Enforce authorization server-side.
- Do not trust frontend payment status.
- Never expose secrets.
- Do not bypass existing abstractions.

## Change Rules

- Make the smallest safe change.
- Preserve existing behavior unless requirement says otherwise.
- Add tests for new behavior.
- Update documentation when contracts change.

## Forbidden

- Hardcoded secrets.
- Fake API responses.
- Silent schema changes.
- Removing tests to make them pass.
- Bypassing authorization.
- Activating subscription from frontend callback.
- Large unrelated refactors.

## Completion Checklist

Before considering a task complete:

- [ ] Code implemented
- [ ] Tests added/updated
- [ ] Tests pass
- [ ] Lint passes
- [ ] API contract updated if needed
- [ ] Database documentation updated if needed
- [ ] Module documentation updated if needed
- [ ] No unrelated changes

## Commands

### Backend
```sh
cd backend

# Run migrations
make migrate-up      # or: go run ./cmd/server (auto-migrates on start)

# Start server (auto-migrates + seeds packages)
go run ./cmd/server

# Build
go build -o bin/server ./cmd/server

# Run tests
make test            # unit + integration + API
make test-unit       # unit tests only
make test-integration

# Lint
make lint            # requires golangci-lint
```

### Frontend
```sh
cd frontend

# Dev server
npm run dev

# Build
npm run build

# Lint
npm run lint

# Test
npm run test
```

### Database (local PostgreSQL on port 5433 via Docker)
```sh
docker exec owndangan-db psql -U postgres -d owndangan -c "SELECT 1;"
```
