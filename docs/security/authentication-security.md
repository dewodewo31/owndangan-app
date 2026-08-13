# Authentication Security

## Password Hashing (bcrypt)

### Configuration
- Algorithm: bcrypt
- Cost factor: **12** (approx 250ms per hash on modern hardware)
- Library: `golang.org/x/crypto/bcrypt`

### Implementation Rules

```go
// Correct — cost 12, one-way hash
hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

// Verification — constant-time comparison
err = bcrypt.CompareHashAndPassword(hash, []byte(inputPassword))
```

- Never store plaintext passwords. Ever.
- Never log passwords or password hashes.
- Never use MD5, SHA1, or SHA256 for password storage.
- On login, compare hash then clear both variables from memory.
- Password minimum length: 8 characters. Maximum: 128 characters (bcrypt truncates at 72 bytes).

### Password Policy
- Minimum 8 characters, recommended 12+.
- No composition rules (uppercase/special char requirements are counterproductive; encourage longer passphrases).
- Rate-limited: 5 failed attempts per minute per IP, 10 per hour per email.
- Account lockout: after 10 consecutive failures, lock for 15 minutes.

## JWT Signing and Expiration

### Token Format
- Algorithm: HS256 (HMAC-SHA256)
- Secret: 256-bit (32-byte) random value from `JWT_SECRET` environment variable
- Library: `github.com/golang-jwt/jwt/v5`

### Access Token
```go
claims := jwt.MapClaims{
    "sub": user.ID,
    "role": user.Role,
    "iat": time.Now().Unix(),
    "exp": time.Now().Add(15 * time.Minute).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
```

- Expiry: **15 minutes** (absolute maximum).
- Include `sub` (user ID), `role` (admin/user), `iat` (issued at).
- Never include password, email, or PII in the token.
- Validate `exp`, `iat`, and signature on every request.

### Refresh Token
- Expiry: **7 days**.
- Stored as a bcrypt hash in the database (not raw).
- Includes a `token_family` UUID to enable rotation detection.
- Returned as HttpOnly, Secure, SameSite=Strict cookie, not in response body.

### Refresh Token Rotation
```
1. Client sends refresh token + access token (both expired).
2. Server looks up refresh token hash in DB.
3. Server verifies hash matches, checks token_family.
4. Server issues new access token + new refresh token.
5. Server invalidates old refresh token (set `used_at` timestamp).
6. If a used refresh token is ever presented again → possible theft → invalidate entire token_family.
```

- On rotation detection (reuse of a revoked token): invalidate all tokens in the family, require re-login.
- Log the security event for audit.

## Token Storage (Client-Side)

- Access token: memory only (React state/context). Never localStorage.
- Refresh token: HttpOnly, Secure, SameSite=Strict cookie.
- On page load: no token in memory → redirect to login.
- On token refresh: new access token returned in response body, new refresh token in Set-Cookie header.

## Endpoints

| Endpoint | Method | Auth | Rate Limit |
|---|---|---|---|
| `/api/v1/auth/register` | POST | None | 3/min per IP |
| `/api/v1/auth/login` | POST | None | 5/min per IP, 10/hr per email |
| `/api/v1/auth/refresh` | POST | Refresh cookie | 10/min per IP |
| `/api/v1/auth/logout` | POST | Bearer | 10/min per IP |

## Logout Flow
1. Client calls `/api/v1/auth/logout` with access token.
2. Server clears refresh token cookie (`Max-Age=0`).
3. Server marks refresh token as revoked in database.
4. Client clears access token from memory.