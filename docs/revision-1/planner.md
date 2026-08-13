# PHASE 1 — FRONTEND FOUNDATION & DESIGN SYSTEM

## OBJECTIVE

Implement Phase 1 dari frontend migration plan berdasarkan hasil audit yang telah selesai.

Target Phase 1:

1. Design tokens
2. Tailwind configuration
3. Reusable UI component library
4. Shared layout foundations
5. Loading / error / empty states
6. Form foundations
7. Responsive foundations

**JANGAN mengerjakan Phase 2–7 pada task ini.**

---

# 1. SOURCE OF TRUTH

Sebelum coding, baca:

```text
docs/context/project-context.md
docs/context/business-context.md
docs/context/architecture-context.md
docs/context/technical-context.md
docs/context/coding-agent-context.md

docs/architecture/overview.md
docs/architecture/frontend-architecture.md

docs/frontend/components.md
docs/frontend/layouts.md
docs/frontend/ui-ux.md
docs/frontend/routing.md

docs/modules/README.md
```

Kemudian inspect implementation existing.

Audit yang sudah dilakukan menunjukkan:

* App Router sudah digunakan
* existing dashboard menggunakan `/dashboard/*`
* `/user/*` belum digunakan
* admin belum tersedia
* marketing layout belum tersedia
* public invitation belum tersedia
* reusable UI component library hampir tidak ada
* design system belum tersedia
* inline styling masih banyak

Jangan mengasumsikan audit selalu benar.

Jika menemukan perbedaan antara audit dan repository aktual, **repository aktual menjadi source of truth untuk implementation**, sementara docs menjadi source of truth untuk architecture/business requirement.

---

# 2. CRITICAL RULE

Jangan melakukan:

```text
delete dashboard
delete existing components
rewrite entire frontend
rewrite API client
rewrite authentication
create fake API
create fake business data
```

Phase ini hanya foundation.

Existing pages harus tetap dapat berjalan.

---

# 3. DESIGN SYSTEM

Buat centralized design tokens.

Minimal:

```text
Colors
Typography
Spacing
Radius
Shadows
Borders
Transitions
Breakpoints
```

Gunakan CSS variables jika sesuai dengan architecture existing.

Design language:

### SaaS

```text
Modern
Clean
Professional
Premium
Minimal
```

### Wedding

```text
Elegant
Warm
Editorial
Romantic
Premium
```

Gunakan dua typography hierarchy:

```text
UI / SaaS → Inter
Wedding / Editorial → Playfair Display
```

Jangan memaksa Playfair Display ke seluruh dashboard.

---

# 4. COLOR SYSTEM

Buat semantic tokens, bukan hardcoded colors.

Contoh conceptual token:

```text
background
foreground
muted
muted-foreground
primary
primary-foreground
secondary
secondary-foreground
accent
accent-foreground
border
input
ring
success
warning
destructive
```

Jangan menyebarkan:

```tsx
bg-[#xxxxxx]
text-[#xxxxxx]
```

ke seluruh component.

Jika project sudah memiliki warna yang digunakan existing pages, pertahankan compatibility terlebih dahulu.

---

# 5. UI COMPONENT LIBRARY

Buat reusable components.

Minimal:

```text
Button
Input
Textarea
Select
Checkbox
Switch
Label
Card
Badge
Avatar
Modal/Dialog
Dropdown
Tabs
Tooltip
Table
Pagination
Toast
Alert
Skeleton
Spinner
EmptyState
ErrorState
LoadingState
Separator
```

Struktur konseptual:

```text
components/
└── ui/
    ├── button.tsx
    ├── input.tsx
    ├── textarea.tsx
    ├── select.tsx
    ├── checkbox.tsx
    ├── switch.tsx
    ├── label.tsx
    ├── card.tsx
    ├── badge.tsx
    ├── avatar.tsx
    ├── dialog.tsx
    ├── dropdown-menu.tsx
    ├── tabs.tsx
    ├── tooltip.tsx
    ├── table.tsx
    ├── pagination.tsx
    ├── toast.tsx
    ├── alert.tsx
    ├── skeleton.tsx
    ├── spinner.tsx
    ├── empty-state.tsx
    ├── error-state.tsx
    ├── loading-state.tsx
    └── separator.tsx
```

