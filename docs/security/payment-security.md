# Payment Security

## Core Principle: Never Trust the Frontend

The frontend (Next.js) is untrusted for payment status. All payment confirmation **must** come from server-side verification.

### What the Frontend May Do
- Request a Snap token from the backend.
- Open Midtrans Snap popup with the token.
- Display success/failure UI after redirect.

### What the Frontend Must NOT Do
- Report payment status to the backend.
- Send `transaction_status` or `order_id` from browser to backend.
- Activate a subscription or upgrade a plan based on Snap callback.

## Architecture: Midtrans Snap (Not Core API)

### Why Snap, Not Core API
| Aspect | Snap | Core API |
|---|---|---|
| PCI DSS scope | Reduced (Midtrans handles card data) | Full PCI compliance required |
| Tokenization | Handled by Midtrans | Must implement yourself |
| 3DS handling | Automatic | Manual |
| Fraud screening | Built-in | Must implement |
| **Recommendation** | **Use this** | Avoid |

### Snap Flow

```
1. Frontend → Backend: POST /api/v1/payments/create (plan ID, event ID)
2. Backend → Midtrans: POST /v1/charge (order_id, gross_amount, customer_details)
3. Midtrans → Backend: Snap token + redirect URL
4. Backend → Frontend: { snap_token, redirect_url }
5. Frontend: Open Snap popup (snap_token)
6. User: Completes payment in Snap popup
7. Midtrans → Backend: HTTP POST webhook (transaction_status, order_id, etc.)
8. Backend: Verify webhook signature, verify status with Midtrans API, update order
9. Backend: Activate subscription / unlock feature
10. Frontend: Poll backend for status or receive via WebSocket
```

## Creating a Transaction (Server-Side)

```go
type TransactionRequest struct {
    OrderID     string `json:"order_id" validate:"required"`
    GrossAmount int64  `json:"gross_amount" validate:"required,min=1000"`
    CustomerDetails map[string]interface{} `json:"customer_details"`
}

func (s *PaymentService) CreateTransaction(ctx context.Context, req TransactionRequest) (*SnapResponse, error) {
    // 1. Generate unique order_id: "INV-" + timestamp + random suffix
    // 2. Store pending transaction in database
    // 3. Call Midtrans Snap API
    // 4. Return snap_token to frontend
}
```

### Order ID Format
`INV-{user_id}-{timestamp}-{random_4_digits}` — unique, traceable, no sequential guessing.

## Webhook as Authoritative Source

### Acceptable Transaction Statuses

| Midtrans Status | Action | Notes |
|---|---|---|
| `settlement` | ✅ Activate subscription | Payment captured |
| `capture` | ✅ Activate (if fraud_status=accept) | Credit card |
| `success` | ✅ Activate | Bank transfer confirmed |
| `pending` | ⏳ Wait | Awaiting payment |
| `deny` | ❌ Reject | Payment denied |
| `cancel` | ❌ Reject | User cancelled |
| `expire` | ❌ Reject | Payment expired |
| `failure` | ❌ Reject | Payment failed |
| `refund` | ⚠️ Revoke access | Refund processed |
| `chargeback` | ⚠️ Revoke access | Dispute filed |

### Idempotency

```go
// Before processing webhook:
// 1. Check if transaction_id already processed
// 2. If yes, return 200 OK (do not reprocess)
// 3. If no, proceed with status update
// 4. Use database transaction to atomically update status + activate subscription
```

## Transaction Status Verification

After receiving a webhook, always verify the status by calling Midtrans API:

```go
// GET /v2/{order_id}/status
// Compare webhook status with API status
// Only accept if both match AND signature is valid
```

## Testing

### Sandbox Testing
- Use Midtrans sandbox environment for all development.
- Test card numbers: `4811 1111 1111 1114` (success), `4911 1111 1111 1113` (denied).
- 3DS simulation: password `112233`.
- Never use real payment credentials in development.

### Test Scenarios
- [ ] Payment success flow (webhook received, subscription activated).
- [ ] Payment failure flow (webhook received, subscription not activated).
- [ ] Duplicate webhook (idempotency prevents double activation).
- [ ] Expired transaction (pending → expire, no activation).
- [ ] Refund after settlement (access revoked).
- [ ] Invalid signature webhook (rejected).