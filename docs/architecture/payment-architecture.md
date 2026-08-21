# Payment Architecture

## Overview

Payment processing uses Midtrans Snap. The frontend never directly handles payment status. The backend receives webhooks from Midtrans, verifies the signature, updates the transaction, and activates the subscription.

## Payment Flow

```
User
 ↓
Select Package (user/billing)
 ↓
POST /api/v1/payments/snap
 ↓
Backend creates Midtrans transaction (Snap Token)
 ↓
Frontend launches Snap.js with Snap Token
 ↓
User completes payment in Midtrans hosted page
 ↓
Midtrans sends HTTP Webhook Notification
 ↓
POST /api/v1/webhook/midtrans
 ↓
Backend verifies signature (SHA512)
 ↓
Backend updates transaction status
 ↓
If settlement → activate subscription
 ↓
User sees activated subscription in dashboard
```

## Key Principle

> Frontend payment callback is NOT the authoritative source for subscription activation.
> Subscription activation happens ONLY after the backend receives and verifies a Midtrans notification webhook.

## Midtrans Integration Details

### Transaction Creation (Backend)

```go
// Create Snap transaction
req := &midtrans.CreateSnapRequest{
    CustomerDetails: &midtrans.CustomerDetails{
        Email: user.Email,
        Name:  user.Name,
    },
    TransactionDetails: midtypes.TransactionDetails{
        OrderID:  orderID,
        GrossAmt: package.Price,
    },
    EnabledPayments: []midtrans.PaymentType{
        midtypes.PaymentTypeQRIS,
        midtypes.PaymentTypeBankTransfer,
        midtypes.PaymentTypeGopay,
        midtypes.PaymentTypeCreditCard,
    },
}

snapToken, err := client.CreateToken(req)
```

### Order ID Format

`ORDER-{user_id_short}-{timestamp}` — must be unique and 100 chars max.

### Webhook Notification

Midtrans POSTs JSON to `POST /api/v1/webhook/midtrans`:

```json
{
  "transaction_status": "settlement",
  "order_id": "ORDER-123-20260811120000",
  "payment_type": "qris",
  "gross_amount": "99000.00",
  "signature_key": "sha512_hash_here",
  "transaction_id": "midtrans-tx-id",
  "status_code": "200"
}
```

### Signature Verification

```
signature = sha512(order_id + status_code + gross_amount + server_key)
```

Backend recomputes and compares. Mismatch → reject webhook.

### Transaction Status Mapping

| Midtrans Status | Transaction Status | Action |
|-----------------|--------------------|--------|
| pending | pending | Wait |
| settlement | settlement | Activate subscription |
| capture | settlement | Activate subscription |
| expire | expire | Mark expired |
| deny | deny | Mark denied |
| cancel | cancel | Mark cancelled |
| refund | refund | Deactivate subscription (if active) |
| failure | deny | Mark denied |

## Idempotency

- Webhook notifications may be delivered multiple times.
- The `order_id` uniqueness in the transactions table prevents duplicate processing.
- If a transaction is already in `settlement` state, re-processing a settlement webhook is a no-op.
- Subscription activation is idempotent (only activate if not already active for that transaction).

## Related Documentation

- `docs/modules/payments.md`
- `docs/internal/integrations/midtrans.md`
- `docs/security/payment-security.md`
- `docs/security/webhook-security.md`
- `docs/testing/payment-testing.md`
- `docs/api/payments.md`