Sesuaikan dengan existing architecture.

---

# 6. COMPONENT REQUIREMENTS

Semua UI components harus:

* TypeScript typed
* reusable
* accessible
* responsive
* composable
* support className customization
* tidak bergantung pada business logic
* tidak melakukan API call

Contoh prinsip:

```text
UI Component
    ↓
Presentation only

Business Component
    ↓
Uses UI Component

Page
    ↓
Uses Business Components
```

Jangan membuat:

```text
Button → API call
Card → fetch data
Modal → business logic
```

---

# 7. BUTTON

Button minimal mendukung:

```text
variant
size
loading
disabled
fullWidth
className
```

Variants minimal:

```text
primary
secondary
outline
ghost
destructive
link
```

Wedding-specific styling boleh dibuat melalui variant atau wrapper.

Jangan membuat button component terpisah untuk setiap page jika hanya berbeda styling.

---

# 8. FORM COMPONENTS

Buat foundation untuk:

```text
Label
Input
Textarea
Select
Checkbox
Switch
```

Support:

```text
label
description
error
disabled
required
```

Accessibility:

```text
aria-invalid
aria-describedby
label htmlFor
keyboard navigation
focus state
```

---

# 9. CARD

Card harus dapat digunakan untuk:

```text
Dashboard metrics
Pricing
Template
Invitation sections
Forms
Settings
```

Gunakan composition:

```tsx
<Card>
  <CardHeader />
  <CardTitle />
  <CardDescription />
  <CardContent />
  <CardFooter />
</Card>
```

Jangan membuat Card hanya untuk satu use case.

---

# 10. TABLE

Table wajib memperhatikan masalah existing:

> Jangan membuat seluruh halaman horizontal scroll.

Gunakan:

```text
page
└── table container
    └── horizontal overflow
```

Jika table membutuhkan width besar, hanya container table yang boleh scroll.

Support foundation untuk:

```text
header
body
row
cell
loading
empty
```

Pagination dibuat reusable.

---

# 11. LOADING STATES

Buat reusable:

```text
Skeleton
Spinner
LoadingState
```

Dashboard tidak boleh menggunakan:

```text
Loading...
```

sebagai satu-satunya loading UI jika skeleton lebih sesuai.

---

# 12. EMPTY STATES

Buat reusable EmptyState.

Support:

```text
icon
title
description
action
```

Contoh:

```text
Belum ada tamu

Tambahkan daftar tamu untuk mulai mengelola undangan Anda.

[Tambah Tamu]
```

Jangan membuat empty state palsu dengan fake data.

---

# 13. ERROR STATES

Buat reusable ErrorState.

Support:

```text
title
description
retry action
```

Jangan menampilkan:

```text
SQLSTATE
stack trace
JWT error
database error
raw API response
```

kepada user.

---

# 14. RESPONSIVE FOUNDATION

Pastikan UI components bekerja minimal pada:

```text
320px
375px
390px
768px
1024px
1280px
1440px
```

Tidak boleh ada:

```text
horizontal page overflow
```

Komponen table boleh memiliki internal horizontal scrolling.

---

# 15. ACCESSIBILITY

Semua components harus memperhatikan:

```text
Semantic HTML
Keyboard navigation
Focus visible
ARIA
Screen readers
Disabled state
Loading state
Error state
```

Jangan menggunakan:

```text
<div onClick>
```

untuk sesuatu yang seharusnya merupakan button.

---

# 16. DARK MODE

Jika docs/frontend/ui-ux.md memang mensyaratkan dark mode, siapkan foundation.

Namun:

**Jangan merusak wedding invitation aesthetic hanya demi dark mode.**

SaaS dashboard dapat mendukung:

```text
light
dark
system
```

Wedding invitation dapat memiliki theme behavior sendiri jika architecture mengizinkan.

---

# 17. LAYOUT FOUNDATION

Buat reusable layout primitives untuk:

```text
Marketing
User Dashboard
Admin Dashboard
Public Invitation
```

Namun pada Phase 1:

**Jangan membuat seluruh halaman marketing/admin/invitation.**

Hanya buat foundation/layout shell yang diperlukan.

