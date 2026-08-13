# Database Deployment

## Migration Strategy

We use **goose** for database migrations. Migrations are SQL files in `backend/migrations/`.

### Migration Rules

- **Always forward-compatible**: New columns must be nullable or have defaults.
- **Never delete a column** in the same deployment that removes code referencing it.
- **Rollback migrations** are required for every forward migration.
- **Zero-downtime** migrations: Add columns/ tables first, remove old code in a follow-up.

### Running Migrations

```bash
# Apply all pending migrations
goose postgres "$DATABASE_URL" up

# Rollback the last migration
goose postgres "$DATABASE_URL" down

# Check status
goose postgres "$DATABASE_URL" status
```

### CI/CD Pipeline

Migrations run automatically in the deployment pipeline:

1. Before the new backend version starts.
2. In a transaction (one migration at a time).
3. If migration fails, the deployment is rolled back.

```yaml
# deploy step
- name: Run migrations
  run: goose postgres "$DATABASE_URL" up
- name: Deploy backend
  run: docker compose up -d backend
```

## Backup Strategy

### Production

| Frequency | Type | Retention | Storage |
|-----------|------|-----------|---------|
| Hourly | WAL archiving | 24 hours | Cloud storage |
| Daily | Full dump | 30 days | Cloud storage |
| Weekly | Full dump | 90 days | Cloud storage |

### Backup Command

```bash
# Full backup
pg_dump -Fc -h localhost -U postgres owndangan > backup/owndangan_$(date +%Y%m%d_%H%M%S).dump

# Restore
pg_restore -Fc -h localhost -U postgres -d owndangan backup.dump
```

### Point-in-Time Recovery (PITR)

Requires WAL archiving to be configured:

```bash
# recovery.conf
restore_command = 'aws s3 cp s3://backups/wal/%f %p'
recovery_target_time = '2025-01-15 14:30:00+07'
```

## Monitoring

### Key Metrics

- **Connection count**: Should stay below `max_connections` (default: 100).
- **Active queries**: Long-running queries (>5s) should be investigated.
- **Cache hit ratio**: Should be >99% for production workloads.
- **Replication lag**: For read replicas, should be <1 second.
- **Disk usage**: Alert at 80% capacity.

### Queries to Check

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity;

-- Long-running queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query
FROM pg_stat_activity
WHERE state != 'idle' AND now() - pg_stat_activity.query_start > interval '5 seconds';

-- Cache hit ratio
SELECT
  'index hit rate' AS name,
  (sum(idx_blks_hit)) / nullif(sum(idx_blks_hit + idx_blks_read), 0) AS ratio
FROM pg_statio_user_indexes
UNION ALL
SELECT
  'table hit rate',
  sum(heap_blks_hit) / nullif(sum(heap_blks_hit + heap_blks_read), 0)
FROM pg_statio_user_tables;
```

## Connection Pooling

The backend uses `pgxpool` for connection pooling. Configuration:

```go
config.MaxConns = 25        // Max connections in pool
config.MinConns = 5         // Idle pool size
config.MaxConnLifetime = 5 * time.Minute
config.MaxConnIdleTime = 1 * time.Minute
```

For production, consider using **PgBouncer** as a connection pooler between the app and database:

```bash
docker run -d --name pgbouncer \
  -e DATABASE_URL="postgres://user:pass@db:5432/owndangan" \
  -p 6432:6432 \
  edoburu/pgbouncer
```

This allows the backend to maintain fewer connections and reduces load on PostgreSQL.