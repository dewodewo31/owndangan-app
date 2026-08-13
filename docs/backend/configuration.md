# Configuration

## Approach

All configuration is loaded from environment variables at startup. There are no configuration files committed to the repository. The `config` package reads env vars and produces a strongly-typed struct that is passed to all components via dependency injection.

## Configuration Struct

```go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/joho/godotenv"
)

type Config struct {
    Env      string         // "development", "staging", "production"
    Port     string         // HTTP listen port (default: "8080")
    Database DatabaseConfig
    JWT      JWTConfig
    Midtrans MidtransConfig
    Storage  StorageConfig
    CORS     CORSConfig
}

type DatabaseConfig struct {
    Host            string
    Port            string
    User            string
    Password        string
    Name            string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

type JWTConfig struct {
    Secret             string
    AccessTokenExpiry  time.Duration // default: 15m
    RefreshTokenExpiry time.Duration // default: 7d
}

type MidtransConfig struct {
    ServerKey    string
    ClientKey    string
    IsProduction bool
}

type StorageConfig struct {
    Provider     string // "s3", "gcs", "local"
    Bucket       string
    Region       string
    AccessKey    string
    SecretKey    string
    Endpoint     string // custom endpoint for S3-compatible (MinIO)
    PublicURL    string
}

type CORSConfig struct {
    AllowedOrigins []string
}
```

## Loading Configuration

```go
func Load() (*Config, error) {
    // Load .env file in non-production environments
    if os.Getenv("ENV") != "production" {
        _ = godotenv.Load() // ignore error if .env does not exist
    }

    cfg := &Config{
        Env:  getEnv("ENV", "development"),
        Port: getEnv("PORT", "8080"),
        Database: DatabaseConfig{
            Host:            getEnv("DB_HOST", "localhost"),
            Port:            getEnv("DB_PORT", "5432"),
            User:            getEnv("DB_USER", "postgres"),
            Password:        getEnv("DB_PASSWORD", ""),
            Name:            getEnv("DB_NAME", "owndangan"),
            SSLMode:         getEnv("DB_SSLMODE", "disable"),
            MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
            MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
            ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour),
        },
        JWT: JWTConfig{
            Secret:             getEnv("JWT_SECRET", ""),
            AccessTokenExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
            RefreshTokenExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
        },
        Midtrans: MidtransConfig{
            ServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
            ClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
            IsProduction: getEnvBool("MIDTRANS_IS_PRODUCTION", false),
        },
        Storage: StorageConfig{
            Provider:  getEnv("STORAGE_PROVIDER", "local"),
            Bucket:    getEnv("STORAGE_BUCKET", "uploads"),
            Region:    getEnv("STORAGE_REGION", ""),
            AccessKey: getEnv("STORAGE_ACCESS_KEY", ""),
            SecretKey: getEnv("STORAGE_SECRET_KEY", ""),
            Endpoint:  getEnv("STORAGE_ENDPOINT", ""),
            PublicURL: getEnv("STORAGE_PUBLIC_URL", ""),
        },
        CORS: CORSConfig{
            AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
        },
    }

    // Validate required secrets
    if cfg.JWT.Secret == "" && cfg.Env == "production" {
        return nil, fmt.Errorf("JWT_SECRET must be set in production")
    }
    if cfg.Midtrans.ServerKey == "" && cfg.Env == "production" {
        return nil, fmt.Errorf("MIDTRANS_SERVER_KEY must be set in production")
    }
    return cfg, nil
}
```

## Helper Functions

```go
func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}

func getEnvInt(key string, fallback int) int {
    if val := os.Getenv(key); val != "" {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return fallback
}

func getEnvBool(key string, fallback bool) bool {
    if val := os.Getenv(key); val != "" {
        b, err := strconv.ParseBool(val)
        return err == nil && b
    }
    return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
    if val := os.Getenv(key); val != "" {
        d, err := time.ParseDuration(val)
        if err == nil {
            return d
        }
    }
    return fallback
}

func getEnvSlice(key string, fallback []string) []string {
    if val := os.Getenv(key); val != "" {
        return strings.Split(val, ",")
    }
    return fallback
}
```

## Required Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ENV` | No | `development` | Runtime environment |
| `PORT` | No | `8080` | HTTP server port |
| `DB_HOST` | No | `localhost` | PostgreSQL host |
| `DB_PORT` | No | `5432` | PostgreSQL port |
| `DB_USER` | No | `postgres` | Database user |
| `DB_PASSWORD` | Yes* | — | Database password (*required in staging/production) |
| `DB_NAME` | No | `owndangan` | Database name |
| `JWT_SECRET` | Yes* | — | JWT signing secret (*required in production) |
| `MIDTRANS_SERVER_KEY` | Yes* | — | Midtrans server key (*required in production) |
| `MIDTRANS_CLIENT_KEY` | Yes* | — | Midtrans client key for Snap (*required in production) |
| `STORAGE_PROVIDER` | No | `local` | File storage backend |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000` | Comma-separated allowed origins |

## Usage in main.go

```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal().Err(err).Msg("failed to load config")
    }
    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to connect to database")
    }
    // ... wire dependencies, start server
}
```

## Environment-Specific .env Files

- `.env.example` — committed to repo with placeholder values.
- `.env` — local development, gitignored.
- Staging/production: env vars set via container orchestration (Docker Compose, Kubernetes secrets).
