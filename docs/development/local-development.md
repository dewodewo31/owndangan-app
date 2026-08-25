# Local Development Workflow

## Hot Reload

### Backend (Air)

[Air](https://github.com/air-verse/air) watches for `.go` file changes and rebuilds/restarts the server automatically.

```bash
# Install Air (once)
go install github.com/air-verse/air@latest

# Start with hot reload
cd backend
air
```

Air configuration is in `backend/.air.toml`. It excludes `tmp/`, `vendor/`, and test files from the watch list.

### Frontend (Next.js)

```bash
cd frontend
npm run dev
```

Next.js dev server provides:
- Fast Refresh (preserves component state on edits).
- Fast compilation with SWC.
- Error overlay for build/runtime errors.

## Database Management

### Local PostgreSQL

```bash
# Start PostgreSQL (if installed via package manager)
sudo systemctl start postgresql

# Or via Docker
docker start owndangan-db
```

### Run Migrations

Migrations and default package seeds run automatically when the backend starts (`db.AutoMigrate()` + `SeedPackages` in `cmd/server/main.go`). Just start the server:

```bash
cd backend
go run ./cmd/server
```

> `make migrate-up` currently points at a missing `cmd/migrate` binary and will fail; rely on the server's auto-migrate.

## Running Tests

### Backend

```bash
# All tests
cd backend && go test ./...

# With race detection
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test ./internal/service/...

# Integration tests (requires test database)
go test -tags=integration ./internal/repository/...
```

### Frontend

```bash
cd frontend

# All tests
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage
```

## Linting

### Backend

```bash
# Install golangci-lint (once)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
cd backend && golangci-lint run ./...

# Fix auto-fixable issues
golangci-lint run --fix ./...
```

### Frontend

```bash
cd frontend

# ESLint
npm run lint

# TypeScript check
npm run typecheck

# Format
npm run format
```

## Pre-commit Hooks

We use [pre-commit](https://pre-commit.com) or Lefthook:

```yaml
# .lefthook.yml
pre-commit:
  parallel: true
  commands:
    go-lint:
      run: golangci-lint run ./...
    ts-lint:
      run: npm run lint --prefix frontend
    go-test:
      run: go test ./internal/... -short
```

## Workflow Summary

```bash
# Start everything
docker start owndangan-db  # If using Docker
cd backend && air &         # Backend hot reload
cd frontend && npm run dev  # Frontend dev server

# In another terminal, run tests
cd backend && go test ./...
cd frontend && npm test
```