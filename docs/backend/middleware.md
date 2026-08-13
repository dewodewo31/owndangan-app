# Middleware

## Middleware Stack Order

Middleware is applied in a specific order. The outer layers wrap the inner layers. Each middleware should do one thing and do it well.

```
Request
  │
  1. RequestID        ← Inject or preserve X-Request-ID
  2. StructuredLogger ← Log method, path, status, duration
  3. Recovery         ← Recover panics, log stack trace, return 500
  4. CORS             ← Set CORS headers for frontend origin
  5. RateLimiter      ← Per-IP or per-user rate limiting
  6. Auth (JWT)       ← Parse token, inject user claims into context
  7. RoleAuthorizer   ← Check role for admin-only routes
  8. ✅ Handler       ← Business logic
```

## Request ID Middleware

```go
package middleware

import (
    "context"
    "net/http"
    "github.com/google/uuid"
)

type contextKey string
const RequestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            id = uuid.New().String()
        }
        ctx := context.WithValue(r.Context(), RequestIDKey, id)
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetRequestID(ctx context.Context) string {
    id, _ := ctx.Value(RequestIDKey).(string)
    return id
}
```

## JWT Auth Middleware

```go
func Authenticate(jwtService *jwt.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenStr := extractBearerToken(r)
            if tokenStr == "" {
                response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token")
                return
            }
            claims, err := jwtService.ValidateToken(tokenStr)
            if err != nil {
                response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
                return
            }
            ctx := context.WithValue(r.Context(), auth.UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, auth.UserRoleKey, claims.Role)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## CORS Middleware

```go
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            for _, allowed := range allowedOrigins {
                if origin == allowed {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    break
                }
            }
            w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-ID")
            w.Header().Set("Access-Control-Max-Age", "86400")
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Rate Limiter Middleware

Rate limiting is applied per route group. Public endpoints get stricter limits than authenticated routes.

```go
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
    store := sync.Map{}
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := r.RemoteAddr // per-IP; use userID for authenticated routes
            counter, _ := store.LoadOrStore(key, &rateCounter{})
            rc := counter.(*rateCounter)
            rc.mu.Lock()
            now := time.Now()
            if now.After(rc.windowEnd) {
                rc.count = 0
                rc.windowEnd = now.Add(window)
            }
            rc.count++
            if rc.count > limit {
                rc.mu.Unlock()
                w.Header().Set("Retry-After", fmt.Sprintf("%.0f", rc.windowEnd.Sub(now).Seconds()))
                response.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests")
                return
            }
            rc.mu.Unlock()
            next.ServeHTTP(w, r)
        })
    }
}
```

## Recovery Middleware

Always the outermost middleware (applied last in chain, runs first):

```go
func Recovery(logger *zerolog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    logger.Error().
                        Str("request_id", GetRequestID(r.Context())).
                        Str("panic", fmt.Sprintf("%v", rec)).
                        Stack().
                        Msg("handler panic recovered")
                    response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
```
