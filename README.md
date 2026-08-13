# Wedding Invitation Platform

SaaS platform for digital & print wedding invitations, Indonesia market.

## Overview

Decoupled architecture:
- **Backend:** Golang + GORM + PostgreSQL
- **Frontend:** Next.js App Router + TypeScript + Tailwind CSS
- **Payment:** Midtrans Snap
- **Auth:** JWT

Dual dashboard system:
- **Admin Dashboard** — platform management
- **User Dashboard** — wedding invitation management

Public invitation pages via SEO-friendly dynamic slugs.

## Architecture

```
Backend (Golang REST API)  <-->  Database (PostgreSQL)
      |
    API Gateway / Reverse Proxy
      |
Frontend (Next.js App Router)
  ├── (public)/  Landing, Pricing, [slug]
  ├── admin/     Admin dashboard
  ├── user/      User dashboard
  └── [slug]/    Public invitation
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go |
| API | RESTful |
| ORM | GORM |
| Database | PostgreSQL |
| Frontend | Next.js App Router |
| Language | TypeScript |
| Styling | Tailwind CSS |
| Auth | JWT |
| Payment | Midtrans Snap |

## Repository Structure

```
/
├── backend/          # Golang REST API
├── frontend/         # Next.js application
├── docs/             # Project documentation
├── prd.md            # Product Requirement Document
└── AGENTS.md         # AI coding agent instructions
```

## Getting Started

> TODO: Define setup commands after implementation begins.

## Development

See `docs/development/setup.md` and `docs/development/local-development.md`.

## Testing

See `docs/testing/strategy.md`.

## Documentation

Full documentation is in `/docs`. Start at `docs/README.md`.

## Roadmap

See `docs/roadmap/roadmap.md`.

## License

All rights reserved.
