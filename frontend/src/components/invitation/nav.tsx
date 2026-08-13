"use client"

import { useState } from "react"
import type { NavStyle, SectionSpec } from "@/templates/types"
import { cn } from "@/lib/cn"

const LABELS: Record<string, string> = {
  cover: "Cover",
  quote: "Sambutan",
  couple: "Mempelai",
  parents: "Orang Tua",
  countdown: "Hitungan",
  events: "Acara",
  gallery: "Galeri",
  location: "Lokasi",
  rsvp: "RSVP",
  gift: "Hadiah",
  guestbook: "Buku Tamu",
  closing: "Penutup",
}

export function Nav({
  style,
  sections,
}: {
  style: NavStyle
  sections: SectionSpec[]
}) {
  const items = sections.map((s) => ({ key: s.key, label: LABELS[s.key] || s.key }))
  const [open, setOpen] = useState(false)

  const linkCls =
    "text-xs uppercase tracking-wider transition-opacity hover:opacity-100"

  if (style === "side") {
    return (
      <nav className="fixed right-3 top-1/2 z-40 hidden -translate-y-1/2 flex-col gap-3 md:flex" aria-label="Navigasi undangan">
        {items.map((it) => (
          <a key={it.key} href={`#${it.key}`} className={cn(linkCls, "opacity-60")} style={{ writingMode: "vertical-rl" }}>
            {it.label}
          </a>
        ))}
      </nav>
    )
  }

  if (style === "overlay") {
    return (
      <>
        <button onClick={() => setOpen((v) => !v)} className="fixed right-4 top-4 z-50 flex h-11 w-11 items-center justify-center rounded-full bg-black text-white shadow-lg" aria-label="Buka menu">
          <span className="text-xl">{open ? "✕" : "☰"}</span>
        </button>
        {open && (
          <div className="fixed inset-0 z-40 flex flex-col items-center justify-center gap-5 bg-white/95">
            {items.map((it) => (
              <a key={it.key} href={`#${it.key}`} onClick={() => setOpen(false)} className="text-2xl font-semibold">
                {it.label}
              </a>
            ))}
          </div>
        )}
      </>
    )
  }

  if (style === "decorative-bottom") {
    return (
      <nav className="fixed inset-x-0 bottom-0 z-40 flex justify-around border-t bg-white/90 px-2 py-2 backdrop-blur" aria-label="Navigasi undangan">
        {items.slice(0, 5).map((it) => (
          <a key={it.key} href={`#${it.key}`} className="flex flex-col items-center gap-1 text-[10px] uppercase">
            <span className="h-1.5 w-1.5 rounded-full bg-current opacity-40" />
            {it.label}
          </a>
        ))}
      </nav>
    )
  }

  if (style === "floating-menu") {
    return (
      <nav className="fixed bottom-4 left-1/2 z-40 flex -translate-x-1/2 gap-1 rounded-full bg-white/90 px-3 py-2 shadow-lg backdrop-blur" aria-label="Navigasi undangan">
        {items.map((it) => (
          <a key={it.key} href={`#${it.key}`} className="rounded-full px-3 py-1 text-xs">
            {it.label}
          </a>
        ))}
      </nav>
    )
  }

  // bottom-floating (default)
  return (
    <nav className="fixed bottom-4 left-1/2 z-40 flex -translate-x-1/2 gap-1 rounded-full bg-black/80 px-3 py-2 text-white shadow-lg backdrop-blur" aria-label="Navigasi undangan">
      {items.map((it) => (
        <a key={it.key} href={`#${it.key}`} className={cn(linkCls, "opacity-70 px-2")}>
          {it.label}
        </a>
      ))}
    </nav>
  )
}
