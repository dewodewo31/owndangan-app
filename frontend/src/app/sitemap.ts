import type { MetadataRoute } from "next"

interface PublicInvitation {
  slug: string
  updated_at: string
}

const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"
const API_URL = process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const staticRoutes = ["", "/cara-order", "/faq", "/testimoni"]
  const staticEntries: MetadataRoute.Sitemap = staticRoutes.map((path) => ({
    url: `${SITE_URL}${path}`,
    lastModified: new Date(),
  }))

  let invitationEntries: MetadataRoute.Sitemap = []
  try {
    const res = await fetch(`${API_URL}/invitations/public`, { next: { revalidate: 3600 } })
    if (res.ok) {
      const json = (await res.json()) as { data?: PublicInvitation[] }
      const invitations = json.data ?? []
      invitationEntries = invitations.map((inv) => ({
        url: `${SITE_URL}/${inv.slug}`,
        lastModified: inv.updated_at ? new Date(inv.updated_at) : new Date(),
      }))
    }
  } catch {
    // Backend unreachable — expose static routes only.
  }

  return [...staticEntries, ...invitationEntries]
}
