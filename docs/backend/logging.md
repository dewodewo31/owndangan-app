# Logging

## Library Choice

The project uses **zerolog** for structured JSON logging. It provides zero-allocation, high-performance logging with a clean API. Logs are always output as JSON for production consumption by log aggregators (ELK, Datadog, Grafana Loki).

## Logger Initialisation

```go
package logger

import (
    "os"
    "time"
    "github.com/rs/zerolog"
)

func New(env string) zerolog.Logger {
    if env == "development" {
        // Human-readable console output for local dev
        return zerolog.New(zerolog.ConsoleWriter{
            Out:        os.Stdout,
            TimeFormat: time.RFC3339,
        }).With().Timestamp().Logger()
    }
    // Production: JSON output
    return zerolog.New(os.Stdout).
        With().
        Timestamp().
        Caller().
        Logger().
        Level(zerolog.InfoLevel)
}
```

## Log Levels

| Level | Method | When to Use |
|-------|--------|-------------|
| `debug` | `logger.Debug()` | Development details, verbose SQL traces. Not enabled in production. |
| `info` | `logger.Info()` | Request completion, successful operations, service start/stop. |
| `warn` | `logger.Warn()` | Unexpected but handled situations: rate limit接近, slow queries, deprecated endpoint usage. |
| `error` | `logger.Error()` | Recoverable errors: validation failures in service, external API call failures, DB connection issues. |
| `fatal` | `logger.Fatal()` | Unrecoverable errors: config loading failure, DB migration failure, port binding failure. Calls `os.Exit(1)`. |
| `panic` | `logger.Panic()` | Bugs that should never happen: nil pointer dereference, assertion failures. Calls `panic()`. |

## Structured Fields

Every log entry includes context fields. The middleware attaches request-scoped fields automatically:

```go
// Structured logging middleware
func StructuredLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            // Capture status code via wrapper
            lw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
            next.ServeHTTP(lw, r)

            logger.Info().
                Str("request_id", middleware.GetRequestID(r.Context())).
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Str("query", r.URL.RawQuery).
                Int("status", lw.status).
                Str("remote_ip", r.RemoteAddr).
                Str("user_agent", r.UserAgent()).
                Dur("duration_ms", time.Since(start)).
                Msg("request completed")
        })
    }
}
```

## Context-Level Logger

Use zerolog's `Ctx` method to attach a logger to the request context with pre-populated fields:

```go
func ContextLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            l := logger.With().
                Str("request_id", middleware.GetRequestID(r.Context())).
                Logger()
            ctx := l.WithContext(r.Context())
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Usage in handlers/services:
func (s *EventService) Create(ctx context.Context, req ...) {
    log := zerolog.Ctx(ctx)
    log.Info().Str("slug", req.Slug).Msg("creating event")
}
```

## Logging in Services

Services log decisions and errors. They never log at `Info` for normal CRUD operations (the middleware already logs every request).

```go
func (s *PaymentService) HandleWebhook(ctx context.Context, payload dto.MidtransPayload) error {
    log := zerolog.Ctx(ctx)

    log.Info().
        Str("order_id", payload.OrderID).
        Str("transaction_status", payload.TransactionStatus).
        Msg("midtrans webhook received")

    tx, err := s.txRepo.GetByOrderID(ctx, payload.OrderID)
    if err != nil {
        log.Error().Err(err).Str("order_id", payload.OrderID).Msg("transaction lookup failed")
        return fmt.Errorf("lookup transaction: %w", err)
    }
    if tx == nil {
        log.Warn().Str("order_id", payload.OrderID).Msg("unknown order_id in webhook")
        return nil // acknowledge webhook to Midtrans, skip processing
    }
    // ... process webhook
}
```

## Recommended Conventions

- Always include `request_id` in every log line.
- Use `Err(err)` for errors, not string interpolation: `log.Error().Err(err).Msg("...")`.
- Log in English, consistently.
- Do not log sensitive data: passwords, tokens, full Midtrans responses with card numbers.
- Mask PII in log entries: `log.Info().Str("email", maskEmail(email)).Msg("...")`.
- Use `Dur()` for durations, not floats: `Dur("query_ms", duration)`.
- Keep log messages static; put dynamic values in structured fields.
- Never log at `Info` inside a hot loop — use `Debug` for high-frequency events.
