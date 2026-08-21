"use client"

import { useState } from "react"
import { Share2, Check, MessageCircle, Send, Link2 } from "lucide-react"
import { buildInvitationUrl, DEFAULT_GUEST_NAME } from "@/lib/guest"
import { track } from "./sections"
import type { ThemeTokens } from "@/templates/types"

/**
 * "Bagikan Undangan" — Web Share API when available, otherwise copy URL /
 * WhatsApp / Telegram. The guest name (`?to=`) is preserved in the shared URL.
 */
export function ShareButton({
  slug,
  guestName,
  theme,
  eventId,
}: {
  slug: string
  guestName?: string
  theme: ThemeTokens
  eventId?: string
}) {
  const [copied, setCopied] = useState(false)
  const url = buildInvitationUrl(slug, guestName && guestName !== DEFAULT_GUEST_NAME ? guestName : undefined)

  async function handleShare() {
    if (typeof navigator !== "undefined" && navigator.share) {
      try {
        await navigator.share({ title: "Undangan Pernikahan", text: "Kami mengundang Anda", url })
        return
      } catch {
        // user cancelled or share unavailable — fall through to the fallbacks
      }
    }
    await copyUrl()
  }

  async function copyUrl() {
    try {
      await navigator.clipboard.writeText(url)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // clipboard may be unavailable (http); nothing to do
    }
  }

  const whatsapp = `https://wa.me/?text=${encodeURIComponent(`${url}`)}`
  const telegram = `https://t.me/share/url?url=${encodeURIComponent(url)}`

  const btn =
    "flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition-opacity hover:opacity-80"

  return (
    <div className="flex flex-wrap items-center justify-center gap-3">
      <button
        onClick={handleShare}
        className={btn}
        style={{ background: "var(--t-primary)", color: "var(--t-on-primary)" }}
        aria-label="Bagikan undangan"
      >
        <Share2 className="h-4 w-4" />
        Bagikan Undangan
      </button>

      <button onClick={copyUrl} className={btn} style={{ border: `1px solid var(--t-border)`, color: "var(--t-text)" }} aria-label="Salin tautan">
        {copied ? <Check className="h-4 w-4" /> : <Link2 className="h-4 w-4" />}
        {copied ? "Tautan disalin" : "Salin Tautan"}
      </button>

      <a href={whatsapp} target="_blank" rel="noreferrer" onClick={() => track("whatsapp_click", eventId)} className={btn} style={{ background: "#25D366", color: "#fff" }} aria-label="Bagikan via WhatsApp">
        <MessageCircle className="h-4 w-4" />
        WhatsApp
      </a>

      <a href={telegram} target="_blank" rel="noreferrer" className={btn} style={{ background: "#0088cc", color: "#fff" }} aria-label="Bagikan via Telegram">
        <Send className="h-4 w-4" />
        Telegram
      </a>
    </div>
  )
}
