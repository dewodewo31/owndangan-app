# API Testing

## Overview

API tests validate the full HTTP contract between frontend and backend. They run against a real backend instance (in-process or deployed) and verify request/response shape, status codes, headers, and error payloads.

## Types of API Tests

### 1. Contract Tests

Ensure the API matches the OpenAPI specification.

```go
func TestGetInvitation_ResponseShape(t *testing.T) {
  resp := doRequest(t, "GET", "/api/invitations/inv-123", nil)
  assert.Equal(t, 200, resp.StatusCode)

  var body domain.InvitationResponse
  json.NewDecoder(resp.Body).Decode(&body)

  assert.Equal(t, "inv-123", body.ID)
  assert.NotEmpty(t, body.Name)
  assert.NotEmpty(t, body.Date)
  assert.NotEmpty(t, body.CreatedAt)
  assert.True(t, body.IsValid())
}
```

- Every endpoint must have a contract test for the 200 response shape.
- Every error response must have a test verifying the error envelope `{ "error": { "code": "...", "message": "..." } }`.

### 2. Integration Tests (Full Backend)

Tests that exercise the full Go stack: router → middleware → handler → service → repository → database.

- Use `httptest` with a real router and a real test database.
- Seed known data, then make requests and assert responses.
- Cover all CRUD operations for each resource.

```go
func TestCreateInvitation_Integration(t *testing.T) {
  db := testutil.NewTestDB(t)
  app := bootstrap.NewApp(db) // builds router, handlers, services, repos

  body := `{"name":"Wedding A","date":"2025-08-15T10:00:00Z"}`
  resp := doAuthenticatedRequest(t, app.Router, "POST", "/api/invitations", body)

  require.Equal(t, 201, resp.StatusCode)
  var created domain.Invitation
  json.NewDecoder(resp.Body).Decode(&created)
  assert.NotEmpty(t, created.ID)

  // Verify it was persisted
  getResp := doAuthenticatedRequest(t, app.Router, "GET", "/api/invitations/"+created.ID, nil)
  assert.Equal(t, 200, getResp.StatusCode)
}
```

### 3. Request/Response Validation Tests

Test every validation rule on every endpoint:

| Test Case | Method | Input | Expected Status |
|-----------|--------|-------|-----------------|
| Missing required field | POST | `{}` | 400 |
| Invalid field type | POST | `{"date": "not-a-date"}` | 400 |
| Field too long | POST | `{"name": "a... (256 chars)"}` | 400 |
| Missing auth header | GET | no header | 401 |
| Expired token | GET | expired JWT | 401 |
| Wrong role | DELETE | non-admin token | 403 |
| Non-existent resource | GET | UUID that doesn't exist | 404 |
| Duplicate creation | POST | same unique fields | 409 |

## Test Helpers

```go
func doRequest(t *testing.T, router http.Handler, method, path string, body io.Reader) *http.Response {
  req := httptest.NewRequest(method, path, body)
  req.Header.Set("Content-Type", "application/json")
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)
  return w.Result()
}

func doAuthenticatedRequest(t *testing.T, router http.Handler, method, path string, body io.Reader) *http.Response {
  req := httptest.NewRequest(method, path, body)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", "Bearer "+testutil.TestToken(t))
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)
  return w.Result()
}
```

## Running API Tests

```bash
# In-process (no external dependencies)
go test ./internal/handler/... -tags=integration

# Against deployed environment
TODO: Add script for running against staging
```