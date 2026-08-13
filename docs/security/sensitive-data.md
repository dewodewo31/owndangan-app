# Sensitive Data Handling

## Data Classification

| Category | Examples | Classification | Storage Requirements |
|---|---|---|---|
| PII | Name, email, phone, address | Sensitive | Encrypted at rest, access-logged |
| Authentication | Password hash, refresh token hash | Critical | bcrypt hashed, never logged |
| Payment | Midtrans transaction ID, order ID | Sensitive | Logged in audit trail, no raw card data |
| Business | Event details, guest list, RSVPs | Internal | Standard access control |
| Secrets | API keys, JWT secret, DB password | Critical | Environment variables, never in code |

## Personally Identifiable Information (PII)

### PII Fields
- `users.name` — Full name
- `users.email` — Email address
- `users.phone` — Phone number
- `events.groom_name`, `events.bride_name` — Wedding couple names
- `guests.name`, `guests.phone`, `guests.email` — Guest information

### PII Handling Rules
1. **Minimization**: Only collect PII that is strictly necessary for the feature.
2. **Encryption at rest**: Database storage encryption enabled (e.g., RDS encryption, or column-level encryption for phone numbers).
3. **Access control**: PII fields are only returned in API responses when explicitly needed (e.g., guest list download). List endpoints return minimal data.
4. **No PII in logs**: Never log names, emails, or phone numbers. Use user IDs instead.
5. **No PII in URLs**: Never include emails or phone numbers in URL paths or query parameters.
6. **No PII in JWT**: JWT claims contain only `sub` (user ID) and `role`.

### PII in API Responses
```json
// GET /api/v1/guests — list endpoint (safe)
{
  "guests": [
    {"id": "gst_abc", "name": "A*** B***", "rsvp_status": "confirmed"}
    // No email, no phone in list view
  ]
}

// GET /api/v1/guests/:id — detail endpoint (authorized)
{
  "id": "gst_abc",
  "name": "Alice Budiman",
  "email": "alice@example.com",
  "phone": "+6281234567890"
}
```

## Password Hashing

See [authentication-security.md](./authentication-security.md) for full details.

Summary:
- Algorithm: bcrypt, cost 12.
- Never store plaintext passwords.
- Never log password hashes.
- Never use MD5, SHA1, SHA256 for passwords.
- Minimum 8 characters, maximum 128.

## API Keys and Secrets

### Where Secrets Live
- **Development**: `.env` file (never committed to git, added to `.gitignore`).
- **Staging/Production**: Environment variables set via deployment platform (e.g., Docker secrets, Kubernetes secrets, Vault).
- **Never**: Hardcoded in source code, committed to repository, or logged.

### Secrets List
| Secret | Environment Variable | Rotation |
|---|---|---|
| JWT signing key | `JWT_SECRET` | Every 90 days |
| Midtrans server key | `MIDTRANS_SERVER_KEY` | Every 180 days |
| Database password | `DB_PASSWORD` | Every 90 days |
| S3 access key | `S3_ACCESS_KEY` | Every 180 days |
| SMTP password | `SMTP_PASSWORD` | Every 180 days |
| Refresh token key | `REFRESH_TOKEN_SECRET` | Every 90 days |

### Key Rotation Procedure
1. Generate new secret value.
2. Update environment variable in deployment.
3. Restart application (zero-downtime deploy).
4. For JWT: old tokens remain valid until expiry (grace period).
5. For Midtrans: update in Midtrans dashboard, then update env var.

## Audit Logging of Sensitive Operations

### Events to Audit

| Event | Data Logged | Retention |
|---|---|---|
| User login | user_id, IP, timestamp, success/failure | 90 days |
| User registration | user_id, IP, timestamp | 90 days |
| Password change | user_id, timestamp | 1 year |
| Email change | user_id, old_email_hash, new_email_hash, timestamp | 1 year |
| Payment webhook | order_id, transaction_id, status, IP, signature_valid | 1 year |
| Subscription change | user_id, plan, action, timestamp | 1 year |
| Admin action | admin_id, action, resource, timestamp | 1 year |
| Failed auth attempt | IP, email_hash, timestamp | 30 days |

### Audit Log Format
```json
{
  "timestamp": "2025-01-15T10:30:00Z",
  "event": "user.login",
  "actor_id": "usr_abc123",
  "actor_ip": "203.0.113.42",
  "resource_type": "session",
  "resource_id": "sess_xyz",
  "outcome": "success",
  "metadata": {}
}
```

### Audit Implementation
- Write to a separate `audit_logs` table (not mixed with business data).
- Use `SERIALIZABLE` isolation level for critical audit events.
- Audit logs are append-only (no UPDATE, no DELETE).
- Audit logs are accessible only to admins.

## Data Retention and Deletion

### Retention Schedule

| Data Type | Retention Period | Action After |
|---|---|---|
| Active user accounts | Indefinite (while active) | — |
| Deleted user accounts | 30 days (soft delete) | Hard delete + anonymize |
| Event data | Indefinite (while account active) | — |
| Deleted events | 30 days (soft delete) | Hard delete |
| Guest data | 30 days after event completion | Hard delete |
| Payment logs | 5 years (regulatory) | Archive, then delete |
| Audit logs | 1 year | Archive, then delete |
| Session tokens | 7 days (refresh token expiry) | Auto-deleted |

### Account Deletion Flow
1. User requests account deletion.
2. Soft-delete user (set `deleted_at`).
3. Soft-delete all user events.
4. Anonymize PII in audit logs (replace with hashes).
5. After 30 days: hard delete user data, retain anonymized payment records.
6. Send confirmation email to user.

### Anonymization
When hard-deleting, replace PII with irreversible hashes:
```sql
UPDATE users SET
    email = CONCAT('deleted-', id, '@anonymized'),
    name = 'Deleted User',
    phone = NULL
WHERE id = 'usr_abc123';
```