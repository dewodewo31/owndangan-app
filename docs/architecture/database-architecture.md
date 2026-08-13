# Database Architecture

## Overview

PostgreSQL is the single source of truth. All data is persisted in a single database instance (read replicas planned for future). GORM is the ORM for Go backend.

## Connection Management

```go
// Database connection configuration
dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logLevel),
    SkipDefaultTransaction: true, // performance
    PrepareStmt: true,            // prepared statement cache
})

sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Data Flow

```
Handler → Service → Repository → GORM → PostgreSQL
```

## Migration Strategy

See `docs/database/migrations.md`.

## Entity Ownership

| Entity | Owner | Cascade |
|--------|-------|---------|
| users | self | — |
| subscriptions | users | CASCADE delete |
| transactions | users | RESTRICT delete |
| events | users | CASCADE delete |
| event_sections | events | CASCADE delete |
| guests | events | CASCADE delete |
| rsvps | guests | CASCADE delete |
| guestbook_messages | events | CASCADE delete |
| digital_gifts | events | CASCADE delete |
| gallery_photos | events | CASCADE delete |
| music | events | SET NULL |
| templates | admin | RESTRICT delete |
| packages | admin | RESTRICT delete |
| audit_logs | system | — |
| analytics_events | system | — |

## Backup Strategy

> TODO: Define backup schedule (daily, weekly).

## Key Design Decisions

1. **UUID primary keys** — Avoids sequential ID enumeration, supports distributed systems.
2. **JSONB for flexible data** — Features (packages), bank accounts (digital_gifts), midtrans response (transactions).
3. **Soft deletes** — For user data that may need recovery (users, events, guests).
4. **Hard deletes** — For immutable audit data (transactions, audit_logs, analytics_events).
5. **Timestamps on all tables** — `created_at`, `updated_at` for all tables; `deleted_at` for soft-delete tables.
6. **BIGINT for money** — Avoids floating-point precision issues, stores in IDR (smallest unit).

## Related Documentation

- `docs/database/schema.md`
- `docs/database/relationships.md`
- `docs/database/indexes.md`
- `docs/database/conventions.md`
- `docs/database/migrations.md`