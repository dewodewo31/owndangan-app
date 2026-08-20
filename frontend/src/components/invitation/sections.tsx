"use client"

import { useState, useEffect, type FormEvent } from "react"
import { z } from "zod"
import {
  CalendarHeart,
  CalendarPlus,
  Check,
  Church,
  Copy,
  Download,
  Heart,
  MapPin,
  Navigation,
  PartyPopper,
  type LucideIcon,
} from "lucide-react"
import type { InvitationModel, EventBlock, SectionSpec, ThemeTokens, TimelineStyle } from "@/templates/types"
import { cn } from "@/lib/cn"
import { Lightbox } from "./lightbox"
import { Timeline } from "./timeline"

export interface SectionProps {
  model: InvitationModel
  theme: ThemeTokens
  spec?: SectionSpec
}

function apiBase() {
  return process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"
}

function fmtDate(d?: string) {
  if (!d) return ""
  const dt = new Date(d)
  if (Number.isNaN(dt.getTime())) return d
  return dt.toLocaleDateString("id-ID", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  })
}

function photo(model: InvitationModel, i = 0) {
  return model.gallery[i]?.image_url
}

function InstagramIcon({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <rect width="20" height="20" x="2" y="2" rx="5" ry="5" />
      <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z" />
      <line x1="17.5" x2="17.51" y1="6.5" y2="6.5" />
    </svg>
  )
}

const narrow = { maxWidth: "var(--t-content-width)", margin: "0 auto", width: "100%" }

/* ---------------------------------------------------------------- COVER */

export function Cover({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as string) || "centered"
  const img = photo(model, 0)
  const names = model.names.full

  if (variant === "minimal") {
    return (
      <header className="flex min-h-[70vh] flex-col items-center justify-center px-6 text-center" style={{ background: theme.background }}>
        <p className="mb-4 text-sm uppercase tracking-[0.4em]" style={{ color: theme.muted }}>
          The Wedding Of
        </p>
        <h1 className="text-5xl leading-tight sm:text-7xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>
          {names}
        </h1>
        <div className="mt-8 h-px w-24" style={{ background: theme.accent }} />
        <p className="mt-6 text-lg" style={{ color: theme.muted }}>{fmtDate(model.date)}</p>
      </header>
    )
  }

  if (variant === "split") {
    return (
      <header className="grid min-h-[88vh] grid-cols-1 md:grid-cols-2">
        <div className="relative min-h-[40vh] md:min-h-full" style={{ backgroundImage: img ? `url(${img})` : undefined, backgroundColor: theme.surface, backgroundSize: "cover", backgroundPosition: "center" }} />
        <div className="flex flex-col justify-center px-8 py-16 md:px-14" style={{ background: theme.background }}>
          <p className="mb-3 text-xs uppercase tracking-[0.35em]" style={{ color: theme.muted }}>Undangan Pernikahan</p>
          <h1 className="text-5xl leading-tight sm:text-6xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{names}</h1>
          <div className="my-6 h-px w-20" style={{ background: theme.accent }} />
          <p className="text-lg" style={{ color: theme.muted }}>{fmtDate(model.date)}</p>
          <p className="mt-2" style={{ color: theme.muted }}>{model.time}</p>
        </div>
      </header>
    )
  }

  if (variant === "cinematic" || variant === "framed" || variant === "editorial") {
    const overlay = variant === "cinematic" ? "rgba(0,0,0,0.55)" : variant === "editorial" ? "rgba(0,0,0,0.35)" : "rgba(0,0,0,0.3)"
    const align = variant === "editorial" ? "items-start text-left pl-8 md:pl-24" : "items-center text-center"
    return (
      <header className={cn("relative flex min-h-[var(--t-hero-height)] flex-col justify-end px-6 pb-20 pt-32 text-white", align)} style={{ backgroundImage: img ? `url(${img})` : undefined, backgroundColor: theme.background, backgroundSize: "cover", backgroundPosition: "center" }}>
        <div className="absolute inset-0" style={{ background: overlay }} />
        {variant === "framed" && (
          <div className="pointer-events-none absolute inset-5 border" style={{ borderColor: "rgba(255,255,255,0.5)" }} />
        )}
        <div className="relative z-10" style={{ maxWidth: "var(--t-content-width)" }}>
          {variant === "editorial" && <p className="mb-2 text-xs uppercase tracking-[0.4em] text-white/70">Wedding</p>}
          <h1 className={cn("text-6xl leading-none sm:text-8xl", variant === "editorial" ? "font-black" : "")} style={{ fontFamily: "var(--t-font-heading)" }}>{names}</h1>
          <p className="mt-5 text-lg text-white/85">{fmtDate(model.date)} · {model.time}</p>
        </div>
      </header>
    )
  }

  // default: centered
  return (
    <header className="relative flex min-h-[var(--t-hero-height)] flex-col items-center justify-center px-6 text-center text-white" style={{ backgroundImage: img ? `url(${img})` : undefined, backgroundColor: theme.background, backgroundSize: "cover", backgroundPosition: "center" }}>
      <div className="absolute inset-0 bg-black/40" />
      <div className="relative z-10" style={{ maxWidth: "var(--t-content-width)" }}>
        <p className="mb-3 text-sm uppercase tracking-[0.4em] text-white/80">The Wedding Of</p>
        <h1 className="text-6xl leading-tight sm:text-7xl" style={{ fontFamily: "var(--t-font-heading)" }}>{names}</h1>
        <p className="mt-6 text-lg text-white/85">{fmtDate(model.date)}</p>
      </div>
    </header>
  )
}

