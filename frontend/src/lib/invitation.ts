import type { InvitationModel } from "@/templates/types"
import { DEFAULT_GUEST_NAME, getGuestName, getGuestToken } from "./guest"

/** Shape of GET /api/v1/e/{slug} (PublicEventResponse). */
export interface PublicEventResponse {
  event: {
    id: string
    title: string
    couple_name?: string
    groom_name?: string
    bride_name?: string
    groom_parents?: string
    bride_parents?: string
    wedding_date?: string
    wedding_time?: string
    ceremony_venue?: string
    ceremony_address?: string
    ceremony_map_url?: string
    reception_venue?: string
    reception_address?: string
    reception_map_url?: string
    video_url?: string
    view_count?: number
  }
  template?: {
    name?: string
    group_name?: string
    css_config?: Record<string, unknown>
    layout_config?: Record<string, unknown>
  } | null
  sections?: {
    hero_enabled?: boolean
    couple_enabled?: boolean
    event_details_enabled?: boolean
    gallery_enabled?: boolean
    video_enabled?: boolean
    rsvp_enabled?: boolean
    guestbook_enabled?: boolean
    love_story_enabled?: boolean
    digital_gifts_enabled?: boolean
    dress_code?: string
    opening_message?: string
    closing_message?: string
    verse_enabled?: boolean
    verse_religion?: string
    verse_text?: string
    verse_source?: string
    music?: { title: string; file_url?: string; preset?: string; is_preset?: boolean } | null
  } | null
  gallery?: { image_url: string; caption?: string; sort_order?: number }[]
  love_stories?: { id: string; title: string; story: string; year?: string; date?: string; image_url?: string; sort_order?: number }[]
  guestbook?: { name: string; message: string; created_at?: string }[]
  digital_gift?: {
    bank_accounts?: Array<Record<string, unknown>>
    ewallet?: Record<string, unknown>
    qris_image_url?: string
    gift_message?: string
  } | null
}

export interface InvitationPageData {
  model: InvitationModel
  title: string
  description: string
  heroImage: string
  primaryColor: string
}

/** Fetch the public event by slug. Returns null on any failure. */
export async function fetchPublicEvent(slug: string): Promise<PublicEventResponse | null> {
  try {
    const base = process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"
    const res = await fetch(`${base}/e/${slug}`, { next: { revalidate: 0 } })
    if (!res.ok) return null
    const json = await res.json()
    return (json.data ?? null) as PublicEventResponse
  } catch {
    return null
  }
}

/** Map the public API payload + URL params into the template model. */
export function buildInvitationModel(
  data: PublicEventResponse,
  searchParams: URLSearchParams,
  slug: string
): InvitationModel {
  const ev = data.event
  const couple =
    ev.couple_name || [ev.groom_name, ev.bride_name].filter(Boolean).join(" & ")
  const groomParents = ev.groom_parents
  const brideParents = ev.bride_parents
  const s = data.sections ?? {}

  const events: InvitationModel["events"] = {}
  if (ev.ceremony_venue || ev.ceremony_address) {
    events.akad = {
      label: "Akad Nikah",
      venue: ev.ceremony_venue,
      address: ev.ceremony_address,
      map_url: ev.ceremony_map_url,
      date: ev.wedding_date,
      time: ev.wedding_time,
    }
  }
  if (ev.reception_venue || ev.reception_address) {
    events.resepsi = {
      label: "Resepsi",
      venue: ev.reception_venue,
      address: ev.reception_address,
      map_url: ev.reception_map_url,
      date: ev.wedding_date,
      time: ev.wedding_time,
    }
  }

  return {
    slug,
    eventId: ev.id,
    names: {
      groom: ev.groom_name,
      bride: ev.bride_name,
      full: couple,
    },
    parents: { groom: groomParents, bride: brideParents },
    date: ev.wedding_date,
    time: ev.wedding_time,
    events,
    opening: s.opening_message,
    closing: s.closing_message,
    verse: {
      enabled: !!s.verse_enabled,
      religion: s.verse_religion,
      text: s.verse_text,
      source: s.verse_source,
    },
    dressCode: s.dress_code,
    gallery: (data.gallery ?? []).map((g) => ({
      image_url: g.image_url,
      caption: g.caption,
    })),
    loveStories: (data.love_stories ?? []).map((st) => ({
      title: st.title,
      story: st.story,
      year: st.year,
      date: st.date,
      image_url: st.image_url,
    })),
    video: ev.video_url || null,
    guestbook: (data.guestbook ?? []).map((g) => ({
      name: g.name,
      message: g.message,
      created_at: g.created_at,
    })),
    gift: data.digital_gift ?? undefined,
    music: s.music ?? null,
    sections: {
      hero_enabled: s.hero_enabled ?? true,
      couple_enabled: s.couple_enabled ?? true,
      event_details_enabled: s.event_details_enabled ?? true,
      gallery_enabled: s.gallery_enabled ?? true,
      video_enabled: s.video_enabled ?? false,
      love_story_enabled: s.love_story_enabled ?? false,
      rsvp_enabled: s.rsvp_enabled ?? true,
      guestbook_enabled: s.guestbook_enabled ?? true,
      digital_gifts_enabled: s.digital_gifts_enabled ?? false,
    },
    token: getGuestToken(searchParams),
    guestName: getGuestName(searchParams),
  }
}

/** Assemble the render-ready page data (model + SEO bits). */
export function buildPageData(
  data: PublicEventResponse,
  searchParams: URLSearchParams,
  slug: string
): InvitationPageData {
  const model = buildInvitationModel(data, searchParams, slug)
  const css = (data.template?.css_config ?? {}) as Record<string, unknown>
  const primary = (css.primary_color as string) || "#b22234"
  const heroImage =
    (css.hero_image as string) ||
    model.gallery[0]?.image_url ||
    "https://images.unsplash.com/photo-1519046904744-6fd9aeda9e0e?auto=format&fit=crop&w=1600&q=60"
  const names = model.names.full || "Undangan Pernikahan"
  const date = model.date
    ? new Date(model.date).toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
      })
    : ""

  const guest = model.guestName && model.guestName !== DEFAULT_GUEST_NAME ? ` untuk ${model.guestName}` : ""
  return {
    model,
    title: `The Wedding of ${names}`,
    description: `Kami mengundang Anda${guest} untuk hadir di acara pernikahan ${names}${date ? ` pada ${date}` : ""}.`,
    heroImage,
    primaryColor: primary,
  }
}
