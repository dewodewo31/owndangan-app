"use client"

import { useEffect, useState, type MouseEvent } from "react"
import { DEFAULT_GUEST_NAME } from "@/lib/guest"
import type { InvitationModel } from "@/templates/types"

/**
 * Full-screen invitation gate.
 *
 * While closed:
 *  - page scroll is locked (the invitation below must not be reachable)
 *  - a cover shows the couple, date and the guest's name
 *
 * Tapping "Buka Undangan" unlocks scroll and fires `invitation:open` so the
 * music player can attempt autoplay from a real user gesture (browser policy).
 */
export function CoverGate({
  model,
  heroImage,
  primaryColor,
  onOpen,
  lockBodyScroll = true,
}: {
  model: InvitationModel
  heroImage: string
  primaryColor: string
  onOpen?: () => void
  lockBodyScroll?: boolean
}) {
  const [open, setOpen] = useState(false)
  const guest = model.guestName || DEFAULT_GUEST_NAME

  useEffect(() => {
    if (!lockBodyScroll) return
    document.body.style.overflow = "hidden"
    return () => {
      document.body.style.overflow = ""
    }
  }, [lockBodyScroll])

  useEffect(() => {
    if (!open) return
    document.body.style.overflow = ""
    window.dispatchEvent(new Event("invitation:open"))
    onOpen?.()
  }, [open, onOpen])

  function handleOpen(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault()
    setOpen(true)
  }

  if (open) return null

  return (
    <div className="fixed inset-0 z-[60] overflow-y-auto">
      <div
        className="relative flex min-h-full flex-col items-center justify-center bg-cover bg-center px-6 py-16 text-center"
        style={{ backgroundImage: `url('${heroImage}')`, backgroundColor: primaryColor }}
      >
        <div className="absolute inset-0 bg-black/50" aria-hidden="true" />
        <div className="relative z-10 flex flex-col items-center text-white">
          <p className="text-xs uppercase tracking-[0.45em] text-white/80">The Wedding Of</p>
          <h1 className="mt-4 text-5xl leading-tight sm:text-6xl" style={{ fontFamily: "var(--t-font-heading, serif)" }}>
            {model.names.full}
          </h1>
          <div className="my-6 h-px w-24 bg-white/50" aria-hidden="true" />
          <p className="text-lg text-white/90">
            {model.date
              ? new Date(model.date).toLocaleDateString("id-ID", {
                  weekday: "long",
                  year: "numeric",
                  month: "long",
                  day: "numeric",
                })
              : ""}
          </p>
          <div className="mt-10">
            <p className="text-sm uppercase tracking-widest text-white/70">Kepada Yth.</p>
            <p className="mt-1 text-xl font-medium">{guest}</p>
          </div>
          <button
            onClick={handleOpen}
            className="mt-12 rounded-full px-10 py-3.5 text-sm font-semibold uppercase tracking-widest text-white shadow-xl transition-transform hover:scale-105 focus:outline-none focus-visible:ring-2 focus-visible:ring-white"
            style={{ backgroundColor: primaryColor }}
          >
            Buka Undangan
          </button>
        </div>
      </div>
    </div>
  )
}
