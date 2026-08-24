"use client"

import type {
  WeddingEvent,
  EventSections,
  GalleryPhoto,
  Music,
  DigitalGift,
  TemplateSummary,
  LoveStory,
} from "@/lib/types"
import { Calendar, MapPin, Clock, Music as MusicIcon } from "lucide-react"

const FONT_MAP: Record<string, string> = {
  serif: "Georgia, 'Times New Roman', serif",
  "sans-serif": "ui-sans-serif, system-ui, sans-serif",
  handwritten: "'Dancing Script', cursive",
}

function resolveFont(f?: string): string {
  if (!f) return "ui-sans-serif, system-ui, sans-serif"
  return FONT_MAP[f] || f
}

function formatDate(d?: string | null) {
  if (!d) return ""
  const date = new Date(d)
  if (Number.isNaN(date.getTime())) return d
  return date.toLocaleDateString("id-ID", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  })
}

function colorStyle(primary: string) {
  return { color: primary, borderColor: primary } as const
}

type Props = {
  event: WeddingEvent | null
  sections: EventSections | null
  gallery: GalleryPhoto[]
  music: Music | null
  gift: DigitalGift | null
  loveStories?: LoveStory[]
  template?: TemplateSummary | null
}

export default function InvitationPreview({
  event,
  sections,
  gallery,
  music,
  gift,
  loveStories = [],
  template,
}: Props) {
  if (!event) {
    return (
      <div className="rounded-xl border border-dashed border-input p-10 text-center text-sm text-muted-foreground">
        Pratinjau akan muncul setelah undangan dibuat.
      </div>
    )
  }

  const css = (template?.css_config ?? {}) as Record<string, unknown>
  const primary = (css.primary_color as string) || "#b22234"
  const background = (css.background_color as string) || "#faf7f5"
  const heroImage =
    (css.hero_image as string) ||
    "https://images.unsplash.com/photo-1519046904744-6fd9aeda9e0e?auto=format&fit=crop&w=1600&q=60"
  const fontFamily = resolveFont(css.font_family as string | undefined)

  const sec: EventSections = sections || {
    id: "",
    event_id: "",
    hero_enabled: true,
    couple_enabled: true,
    event_details_enabled: true,
    gallery_enabled: true,
    video_enabled: false,
    rsvp_enabled: true,
    guestbook_enabled: true,
    digital_gifts_enabled: false,
  }

  const couple =
    event.couple_name ||
    [event.groom_name, event.bride_name].filter(Boolean).join(" & ") ||
    "Mempelai"
  const names = event.groom_name && event.bride_name
    ? `${event.groom_name} & ${event.bride_name}`
    : couple

  return (
    <div className="space-y-3">
      <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
        Pratinjau Undangan
      </p>
      <div
        className="overflow-hidden rounded-xl border shadow-sm"
        style={{ fontFamily, backgroundColor: background }}
      >
        {sec.hero_enabled && (
          <section
            className="relative flex h-56 items-center justify-center bg-cover bg-center px-6 text-center"
            style={{
              backgroundImage: `url('${heroImage}')`,
              backgroundColor: primary,
            }}
          >
            <div className="absolute inset-0 bg-black/35" />
            <div className="relative z-10 max-w-2xl text-white">
              <h1 className="text-2xl font-bold sm:text-3xl">{names}</h1>
              <p className="mt-2 text-sm sm:text-base">
                Pernikahan {formatDate(event.wedding_date)}
              </p>
            </div>
          </section>
        )}

        <div className="space-y-8 px-6 py-8 text-foreground">
          {sec.opening_message && (
            <p className="text-center text-sm italic leading-relaxed">
              {sec.opening_message}
            </p>
          )}

          {sec.couple_enabled && (
            <section className="space-y-1 text-center">
              <h2 className="text-2xl font-bold" style={{ color: primary }}>
                {names}
              </h2>
              {event.couple_name && (
                <p className="text-sm text-muted-foreground">
                  {event.couple_name}
                </p>
              )}
              {event.groom_parents && (
                <p className="text-sm text-muted-foreground">
                  {event.groom_parents}
                </p>
              )}
              {event.bride_parents && (
                <p className="text-sm text-muted-foreground">
                  {event.bride_parents}
                </p>
              )}
            </section>
          )}

          {sec.verse_enabled && sec.verse_text && (
            <section className="space-y-2 text-center">
              <p className="text-xs font-semibold uppercase tracking-wider" style={{ color: primary }}>
                {sec.verse_religion === "alkitab" ? "Alkitab" : "Al-Quran"}
              </p>
              <p className="text-sm italic leading-relaxed">{sec.verse_text}</p>
              {sec.verse_source && (
                <p className="text-xs" style={{ color: primary }}>
                  {sec.verse_source}
                </p>
              )}
            </section>
          )}

          {sec.event_details_enabled && (
            <section className="space-y-6 text-sm">
              <div className="flex items-start gap-3">
                <Calendar className="mt-1 h-5 w-5" style={colorStyle(primary)} />
                <div>
                  <h3 className="font-semibold">Tanggal & Waktu</h3>
                  <p>{formatDate(event.wedding_date)}</p>
                  <p className="mt-1 flex items-center gap-1 text-xs">
                    <Clock className="h-4 w-4" /> {event.wedding_time || "—"}
                  </p>
                </div>
              </div>

              {event.ceremony_venue && (
                <div className="flex items-start gap-3">
                  <MapPin className="mt-1 h-5 w-5" style={colorStyle(primary)} />
                  <div>
                    <h3 className="font-semibold">Akad</h3>
                    <p>{event.ceremony_venue}</p>
                    {event.ceremony_address && (
                      <p className="text-xs text-muted-foreground">
                        {event.ceremony_address}
                      </p>
                    )}
                    {event.ceremony_map_url && (
                      <a
                        href={event.ceremony_map_url}
                        target="_blank"
                        rel="noreferrer"
                        className="text-xs text-blue-600"
                      >
                        Lihat peta
                      </a>
                    )}
                  </div>
                </div>
              )}

              {event.reception_venue && (
                <div className="flex items-start gap-3">
                  <MapPin className="mt-1 h-5 w-5" style={colorStyle(primary)} />
                  <div>
                    <h3 className="font-semibold">Resepsi</h3>
                    <p>{event.reception_venue}</p>
                    {event.reception_address && (
                      <p className="text-xs text-muted-foreground">
                        {event.reception_address}
                      </p>
                    )}
                    {event.reception_map_url && (
                      <a
                        href={event.reception_map_url}
                        target="_blank"
                        rel="noreferrer"
                        className="text-xs text-blue-600"
                      >
                        Lihat peta
                      </a>
                    )}
                  </div>
                </div>
              )}

              {sec.dress_code && (
                <p className="text-sm">
                  <span className="font-semibold">Dress code:</span>{" "}
                  {sec.dress_code}
                </p>
              )}
            </section>
          )}

          {music && music.file_url && (
            <section className="flex items-center gap-3 text-sm">
              <MusicIcon className="h-5 w-5" style={colorStyle(primary)} />
              <span>{music.title}</span>
              <audio controls src={music.file_url} className="h-8" />
            </section>
          )}

          {sec.gallery_enabled && gallery.length > 0 && (
            <section className="space-y-3">
              <h3 className="text-center font-semibold">Galeri</h3>
              <div className="grid grid-cols-2 gap-2">
                {gallery.map((p) => (
                  <img
                    key={p.id}
                    src={p.image_url}
                    alt={p.caption || ""}
                    className="aspect-square w-full rounded-lg object-cover shadow"
                  />
                ))}
              </div>
            </section>
          )}

          {sec.love_story_enabled && loveStories.length > 0 && (
            <section className="space-y-3">
              <h3 className="text-center font-semibold">Kisah Kami</h3>
              {loveStories.map((st) => (
                <div key={st.id} className="space-y-1 text-sm">
                  {st.image_url && (
                    <img src={st.image_url} alt={st.title} className="h-24 w-full rounded-lg object-cover" />
                  )}
                  <p className="font-medium" style={{ color: primary }}>{st.title}</p>
                  <p className="text-muted-foreground">{st.story}</p>
                </div>
              ))}
            </section>
          )}

          {sec.video_enabled && event.video_url && (
            <section className="space-y-2">
              <h3 className="text-center font-semibold">Video</h3>
              <a href={event.video_url} target="_blank" rel="noreferrer" className="block rounded-lg border p-3 text-center text-sm text-blue-600">
                Tonton video pernikahan
              </a>
            </section>
          )}

          {sec.digital_gifts_enabled && gift && (
            <section className="space-y-3 rounded-xl bg-white p-5 shadow-sm">
              <h3 className="font-semibold">Hadiah Digital</h3>
              {gift.gift_message && (
                <p className="text-sm italic">{gift.gift_message}</p>
              )}
              {(gift.bank_accounts || []).map((acc, i) => (
                <div key={i} className="rounded border p-3 text-sm">
                  <p className="font-medium">
                    {(acc.bank as string) || "Bank"} - {(acc.name as string)}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {acc.account as string}
                  </p>
                </div>
              ))}
              {gift.qris_image_url && (
                <img
                  src={gift.qris_image_url}
                  alt="QRIS"
                  className="h-28 w-28"
                />
              )}
            </section>
          )}

          {sec.closing_message && (
            <p className="text-center text-sm leading-relaxed">
              {sec.closing_message}
            </p>
          )}
        </div>
      </div>
    </div>
  )
}
