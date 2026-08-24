# OWNDANGAN — Design System & UI/UX Specification

> Design specification for **owndangan-app**, a digital wedding invitation platform with a public invitation experience, invitation editor, guest management, RSVP, digital gift, QR check-in, and administration dashboard.
>
> This document adapts the supplied Material Design 3 specification to the OWNDANGAN product context. The supplied source establishes the visual language, editor structure, workflows, technical direction, responsive breakpoints, and performance targets. fileciteturn0file0L1-L21

---

## 1. Design Goals

OWNDANGAN should feel:

- **Elegant** — suitable for wedding invitations without looking overly corporate.
- **Modern** — clean Material 3 foundations with contemporary wedding aesthetics.
- **Personal** — every invitation should feel different through themes, typography, imagery, colors, and section arrangements.
- **Easy to edit** — non-technical users should be able to build an invitation through guided controls.
- **Fast** — public invitation pages are content-heavy and must remain lightweight.
- **Responsive** — the public invitation is mobile-first while the editor is optimized for desktop/tablet workflows.
- **Consistent** — dashboard, editor, public invitation, guest management, and check-in use the same underlying design tokens.

---

# 2. Design System

## 2.1 Design Foundation

Use **Material Design 3** as the structural foundation while allowing wedding-specific visual themes.

Material 3 should control:

- spacing
- component behavior
- form controls
- dialogs
- navigation
- elevation
- responsive layout
- accessibility states

Wedding themes should control:

- decorative typography
- palette
- background treatment
- ornaments
- image treatment
- section decoration
- animation style

The product must **not** force every wedding invitation into one visual style. The editor should allow multiple invitation themes while keeping dashboard and editor controls consistent.

---

## 2.2 Typography

OWNDANGAN uses two clearly separated typography systems:

1. **Product UI Typography** — landing page, dashboard, editor, forms, tables, navigation, and administrative interfaces.
2. **Invitation Canvas Typography** — the actual wedding invitation preview/public page, where typography is part of the visual identity of each template.

Product UI prioritizes readability and consistency. Invitation themes prioritize emotion, elegance, cultural character, and visual storytelling.

### 2.2.1 Product UI Font System

#### Primary: Plus Jakarta Sans

**Role:** Primary UI font for the landing page, dashboard, editor, headings, buttons, navigation, and general interface copy.

**Character:** modern geometric-humanist sans-serif, clean, contemporary, and friendly while remaining professional.

Plus Jakarta Sans is the default OWNDANGAN product font.

#### Secondary: Inter

**Role:** Dense UI content such as body text, form labels, data tables, metadata, statistics, and compact editor controls.

**Character:** neutral, highly readable at small sizes, and suitable for dense numerical/data interfaces.

Inter should be preferred when information density is more important than brand personality.

#### Display / Marketing: Outfit

**Role:** Landing-page display headings, promotional banners, feature headlines, and selected marketing callouts.

**Character:** geometric, contemporary, visually strong, and suitable for large display text.

#### Optional Display Alternative: Montserrat

Montserrat may be used for selected promotional or display sections when a stronger geometric character is required. Do not mix Outfit and Montserrat indiscriminately on the same page.

### 2.2.2 Product UI Type Scale

| Token | Font | Size | Line Height | Weight | Usage |
|---|---|---:|---:|---:|---|
| Display Large | Outfit / Plus Jakarta Sans | 57px | 64px | 600 | Marketing hero |
| Heading Large | Plus Jakarta Sans | 22px | 28px | 600 | Page headings |
| Title Medium | Plus Jakarta Sans | 16px | 24px | 500 | Card titles / controls |
| Body Medium | Inter / Plus Jakarta Sans | 14px | 20px | 400 | General content |
| Label Small | Inter | 11px | 16px | 500 | Labels / metadata |

The original Material 3 type-scale values are retained, while font families are assigned according to the researched font-pairing strategy.

### 2.2.3 Invitation Canvas Font System

Invitation templates must not be limited to one global font pairing. Each template defines:

