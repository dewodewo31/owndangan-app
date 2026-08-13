# Security Overview

## Threat Model

| Threat | Impact | Likelihood | Mitigation |
|---|---|---|---|
| SQL injection | Data exfiltration, deletion | Medium | GORM parameterized queries, input validation |
| XSS (stored/reflected) | Session theft, defacement | Medium | Output encoding, CSP headers, template escaping |
| CSRF | State-changing actions on behalf of user | Medium | CSRF tokens (Go `nosurf`), SameSite cookies |
| Auth bypass | Unauthorized access to any account | High | JWT verification, rate limiting, bcrypt cost 12 |
| Payment fraud | Free upgrades, revenue loss | Critical | Server-side verification, webhook as source of truth |
| Privilege escalation | User accesses admin functions | High | Role middleware, ownership checks |
| Token theft | Account takeover | High | Short-lived JWT (15 min), refresh rotation, HTTPS only |
| Mass assignment | User sets admin=true on register | Medium | Explicit DTOs, never bind raw request to model |

## Security Principles

1. **Defense in depth** — Every layer validates independently. Frontend validation is convenience; backend validation is security.
2. **Least privilege** — Each component has minimum required permissions. Database users get SELECT/INSERT/UPDATE only where needed.
3. **Never trust client input** — All request data (headers, body, query params) is untrusted until validated.
4. **Fail securely** — Default-deny access. Errors reveal no internals (no stack traces, no SQL errors).
5. **Secure by default** — HTTPS enforced, HSTS enabled, cookies set HttpOnly+Secure+SameSite.
6. **Separation of concerns** — Handlers parse, services contain logic, repositories access DB. Auth never bypasses middleware.

## Defense in Depth Layers

```
Internet
  ├── TLS 1.3 (HTTPS enforcement, HSTS)
  ├── WAF / rate limiter (nginx/Cloudflare)
  ├── Go HTTP router (middleware stack)
  │   ├── CORS middleware
  │   ├── CSRF middleware (nosurf)
  │   ├── Rate limiter (token bucket per IP/user)
  │   ├── Auth middleware (JWT verification)
  │   └── Role/ownership middleware
  ├── Handler (input parsing, validation)
  ├── Service (business logic, authorization)
  ├── Repository (parameterized queries)
  └── Database (least privilege user, encrypted at rest)
```

## Key Security Risks

### Injection Attacks
- SQL: Prevented by GORM parameterized queries. Never use `Raw` or `Exec` with string concatenation.
- NoSQL: If MongoDB is used, sanitize query operators (`$gt`, `$ne`).
- Template injection: Go templates auto-escape HTML. Never use `template.HTML` with user input.

### Cross-Site Scripting (XSS)
- Go templates escape by default. Verify `text/template` is not used for HTML output.
- Next.js JSX auto-escapes. Avoid `dangerouslySetInnerHTML`.
- Set `Content-Security-Policy` header: `default-src 'self'`.
- Sanitize rich text input (if any) with a library like `bluemonday`.

### Cross-Site Request Forgery (CSRF)
- All state-changing requests (POST, PUT, DELETE) require CSRF token.
- Use `nosurf` middleware in Go. Exempt webhook endpoints (no cookie auth).
- Set `SameSite=Strict` on session cookies.

### Authentication Bypass
- JWT signed with HS256 using a 256-bit secret from environment variable.
- Access token expiry: 15 minutes. Refresh token expiry: 7 days.
- Refresh token rotation: old token invalidated on each use.
- Rate limit login: 5 attempts per minute per IP, 10 per hour per email.

### Payment Fraud
- Frontend never reports payment status to backend.
- Midtrans Snap token is generated server-side only.
- Webhook notification is the sole authoritative payment confirmation.
- Transaction status verified against Midtrans API before granting access.
- Idempotency key prevents duplicate processing.