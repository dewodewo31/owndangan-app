# UI/UX Guidelines

## Design System

The frontend uses **Tailwind CSS** with a design token approach. All colors, fonts, spacing, and shadows are defined in `tailwind.config.ts` so the UI remains consistent across the landing page, dashboards, and public invitations.

## Tailwind Configuration

```ts
// tailwind.config.ts
export default {
  content: ['./src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#6366f1', 50: '#eef2ff', ... 900: '#312e81' },
        accent: { DEFAULT: '#f59e0b' },
        surface: { DEFAULT: '#ffffff' },
      },
      fontFamily: {
        sans: ['var(--font-inter)', 'system-ui', 'sans-serif'],
        display: ['var(--font-playfair)', 'serif'],
      },
    },
  },
};
```

**Rules:**
- Do NOT hardcode colors in components. Use design tokens (`bg-primary`, `text-slate-600`).
- Do NOT add new color tokens without updating this config file.
- Use `font-display` (Playfair Display) for headings and `font-sans` (Inter) for body text.
- Dark mode: toggled via `darkMode: 'class'`, driven by a ThemeProvider that respects `prefers-color-scheme`.

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

## Content Loading

Use `next/image` for all images with explicit width/height to prevent CLS (Cumulative Layout Shift). Lazy-load images below the fold. Gallery images use `loading="lazy"` and `decoding="async"`.