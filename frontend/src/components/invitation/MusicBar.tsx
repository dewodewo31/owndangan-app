"use client"

import { useEffect, useRef, useState } from "react"
import { Volume2, VolumeX, Play, Pause, Music } from "lucide-react"
import type { ThemeTokens } from "@/templates/types"
import { cn } from "@/lib/cn"

export function MusicBar({
  music,
  theme,
  placement = "top-left",
}: {
  music: { title: string; file_url?: string }
  theme: ThemeTokens
  placement?: "top-left" | "bottom-left"
}) {
  const ref = useRef<HTMLAudioElement>(null)
  const [playing, setPlaying] = useState(false)
  const [muted, setMuted] = useState(false)

  // `invitation:open` fires from the cover gate's tap — the legal autoplay gesture
  useEffect(() => {
    function tryPlay() {
      const a = ref.current
      if (!a) return
      a.play()
        .then(() => setPlaying(true))
        .catch(() => setPlaying(false))
    }
    window.addEventListener("invitation:open", tryPlay)
    return () => window.removeEventListener("invitation:open", tryPlay)
  }, [])

  function toggle() {
    const a = ref.current
    if (!a) return
    if (playing) {
      a.pause()
      setPlaying(false)
    } else {
      a.play()
        .then(() => setPlaying(true))
        .catch(() => setPlaying(false))
    }
  }

  function toggleMute() {
    const a = ref.current
    if (!a) return
    a.muted = !a.muted
    setMuted(a.muted)
  }

  return (
    <div
      className={cn(
        "fixed z-50 flex items-center gap-2.5 rounded-full py-1.5 pl-1.5 pr-2 shadow-lg",
        placement === "bottom-left" ? "bottom-4 left-4" : "left-4 top-4"
      )}
      style={{ background: "rgba(255,255,255,0.92)", border: `1px solid ${theme.border}` }}
    >
      <span
        className={cn(
          "flex h-10 w-10 flex-none items-center justify-center rounded-full",
          playing && "animate-spin motion-reduce:animate-none"
        )}
        style={{ background: "var(--t-primary)", color: "var(--t-on-primary)" }}
        aria-hidden
      >
        <Music className="h-4 w-4" />
      </span>
      <span className="flex min-w-0 flex-col">
        <span className="text-[10px] uppercase tracking-wider" style={{ color: "var(--t-muted)" }}>
          Wedding Music
        </span>
        <span className="max-w-[10rem] truncate text-sm font-medium" style={{ color: "var(--t-text)" }}>
          {music.title}
        </span>
      </span>
      <button
        onClick={toggle}
        className="flex h-10 w-10 flex-none items-center justify-center rounded-full transition-transform hover:scale-105"
        style={{ background: "var(--t-primary)", color: "var(--t-on-primary)" }}
        aria-label={playing ? "Hentikan musik" : "Putar musik"}
        title={music.title}
      >
        {playing ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
      </button>
      <button
        onClick={toggleMute}
        className="flex h-8 w-8 flex-none items-center justify-center rounded-full text-black/70 transition-colors hover:text-black"
        aria-label={muted ? "Aktifkan suara" : "Matikan suara"}
      >
        {muted ? <VolumeX className="h-4 w-4" /> : <Volume2 className="h-4 w-4" />}
      </button>
      <audio ref={ref} src={music.file_url} loop preload="none" />
    </div>
  )
}