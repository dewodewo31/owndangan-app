# Development Setup

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.22+ | Backend runtime |
| Node.js | 20 LTS+ | Frontend runtime |
| PostgreSQL | 16+ | Primary database |
| Docker (optional) | Latest | Local PostgreSQL, CI |
| Air (optional) | Latest | Go hot reload |
| goose | Latest | Database migrations |

## Clone the Repository

```bash
git clone <repository-url> owndangan
cd owndangan
```

## Install Dependencies

### Backend

```bash
cd backend
go mod download
go mod tidy
```

### Frontend

```bash
cd frontend
npm install
```

## Configure Environment Variables

### Backend

```bash
cp backend/.env.example backend/.env
# Edit backend/.env with your values
```

The backend reads individual `DB_*` variables (not a single `DATABASE_URL`). Minimum required:

```
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=owndangan
JWT_SECRET=<random-64-char-string>
MIDTRANS_SERVER_KEY=<sandbox-server-key>   # optional for local boot
MIDTRANS_CLIENT_KEY=<sandbox-client-key>   # optional for local boot
MIDTRANS_IS_PRODUCTION=false
```

### Frontend

```bash
cp frontend/.env.example frontend/.env.local
```

Minimum required variables:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=<sandbox-client-key>
```

## Database Setup

### Create Database

```bash
createdb owndangan
createdb owndangan_test
```

Or via Docker (the repo convention uses port 5433, matching `backend/.env`):

```bash
docker run -d --name owndangan-db -e POSTGRES_PASSWORD=password -p 5433:5432 postgres:16
```

### Migrations & Seeds

Migrations and default package seeding run **automatically when the backend starts** (`db.AutoMigrate()` + `SeedPackages` in `cmd/server/main.go`). No separate migration command is required — just start the server (see below).

> Note: `make migrate-up` currently points at a `cmd/migrate` binary that does not exist and will fail. Use the server's auto-migrate instead.

## Start Development Servers

### Backend (with hot reload)

```bash
cd backend
air
```

Without hot reload:

```bash
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm run dev
```

## Verify Setup

1. Backend health check: `curl http://localhost:8080/health` → `200 OK`
2. Frontend: `http://localhost:3000` → Login page loads
3. API: `curl http://localhost:8080/api/v1/invitations` → 401 (auth required, meaning the server is running)

## Common Issues

- **Port already in use**: Kill existing process or change `PORT` in `.env`.
- **Database connection refused**: Ensure PostgreSQL is running and the `DB_*` values in `backend/.env` (host/port/user/password/name) are correct. The repo convention is PostgreSQL on port `5433`.
- **Missing migrations**: Not applicable — the backend auto-migrates the schema on startup. Just start the server.