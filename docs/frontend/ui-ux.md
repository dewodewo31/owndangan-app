# UI/UX Guidelines

## Design System

The frontend follows [`docs/design.md`](../design.md). It is built on **Material 3** foundations
with the OWNDANGAN product identity (**Deep Plum / Rose Gold**) expressed through semantic design
tokens.

All colors, fonts, radius, and shadows are defined as M3 semantic tokens in `src/globals.css`
(Tailwind v4 `@theme`). Components must reference tokens (`bg-primary`, `bg-primary-container`,
`text-on-surface-variant`, `from-plum to-rosegold`) and never hard-code colors.

## Tokens

- **Color:** full M3 semantic set — `primary`/`primary-container`/`on-primary-container`,
  `secondary`/`secondary-container`, `tertiary`/`tertiary-container`, `surface` variants
  (`surface-container-low/container/container-high`, `surface-variant`), `on-surface-variant`,
  `outline`/`outline-variant`, plus `error`, `success`, `warning` with containers. Light + dark
  palettes via `:root` and `.dark`.
- **Product brand:** `--color-plum` (Deep Plum `#4a1b33`) and `--color-rosegold` (Rose Gold
  `#a75e6e`) are fixed, theme-independent brand colors for decorative panels/gradients that carry
  white text (`from-plum to-rosegold`). Interactive UI uses the semantic M3 tokens instead.
- **Shape:** small = `8px` inputs/chips/buttons (`rounded-lg`), medium = `16px` cards/dialogs
  (`rounded-2xl`), extra-large = `28px` panels/sheets.
- **Elevation:** `shadow-elevation-0..3` mapping to the M3 0/1/3/6dp levels.
- **Typography:** `--text-display/heading/title/body/label` scale (57/22/16/14/11px).

## Typography

- **UI font:** **Plus Jakarta Sans** (loaded site-wide in `app/layout.tsx`, fallback to system sans).
- **Invitation fonts:** Playfair Display and Cormorant Garamond are loaded for templated public
  invitations; each invitation template owns its decorative pair (see `src/templates/`).
- Old guidance to use Inter + Playfair-for-headings is obsolete.

## Responsive Design

Mobile-first approach. All pages are designed for a 375px viewport first, then scaled up.

**Breakpoint strategy:**

| Breakpoint | Usage |
|-----------|-------|
| `sm` (640px) | Phones → tablets. Sidebar switches to collapsible. |
| `md` (768px) | Tablets. Invitation grids go from 1 to 2 columns. |
| `lg` (1024px) | Laptops. Dashboards show full sidebar + content. |
| `xl` (1280px) | Desktops. Max-width containers. |

**Patterns:**
- Grids: `grid-cols-1 sm:grid-cols-2 lg:grid-cols-3` for cards.
- Tables: horizontally scrollable on mobile (`overflow-x-auto`).
- Sidebars: hidden behind a hamburger toggle below `md`.
- Buttons: full-width on mobile (`w-full sm:w-auto`) for thumb-friendly targets.
- Touch targets: minimum 44px height for all tappable elements.

## Indonesian Market UX Patterns

The platform targets the Indonesian market. Several UX patterns are specific to this audience:

- **Language:** All UI copy is in Indonesian. Button labels: "Simpan", "Batal", "Hapus", "Publikasikan", "Buka Undangan".
- **Date & time format:** `dd MMMM yyyy` (e.g., "15 Juni 2025") and 24-hour time (`09.00 WIB`). Use `Intl.DateTimeFormat('id-ID')`.
- **Currency format:** Rupiah — `Rp 149.000`. Use `Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' })`.
- **Phone numbers:** Displayed as `+62 812-3456-7890`. The guest form accepts both `08xx` and `62xx` formats, normalized to `62` before storage.
- **WhatsApp-first communication:** Wherever possible, actions open WhatsApp links (e.g., "Kirim Undangan via WhatsApp") rather than email.
- **WIB timezone:** The default timezone for wedding dates is Asia/Jakarta (WIB), with WITA and WIT selectable.
- **Religious context:** Wedding invitations in Indonesia often include Islamic or religious greeting phrases. The opening section supports this.
- **Polite tone:** Use formal Indonesian address terms (Bapak/Ibu/Saudara/i) in invitation copy and confirmation messages.

## Loading States

Every data-fetching UI has a defined loading state:

- **Skeletons:** Dashboard tables and cards use skeleton placeholders that match the final layout dimensions (prevent layout shift). Use `<Skeleton className="h-4 w-24" />`.
- **Spinners:** Buttons show a spinner inside when `isLoading` (replace label content, keep width stable).
- **Page loading:** Each segment can have `loading.tsx` (Suspense fallback). The `[slug]` invitation shows a skeleton of couple names + hero image.
- **Initial page load:** Avoid spinners for full pages; prefer skeletons to reduce perceived latency.

## Error States

- **API errors:** Catch `ApiError` and display a friendly message. Use toast notifications for background actions (save, delete, publish).
- **Inline validation errors:** Form fields show red borders and Indonesian messages below the field ("Nama wajib diisi", "Format email tidak valid").
- **Empty states:** Every list has a custom empty state with an illustration, a short explanation, and a CTA. Example: "Belum ada tamu. Tambahkan tamu pertama Anda." with a "Tambah Tamu" button.
- **404 (notFound):** The `[slug]` page renders a friendly 404: "Undangan tidak ditemukan atau telah dihapus."
- **Error boundary:** `app/error.tsx` shows "Terjadi kesalahan. Silakan coba lagi." with a retry button.

## Notifications (Toasts)

- Success: green, "Perubahan berhasil disimpan."
- Error: red, "Gagal menyimpan. Silakan coba lagi."
- Warning: amber, "Masa aktif paket Anda akan berakhir."
- Info: blue, "Proses import sedang berjalan."

Toasts auto-dismiss after 4 seconds (8 for errors). Stacked at bottom-right on desktop, top-center on mobile.

## Accessibility

- All form fields have associated `<label>` elements.
- Modals are focus-trapped with escape-to-close.
- Icon buttons have `aria-label` attributes.
- Color contrast meets WCAG AA (minimum 4.5:1 for text).
- Focus outlines are visible (`focus-visible:ring-2`).
- All interactive elements are reachable via keyboard.
- Animations respect `prefers-reduced-motion`.

## Content Loading

Use `next/image` for all images with explicit width/height to prevent CLS (Cumulative Layout Shift). Lazy-load images below the fold. Gallery images use `loading="lazy"` and `decoding="async"`.

## Dashboard

The user dashboard follows the §17 hierarchy of `docs/design.md`: time-aware greeting, overview
stat cards (Undangan / Tamu / RSVP / Dilihat), "Undangan Kamu" invitation cards with status and
Edit/Pratinjau actions, and a "Aktivitas Terbaru" list. Avoid a generic enterprise admin look —
the goal is helping the wedding owner complete and manage an invitation.