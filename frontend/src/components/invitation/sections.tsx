import { useState, useEffect, type FormEvent } from "react"
import type { InvitationModel, SectionSpec, ThemeTokens } from "@/templates/types"
import { cn } from "@/lib/cn"

export interface SectionProps {
  model: InvitationModel
  theme: ThemeTokens
  spec?: SectionSpec
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
  if (!model.date) return null
  const target = new Date(model.date + "T" + (model.time || "09:00:00"))
  const now = new Date()
  const diff = Math.max(0, target.getTime() - now.getTime())
  const days = Math.floor(diff / 86400000)
  const hours = Math.floor((diff % 86400000) / 3600000)
  const mins = Math.floor((diff % 3600000) / 60000)
  const secs = Math.floor((diff % 60000) / 1000)
  const items = [
    { v: days, l: "Hari" },
    { v: hours, l: "Jam" },
    { v: mins, l: "Menit" },
    { v: secs, l: "Detik" },
  ]
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.background }}>
      <div style={narrow} className="mx-auto">
        <p className="mb-5 text-sm uppercase tracking-[0.3em]" style={{ color: theme.muted }}>Menuju Hari Bahagia</p>
        <div className="flex justify-center gap-4">
          {items.map((it) => (
            <div key={it.l} className="min-w-[64px] rounded-[var(--t-radius)] px-3 py-4" style={{ background: theme.surface, border: `1px solid var(--t-border)` }}>
              <div className="text-3xl font-semibold" style={{ color: theme.primary }}>{String(it.v).padStart(2, "0")}</div>
              <div className="text-xs" style={{ color: theme.muted }}>{it.l}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- EVENTS */

export function Events({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as string) || "cards"
  const blocks = [model.events.akad, model.events.resepsi].filter(Boolean) as NonNullable<InvitationModel["events"]["akad"]>[]
  if (blocks.length === 0) return null

  if (variant === "timeline") {
    return (
      <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto space-y-6">
          <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Akad & Resepsi</h3>
          {blocks.map((b, i) => (
            <div key={i} className="flex gap-4 border-l-2 pl-4" style={{ borderColor: theme.accent }}>
              <div>
                <p className="font-semibold" style={{ color: theme.primary }}>{b.label}</p>
                {b.venue && <p style={{ color: theme.text }}>{b.venue}</p>}
                {(b.date || b.time) && <p className="text-sm" style={{ color: theme.muted }}>{b.date || fmtDate(model.date)} · {b.time || model.time}</p>}
                {b.address && <p className="text-sm" style={{ color: theme.muted }}>{b.address}</p>}
                {b.map_url && <a href={b.map_url} target="_blank" rel="noreferrer" className="text-sm underline" style={{ color: theme.primary }}>Lihat peta</a>}
              </div>
            </div>
          ))}
        </div>
      </section>
    )
  }

  if (variant === "side") {
    return (
      <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
        <div style={narrow} className="mx-auto grid grid-cols-1 gap-6 sm:grid-cols-2">
          {blocks.map((b, i) => (
            <div key={i} className="rounded-[var(--t-radius)] p-6 text-center" style={{ background: theme.background, border: `1px solid var(--t-border)` }}>
              <p className="text-lg font-semibold" style={{ color: theme.primary }}>{b.label}</p>
              {b.venue && <p className="mt-2" style={{ color: theme.text }}>{b.venue}</p>}
              <p className="text-sm" style={{ color: theme.muted }}>{b.date || fmtDate(model.date)} · {b.time || model.time}</p>
              {b.address && <p className="text-sm" style={{ color: theme.muted }}>{b.address}</p>}
            </div>
          ))}
        </div>
      </section>
    )
  }

  // cards (default)
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto grid grid-cols-1 gap-6 sm:grid-cols-2">
        {blocks.map((b, i) => (
          <div key={i} className="rounded-[var(--t-radius)] p-6 text-center shadow-sm" style={{ background: theme.background, border: `1px solid var(--t-border)` }}>
            <p className="text-xs uppercase tracking-widest" style={{ color: theme.muted }}>{b.label}</p>
            <h4 className="mt-2 text-2xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>{b.venue || "—"}</h4>
            <p className="mt-2 text-sm" style={{ color: theme.muted }}>{b.date || fmtDate(model.date)} · {b.time || model.time}</p>
            {b.address && <p className="mt-1 text-sm" style={{ color: theme.muted }}>{b.address}</p>}
            {b.map_url && <a href={b.map_url} target="_blank" rel="noreferrer" className="mt-2 inline-block text-sm underline" style={{ color: theme.primary }}>Lihat peta</a>}
          </div>
        ))}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- GALLERY */

export function Gallery({ model, theme, spec }: SectionProps) {
  const variant = (spec?.variant as string) || "grid"
  if (model.gallery.length === 0) return null

  if (variant === "columns") {
    return (
      <section className="px-4 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Galeri</h3>
        <div className="columns-2 gap-3 sm:columns-3" style={{ maxWidth: "var(--t-content-width)", margin: "0 auto" }}>
          {model.gallery.map((g, i) => (
            <img key={i} src={g.image_url} alt={g.caption || "Galeri"} className="mb-3 w-full break-inside-avoid rounded-[var(--t-radius)] object-cover" loading="lazy" />
          ))}
        </div>
      </section>
    )
  }

  if (variant === "masonry") {
    return (
      <section className="px-4 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <div className="columns-2 gap-4" style={{ maxWidth: "var(--t-content-width)", margin: "0 auto" }}>
          {model.gallery.map((g, i) => (
            <figure key={i} className="mb-4 break-inside-avoid">
              <img src={g.image_url} alt={g.caption || "Galeri"} className="w-full rounded-[var(--t-radius)] object-cover" loading="lazy" />
              {g.caption && <figcaption className="mt-1 text-center text-xs" style={{ color: theme.muted }}>{g.caption}</figcaption>}
            </figure>
          ))}
        </div>
      </section>
    )
  }

  if (variant === "horizontal") {
    return (
      <section className="py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
        <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Galeri</h3>
        <div className="flex snap-x gap-4 overflow-x-auto px-6 pb-2">
          {model.gallery.map((g, i) => (
            <img key={i} src={g.image_url} alt={g.caption || "Galeri"} className="h-56 w-44 flex-none snap-center rounded-[var(--t-radius)] object-cover" loading="lazy" />
          ))}
        </div>
      </section>
    )
  }

  // grid (default)
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.background }}>
      <h3 className="mb-6 text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Galeri</h3>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3" style={{ maxWidth: "var(--t-content-width)", margin: "0 auto" }}>
        {model.gallery.map((g, i) => (
          <img key={i} src={g.image_url} alt={g.caption || "Galeri"} className="aspect-square w-full rounded-[var(--t-radius)] object-cover" loading="lazy" />
        ))}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- LOCATION */

export function Location({ model, theme }: SectionProps) {
  const blocks = [model.events.akad, model.events.resepsi].filter(Boolean) as NonNullable<InvitationModel["events"]["akad"]>[]
  if (blocks.length === 0) return null
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto space-y-4">
        <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Lokasi</h3>
        {blocks.map((b, i) => (
          <div key={i} className="rounded-[var(--t-radius)] p-5 text-center" style={{ background: theme.background, border: `1px solid var(--t-border)` }}>
            <p className="font-semibold" style={{ color: theme.primary }}>{b.label}</p>
            {b.venue && <p style={{ color: theme.text }}>{b.venue}</p>}
            {b.address && <p className="text-sm" style={{ color: theme.muted }}>{b.address}</p>}
            {b.map_url && <a href={b.map_url} target="_blank" rel="noreferrer" className="mt-2 inline-block text-sm underline" style={{ color: theme.primary }}>Buka di Maps</a>}
          </div>
        ))}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- RSVP */

export function RSVP({ model, theme }: SectionProps) {
  if (!model.sections.rsvp_enabled) return null
  if (!model.token) {
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
  const [attendance, setAttendance] = useState("yes")
  const [count, setCount] = useState(1)
  const [message, setMessage] = useState("")
  const [done, setDone] = useState(false)
  const [err, setErr] = useState("")
  async function submit(e: FormEvent) {
    e.preventDefault()
    setErr("")
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"}/public/rsvps`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token: model.token, attendance, guest_count: count, message }),
      })
      if (!res.ok) {
        const j = await res.json().catch(() => ({}))
        throw new Error(j?.error?.message || "Gagal mengirim RSVP")
      }
      setDone(true)
    } catch (e2) {
      setErr((e2 as Error).message)
    }
  }
  return (
    <section className="px-6 py-[var(--t-section-spacing)] text-center" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto">
        <h3 className="text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>RSVP</h3>
        {done ? (
          <p className="mt-4" style={{ color: theme.text }}>Terima kasih, konfirmasi Anda telah kami terima.</p>
        ) : (
          <form onSubmit={submit} className="mt-5 space-y-4 text-left">
            <div className="flex gap-2">
              {["yes", "no", "maybe"].map((a) => (
                <button key={a} type="button" onClick={() => setAttendance(a)} className="flex-1 rounded-[var(--t-radius)] border px-3 py-2 text-sm" style={{ borderColor: attendance === a ? theme.primary : "var(--t-border)", color: attendance === a ? theme.primary : theme.muted }}>
                  {a === "yes" ? "Hadir" : a === "no" ? "Tidak" : "Ragu"}
                </button>
              ))}
            </div>
            <input type="number" min={1} value={count} onChange={(e) => setCount(Number(e.target.value))} className="w-full rounded-[var(--t-radius)] border px-3 py-2 text-sm" style={{ borderColor: "var(--t-border)", background: theme.background }} placeholder="Jumlah tamu" />
            <textarea value={message} onChange={(e) => setMessage(e.target.value)} className="w-full rounded-[var(--t-radius)] border px-3 py-2 text-sm" style={{ borderColor: "var(--t-border)", background: theme.background }} placeholder="Ucapan (opsional)" rows={3} />
            {err && <p className="text-sm text-red-500">{err}</p>}
            <button type="submit" className="w-full rounded-[var(--t-radius)] px-4 py-2 text-sm font-semibold text-white" style={{ background: theme.primary }}>Kirim</button>
          </form>
        )}
      </div>
    </section>
  )
}

/* ---------------------------------------------------------------- GIFT */

export function Gift({ model, theme }: SectionProps) {
  const g = model.gift
  if (!model.sections.digital_gifts_enabled || !g) return null
  return (
    <section className="px-6 py-[var(--t-section-spacing)]" style={{ background: theme.surface }}>
      <div style={narrow} className="mx-auto space-y-4">
        <h3 className="text-center text-3xl" style={{ fontFamily: "var(--t-font-heading)", color: theme.primary }}>Hadiah</h3>
        {g.gift_message && <p className="text-center" style={{ color: theme.muted }}>{g.gift_message}</p>}
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          {(g.bank_accounts || []).map((acc, i) => (
            <div key={i} className="rounded-[var(--t-radius)] border p-4 text-sm" style={{ borderColor: "var(--t-border)", background: theme.background }}>
              <p className="font-semibold" style={{ color: theme.primary }}>{(acc.bank as string) || (acc.bank_name as string) || "Bank"} - {(acc.name as string) || (acc.account_holder as string)}</p>
              <p style={{ color: theme.muted }}>{acc.account as string || acc.account_number as string}</p>
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
  const [err, setErr] = useState("")
  useEffect(() => setItems(model.guestbook), [model.guestbook])
  async function submit(e: FormEvent) {
    e.preventDefault()
    setErr("")
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"}/public/guestbook`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ slug: model.slug, name, message }),
      })
      if (!res.ok) {
        const j = await res.json().catch(() => ({}))
        throw new Error(j?.error?.message || "Gagal mengirim pesan")
      }
      setName(""); setMessage("")
      setItems([{ name, message, created_at: new Date().toISOString() }, ...items])
    } catch (e2) {
      setErr((e2 as Error).message)
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
          {err && <p className="text-sm text-red-500">{err}</p>}
          <button type="submit" className="w-full rounded-[var(--t-radius)] px-4 py-2 text-sm font-semibold text-white" style={{ background: theme.primary }}>Kirim</button>
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
