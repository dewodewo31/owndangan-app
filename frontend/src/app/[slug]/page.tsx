import { notFound } from "next/navigation"
import { Calendar, MapPin, Clock, Music } from "lucide-react"
import type { PublicEvent } from "@/lib/types"

export const revalidate = 0

interface InvitationData {
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
  view_count?: number
}

interface MusicDTO {
  title: string
  file_url?: string
  preset?: string
  is_preset: boolean
}

interface Sections {
  hero_enabled: boolean
  couple_enabled: boolean
  event_details_enabled: boolean
  gallery_enabled: boolean
  video_enabled: boolean
  rsvp_enabled: boolean
  guestbook_enabled: boolean
  digital_gifts_enabled: boolean
  dress_code?: string
  opening_message?: string
  closing_message?: string
  verse_enabled?: boolean
  verse_religion?: string
  verse_text?: string
  verse_source?: string
  music?: MusicDTO | null
}

interface GalleryPhoto {
  image_url: string
  caption?: string
  sort_order: number
}

interface GuestbookMsg {
  name: string
  message: string
  created_at: string
}

interface DigitalGift {
  bank_accounts?: Array<Record<string, unknown>>
  ewallet?: Record<string, unknown>
  qris_image_url?: string
  gift_message?: string
}

interface Template {
  name?: string
  group_name?: string
  css_config?: Record<string, unknown>
  layout_config?: Record<string, unknown>
}

interface Invitation {
  event: InvitationData
  template?: Template
  sections?: Sections
  gallery?: GalleryPhoto[]
  guestbook?: GuestbookMsg[]
  digital_gift?: DigitalGift
}

async function fetchInvitation(slug: string): Promise<Invitation | null> {
  try {
    const base = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"
    const res = await fetch(`${base}/e/${slug}`, { next: { revalidate: 0 } })
    if (!res.ok) return null
    const json = await res.json()
    return (json.data ?? null) as Invitation
  } catch {
    return null
  }
}

const FONT_MAP: Record<string, string> = {
  serif: "serif",
  "sans-serif": "ui-sans, system-ui, sans-serif",
  handwritten: "'Dancing Script', cursive",
}

function resolveFont(f?: string): string {
  if (!f) return "ui-sans, system-ui, sans-serif"
  return FONT_MAP[f] || f
}

function formatDate(d?: string) {
  if (!d) return ""
  const date = new Date(d)
  if (Number.isNaN(date.getTime())) return d
  return date.toLocaleDateString("id-ID", { weekday: "long", year: "numeric", month: "long", day: "numeric" })
}

function colorStyle(primary: string) {
  return { color: primary, borderColor: primary } as const
}

