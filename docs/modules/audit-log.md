# Module: Audit Log

## Purpose

Provide an immutable, tamper-evident record of all important actions performed on the platform. Supports compliance, security monitoring, debugging, and admin oversight. Every administrative action and key user action is logged.

## Responsibilities

- Log all admin actions (user management, package changes, template changes, subscription overrides).
- Log key user actions (event create, publish, unpublish, delete).
- Log payment-related events (transaction creation, webhook receipt, settlement).
- Log authentication events (login, failed login, logout, refresh).
- Store immutable audit records with: actor, action, entity, metadata, IP address, and timestamp.
- Provide audit log querying for admins (filterable by action, entity, user, date range).

## Non-Responsibilities

- Application debug logging (structured logging for developers, not audit).
- Tracking read operations (audit logs are for mutations that change state).
- Analytics or user behavior tracking (handled by Analytics module).
- Storing full request/response payloads (stored in application logs, not audit log).
- Log retention/purging (immutable by design; retention policy is TBD).

## Actors

- **Admin:** Reads audit logs; all admin actions are automatically logged.
- **User (couple):** Some user actions are logged (event management).
- **System (background):** Automated actions are logged (webhook processing, etc.).

## Business Rules

- Audit log entries are immutable — once written, they cannot be modified or deleted.
- Each audit entry captures: `user_id` (who did it, NULL if system), `action` (what was done), `entity_type` (what was affected), `entity_id` (which record), `metadata` (additional context as JSONB), `ip_address` (source IP), `created_at` (when).
- Action naming convention: `{entity}.{action}`, e.g., `event.created`, `subscription.activated`, `user.suspended`.
- All admin-initiated mutations must be logged.
- System actions (e.g., webhook processing) are logged with `user_id = NULL` and metadata indicating the source.
- Audit logs are read-only via the admin API.
- No user-facing audit log access (admin only).
- IP address is captured from the request (X-Forwarded-For or remote address).
- Metadata JSONB can contain additional context, e.g., previous values, reason, or request ID.

## Entities

- **AuditLog:** `{ id, user_id, action, entity_type, entity_id, metadata, ip_address, created_at }`

## Database

- Table: `audit_logs`
- Append-only — no UPDATE, no DELETE.
- Indexes: `user_id`, `action`, `created_at`.
- No soft delete (immutable by design).
- Metadata: JSONB for flexible context storage.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/admin/audit-logs` | Admin | Query audit logs (paginated, filterable) |

### Filter parameters

- `user_id` — filter by actor.
- `action` — filter by action type (e.g., `event.created`).
- `entity_type` — filter by entity (e.g., `event`, `subscription`).
- `entity_id` — filter by specific entity.
- `from` / `to` — date range filter (ISO 8601).
- `page` / `per_page` — pagination.

## Request Flow

```
Audit log write (called by other services)
  → Service (any module): after successful mutation, call AuditService.Log()
  → AuditService: build entry (user_id from context, action, entity, metadata, IP)
  → AuditRepository: insert record (immutable, no update)
  → (no return to caller — fire-and-forget or synchronous per requirement)
```

```
GET /admin/audit-logs
  → Handler: parse query filters
  → Service: verify admin role
  → Service: build query with filters, paginate
  → Handler: return audit log entries
```

## Validation

- action: required, must follow `{entity}.{action}` convention, max 100 chars.
- entity_type: optional, max 50 chars.
- entity_id: optional, valid UUID.
- metadata: optional, valid JSON object.
- ip_address: captured server-side, validated format.

## Authorization

- Audit log writing: internal service call (no HTTP endpoint for writing).
- Audit log reading: JWT with admin role only.
- Users cannot view audit logs.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing/invalid JWT |
| 403 | FORBIDDEN | Not admin |
| 422 | VALIDATION_ERROR | Invalid filter parameters |

## Security Considerations

- Audit logs are immutable — this is a security feature to prevent tampering.
- Consider append-only database permissions for the audit_logs table (INSERT only for the application user).
- Metadata must not contain secrets (passwords, tokens, full payment data).
- Log entry integrity: consider a hash chain for tamper evidence (future).
- IP address is PII in some jurisdictions — consider retention policy.
- Large metadata payloads can bloat the table — enforce a size limit (e.g., 4KB JSONB).
- Audit log access is admin-only, but should be logged itself (meta-audit) or at least logged in application logs.

## Testing Requirements

- Unit tests for audit log entry creation (action naming, metadata).
- Integration tests for admin log querying (filtering, pagination).
- Integration tests that admin actions produce audit entries.
- Test that user actions (event publish) produce audit entries.
- Test that system actions (webhook) produce audit entries with NULL user_id.
- Test immutability (no update/delete API).
- Test metadata size limit.

## Dependencies

- Authentication module — user context (user_id, IP) for all services.
- No external dependencies beyond the database.

## Related Modules

- **All modules** — Every module calls AuditLog for important mutations.
- **Admin** — Primary consumer of audit log queries.
- **Security** — Audit logs support security monitoring and incident response.

## Known Limitations

- No integrity verification (hash chain) for tamper evidence.
- No automated alerting on suspicious audit events.
- No retention policy (log grows indefinitely).
- Metadata size limit may be too small for some complex operations.
- No user-facing audit (users cannot see their own action history).
- No export of audit logs.

## TODO

- [ ] Implement hash chain for tamper-evident audit log.
- [ ] Add retention policy and archival strategy.
- [ ] Add automated alerting on suspicious actions (e.g., multiple admin logins, bulk user suspension).
- [ ] Add user-facing activity log (limited to own actions).
- [ ] Add audit log export for compliance.
- [ ] Add metadata size limit enforcement.