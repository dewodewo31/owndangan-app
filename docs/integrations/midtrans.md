# Midtrans Integration

## Overview

Midtrans is the payment gateway for the platform. We use **Midtrans Snap** (not Core API) to minimize PCI DSS scope. Snap handles card data capture, 3DS authentication, and fraud screening on Midtrans's side.

## Environment Configuration

| Variable | Description | Example |
|---|---|---|
| `MIDTRANS_SERVER_KEY` | Secret key for API calls (server-side only) | `SB-Mid-server-abc123...` |
| `MIDTRANS_CLIENT_KEY` | Public key for Snap frontend (safe to expose) | `SB-Mid-client-abc123...` |
| `MIDTRANS_IS_PRODUCTION` | Toggle sandbox/production | `false` (dev), `true` (prod) |
| `MIDTRANS_SANDBOX_BASE_URL` | Sandbox API endpoint | `https://api.sandbox.midtrans.com` |
| `MIDTRANS_PROD_BASE_URL` | Production API endpoint | `https://api.midtrans.com` |

**Never expose `MIDTRANS_SERVER_KEY` to the frontend.** Only `MIDTRANS_CLIENT_KEY` is used in the browser.

## Snap Transaction Flow

### 1. Backend Creates Transaction

```go
func (s *PaymentService) CreateSnapTransaction(ctx context.Context, user *User, plan *Plan) (*SnapResponse, error) {
    orderID := fmt.Sprintf("INV-%s-%d-%04d", user.ID, time.Now().Unix(), rand.Intn(10000))

    payload := map[string]interface{}{
        "transaction_details": map[string]interface{}{
            "order_id":     orderID,
            "gross_amount": plan.Price,  // int64 (in IDR, no decimals)
        },
        "customer_details": map[string]interface{}{
            "first_name": user.Name,
            "email":      user.Email,
            "phone":      user.Phone,
        },
        "credit_card": map[string]interface{}{
            "secure": true,  // Enable 3DS
        },
    }

    // POST /v1/charge to Midtrans Snap API
    resp, err := s.midtransClient.CreateTransaction(payload)
    if err != nil {
        return nil, fmt.Errorf("midtrans create transaction: %w", err)
    }

    // Store pending transaction in DB
    s.paymentRepo.Create(ctx, &Payment{
        OrderID:           orderID,
        UserID:            user.ID,
        PlanID:            plan.ID,
        GrossAmount:       plan.Price,
        Status:            "pending",
        SnapToken:         resp.Token,
        SnapRedirectURL:   resp.RedirectURL,
    })

    return &SnapResponse{
        Token:       resp.Token,
        RedirectURL: resp.RedirectURL,
    }, nil
}
```

### 2. Frontend Opens Snap

```javascript
// Next.js — client-side only
const snap = window.snap; // Loaded from Midtrans script
snap.pay(token, {
  onSuccess: function(result) { /* UI only — do NOT send to backend */ },
  onPending: function(result) { /* UI only */ },
  onError: function(result)  { /* UI only */ },
  onClose:  function()       { /* UI only */ },
});
```

### 3. Midtrans Sends Webhook

See [webhook-security.md](../security/webhook-security.md) for webhook handling.

## Webhook Notification Format

```json
{
  "transaction_time": "2025-01-15 10:30:25",
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "transaction_status": "settlement",
  "status_code": "200",
  "status_message": "Success, transaction is found",
  "order_id": "INV-usr_abc-1712345678-4829",
  "gross_amount": "150000.00",
  "payment_type": "bank_transfer",
  "signature_key": "sha512_hex_hash",
  "fraud_status": "accept",
  "bank": "bca",
  "va_numbers": [{"bank": "bca", "va_number": "12345678901"}]
}
```

## Status Mapping

| Midtrans Status | Internal Status | Action |
|---|---|---|
| `settlement` | `paid` | Subscription activated, email sent |
| `capture` | `paid` | (if fraud_status=accept) Subscription activated |
| `pending` | `pending` | Wait for completion |
| `deny` | `failed` | Notify user payment denied |
| `cancel` | `cancelled` | Notify user payment cancelled |
| `expire` | `expired` | Notify user payment expired |
| `failure` | `failed` | Notify user payment failed |
| `refund` | `refunded` | Subscription revoked, notify user |
| `chargeback` | `chargeback` | Subscription revoked, notify user |

## Sandbox vs Production

### Switching Environments

```go
func (c *MidtransClient) baseURL() string {
    if c.isProduction {
        return "https://api.midtrans.com"
    }
    return "https://api.sandbox.midtrans.com"
}
```

### Sandbox Test Cards

| Card Number | Bank | Result |
|---|---|---|
| `4811 1111 1111 1114` | BCA | Success |
| `4911 1111 1111 1113` | BCA | Denied |
| `5211 1111 1111 1117` | Mandiri | Success |
| `5111 1111 1111 1118` | Mandiri | Denied |
| 3DS PIN: `112233` | — | 3DS simulation |

## Error Handling

| HTTP Status | Meaning | Action |
|---|---|---|
| `200` | Success | Process response |
| `400` | Bad request (invalid param) | Fix payload |
| `401` | Unauthorized (invalid server key) | Check credentials |
| `402` | Cannot process (fraud/etc) | Log and notify |
| `4xx` | Client error | Log full response |
| `5xx` | Midtrans server error | Retry with backoff (max 3) |

## Testing Checklist

- [ ] Snap popup opens with correct amount.
- [ ] Payment success → webhook received → subscription activated.
- [ ] Payment failure → webhook received → subscription not activated.
- [ ] Expired transaction → webhook received → status updated.
- [ ] Refund → webhook received → subscription revoked.
- [ ] Invalid signature webhook → rejected.
- [ ] Sandbox → production switch works (env var toggle).
- [ ] No server key exposed to frontend.