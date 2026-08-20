import { editorStateToPublicEvent, buildPreviewModel } from "../editor-preview"
import type {
  DigitalGift,
  EventSections,
  GalleryPhoto,
  LoveStory,
  Music,
  TemplateSummary,
  WeddingEvent,
} from "../types"

/** Realistic full wedding fixture (ceremony + reception, all sections on). */
const event: WeddingEvent = {
  id: "evt_001",
  title: "Undangan Pernikahan Rina & Pika",
  slug: "rina-pika",
  couple_name: "Rina & Pika",
  groom_name: "Pika",
  bride_name: "Rina",
  groom_parents: "Mr. & Mrs. Bayu",
  bride_parents: "Mr. & Mrs. Rahat",
  wedding_date: "2026-12-24",
  wedding_time: "15:00",
  ceremony_venue: "Masjid Al-Nur",
  ceremony_address: "Jl. Imam Reza 12, Jakarta",
  ceremony_map_url: "https://maps.example.com/ceremony",
  reception_venue: "Hotel Astoria Ballroom",
  reception_address: "Blvd. Oscar 45, Jakarta",
  reception_map_url: "https://maps.example.com/reception",
  video_url: "https://cdn.example.com/video.mp4",
  status: "draft",
  published_at: null,
  view_count: 42,
}

const sections: EventSections = {
  id: "sec_001",
  event_id: "evt_001",
  hero_enabled: true,
  couple_enabled: true,
  event_details_enabled: true,
  gallery_enabled: true,
  video_enabled: true,
  rsvp_enabled: true,
  guestbook_enabled: true,
  love_story_enabled: true,
  digital_gifts_enabled: true,
  dress_code: "Formal",
  opening_message: "Bismillah",
  closing_message: "Salam",
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
  {
    id: "gal_001",
    image_url: "https://cdn.example.com/photo1.jpg",
    caption: "First look",
    sort_order: 0,
  },
  {
    id: "gal_002",
    image_url: "https://cdn.example.com/photo2.jpg",
    caption: "Reception dance",
    sort_order: 1,
  },
]

const loveStories: LoveStory[] = [
  {
    id: "ls_001",
    title: "How We Met",
    story: "Met in university library.",
    year: "2019",
    date: "2019-03-14",
    image_url: "https://cdn.example.com/story1.jpg",
    sort_order: 0,
  },
  {
    id: "ls_002",
    title: "The Proposal",
    story: "Proposed on the beach.",
    year: "2025",
    date: "2025-06-01",
    image_url: "https://cdn.example.com/story2.jpg",
    sort_order: 1,
  },
]

const gift: DigitalGift = {
  id: "gft_001",
  event_id: "evt_001",
  bank_accounts: [
    {
      bank_name: "BCA",
      account_holder: "Rina & Pika",
      account_number: "1234-5678",
    },
  ],
  ewallet: { provider: "PayPal", handle: "rina@example.com" },
  qris_image_url: "https://cdn.example.com/qris.png",
  gift_message: "Thank you for your generosity",
}

const templateSummary: TemplateSummary = {
  id: "tpl_001",
  name: "Modern Minimalist",
  group_name: "standard",
  css_config: { primary_color: "#b22234", hero_image: "https://cdn.example.com/hero.jpg" },
  layout_config: { columns: 2 },
}

const full = editorStateToPublicEvent(event, sections, gallery, music, gift, loveStories, templateSummary)

