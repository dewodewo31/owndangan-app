# Module: Analytics

## Purpose

Track platform-wide and per-invitation analytics data. Provides insights on invitation views, RSVP statistics, page visit patterns, and platform growth metrics. Enables data-driven decision-making for the platform owner and insight for the couple.

## Responsibilities

- Track page views for public invitation pages.
- Track RSVP submission events.
- Track digital gift link clicks.
- Provide per-invitation view counts.
- Provide platform-wide dashboard metrics (total users, active subscriptions, revenue, active events).
- Store analytics events immutably for historical analysis.
- Support basic aggregation queries (daily, weekly, monthly).

## Non-Responsibilities

- Real-time analytics streaming (batch processing, not real-time).
- User behavior tracking beyond page views (clicks, scroll depth, time on page).
- A/B testing or experimentation platform.
- Third-party analytics integration (Google Analytics, etc. — separate concern).
- Heatmaps or session recordings.

## Actors

- **Admin:** Views platform-wide analytics dashboard.
- **User (couple):** Views invitation-level view counts and RSVP stats.
- **System (automated):** Records analytics events on page loads and guest actions.

## Business Rules

- Page view tracking: increment `events.view_count` on each unique visit (or rate-limited).
- View deduplication: same IP + same event within 30 minutes counts as one view (optional, configurable).
- RSVP events are tracked separately from page views.
- Analytics events are immutable — once written, they cannot be deleted or modified.
- Analytics events are stored for platform metrics (not deleted on event expiry).
- Per-invitation view count is available to the couple.
- Platform-wide metrics are available to the admin dashboard.
- Metrics aggregation: total users, total active subscriptions, total revenue (settlement sums), total active events.
- Revenue metric is the sum of `gross_amount` from transactions with status `settlement`.
- Revenue is not adjusted for refunds (future).
- Analytics event types: `page_view`, `rsvp_submitted`, `gift_clicked`, `whatsapp_link_clicked`.

## Entities

- **AnalyticsEvent:** `{ id, event_id, event_type, metadata, ip_address, user_agent, created_at }`
- **Event (view_count):** The `events.view_count` column is a cached counter derived from analytics events.

## Database

- Table: `analytics_events`
- Append-only, immutable (no UPDATE, no DELETE).
- Index on `event_type`, `created_at`.
- `events.view_count` is a counter column updated on page view events.
- Future: materialized views or aggregated tables for dashboard metrics.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/public/analytics/page-view` | Public | Record a page view |
| GET | `/api/v1/events/:id/analytics` | JWT | Get invitation-level analytics |
| GET | `/api/v1/admin/dashboard` | Admin | Platform-wide dashboard metrics |

## Request Flow

```
POST /public/analytics/page-view
  → Handler: parse body { slug, referrer }
  → Service: resolve event by slug, verify published
  → Service: check deduplication window (optional)
  → Service: increment events.view_count
  → Service: insert analytics_event (type=page_view, ip, user_agent, referrer)
  → Handler: return 204 No Content (no data returned)
```

```
GET /admin/dashboard
  → Handler: parse query params (date range)
  → Service: aggregate total users, active subscriptions, revenue, active events
  → Service: query analytics_events for visitor counts in range
  → Service: compute growth rates (vs. previous period)
  → Handler: return dashboard summary
```

## Validation

- slug: must match a published event.
- event_type: one of known types (`page_view`, `rsvp_submitted`, `gift_clicked`, `whatsapp_link_clicked`).
- IP address: captured server-side, not client-supplied (trusted).
- Referrer: optional, max 500 chars, validated URL format.

## Authorization

- Page view recording: public (rate-limited).
- Invitation analytics: JWT + event ownership.
- Dashboard: JWT + admin role.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 404 | NOT_FOUND | Event slug not found |
| 429 | RATE_LIMITED | Too many page view requests |

## Security Considerations

- Analytics events should not store PII beyond IP address (and IP is retained only for deduplication, not for tracking individuals).
- Consider IP anonymization (truncate last octet) for GDPR/UU PDP compliance.
- Rate-limit the page view endpoint to prevent view-count manipulation.
- View count from the events table is a cached approximation — not real-time.
- Analytics events cannot be deleted — ensure no PII is written to them.
- User agent is stored for device/browser breakdown — ensure it is not used for fingerprinting.

## Testing Requirements

- Unit tests for view deduplication logic.
- Integration tests for page view recording.
- Integration tests for dashboard aggregation.
- Test that view count increments correctly.
- Test rate limiting on page view endpoint.
- Test that analytics events are immutable.

## Dependencies

- Events module — view_count, slug resolution.
- Subscriptions module — active subscription count.
- Users module — total user count.
- Transactions module — revenue aggregation.

## Related Modules

- **Events** — View count tracking target.
- **RSVP** — RSVP analytics events.
- **Admin** — Dashboard analytics consumption.
- **Digital Gifts** — Gift click tracking.

## Known Limitations

- No real-time analytics (batching/reporting delay).
- No per-event detailed analytics (e.g., visits by day, device breakdown, location).
- No export of analytics data.
- View deduplication is basic (IP-based, 30-min window) — does not account for shared IPs (NAT).
- No analytics for admin actions (audit log covers that).
- No anomaly detection (e.g., sudden traffic spike detection).

## TODO

- [ ] Implement IP anonymization for UU PDP compliance.
- [ ] Add per-event detailed analytics (visits by day, device, location).
- [ ] Add analytics data export.
- [ ] Add materialized views for dashboard aggregates.
- [ ] Consider event store vs. relational DB for analytics (high volume).
- [ ] Add bot detection (exclude known crawlers from view count).