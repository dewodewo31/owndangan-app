# SEO

## Strategy

SEO is critical for this platform because each invitation page must rank for the couple's names and wedding-related queries. The strategy is built on Next.js Metadata API, server-side rendering, and dynamic metadata generation per slug. Every public page is server-rendered; no client-side rendering for SEO-sensitive content.

## Metadata API

All pages use the `Metadata` export from Next.js to set `<title>`, `<meta>`, and social tags.

### Root Layout Metadata

```tsx
// app/layout.tsx
export const metadata: Metadata = {
  title: {
    template: '%s | Undangan Pernikahan Digital',
    default: 'Undangan Pernikahan Digital - Buat Undangan Online',
  },
  description: 'Buat undangan pernikahan digital dengan mudah. Pilih template, atur tamu, dan bagikan dengan tamu undangan.',
  openGraph: {
    type: 'website',
    locale: 'id_ID',
    siteName: 'Undangan Pernikahan Digital',
  },
};
```

The `title.template` pattern means every child page that sets a `title` will have it prepended to the template. Example: setting `title: 'Dashboard'` on `/user` produces `<title>Dashboard | Undangan Pernikahan Digital</title>`.

## Dynamic [slug] Metadata

The `[slug]` page generates metadata dynamically based on the invitation data. This is the most important SEO optimization.

```tsx
// app/[slug]/page.tsx
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const event = await getInvitation(params.slug);
  const coupleName = `${event.groom_name} & ${event.bride_name}`;
  const weddingDate = new Date(event.wedding_date).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'long', year: 'numeric',
  });
  const title = `Undangan Pernikahan ${coupleName}`;
  const description = `Kami mengundang Bapak/Ibu/Saudara/i pada acara pernikahan ${coupleName} yang akan dilaksanakan pada ${weddingDate}.`;

  return {
    title,
    description,
    openGraph: {
      title,
      description,
      type: 'website',
      url: `https://domain.com/${params.slug}`,
      images: [{
        url: event.gallery?.[0]?.image_url || '/images/og-default.jpg',
        width: 1200,
        height: 630,
        alt: `Undangan Pernikahan ${coupleName}`,
      }],
    },
    twitter: {
      card: 'summary_large_image',
      title,
      description,
      images: [event.gallery?.[0]?.image_url || '/images/og-default.jpg'],
    },
    robots: {
      index: true,
      follow: true,
    },
    alternates: {
      canonical: `https://domain.com/${params.slug}`,
    },
  };
}
```

## OpenGraph & Twitter Cards

The metadata includes:

- **og:title** — "Undangan Pernikahan [Couple Names]"
- **og:description** — "Kami mengundang [Bapak/Ibu/Saudara/i] pada acara pernikahan [Couple Names]..."
- **og:image** — First gallery photo or a default OG image (1200x630, with a subtle wedding-themed background and the couple's names overlaid).
- **og:url** — Canonical URL of the invitation.
- **og:type** — `website`.
- **twitter:card** — `summary_large_image` for rich previews.

The default OG image (`/images/og-default.jpg`) is a branded template with the platform logo and a placeholder couple silhouette. User-uploaded gallery images are preferred when available.

## Sitemap Generation

A `sitemap.ts` file at the app root generates a dynamic sitemap of all published invitations.

```tsx
// app/sitemap.ts
export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const baseUrl = 'https://domain.com';
  const events = await getPublishedEvents(); // GET /api/v1/public/events/sitemap

  const eventUrls = events.map((event) => ({
    url: `${baseUrl}/${event.slug}`,
    lastModified: event.updated_at,
    changeFrequency: 'weekly' as const,
    priority: 0.9,
  }));

  return [
    { url: baseUrl, lastModified: new Date(), changeFrequency: 'daily', priority: 1.0 },
    { url: `${baseUrl}/pricing`, lastModified: new Date(), changeFrequency: 'monthly', priority: 0.7 },
    ...eventUrls,
  ];
}
```

The sitemap is submitted to Google Search Console and Bing Webmaster Tools. The backend provides a dedicated endpoint (`GET /api/v1/public/events/sitemap`) that returns all published slugs with their last-modified timestamps.

## robots.txt

A `robots.ts` file generates the robots.txt:

```tsx
// app/robots.ts
export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      { userAgent: '*', allow: '/' },
      { userAgent: '*', disallow: ['/admin/', '/user/', '/api/'] },
    ],
    sitemap: 'https://domain.com/sitemap.xml',
  };
}
```

Admin and user dashboard routes are disallowed from indexing. Public invitation pages and landing pages are fully indexable.

## Structured Data (JSON-LD)

The `[slug]` page includes structured data for the event:

```json
{
  "@context": "https://schema.org",
  "@type": "Event",
  "name": "Pernikahan Andi & Sinta",
  "description": "Undangan pernikahan Andi Pratama dan Sinta Dewi",
  "startDate": "2025-06-15T09:00:00+07:00",
  "endDate": "2025-06-15T17:00:00+07:00",
  "location": {
    "@type": "Place",
    "name": "Gedung Serbaguna Jakarta",
    "address": "Jl. Merdeka No. 1, Jakarta"
  }
}
```

This is injected via a `<script type="application/ld+json">` tag in the page component.

## Performance Considerations

- Server components eliminate client-side JavaScript for SEO-critical pages.
- Images use `next/image` with WebP format and responsive breakpoints.
- Fonts are preloaded via `next/font` to avoid layout shift.
- The `/[slug]` page uses `fetch()` with `next: { revalidate: 60 }` for ISR — the page is regenerated at most every 60 seconds, balancing freshness with cache hit rate.