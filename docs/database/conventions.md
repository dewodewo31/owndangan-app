# Database Conventions

## Naming

- Tables: `snake_case` plural nouns (e.g., `guestbook_messages`, not `guestbookmessage`).
- Columns: `snake_case`.
- Primary key: always `id` (UUID type).
- Foreign key: `{referenced_table_singular}_id` (e.g., `user_id`, `event_id`).
- Timestamps: `created_at`, `updated_at`, `deleted_at`.
- Boolean columns: `is_` prefix (e.g., `is_active`, `is_approved`).
- Status columns: `status` (VARCHAR), not a boolean.

## Types

- UUID for all primary keys (`gen_random_uuid()`).
- `TIMESTAMP WITH TIME ZONE` for all timestamps.
- `BIGINT` for monetary amounts (IDR, stored as smallest unit — Rupiah, no decimal).
- `VARCHAR(255)` for short strings, `TEXT` for long strings.
- `JSONB` for flexible data (features, bank accounts, raw webhook responses).
- `DECIMAL` only if needed; prefer BIGINT for money.

## Soft Deletes

Use `deleted_at TIMESTAMP WITH TIME ZONE` for tables where data recovery matters:
- users, subscriptions, events, guests

Hard delete for audit/immutable data:
- transactions, audit_logs, analytics_events, rsvps

GORM convention: use `gorm.DeletedAt` type for soft delete fields.

## Timestamps

GORM auto-manages `CreatedAt`, `UpdatedAt`, and `DeletedAt` when field names match.

```go
CreatedAt time.Time
UpdatedAt time.Time
DeletedAt gorm.DeletedAt
```

## Default Values

- `id`: `gen_random_uuid()` (UUID v4)
- `created_at`: `CURRENT_TIMESTAMP`
- `status`: Default to most common value
- `is_*`: `false` (safe default)
- `sort_order`: `0`

## Migrations

- Development: GORM AutoMigrate is acceptable.
- Production: versioned SQL migrations (golang-migrate or similar).
- Migration files: `YYYYMMDDHHMMSS_description.sql` (timestamp-prefixed).
- Each migration: one up and one down file.

## Auditing

Tables with audit fields:
- All tables: `created_at`, `updated_at`
- Soft-delete tables: `deleted_at`
- Audit log table: `audit_logs` for important actions (subscription activation, payment settlement, admin actions).

## JSONB Usage

JSONB columns (in `packages` and `transactions`):

```go
type Package struct {
    Features datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
}
```

In queries, use GORM JSON query functions:
```go
db.Where("features->>'custom_domain' = ?", "true").Find(&packages)
```

## Query Conventions

- Use GORM query methods instead of raw SQL where possible.
- Raw SQL is acceptable for complex reports or bulk operations.
- Always use parameterized queries (GORM handles this by default).
- Avoid N+1 queries — use `Preload` or `Joins`.
- Use `Scopes` for reusable query filters (e.g., ActiveSubscription, PublishedEvents).

## Connection Pool

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```
