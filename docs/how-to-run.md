# How to Run Owndangan

Best-practice instructions to get the full stack running locally. Two supported paths:

- **Native dev (recommended)** — Docker runs only PostgreSQL; backend and frontend run on your machine for instant hot reload.
- **Full Docker stack** — `docker compose` builds and runs everything (postgres, backend, frontend, optional ngrok tunnels).

Both assume a fresh clone. Commands are run from the repo root unless noted.

## Prerequisites

| Tool | Version | Used for |
|------|---------|----------|
| Go | 1.22+ (Dockerfile builds with 1.26) | Backend |
| Node.js | 20 LTS+ | Frontend |
| Docker | latest | Local PostgreSQL, or the full stack |
| npm | bundled with Node | Frontend deps |

## Path A — Native dev (recommended)

### 1. Start PostgreSQL

The repo convention is PostgreSQL on **port 5433** (see `backend/.env` and `AGENTS.md`). Use the provided compose file so the container name and port match what the code expects:

```bash
docker compose -f docker/docker-compose.yml up -d postgres
```

Verify it is healthy:

```bash
docker exec owndangan-db psql -U postgres -d owndangan -c "SELECT 1;"
```

### 2. Backend

```bash
cd backend

# Install deps (first time, or after go.mod changes)
go mod download

# Ensure env exists; edit values if needed
cp .env.example .env    # skip if backend/.env already present

# Run — auto-migrates the schema AND seeds default packages on boot
go run ./cmd/server
```

The server listens on `http://localhost:8080`. On startup it runs `db.AutoMigrate()` and `SeedPackages`, so **no separate migration step is required**.

> `make migrate-up` is currently a no-op/broken (it points at a `cmd/migrate` binary that does not exist). Do not use it; the server handles migrations on start.

### 3. Frontend

```bash
cd frontend

# Install deps (first time, or after package.json changes)
npm install

# Ensure env exists; edit values if needed
cp .env.example .env.local    # skip if frontend/.env.local already present

# Start dev server with Fast Refresh
npm run dev
```

The frontend listens on `http://localhost:3000` and calls the backend at `http://localhost:8080/api/v1` (set via `NEXT_PUBLIC_API_URL` in `.env.local`).

### 4. Verify

```bash
curl -i http://localhost:8080/health          # -> 200 OK
curl -i http://localhost:8080/api/v1/invitations  # -> 401 (auth required = server is up)
# Browser: http://localhost:3000  -> landing/login loads
```

## Path B — Full Docker stack

From the `docker/` directory. This builds the backend and frontend images and wires them together (plus optional ngrok tunnels for public URLs).

```bash
cd docker
cp .env.example .env    # fill NGROK_AUTHTOKEN only if you want public tunnels
docker compose up --build
```

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- PostgreSQL runs inside the `postgres` service on port `5433` (host mapping from `.env` `DB_PORT`, default `5433`).

The compose frontend is configured to proxy `/api/v1/*` and `/uploads/*` to the backend at runtime (via `API_BASE_URL`), so it works same-origin without CORS tuning. Frontend env (`NEXT_PUBLIC_API_URL`, Midtrans keys) is **baked at build time** — change it in `docker/.env` and rebuild.

Stop and clean up:

```bash
docker compose down          # keep data volume
docker compose down -v       # also drop the postgres volume
```

## Environment variables (minimum to run)

### Backend — `backend/.env`

Copy from `backend/.env.example`. The values below are what the repo ships with and are sufficient to boot locally:

```
APP_ENV=development
PORT=8080

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=owndangan
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

Midtrans keys are optional for local runs (payments will not complete without them, but the app boots). Leave `MIDTRANS_*` empty unless you have sandbox keys.

### Frontend — `frontend/.env.local`

Copy from `frontend/.env.example`:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=
NEXT_PUBLIC_MIDTRANS_IS_PRODUCTION=false
```

## Database

- Engine: PostgreSQL 16.
- Local port: **5433** (matches `backend/.env` and `AGENTS.md`).
- Migrations + package seeding happen automatically when the backend starts. You do not run a separate migration tool.
- For tests, a separate database `owndangan_test` on the same host/port is used (see "Running tests").

If you need a clean Postgres without the compose file:

```bash
docker run -d --name owndangan-db -e POSTGRES_PASSWORD=password -p 5433:5432 postgres:16
```

## Running tests

### Backend

```bash
cd backend

make test          # unit + integration + API (requires DB on 5433)
make test-unit     # unit only, no DB needed
```

Integration and API tests need a test database reachable on `localhost:5433` named `owndangan_test` (user `postgres`, password `password`). Create it once:

```bash
docker exec -it owndangan-db psql -U postgres -c "CREATE DATABASE owndangan_test;"
```

### Frontend

```bash
cd frontend
npm test           # Jest
npm run lint       # ESLint (next lint)
npm run build      # production build / type check
```

## Common issues

- **Backend fails to connect to DB** — confirm the `postgres` container is healthy and `DB_PORT` in `backend/.env` matches the exposed port (default `5433`).
- **401 on `/api/v1/...` is expected** for unauthenticated calls; it means the server is up.
- **Frontend can't reach backend** — check `NEXT_PUBLIC_API_URL` in `frontend/.env.local` equals `http://localhost:8080/api/v1` (note the `/v1`).
- **Port in use** — change `PORT` in `backend/.env` / `FRONTEND_PORT` in `docker/.env`, or stop the conflicting process.
- **Full-stack Docker: CORS / wrong API host** — frontend API URL is baked at build time; edit `docker/.env` (`NEXT_PUBLIC_API_URL`, `API_BASE_URL`) and `docker compose up --build` again.
