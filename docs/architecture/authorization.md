# Authorization Architecture

## Authorization Model

Two-layer authorization:
1. **Role-based (RBAC):** Determines which areas a user can access (admin vs user).
2. **Ownership-based:** Determines which resources a user can modify (own data only).

## Roles

| Role | Access | Routes |
|------|--------|--------|
| `user` | Own profile, own events, own guests, own subscriptions | `/api/v1/users/*`, `/api/v1/events/*`, `/api/v1/guests/*` |
| `admin` | All user data, platform management, override capabilities | `/api/v1/admin/*` + all user routes |

## Authorization Layers

### Layer 1: JWT Authentication Middleware

Applied to all protected routes. Validates token and extracts user context.

### Layer 2: Role Authorization Middleware

Applied to admin routes. Checks `role == "admin"` in JWT payload.

### Layer 3: Service-Level Authorization

Applied within service methods. Verifies resource ownership.

```go
func (s *EventService) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
    userID := auth.GetUserID(ctx)
    event, err := s.eventRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if event.UserID != userID && !auth.IsAdmin(ctx) {
        return nil, ErrForbidden
    }
    return event, nil
}
```

### Layer 4: Repository-Level Scoping

Repository methods always scope queries by user context.

```go
func (r *GuestRepository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Guest, error) {
    var guests []model.Guest
    err := r.db.WithContext(ctx).
        Where("event_id = ?", eventID).
        Where("deleted_at IS NULL").
        Find(&guests).Error
    return guests, err
}
```

## Ownership Rules

| Resource | Owner | Access Rule |
|----------|-------|-------------|
| Event | User who created it | `event.user_id == user_id` OR admin |
| Guest | Event owner | `guest.event.user_id == user_id` OR admin |
| RSVP | Event owner or guest | Public submission (guest), owner view (user) |
| Guestbook message | Event owner | Owner can moderate |
| Digital gift config | Event owner | `gift.event.user_id == user_id` |
| Gallery photo | Event owner | `photo.event.user_id == user_id` |
| Subscription | User | `sub.user_id == user_id` OR admin |
| Transaction | User | `tx.user_id == user_id` OR admin |

## Public Access

No authorization required for:
- Package listing (`GET /api/v1/packages`)
- Public invitation (`GET /api/v1/public/events/:slug`)
- RSVP submission (`POST /api/v1/public/rsvps`)
- Guestbook message submission (`POST /api/v1/public/guestbook`)
- Approved guestbook listing (`GET /api/v1/public/guestbook/:slug`)

## Admin Override

Admin can:
- View any user's events, guests, RSVPs
- Create/extend/terminate any subscription
- Suspend/activate any user account
- Moderate any guestbook message
- Manage packages and templates
- View all transactions

All admin overrides are logged in audit_logs.

## Implementation Pattern

```go
// In handler: extract user context
userID := auth.GetUserID(r.Context())
role := auth.GetRole(r.Context())

// In service: check ownership
func (s *EventService) Update(ctx, eventID, userID, dto) {
    event, err := s.repo.GetByID(ctx, eventID)
    if event.UserID != userID {
        return nil, ErrForbidden
    }
    // proceed with update
}
```

## Related Documentation

- `docs/security/authorization.md`
- `docs/security/authentication-security.md`
- `docs/backend/middleware.md`
- `docs/modules/admin.md`