/* ---------------------------------------------------------------- QUOTE */

export function Quote({ model, theme, spec }: SectionProps) {
  if (model.verse?.enabled && model.verse.text) {
    return (
      <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto">
          <p className="text-xs font-semibold uppercase tracking-[0.3em]" style={{ color: theme.primary }}>
            {model.verse.religion === "alkitab" ? "Alkitab" : "Al-Quran"}
          </p>
          <p className="mt-4 text-xl italic leading-relaxed" style={{ fontFamily: "var(--t-font-accent)" }}>{model.verse.text}</p>
          {model.verse.source && <p className="mt-3 text-sm" style={{ color: theme.primary }}>{model.verse.source}</p>}
        </div>
      </section>
    )
  }
  if (model.opening) {
    return (
      <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto">
          <p className="text-lg italic leading-relaxed" style={{ color: theme.muted }}>{model.opening}</p>
        </div>
      </section>
    )
  }
  return null
}

/* ---------------------------------------------------------------- COUPLE */

export function Couple({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as string) || "portrait"
  const groom = model.couple?.groom
  const bride = model.couple?.bride

  if (groom || bride) {
    const p1 = groom?.photo || photo(model, 0)
    const p2 = bride?.photo || photo(model, 1) || p1
    const Person = ({ p, role, img }: { p: NonNullable<InvitationModel["couple"]>["groom"]; role: string; img?: string }) => (
      <div className="text-center">
        <div className="mx-auto aspect-[3/4] w-full max-w-xs overflow-hidden rounded-[var(--t-radius)]" style={{ background: theme.surface }}>
          {img && <img src={img} alt={p?.name || role} className="h-full w-full object-cover" loading="lazy" />}
        </div>
        <p className="mt-5 text-xs uppercase tracking-widest" style={{ color: theme.muted }}>{role}</p>
        {p?.name && (
          <h4 className="mt-1 text-2xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{p.name}</h4>
        )}
        {p?.nickname && <p className="text-sm" style={{ color: theme.accent }}>{p.nickname}</p>}
        {p?.description && <p className="mx-auto mt-3 max-w-sm text-sm leading-relaxed" style={{ color: theme.muted }}>{p.description}</p>}
        {(p?.childOrder || p?.parents) && (
          <p className="mt-3 text-sm" style={{ color: theme.muted }}>
            {p?.childOrder}{p?.childOrder && p?.parents ? " dari pasangan " : ""}{p?.parents}
          </p>
        )}
        {p?.instagram && (
          <a
            href={`https://instagram.com/${p.instagram.replace(/^@/, "")}`}
            target="_blank"
            rel="noreferrer"
            className="mt-4 inline-flex items-center gap-1.5 text-sm"
            style={{ color: theme.primary }}
          >
            <InstagramIcon className="h-4 w-4" /> @{p.instagram.replace(/^@/, "")}
          </a>
        )}
      </div>
    )
    return (
      <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <div style={narrow} className="mx-auto">
          <p className="text-center text-xs uppercase tracking-[0.3em]" style={{ color: theme.muted }}>Mempelai</p>
          <h3 className="mt-2 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Kedua Mempelai</h3>
          <div className="mt-10 grid grid-cols-1 gap-10 md:grid-cols-2">
            {groom && <Person p={groom} role="Mempelai Pria" img={p1} />}
            {bride && <Person p={bride} role="Mempelai Wanita" img={p2} />}
          </div>
        </div>
      </section>
    )
  }

  const p1 = photo(model, 0)
  const p2 = photo(model, 1) || p1

  if (variant === "stacked") {
    return (
      <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.background }}>
        <div style={narrow} className="mx-auto space-y-8">
          {p1 && <img src={p1} alt={model.names.full} className="mx-auto aspect-[4/5] w-full max-w-md rounded-[var(--t-radius)] object-cover shadow-lg" loading="lazy" />}
          <div>
            <h2 className="text-4xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{model.names.full}</h2>
            <p className="mt-3" style={{ color: theme.muted }}>{model.names.groom} &amp; {model.names.bride}</p>
          </div>
        </div>
      </section>
    )
  }

  if (variant === "editorial") {
    return (
      <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <div style={narrow} className="mx-auto grid grid-cols-1 items-center gap-8 md:grid-cols-[1fr_auto_1fr]">
          {p1 && <img src={p1} alt={model.names.groom || "Mempelai"} className="aspect-[3/4] w-full rounded-[var(--t-radius)] object-cover" loading="lazy" />}
          <div className="text-center">
            <h2 className="text-4xl leading-tight" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{model.names.groom}</h2>
            <span className="my-3 block text-2xl" style={{ color: theme.accent }}>&amp;</span>
            <h2 className="text-4xl leading-tight" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{model.names.bride}</h2>
          </div>
          {p2 && <img src={p2} alt={model.names.bride || "Mempelai"} className="aspect-[3/4] w-full rounded-[var(--t-radius)] object-cover" loading="lazy" />}
        </div>
      </section>
    )
  }

  // portrait (default)
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.background }}>
      <div style={narrow} className="mx-auto">
        <h2 className="mb-8 text-4xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{model.names.full}</h2>
        <div className="grid grid-cols-2 gap-4">
          {[p1, p2].map((p, i) => (
            <div key={i} className="overflow-hidden rounded-[var(--t-radius)]">
              {p ? (
                <img src={p} alt={(i === 0 ? model.names.groom : model.names.bride) || "Mempelai"} className="aspect-square w-full object-cover" loading="lazy" />
              ) : (
                <div className="aspect-square w-full" style={{ background: theme.surface }} />
              )}
            </div>
          ))}
        </div>
        <p className="mt-6" style={{ color: theme.muted }}>{model.names.groom} &amp; {model.names.bride}</p>
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- PARENTS */

