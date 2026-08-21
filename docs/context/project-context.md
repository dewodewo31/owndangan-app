# Project Context

## Overview

Platform Undangan Pernikahan Digital & Cetak is a SaaS platform enabling couples (pengantin) to create, customize, and share digital wedding invitations. The platform also supports printed invitation orders (future). Built for the Indonesia market with Indonesian payment methods (QRIS, Bank Transfer, e-Wallet) via Midtrans.

## Target Users

1. **Couples (Pengantin)** — Users who create and manage their wedding invitations.
2. **Guests (Tamu)** — Invitees who RSVP, send messages, and view invitation details.
3. **Admin (Platform Owner)** — Manages users, transactions, packages, templates.

## Problem

Traditional wedding invitations are costly, environmentally wasteful, and difficult to manage for guest responses. Digital invitations solve these problems while adding interactivity (music, gallery, RSVP tracking, digital gifts).

## Core Value

- Easy-to-use invitation editor
- Beautiful templates
- Real-time RSVP tracking
- Integrated digital gifts (amplop digital)
- WhatsApp sharing
- SEO-friendly public invitation pages

## Business Model

Freemium SaaS with 3 paid tiers:
- **Starter** (Rp 99K, 30 days)
- **Premium** (Rp 299K, 60 days)
- **All Access** (Rp 999K, lifetime, unlimited guests)

## Platform Flow

```
User → Register/Login → Choose Package → Create Payment (Midtrans) → 
Payment Settlement (Webhook) → Subscription Activated → Create Invitation → 
Configure 9 Sections → Publish → Share → Guest RSVP
```

## Roles

| Role | Scope |
|------|-------|
| `admin` | Full platform management |
| `user` | Own invitations and guest management |
| `guest` | View public invitation, RSVP, send messages |

## Public vs Private Area

| Area | Routes | Auth |
|------|--------|------|
| Public | Landing, Pricing, `/[slug]` | None |
| User Dashboard | `/user/*` | JWT (role: user) |
| Admin Dashboard | `/admin/*` | JWT (role: admin) |

## Core Entities

- User
- Subscription
- Transaction (Midtrans)
- Event (Wedding Invitation)
- Guest
- RSVP
- Guestbook Message
- Template
- Digital Gift Info

## Core Modules

Authentication, Users, Subscriptions, Packages, Payments, Events (invitations), Invitation Editor, Templates, Guests, WhatsApp, RSVP, Guestbook, Digital Gifts, Gallery, Music, SEO, Admin, Analytics, Audit Log, Export.

## Status

> CONFIRMED — PRD v2.0 approved, documentation in progress, no code implemented yet.
