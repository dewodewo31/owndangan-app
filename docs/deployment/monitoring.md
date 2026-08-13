# Monitoring

## Overview

Production monitoring covers four pillars: **Uptime**, **Error Tracking**, **Performance**, and **Database Metrics**. All dashboards are accessible to the engineering team.

## Uptime Monitoring

| Tool | Purpose | Frequency |
|------|---------|-----------|
| Better Uptime (or Uptime Robot) | Public endpoint checks | Every 1 minute |
| Kubernetes liveness probe | Container health | Every 30 seconds |
| External pingdom | Global availability | Every 5 minutes |

### Alerting

- **Critical**: Backend returns 5xx for >1 minute → PagerDuty call.
- **Warning**: Backend returns 5xx for >10 seconds → Slack notification.
- **Info**: SSL cert expires in <14 days → Slack reminder.

### Health Endpoint

```
GET /health → 200 {"status":"ok"}
```

Monitored: HTTP status, response time (<500ms), response body validity.

## Error Tracking (Sentry)

**Sentry** is configured for both backend and frontend.

### Backend (Go)

```go
sentry.Init(sentry.ClientOptions{
  Dsn: os.Getenv("SENTRY_DSN"),
  Environment: "production",
  Release: version,
})
```

- Captures all unhandled panics via middleware.
- Captures 5xx errors with full stack trace.
- Captures slow requests (>2s) as performance spans.
- Excludes 4xx errors (client errors are not bugs).

### Frontend (Next.js)

```ts
Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  environment: process.env.NODE_ENV,
  tracesSampleRate: 0.2,
});
```

- Captures unhandled React errors.
- Captures API call failures (4xx/5xx).
- Captures performance traces for page loads.

### Alerting

- **Critical**: New error affecting >10 users in 5 minutes.
- **Warning**: New error type seen in the last hour.
- **Info**: Error count spikes >2x daily average.

## Performance Monitoring

### Backend

| Metric | Tool | Threshold | Action |
|--------|------|-----------|--------|
| Request latency (p95) | Prometheus + Grafana | <500ms | Alert if >1s |
| Request rate | Prometheus | — | Capacity planning |
| Error rate | Prometheus | <1% | Alert if >5% |
| Go goroutine count | Prometheus | <1000 | Alert if >5000 |
| Go memory usage | Prometheus | <500MB | Alert if >1GB |

### Frontend

| Metric | Tool | Threshold | Action |
|--------|------|-----------|--------|
| LCP (Largest Contentful Paint) | Web Vitals | <2.5s | Alert if >4s |
| FID (First Input Delay) | Web Vitals | <100ms | Alert if >300ms |
| CLS (Cumulative Layout Shift) | Web Vitals | <0.1 | Alert if >0.25 |
| API call success rate | Custom | >99% | Alert if <95% |

## Database Monitoring

| Metric | Tool | Threshold | Action |
|--------|------|-----------|--------|
| Connections | Grafana + `pg_stat_activity` | <80% of max | Alert if >80% |
| Query latency (p95) | pg_stat_statements | <100ms | Alert if >500ms |
| Cache hit ratio | pg_statio | >99% | Investigate if <95% |
| Disk usage | Cloud provider | <80% | Alert if >80% |
| Replication lag | pg_stat_replication | <1s | Alert if >5s |

## Payment Monitoring

| Metric | Tool | Threshold | Action |
|--------|------|-----------|--------|
| Midtrans API errors | Backend logs | 0 | Alert on any error |
| Payment success rate | Custom metric | >95% | Alert if <90% |
| Webhook processing time | Backend metrics | <2s | Alert if >5s |
| Failed transactions | Midtrans dashboard | — | Daily review |

## Dashboards

- **Grafana**: Backend + database metrics.
- **Sentry**: Error tracking and performance.
- **Better Uptime**: Uptime status page.
- **Midtrans Dashboard**: Payment transaction monitoring.

## On-Call

- On-call engineer receives critical alerts via PagerDuty.
- Response time: 15 minutes for critical, 1 hour for warning.
- Post-incident RCA is documented within 48 hours.