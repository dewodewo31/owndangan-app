import { notFound } from "next/navigation"
import type { Metadata } from "next"
import { TemplateShell } from "@/components/invitation/shell"
import { CoverGate } from "@/components/invitation/cover-gate"
import { fetchPublicEvent, buildPageData } from "@/lib/invitation"
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
    alternates: { canonical: url.toString() },
  }
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

  return (
    <div className="relative">
      <TemplateShell definition={definition} model={page.model} />
      <CoverGate model={page.model} heroImage={page.heroImage} primaryColor={page.primaryColor} />
    </div>
  )
}
