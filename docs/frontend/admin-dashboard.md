# Admin Dashboard

## Overview

The admin dashboard (`/admin/*`) is the platform owner's control panel. It provides a centralized view of the entire platform: users, transactions, packages, templates, and analytics. All admin pages are protected by the admin layout, which requires JWT authentication with `role === 'admin'`.

## Dashboard Overview (`/admin`)

The main admin page displays a summary of key platform metrics. It is a client component that fetches data from `GET /api/v1/admin/dashboard`.

**Stats cards (top row):**
- Total registered users (count)
- Active subscriptions (count)
- Total revenue this month (formatted in Rupiah)
- Active invitations (count)

**Charts:**
- Monthly revenue trend (line chart, last 6 months)
- New user registrations (bar chart, last 30 days)
- Subscription plan distribution (pie chart: Basic vs Premium vs Pro)

**Recent activity:**
- Latest 10 transactions (type, amount, status, date)
- Latest 5 registered users

Each stat card and chart component is a separate client component for independent loading using Suspense.

## User Management (`/admin/users`)

User management page provides a searchable, paginated table of all platform users.

**Features:**
- Search by name, email, or phone number.
- Filters: role (user/admin), status (active/suspended), subscription plan.
- Sort by: name, email, created_at, last_login.
- Pagination (20 per page).
- Row actions: view user details, suspend/activate account, reset password.

**User detail modal** (or page) shows:
- User profile info (name, email, phone, role).
- Current subscription (plan, expiry, status).
- Events owned by this user (count, list).
- Transaction history.
- Account status toggle (active/suspended).

**Suspend action:**
```tsx
const handleSuspend = async (userId: string) => {
  await apiClient.put(`/admin/users/${userId}/status`, { status: 'suspended' });
  mutate(); // revalidate user list
  toast.success('Akun pengguna telah dinonaktifkan.');
};
```

## Transaction Monitoring (`/admin/transactions`)

Transaction page lists all Midtrans payment records with real-time status.

**Table columns:**
- Order ID (linked to Midtrans dashboard for reference)
- User name and email
- Package purchased (Basic / Premium / Pro)
- Amount (Rp)
- Payment method (Bank Transfer, QRIS, GoPay, etc.)
- Status (pending / settlement / expire / cancel) — color-coded badge
- Date

**Filters:**
- Status multi-select (pending, settlement, expire, cancel).
- Date range picker.
- Payment method filter.

**Actions:**
- View raw Midtrans response (JSON expandable).
- Manually trigger status check (for reconciliation when webhook is delayed).

## Package Management (`/admin/packages`)

Package management page controls the subscription plans offered to users.

**List view:**
- Table of packages (Basic, Premium, Pro, and any custom packages).
- Columns: name, price (Rp), duration, guest limit, template count, status.

**Edit form (modal or inline):**
```tsx
interface PackageForm {
  name: string;
  price: number;
  duration_days: number;
  max_guests: number;
  max_templates: number;
  max_gallery: number;
  has_music_upload: boolean;
  has_video: boolean;
  has_whatsapp_broadcast: boolean;
  has_qris: boolean;
  has_custom_domain: boolean;
  has_remove_watermark: boolean;
  is_active: boolean;
}
```

**Validation:** price must be positive, duration must be > 0, guest limit must be >= 10.

## Template Management (`/admin/templates`)

Template management page controls the invitation themes available to users.

**List view:**
- Grid of template cards with thumbnail previews.
- Each card shows: template name, group (basic/premium/pro), status (active/inactive).
- Sort by: name, group, created_at.

**Upload flow:**
1. Admin clicks "Upload Template".
2. Modal opens with: name, group selector, CSS config fields (primary color, secondary color, font), and a ZIP file upload (HTML + assets).
3. Backend validates the ZIP, extracts assets, stores them in object storage, and creates a template record.
4. The template appears in the list with an "Inactive" status until activated.

**Activation:**
- Toggle switch to activate/deactivate a template.
- Deactivated templates are hidden from users in the editor.

## Related API Documentation

See `docs/api/admin.md` for the full admin API contract.