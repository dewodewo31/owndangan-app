"use client"

import { useRef, useState } from "react"
import type { ThemeTokens } from "@/templates/types"

export function MusicBar({
  music,
  theme,
}: {
  music: { title: string; file_url?: string }
  theme: ThemeTokens
}) {
  const ref = useRef<HTMLAudioElement>(null)
  const [playing, setPlaying] = useState(false)
  function toggle() {
    const a = ref.current
    if (!a) return
    if (playing) {
      a.pause()
      setPlaying(false)
    } else {
      a.play().catch(() => setPlaying(false))
      setPlaying(true)
    }
  }
  return (
    <button
      onClick={toggle}
      className="fixed left-4 top-4 z-50 flex h-11 w-11 items-center justify-center rounded-full text-white shadow-lg"
      style={{ background: theme.primary }}
      aria-label={playing ? "Hentikan musik" : "Putar musik"}
      title={music.title}
    >
      <span className="text-lg">{playing ? "⏸" : "♪"}</span>
      <audio ref={ref} src={music.file_url} loop preload="none" />
    </button>
  )
}
