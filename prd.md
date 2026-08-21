# Product Requirement Document (PRD): Platform Undangan Pernikahan Digital & Cetak (Indonesia Edition)

**Dokumen Versi:** 2.0.0  
**Tanggal:** 11 Agustus 2026  
**Status:** Approved / Updated  
**Penulis:** Product Management Team  

---

## 1. Ringkasan Eksekutif (Executive Summary)

Sistem Undangan Pernikahan Digital ini dirancang sebagai platform SaaS (*Software as a Service*) berbasis arsitektur terpisah (*decoupled/headless architecture*):
* **Backend:** Golang (RESTful API / gRPC) yang tangguh, cepat, dan terorganisir secara modular.
* **Frontend & Landing Page:** Next.js (App Router) dengan pendekatan **SEO-Friendly via Dynamic Slugs** di seluruh modul.
* **Payment Gateway:** Integrasi **Midtrans** (Sandbox mode secara default, dengan saklar konfigurasi lingkungan yang mudah dialihkan ke Production).
* **Dual Dashboard System:** 
  1. **User Dashboard (Pasangan Pengantin):** Mengelola isi undangan, daftar tamu, kirim pesan WhatsApp, dan rekapitulasi RSVP/amplop.
  2. **Admin Dashboard (Pemilik Platform/Sistem):** Mengelola pengguna, transaksi, paket berlangganan, template, serta analitik platform secara keseluruhan.

---

## 2. Struktur Arsitektur & Folder Frontend (Next.js)

Aplikasi frontend dibangun menggunakan **Next.js** yang dibagi menjadi dua area utama untuk mendukung performa, SEO, dan separasi tanggung jawab (*separation of concerns*):

```text
frontend/
├── app/
│   ├── (public)/                 # Landing Page Publik, Showcase, Pricing, Blog
│   │   ├── page.tsx
│   │   ├── pricing/
│   │   └── [slug]/              # Dynamic Slug untuk SEO-Friendly Undangan Digital (/romeo-juliet)
│   ├── admin/                    # Dashboard Pemilik Platform (Admin Root)
│   │   ├── page.tsx              # Analytics Overview
│   │   ├── users/                # Manajemen User & Pengantin
│   │   ├── transactions/         # Verifikasi & Log Midtrans
│   │   ├── packages/             # Setting Harga & Fitur Paket (Basic, Premium, Pro)
│   │   └── templates/            # Pengaturan Template Undangan
│   └── user/                     # Dashboard Pasangan Pengantin (User Root)
│       ├── page.tsx              # Overview Undangan Saya
│       ├── editor/               # Editor Isian Undangan (Section 1-9)
│       ├── guests/               # Guest List & WhatsApp Generator
│       ├── rsvp/                 # Rekap RSVP & Buku Tamu
│       └── billing/              # Pembelian Paket & Riwayat Transaksi Midtrans
```

---

## 3. Skema Paket Berlangganan (Tiered Pricing Plans)

Berdasarkan analisis pasar kompetitor undangan pernikahan digital di Indonesia (seperti AkanNikah, Satumomen, SebarIn, dll.), platform ini membagi paket ke dalam 3 tingkatan (*tier*):

| Fitur / Komponen | Paket Starter | Paket Premium | Paket All Access |
| :--- | :--- | :--- | :--- |
| **Harga Acuan** | Rp 99.000 | Rp 299.000 | Rp 999.000 |
| **Masa Aktif** | 30 Hari | 60 Hari | Selamanya / Custom |
| **Pilihan Tema** | Standard (5 Theme) | Premium (20+ Theme) | All Themes + Request Custom |
| **Batas Tamu Undangan** | Max 100 Tamu | Max 500 Tamu | Unlimited Tamu |
| **Subdomain / Slug** | Custom Slug (`/mempelai`) | Custom Slug (`/mempelai`) | Custom Domain (`pasangan.com`) / Custom Slug |
| **Background Music** | Musik Galeri Preset | Upload MP3 Sendiri | Upload MP3 Sendiri |
| **Galeri Foto & Video** | Max 5 Foto | Max 20 Foto + Video YouTube | Unlimited Foto + Video HD |
| **RSVP & Ucapan Real-Time** | Ya | Ya | Ya (+ Export Excel) |
| **Amplop Digital / QRIS** | Rekening Bank | Rekening + e-Wallet | QRIS Statis + Direct Transfer |
| **WhatsApp Broadcast** | Manual (WhatsApp Link) | Auto Broadcast / Direct Tool | Integrated WhatsApp API / Bulk Sender |
| **Buku Tamu QR Code (Lokasi)** | Tidak | Tidak | **Ya (Scan QR Tamu di Gedung)** |
| **Remove Watermark Platform** | Tidak | Ya | Ya |

---

## 4. Sistem Pembayaran: Integrasi Midtrans