export function Parents({ model, theme }: SectionProps) {
  if (!model.parents.groom && !model.parents.bride) return null
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto grid grid-cols-1 gap-6 sm:grid-cols-2">
        {model.parents.groom && (
          <div>
            <p className="text-xs uppercase tracking-widest" style={{ color: theme.muted }}>Mempelai Pria</p>
            <p className="mt-2" style={{ color: theme.text }}>{model.parents.groom}</p>
          </div>
        )}
        {model.parents.bride && (
          <div>
            <p className="text-xs uppercase tracking-widest" style={{ color: theme.muted }}>Mempelai Wanita</p>
            <p className="mt-2" style={{ color: theme.text }}>{model.parents.bride}</p>
          </div>
        )}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- COUNTDOWN */

export function Countdown({ model, theme }: SectionProps) {
  // null-first: server and client first render agree, ticking starts post-mount
  const [now, setNow] = useState<number | null>(null)
  useEffect(() => {
    if (!model.date) return
    setNow(Date.now())
    const id = setInterval(() => setNow(Date.now()), 1000)
    return () => clearInterval(id)
  }, [model.date])
  if (!model.date) return null
  const target = new Date(model.date + "T" + (model.time || "09:00:00")).getTime()

  const diff = now === null ? 0 : Math.max(0, target - now)
  const done = now !== null && target <= now
  const days = Math.floor(diff / 86400000)
  const hours = Math.floor((diff % 86400000) / 3600000)
  const mins = Math.floor((diff % 3600000) / 60000)
  const secs = Math.floor((diff % 60000) / 1000)
  const pad = (v: number) => (now === null ? "--" : String(v).padStart(2, "0"))
  const items = [
    { v: pad(days), l: "Hari" },
    { v: pad(hours), l: "Jam" },
    { v: pad(mins), l: "Menit" },
    { v: pad(secs), l: "Detik" },
  ]
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.background }}>
      <div style={narrow} className="mx-auto">
        <p className="mb-5 text-sm uppercase tracking-[0.3em]" style={{ color: theme.muted }}>
          {done ? "Hari Bahagia Telah Tiba" : "Menuju Hari Bahagia"}
        </p>
        {done ? (
          <p className="text-2xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>
            The Day Has Come
          </p>
        ) : (
          <div className="flex justify-center gap-4">
            {items.map((it) => (
              <div key={it.l} className="min-w-[64px] rounded-[var(--t-radius)] px-3 py-4" style={{ background: theme.surface, border: `1px solid var(--t-border)` }}>
                <div className="text-3xl font-semibold tabular-nums" style={{ color: theme.primary }}>{it.v}</div>
                <div className="text-xs" style={{ color: theme.muted }}>{it.l}</div>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- ADD TO CALENDAR */

function eventBlocks(model: InvitationModel): NonNullable<InvitationModel["events"]["akad"]>[] {
  return [model.events.akad, model.events.resepsi, ...(model.events.extra ?? [])].filter(
    (b): b is NonNullable<InvitationModel["events"]["akad"]> => !!b
  )
}

function endTime(b: EventBlock): string {
  if (b.end_time) return b.end_time
  if (!b.time) return "11:00:00"
  const [h, m] = b.time.split(":").map(Number)
  const d = new Date(0, 0, 0, (h || 0) + 2, m || 0)
  return `${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}:00`
}

function fmtTime(t?: string) {
  if (!t) return ""
  const [h, m] = t.split(":").map(Number)
  return `${String(h).padStart(2, "0")}.${String(m ?? 0).padStart(2, "0")}`
}

function gcalStamp(date?: string, time?: string): string {
  if (!date) return ""
  const [h = "09", m = "00", s = "00"] = (time || "09:00:00").split(":")
  const dt = new Date(`${date}T${h.padStart(2, "0")}:${m.padStart(2, "0")}:${s.padStart(2, "0")}`)
  if (Number.isNaN(dt.getTime())) return ""
  return dt.toISOString().replace(/[-:]/g, "").replace(/\.\d{3}/, "")
}

function googleCalendarUrl(b: EventBlock, coupleName: string): string {
  const params = new URLSearchParams({
    action: "TEMPLATE",
    text: `${b.label} - ${coupleName}`,
    dates: `${gcalStamp(b.date, b.time)}/${gcalStamp(b.date, endTime(b))}`,
    details: b.description || "",
    location: [b.venue, b.address].filter(Boolean).join(", "),
  })
  const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
  if (tz) params.set("ctz", tz)
  return `https://calendar.google.com/calendar/render?${params.toString()}`
}

function escapeIcs(v: string): string {
  return v.replace(/\\/g, "\\\\").replace(/\n/g, "\\n").replace(/,/g, "\\,").replace(/;/g, "\\;")
}

export interface IcsEventInput {
  date?: string
  time?: string
  end?: string
  label: string
  coupleName: string
  venue?: string
  address?: string
  description?: string
  timeZone?: string
}

export function buildIcs(opts: IcsEventInput): string {
  const pad6 = (t: string) => t.replace(/:/g, "").padEnd(6, "0")
  const datePart = (opts.date || "").replace(/-/g, "")
  const tzid = opts.timeZone ? `;TZID=${opts.timeZone}` : ""
  const start = datePart ? `DTSTART${tzid}:${datePart}T${pad6(opts.time || "09:00:00")}` : ""
  const end = datePart ? `DTEND${tzid}:${datePart}T${pad6(opts.end || opts.time || "09:00:00")}` : ""
  const stamp = new Date().toISOString().replace(/[-:]/g, "").replace(/\.\d{3}/, "")
  const loc = [opts.venue, opts.address].filter(Boolean).join(", ")
  const lines = [
    "BEGIN:VCALENDAR",
    "VERSION:2.0",
    "PRODID:-//Owndangan//Undangan Pernikahan//ID",
    opts.timeZone ? `X-WR-TIMEZONE:${opts.timeZone}` : "",
    "BEGIN:VEVENT",
    `UID:${opts.label}-${opts.date || "undangan"}@owndangan`,
    `DTSTAMP:${stamp}`,
    start,
    end,
    `SUMMARY:${escapeIcs(`${opts.label} - ${opts.coupleName}`)}`,
    loc ? `LOCATION:${escapeIcs(loc)}` : "",
    opts.description ? `DESCRIPTION:${escapeIcs(opts.description)}` : "",
    "BEGIN:VALARM",
    "TRIGGER:-PT1H",
    "ACTION:DISPLAY",
    `DESCRIPTION:Reminder ${escapeIcs(opts.label)}`,
    "END:VALARM",
    "END:VEVENT",
    "END:VCALENDAR",
  ]
    .filter(Boolean)
    .join("\r\n")
  return lines + "\r\n"
}

function icsFile(b: EventBlock, coupleName: string): string {
  return buildIcs({
    date: b.date,
    time: b.time,
    end: endTime(b),
    label: b.label,
    coupleName,
    venue: b.venue,
    address: b.address,
    description: b.description,
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  })
}

function downloadIcs(b: EventBlock, coupleName: string) {
  const blob = new Blob([icsFile(b, coupleName)], { type: "text/calendar;charset=utf-8" })
  const url = URL.createObjectURL(blob)
  const a = document.createElement("a")
  a.href = url
  a.download = `undangan-${b.label.toLowerCase().replace(/\s+/g, "-")}.ics`
  a.click()
  URL.revokeObjectURL(url)
}

function eventIcon(b: EventBlock): LucideIcon {
  const l = (b.label || "").toLowerCase()
  if (l.includes("akad") || l.includes("nikah")) return Heart
  if (l.includes("resepsi") || l.includes("pesta") || l.includes("party")) return PartyPopper
  if (l.includes("berkat") || l.includes("gereja")) return Church
  return CalendarPlus
}

export function AddToCalendar({ model, theme }: SectionProps) {
  const blocks = eventBlocks(model)
  if (blocks.length === 0) return null
  const couple = model.names.full || "Undangan Pernikahan"
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.background }}>
      <div style={narrow} className="mx-auto">
        <h3 className="text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>
          Simpan Tanggalnya
        </h3>
        <p className="mt-2 text-sm" style={{ color: theme.muted }}>
          Tambahkan acara ke kalender agar tidak terlewat.
        </p>
        <div className="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2">
          {blocks.map((b, i) => {
            const Icon = eventIcon(b)
            return (
              <div key={i} className="rounded-[var(--t-radius)] border p-6 text-left" style={{ borderColor: "var(--t-border)", background: theme.surface }}>
                <div className="flex items-start gap-3">
                  <span className="flex h-10 w-10 flex-none items-center justify-center rounded-full" style={{ background: `${theme.primary}14`, color: theme.primary }}>
                    <Icon className="h-5 w-5" />
                  </span>
                  <div>
                    <p className="font-semibold" style={{ color: theme.text }}>{b.label}</p>
                    <p className="mt-1 text-sm" style={{ color: theme.muted }}>
                      {b.date ? fmtDate(b.date) : ""} · {fmtTime(b.time)}{b.end_time ? ` – ${fmtTime(b.end_time)}` : ""}
                    </p>
                    {(b.venue || b.address) && (
                      <p className="mt-1 text-sm" style={{ color: theme.muted }}>{[b.venue, b.address].filter(Boolean).join(", ")}</p>
                    )}
                  </div>
                </div>
                <div className="mt-4 flex flex-wrap gap-2">
                  <a
                    href={googleCalendarUrl(b, couple)}
                    target="_blank"
                    rel="noreferrer"
                    className="inline-flex flex-1 items-center justify-center gap-1.5 rounded-[var(--t-radius)] px-3 py-2 text-xs font-semibold text-white"
                    style={{ background: theme.primary }}
                  >
                    <CalendarHeart className="h-4 w-4" /> Google Calendar
                  </a>
                  <button
                    type="button"
                    onClick={() => downloadIcs(b, couple)}
                    className="inline-flex flex-1 items-center justify-center gap-1.5 rounded-[var(--t-radius)] border px-3 py-2 text-xs font-semibold"
                    style={{ borderColor: "var(--t-border)", color: theme.primary }}
                  >
                    <Download className="h-4 w-4" /> Unduh .ics
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- EVENTS */

export function Events({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as string) || "cards"
  const blocks = eventBlocks(model)
  if (blocks.length === 0) return null

  if (variant === "timeline") {
    return (
      <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto space-y-6">
          <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Rangkaian Acara</h3>
          {blocks.map((b, i) => {
            const Icon = eventIcon(b)
            return (
              <div key={i} className="flex gap-4 border-l-2 pl-4" style={{ borderColor: theme.accent }}>
                <span className="flex h-9 w-9 flex-none items-center justify-center rounded-full" style={{ background: `${theme.primary}14`, color: theme.primary }}>
                  <Icon className="h-4 w-4" />
                </span>
                <div>
                  <p className="font-semibold" style={{ color: theme.primary }}>{b.label}</p>
                  {b.venue && <p style={{ color: theme.text }}>{b.venue}</p>}
                  {(b.date || b.time) && <p className="text-sm" style={{ color: theme.muted }}>{b.date || fmtDate(model.date)} · {fmtTime(b.time || model.time)}{b.end_time ? ` – ${fmtTime(b.end_time)}` : ""}</p>}
                  {b.address && <p className="text-sm" style={{ color: theme.muted }}>{b.address}</p>}
                  {b.description && <p className="mt-1 text-sm" style={{ color: theme.muted }}>{b.description}</p>}
                  {b.map_url && <a href={b.map_url} target="_blank" rel="noreferrer" className="text-sm underline" style={{ color: theme.primary }}>Lihat peta</a>}
                </div>
              </div>
            )
          })}
        </div>
      </section>
    )
  }

  if (variant === "side") {
    return (
      <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto grid grid-cols-1 gap-6 sm:grid-cols-2">
          {blocks.map((b, i) => {
            const Icon = eventIcon(b)
            return (
              <div key={i} className="rounded-[var(--t-radius)] p-6 text-center" style={{ background: theme.background, border: `1px solid var(--t-border)` }}>
                <span className="mx-auto flex h-10 w-10 items-center justify-center rounded-full" style={{ background: `${theme.primary}14`, color: theme.primary }}>
                  <Icon className="h-5 w-5" />
                </span>
                <p className="mt-3 text-lg font-semibold" style={{ color: theme.primary }}>{b.label}</p>
                {b.venue && <p className="mt-2" style={{ color: theme.text }}>{b.venue}</p>}
                <p className="text-sm" style={{ color: theme.muted }}>{b.date || fmtDate(model.date)} · {fmtTime(b.time || model.time)}{b.end_time ? ` – ${fmtTime(b.end_time)}` : ""}</p>
                {b.address && <p className="text-sm" style={{ color: theme.muted }}>{b.address}</p>}
                {b.description && <p className="mx-auto mt-2 max-w-xs text-sm" style={{ color: theme.muted }}>{b.description}</p>}
              </div>
            )
          })}
        </div>
      </section>
    )
  }

  // cards (default)
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto grid grid-cols-1 gap-6 sm:grid-cols-2">
        {blocks.map((b, i) => {
          const Icon = eventIcon(b)
          return (
            <div key={i} className="rounded-[var(--t-radius)] p-6 text-center shadow-sm" style={{ background: theme.background, border: `1px solid var(--t-border)` }}>
              <span className="mx-auto flex h-11 w-11 items-center justify-center rounded-full" style={{ background: `${theme.primary}14`, color: theme.primary }}>
                <Icon className="h-5 w-5" />
              </span>
              <p className="mt-3 text-xs uppercase tracking-widest" style={{ color: theme.muted }}>{b.label}</p>
              <h4 className="mt-1 text-2xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{b.venue || "—"}</h4>
              <p className="mt-2 text-sm" style={{ color: theme.muted }}>{b.date || fmtDate(model.date)} · {fmtTime(b.time || model.time)}{b.end_time ? ` – ${fmtTime(b.end_time)}` : ""}</p>
              {b.address && <p className="mt-1 text-sm" style={{ color: theme.muted }}>{b.address}</p>}
              {b.description && <p className="mx-auto mt-3 max-w-xs text-sm" style={{ color: theme.muted }}>{b.description}</p>}
              {b.map_url && <a href={b.map_url} target="_blank" rel="noreferrer" className="mt-2 inline-block text-sm underline" style={{ color: theme.primary }}>Lihat peta</a>}
            </div>
          )
        })}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- LOVE STORY */

export function LoveStory({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as TimelineStyle) || "classic"
  if (model.loveStories.length === 0) return null
  const items = model.loveStories.map((s) => ({
    year: s.year,
    date: s.date,
    title: s.title,
    description: s.story,
    image: s.image_url,
  }))
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
      <h3 className="mb-10 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>
        Kisah Kami
      </h3>
      <Timeline variant={variant} items={items} />
    </section>
  )
}

/* ---------------------------------------------------------------- VIDEO */

function embedVideoUrl(url: string): string {
  try {
    const u = new URL(url)
    if (u.hostname === "youtu.be") return `https://www.youtube.com/embed/${u.pathname.slice(1)}`
    if (u.hostname.includes("youtube.com")) {
      const v = u.searchParams.get("v")
      if (v) return `https://www.youtube.com/embed/${v}`
    }
    return url
  } catch {
    return url
  }
}

export function Video({ model, theme }: SectionProps) {
  if (!model.video || !model.sections?.video_enabled) return null
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
      <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Video</h3>
      <div className="mx-auto aspect-video w-full max-w-3xl overflow-hidden rounded-[var(--t-radius)]" style={{ border: `1px solid ${theme.border}` }}>
        <iframe
          src={embedVideoUrl(model.video)}
          title="Video Pernikahan"
          className="h-full w-full"
          loading="lazy"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowFullScreen
        />
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- GALLERY */

export function Gallery({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as string) || "grid"
  const [lightboxIndex, setLightboxIndex] = useState<number | null>(null)
  if (model.gallery.length === 0) return null

  const imgCls = "w-full cursor-pointer rounded-[var(--t-radius)] object-cover transition-opacity hover:opacity-90"
  const openAt = (i: number) => () => setLightboxIndex(i)
  const thumb = (g: { image_url: string; caption?: string }, i: number, cls: string) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img key={g.image_url + i} src={g.image_url} alt={g.caption || "Galeri"} className={cls} loading="lazy" onClick={openAt(i)} />
  )

  if (variant === "columns") {
    return (
      <section className="px-4 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Galeri</h3>
        <div className="columns-2 gap-3 sm:columns-3 lg:columns-4" style={{ maxWidth: "var(--t-content-width)", margin: "0 auto" }}>
          {model.gallery.map((g, i) => thumb(g, i, `${imgCls} mb-3 break-inside-avoid`))}
        </div>
        <Lightbox images={model.gallery} index={lightboxIndex} onClose={() => setLightboxIndex(null)} onIndexChange={setLightboxIndex} />
      </section>
    )
  }

  if (variant === "masonry") {
    return (
      <section className="px-4 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <div className="columns-2 gap-4" style={{ maxWidth: "var(--t-content-width)", margin: "0 auto" }}>
          {model.gallery.map((g, i) => (
            <figure key={g.image_url + i} className="mb-4 break-inside-avoid">
              {thumb(g, i, "w-full rounded-[var(--t-radius)] object-cover")}
              {g.caption && <figcaption className="mt-1 text-center text-xs" style={{ color: theme.muted }}>{g.caption}</figcaption>}
            </figure>
          ))}
        </div>
        <Lightbox images={model.gallery} index={lightboxIndex} onClose={() => setLightboxIndex(null)} onIndexChange={setLightboxIndex} />
      </section>
    )
  }

  if (variant === "horizontal") {
    return (
      <section className="py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Galeri</h3>
        <div className="flex snap-x gap-4 overflow-x-auto px-6 pb-2">
          {model.gallery.map((g, i) => thumb(g, i, "h-56 w-44 flex-none snap-center rounded-[var(--t-radius)] object-cover"))}
        </div>
        <Lightbox images={model.gallery} index={lightboxIndex} onClose={() => setLightboxIndex(null)} onIndexChange={setLightboxIndex} />
      </section>
    )
  }

  // grid (default)
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
      <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Galeri</h3>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3" style={{ maxWidth: "var(--t-content-width)", margin: "0 auto" }}>
        {model.gallery.map((g, i) => thumb(g, i, `${imgCls} aspect-square`))}
      </div>
      <Lightbox images={model.gallery} index={lightboxIndex} onClose={() => setLightboxIndex(null)} onIndexChange={setLightboxIndex} />
    </section>
  )
}

/* ---------------------------------------------------------------- LOCATION */

function mapQuery(b: EventBlock): string {
  return encodeURIComponent([b.address, b.venue].filter(Boolean).join(", ") || "Indonesia")
}

function embedMapUrl(b: EventBlock): string {
  return `https://www.google.com/maps?q=${mapQuery(b)}&output=embed`
}

function mapsLink(b: EventBlock): string {
  return `https://www.google.com/maps?q=${mapQuery(b)}`
}

function directionsLink(b: EventBlock): string {
  return `https://www.google.com/maps/dir/?api=1&destination=${mapQuery(b)}`
}

export function Location({ model, theme }: SectionProps) {
  const blocks = eventBlocks(model)
  if (blocks.length === 0) return null
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto space-y-4">
        <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Lokasi</h3>
        <div className="overflow-hidden rounded-[var(--t-radius)] border" style={{ borderColor: "var(--t-border)" }}>
          <iframe
            src={embedMapUrl(blocks[0])}
            title="Peta lokasi acara"
            className="h-64 w-full"
            loading="lazy"
            referrerPolicy="no-referrer-when-downgrade"
            allowFullScreen
          />
        </div>
        {blocks.map((b, i) => (
          <div key={i} className="rounded-[var(--t-radius)] p-5 text-center" style={{ background: theme.background, border: `1px solid var(--t-border)` }}>
            <p className="font-semibold" style={{ color: theme.primary }}>{b.label}</p>
            {b.venue && <p style={{ color: theme.text }}>{b.venue}</p>}
            {b.address && <p className="text-sm" style={{ color: theme.muted }}>{b.address}</p>}
            <div className="mt-3 flex flex-wrap justify-center gap-2">
              <a
                href={b.map_url || mapsLink(b)}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1.5 rounded-[var(--t-radius)] px-4 py-2 text-xs font-semibold text-white"
                style={{ background: theme.primary }}
              >
                <MapPin className="h-4 w-4" /> Buka Google Maps
              </a>
              <a
                href={directionsLink(b)}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1.5 rounded-[var(--t-radius)] border px-4 py-2 text-xs font-semibold"
                style={{ borderColor: "var(--t-border)", color: theme.primary }}
              >
                <Navigation className="h-4 w-4" /> Petunjuk Arah
              </a>
            </div>
          </div>
        ))}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- RSVP */

const rsvpSchema = z.object({
  attendance: z.enum(["attending", "not_attending", "maybe"]),
  guest_count: z.coerce.number().int().min(1).max(10),
  message: z.string().max(1000).optional(),
})

const ATTENDANCE_OPTIONS = [
  { value: "attending", label: "Hadir" },
  { value: "not_attending", label: "Tidak Hadir" },
  { value: "maybe", label: "Masih Ragu" },
] as const

export function RSVP({ model, theme }: SectionProps) {
  if (!model.sections.rsvp_enabled) return null
  if (!model.token || !model.eventId) {
    return (
      <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto">
          <h3 className="text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>RSVP</h3>
          <p className="mt-3" style={{ color: theme.muted }}>Konfirmasi kehadiran melalui tautan undangan pribadi Anda.</p>
        </div>
      </section>
    )
  }
  return <RsvpForm model={model} theme={theme} />
}

function RsvpForm({ model, theme }: SectionProps) {
  const [attendance, setAttendance] = useState("attending")
  const [count, setCount] = useState(1)
  const [message, setMessage] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [done, setDone] = useState(false)
  const [err, setErr] = useState("")

  async function submit(e: FormEvent) {
    e.preventDefault()
    setErr("")
    const parsed = rsvpSchema.safeParse({
      attendance,
      guest_count: count,
      message: message || undefined,
    })
    if (!parsed.success) {
      setErr(parsed.error.issues[0]?.message || "Data tidak valid")
      return
    }
    setSubmitting(true)
    try {
      const res = await fetch(`${apiBase()}/rsvp/${model.eventId}/submit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token: model.token, ...parsed.data }),
      })
      if (!res.ok) {
        const j = await res.json().catch(() => ({}))
        throw new Error(j?.error?.message || "Gagal mengirim RSVP")
      }
      setDone(true)
    } catch (e2) {
      setErr((e2 as Error).message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto">
        <h3 className="text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>RSVP</h3>
        {done ? (
          <div className="mt-6">
            <p className="text-lg" style={{ color: theme.text }}>
              Terima kasih, konfirmasi Anda telah kami terima.
            </p>
            {model.guestName && (
              <p className="mt-2 text-sm" style={{ color: theme.muted }}>Sampai jumpa, {model.guestName}!</p>
            )}
          </div>
        ) : (
          <form onSubmit={submit} className="mt-5 space-y-4 text-left">
            <div className="flex gap-2">
              {ATTENDANCE_OPTIONS.map((o) => (
                <button
                  key={o.value}
                  type="button"
                  onClick={() => setAttendance(o.value)}
                  className="flex-1 rounded-[var(--t-radius)] border px-3 py-2 text-sm"
                  style={{
                    borderColor: attendance === o.value ? theme.primary : "var(--t-border)",
                    color: attendance === o.value ? theme.primary : theme.muted,
                    background: attendance === o.value ? `${theme.primary}14` : "transparent",
                  }}
                >
                  {o.label}
                </button>
              ))}
            </div>
            <div>
              <label htmlFor="rsvp-count" className="block text-xs uppercase tracking-wider" style={{ color: theme.muted }}>
                Jumlah tamu
              </label>
              <input
                id="rsvp-count"
                type="number"
                min={1}
                max={10}
                value={count}
                onChange={(e) => setCount(Number(e.target.value))}
                className="w-full rounded-[var(--t-radius)] border px-3 py-2 text-sm"
                style={{ borderColor: "var(--t-border)", background: theme.background }}
              />
            </div>
            <div>
              <label htmlFor="rsvp-message" className="block text-xs uppercase tracking-wider" style={{ color: theme.muted }}>
                Ucapan (opsional)
              </label>
              <textarea
                id="rsvp-message"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                className="w-full rounded-[var(--t-radius)] border px-3 py-2 text-sm"
                style={{ borderColor: "var(--t-border)", background: theme.background }}
                rows={3}
              />
            </div>
            {err && <p role="alert" className="text-sm text-red-500">{err}</p>}
            <button
              type="submit"
              disabled={submitting}
              className="w-full rounded-[var(--t-radius)] px-4 py-2 text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: theme.primary }}
            >
              {submitting ? "Mengirim..." : "Kirim"}
            </button>
          </form>
        )}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- GIFT */

function CopyAccount({ value }: { value: string }) {
  const [copied, setCopied] = useState(false)
  async function copy() {
    try {
      await navigator.clipboard.writeText(value)
    } catch {
      const ta = document.createElement("textarea")
      ta.value = value
      ta.style.position = "fixed"
      ta.style.opacity = "0"
      document.body.appendChild(ta)
      ta.select()
      document.execCommand("copy")
      document.body.removeChild(ta)
    }
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <button
      type="button"
      onClick={copy}
      className="inline-flex items-center gap-1.5 rounded-[var(--t-radius)] border px-3 py-1.5 text-xs font-semibold"
      style={{ borderColor: "var(--t-border)", color: "var(--t-primary)" }}
      aria-label={copied ? "Nomor rekening tersalin" : "Salin nomor rekening"}
    >
      {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
      {copied ? "Tersalin" : "Salin"}
    </button>
  )
}

export function Gift({ model, theme, spec }: SectionProps) {
  const g = model.gift
  if (!model.sections.digital_gifts_enabled || !g) return null
  const title = (spec?.variant as string) === "envelope" ? "Amplop Digital" : "Hadiah"
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto space-y-4">
        <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{title}</h3>
        {g.gift_message && <p className="text-center" style={{ color: theme.muted }}>{g.gift_message}</p>}
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          {(g.bank_accounts || []).map((acc, i) => (
            <div key={i} className="rounded-[var(--t-radius)] border p-4 text-sm" style={{ borderColor: "var(--t-border)", background: theme.background }}>
              <p className="font-semibold" style={{ color: theme.primary }}>{(acc.bank as string) || (acc.bank_name as string) || "Bank"} - {(acc.name as string) || (acc.account_holder as string)}</p>
              <div className="mt-1 flex items-center justify-between gap-3">
                <p className="font-mono tracking-wide" style={{ color: theme.muted }}>{acc.account as string || acc.account_number as string}</p>
                <CopyAccount value={(acc.account as string) || (acc.account_number as string) || ""} />
              </div>
            </div>
          ))}
        </div>
        {g.ewallet && Object.keys(g.ewallet).length > 0 && (
          <div className="flex flex-wrap gap-3 text-sm">
            {Object.entries(g.ewallet).map(([k, v]) => (
              <span key={k} className="rounded-[var(--t-radius)] border px-3 py-1" style={{ borderColor: "var(--t-border)" }}>{(k as string).toUpperCase()}: {String(v)}</span>
            ))}
          </div>
        )}
        {g.qris_image_url && <img src={g.qris_image_url} alt="QRIS" className="mx-auto h-32 w-32 rounded-[var(--t-radius)]" loading="lazy" />}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- GUESTBOOK */

export function Guestbook({ model, theme }: SectionProps) {
  if (!model.sections.guestbook_enabled) return null
  return <GuestbookInner model={model} theme={theme} />
}

function GuestbookInner({ model, theme }: SectionProps) {
  const [items, setItems] = useState(model.guestbook)
  const [name, setName] = useState("")
  const [message, setMessage] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [err, setErr] = useState("")
  useEffect(() => setItems(model.guestbook), [model.guestbook])
  const canSubmit = !!model.eventId
  async function submit(e: FormEvent) {
    e.preventDefault()
    setErr("")
    setSubmitting(true)
    try {
      const res = await fetch(`${apiBase()}/guestbook/${model.eventId}/submit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, message }),
      })
      if (!res.ok) {
        const j = await res.json().catch(() => ({}))
        throw new Error(j?.error?.message || "Gagal mengirim pesan")
      }
      setName(""); setMessage("")
      setItems([{ name, message, created_at: new Date().toISOString() }, ...items])
    } catch (e2) {
      setErr((e2 as Error).message)
    } finally {
      setSubmitting(false)
    }
  }
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto space-y-4">
        <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Buku Tamu</h3>
        <ul className="space-y-3">
          {items.map((m, i) => (
            <li key={i} className="rounded-[var(--t-radius)] border p-4" style={{ borderColor: "var(--t-border)", background: theme.background }}>
              <p className="font-medium" style={{ color: theme.text }}>{m.name}</p>
              <p className="text-sm" style={{ color: theme.muted }}>{m.message}</p>
            </li>
          ))}
        </ul>
        <form onSubmit={submit} className="space-y-3">
          <input value={name} onChange={(e) => setName(e.target.value)} className="w-full rounded-[var(--t-radius)] border px-3 py-2 text-sm" style={{ borderColor: "var(--t-border)", background: theme.background }} placeholder="Nama" />
          <textarea value={message} onChange={(e) => setMessage(e.target.value)} className="w-full rounded-[var(--t-radius)] border px-3 py-2 text-sm" style={{ borderColor: "var(--t-border)", background: theme.background }} placeholder="Ucapan" rows={3} />
        {err && <p role="alert" className="text-sm text-red-500">{err}</p>}
        {canSubmit ? (
          <button type="submit" disabled={submitting} className="w-full rounded-[var(--t-radius)] px-4 py-2 text-sm font-semibold text-white disabled:opacity-60" style={{ background: theme.primary }}>
            {submitting ? "Mengirim..." : "Kirim"}
          </button>
        ) : (
          <p className="text-center text-sm" style={{ color: theme.muted }}>Ucapan dapat dikirim melalui tautan undangan.</p>
        )}
        </form>
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- CLOSING */

export function Closing({ model, theme }: SectionProps) {
  if (!model.closing) return null
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.background }}>
      <div style={narrow} className="mx-auto">
        <p className="text-lg leading-relaxed" style={{ color: theme.text }}>{model.closing}</p>
      </div>
    </section>
  )
}