```text
heading_font
body_font
accent_font
```

Approved invitation font families:

#### Playfair Display

**Role:** bride/groom hero names, wedding titles, event headings, and formal section titles.

**Character:** high-contrast serif, luxurious, formal, classical, and elegant.

Best suited to **Classic Elegance** and refined romantic themes.

#### Cinzel

**Role:** luxury headings, couple names, chapter/section headings, and formal display typography.

**Character:** Roman inscription-inspired, structured, ceremonial, and premium.

Best suited to **Modern Luxury** and sophisticated minimal themes.

#### Cormorant Garamond

**Role:** quotes, religious or ceremonial text, love story, opening messages, and editorial typography.

**Character:** refined, poetic, calm, literary, and traditional.

#### Alex Brush

**Role:** monograms, initials, decorative couple names, and small script accents.

Use sparingly. Never use it for long paragraphs or critical information.

#### Great Vibes

**Role:** romantic names, monograms, decorative headings, and signature-style accents.

Use as a visual accent rather than a primary reading font.

#### Allura

**Role:** subtle handwritten accents, initials, decorative names, and small invitation ornaments.

Allura remains secondary to a readable serif or sans-serif family.

#### Lora

**Role:** readable serif body text, classic invitation copy, secondary headings, and long-form story content.

#### Bodoni Moda

**Role:** fashion/editorial-inspired headings, luxury themes, and minimalist luxury typography.

Use selectively at display sizes.

### 2.2.4 Approved Font Pairing Matrix

| Style / Theme | Heading | Body | Accent |
|---|---|---|---|
| Landing Page & SaaS UI | Plus Jakarta Sans Bold | Plus Jakarta Sans / Inter | Outfit |
| Modern Luxury | Cinzel / Cormorant Garamond | Plus Jakarta Sans | Alex Brush |
| Classic Elegance | Playfair Display | Lora / Inter | Great Vibes |
| Minimalist & Clean | Montserrat | Inter | Bodoni Moda |
| Romantic Editorial | Cormorant Garamond | Plus Jakarta Sans / Lora | Allura |
| Modern Classic | Playfair Display | Plus Jakarta Sans | Great Vibes |
| Luxury Editorial | Bodoni Moda | Inter / Plus Jakarta Sans | Alex Brush |
| Cultural / Ceremonial | Cormorant Garamond / Cinzel | Lora | Allura |

The first four rows are the primary pairings from the provided font research. The additional rows are approved extensions for OWNDANGAN's template system.

### 2.2.5 Font Pairing Rules

1. Maximum three font families per invitation theme.
2. One family is the dominant heading/display family.
3. One family provides readable body copy.
4. Script fonts are accent fonts only.
5. Never use script fonts for paragraphs, RSVP forms, event details, addresses, or payment information.
6. UI fonts remain independent from the selected invitation theme.
7. Changing an invitation template may change typography substantially.
8. The editor must preview typography using the actual selected theme fonts.
9. Use real font weights/styles when available; avoid faux-bold and faux-italic.
10. Avoid excessive font-size variation inside one section.
11. Names and headings may use decorative fonts, but dates, addresses, times, RSVP instructions, and payment information must prioritize readability.
12. Templates must differ through typography, spacing, composition, imagery, and decoration—not only color.

### 2.2.6 Typography by Product Surface

| Surface | Recommended Typography |
|---|---|
| Landing page navigation | Plus Jakarta Sans |
| Landing page hero | Outfit + Plus Jakarta Sans |
| Landing page feature cards | Plus Jakarta Sans |
| Pricing / package tables | Plus Jakarta Sans + Inter |
| User dashboard | Plus Jakarta Sans |
| Admin dashboard | Plus Jakarta Sans + Inter |
| Data tables | Inter |
| Editor controls | Plus Jakarta Sans + Inter |
| Editor property labels | Inter |
| Public invitation hero | Theme heading font |
| Couple names | Theme heading / display font |
| Wedding story | Theme body font |
| Event information | Theme body font |
| RSVP form | Theme body font + product UI conventions |
| Digital gift/payment information | Theme body font with strong readability |
| Decorative monogram | Theme accent/script font |

