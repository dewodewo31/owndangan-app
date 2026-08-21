// Shared content schema consumed by every template.
// All templates render the SAME data — only presentation differs.

export type SectionKey =
  | "cover"
  | "quote"
  | "couple"
  | "parents"
  | "countdown"
  | "events"
  | "add-to-calendar"
  | "love-story"
  | "video"
  | "gallery"
  | "location"
  | "rsvp"
  | "gift"
  | "guestbook"
  | "closing"

export interface EventBlock {
  label: string
  venue?: string
  address?: string
  map_url?: string
  time?: string
  end_time?: string
  date?: string
  description?: string
  icon?: string
}

export interface CoupleProfile {
  name?: string
  nickname?: string
  description?: string
  childOrder?: string
  parents?: string
  instagram?: string
  photo?: string
}

export interface InvitationModel {
  slug: string
  eventId?: string
  names: { groom?: string; bride?: string; full: string }
  parents: { groom?: string; bride?: string }
  date?: string // YYYY-MM-DD
  time?: string
  events: { akad?: EventBlock; resepsi?: EventBlock; extra?: EventBlock[] }
  couple?: { groom?: CoupleProfile; bride?: CoupleProfile }
  opening?: string
  closing?: string
  verse?: { enabled: boolean; religion?: string; text?: string; source?: string }
  dressCode?: string
  gallery: { image_url: string; caption?: string }[]
  loveStories: { title: string; story: string; year?: string; date?: string; image_url?: string }[]
  video?: string | null
  guestbook: { name: string; message: string; created_at?: string }[]
  gift?: {
    bank_accounts?: Array<Record<string, unknown>>
    ewallet?: Record<string, unknown>
    qris_image_url?: string
    gift_message?: string
  }
  music?: { title: string; file_url?: string; preset?: string } | null
  sections: Record<string, boolean>
  token?: string | null
  guestName?: string
}

export interface ThemeTokens {
  primary: string
  secondary: string
  background: string
  surface: string
  text: string
  muted: string
  accent: string
  border: string
  radius: string
  fontHeading: string
  fontBody: string
  fontAccent: string
  sectionSpacing: string
  contentWidth: string
  heroHeight: string
  animationDuration: string
  animationEasing: string
  revealDistance: string
}

export type NavStyle =
  | "bottom-floating"
  | "floating-menu"
  | "side"
  | "overlay"
  | "decorative-bottom"

export type RevealVariant = "fade" | "fade-up" | "scale" | "image-reveal" | "text-reveal"

export interface SectionSpec {
  key: SectionKey
  variant?: string
  [k: string]: unknown
}

export type DecorationStyle =
  | "none"
  | "floral"
  | "batik"
  | "botanical"
  | "geometric"
  | "gold"
  | "zen"
  | "garden"
  | "editorial"

export type TimelineStyle =
  | "classic"
  | "editorial"
  | "minimal"
  | "romantic"
  | "luxury"
  | "botanical"
  | "modern"
  | "zen"

export interface TemplateDefinition {
  kind: string
  name: string
  category: string
  occasions: string[]
  description: string
  tags: string[]
  thumbnail: string
  theme: ThemeTokens
  night?: Partial<ThemeTokens>
  nav: NavStyle
  animation: { variant: RevealVariant; stagger: boolean; parallax: boolean }
  sections: SectionSpec[]
  decoration: DecorationStyle
  timelineStyle: TimelineStyle
  musicPlacement?: "top-left" | "bottom-left"
}
