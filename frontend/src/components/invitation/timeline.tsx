"use client"

import { useEffect, useRef, useState } from "react"
import { motion, useReducedMotion } from "framer-motion"
import { Flower2, Heart } from "lucide-react"
import { cn } from "@/lib/utils"
import type { TimelineStyle } from "@/templates/types"

export interface TimelineItem {
  year?: string
  date?: string
  title: string
  description?: string
  image?: string
}

interface TimelineProps {
  variant?: TimelineStyle
  items: TimelineItem[]
}

/* ------------------------------------------------------------------ */
/* Node                                                               */
/* ------------------------------------------------------------------ */

const NODE_GLYPH: Record<TimelineStyle, string> = {
  classic: "●",
  editorial: "◆",
  minimal: "·",
  romantic: "",
  luxury: "◇",
  botanical: "",
  modern: "●",
  zen: "○",
}

function Node({
  variant,
  active,
  index,
}: {
  variant: TimelineStyle
  active: boolean
  index: number
}) {
  const base = cn(
    "relative z-10 flex items-center justify-center transition-all duration-500",
    "text-[var(--t-primary)]",
    active ? "scale-125" : "scale-100"
  )
  const sizes: Record<TimelineStyle, string> = {
    classic: "h-4 w-4 rounded-full bg-[var(--t-primary)] text-transparent shadow-[0_0_0_4px_var(--t-surface)]",
    editorial: "h-3.5 w-3.5 rotate-45 bg-[var(--t-primary)] text-transparent",
    minimal: "h-2.5 w-2.5 rounded-full bg-[var(--t-primary)]",
    romantic: "h-5 w-5 text-[var(--t-primary)]",
    luxury: "h-5 w-5 text-[var(--t-primary)]",
    botanical: "h-5 w-5 text-[var(--t-primary)]",
    modern: "h-4 w-4 rounded-full bg-[var(--t-primary)] shadow-[0_0_14px_var(--t-primary)]",
    zen: "h-3 w-3 rounded-full border-2 border-[var(--t-primary)] bg-transparent",
  }
  return (
    <span
      className={cn(base, sizes[variant], active && "glow-node")}
      style={active ? { boxShadow: "0 0 0 6px color-mix(in srgb, var(--t-primary) 18%, transparent)" } : undefined}
      aria-hidden
    >
      {variant === "romantic" && (
        <Heart className="text-[inherit] h-4 w-4 fill-current" aria-hidden />
      )}
      {variant === "botanical" && (
        <Flower2 className="text-[inherit] h-4 w-4" aria-hidden />
      )}
      {variant === "luxury" && (
        <span className="text-[inherit] text-lg leading-none">{NODE_GLYPH[variant]}</span>
      )}
      {variant === "modern" && <span className="sr-only">{index + 1}</span>}
    </span>
  )
}

/* ------------------------------------------------------------------ */
/* Line                                                               */
/* ------------------------------------------------------------------ */

