"use client"

import { motion, useReducedMotion, type Variants } from "framer-motion"
import type { ReactNode } from "react"
import type { RevealVariant } from "./types"

function buildVariants(variant: RevealVariant, distance: number): Variants {
  switch (variant) {
    case "scale":
      return {
        hidden: { opacity: 0, scale: 0.92 },
        show: { opacity: 1, scale: 1, transition: { duration: 0.6 } },
      }
    case "image-reveal":
      return {
        hidden: { opacity: 0, scale: 1.08 },
        show: { opacity: 1, scale: 1, transition: { duration: 0.9, ease: "easeOut" } },
      }
    case "text-reveal":
      return {
        hidden: { opacity: 0, y: distance, filter: "blur(6px)" },
        show: { opacity: 1, y: 0, filter: "blur(0px)", transition: { duration: 0.7 } },
      }
    case "fade":
      return {
        hidden: { opacity: 0 },
        show: { opacity: 1, transition: { duration: 0.6 } },
      }
    case "fade-up":
    default:
      return {
        hidden: { opacity: 0, y: distance },
        show: { opacity: 1, y: 0, transition: { duration: 0.6, ease: "easeOut" } },
      }
  }
}

export function Reveal({
  children,
  variant = "fade-up",
  distance = 40,
  delay = 0,
  as = "div",
  className,
  id,
}: {
  children: ReactNode
  variant?: RevealVariant
  distance?: number
  delay?: number
  as?: "div" | "section" | "li" | "header" | "footer"
  className?: string
  id?: string
}) {
  const reduce = useReducedMotion()
  const Tag = motion[as]
  if (reduce) {
    const Static = as
    return (
      <Static id={id} className={className}>
        {children}
      </Static>
    )
  }
  return (
    <Tag
      id={id}
      className={className}
      initial="hidden"
      whileInView="show"
      viewport={{ once: true, margin: "-80px" }}
      variants={buildVariants(variant, distance)}
      transition={{ delay }}
    >
      {children}
    </Tag>
  )
}

export function Stagger({
  children,
  className,
  id,
}: {
  children: ReactNode
  className?: string
  id?: string
}) {
  const reduce = useReducedMotion()
  if (reduce) return <div className={className}>{children}</div>
  return (
    <motion.div
      id={id}
      className={className}
      initial="hidden"
      whileInView="show"
      viewport={{ once: true, margin: "-60px" }}
      variants={{ show: { transition: { staggerChildren: 0.12 } } }}
    >
      {children}
    </motion.div>
  )
}

export function StaggerItem({
  children,
  variant = "fade-up",
  distance = 30,
  className,
}: {
  children: ReactNode
  variant?: RevealVariant
  distance?: number
  className?: string
}) {
  const reduce = useReducedMotion()
  if (reduce) return <div className={className}>{children}</div>
  return (
    <motion.div className={className} variants={buildVariants(variant, distance)}>
      {children}
    </motion.div>
  )
}