### 4.1 Modul Konfigurasi (Golang Backend)
Sistem menyediakan abstraksi gateway pembayaran yang fleksibel dengan pengaturan *environment variables* (`.env`):

* `MIDTRANS_SERVER_KEY`: Server Key dari Dashboard Midtrans.
* `MIDTRANS_CLIENT_KEY`: Client Key untuk Snap JS Frontend.
* `MIDTRANS_IS_PRODUCTION`: Boolean (`false` untuk Sandbox, `true` untuk Production).

### 4.2 Flow Transaksi
1. User memilih paket di `user/billing`.
2. Backend Golang menggenerasi `Snap Token` via SDK Midtrans Golang.
3. Next.js menampilkan widget **Midtrans Snap Pop-up**.
4. Pembayaran diproses (QRIS, Bank Transfer Virtual Account, GoPay, Credit Card).
5. Midtrans mengirimkan **HTTP Webhook Notification** ke Golang Backend (`POST /api/v1/payments/midtrans-notification`).
6. Backend memverifikasi *signature key*, memperbarui status transaksi di DB, dan otomatis mengaktifkan paket user.

---

## 5. Fitur Dashboard Dual-Role

### 5.1 Dashboard Admin (`/admin`)
* **Ringkasan Ringkas & Analitik:** Grafik total pendaftar, total transaksi sukses, omzet bulanan, dan jumlah undangan aktif.
* **Manajemen Pengguna (User Management):** Akses suspend/activate account user, reset password, dan impersonate login user.
* **Manajemen Paket (Pricing Manager):** Mengubah harga, batas fitur, dan status diskon paket Basic, Premium, dan Pro.
* **Manajemen Template:** Upload & atur template baru (HTML/React-based layout) yang bisa digunakan user.
* **Audit Log & Pembayaran:** Monitoring seluruh callback Webhook Midtrans dan log error.

### 5.2 Dashboard User (`/user`)
* **Wizards & Editor Undangan:** Pengisian formulir terstruktur untuk 9 Section Undangan Pernikahan (Cover, Pembuka, Profil, Acara, Galeri, RSVP, Amplop, Dress Code, Penutup).
* **SEO & Slug Settings:** Pengaturan URL slug unik (misal: `domain.com/budi-ani`).
* **Manajemen Daftar Tamu (Guest List & WhatsApp Generator):** Import CSV, penambahan nama tamu, serta tombol cepat kirim pesan ke WhatsApp dengan template kata-kata otomatis.
* **Laporan RSVP & Buku Tamu:** Melihat siapa saja yang akan hadir, jumlah rombongan, serta moderasi ucapan doa.
* **Manajemen Amplop Digital:** Mengatur nomor rekening/QRIS dan melihat konfirmasi ucapan hadiah.

---

## 6. Arsitektur Data & Skema Database Golang (GORM / SQL)

```sql
-- Tabel Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user', -- 'admin' atau 'user'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Subscriptions / Paket User
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    plan_type VARCHAR(20) NOT NULL, -- 'basic', 'premium', 'pro'
    status VARCHAR(20) NOT NULL, -- 'active', 'expired', 'pending'
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Transactions / Midtrans
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(100) UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id),
    plan_type VARCHAR(20) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    payment_type VARCHAR(50),
    transaction_status VARCHAR(50) NOT NULL, -- 'pending', 'settlement', 'expire', 'cancel'
    snap_token VARCHAR(255),
    raw_response JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Events / Undangan
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(100) UNIQUE NOT NULL, -- SEO Slug
    groom_name VARCHAR(150) NOT NULL,
    groom_parents VARCHAR(255) NOT NULL,
    bride_name VARCHAR(150) NOT NULL,
    bride_parents VARCHAR(255) NOT NULL,
    akad_date TIMESTAMP WITH TIME ZONE NOT NULL,
    akad_location TEXT NOT NULL,
    resepsi_date TIMESTAMP WITH TIME ZONE NOT NULL,
    resepsi_location TEXT NOT NULL,
    music_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## 7. Roadmap & Milestone Milestone

| Tahap | Durasi | Target Capaian |
| :--- | :--- | :--- |
| **Fase 1: Backend Golang & DB** | Minggu 1–2 | Arsitektur REST API Golang, Auth JWT (User & Admin), CRUD Events & Guest, Webhook Midtrans Sandbox. |
| **Fase 2: Next.js Frontend Structure** | Minggu 3–4 | Setup App Router Next.js, Layout `/admin` & `/user`, Dynamic Slug Routing `/[slug]`, Integrasi Midtrans Snap JS. |
| **Fase 3: Template & Multi-Tier Features** | Minggu 5–6 | Pemisahan fitur Basic/Premium/Pro, Editor Undangan Interaktif, WA Link Generator, Export RSVP Excel. |
| **Fase 4: Production Readiness** | Minggu 7–8 | Pengalihan Midtrans ke Production, Optimasi SEO Next.js (OpenGraph, Meta Tags per Slug), Load Testing Backend Golang. |

---