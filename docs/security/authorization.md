# Authorization

## Role-Based Access Control (RBAC)

### Roles
| Role | Scope | Permissions |
|---|---|---|
| `admin` | Global | Read/write all events, users, settings. Access admin dashboard. |
| `user` | Own resources | Create/edit/delete own events, manage own guests, view own payments. |

### Middleware Chain

```
Request → CORS → CSRF → Rate Limit → Auth → Role → Ownership → Handler
```

- **Auth middleware**: Verifies JWT signature, expiry, extracts `sub` (user ID) and `role`.
- **Role middleware**: Checks if `role` claim matches required role(s).
- **Ownership middleware**: Checks if the resource being accessed belongs to the authenticated user.

### Role Middleware Implementation

```go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("role")
        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
            Code:    "FORBIDDEN",
            Message: "insufficient permissions",
        })
    }
}
```

Usage:
```go
router.GET("/api/v1/admin/users", authMiddleware, RequireRole("admin"), adminHandler.ListUsers)
router.POST("/api/v1/events", authMiddleware, RequireRole("user"), eventHandler.Create)
```

## Ownership-Based Access Control

### Principle
A user may only access resources they own. Admins may access all resources.

### Ownership Check Pattern

```go
func RequireEventOwnership(eventRepo repository.EventRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("sub")
        userRole := c.GetString("role")
        eventID := c.Param("eventID")

        // Admins bypass ownership check
        if userRole == "admin" {
            c.Next()
            return
        }

        event, err := eventRepo.FindByID(eventID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
                Code: "NOT_FOUND",
            })
            return
        }

        if event.UserID != userID {
            c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
                Code: "FORBIDDEN",
                Message: "you do not own this resource",
            })
            return
        }

        // Store event in context for downstream handlers
        c.Set("event", event)
        c.Next()
    }
}
```

### Resource Ownership Rules

| Resource | Owner Field | Admin Access |
|---|---|---|
| Events | `events.user_id` | Yes |
| Guests | `guests.event_id → events.user_id` (chain) | Yes |
| Payments | `payments.user_id` | Yes |
| Templates | `templates.user_id` (null if global) | Yes |
| Subscriptions | `subscriptions.user_id` | Yes |

### Guest Ownership (Chained)
Guest resources require a two-step check:
1. Resolve the event from `guest.event_id`.
2. Check `event.user_id == authenticated user`.

## Repository-Level Enforcement

Never expose a repository method that returns resources without a user ID filter:

```go
// Correct — scoped to user
func (r *EventRepository) FindByUserID(userID string) ([]Event, error)

// Wrong — never expose this without ownership check
func (r *EventRepository) FindAll() ([]Event, error)  // ADMIN ONLY
```

`FindAll()` must be gated behind admin role middleware.

## Authorization Testing Checklist

- [ ] User cannot access another user's events via ID guessing.
- [ ] User cannot modify another user's event.
- [ ] User cannot delete another user's event.
- [ ] Admin can access all resources.
- [ ] Unauthenticated requests return 401, not 403.
- [ ] Deleted/archived resources return 404, not 403 (don't reveal existence).
- [ ] Guest endpoints respect event ownership chain.