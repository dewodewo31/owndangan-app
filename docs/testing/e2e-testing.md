# E2E Testing

## Stack

- **Tool**: Playwright with TypeScript.
- **Location**: `tests/e2e/` in the frontend project.
- **CI**: Runs nightly and on demand before production releases.

## Test Setup

### Configuration

```ts
// playwright.config.ts
export default defineConfig({
  testDir: './tests/e2e',
  timeout: 30000,
  retries: 1,
  use: {
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',
    ignoreHTTPSErrors: true,
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'mobile', use: { ...devices['Pixel 5'] } },
  ],
});
```

### Test Accounts

| Role | Email | Password | Purpose |
|------|-------|----------|---------|
| Free user | `e2e-free@test.com` | `Test123!` | Registration flow, basic features |
| Premium user | `e2e-premium@test.com` | `Test123!` | Payment flow, premium features |
| Admin | `e2e-admin@test.com` | `Test123!` | Admin panel tests |

These accounts are seeded in the E2E test database.

## Critical User Flows

### 1. Registration → Login

```ts
test('new user can register and login', async ({ page }) => {
  await page.goto('/register');
  await page.fill('[name=email]', `e2e-${Date.now()}@test.com`);
  await page.fill('[name=password]', 'Test123!');
  await page.fill('[name=confirmPassword]', 'Test123!');
  await page.click('button[type=submit]');
  await expect(page).toHaveURL('/dashboard');
  await expect(page.getByText(/selamat datang/i)).toBeVisible();
});
```

### 2. Payment → Premium Upgrade

```ts
test('user can upgrade to premium', async ({ page }) => {
  await page.goto('/login');
  await page.fill('[name=email]', 'e2e-premium@test.com');
  await page.fill('[name=password]', 'Test123!');
  await page.click('button[type=submit]');

  await page.goto('/pricing');
  await page.click('button:has-text("Pilih Premium")');
  await page.waitForSelector('#snap-container');
  // Fill Midtrans sandbox card
  // ... (handled by Playwright in iframe)
  await expect(page.getByText(/pembayaran berhasil/i)).toBeVisible();
});
```

### 3. Create Invitation

```ts
test('premium user can create invitation', async ({ page }) => {
  await loginAs(page, 'e2e-premium@test.com');
  await page.goto('/invitations/new');
  await page.fill('[name=title]', 'Wedding Test');
  await page.fill('[name=date]', '2025-12-25');
  await page.fill('[name=location]', 'Jakarta');
  await page.click('button[type=submit]');
  await expect(page.getByText(/undangan berhasil dibuat/i)).toBeVisible();
});
```

### 4. RSVP

```ts
test('guest can RSVP to invitation', async ({ page }) => {
  await page.goto('/invite/abc-123');
  await page.click('button:has-text("Hadir")');
  await page.fill('[name=guestCount]', '2');
  await page.click('button[type=submit]');
  await expect(page.getByText(/konfirmasi tersimpan/i)).toBeVisible();
});
```

## Running E2E Tests

```bash
# Start backend and frontend, then run
npm run e2e

# With Playwright UI
npx playwright test --ui

# Specific file
npx playwright test tests/e2e/payment.spec.ts

# Report
npx playwright show-report
```

## CI Pipeline

E2E tests in CI:
1. Spin up test database (Docker).
2. Run migrations and seed E2E test accounts.
3. Start backend binary.
4. Start Next.js frontend.
5. Run Playwright tests.
6. Teardown all containers.