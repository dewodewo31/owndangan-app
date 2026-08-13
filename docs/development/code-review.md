# Code Review Checklist

## Architecture

- [ ] Does the change follow the project architecture (handler → service → repository)?
- [ ] Is business logic in services, not handlers or repositories?
- [ ] Are handlers thin (parse request, call service, return response)?
- [ ] Are there no unnecessary abstractions or layers?
- [ ] Does the change avoid circular dependencies between packages?
- [ ] Are new dependencies justified? Could the feature be built with existing tools?
- [ ] Is the change small enough to review easily? If not, should it be split?

## Security

- [ ] Is all external input validated (request body, query params, URL params)?
- [ ] Is authorization enforced server-side? (Never trust frontend role/permissions.)
- [ ] Are SQL queries parameterized (no string concatenation)?
- [ ] Are secrets and API keys never exposed in client bundles?
- [ ] Is payment status never trusted from frontend callbacks?
- [ ] Are JWT tokens validated (signature, expiry, issuer)?
- [ ] Is there proper rate limiting on auth endpoints?
- [ ] Are error messages free of sensitive information (stack traces, DB details)?

## Error Handling

- [ ] Are all errors checked and handled (not silently ignored)?
- [ ] Do errors returned to the client follow the standard error envelope format?
- [ ] Are internal errors logged but not exposed to the client?
- [ ] Are database errors handled gracefully (connection failures, timeouts)?
- [ ] Are edge cases handled (empty lists, nil pointers, missing records)?
- [ ] Is the idempotency key pattern used for payment webhooks?

## Testing

- [ ] Are unit tests added for new service logic?
- [ ] Are handler tests added for new endpoints?
- [ ] Are repository tests added for new queries?
- [ ] Do tests cover both success and error paths?
- [ ] Are tests isolated (no shared state, no test ordering dependencies)?
- [ ] Are mocks used correctly (asserting expected calls)?
- [ ] Are test cases deterministic?
- [ ] Are new E2E tests added for critical user flows?

## Documentation

- [ ] Is the API contract updated (OpenAPI spec, endpoint docs)?
- [ ] Are environment variables documented in `.env.example`?
- [ ] Are database schema changes documented?
- [ ] Is the module documentation updated?
- [ ] Are complex logic and business rules explained in comments?
- [ ] Is the ADR updated for significant architecture decisions?

## Performance

- [ ] Are N+1 queries avoided?
- [ ] Are database queries indexed appropriately?
- [ ] Are large payloads paginated?
- [ ] Are unnecessary API calls avoided?
- [ ] Are expensive operations (email, payment) handled asynchronously?

## Code Quality

- [ ] Does the code follow Go/TypeScript conventions for the project?
- [ ] Are variable names clear and descriptive?
- [ ] Are functions focused on a single responsibility?
- [ ] Is there dead code or commented-out code that should be removed?
- [ ] Do imports follow the project convention (stdlib, third-party, internal)?
- [ ] Does the linter pass without new warnings?