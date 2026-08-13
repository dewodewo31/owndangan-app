# MASTER PROMPT — Production Implementation Phase 1–17

Saya sedang membangun platform **SaaS Undangan Pernikahan Digital & Cetak Indonesia Edition**.

Gunakan seluruh dokumentasi project, PRD, `AGENTS.md`, architecture documentation, module documentation, API documentation, database documentation, dan context documentation sebagai source of truth.

Tujuan akhir:

> Menghasilkan aplikasi yang benar-benar **production-ready**, secure, tested, maintainable, observable, scalable secara wajar, dan siap dideploy ke production.

Jangan hanya membuat prototype atau mockup.

Jangan membuat fake implementation.

Jangan menggunakan dummy API response untuk membuat test terlihat berhasil.

Jangan menandai phase sebagai selesai jika exit criteria belum terpenuhi.

---

# 0. PRINCIPLE UTAMA

Ikuti prinsip:

```text
Understand
    ↓
Inspect
    ↓
Plan
    ↓
Implement
    ↓
Test
    ↓
Audit
    ↓
Document
    ↓
Verify
    ↓
Next Phase
```

Untuk SETIAP phase:

1. Baca dokumentasi yang relevan.
2. Inspect repository dan implementation yang sudah ada.
3. Identifikasi kondisi aktual.
4. Bandingkan dengan requirement.
5. Buat implementation plan.
6. Implementasikan.
7. Buat/update tests.
8. Jalankan tests.
9. Jalankan lint/static analysis.
10. Periksa regression.
11. Update documentation.
12. Lakukan security review jika relevan.
13. Lakukan architecture review.
14. Pastikan exit criteria terpenuhi.
15. Baru lanjut ke phase berikutnya.

---

# 1. GLOBAL RULES

## 1.1 Jangan mengarang

Jika requirement belum jelas:

```text
UNKNOWN
```

Jika keputusan belum dibuat:

```text
TODO
```

Jika sebuah solusi adalah rekomendasi:

```text
PROPOSED
```

Jangan mengubah `PROPOSED` menjadi `CONFIRMED` tanpa dasar.

---

# 1.2 Jangan melakukan perubahan tidak terkait

Jika task Phase 5 adalah payment:

Jangan melakukan refactor besar frontend yang tidak diperlukan.

Gunakan:

> Smallest safe change.

---

# 1.3 Jangan bypass architecture

Backend mengikuti:

```text
HTTP Handler
      ↓
Service
      ↓
Repository
      ↓
Database
```

Handler tidak boleh menjadi tempat business logic utama.

---

# 1.4 Jangan percaya frontend

Frontend tidak boleh menjadi authority untuk:

* payment settlement
* subscription activation
* authorization
* ownership
* feature entitlement
* administrative permission

Semua harus diverifikasi server-side.

---

# 1.5 Jangan hardcode secret

Dilarang:

```text
password
API key
JWT secret
Midtrans server key
database password
storage credential
```

di source code.

Gunakan environment/configuration.

---

# 1.6 Jangan menghapus test

Jika test gagal:

```text
Investigate
→ Fix root cause
→ Re-run
```

Jangan:

```text
Delete test
Skip test
Disable test
```

hanya agar pipeline hijau.

---

# 1.7 Documentation wajib sinkron

Jika implementation berubah:

```text
Code
API
Database
Architecture
Module documentation
```

harus dievaluasi.

---

# 1.8 Production quality

Semua implementation harus mempertimbangkan:

* security
* performance
* observability
* maintainability
* error handling
* validation
* authorization
* testing
* migrations
* rollback
* monitoring
* logging
* backup
* recovery

---

# 2. PHASE 1 — ARCHITECTURE VALIDATION & PROJECT FOUNDATION

## Objective

Memastikan architecture dan repository siap digunakan untuk development production.

## Tasks

Audit:

```text
Repository
Backend
Frontend
Database
Environment
Documentation
Build system
Testing
```

Pastikan struktur backend dan frontend sesuai dokumentasi.

Backend minimal:

```text
cmd/
internal/
  config/
  handler/
  middleware/
  service/
  repository/
  model/
  validator/
  logger/
migrations/
tests/
```

Gunakan struktur aktual project jika sudah ada dan jangan melakukan restrukturisasi besar tanpa alasan.

Frontend:

```text
app/
components/
lib/
services/
hooks/
types/
```

sesuaikan dengan architecture yang sudah disepakati.

