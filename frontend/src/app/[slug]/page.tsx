import { notFound } from "next/navigation"
import type { Metadata } from "next"
import { TemplateShell } from "@/components/invitation/shell"
import { CoverGate } from "@/components/invitation/cover-gate"
import { fetchPublicEvent, buildPageData, type PublicEventResponse, type InvitationPageData } from "@/lib/invitation"
import { selectTemplate } from "@/templates"

export const revalidate = 0
export const dynamic = "force-dynamic"

interface PageProps {
  params: { slug: string }
  searchParams: { to?: string; token?: string }
}

export async function generateMetadata({ params, searchParams }: PageProps): Promise<Metadata> {
  const data = await fetchPublicEvent(params.slug)
  if (!data) {
    return { title: "Undangan Tidak Ditemukan" }
  }
  const page = buildPageData(data, new URLSearchParams(searchParams), params.slug)
  const url = new URL(`/${params.slug}`, process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000")
  if (searchParams.to) url.searchParams.set("to", searchParams.to)

  return {
    title: page.title,
    description: page.description,
    openGraph: {
      title: page.title,
      description: page.description,
      type: "website",
      url: url.toString(),
      images: [{ url: page.heroImage, width: 1600, height: 1067, alt: page.title }],
    },
    twitter: {
      card: "summary_large_image",
      title: page.title,
      description: page.description,
      images: [page.heroImage],
    },
    alternates: { canonical: url.toString() },
  }
}

function buildJsonLd(
  data: PublicEventResponse,
  page: InvitationPageData,
  slug: string
): Record<string, unknown> {
  const ev = data.event
  const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"
  const canonical = `${siteUrl}/${slug}`

  const jsonLd: Record<string, unknown> = {
    "@context": "https://schema.org",
    "@type": "Event",
    name: page.title,
    description: page.description,
    url: canonical,
    eventAttendanceMode: "https://schema.org/OfflineEventAttendanceMode",
    eventStatus: "https://schema.org/EventScheduled",
  }

  if (ev.wedding_date) {
    jsonLd.startDate = ev.wedding_date
  }
  if (page.heroImage) {
    jsonLd.image = [page.heroImage]
  }

  const location = page.model.events.resepsi || page.model.events.akad
  if (location) {
    const place: Record<string, unknown> = { "@type": "Place" }
    if (location.venue) place.name = location.venue
    if (location.address) place.address = location.address
    if (Object.keys(place).length > 1) jsonLd.location = place
  }

  if (page.model.names.full) {
    jsonLd.organizer = {
      "@type": "Person",
      name: page.model.names.full,
    }
  }

  return jsonLd
}

export default async function InvitationPage({ params, searchParams }: PageProps) {
  const data = await fetchPublicEvent(params.slug)
  if (!data) {
    notFound()
  }

  const page = buildPageData(data, new URLSearchParams(searchParams), params.slug)
  const definition = selectTemplate(
    data.template?.group_name,
    data.template?.name
  )
  const jsonLd = buildJsonLd(data, page, params.slug)

  return (
    <div className="relative">
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />
      <TemplateShell definition={definition} model={page.model} />
      <CoverGate model={page.model} heroImage={page.heroImage} primaryColor={page.primaryColor} />
    </div>
  )
}
