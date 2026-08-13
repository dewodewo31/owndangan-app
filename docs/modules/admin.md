# Module: Admin

## Purpose

Provide platform administration capabilities: manage users, packages, templates, transactions, and monitor platform analytics. This is the root-level management interface for the platform owner.

## Responsibilities

- Display platform-wide dashboard analytics (registered users, active subscriptions, revenue, active invitations).
- Manage user accounts: list, view, suspend, activate, change role.
- Manage packages: create, update pricing, toggle feature flags, activate/deactivate.
- Manage templates: upload, update, activate/deactivate, delete.
- View all transactions with filtering and search.
- Moderate guestbook messages across all events.
- View audit logs for platform activity.
- Override subscription status for any user (admin intervention).

## Non-Responsibilities

- Creating invitations or managing guest lists (user-only functionality).
- Processing payments (handled by Payments module).
- Managing own admin profile (handled by Users module).
- Generating reports (future feature).

## Actors

- **Admin (platform owner):** Full access to all admin endpoints. Must have `role = 'admin'`.

## Business Rules

- Only users with `role = 'admin'` can access admin endpoints.
- Admin can view any user's data but cannot impersonate the user for actions.
- Admin can suspend/activate any user account.
- Admin can change a user's role (promote to admin or demote to user).
- Admin can create, update, or delete packages.
- Admin can upload, activate, or deactivate templates.
- Admin can view all transactions regardless of status.
- Admin can manually extend or terminate any subscription.
- All admin actions are logged in the audit log with the admin's user ID.
- Admin cannot delete their own admin account.
- At least one admin must always exist in the system.

## Entities

- No unique entities; this module operates on existing entities: users, packages, templates, transactions, subscriptions, audit_logs.

## Database

- Reads from: `users`, `packages`, `templates`, `transactions`, `subscriptions`, `audit_logs`, `events`, `analytics_events`.
- Writes to: `users` (status updates), `packages`, `templates`, `subscriptions` (admin overrides).

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/dashboard` | Platform-wide analytics summary |
| GET | `/api/v1/admin/users` | List users (paginated, filterable) |
| GET | `/api/v1/admin/users/:id` | Get user details |
| PUT | `/api/v1/admin/users/:id/status` | Suspend or activate user |
| PUT | `/api/v1/admin/users/:id/role` | Change user role |
| GET | `/api/v1/admin/packages` | List all packages |
| POST | `/api/v1/admin/packages` | Create new package |
| PUT | `/api/v1/admin/packages/:id` | Update package |
| GET | `/api/v1/admin/transactions` | List all transactions |
| GET | `/api/v1/admin/transactions/:id` | Get transaction details |
| GET | `/api/v1/admin/templates` | List templates |
| POST | `/api/v1/admin/templates` | Create/upload template |
| PUT | `/api/v1/admin/templates/:id` | Update template |
| DELETE | `/api/v1/admin/templates/:id` | Delete template |
| GET | `/api/v1/admin/subscriptions` | List all subscriptions |
| POST | `/api/v1/admin/subscriptions/:id/extend` | Extend a subscription |
| POST | `/api/v1/admin/subscriptions/:id/terminate` | Terminate a subscription |
| GET | `/api/v1/admin/audit-logs` | View audit logs |
| GET | `/api/v1/admin/guestbook` | Moderate guestbook messages across events |

## Request Flow

```
GET /admin/dashboard
  → Handler: parse query params (date range)
  → Service: aggregate total users, active subscriptions, revenue, active events
  → Service: query analytics_events for visitor stats
  → Handler: return dashboard summary
```

```
PUT /admin/users/:id/status
  → Handler: parse user ID, extract new status
  → Service: verify user exists, check not last admin if suspending admin
  → Service: update user status, log action in audit log
  → Handler: return updated user
```

## Validation

- User status must be one of: `active`, `suspended`.
- Role must be one of: `user`, `admin`.
- Package price must be a positive integer (IDR, no decimals).
- Package duration_days must be positive integer or null (lifetime).
- Template config must be valid JSON.

## Authorization

- JWT required with `role = 'admin'` in token payload.
- Admin middleware checks role before handler execution.
- Admin cannot be decoded from user-submitted data; only JWT token is authoritative.

## Error Cases

| HTTP | Code | Condition |
|------|------|-----------|
| 401 | UNAUTHORIZED | Missing or invalid JWT |
| 403 | FORBIDDEN | User is not admin |
| 404 | NOT_FOUND | User/package/template not found |
| 422 | VALIDATION_ERROR | Invalid input data |
| 409 | CONFLICT | Cannot delete last admin, package name conflict |

## Security Considerations

- Admin endpoints are high-value targets; rate limiting is stricter than user endpoints.
- All admin mutations are logged in audit log with admin identity, timestamp, IP, and action.
- Admin cannot create subscriptions without a valid transaction (no free subscriptions).
- Admin suspension of a user must not affect the user's published invitations.
- Admin role should be assigned sparingly; only platform owners.
- Admin impersonation is explicitly forbidden — no API to "login as user."

## Testing Requirements

- Unit tests for all admin service operations.
- Integration tests for authorization (non-admin users get 403).
- Test suspend/activate user flow.
- Test package CRUD operations.
- Test template management.
- Test dashboard analytics aggregation.
- Test that admin actions appear in audit log.
- Test that last admin cannot be deleted.

## Dependencies

- Authentication module for JWT middleware.
- Users module for user management.
- Packages module for package management.
- Templates module for template management.
- Subscriptions module for subscription override.
- Payments module for transaction viewing.
- Audit Log module for action logging.
- Analytics module for dashboard metrics.

## Related Modules

- All modules — admin has read access across the entire platform.

## Known Limitations

- No role-based sub-permissions (admin is all-or-nothing).
- No admin notification system (e.g., alert on failed transactions).
- No bulk operations (admin must manage users/packages individually).
- No export of analytics data.
- No admin activity dashboard (who did what, when).

## TODO

- [ ] Implement admin notification system (low stock of templates, high error rate).
- [ ] Add role-based sub-permissions (e.g., support agent can view but not edit packages).
- [ ] Add analytics data export.
- [ ] Add admin activity timeline.
- [ ] Implement webhook testing tool for admin to simulate Midtrans notifications.