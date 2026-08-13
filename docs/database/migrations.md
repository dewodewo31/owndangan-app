# Database Migrations

## Migration Strategy

### Development
- Use GORM AutoMigrate for rapid prototyping.
- GORM creates/updates tables based on model definitions.
- Do NOT rely on AutoMigrate for production schema changes.

### Production
- Use versioned SQL migrations.
- Migration tool: > TODO (golang-migrate / goose / pressly/goose)
- Migration files live in `backend/migrations/`.

## Migration File Format

```
backend/migrations/
├── 20260811000001_create_users.up.sql
├── 20260811000001_create_users.down.sql
├── 20260811000002_create_packages.up.sql
├── 20260811000002_create_packages.down.sql
├── 20260811000003_create_transactions.up.sql
├── 20260811000003_create_transactions.down.sql
├── 20260811000004_create_events.up.sql
├── 20260811000004_create_events.down.sql
...
```

## Migration Rules

1. All migrations must be reversible (up AND down).
2. Never modify an already-applied migration.
3. To change a table: create a new migration with ALTER TABLE.
4. Test migrations on a staging database before production.
5. Back up production database before running migrations.
6. Use transactions in migration files for atomicity.

## Migration Naming

Format: `YYYYMMDDHHMMSS_description.up.sql`

Example: `20260811000001_create_users.up.sql`

## Seed Data

> TODO: Define seed data approach.

Required seed data:
- Default packages (Free, Basic, Premium, Pro)
- Admin user account
- Default templates (themes)

## Migration Commands

> TODO: Define after choosing migration tool.

## GORM AutoMigrate Usage (Development Only)

```go
db.AutoMigrate(
    &model.User{},
    &model.Package{},
    &model.Subscription{},
    &model.Transaction{},
    &model.Event{},
    &model.EventSection{},
    &model.Template{},
    &model.Guest{},
    &model.RSVP{},
    &model.GuestbookMessage{},
    &model.DigitalGift{},
    &model.GalleryPhoto{},
    &model.Music{},
    &model.AuditLog{},
    &model.AnalyticsEvent{},
)
```