### 2.2.7 Font Loading & Performance

Public invitation pages must avoid loading every available font. Each invitation should load only the font families and weights required by its selected theme.

Example:

```text
Modern Luxury
├── Cinzel: 400, 600
├── Plus Jakarta Sans: 400, 500
└── Alex Brush: 400
```

Do not globally load every approved font for every invitation. This protects the public invitation performance target and reduces unnecessary network requests.

## 2.3 Color System

Use Material 3 semantic color tokens rather than hard-coded colors.

### Core Tokens

```text
primary
on-primary
primary-container
on-primary-container

secondary
on-secondary
secondary-container
on-secondary-container

tertiary
on-tertiary
tertiary-container
on-tertiary-container

surface
surface-container
surface-container-low
surface-container-high
surface-variant

on-surface
on-surface-variant

outline
outline-variant

error
on-error
error-container
on-error-container
```

### Default OWNDANGAN Theme Direction

The default product interface may use an elegant **Deep Plum / Rose Gold** direction while supporting light and dark modes.

Important:

- Do not hard-code the wedding invitation itself to Deep Plum / Rose Gold.
- Each invitation theme owns its own palette.
- Dashboard/editor colors should remain product-consistent.
- Public invitation themes can be dramatically different from each other.

---

## 2.4 Shape Tokens

Use the supplied Material 3 shape hierarchy:

| Component | Radius |
|---|---:|
| Inputs / chips / small controls | 8px |
| Cards / dialogs / medium components | 16px |
| Large panels / sidebars / sheets | 28px |

These correspond to the supplied `small`, `medium`, and `extra-large` shape levels. fileciteturn0file0L7-L12

For public wedding templates, theme-specific decorative shapes may override visual treatment while preserving usable interaction targets.

---

## 2.5 Elevation

| Level | Usage |
|---|---|
| Level 0 | Main canvas / page background |
| Level 1 | Cards, sidebar panels, dashboard grids |
| Level 2 | Floating editor toolbar / floating controls |
| Level 3 | Dialogs, preview overlays, modal surfaces |

The supplied design specifies 0dp, 1dp, 3dp, and 6dp elevation layers respectively. fileciteturn0file0L13-L17

Avoid excessive shadows. Wedding pages should generally prefer layering, contrast, borders, and imagery over heavy elevation.

---

# 3. Product Architecture

OWNDANGAN consists of the following major experiences:

```text
OWNDANGAN
│
├── Public Website
│   ├── Landing Page
│   ├── Template Showcase
│   ├── Pricing / Package
│   └── Public Invitation
│
├── User Dashboard
│   ├── Overview
│   ├── My Invitations
│   ├── Create Invitation
│   ├── Invitation Editor
│   ├── Guest Management
│   ├── RSVP
│   ├── Digital Gift
│   ├── Analytics
│   └── Account Settings
│
├── Event Operations
│   ├── QR Check-in
│   ├── Guest Book
│   ├── Seat / Table Assignment
│   └── Welcome Screen
│
└── Admin Dashboard
    ├── Users
    ├── Invitations
    ├── Templates
    ├── Payments
    ├── Orders / Transactions
    ├── Media
    └── System Settings
```

The supplied architecture identifies the editor, guest/distribution management, check-in interface, analytics dashboard, media processing, WhatsApp gateway, payment/QRIS services, and domain routing as core platform areas. fileciteturn0file0L21-L28

---

# 4. Public Invitation Design

The public invitation is the most visually expressive part of OWNDANGAN.

## 4.1 Standard Section Model

An invitation can contain:

1. Cover / Hero
2. Opening / Quote
3. Groom & Bride
4. Story / Timeline
5. Event Details
6. Add to Calendar
7. Gallery
8. RSVP
9. Guestbook
10. Digital Gift
11. Closing
12. Footer

Sections must be:

- reorderable in the editor
- individually configurable
- individually hideable
- theme-aware
- responsive