describe("editorStateToPublicEvent — full event", () => {
it("maps every event field", () => {
    const ev = full!.event
    expect(ev.id).toBe("evt_001")
    expect(ev.title).toBe("Undangan Pernikahan Rina & Pika")
    expect(ev.couple_name).toBe("Rina & Pika")
    expect(ev.groom_name).toBe("Pika")
    expect(ev.bride_name).toBe("Rina")
    expect(ev.groom_parents).toBe("Mr. & Mrs. Bayu")
    expect(ev.bride_parents).toBe("Mr. & Mrs. Rahat")
    expect(ev.wedding_date).toBe("2026-12-24")
    expect(ev.wedding_time).toBe("15:00")
    expect(ev.ceremony_venue).toBe("Masjid Al-Nur")
    expect(ev.ceremony_address).toBe("Jl. Imam Reza 12, Jakarta")
    expect(ev.ceremony_map_url).toBe("https://maps.example.com/ceremony")
    expect(ev.reception_venue).toBe("Hotel Astoria Ballroom")
    expect(ev.reception_address).toBe("Blvd. Oscar 45, Jakarta")
    expect(ev.reception_map_url).toBe("https://maps.example.com/reception")
    expect(ev.video_url).toBe("https://cdn.example.com/video.mp4")
    expect(ev.view_count).toBe(42)
  })

  it("maps the template summary", () => {
    expect(full!.template).toEqual({
      name: "Modern Minimalist",
      group_name: "standard",
      css_config: { primary_color: "#b22234", hero_image: "https://cdn.example.com/hero.jpg" },
      layout_config: { columns: 2 },
    })
  })

  it("maps every section flag and text field", () => {
    const s = full!.sections
    expect(s).not.toBeUndefined()
    expect(s!.hero_enabled).toBe(true)
    expect(s!.couple_enabled).toBe(true)
    expect(s!.event_details_enabled).toBe(true)
    expect(s!.gallery_enabled).toBe(true)
    expect(s!.video_enabled).toBe(true)
    expect(s!.rsvp_enabled).toBe(true)
    expect(s!.guestbook_enabled).toBe(true)
    expect(s!.love_story_enabled).toBe(true)
    expect(s!.digital_gifts_enabled).toBe(true)
    expect(s!.dress_code).toBe("Formal")
    expect(s!.opening_message).toBe("Bismillah")
    expect(s!.closing_message).toBe("Salam")
    expect(s!.verse_enabled).toBe(true)
    expect(s!.verse_religion).toBe("islam")
    expect(s!.verse_text).toBe("Wa min ayatihi an khalaqakum min nafs wahidah")
    expect(s!.verse_source).toBe("Ar-Rum 21")
  })

  it("maps nested music inside sections", () => {
    expect(full!.sections!.music).toEqual({
      title: "Bismillah Anthem",
      file_url: "https://cdn.example.com/music.mp3",
      preset: "anthem",
      is_preset: true,
    })
  })

  it("maps gallery preserving order and caption", () => {
    expect(full!.gallery).toHaveLength(2)
    expect(full!.gallery![0]).toEqual({
      image_url: "https://cdn.example.com/photo1.jpg",
      caption: "First look",
      sort_order: 0,
    })
    expect(full!.gallery![1].sort_order).toBe(1)
  })

  it("maps every love story field", () => {
    expect(full!.love_stories).toHaveLength(2)
    expect(full!.love_stories![0]).toEqual({
      id: "ls_001",
      title: "How We Met",
      story: "Met in university library.",
      year: "2019",
      date: "2019-03-14",
      image_url: "https://cdn.example.com/story1.jpg",
      sort_order: 0,
    })
    expect(full!.love_stories![1].date).toBe("2025-06-01")
  })

  it("maps digital gift fields", () => {
    expect(full!.digital_gift).toEqual({
      bank_accounts: [
        {
          bank_name: "BCA",
          account_holder: "Rina & Pika",
          account_number: "1234-5678",
        },
      ],
      ewallet: { provider: "PayPal", handle: "rina@example.com" },
      qris_image_url: "https://cdn.example.com/qris.png",
      gift_message: "Thank you for your generosity",
    })
  })
})

describe("editorStateToPublicEvent — partial event", () => {
  const partial = editorStateToPublicEvent(
    {
      id: "evt_999",
      title: "Titulo",
      slug: "t",
      couple_name: "Solo nombre",
      status: "draft",
      wedding_date: null,
    },
    {
      id: "sec_999",
      event_id: "evt_999",
      hero_enabled: false,
      couple_enabled: false,
      event_details_enabled: false,
      gallery_enabled: false,
      video_enabled: false,
      rsvp_enabled: false,
      guestbook_enabled: false,
      digital_gifts_enabled: false,
    },
    null,
    null,
    null,
    null,
    null
  )

  it("keeps present fields, drops absent event fields", () => {
    expect(partial!.event.id).toBe("evt_999")
    expect(partial!.event.title).toBe("Titulo")
    expect(partial!.event.couple_name).toBe("Solo nombre")
    expect(partial!.event.groom_name).toBeUndefined()
    expect(partial!.event.wedding_date).toBeUndefined()
  })

  it("defaults missing optional section flags to false and music to null", () => {
    expect(partial!.sections!.love_story_enabled).toBe(false)
    expect(partial!.sections!.music).toBeNull()
  })

  it("coerces null gallery/loveStories to empty arrays", () => {
    expect(partial!.gallery).toEqual([])
    expect(partial!.love_stories).toEqual([])
  })

  it("leaves optional subsidiaries null", () => {
    expect(partial!.digital_gift).toBeNull()
    expect(partial!.template).toBeNull()
  })
})

describe("editorStateToPublicEvent — null event", () => {
  it("returns null without throwing for nullish event", () => {
    expect(editorStateToPublicEvent(null, sections, gallery, music, gift, loveStories, templateSummary)).toBeNull()
    expect(editorStateToPublicEvent(undefined, undefined, undefined, undefined, undefined, undefined, undefined)).toBeNull()
  })
})

describe("buildPreviewModel — real invitation chain", () => {
  it("maps the full event into a render-ready InvitationModel", () => {
    const model = buildPreviewModel(event, sections, gallery, music, gift, loveStories, templateSummary)
    expect(model).not.toBeNull()
    expect(model!.slug).toBe("rina-pika")
    expect(model!.names.full).toBe("Rina & Pika")
    expect(model!.date).toBe("2026-12-24")
    expect(model!.gallery).toHaveLength(2)
    expect(model!.loveStories).toHaveLength(2)
    expect(model!.verse?.enabled).toBe(true)
    expect(model!.events.akad?.venue).toBe("Masjid Al-Nur")
    expect(model!.events.resepsi?.venue).toBe("Hotel Astoria Ballroom")
    expect(model!.video).toBe("https://cdn.example.com/video.mp4")
    expect(model!.gift?.gift_message).toBe("Thank you for your generosity")
  })

  it("returns null for a null event", () => {
    expect(buildPreviewModel(null, null, null, null, null, null, null)).toBeNull()
  })
})