## Deliverables

* architecture audit
* dependency audit
* project structure
* environment configuration
* development scripts
* baseline tests
* baseline lint
* baseline build

## Exit Criteria

```text
[ ] Backend starts
[ ] Frontend starts
[ ] Database connection works
[ ] Migrations can run
[ ] Tests execute
[ ] Lint executes
[ ] Build succeeds
[ ] No critical architecture contradiction
```

---

# 3. PHASE 2 — DATABASE & PERSISTENCE

## Objective

Membangun PostgreSQL schema production-ready.

Implement:

```text
users
subscriptions
transactions
events
packages
package_features
templates
guests
rsvps
guestbook_entries
digital_gifts
galleries
event_sections
audit_logs
```

Gunakan entity tambahan hanya jika benar-benar dibutuhkan oleh requirement.

Pastikan:

* UUID
* foreign key
* indexes
* unique constraints
* timestamps
* soft delete jika diperlukan
* transaction boundaries
* migration strategy

Review PRD schema dan perbaiki jika schema awal tidak cukup untuk requirement.

## Deliverables

* migrations
* models
* repositories
* database indexes
* seed data untuk development/testing
* schema documentation
* ERD

## Exit Criteria

```text
[ ] Fresh database migration works
[ ] Migration rollback works where supported
[ ] Relationships tested
[ ] Unique constraints tested
[ ] Foreign keys tested
[ ] Important indexes exist
[ ] Repository tests pass
```

---

# 4. PHASE 3 — AUTHENTICATION & AUTHORIZATION

## Objective

Membangun authentication dan authorization production-ready.

Implement:

```text
Register
Login
Logout
JWT
Token expiration
Password hashing
User profile
Role
Ownership authorization
```

Roles minimal:

```text
admin
user
```

Authorization harus membedakan:

```text
Authentication
Authorization
Resource ownership
Admin privilege
```

Contoh:

User A tidak boleh mengakses:

```text
User B's invitation
User B's guests
User B's transaction
```

meskipun mengetahui ID resource.

## Security

Implement:

* secure password hashing
* JWT validation
* expiration
* brute force protection/rate limit jika infrastructure mendukung
* CORS
* input validation
* secure error messages

## Exit Criteria

```text
[ ] Register works
[ ] Login works
[ ] Invalid credentials rejected
[ ] Expired token rejected
[ ] Unauthorized resource access rejected
[ ] Admin endpoint protected
[ ] Password never returned
[ ] Authentication tests pass
[ ] Authorization tests pass
```

---

# 5. PHASE 4 — PACKAGE, SUBSCRIPTION & ENTITLEMENT

## Objective

Membangun sistem pricing dan feature entitlement.

Implement:

```text
Basic
Premium
Pro
```

Jangan membuat business logic berdasarkan banyak hardcoded:

```go
if plan == "premium"
```

Gunakan capability/feature system jika architecture memungkinkan.

Contoh:

```text
guest.max
gallery.max
video.enabled
custom_domain.enabled
watermark.removable
whatsapp.bulk.enabled
guestbook.qr.enabled
rsvp.export.enabled
```

Implement lifecycle:

```text
pending
active
expired
cancelled
```

Handle:

* expiration
* renewal
* upgrade
* downgrade
* feature access

## Exit Criteria

```text
[ ] Package CRUD/admin management
[ ] Subscription creation
[ ] Subscription activation
[ ] Expiration handled
[ ] Entitlement resolver works
[ ] Guest limits enforced
[ ] Feature restrictions enforced
[ ] Subscription tests pass
```

---

# 6. PHASE 5 — MIDTRANS PAYMENT & WEBHOOK

## Objective

Implement real Midtrans integration.

Environment:

```text
MIDTRANS_SERVER_KEY
MIDTRANS_CLIENT_KEY
MIDTRANS_IS_PRODUCTION
```

Support sandbox first.

Flow:

```text
User
 ↓
Select Package
 ↓
Create Transaction
 ↓
Backend
 ↓
Midtrans
 ↓
Snap Token
 ↓
Frontend Snap
 ↓
Payment
 ↓
Midtrans Notification
 ↓
Signature Verification
 ↓
Transaction State
 ↓
Subscription Activation
```

CRITICAL:

> Frontend callback MUST NOT directly activate subscription.

Only verified backend webhook may perform settlement logic.

Implement:

* Snap token
* transaction creation
* notification endpoint
* signature verification
* idempotency
* transaction state machine
* duplicate webhook protection
* raw notification logging
* safe retry behavior

Transaction states harus konsisten.

## Exit Criteria

```text
[ ] Sandbox transaction works
[ ] Snap token generated
[ ] Notification received
[ ] Signature verified
[ ] Invalid signature rejected
[ ] Duplicate webhook handled
[ ] Settlement activates subscription
[ ] Pending does not activate subscription
[ ] Expire handled
[ ] Cancel handled
[ ] Payment tests pass
```

---

# 7. PHASE 6 — INVITATION CORE ENGINE

## Objective

Membangun core wedding invitation.

Implement:

```text
Create invitation
Update invitation
Delete/archive invitation
Publish
Unpublish
Preview
Slug
Active status
```

Lifecycle:

```text
Draft
 ↓
Configured
 ↓
Published
 ↓
Active
 ↓
Disabled/Expired
```

Implement slug:

```text
domain.com/{slug}
```

Rules:

* globally unique
* validated
* safe characters
* collision handling
* reserved slug protection

## 9 Sections

```text
1. Cover
2. Opening
3. Couple Profile
4. Events
5. Gallery
6. RSVP
7. Digital Gift
8. Dress Code
9. Closing
```

## Exit Criteria

```text
[ ] Invitation CRUD works
[ ] Slug works
[ ] Slug uniqueness enforced
[ ] Publish/unpublish works
[ ] Ownership enforced
[ ] Invitation sections persist
[ ] API tests pass
```

---

# 8. PHASE 7 — GUEST, RSVP & GUESTBOOK

## Guest

Implement:

```text
Create
Read
Update
Delete
Import CSV
Validation
Duplicate detection
Guest limit
```

Fields dapat mencakup:

```text
name
phone
group
plus_one
invitation_token
attendance_status
```

hanya jika dibutuhkan.

## RSVP

Implement:

```text
attendance
guest_count
message
submitted_at
```

## Guestbook

Implement:

```text
submit message
moderation
publish/unpublish
```

## QR Guest

Jika Pro entitlement tersedia:

```text
generate QR
scan
check attendance
```

## Exit Criteria

```text
[ ] Guest CRUD works
[ ] Guest limit enforced
[ ] CSV import validated
[ ] Duplicate handling works
[ ] RSVP works
[ ] Guestbook moderation works
[ ] QR attendance works where entitled
[ ] Tests pass
```

---

# 9. PHASE 8 — TEMPLATE ENGINE & MEDIA

## Objective

Membangun template architecture.

Jangan hardcode:

```text
if template == "template1"
```

di seluruh application.

Gunakan:

```text
Invitation Data
      ↓
Template
      ↓
Renderer
      ↓
Public Invitation
```

Implement:

```text
Templates
Template versions
Preview
Activation
Package entitlement
Theme assets
Gallery
Music
Video embeds
```

Media upload harus melakukan:

* MIME validation
* file size validation
* extension validation
* safe filename
* storage abstraction

Jangan mengikat domain logic langsung ke provider storage tertentu.

## Exit Criteria

```text
[ ] Template selection works
[ ] Template entitlement works
[ ] Template rendering works
[ ] Template versioning works
[ ] Gallery works
[ ] Music works
[ ] Invalid upload rejected
[ ] Media tests pass
```

---

# 10. PHASE 9 — PUBLIC NEXT.JS INVITATION

## Objective

Membangun public invitation yang production-ready.

Route:

```text
/{slug}
```

Optimalkan:

* mobile
* SEO
* performance
* accessibility
* OpenGraph
* metadata
* dynamic rendering
* loading
* error
* not-found

Implement:

```text
Cover
Opening
Couple
Events
Gallery
RSVP
Gift
Dress Code
Closing
```

Pastikan public page tidak membutuhkan login.

## SEO

Generate dynamic:

```text
title
description
OpenGraph
Twitter metadata
canonical URL
```

berdasarkan invitation.

## Exit Criteria

```text
[ ] Public slug works
[ ] SSR/appropriate rendering works
[ ] SEO metadata works
[ ] OpenGraph works
[ ] Mobile layout works
[ ] 404 works
[ ] Loading/error states work
[ ] Performance acceptable
```

---

# 11. PHASE 10 — USER DASHBOARD

## Objective

Membangun dashboard pasangan pengantin.

Routes:

```text
/user
/user/editor
/user/guests
/user/rsvp
/user/billing
```