The supplied specification explicitly defines hero, couple details, event timeline, gallery/story, audio, digital gift, and RSVP/guestbook modules. fileciteturn0file0L30-L44

---

## 4.2 Hero / Cover

Required capabilities:

- bride name
- groom name
- wedding date
- background image
- optional background video
- opening button
- optional guest name
- animated entrance

Example:

```text
────────────────────────────
        THE WEDDING

       SITI & BUDI

       20.11.2026

      [ BUKA UNDANGAN ]

        Kepada:
        {{guest_name}}
────────────────────────────
```

The hero should immediately establish the visual identity of the selected template.

---

## 4.3 Couple Section

Display:

- bride name
- groom name
- parent/family information
- profile image
- optional social links

Recommended visual variants:

- side-by-side
- stacked
- editorial
- floral
- minimalist
- traditional

---

# 5. Story / Timeline

The timeline must be visually different between templates.

Do **not** implement every theme as the same generic vertical list.

Possible layouts:

### Classic Vertical

```text
o Pertemuan Pertama
│
│  Cerita pertama...
│
o Lamaran
│
│  Cerita berikutnya...
│
o Pernikahan
```

### Alternating Timeline

```text
Pertemuan Pertama       o
                       │
                 o     Lamaran
                       │
Pernikahan              o
```

### Editorial

Large year/date markers with image and text blocks.

### Minimal

Thin line with compact milestones.

The content model remains consistent while the template controls presentation.

---

# 6. Event Details

Support multiple events:

- Akad
- Pemberkatan
- Resepsi
- After Party
- Custom event

Each event supports:

```text
title
date
start_time
end_time
venue_name
address
maps_url
description
```

The supplied source requires multi-event support, date/time configuration, venue/address details, and Google Maps integration. fileciteturn0file0L38-L41

UI should make multiple events easy to understand without overwhelming the guest.

---

# 7. Gallery

Gallery supports:

- grid
- masonry
- single featured image
- horizontal carousel
- lightbox
- lazy loading

Images should use responsive sizing and compressed WebP/modern image formats where possible.

Gallery layout is controlled by the selected template.

---

# 8. Audio

Audio manager supports:

- background music
- play/pause
- loop
- autoplay configuration
- floating player

Autoplay must never be assumed to work on every browser. The invitation should provide an obvious manual play interaction.

The supplied source includes background music selection, autoplay, loop, and floating audio controls. fileciteturn0file0L41-L43

---

# 9. Digital Gift

Digital gift supports:

- bank account
- account holder
- QRIS image
- physical gift address
- copy account number
- QRIS display

Example:

```text
       Hadiah Pernikahan

        [ QRIS IMAGE ]

     Bank BCA
     1234567890
     Siti Rahma

     [ Salin Nomor Rekening ]
```

The supplied source defines bank information, QRIS upload, and physical gift/shipping configuration. fileciteturn0file0L42-L44

---

# 10. RSVP & Guestbook

RSVP fields should be configurable.

Possible fields:

- attendance status
- number of guests
- meal/diet notes
- guest message

Guestbook should support:

- name
- message
- timestamp
- moderation status

The editor must allow wedding owners to enable or disable optional fields. fileciteturn0file0L43-L44

---

# 11. Guest Management

## 11.1 Guest Data

Minimum guest structure:

```text
id
name
category
phone
table_number
slug
rsvp_status
guest_count
qr_token
created_at
```

Import must support CSV/Excel and map columns such as:

- Name
- Category
- Phone Number
- Table Number

This follows the supplied guest-management specification. fileciteturn0file0L46-L53

---

## 11.2 Personalized Invitation URL

Default pattern:

```text
/v/{bride-groom}?to={guest_name}&cat={category}
```

Example:

```text
/v/siti-and-budi?to=Bapak+Ahmad&cat=VVIP
```

The system should support placeholders:

```text
{guest_name}
{event_date}
{venue_name}
{unique_link}
```

fileciteturn0file0L51-L56

---

