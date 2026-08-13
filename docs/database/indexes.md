# Database Indexes

## Index Strategy

### Primary Indexes (PK)
All tables: `id` (UUID) — auto-indexed by GORM.

### Unique Indexes
| Table | Columns | Condition | Purpose |
|-------|---------|-----------|---------|
| users | email | WHERE deleted_at IS NULL | Login, duplicate prevention |
| users | phone | WHERE deleted_at IS NULL | WhatsApp, duplicate prevention |
| events | slug | WHERE deleted_at IS NULL | Public URL uniqueness |
| guests | token | (always) | RSVP token uniqueness |
| transactions | order_id | (always) | Midtrans order reference |
| packages | code | (always) | Internal code reference |
| packages | name | (always) | Admin display |
| rsvps | guest_id | (always) | One RSVP per guest |
| digital_gifts | event_id | (always) | One gift config per event |
| event_sections | event_id | (always) | One section config per event |

### Performance Indexes
| Table | Column(s) | Purpose |
|-------|-----------|---------|
| subscriptions | user_id | Find user's subscriptions |
| subscriptions | status, expires_at | Find expired subscriptions (cron) |
| transactions | user_id | Find user's transactions |
| events | user_id | Find user's events |
| events | status | Filter draft/published events |
| guests | event_id | Find event's guests |
| rsvps | event_id | RSVP recap per event |
| guestbook_messages | event_id, is_approved | Moderated listing |
| gallery_photos | event_id, sort_order | Gallery display order |
| music | event_id | Find event's music |
| audit_logs | user_id | User action history |
| audit_logs | action | Filter by action type |
| audit_logs | created_at | Time-based queries |
| analytics_events | event_type, created_at | Analytics queries |

### Composite Indexes
| Table | Columns | Purpose |
|-------|---------|---------|
| guestbook_messages | event_id, is_approved, created_at | Sorted moderated messages |
| subscriptions | user_id, status | Active subscription lookup |
| transactions | user_id, status | User's payment history |
| events | user_id, status | User's filtered events |

## Index Naming Convention

```
idx_{table}_{column(s)}
unq_{table}_{column(s)}  (unique)
```

Examples: `idx_events_user_id`, `unq_events_slug`

## GORM Index Tags

```go
type Event struct {
    Slug string `gorm:"uniqueIndex:unq_events_slug,where:deleted_at IS NULL"`
    UserID uuid.UUID `gorm:"index:idx_events_user_id"`
}
```

## Future Index Considerations

- Full-text search index on events for search feature (GIN index).
- Partial index on analytics_events for faster aggregation queries.
- BRIN index on created_at for time-series queries if volume grows.