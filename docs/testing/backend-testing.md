# Backend Testing

## Stack

- **Framework**: Go standard `testing` package.
- **Assertions**: `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`.
- **Mocks**: `github.com/stretchr/testify/mock` for interface mocking.
- **HTTP**: `net/http/httptest` for handler tests.
- **Database**: Test PostgreSQL with transaction rollback isolation.

## Directory Layout

```
backend/
  internal/
    service/     *_test.go    (mocked repos)
    handler/     *_test.go    (httptest + mocked services)
    repository/  *_test.go    (real test DB)
  testutil/
    factories.go              (domain object builders)
    db.go                     (test DB connection, migration runner)
    mock_*.go                 (generated testify mocks)
```

## Service Tests (with Mocks)

```go
func TestCreateInvitation_Success(t *testing.T) {
  mockRepo := new(mocks.InvitationRepository)
  mockRepo.On("Create", mock.Anything, mock.Anything).
    Return(&domain.Invitation{ID: "inv-1"}, nil)

  svc := service.NewInvitationService(mockRepo)
  result, err := svc.CreateInvitation(ctx, validInput)

  assert.NoError(t, err)
  assert.Equal(t, "inv-1", result.ID)
  mockRepo.AssertExpectations(t)
}
```

- One `mock.Mock` expectation per repository call.
- Test both success and error paths (validation errors, DB failures, not-found).
- Use `mock.Anything` for context; use concrete values for business inputs.

## Handler Tests (with httptest)

```go
func TestHandleCreateInvitation(t *testing.T) {
  mockSvc := new(mocks.InvitationService)
  mockSvc.On("CreateInvitation", mock.Anything, mock.Anything).
    Return(&domain.Invitation{ID: "inv-1"}, nil)

  handler := handler.NewInvitationHandler(mockSvc)
  router := chi.NewRouter()
  router.Post("/api/invitations", handler.Create)

  body := `{"name":"Test","date":"2025-12-01"}`
  req := httptest.NewRequest("POST", "/api/invitations", strings.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", "Bearer test-token")
  rec := httptest.NewRecorder()
  router.ServeHTTP(rec, req)

  assert.Equal(t, http.StatusCreated, rec.Code)
  assert.JSONEq(t, `{"id":"inv-1"}`, rec.Body.String())
}
```

- Set up a minimal router with only the endpoint under test.
- Test request validation (missing body, bad JSON, missing auth header).
- Test response codes: 200, 201, 400, 401, 403, 404, 500.

## Repository Tests (with test DB)

```go
func TestInvitationRepository_FindByID(t *testing.T) {
  db := testutil.NewTestDB(t)
  defer db.Close()

  repo := repository.NewInvitationRepository(db)
  created := seedInvitation(t, db, &domain.Invitation{...})

  found, err := repo.FindByID(ctx, created.ID)
  assert.NoError(t, err)
  assert.Equal(t, created.ID, found.ID)
}
```

- Each test suite gets its own transaction; rollback happens via `t.Cleanup`.
- Seed necessary data inline; do not depend on global seed state.
- Run migrations once per package via `TestMain`.

## Running Tests

```bash
# All tests
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```