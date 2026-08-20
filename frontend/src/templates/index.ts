import type { TemplateDefinition } from "./types"
import { definition as contemporaryEditorial } from "./contemporary-editorial"
import { definition as corporate } from "./corporate"
import { definition as islamic } from "./islamic"
import { definition as japaneseZen } from "./japanese-zen"
import { definition as javanese } from "./javanese"
import { definition as luxuryBlackGold } from "./luxury-black-gold"
import { definition as modernBotanical } from "./modern-botanical"
import { definition as modernMinimalist } from "./modern-minimalist"
import { definition as romanticElegant } from "./romantic-elegant"
import { definition as rusticBohemian } from "./rustic-bohemian"
import { definition as sundanese } from "./sundanese"

export const TEMPLATES: TemplateDefinition[] = [
  contemporaryEditorial,
  corporate,
  islamic,
  japaneseZen,
  javanese,
  luxuryBlackGold,
  modernBotanical,
  modernMinimalist,
  romanticElegant,
  rusticBohemian,
  sundanese,
]

export const OCCASION_LABELS: Record<string, string> = {
  pernikahan: "Pernikahan",
  prewedding: "Prewedding",
  "ulang-tahun": "Ulang Tahun",
  aqiqah: "Aqiqah",
  khitanan: "Khitanan",
  anniversary: "Anniversary",
  graduation: "Graduation",
  corporate: "Corporate/Meeting",
  event: "Event",
}

export function templateOccasions(): string[] {
  const seen: string[] = []
  for (const t of TEMPLATES) for (const o of t.occasions) if (!seen.includes(o)) seen.push(o)
  return seen
}

/**
 * Pick a template for an event.
 *
 * Preference order:
 *  1. exact match on the backend template name (e.g. "Elegan Klasik")
 *  2. fallback by backend template group: standard/premium/all -> a curated
 *     template in the same visual tier
 *  3. modern-minimalist (safe default, exists in every data set)
 */
export function selectTemplate(groupName?: string, templateName?: string): TemplateDefinition {
  const byName = templateName
    ? TEMPLATES.find((t) => t.name.toLowerCase() === templateName.toLowerCase())
    : undefined
  if (byName) return byName

  const tier: Record<string, string> = {
    standard: "modern-minimalist",
    premium: "romantic-elegant",
    all: "luxury-black-gold",
  }
  const fallbackKind = groupName ? tier[groupName] : undefined
  if (fallbackKind) {
    const t = TEMPLATES.find((x) => x.kind === fallbackKind)
    if (t) return t
  }
  return TEMPLATES.find((x) => x.kind === "modern-minimalist") ?? TEMPLATES[0]
}