# 12. WhatsApp Distribution

Guest messaging UI should provide:

- message template
- preview
- guest selection
- bulk sending
- sending queue
- delivery status
- configurable delay

Statuses:

```text
Pending
Sent
Delivered
Opened
Failed
```

The source specifies batch dispatch, configurable delays, custom message templates, and delivery tracking. fileciteturn0file0L54-L57

The UI must make it clear that the system is managing a queue rather than instantly sending every message.

---

# 13. QR Check-in

Every guest receives a unique QR token.

Check-in interface:

```text
┌─────────────────────────────┐
│       SCAN QR GUEST         │
│                             │
│       ┌─────────────┐       │
│       │             │       │
│       │   CAMERA    │       │
│       │             │       │
│       └─────────────┘       │
│                             │
│      [ Enter Manually ]     │
└─────────────────────────────┘
```

After successful scan:

```text
[check_circle] CHECK-IN BERHASIL

Bapak/Ibu Ahmad
Table A-12

19:42
```

The supplied specification requires encrypted guest QR tokens, a mobile web scanner, arrival timestamp logging, seat/table display, and optional label printing. fileciteturn0file0L59-L64

---

# 14. Welcome Screen

A separate full-screen route should be available for venue displays.

Example:

```text
┌─────────────────────────────────────┐
│                                     │
│           SELAMAT DATANG            │
│                                     │
│            BAPAK AHMAD              │
│                                     │
│          Meja / Table A-12          │
│                                     │
│        Siti & Budi Wedding          │
│                                     │
└─────────────────────────────────────┘
```

The supplied source describes a standalone full-screen output that reacts to successful guest scans. fileciteturn0file0L61-L64

---

# 15. Invitation Editor

## 15.1 Desktop Layout

Use the supplied three-column editor structure:

```text
┌──────────────┬──────────────────────────────┬──────────────────┐
│ LEFT PANEL   │       CENTER CANVAS          │ RIGHT PANEL      │
│ 280px        │       Flexible               │ 320px            │
├──────────────┼──────────────────────────────┼──────────────────┤
│ Sections     │                              │ Properties       │
│ Layers       │     Invitation Preview      │ Typography       │
│ Assets       │                              │ Colors           │
│              │                              │ Spacing          │
└──────────────┴──────────────────────────────┴──────────────────┘
```

This three-column model is directly based on the supplied editor workspace specification. fileciteturn0file0L68-L79

---

## 15.2 Left Panel

Contains:

- section tree
- drag handles
- visibility toggle
- add section button
- asset library

Example:

```text
SECTIONS

[menu] Hero                 [visibility]
[menu] Quote                [visibility]
[menu] Mempelai             [visibility]
[menu] Story                [visibility]
[menu] Event                [visibility]
[menu] Gallery              [visibility]
[menu] RSVP                 [visibility]
[menu] Gift                 [visibility]
[menu] Footer               [visibility]

+ Add Section
```

The supplied source explicitly requires section navigation, reordering, and visibility controls. fileciteturn0file0L81-L84

---

## 15.3 Center Canvas

The canvas must support:

- mobile portrait
- mobile landscape
- tablet
- desktop
- zoom
- pan
- undo
- redo
- preview

The supplied source specifies responsive viewport controls and a mobile frame preview. fileciteturn0file0L86-L89

### Important

The editor preview is **not** the final public page.

The public page must be rendered using the same configuration but optimized independently for production performance.

---

## 15.4 Right Property Inspector

Controls change according to the selected section.

Examples:

### Hero selected

```text
HERO
─────────────────
Title
[ Siti & Budi ]

Date
[ 20 November 2026 ]

Background
[ Image ]

Button Text
[ Buka Undangan ]

Animation
[ Fade [expand_more] ]
```

### Event selected

```text
EVENT
─────────────────
Event Name
[ Resepsi ]

Date
[ 20/11/2026 ]

Start
[ 10:00 ]

End
[ 14:00 ]

Venue
[ Gedung ABC ]

Address
[ ... ]

Maps
[ URL ]
```