function Line({ variant, active }: { variant: TimelineStyle; active: boolean }) {
  const reduce = useReducedMotion()
  const color: Record<TimelineStyle, string> = {
    classic: "bg-[var(--t-border)]",
    editorial: "bg-[var(--t-border)]",
    minimal: "bg-[var(--t-border)]",
    romantic: "bg-gradient-to-b from-[var(--t-accent)] via-[var(--t-primary)] to-[var(--t-accent)]",
    luxury: "bg-gradient-to-b from-transparent via-[var(--t-accent)] to-transparent",
    botanical: "bg-[var(--t-border)]",
    modern: "bg-[var(--t-border)]",
    zen: "bg-[var(--t-border)]",
  }
  const pos: Record<TimelineStyle, string> = {
    classic: "left-4 md:left-1/2 -translate-x-1/2",
    editorial: "left-4 md:left-1/2 -translate-x-1/2",
    minimal: "left-4 -translate-x-1/2",
    romantic: "left-4 md:left-1/2 -translate-x-1/2",
    luxury: "left-4 md:left-1/2 -translate-x-1/2",
    botanical: "left-4 md:left-1/2 -translate-x-1/2",
    modern: "left-4 md:left-1/2 -translate-x-1/2",
    zen: "left-4 -translate-x-1/2",
  }
  const width: Record<TimelineStyle, string> = {
    classic: "w-px",
    editorial: "w-px",
    minimal: "w-px",
    romantic: "w-[2px]",
    luxury: "w-[2px]",
    botanical: "w-[2px]",
    modern: "w-[2px]",
    zen: "w-px",
  }
  return (
    <div className={cn("absolute top-0 bottom-0", pos[variant], width[variant], color[variant])}>
      <motion.div
        className="absolute inset-0 origin-top bg-[var(--t-primary)]"
        initial={{ scaleY: 0 }}
        animate={reduce ? undefined : { scaleY: 1 }}
        transition={{ duration: 1.4, ease: [0.22, 1, 0.36, 1] }}
        style={reduce ? { scaleY: 1 } : undefined}
      />
      {active && variant === "botanical" && (
        <Flower2 className="absolute left-1/2 top-0 h-4 w-4 -translate-x-1/2 text-[var(--t-accent)]" aria-hidden />
      )}
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Hooks                                                              */
/* ------------------------------------------------------------------ */

function useActiveIndex(containerRef: React.RefObject<HTMLDivElement | null>, count: number): number {
  const [active, setActive] = useState(-1)
  useEffect(() => {
    if (count === 0) return
    const root = containerRef.current
    if (!root) return
    const els = Array.from(root.querySelectorAll<HTMLElement>("[data-idx]"))
    const io = new IntersectionObserver(
      (entries) => {
        entries.forEach((e) => {
          if (e.isIntersecting) setActive(Number((e.target as HTMLElement).dataset.idx))
        })
      },
      { rootMargin: "-40% 0px -40% 0px", threshold: 0 }
    )
    els.forEach((el) => io.observe(el))
    return () => io.disconnect()
  }, [count, containerRef])
  return active
}

/* ------------------------------------------------------------------ */
/* Item row (shared shell, variant classes)                           */
/* ------------------------------------------------------------------ */

function Card({ item, variant, align }: { item: TimelineItem; variant: TimelineStyle; align: "left" | "right" }) {
  const reduce = useReducedMotion()
  const cardCls: Record<TimelineStyle, string> = {
    classic: "rounded-[var(--t-radius)] border border-[var(--t-border)] bg-[var(--t-surface)] p-5 shadow-sm",
    editorial: "border-t-2 border-[var(--t-primary)] bg-transparent p-5",
    minimal: "border-l-2 border-[var(--t-primary)] bg-transparent pl-4",
    romantic:
      "rounded-[1.5rem] border border-[var(--t-border)] bg-[var(--t-surface)] p-5 shadow-[0_10px_30px_-15px_color-mix(in_srgb,var(--t-primary)_30%,transparent)]",
    luxury: "border border-[var(--t-border)] bg-[var(--t-surface)] p-6",
    botanical: "rounded-[2rem] border border-[var(--t-border)] bg-[var(--t-surface)] p-5",
    modern: "rounded-xl border border-[var(--t-border)] bg-[var(--t-surface)] p-5 shadow-sm",
    zen: "bg-transparent",
  }
  const titleAlign = align === "right" ? "md:text-right" : "md:text-left"
  const yearSize: Record<TimelineStyle, string> = {
    classic: "text-lg",
    editorial: "text-2xl",
    minimal: "text-4xl",
    romantic: "text-xl",
    luxury: "text-5xl",
    botanical: "text-xl",
    modern: "text-5xl",
    zen: "text-6xl",
  }
  return (
    <div
      className={cn("group", cardCls[variant], align === "right" && "md:text-right")}
      style={
        variant === "modern" && align === "right"
          ? { marginTop: "2.5rem" }
          : variant === "luxury" && align === "right"
            ? { marginTop: "4rem" }
            : undefined
      }
    >
      {item.image && (
        <div className={cn("overflow-hidden", variant === "luxury" ? "mb-5" : "mb-4", variant === "zen" && "mb-6")}>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={item.image}
            alt={item.title}
            loading="lazy"
            className={cn(
              "h-44 w-full object-cover transition-transform duration-700 group-hover:scale-105",
              variant === "classic" && "rounded-[var(--t-radius)]",
              variant === "romantic" && "rounded-[1.25rem]",
              variant === "botanical" && "rounded-[1.5rem]",
              variant === "luxury" && "h-56 md:h-64",
              variant === "minimal" && "hidden",
              variant === "zen" && "hidden"
            )}
          />
        </div>
      )}
      {(item.year || item.date) && (
        <div
          className={cn(
            "mb-1 flex items-center gap-2 font-[var(--t-font-accent)] text-[var(--t-primary)]",
            yearSize[variant],
            align === "right" && "md:justify-end"
          )}
        >
          {item.year && <span>{item.year}</span>}
          {item.date && (
            <span className={cn("text-xs uppercase tracking-[0.2em] text-[var(--t-muted)]", variant !== "luxury" && variant !== "modern" && variant !== "zen" && "text-base")}>
              {item.date}
            </span>
          )}
        </div>
      )}
      <h4 className={cn("font-[var(--t-font-heading)] text-[var(--t-primary)]", variant === "luxury" || variant === "modern" || variant === "zen" ? "text-2xl" : "text-xl", titleAlign)}>
        {item.title}
      </h4>
      {item.description && (
        <p className={cn("mt-2 text-sm leading-relaxed text-[var(--t-muted)]", align === "right" && "md:text-right")}>{item.description}</p>
      )}
    </div>
  )
}

function Item({
  item,
  variant,
  index,
  active,
  reduce,
}: {
  item: TimelineItem
  variant: TimelineStyle
  index: number
  active: boolean
  reduce: boolean | null
}) {
  const alternating = variant === "classic" || variant === "editorial" || variant === "romantic" || variant === "botanical"
  const single = variant === "minimal" || variant === "zen"
  const isRight = !single && index % 2 === 1

  const sideClass = single
    ? "pl-12 md:pl-12"
    : alternating
      ? cn("pl-12", isRight ? "md:pl-[calc(50%+2.5rem)] md:pr-0" : "md:pr-[calc(50%+2.5rem)] md:pl-0")
      : cn("pl-12", isRight ? "md:pl-[calc(50%+2.5rem)]" : "md:pr-[calc(50%+2.5rem)] md:pl-0")

  return (
    <div
      data-idx={index}
      className={cn("relative", sideClass)}
    >
      <div className="absolute left-4 top-2 -translate-x-1/2 md:left-1/2">
        <Node variant={variant} active={active} index={index} />
      </div>
      <motion.div
        initial={{ opacity: 0, y: 24 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true, margin: "-60px" }}
        transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
        className={cn("transition-transform duration-500", !reduce && "hover:-translate-y-1")}
      >
        <Card item={item} variant={variant} align={isRight ? "right" : "left"} />
      </motion.div>
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Timeline                                                           */
/* ------------------------------------------------------------------ */

export function Timeline({ variant = "classic", items }: TimelineProps) {
  const reduce = useReducedMotion()
  const containerRef = useRef<HTMLDivElement | null>(null)
  const active = useActiveIndex(containerRef, items.length)
  if (items.length === 0) return null

  const wrap: Record<TimelineStyle, string> = {
    classic: "md:max-w-4xl",
    editorial: "md:max-w-5xl",
    minimal: "max-w-2xl",
    romantic: "md:max-w-4xl",
    luxury: "md:max-w-5xl",
    botanical: "md:max-w-4xl",
    modern: "md:max-w-5xl",
    zen: "max-w-2xl",
  }

  return (
    <div ref={containerRef} className={cn("relative mx-auto", wrap[variant])}>
      <Line variant={variant} active={active >= 0} />
      <div className="space-y-10 md:space-y-14">
        {items.map((it, i) => (
          <Item key={it.title + i} item={it} variant={variant} index={i} active={active === i} reduce={reduce} />
        ))}
      </div>
    </div>
  )
}
