import { renderToStaticMarkup } from "react-dom/server"
import { LivePreview } from "../live-preview"
import type {
  DigitalGift,
  EventSections,
  GalleryPhoto,
  LoveStory,
  Music,
  TemplateSummary,
  WeddingEvent,
} from "@/lib/types"

const event: WeddingEvent = {
  id: "evt_001",
  title: "Undangan Pernikahan Rina & Pika",
  slug: "rina-pika",
  couple_name: "Rina & Pika",
  groom_name: "Pika",
  bride_name: "Rina",
  wedding_date: "2026-12-24",
  wedding_time: "15:00",
  ceremony_venue: "Masjid Al-Nur",
  ceremony_address: "Jl. Imam Reza 12, Jakarta",
  reception_venue: "Hotel Astoria Ballroom",
  reception_address: "Blvd. Oscar 45, Jakarta",
  status: "draft",
  published_at: null,
}

const sections: EventSections = {
  id: "sec_001",
  event_id: "evt_001",
  hero_enabled: true,
  couple_enabled: true,
  gallery_enabled: true,
  event_details_enabled: true,
  video_enabled: true,
  rsvp_enabled: true,
  guestbook_enabled: true,
  love_story_enabled: true,
  digital_gifts_enabled: true,
  verse_enabled: true,
  verse_religion: "islam",
  verse_text: "Wa min ayatihi an khalaqakum min nafs wahidah",
  verse_source: "Ar-Rum 21",
}

const music: Music = {
  id: "mus_001",
  event_id: "evt_001",
  title: "Bismillah Anthem",
  file_url: "https://cdn.example.com/music.mp3",
  preset: "anthem",
  is_preset: true,
}

const gallery: GalleryPhoto[] = [
  { id: "gal_001", image_url: "https://cdn.example.com/photo1.jpg", sort_order: 0 },
  { id: "gal_002", image_url: "https://cdn.example.com/photo2.jpg", caption: "Reception", sort_order: 1 },
]

const loveStories: LoveStory[] = [
  {
    id: "ls_001",
    title: "The Story",
    story: "Met at university.",
    date: "2019-03-14",
    image_url: "https://cdn.example.com/story1.jpg",
    sort_order: 0,
  },
]

const gift: DigitalGift = {
  id: "gft_001",
  event_id: "evt_001",
  gift_message: "Thank you for your generosity",
  bank_accounts: [
    {
      bank_name: "BCA",
      account_holder: "Rina & Pika",
      account_number: "1234-5678",
    },
  ],
}

const modernTemplate: TemplateSummary = {
  id: "tpl_a",
  name: "Modern Minimalist",
  group_name: "standard",
  css_config: {},
  layout_config: {},
}

const luxuryTemplate: TemplateSummary = {
  id: "tpl_b",
  name: "Luxury Black & Gold",
  group_name: "premium",
  css_config: {},
  layout_config: {},
}

const filled = {
  event,
  sections,
  gallery,
  music,
  gift,
  loveStories,
  template: modernTemplate,
}

describe("LivePreview — filled event", () => {
  const html = renderToStaticMarkup(<LivePreview {...filled} />)

it("renders the real template with couple names present", () => {
    expect(html).toContain("Rina")
    expect(html).toContain("Pika")
    expect(html).toContain("Pratinjau Undangan")
    expect(html).toContain("border-neutral-400")
    expect(html).not.toContain("border-dashed")
  })

  it("applies template theme vars to the wrapper", () => {
    expect(html).toContain("--t-background:#ffffff")
    expect(html).toContain("--t-primary:#1f2937")
  })
})

describe("LivePreview — switching template changes themeVars", () => {
  const a = renderToStaticMarkup(<LivePreview {...filled} />)
  const b = renderToStaticMarkup(<LivePreview {...filled} template={luxuryTemplate} />)

  it("emits different --t-background between the two templates", () => {
    expect(a).toContain("--t-background:#ffffff")
    expect(b).toContain("--t-background:#0e0e10")
  })

  it("emits different --t-primary between the two templates", () => {
    expect(a).toContain("--t-primary:#1f2937")
    expect(b).toContain("--t-primary:#d4af37")
  })
})

describe("LivePreview — null event / template", () => {
it("shows placeholder instead of crashing when event is null", () => {
    const html = renderToStaticMarkup(<LivePreview event={null} template={modernTemplate} />)
    expect(html).toContain("Pratinjau Undangan")
    expect(html).toContain("border-dashed")
    expect(html).not.toContain("border-neutral-400")
    expect(html).not.toContain("Rina")
  })

  it("shows placeholder when event set but template missing", () => {
    const html = renderToStaticMarkup(<LivePreview {...filled} template={null} />)
    expect(html).toContain("Pratinjau Undangan")
    expect(html).toContain("border-dashed")
    expect(html).not.toContain("border-neutral-400")
  })
})