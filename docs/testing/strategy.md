# Testing Strategy

## Overview

This project follows a 4-tier testing pyramid: **Unit → Integration → API → E2E**.
Each tier has a specific scope, speed, and feedback cadence.

## Testing Pyramid

```
         /\
        /  \         E2E (few, slow, high confidence)
       /    \
      / API  \       API/contract (some, medium)
     /        \
    / Integr.  \     Integration (many, fast)
   /            \
  /    Unit      \   Unit (most, instant)
 /________________\
```

### 1. Unit Tests (fastest, most numerous)

- **Scope**: Single function, method, or pure computation.
- **What to test**:
  - Service business logic (validation, calculations, state transitions).
  - Utility/helper functions.
  - Repository query builders (without DB).
- **What NOT to test**: HTTP handlers, database calls, external APIs.
- **Mocking**: Use `testify/mock` for repository interfaces in service tests.
- **Target**: >70% code coverage.

### 2. Integration Tests

- **Scope**: Interaction between two adjacent layers (service + repository, handler + service).
- **What to test**:
  - Repository against a real PostgreSQL test database.
  - Service with real repository (no mocks) for critical paths.
- **Database**: Ephemeral test database created per test suite; migrations run once.
- **Target**: Cover all data-access paths and edge cases.

### 3. API / Contract Tests

- **Scope**: Full HTTP request → handler → service → repository → response.
- **What to test**:
  - Request validation (missing fields, bad types, auth headers).
  - Response status codes, body shape, error payloads.
  - Contract conformance (OpenAPI spec vs actual response).
- **Tooling**: `httptest` for in-process server, `supertest`-style helpers.
- **Target**: Every endpoint, every non-2xx path.

### 4. End-to-End Tests (slowest, fewest)

- **Scope**: Full user flow across frontend + backend.
- **What to test**:
  - Critical business flows: registration → payment → invitation creation → RSVP.
  - Auth flow (login, logout, protected routes).
- **Tooling**: Playwright.
- **Target**: 5-10 critical user journeys.

## Test Data

- **Seeds**: Deterministic seed data in `db/seeds/` for integration/E2E tests.
- **Factories**: Go test helpers to build domain objects on the fly.
- **Cleanup**: Test suites truncate tables before each run (never share state).

## Test Database

- A dedicated PostgreSQL database named `owndangan_test` is used.
- Migrations are run via `goose up` before the test suite.
- Each test function runs inside a **transaction** that is rolled back after the test (using `go-playground/pgTestTx` or similar pattern).
- This ensures isolation without slow truncation between every test.

## CI Integration

- Unit + Integration tests run on every PR.
- API tests run on PRs targeting `main`.
- E2E tests run nightly and before release.