The supplied source specifies context-sensitive property controls, typography controls, alignment controls, and spacing/radius sliders. fileciteturn0file0L91-L95

---

# 16. Editor Modes

The editor should expose two modes.

## Visual Mode

For users who want direct manipulation:

- drag sections
- reorder blocks
- edit content
- select elements
- adjust colors
- adjust typography
- adjust spacing
- upload media

## Structure Mode

For users who prefer forms:

```text
Hero
├── Bride Name
├── Groom Name
├── Wedding Date
├── Background
└── Button

Event
├── Event Name
├── Date
├── Time
├── Venue
└── Maps
```

The two-mode concept comes directly from the supplied WYSIWYG specification. fileciteturn0file0L30-L37

---

# 17. Dashboard Design

The dashboard should prioritize the user's invitations and important wedding activity.

## Overview

Suggested hierarchy:

```text
Good evening, Siti 

[ Total Invitations ] [ Guests ] [ RSVP ] [ Views ]

Your Invitations
┌───────────────────────┐
│ Siti & Budi            │
│ Published              │
│ 20 Nov 2026            │
│ [ Edit ] [ Preview ]   │
└───────────────────────┘

Recent Activity
─────────────────────────
[check_circle] Guest RSVP
[check_circle] Invitation viewed
[check_circle] Guest added
```

Avoid turning the dashboard into a generic enterprise admin panel. The primary goal is helping the wedding owner complete and manage an invitation.

---

# 18. Admin Dashboard

Admin UI may use a denser layout than the customer dashboard.

Main modules:

```text
Overview
Users
Invitations
Templates
Payments
Transactions
Media
Reports
Settings
```

Admin tables should support:

- search
- filter
- sorting
- pagination
- bulk action
- status badges

---

# 19. Responsive Design

Use the supplied responsive breakpoints:

| Breakpoint | Range | Target |
|---|---:|---|
| Compact | < 600px | Mobile portrait |
| Medium | 600–840px | Tablet / landscape |
| Expanded | > 840px | Desktop |

fileciteturn0file0L176-L181

## Public Invitation

Mobile-first.

## Editor

Desktop-first, but usable on tablet.

## Dashboard

Responsive:

- desktop: sidebar + content
- tablet: compact sidebar
- mobile: bottom navigation / drawer

---

# 20. Accessibility

Every interactive component must provide:

- visible focus state
- keyboard accessibility where applicable
- sufficient contrast
- semantic labels
- accessible form errors
- accessible buttons
- touch targets appropriate for mobile

Never communicate important information through color alone.

---

# 21. Motion & Animation

Animations should enhance the wedding atmosphere without hurting performance.

Allowed:

- fade
- slide
- gentle scale
- reveal
- staggered section entrance
- subtle parallax

Avoid:

- excessive continuous motion
- heavy particle effects
- animations that block interaction
- large JavaScript animation payloads

Respect `prefers-reduced-motion`.

---

# 22. Data Configuration Model

Invitation content should be configuration-driven.

Example:

```json
{
  "invitation_id": "inv_98234723",
  "slug": "siti-and-budi",
  "theme_id": "theme_rustic_m3_01",
  "metadata": {
    "title": "Pernikahan Siti & Budi",
    "og_image_url": "https://cdn.platform.id/assets/og_siti_budi.jpg"
  },
  "config": {
    "audio": {
      "url": "https://cdn.platform.id/audio/backsound_01.mp3",
      "auto_play": true
    },
    "sections": [
      {
        "id": "hero_01",
        "type": "hero",
        "visible": true,
        "content": {
          "groom_nickname": "Budi",
          "bride_nickname": "Siti",
          "event_date": "2026-11-20T09:00:00Z"
        }
      },
      {
        "id": "gift_01",
        "type": "digital_gift",
        "visible": true,
        "content": {
          "qris_image_url": "https://cdn.platform.id/qris/qris_siti.png",
          "bank_accounts": [
            {
              "bank_name": "BCA",
              "account_number": "1234567890",
              "account_holder": "Siti Rahma"
            }
          ]
        }
      }
    ]
  }
}
```

