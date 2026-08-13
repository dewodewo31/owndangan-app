# Frontend Testing

## Stack

- **Framework**: Jest + React Testing Library.
- **Component tests**: `@testing-library/react` with `user-event`.
- **API mocking**: MSW (Mock Service Worker) for handler-level mocking.
- **Auth mocking**: Custom test wrapper that provides a mock auth context.
- **Coverage target**: >60% for UI components, >80% for utility/API modules.

## Component Tests

Test behavior, not implementation. Query by role, label, or text — never by CSS class.

```tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { InvitationForm } from './InvitationForm';

test('shows validation error when name is empty', async () => {
  render(<InvitationForm />);
  await userEvent.click(screen.getByRole('button', { name: /submit/i }));
  expect(screen.getByText(/nama wajib diisi/i)).toBeInTheDocument();
});
```

### What to test in components

- Rendering with required props.
- Form validation messages (empty, invalid format, too long).
- Loading state (spinner, disabled button).
- Empty state ("no invitations yet").
- Error state (API failure banner).
- Success state (confirmation message, redirect).
- User interactions (click, type, select).

### What NOT to test

- Internal state of third-party libraries.
- Pixel-perfect layout.
- CSS animations.

## Form Tests

Use `userEvent` for realistic typing and clicking. Test the full form lifecycle:

1. Render form with initial values.
2. Type invalid data → assert validation errors appear.
3. Type valid data → assert no errors.
4. Submit → assert loading state.
5. Assert success/error callback was called.

## API Client Tests

Mock fetch/MSW handlers in `src/mocks/handlers.ts`.

```ts
import { server } from '@/mocks/server';
import { http, HttpResponse } from 'msw';

test('getInvitations returns parsed data', async () => {
  server.use(
    http.get('/api/invitations', () =>
      HttpResponse.json([{ id: '1', name: 'Test' }])
    )
  );
  const data = await getInvitations();
  expect(data).toHaveLength(1);
});
```

- Test successful response parsing.
- Test 4xx/5xx error handling (thrown exceptions, error messages).
- Test network failure (timeout, offline).

## Auth Flow Tests

Wrap components with a mock `AuthProvider`:

```tsx
function renderWithAuth(ui: ReactElement, { user = null }: { user?: User | null }) {
  return render(
    <AuthContext.Provider value={{ user, login: jest.fn(), logout: jest.fn() }}>
      {ui}
    </AuthContext.Provider>
  );
}
```

- Test protected route redirects to `/login` when unauthenticated.
- Test login form submits email/password.
- Test logout clears user state and redirects.
- Test token expiry triggers re-login flow.

## Running Tests

```bash
# All frontend tests
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage
```