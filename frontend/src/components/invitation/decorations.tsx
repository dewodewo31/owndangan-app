import type { ThemeTokens, DecorationStyle } from "@/templates/types"

// Lightweight SVG/CSS decorative motifs — no external assets.
// Each returns a centered separator/divider using the template's accent color.

export function Divider({ style, theme }: { style: DecorationStyle; theme: ThemeTokens }) {
  if (style === "none") return null
  return (
    <div className="flex items-center justify-center py-2" aria-hidden>
      <DividerInner style={style} color={theme.accent} />
    </div>
  )
}

function DividerInner({ style, color }: { style: DecorationStyle; color: string }) {
  switch (style) {
    case "floral":
    case "garden":
      return (
        <svg width="120" height="20" viewBox="0 0 120 20" fill="none">
          <path d="M10 10 H45" stroke={color} strokeWidth="1" />
          <path d="M75 10 H110" stroke={color} strokeWidth="1" />
          <path d="M60 4 C55 10 55 10 60 16 C65 10 65 10 60 4Z" fill={color} opacity="0.8" />
          <circle cx="60" cy="10" r="2" fill={color} />
        </svg>
      )
    case "batik":
      return (
        <svg width="120" height="20" viewBox="0 0 120 20" fill="none">
          <path d="M10 10 H48" stroke={color} strokeWidth="1" />
          <path d="M72 10 H110" stroke={color} strokeWidth="1" />
          <g fill="none" stroke={color} strokeWidth="1">
            <circle cx="60" cy="10" r="5" />
            <circle cx="60" cy="10" r="2" fill={color} />
          </g>
        </svg>
      )
    case "botanical":
      return (
        <svg width="120" height="20" viewBox="0 0 120 20" fill="none">
          <path d="M10 10 H50" stroke={color} strokeWidth="1" />
          <path d="M70 10 H110" stroke={color} strokeWidth="1" />
          <path d="M60 2 C56 8 56 12 60 18 C64 12 64 8 60 2Z" stroke={color} strokeWidth="1" fill="none" />
        </svg>
      )
    case "geometric":
      return (
        <svg width="120" height="20" viewBox="0 0 120 20" fill="none">
          <path d="M10 10 H50" stroke={color} strokeWidth="1" />
          <path d="M70 10 H110" stroke={color} strokeWidth="1" />
          <rect x="55" y="5" width="10" height="10" transform="rotate(45 60 10)" stroke={color} strokeWidth="1" fill="none" />
        </svg>
      )
    case "gold":
      return (
        <svg width="140" height="22" viewBox="0 0 140 22" fill="none">
          <path d="M5 11 H58" stroke={color} strokeWidth="1.5" />
          <path d="M82 11 H135" stroke={color} strokeWidth="1.5" />
          <path d="M70 3 L78 11 L70 19 L62 11 Z" fill={color} />
        </svg>
      )
    case "zen":
      return (
        <svg width="80" height="24" viewBox="0 0 80 24" fill="none">
          <circle cx="40" cy="12" r="9" stroke={color} strokeWidth="1.5" fill="none" />
          <path d="M40 3 A9 9 0 0 1 49 12" stroke={color} strokeWidth="1.5" fill="none" />
        </svg>
      )
    case "editorial":
      return <div style={{ width: "60px", height: "2px", background: color }} />
    default:
      return null
  }
}

// A batik-inspired corner ornament for traditional templates.
export function BatikCorner({ color }: { color: string }) {
  return (
    <svg width="60" height="60" viewBox="0 0 60 60" fill="none" aria-hidden>
      <path d="M2 2 H30 M2 2 V30" stroke={color} strokeWidth="1.5" />
      <g stroke={color} strokeWidth="1" fill="none">
        <circle cx="14" cy="14" r="5" />
        <path d="M14 2 C14 8 20 14 26 14" />
      </g>
    </svg>
  )
}
