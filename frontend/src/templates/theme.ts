import type { CSSProperties } from "react"
import type { ThemeTokens } from "./types"

// Map theme tokens to CSS custom properties applied on the invitation root.
export function themeVars(theme: ThemeTokens): CSSProperties {
  return {
    "--t-primary": theme.primary,
    "--t-secondary": theme.secondary,
    "--t-background": theme.background,
    "--t-surface": theme.surface,
    "--t-text": theme.text,
    "--t-muted": theme.muted,
    "--t-accent": theme.accent,
    "--t-border": theme.border,
    "--t-radius": theme.radius,
    "--t-font-heading": theme.fontHeading,
    "--t-font-body": theme.fontBody,
    "--t-font-accent": theme.fontAccent,
    "--t-section-spacing": theme.sectionSpacing,
    "--t-content-width": theme.contentWidth,
    "--t-hero-height": theme.heroHeight,
    "--t-anim-duration": theme.animationDuration,
    "--t-anim-easing": theme.animationEasing,
    "--t-reveal-distance": theme.revealDistance,
    backgroundColor: theme.background,
    color: theme.text,
    fontFamily: theme.fontBody,
  } as CSSProperties
}

// All font families referenced by the template library.
export const FONT_FAMILIES = [
  "Playfair Display",
  "Cormorant Garamond",
  "DM Serif Display",
  "Petrona",
  "Lora",
  "Inter",
  "Montserrat",
  "Poppins",
  "Jost",
  "Dancing Script",
  "Great Vibes",
  "Pinyon Script",
  "Shippori Mincho",
  "Zen Maru Gothic",
  "Caveat",
]

export const GOOGLE_FONTS_HREF =
  "https://fonts.googleapis.com/css2?family=" +
  FONT_FAMILIES.map(
    (f) => f.replace(/ /g, "+") + ":wght@400;500;600;700&display=swap"
  ).join("&family=")
