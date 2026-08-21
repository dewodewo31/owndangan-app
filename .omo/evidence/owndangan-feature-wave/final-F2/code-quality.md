# F2 — Code Quality Review (diff 7d31e8c..HEAD vs AGENTS.md rules)

## Layering
- Handlers thin: parse -> service call -> response (analytics_handler, payment/guest handlers unchanged in shape). Business logic in services (analytics_service, payment_service, guest service, email worker). DB access in repositories (new sub/audit methods added at repo layer).

## Server-side authorization
- GET /events/{id}/analytics: ownership checked against authenticated user; entitlement AnalyticsEnabled() resolved server-side (free -> 403). Tests assert both.
- Guest restore/delete: event-owner check -> 403 for foreign users (tested).
- Expiry worker operates via repos only; status flip server-side.

## Input validation & safety
- POST /analytics/events validates type whitelist -> 400 on invalid; unpublished events -> 404 (tested).
- Email renderer uses html/template (auto-escape) — injection test asserts `<script>` escapes.
- Webhook signature verification unchanged; idempotency intact.

## Secrets & suppressions
- SMTP credentials via config/env only (.env.example placeholders); none hardcoded.
- `as any` occurrences exist ONLY in editor/page.tsx — verbatim from initial commit 7d31e8c (pickaxe-verified), i.e. PRE-EXISTING WIP, not introduced by this plan; left untouched per minimal-change rule.
- No @ts-ignore/@ts-expect-error introduced. No empty catch blocks added; swallowed errors are deliberate best-effort paths, logged via zerolog (documented at sendPaymentSuccessEmail).

## Error handling
- Async email failures never block callers (SendAsync goroutine + logged); SendWithRetry bounded (3 attempts); webhook response/idempotency behavior preserved (duplicate-send test green).

Verdict: **APPROVE** (0 new violations; 1 pre-existing suppression documented).