Dashboard:

```text
Invitation status
Subscription
Guest count
RSVP count
Recent messages
Payment status
```

Editor:

```text
Section 1-9
Save
Preview
Publish
```

Guests:

```text
CRUD
Import
WhatsApp generator
```

Billing:

```text
Package
Subscription
Transactions
Payment
```

## Exit Criteria

```text
[ ] Protected routes
[ ] User ownership enforced
[ ] Editor works
[ ] Preview works
[ ] Guest management works
[ ] RSVP dashboard works
[ ] Billing works
[ ] Responsive UI
```

---

# 12. PHASE 11 — ADMIN DASHBOARD

## Objective

Membangun admin platform dashboard.

Routes:

```text
/admin
/admin/users
/admin/transactions
/admin/packages
/admin/templates
```

Implement:

### Analytics

```text
users
active invitations
transactions
revenue
subscriptions
```

### Users

```text
view
activate
suspend
reset password if supported
```

### Packages

```text
CRUD
price
features
discount
active status
```

### Templates

```text
CRUD
version
preview
activation
```

### Payments

```text
transaction
status
webhook log
error
```

## Security

Semua admin API harus server-side authorized.

## Exit Criteria

```text
[ ] Admin authentication
[ ] Admin authorization
[ ] Dashboard analytics
[ ] User management
[ ] Package management
[ ] Template management
[ ] Payment monitoring
[ ] Audit logs
[ ] Admin tests pass
```

---

# 13. PHASE 12 — SECURITY HARDENING

## Objective

Melakukan security audit terhadap seluruh system.

Audit:

```text
Authentication
Authorization
API
Database
Frontend
Payment
Webhook
Uploads
Secrets
Logging
CORS
Rate limiting
```

Cari:

```text
SQL Injection
XSS
IDOR
Broken Access Control
Privilege Escalation
Sensitive Data Exposure
Mass Assignment
Improper Validation
Webhook Forgery
Replay
```

Pastikan:

* secrets tidak masuk Git
* production error tidak expose stack trace
* password tidak masuk log
* payment credential tidak masuk frontend
* user ownership enforced
* admin privilege enforced

## Exit Criteria

```text
[ ] Security audit completed
[ ] Critical vulnerabilities fixed
[ ] High vulnerabilities fixed
[ ] Secrets audit clean
[ ] Authorization audit clean
[ ] Payment security verified
[ ] Webhook security verified
```

Jika ditemukan critical vulnerability:

> STOP.

Jangan lanjut Phase 13 sebelum diperbaiki.

---

# 14. PHASE 13 — AUTOMATED TESTING & QUALITY GATE

## Objective

Membangun quality gate.

Backend:

```text
Unit
Integration
Repository
Service
Handler
Auth
Authorization
Payment
Webhook
Subscription
```

Frontend:

```text
Component
Form
Editor
API integration
Auth
Route protection
```

E2E:

```text
Register
Login
Purchase
Payment
Subscription
Invitation
Guest
RSVP
Guestbook
Admin
```

Pastikan test tidak hanya menguji happy path.

Tambahkan:

```text
invalid input
unauthorized
forbidden
not found
duplicate
expired
race/duplicate webhook
```

## Exit Criteria

```text
[ ] Tests pass
[ ] No disabled critical tests
[ ] No skipped critical tests
[ ] Coverage acceptable
[ ] Build passes
[ ] Lint passes
[ ] Static analysis passes
```

Jangan menetapkan angka coverage tertentu jika belum ditentukan project.

Gunakan coverage sebagai indikator, bukan satu-satunya quality metric.

---

# 15. PHASE 14 — CI/CD & ENVIRONMENT

## Objective

Membangun deployment pipeline.

Environment:

```text
development
staging
production
```

CI harus menjalankan:

```text
install
lint
test
build
security checks
```

CD harus memiliki:

```text
staging deployment
production deployment
rollback strategy
```

Jangan deploy production otomatis tanpa approval jika workflow project belum menetapkan demikian.

Environment secrets harus aman.

## Exit Criteria

```text
[ ] CI pipeline works
[ ] Tests run in CI
[ ] Build runs in CI
[ ] Staging deployment works
[ ] Environment separation exists
[ ] Production configuration documented
[ ] Rollback documented/tested
```

---

# 16. PHASE 15 — PERFORMANCE & LOAD TESTING

## Objective

Memastikan system dapat menangani traffic realistis.

