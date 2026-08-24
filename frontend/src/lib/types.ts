export type AuthResponse = {
  id: string
  name: string
  email: string
  phone?: string
  role: string
  status: string
  created_at: string
  access_token: string
  refresh_token: string
  expires_in: number
}

export type PackageBrief = {
  id: string
  name: string
  code: string
  price: number
  guest_limit?: number
  template_group: string
  features?: Record<string, unknown>
}

export type Package = {
  id: string
  name: string
  code: string
  price: number
  duration_days?: number
  guest_limit?: number
  template_group: string
  features?: Record<string, unknown>
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export type Subscription = {
  id: string
  package: PackageBrief
  status: string
  start_at: string
  expires_at?: string
}

export type ApiResponse<T> = {
  success: boolean
  data: T
  meta: { request_id?: string; pagination?: unknown }
}

export type TokenResponse = {
  access_token: string
  refresh_token: string
  expires_in: number
}

export type WeddingEvent = {
  id: string
  user_id?: string
  template_id?: string | null
  title: string
  slug: string
  couple_name?: string
  groom_name?: string
  bride_name?: string
  groom_parents?: string
  bride_parents?: string
  wedding_date?: string | null
  wedding_time?: string
  ceremony_venue?: string
  ceremony_address?: string
  ceremony_map_url?: string
  reception_venue?: string
  reception_address?: string
  reception_map_url?: string
  music_url?: string
  video_url?: string
  status: string
  published_at?: string | null
  view_count?: number
  guest_count?: number
  rsvp_count?: number
  created_at?: string
  updated_at?: string
}

export type EventSections = {
  id: string
  event_id: string
  hero_enabled: boolean
  couple_enabled: boolean
  event_details_enabled: boolean
  gallery_enabled: boolean
  video_enabled: boolean
  music_id?: string | null
  rsvp_enabled: boolean
  guestbook_enabled: boolean
  love_story_enabled?: boolean
  digital_gifts_enabled: boolean
  dress_code?: string
  closing_message?: string
  opening_message?: string
  verse_enabled?: boolean
  verse_religion?: string
  verse_text?: string
  verse_source?: string
}

export type GalleryPhoto = {
  id: string
  image_url: string
  caption?: string
  sort_order: number
}

export type Music = {
  id: string
  event_id?: string | null
  title: string
  file_url?: string
  preset?: string
  is_preset: boolean
}

export type LoveStory = {
  id: string
  title: string
  story: string
  year?: string
  date?: string
  image_url?: string
  sort_order: number
  created_at?: string
  updated_at?: string
}

export type DigitalGift = {
  id: string
  event_id: string
  bank_accounts?: Array<Record<string, unknown>>
  ewallet?: Record<string, unknown>
  qris_image_url?: string
  gift_message?: string
}

export type TemplateSummary = {
  id: string
  name: string
  group_name: string
  thumbnail_url?: string
  css_config?: Record<string, unknown>
  layout_config?: Record<string, unknown>
}

export type PublicEvent = {
  title: string
  couple_name?: string
  groom_name?: string
  bride_name?: string
  groom_parents?: string
  bride_parents?: string
  wedding_date?: string
  wedding_time?: string
  ceremony_venue?: string
  ceremony_address?: string
  ceremony_map_url?: string
  reception_venue?: string
  reception_address?: string
  reception_map_url?: string
  video_url?: string
  view_count?: number
}

export type PublicEventResponse = {
  event: PublicEvent
  sections?: EventSections
  gallery?: GalleryPhoto[]
  guestbook?: GuestbookMessage[]
  love_stories?: LoveStory[]
  digital_gift?: DigitalGift
}

export type GuestbookMessage = {
  name: string
  message: string
  created_at: string
}
