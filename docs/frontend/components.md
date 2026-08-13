# Components

## Component Architecture

Components are organized in `src/components/` by category. The guiding rule: **server components by default, client components only when interactivity is required**. Pages compose sections; sections compose UI primitives; UI primitives are the smallest reusable units.

## Directory Layout

```
src/components/
├── ui/               ← Shared UI primitives (Button, Input, Modal, ...)
├── forms/            ← Form components (FormField, FileUpload, Select...)
├── layout/           ← Header, Footer, Sidebar, Nav
├── invitation/       ← Invitation section components (rendered on /[slug])
├── admin/            ← Admin dashboard-specific components
└── user/             ← User dashboard-specific components
```

## UI Components (`src/components/ui/`)

Shared, themeable primitives. All are client components (`'use client'`) because they handle events and hover states. Styling uses Tailwind classes; variants are controlled by a `variant` prop.

### Button

```tsx
<Button variant="primary" size="lg" isLoading={submitting}>
  Simpan Perubahan
</Button>
```

- Variants: `primary` (indigo), `secondary` (outline), `ghost`, `danger`, `success`.
- Sizes: `sm`, `md`, `lg`.
- Props: `isLoading` (shows spinner, disables), `asChild` (renders as `next/link` when needed).

### Input, Textarea, Select

Controlled inputs with label, error message, and hint text. All accept `value`/`onChange` and forward refs.

```tsx
<FormField label="Nama Lengkap" error={errors.name}>
  <Input value={form.name} onChange={handleChange} />
</FormField>
```

### Modal

A dialog with backdrop, escape-to-close, focus trap, and scroll lock. Uses a portal (`createPortal`) to avoid stacking-context issues.

```tsx
<Modal open={confirmDelete} onClose={() => setConfirmDelete(false)}>
  <Modal.Header title="Hapus Tamu?" />
  <Modal.Body>Apakah Anda yakin ingin menghapus tamu ini?</Modal.Body>
  <Modal.Footer>
    <Button variant="secondary">Batal</Button>
    <Button variant="danger">Hapus</Button>
  </Modal.Footer>
</Modal>
```

Other UI primitives: `Badge`, `Avatar`, `Spinner`, `Skeleton`, `Toast`, `Table`, `Tabs`, `Tooltip`.

## Form Components (`src/components/forms/`)

- `FormField` — label + input + error message wrapper.
- `FileUpload` — drag-and-drop upload with preview and progress bar (used for gallery, music, QRIS).
- `PhoneInput` — Indonesian phone format mask (`+62` prefix).
- `DatePicker` — calendar picker for wedding dates (akad and resepsi).
- `RichEditor` — text editor for invitation copy.
- `SearchBar` — debounced search input for guest/admin lists.

Forms are controlled with local state or react-hook-form. Validation messages are displayed in Indonesian.

## Invitation Section Components (`src/components/invitation/`)

These render the public invitation on `/[slug]`. Each corresponds to a section in the event configuration. They are **server components** (read-only) except where interactive (RSVP form, guestbook form, music toggle).

| Component | Section | Notes |
|-----------|---------|-------|
| `HeroSection` | Cover | Couple names, date, "Buka Undangan" button |
| `OpeningSection` | Pembuka | Quran verse / greeting text |
| `CoupleSection` | Profil | Groom & bride with parents |
| `EventDetailsSection` | Acara | Akad & resepsi time/venue/map link |
| `GallerySection` | Galeri | Image grid (masonry layout) |
| `RsvpSection` | RSVP | Form (client) + attendance recap |
| `GuestbookSection` | Buku Tamu | Message list + form (client) |
| `DigitalGiftsSection` | Amplop | Bank accounts, e-wallet, QRIS |
| `MusicPlayer` | Music | Floating music toggle button |

Sections are conditionally rendered based on `event.sections` flags from the API.

## Admin Components (`src/components/admin/`)

- `StatCard` — dashboard stat with icon, value, trend indicator.
- `RevenueChart` — monthly revenue line/bar chart.
- `UserTable` — admin user list with status badges and actions.
- `TransactionTable` — payment records with status badges (pending/settlement/expire/cancel).
- `PackageEditor` — form for editing package prices and feature limits.
- `TemplateUploader` — upload and activate invitation templates.

## User Components (`src/components/user/`)

- `EventCard` — invitation card with slug, status, view count.
- `SectionForm` — editor form per invitation section.
- `GuestTable` — guest list with checkboxes, WhatsApp link buttons.
- `RsvpSummary` — attendance counts by status.
- `PlanCard` — subscription plan comparison cards for billing page.

## Composition Rules

1. Pages must not contain business logic — only layout + data fetching + component composition.
2. UI primitives must not import from `hooks/` or `providers/`.
3. Invitation section components must not import from `admin/` or `user/`.
4. Keep component files focused; extract sub-components when a file exceeds ~200 lines.
