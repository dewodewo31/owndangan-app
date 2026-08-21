import type { CSSProperties } from "react"
import type { ThemeTokens } from "./types"

// Map theme tokens to CSS custom properties applied on the invitation root.
// Pass the template's night palette to overlay it over the day tokens.
export function themeVars(theme: ThemeTokens, night?: Partial<ThemeTokens>): CSSProperties {
  const tokens = night ? { ...theme, ...night } : theme
  return {
    "--t-primary": tokens.primary,
    "--t-secondary": tokens.secondary,
    "--t-background": tokens.background,
    "--t-surface": tokens.surface,
    "--t-text": tokens.text,
    "--t-muted": tokens.muted,
    "--t-accent": tokens.accent,
    "--t-border": tokens.border,
    "--t-radius": tokens.radius,
    "--t-font-heading": tokens.fontHeading,
    "--t-font-body": tokens.fontBody,
    "--t-font-accent": tokens.fontAccent,
    "--t-section-spacing": tokens.sectionSpacing,
    "--t-content-width": tokens.contentWidth,
    "--t-hero-height": tokens.heroHeight,
    "--t-anim-duration": tokens.animationDuration,
    "--t-anim-easing": tokens.animationEasing,
    "--t-reveal-distance": tokens.revealDistance,
    backgroundColor: tokens.background,
    color: tokens.text,
    fontFamily: tokens.fontBody,
  } as CSSProperties
}

// Build a Google Fonts URL from the families actually used by a theme's
// typography tokens (§2.2.7 — invitations must not load every available font).
export function themeFontsHref(theme: ThemeTokens): string {
  const families = [theme.fontHeading, theme.fontBody, theme.fontAccent]
    .map((f) => {
      // Extract the first quoted family name, e.g. "'Playfair Display', serif"
      const quoted = f.match(/'([^']+)'/)
      if (quoted) return quoted[1]
      return f.split(",")[0].trim()
    })
    .filter((f) => f.length > 0)
    .filter((f, i, arr) => arr.indexOf(f) === i)

  return (
    "https://fonts.googleapis.com/css2?" +
    families.map((f) => `${f.replace(/ /g, "+")}:wght@400;500;600;700&display=swap`).join("&family=")
  )
}
