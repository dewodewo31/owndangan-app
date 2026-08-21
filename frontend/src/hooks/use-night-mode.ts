"use client"

import { useEffect, useState } from "react"

type Options = { start?: string; end?: string; timezone?: string }

function hourInZone(tz?: string): number {
  if (!tz) return new Date().getHours()
  try {
    const parts = new Intl.DateTimeFormat("en-US", {
      timeZone: tz,
      hour: "numeric",
      hour12: false,
    }).formatToParts(new Date())
    const h = parts.find((p) => p.type === "hour")?.value
    return h ? parseInt(h, 10) % 24 : new Date().getHours()
  } catch {
    return new Date().getHours()
  }
}

function isNight(hour: number, start = "18:00", end = "06:00"): boolean {
  const s = parseInt(start.slice(0, 2), 10)
  const e = parseInt(end.slice(0, 2), 10)
  if (s <= e) return hour >= s && hour < e
  return hour >= s || hour < e
}

// Auto night mode: true between `start` and `end` (default 18:00–06:00, wraps
// midnight). Hydration-safe: first render is always false (no window/Date on
// the server); the real value is set in an effect after mount.
export function useNightMode(opts: Options = {}): { night: boolean } {
  const [night, setNight] = useState(false)
  useEffect(() => {
    const update = () =>
      setNight(isNight(hourInZone(opts.timezone), opts.start, opts.end))
    update()
    const t = setInterval(update, 60_000)
    return () => clearInterval(t)
  }, [opts.start, opts.end, opts.timezone])
  return { night }
}
