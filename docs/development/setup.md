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

Minimum required variables:

```
DATABASE_URL=postgres://postgres:password@localhost:5432/owndangan
JWT_SECRET=<random-64-char-string>
MIDTRANS_SERVER_KEY=<sandbox-server-key>
MIDTRANS_CLIENT_KEY=<sandbox-client-key>
MIDTRANS_IS_PRODUCTION=false
```

### Frontend

```bash
cp frontend/.env.example frontend/.env.local
```

Minimum required variables:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=<sandbox-client-key>
```

## Database Setup

### Create Database

```bash
createdb owndangan
createdb owndangan_test
```

Or via Docker:

```bash
docker run -d --name owndangan-db -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:16
```

### Run Migrations

```bash
cd backend
goose postgres "$DATABASE_URL" up
```

### Seed Data (Optional)

```bash
TODO: go run cmd/seed/main.go
```

## Start Development Servers

### Backend (with hot reload)

```bash
cd backend
air
```

Without hot reload:

```bash
go run cmd/server/main.go
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
- **Database connection refused**: Ensure PostgreSQL is running and `DATABASE_URL` is correct.
- **Missing migrations**: Run `goose up` before starting the backend.