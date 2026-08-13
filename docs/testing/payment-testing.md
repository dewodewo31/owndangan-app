# Payment Testing

## Overview

Payment testing uses **Midtrans Sandbox** environment. No real money is involved. Sandbox mode is enabled by setting `MIDTRANS_IS_PRODUCTION=false` in the backend `.env`.

## Midtrans Sandbox

- **Server Key**: Available in Midtrans dashboard → Settings → Access Keys → Sandbox.
- **Client Key**: Used by frontend for Snap.js.
- **Base URL**: `https://app.sandbox.midtrans.com/snap/v1/transactions`.
- **Dashboard**: `https://dashboard.sandbox.midtrans.com`.

## Test Cards (Midtrans Sandbox)

Use these card numbers to simulate different payment outcomes:

| Card Number | Bank | Result |
|-------------|------|--------|
| `4811 1111 1111 1114` | BCA | Success |
| `4911 1111 1111 1113` | Mandiri | Success |
| `5211 1111 1111 1117` | Permata | Success |
| `5111 1111 1111 1111` | Maybank | Success |
| `4611 1111 1111 1116` | CIMB | Success |
| `4000 0000 0000 0001` | Any | 3D Secure Required |
| `4111 1111 1111 1111` | Any | Denied |

**CVV**: Any 3 digits. **Expiry**: Any future date.

## Testing Payment Flows

### 1. Successful Payment

1. Create invoice/order on backend.
2. Get Snap token from backend `/api/payments/snap-token`.
3. Open Snap popup (mocked in tests, real in E2E).
4. Fill test card details (BCA success card).
5. Complete payment.
6. Verify backend receives transaction notification.
7. Verify order status changes to `paid`.

### 2. Failed Payment

1. Create order, get Snap token.
2. Use denied card (`4111 1111 1111 1111`).
3. Verify payment rejected.
4. Verify order status remains `pending`.

### 3. 3D Secure

1. Use 3DS card (`4000 0000 0000 0001`).
2. Complete 3DS challenge in sandbox.
3. Verify payment succeeds after 3DS.

## Webhook Simulation

Midtrans sends payment notifications to `POST /api/payments/notification`.

```go
func TestPaymentNotification_Settlement(t *testing.T) {
  payload := `{
    "transaction_status": "settlement",
    "order_id": "order-123",
    "transaction_id": "trx-abc",
    "status_code": "200",
    "gross_amount": "500000.00"
  }`

  req := httptest.NewRequest("POST", "/api/payments/notification", strings.NewReader(payload))
  req.Header.Set("Content-Type", "application/json")
  rec := httptest.NewRecorder()
  app.Router.ServeHTTP(rec, req)

  assert.Equal(t, 200, rec.StatusCode)
  // Verify order status changed in DB
  order, _ := orderRepo.FindByID(ctx, "order-123")
  assert.Equal(t, domain.OrderStatusPaid, order.Status)
}
```

- Test all transaction statuses: `settlement`, `capture`, `deny`, `cancel`, `expire`, `pending`.
- Test signature verification for incoming webhooks.
- Test invalid/unsigned payloads return 400.

## Idempotency Testing

Payment notifications may be delivered multiple times. Processing must be idempotent.

```go
func TestPaymentNotification_Idempotency(t *testing.T) {
  payload := `{"transaction_status":"settlement","order_id":"order-123",...}`
  // Send twice
  doNotification(t, payload)
  doNotification(t, payload)

  // Order should be paid only once, no duplicate charges
  order, _ := orderRepo.FindByID(ctx, "order-123")
  assert.Equal(t, domain.OrderStatusPaid, order.Status)
  assert.Equal(t, 1, countPaymentLogs(t, "order-123"))
}
```

## Failure Scenarios

- **Midtrans timeout**: Backend should retry or fail gracefully.
- **Invalid signature**: Reject with 400.
- **Duplicate order_id**: Return error, do not process.
- **Mismatched amount**: Log and reject.
- **Network failure during Snap**: Frontend should show error and allow retry.