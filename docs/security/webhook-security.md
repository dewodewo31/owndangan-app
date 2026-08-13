# Webhook Security

## Overview

Midtrans sends payment notifications to a public webhook endpoint on the backend. This endpoint is the **sole authoritative source** for payment status updates. It must be secured against forgery and replay attacks.

## Webhook Endpoint

```
POST /api/v1/webhooks/midtrans
```

This endpoint:
- Does NOT require authentication (no JWT, no CSRF).
- Does NOT require CSRF token (exempted in middleware).
- Is rate-limited: 30 requests per minute per IP.
- Logs every request for audit.

## Signature Verification

### How Midtrans Signs Notifications

Midtrans sends a JSON payload to the webhook URL. The signature is computed as:

```
SHA512(order_id + status_code + gross_amount + server_key)
```

Where `server_key` is the **Midtrans Server Key** (not Client Key), stored in `MIDTRANS_SERVER_KEY` environment variable.

### Verification Implementation

```go
import (
    "crypto/sha512"
    "encoding/hex"
)

func VerifyMidtransSignature(payload WebhookPayload, serverKey string) bool {
    // Build the hash input: order_id + status_code + gross_amount + server_key
    hashInput := payload.OrderID + payload.StatusCode + payload.GrossAmount + serverKey

    // Compute SHA512 hash
    hash := sha512.Sum512([]byte(hashInput))
    computedSignature := hex.EncodeToString(hash[:])

    // Constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare(
        []byte(computedSignature),
        []byte(payload.SignatureKey),
    ) == 1
}
```

### Signature Fields from Payload

| Field | Type | Example |
|---|---|---|
| `order_id` | string | `INV-usr_abc-1712345678-4829` |
| `status_code` | string | `200` |
| `gross_amount` | string | `150000.00` |
| `signature_key` | string | `a1b2c3...` (SHA512 hex) |
| `server_key` | env var | Not in payload |

**Important**: `gross_amount` is a string with decimal format (e.g. `"150000.00"`). Use the raw string value from the payload, do not format or trim it.

## IP Whitelist

Midtrans sends webhooks from known IP ranges. Verify the source IP:

```go
func IsMidtransIP(remoteIP string) bool {
    allowedCIDRs := []string{
        "103.10.63.0/24",   // Midtrans production
        "103.10.124.0/24",  // Midtrans production
        "35.247.180.176/28", // Midtrans sandbox
        "34.96.192.0/20",   // Midtrans sandbox
    }
    for _, cidr := range allowedCIDRs {
        _, block, _ := net.ParseCIDR(cidr)
        if block.Contains(net.ParseIP(remoteIP)) {
            return true
        }
    }
    return false
}
```

**Note**: IP ranges may change. Verify against Midtrans documentation and update as needed. Consider using a firewall rule instead of (or in addition to) application-level filtering.

## Idempotency Handling

### Why Idempotency Matters
Midtrans may send the same webhook notification multiple times (network retries, duplicate delivery). Processing the same webhook twice would:
- Double-activate a subscription.
- Send duplicate confirmation emails.
- Create duplicate audit log entries.

### Implementation

```go
func (s *WebhookService) ProcessNotification(ctx context.Context, payload WebhookPayload) error {
    // 1. Extract transaction_id from payload
    transactionID := payload.TransactionID

    // 2. Check if already processed
    processed, err := s.idempotencyRepo.IsProcessed(ctx, transactionID)
    if err != nil {
        return err
    }
    if processed {
        // Already processed — return 200 OK, do nothing
        return nil
    }

    // 3. Process in database transaction
    tx := s.db.Begin()
    defer tx.Rollback()

    // 3a. Mark as processed (first thing in the transaction)
    err = s.idempotencyRepo.MarkProcessed(ctx, tx, transactionID)
    if err != nil {
        return err
    }

    // 3b. Update transaction status
    err = s.paymentRepo.UpdateStatus(ctx, tx, payload.OrderID, payload.TransactionStatus)
    if err != nil {
        return err
    }

    // 3c. Activate subscription if settlement
    if payload.TransactionStatus == "settlement" {
        err = s.subscriptionRepo.Activate(ctx, tx, payload.OrderID)
        if err != nil {
            return err
        }
    }

    tx.Commit()
    return nil
}
```

### Idempotency Table

```sql
CREATE TABLE webhook_idempotency (
    transaction_id VARCHAR(100) PRIMARY KEY,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    order_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL
);
```

## Additional Security Measures

### Payload Validation
- Verify all required fields are present before processing.
- Reject payloads with unexpected `order_id` format.
- Validate `gross_amount` is a positive number.

### Response
- Always return `200 OK` for any payload (even invalid) to prevent Midtrans from retrying endlessly.
- Log the reason for rejection internally, but return a generic success response.

### Monitoring
- Alert on: signature verification failures, unknown IP addresses, unexpected status codes.
- Log every webhook payload (minus `server_key`) for forensic analysis.
- Set up Grafana dashboard showing webhook volume, success rate, processing latency.

## Testing

### Sandbox Webhooks
- Midtrans dashboard provides a "Webhook Simulator" for testing.
- Use ngrok or similar to expose local server during development.
- Simulate all statuses: `settlement`, `deny`, `expire`, `cancel`, `refund`, `chargeback`.

### Test Checklist
- [ ] Valid signature: webhook is processed.
- [ ] Invalid signature: webhook is rejected (logged, 200 returned).
- [ ] Duplicate webhook: processed once, subsequent calls return 200 without side effects.
- [ ] Unknown IP: rejected (if IP whitelist enforced).
- [ ] Missing fields: rejected gracefully.
- [ ] Concurrent duplicate: database transaction prevents race condition.