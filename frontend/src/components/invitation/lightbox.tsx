"use client"

import { useCallback, useEffect, useState, type TouchEvent } from "react"
import { ChevronLeft, ChevronRight, X } from "lucide-react"

/**
 * Minimal fullscreen image lightbox with prev/next, Esc/close button and
 * touch swipe support. Used by the gallery sections.
 */
export function Lightbox({
  images,
  index,
  onClose,
  onIndexChange,
}: {
  images: { image_url: string; caption?: string }[]
  index: number | null
  onClose: () => void
  onIndexChange: (i: number) => void
}) {
  const [touchX, setTouchX] = useState<number | null>(null)

  const prev = useCallback(
    () => onIndexChange((index! - 1 + images.length) % images.length),
    [images.length, index, onIndexChange]
  )
  const next = useCallback(
    () => onIndexChange((index! + 1) % images.length),
    [images.length, index, onIndexChange]
  )

  useEffect(() => {
    if (index === null) return
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose()
      if (e.key === "ArrowLeft") prev()
      if (e.key === "ArrowRight") next()
    }
    window.addEventListener("keydown", onKey)
    return () => window.removeEventListener("keydown", onKey)
  }, [index, onClose, prev, next])

  useEffect(() => {
    if (index !== null) {
      document.body.style.overflow = "hidden"
      return () => {
        document.body.style.overflow = ""
      }
    }
  }, [index])

  if (index === null) return null
  const img = images[index]

  function onTouchStart(e: TouchEvent) {
    setTouchX(e.touches[0].clientX)
  }
  function onTouchEnd(e: TouchEvent) {
    if (touchX === null) return
    const dx = e.changedTouches[0].clientX - touchX
    if (Math.abs(dx) > 40) {
      if (dx < 0) next()
      else prev()
    }
    setTouchX(null)
  }

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-label="Pratinjau foto"
      className="fixed inset-0 z-[70] flex flex-col items-center justify-center bg-black/95"
      onClick={onClose}
      onTouchStart={onTouchStart}
      onTouchEnd={onTouchEnd}
    >
      <button
        onClick={onClose}
        aria-label="Tutup"
        className="absolute right-4 top-4 z-10 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white hover:bg-white/20"
      >
        <X className="h-5 w-5" />
      </button>

      <button
        onClick={(e) => {
          e.stopPropagation()
          prev()
        }}
        aria-label="Sebelumnya"
        className="absolute left-2 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white hover:bg-white/20 sm:left-4"
      >
        <ChevronLeft className="h-6 w-6" />
      </button>

      <figure className="flex max-h-full max-w-full flex-col items-center" onClick={(e) => e.stopPropagation()}>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={img.image_url}
          alt={img.caption || "Galeri"}
          className="max-h-[80vh] max-w-full object-contain"
        />
        {img.caption && <figcaption className="mt-3 px-6 text-center text-sm text-white/80">{img.caption}</figcaption>}
        <p className="mt-2 text-xs text-white/50">
          {index + 1} / {images.length}
        </p>
      </figure>

      <button
        onClick={(e) => {
          e.stopPropagation()
          next()
        }}
        aria-label="Berikutnya"
        className="absolute right-2 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white hover:bg-white/20 sm:right-4"
      >
        <ChevronRight className="h-6 w-6" />
      </button>
    </div>
  )
}