export default async function InvitationPage({ params }: { params: { slug: string } }) {
  const data = await fetchInvitation(params.slug)
  if (!data) {
    notFound()
  }

  const ev = data.event
  const sections: Sections = data.sections || {
    hero_enabled: true, couple_enabled: true, event_details_enabled: true,
    gallery_enabled: true, video_enabled: false, rsvp_enabled: true,
    guestbook_enabled: true, digital_gifts_enabled: false,
  }
  const couple = ev.couple_name || `${ev.groom_name || ""} & ${ev.bride_name || ""}`.replace(/^ & $/, "")
  const names = ev.groom_name && ev.bride_name ? `${ev.groom_name} & ${ev.bride_name}` : couple

  const css = (data.template?.css_config as any) || {}
  const primary = css.primary_color || "#b22234"
  const background = css.background_color || "#faf7f5"
  const heroImage =
    css.hero_image ||
    "https://images.unsplash.com/photo-1519046904744-6fd9aeda9e0e?auto=format&fit=crop&w=1600&q=60"
  const fontFamily = resolveFont(css.font_family as string | undefined)
  const muted = "text-muted-foreground"

  return (
    <div className="min-h-screen" style={{ fontFamily, backgroundColor: background }}>
      {sections.hero_enabled && (
        <section
          className="relative h-[60vh] min-h-[420px] flex items-center justify-center text-center px-6 bg-cover bg-center"
          style={{
            backgroundImage: `url('${heroImage}')`,
            backgroundColor: primary,
          }}
        >
          <div className="absolute inset-0 bg-black/35" />
          <div className="relative z-10 max-w-3xl text-white">
            <h1 className="text-3xl sm:text-5xl font-bold mb-3">{names}</h1>
            <p className="text-lg sm:text-xl">Pernikahan {formatDate(ev.wedding_date)}</p>
          </div>
        </section>
      )}

      <main className="max-w-3xl mx-auto px-6 py-12 space-y-12 text-foreground">
        {sections.opening_message && (
          <section className="text-center">
            <p className="text-lg italic leading-relaxed">{sections.opening_message}</p>
          </section>
        )}

        {sections.couple_enabled && (
          <section className="text-center space-y-3">
            <h2 className="text-3xl font-bold" style={{ color: primary }}>{names}</h2>
            {ev.couple_name && <p className={muted}>{ev.couple_name}</p>}
            {ev.groom_parents && <p className={muted}>{ev.groom_parents}</p>}
            {ev.bride_parents && <p className={muted}>{ev.bride_parents}</p>}
          </section>
        )}

        {sections.verse_enabled && sections.verse_text && (
          <section className="text-center space-y-3">
            <p className="text-sm font-semibold uppercase tracking-wider" style={{ color: primary }}>
              {sections.verse_religion === "alkitab" ? "Alkitab" : "Al-Quran"}
            </p>
            <p className="text-lg italic leading-relaxed max-w-2xl mx-auto">{sections.verse_text}</p>
            {sections.verse_source && (
              <p className="text-sm" style={{ color: primary }}>{sections.verse_source}</p>
            )}
          </section>
        )}

        {sections.event_details_enabled && (
          <section className="space-y-8">
            <div className="flex items-start gap-4">
              <Calendar className="h-5 w-5 mt-1" style={colorStyle(primary)} />
              <div>
                <h3 className="font-semibold">Tanggal & Waktu</h3>
                <p>{formatDate(ev.wedding_date)}</p>
                <p className="flex items-center gap-1 mt-1 text-sm"><Clock className="h-4 w-4" /> {ev.wedding_time || "—"}</p>
              </div>
            </div>

            {ev.ceremony_venue && (
              <div className="flex items-start gap-4">
                <MapPin className="h-5 w-5 mt-1" style={colorStyle(primary)} />
                <div>
                  <h3 className="font-semibold">Akad</h3>
                  <p>{ev.ceremony_venue}</p>
                  {ev.ceremony_address && <p className={muted}>{ev.ceremony_address}</p>}
                  {ev.ceremony_map_url && (
                    <a href={ev.ceremony_map_url} target="_blank" rel="noreferrer" className="text-sm text-blue-600">Lihat peta</a>
                  )}
                </div>
              </div>
            )}

            {ev.reception_venue && (
              <div className="flex items-start gap-4">
                <MapPin className="h-5 w-5 mt-1" style={colorStyle(primary)} />
                <div>
                  <h3 className="font-semibold">Resepsi</h3>
                  <p>{ev.reception_venue}</p>
                  {ev.reception_address && <p className={muted}>{ev.reception_address}</p>}
                  {ev.reception_map_url && (
                    <a href={ev.reception_map_url} target="_blank" rel="noreferrer" className="text-sm text-blue-600">Lihat peta</a>
                  )}
                </div>
              </div>
            )}

            {sections.dress_code && (
              <p className="text-sm"><span className="font-semibold">Dress code:</span> {sections.dress_code}</p>
            )}
          </section>
        )}

        {sections.music && sections.music.file_url && (
          <section className="flex items-center gap-3">
            <Music className="h-5 w-5" style={colorStyle(primary)} />
            <span>{sections.music.title}</span>
            <audio controls src={sections.music.file_url} />
          </section>
        )}

        {sections.gallery_enabled && data.gallery && data.gallery.length > 0 && (
          <section className="space-y-4">
            <h3 className="font-semibold text-center">Galeri</h3>
            <div className="columns-2 sm:columns-3 gap-3 space-y-3">
              {data.gallery.map((p) => (
                <div key={p.image_url} className="break-inside-avoid">
                  <img src={p.image_url} alt={p.caption || ""} className="w-full rounded-lg shadow" />
                  {p.caption && <p className="text-xs text-center mt-1 text-muted-foreground">{p.caption}</p>}
                </div>
              ))}
            </div>
          </section>
        )}

        {sections.digital_gifts_enabled && data.digital_gift && (
          <section className="space-y-4 bg-white p-6 rounded-xl shadow">
            <h3 className="font-semibold">Hadiah Pengungkapan</h3>
            {data.digital_gift.gift_message && <p className="italic">{data.digital_gift.gift_message}</p>}
            {data.digital_gift.bank_accounts &&
              data.digital_gift.bank_accounts.map((acc, i) => (
                <div key={i} className="p-3 border rounded">
                  <p className="font-medium">{(acc.bank as string) || "Bank"} - {(acc.name as string)}</p>
                  <p className="text-sm text-muted-foreground">{acc.account as string}</p>
                </div>
              ))}
            {data.digital_gift.ewallet && Object.keys(data.digital_gift.ewallet).length > 0 && (
              <div className="flex gap-3 flex-wrap">
                {Object.entries(data.digital_gift.ewallet).map(([k, v]) => (
                  <span key={k} className="text-sm">
                    {(k as string).toUpperCase()}: {String(v)}
                  </span>
                ))}
              </div>
            )}
            {data.digital_gift.qris_image_url && <img src={data.digital_gift.qris_image_url} alt="QRIS" className="h-32" />}
          </section>
        )}

        {sections.guestbook_enabled && data.guestbook && data.guestbook.length > 0 && (
          <section className="space-y-4">
            <h3 className="font-semibold text-center">Buku Tamu</h3>
            <ul className="space-y-3">
              {data.guestbook.map((g, i) => (
                <li key={i} className="bg-white p-4 rounded-lg shadow">
                  <p className="font-medium">{g.name}</p>
                  <p className="text-sm text-muted-foreground">{g.message}</p>
                </li>
              ))}
            </ul>
          </section>
        )}

        {sections.closing_message && (
          <section className="text-center">
            <p className="text-lg leading-relaxed">{sections.closing_message}</p>
          </section>
        )}
      </main>
    </div>
  )
}
