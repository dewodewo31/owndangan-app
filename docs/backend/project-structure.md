# Project Structure

## Directory Layout

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point — config init, DI wiring, server start
│
├── internal/
│   ├── api/
│   │   ├── router.go            # Route registration via chi router
│   │   ├── handler/             # One file per module (auth_handler.go, event_handler.go, ...)
│   │   ├── middleware/          # Auth, logger, CORS, rate limiter, recovery, request ID
│   │   └── dto/                 # Request/response structs — one file per module
│   │
│   ├── service/                 # Business logic — one file per module (auth_service.go, event_service.go)
│   │
│   ├── repository/             # GORM database access — one file per model (user_repo.go, event_repo.go)
│   │
│   ├── model/                  # GORM model definitions — one file per model (user.go, event.go)
│   │
│   ├── config/                 # Env-based configuration loading (config.go)
│   │
│   └── pkg/                    # Shared utilities
│       ├── errors/             # Custom error types & HTTP status mapping
│       ├── validator/          # go-playground/validator setup & custom validators
│       ├── jwt/                # Token generation & validation helpers
│       ├── pagination/         # Pagination param parsing & response formatting
│       ├── response/           # Standard JSON response writers
│       └── logger/             # Structured logger setup (zerolog)
│
├── migrations/                 # Versioned SQL migration files (golang-migrate)
│
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

## Package Naming Conventions

- **Use lowercase, no underscores** for package names (Go convention: `service`, not `service_layer`).
- Match package name to directory name (e.g., `internal/api/handler/` has `package handler`).
- Avoid generic names like `utils` or `common` — use specific names like `pagination`, `response`.
- One package purpose: `repository` only contains database access code.

## Folder Purpose Summary

| Directory | Purpose |
|-----------|---------|
| `cmd/server/` | Application entry, dependency injection, server startup |
| `internal/api/handler/` | HTTP request parsing, validation, response writing |
| `internal/api/middleware/` | Cross-cutting concerns (auth, CORS, logging) |
| `internal/api/dto/` | Request validation structs and response shapes |
| `internal/service/` | Business logic, orchestration, authorization |
| `internal/repository/` | GORM queries, no business logic |
| `internal/model/` | Database table mapping via GORM tags |
| `internal/config/` | Environment variable loading |
| `internal/pkg/` | Reusable cross-module utilities |
| `migrations/` | SQL migration files (timestamp-prefixed) |

## Import Path Convention

```
github.com/yourorg/app-owndangan/internal/{layer}/{module}
```

Example:

```go
import (
    "github.com/yourorg/app-owndangan/internal/api/handler"
    "github.com/yourorg/app-owndangan/internal/service"
    "github.com/yourorg/app-owndangan/internal/repository"
)
```

## Module Organization

Each business domain (auth, events, guests, etc.) has files spread across layers:

```
internal/api/handler/auth_handler.go
internal/api/dto/auth_dto.go
internal/service/auth_service.go
internal/repository/user_repo.go    (also used by user service)
internal/repository/refresh_token_repo.go
internal/model/user.go
internal/model/refresh_token.go
```

This keeps domain logic discoverable while maintaining strict layer separation and testability.