Test area:

```text
Public invitation
API
Login
Invitation rendering
Guest RSVP
Payment webhook
Admin dashboard
```

Cari:

```text
N+1 query
slow query
missing index
unnecessary API requests
large payload
large images
memory leak
connection exhaustion
```

Optimasi:

```text
database indexes
pagination
caching
image optimization
API payload
connection pooling
Next.js rendering
```

Jangan melakukan premature optimization.

## Exit Criteria

```text
[ ] Load test completed
[ ] Bottlenecks identified
[ ] Critical bottlenecks fixed
[ ] Database queries reviewed
[ ] Public invitation performance reviewed
[ ] API performance reviewed
```

Catat hasil benchmark.

---

# 17. PHASE 16 — PRODUCTION DEPLOYMENT

## Objective

Deploy actual application ke production environment.

Pastikan:

```text
Backend
Frontend
Database
Reverse Proxy
HTTPS
DNS
Environment variables
Storage
Monitoring
Logs
```

Production deployment harus reproducible.

Implement jika sesuai infrastructure:

```text
systemd
Docker
Nginx
reverse proxy
TLS
health checks
```

Jangan mengasumsikan provider tertentu jika belum ditentukan.

## Database

Pastikan:

```text
backup
migration
rollback strategy
```

## Monitoring

Minimal:

```text
application logs
error logs
health endpoint
database monitoring
resource monitoring
```

## Exit Criteria

```text
[ ] Production backend reachable
[ ] Production frontend reachable
[ ] HTTPS works
[ ] Database connected
[ ] Migration completed
[ ] Authentication works
[ ] Invitation works
[ ] Payment configured
[ ] Webhook reachable
[ ] Logs available
[ ] Health check works
[ ] Backup configured
[ ] Rollback documented
```

---

# 18. PHASE 17 — FINAL PRODUCTION AUDIT

Ini adalah phase terakhir dan WAJIB dilakukan.

Jangan langsung menyatakan:

```text
Production Ready
```

sebelum audit.

Lakukan audit:

## Architecture

```text
[ ] Architecture consistency
[ ] Module boundaries
[ ] No unnecessary coupling
```

## Backend

```text
[ ] API
[ ] Services
[ ] Repositories
[ ] Validation
[ ] Errors
[ ] Logging
```

## Frontend

```text
[ ] Routing
[ ] Authentication
[ ] UX
[ ] Responsive
[ ] SEO
[ ] Performance
```

## Database

```text
[ ] Schema
[ ] Index
[ ] FK
[ ] Migration
[ ] Backup
```

## Payment

```text
[ ] Midtrans
[ ] Signature
[ ] Webhook
[ ] Idempotency
[ ] Subscription
```

## Security

```text
[ ] Authentication
[ ] Authorization
[ ] IDOR
[ ] XSS
[ ] Injection
[ ] Secrets
[ ] CORS
[ ] Rate limiting
```

## Testing

```text
[ ] Unit
[ ] Integration
[ ] API
[ ] E2E
[ ] Payment
[ ] Regression
```

## Operations

```text
[ ] Logging
[ ] Monitoring
[ ] Backup
[ ] Recovery
[ ] Rollback
[ ] Health check
```

## Documentation

```text
[ ] README
[ ] AGENTS.md
[ ] API docs
[ ] Module docs
[ ] Architecture docs
[ ] Deployment docs
[ ] Troubleshooting
```

---

# 19. PRODUCTION READINESS SCORE

Buat final report:

```text
Production Readiness
====================

Architecture:       PASS/FAIL
Backend:            PASS/FAIL
Frontend:           PASS/FAIL
Database:           PASS/FAIL
Authentication:     PASS/FAIL
Authorization:      PASS/FAIL
Subscription:       PASS/FAIL
Payment:            PASS/FAIL
Invitation:         PASS/FAIL
Guest:              PASS/FAIL
RSVP:               PASS/FAIL
Templates:          PASS/FAIL
Security:           PASS/FAIL
Testing:            PASS/FAIL
CI/CD:              PASS/FAIL
Performance:        PASS/FAIL
Deployment:         PASS/FAIL
Monitoring:         PASS/FAIL
Backup:             PASS/FAIL
Documentation:      PASS/FAIL
```

Gunakan status:

```text
PASS
PARTIAL
FAIL
BLOCKED
NOT TESTED
```

---

# 20. BLOCKING ISSUES

Buat:

```md
## Production Blockers

- CRITICAL:
- HIGH:
- MEDIUM:
- LOW:
```

Rules:

### CRITICAL

Tidak boleh production.

### HIGH

Tidak boleh production tanpa explicit approval.

### MEDIUM

Boleh production hanya jika risk diterima dan terdokumentasi.

### LOW

Boleh menjadi post-production backlog.

---

# 21. FINAL CHANGE REPORT

Setiap phase harus menghasilkan report:

```text
Phase:
Status:

Implemented:
- ...

Modified:
- ...

Created:
- ...

Deleted:
- ...

Tests:
- ...

Lint:
- ...

Build:
- ...

Security:
- ...

Documentation:
- ...

Known Issues:
- ...

TODO:
- ...

Next Phase:
- ...
```

---

# 22. PHASE EXECUTION RULE

JANGAN mengerjakan seluruh phase secara membabi buta dalam satu batch.

Kerjakan:

```text
Phase 1
 ↓
Validate
 ↓
Phase 2
 ↓
Validate
 ↓
Phase 3
 ↓
Validate
...
 ↓
Phase 17
```

Setelah setiap phase:

1. Run tests.
2. Run lint.
3. Run build jika relevan.
4. Review changed files.
5. Check documentation.
6. Check architecture.
7. Check security.
8. Evaluate exit criteria.

Jika exit criteria gagal:

```text
STOP
FIX
RETEST
```

Jangan lanjut.

---

# 23. CHANGE MANAGEMENT

Sebelum perubahan besar:

```text
Affected Modules:
Affected APIs:
Affected Tables:
Affected Frontend:
Affected Tests:
Affected Documentation:
Security Impact:
Performance Impact:
```

Jika database berubah:

```text
Migration required: YES
Rollback strategy:
Data migration required:
```

Jika API berubah:

```text
Breaking change: YES/NO
Frontend impact:
Backward compatibility:
```

---

# 24. NO FAKE COMPLETION

Dilarang menyatakan:

```text
Done
Complete
Production Ready
```

jika hanya:

```text
file dibuat
mock dibuat
endpoint dibuat
test belum dijalankan
integration belum diverifikasi
```

"Implemented" tidak sama dengan "Verified".

Gunakan:

```text
Implemented
Verified
Production Ready
```

sebagai status berbeda.

---

# 25. FINAL REQUIREMENT

Pada akhir Phase 17, hasil yang saya inginkan adalah:

```text
                    ┌─────────────────────┐
                    │       USERS         │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │     NEXT.JS         │
                    │                     │
                    │ Public Invitation   │
                    │ User Dashboard      │
                    │ Admin Dashboard     │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │      GO API         │
                    │                     │
                    │ Auth                │
                    │ Users               │
                    │ Invitations         │
                    │ Guests              │
                    │ RSVP                │
                    │ Subscriptions       │
                    │ Payments            │
                    │ Admin               │
                    └──────────┬──────────┘
                               │
                 ┌─────────────┼─────────────┐
                 ▼             ▼             ▼
          ┌────────────┐ ┌───────────┐ ┌────────────┐
          │ PostgreSQL │ │  Midtrans │ │   Storage  │
          └────────────┘ └───────────┘ └────────────┘
```

Target akhir:

> **Aplikasi bukan sekadar bisa berjalan, tetapi dapat dipertanggungjawabkan untuk deployment production.**

---

# 26. START NOW

Mulai dari **PHASE 1**.

Jangan mengimplementasikan Phase 2 sebelum Phase 1 memenuhi exit criteria.

Pada awal Phase 1:

1. Inspect seluruh repository.
2. Inspect seluruh documentation.
3. Inspect `AGENTS.md`.
4. Inspect PRD.
5. Inspect existing backend.
6. Inspect existing frontend.
7. Inspect database.
8. Inspect environment configuration.
9. Buat architecture gap analysis.
10. Buat implementation plan Phase 1.

Kemudian implementasikan Phase 1.

Setelah Phase 1 selesai, tampilkan:

```text
PHASE 1 REPORT

Status:
Exit Criteria:
Passed:
Failed:
Files Created:
Files Modified:
Tests:
Lint:
Build:
Security Findings:
Documentation Updated:
Known Issues:
Next Phase:
```

Jika ada blocker, **STOP dan jelaskan blocker tersebut**.

Jangan melompati blocker hanya untuk melanjutkan ke phase berikutnya.