The configuration-driven section model follows the supplied data-model example. fileciteturn0file0L121-L174

---

# 23. Frontend Implementation Direction

For the frontend:

- Next.js / React
- SSR for public invitation pages
- Tailwind CSS
- Material 3 semantic CSS variables
- Zustand for editor state
- reusable component system
- server/client boundaries chosen intentionally

The supplied source specifically recommends Next.js with SSR, Tailwind mapped to Material 3 tokens, and Zustand for editor state, undo/redo, and property synchronization. fileciteturn0file0L121-L126

### Recommended frontend separation

```text
app/
├── (marketing)/
├── (dashboard)/
├── (editor)/
├── (admin)/
├── v/
│   └── [slug]/
└── check-in/
```

The exact route structure may follow the existing OWNDANGAN repository architecture where already established.

---

# 24. Component Architecture

Suggested component hierarchy:

```text
components/
├── ui/
│   ├── Button
│   ├── Input
│   ├── Select
│   ├── Dialog
│   ├── Sheet
│   ├── Tabs
│   ├── Badge
│   └── Card
│
├── dashboard/
│   ├── DashboardSidebar
│   ├── StatsCard
│   ├── InvitationCard
│   └── ActivityList
│
├── editor/
│   ├── EditorShell
│   ├── SectionPanel
│   ├── Canvas
│   ├── PropertyInspector
│   ├── DeviceToolbar
│   ├── AssetLibrary
│   └── PreviewButton
│
├── invitation/
│   ├── InvitationRenderer
│   ├── HeroSection
│   ├── CoupleSection
│   ├── StorySection
│   ├── EventSection
│   ├── GallerySection
│   ├── RSVPSection
│   ├── GiftSection
│   └── FooterSection
│
└── check-in/
    ├── QRScanner
    ├── GuestResult
    └── WelcomeDisplay
```

---

# 25. Theme Architecture

A theme should define presentation, not business logic.

Example:

```ts
type InvitationTheme = {
  id: string
  name: string

  typography: {
    headingFont: string
    bodyFont: string
    accentFont?: string
  }

  colors: {
    primary: string
    secondary: string
    background: string
    surface: string
    text: string
  }

  typography: {
    heading: string
    body: string
    accent?: string
  }

  shapes: {
    cardRadius: number
    buttonRadius: number
  }

  sections: {
    hero: string
    story: string
    event: string
    gallery: string
    gift: string
  }

  motion?: {
    entrance: string
    sectionReveal: string
  }
}
```

Business data must remain separate from theme presentation.

---

# 26. Template Philosophy

OWNDANGAN should support multiple visual identities.

Initial template categories can include:

1. Minimalist
2. Elegant
3. Floral
4. Modern
5. Luxury
6. Traditional
7. Islamic-inspired
8. Rustic
9. Editorial
10. Dark / cinematic

Each template should have its own visual composition.

For example:

```text
Template A — Minimalist
- large whitespace
- thin borders
- restrained typography

Template B — Floral
- botanical decoration
- organic shapes
- soft image framing

Template C — Luxury
- dark surface
- metallic-inspired accents
- serif typography

Template D — Traditional
- ornamental motifs
- structured sections
- heritage-inspired composition
```

The editor's data model stays stable while the template renderer changes the presentation.

---

# 27. Invitation Publishing Flow

```text
[Start]
   ↓
Select Template
   ↓
Enter Wedding Information
   ↓
Customize Invitation
   ├── Edit Sections
   ├── Upload Media
   ├── Configure Event
   ├── Configure Gift
   └── Configure Music
   ↓
Add Guests
   ↓
Configure Domain / Slug
   ↓
Preview
   ↓
Publish
   ↓
Generate Guest Links
   ↓
Distribute Invitations
```

This follows the supplied creation/publishing flow while adapting it to the OWNDANGAN product journey. fileciteturn0file0L99-L109

---

# 28. Check-in Flow

