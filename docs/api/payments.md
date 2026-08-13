# Payments

Integrates with Midtrans Snap. Webhook-based status updates.

## Endpoints

### POST /payments/snap

Create a Midtrans Snap transaction. **Auth:** Required

**Request:** `{ "package_id": "770e8400-...", "redirect_url": "https://app.example.com/payment/finish" }`

**Response (201):**
```json
{
  "success": true, "data": {
    "transaction_id": "660e8400-...", "order_id": "INV-20250215-XXXX-ABC123",
    "snap_token": "3c5c8f52-9c0e-4a1a-8b5a-1b2c3d4e5f6a",
    "snap_redirect_url": "https://app.sandbox.midtrans.com/snap/v3/redirection/3c5c8f52-...",
    "gross_amount": 150000,
    "payment_options": { "bank_transfer": ["bca","bni","bri","mandiri"], "e_wallet": ["gopay","shopeepay","qris"] }
  },
  "meta": { "request_id": "req-abc123" }
}
```
**Errors:** VALIDATION_ERROR 422, NOT_FOUND 404 (package), INTERNAL_ERROR 500 (Midtrans API failure).

**Business rules:** Order ID format: `INV-YYYYMMDD-{short_user_hash}-{random8}`. Snap token has 1-hour expiry. Only `is_active = true` packages purchasable. Transaction created `pending` before Snap call.

---

### POST /webhook/midtrans

Midtrans payment notification. **Auth:** None (signature-verified) | **Rate limit:** IP whitelist planned

**Request (Midtrans payload):**
```json
{
  "transaction_status": "settlement", "transaction_id": "3c5c8f52-...", "order_id": "INV-20250215-XXXX-ABC123",
  "gross_amount": "150000.00", "payment_type": "bank_transfer",
  "transaction_time": "2025-02-15 08:35:00", "settlement_time": "2025-02-15 08:36:00",
  "status_code": "200",
  "signature_key": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
  "bank": "bca", "va_numbers": [{ "bank": "bca", "va_number": "12345678901" }],
  "fraud_status": "accept"
}
```
**Response (200):** `{ "success": true, "data": null, "meta": { "request_id": "req-abc123" } }`

Always returns 200 OK to acknowledge, even if processing fails.

**Signature verification:** `sha512(order_id + status_code + gross_amount + server_key)` — computed signature MUST match `signature_key`. Reject with 401 if mismatch.

**Transaction status mapping:**

| Midtrans status | Internal | Action |
|----------------|----------|--------|
| settlement/capture | settlement | Activate subscription |
| pending | pending | Wait for settlement |
| deny/cancel | deny | Mark failed |
| expire | expire | Mark expired |
| refund/partial_refund | refund | Cancel subscription |

**Business rules:** Verify signature before processing. If `order_id` not found, return 200 + log warning. If internal status already `settlement`, ignore (idempotent). On settlement: update transaction + activate subscription, calculate `expires_at` from `duration_days`. Override Free subscription on settlement. Do NOT trust frontend payment status.

---

### GET /admin/transactions

List all transactions (admin). **Query:** `?page=1&per_page=20&status=settlement&q=keyword&sort=created_at&order=desc`

**Response (200):**
```json
{
  "success": true, "data": [
    { "id": "660e8400-...", "order_id": "INV-...", "user": { "id": "550e8400-...", "name": "Andi Pratama", "email": "andi@example.com" }, "package": { "id": "770e8400-...", "name": "Premium" }, "gross_amount": 150000, "status": "settlement", "payment_type": "bank_transfer", "transaction_time": "2025-02-15T08:35:00Z", "settlement_time": "2025-02-15T08:36:00Z", "created_at": "2025-02-15T08:30:00Z" }
  ],
  "meta": { "pagination": { "page": 1, "per_page": 20, "total": 520, "total_pages": 26 }, "request_id": "req-abc123" }
}
```
**Errors:** UNAUTHORIZED 401, FORBIDDEN 403.