Conceptual:

```text
components/
├── ui/
└── layouts/
    ├── marketing/
    ├── dashboard/
    ├── admin/
    └── invitation/
```

Sesuaikan dengan architecture aktual.

---

# 18. DO NOT MIGRATE ROUTES YET

Jangan mengubah:

```text
/dashboard/*
```

menjadi:

```text
/user/*
```

pada Phase 1.

Route migration akan dilakukan pada phase berikutnya setelah foundation stabil.

Jangan membuat redirect `/dashboard → /user` sekarang kecuali benar-benar diperlukan untuk mencegah broken behavior.

---

# 19. DO NOT IMPLEMENT

Pada Phase 1 jangan implement:

```text
Admin Dashboard pages
Marketing Landing Page
Pricing Page
Template Showcase
Public Wedding Invitation
Midtrans
Guest CSV import
Gallery
Digital Gift
Advanced RSVP
Token Refresh
Role-based middleware
```

Semua itu adalah phase berikutnya.

---

# 20. REFACTOR EXISTING PAGES CAREFULLY

Setelah UI library selesai, refactor beberapa existing pages secara terbatas untuk membuktikan component library bekerja.

Prioritaskan:

```text
/dashboard
/dashboard/editor
/dashboard/guests
/dashboard/rsvp
/dashboard/billing
```

Replace obvious duplicated UI dengan reusable components.

Namun jangan mengubah business behavior.

---

# 21. VALIDATION

Jalankan command yang tersedia di `package.json`.

Minimal:

```text
lint
typecheck
test
build
```

Jika command tidak tersedia, jangan membuat command baru.

Periksa:

```text
TypeScript errors
ESLint errors
Build errors
Broken imports
Hydration errors
Responsive issues
```

---

# 22. DOCUMENTATION

Update:

```text
docs/frontend/components.md
docs/frontend/layouts.md
docs/frontend/ui-ux.md
```

Dokumentasikan:

```text
Component architecture
Design tokens
Naming convention
Usage convention
Layout convention
Responsive rules
Accessibility rules
```

Jika ada architecture decision baru yang signifikan:

```text
docs/decisions/
```

buat ADR baru.

---

# 23. GIT SAFETY

Sebelum perubahan:

```text
inspect git status
```

Jangan menghapus uncommitted user changes.

Jangan melakukan:

```text
git reset --hard
git clean -fd
```

Jangan overwrite pekerjaan existing.

---

# 24. COMPLETION CRITERIA

Phase 1 dianggap selesai hanya jika:

```text
[ ] Design tokens tersedia
[ ] Tailwind/CSS architecture konsisten
[ ] UI components tersedia
[ ] Components reusable
[ ] Components typed
[ ] Components accessible
[ ] Loading states tersedia
[ ] Empty states tersedia
[ ] Error states tersedia
[ ] Table tidak menyebabkan page overflow
[ ] Responsive foundation selesai
[ ] Existing dashboard tetap berjalan
[ ] Existing editor tetap berjalan
[ ] Existing guests tetap berjalan
[ ] Existing RSVP tetap berjalan
[ ] Existing billing tetap berjalan
[ ] No TypeScript errors
[ ] No lint errors
[ ] Production build berhasil
[ ] Documentation diperbarui
```

---

# 25. FINAL REPORT

Setelah selesai berikan:

```text
PHASE 1 COMPLETION REPORT

## Implemented
...

## Components Created
...

## Components Refactored
...

## Design Tokens
...

## Layout Foundation
...

## Accessibility
...

## Responsive
...

## Existing Features Verified
...

## Documentation Updated
...

## Tests
...

## Build
...

## Known Issues
...

## Deferred To Next Phase
...

## Production Readiness
PASS / PARTIAL / FAIL
```

---

# FINAL RULE

**Implement only Phase 1.**

Jangan melanjutkan otomatis ke Phase 2.

Setelah Phase 1 selesai, berhenti dan tunggu instruksi berikutnya.

Setelah Phase 1 selesai

Baru lanjut dengan pola:

Phase 1 → audit → test → commit
                         ↓
Phase 2 → audit → test → commit
                         ↓
Phase 3 → audit → test → commit