```text
Guest Arrives
      ↓
Show Personal QR
      ↓
Staff Scans QR
      ↓
Validate Guest
      ↓
┌───────────────┬──────────────────┐
│ Valid         │ Invalid          │
│ ↓             │ ↓                │
│ Check-in      │ Show Error       │
│ recorded      │                  │
└───────┬───────┴──────────────────┘
        ↓
Display Guest Information
        ↓
Optional Welcome Screen
        ↓
Optional Table / Label Output
```

This follows the supplied day-of-event check-in flow. fileciteturn0file0L110-L119

---

# 29. Performance Requirements

Public invitation pages must target:

| Metric | Target |
|---|---:|
| FCP | < 1.2s |
| LCP | < 2.0s |
| CLS | < 0.05 |
| Total initial asset weight | < 2.5MB |

These targets are specified in the supplied source. fileciteturn0file0L176-L186

### Implementation priorities

- optimize hero images
- use responsive image sizes
- lazy-load galleries
- compress images
- avoid loading editor-only code on public pages
- split heavy editor bundles
- defer non-critical scripts
- avoid unnecessary client-side rendering
- preload only critical assets
- keep public invitation HTML/SSR lightweight

---

# 30. UX Rules

## Rule 1 — Public invitation is the product

The invitation must look polished even if the user only fills the default fields.

## Rule 2 — Editing should be predictable

Every editor action should immediately reflect in preview.

## Rule 3 — Never hide important settings

Wedding date, venue, RSVP, and publishing status should be easy to find.

## Rule 4 — Templates must actually look different

Changing templates should produce a meaningful visual transformation, not only a color swap.

## Rule 5 — Mobile comes first for guests

Most invitation visitors are expected to open the invitation through a phone/chat link, so public invitation interaction must be designed around mobile.

## Rule 6 — Dashboard and invitation are different experiences

Dashboard:

```text
functional + efficient
```

Invitation:

```text
emotional + visual + memorable
```

---

# 31. Design Tokens Example

```css
:root {
  --radius-sm: 8px;
  --radius-md: 16px;
  --radius-lg: 28px;

  --elevation-0: none;
  --elevation-1: 0 1px 3px rgb(0 0 0 / 0.08);
  --elevation-2: 0 3px 8px rgb(0 0 0 / 0.12);
  --elevation-3: 0 6px 20px rgb(0 0 0 / 0.16);

  --font-ui: "Plus Jakarta Sans", sans-serif;

  --text-display: 57px;
  --text-heading: 22px;
  --text-title: 16px;
  --text-body: 14px;
  --text-label: 11px;
}
```

Actual colors should be implemented through semantic Material 3 tokens rather than hard-coded product colors.

---

# 32. Definition of Done — UI/UX

A feature is considered visually complete when:

- [ ] Desktop layout works
- [ ] Tablet layout works
- [ ] Mobile layout works
- [ ] Loading state exists
- [ ] Empty state exists
- [ ] Error state exists
- [ ] Success state exists
- [ ] Hover state exists where applicable
- [ ] Focus state exists
- [ ] Disabled state exists
- [ ] Dark mode is considered where applicable
- [ ] Components use shared design tokens
- [ ] No unnecessary hard-coded colors
- [ ] Typography follows the design system
- [ ] Radius follows the shape system
- [ ] Animations respect reduced-motion preferences

---

# 33. Final Product Direction

OWNDANGAN should be treated as **two connected products**:

### Wedding Creation Platform

```text
Dashboard
    ↓
Template
    ↓
Editor
    ↓
Guest Management
    ↓
Publishing
```

### Wedding Day Experience

```text
Guest
   ↓
Personal Invitation
   ↓
RSVP
   ↓
QR
   ↓
Check-in
   ↓
Welcome / Table
   ↓
Guestbook / Gift
```

The visual system keeps these experiences connected through shared design tokens, but their UX priorities remain different.

The final implementation should preserve the supplied Material 3 foundation while making the **public wedding invitation highly customizable and visually distinctive per template**.
