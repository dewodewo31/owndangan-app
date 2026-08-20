import type { InvitationModel } from "@/templates/types"
import { buildInvitationModel, type PublicEventResponse } from "./invitation"
import type {
  DigitalGift,
  EventSections,
  GalleryPhoto,
  LoveStory,
  Music,
  TemplateSummary,
  WeddingEvent,
} from "./types"

/**
 * Pure adapter: editor state -> PublicEventResponse (see invitation.ts:5-59).
 * Read-only over its inputs; never mutates them. Returns null for a nullish
 * event so a not-yet-loaded editor cannot crash the preview.
 */
export function editorStateToPublicEvent(
  event: WeddingEvent | null | undefined,
  sections: EventSections | null | undefined,
  gallery: GalleryPhoto[] | null | undefined,
  music: Music | null | undefined,
  gift: DigitalGift | null | undefined,
  loveStories: LoveStory[] | null | undefined,
  template: TemplateSummary | null | undefined
): PublicEventResponse | null {
  if (!event) return null

  return {
    event: {
      id: event.id,
      title: event.title,
      couple_name: event.couple_name,
      groom_name: event.groom_name,
      bride_name: event.bride_name,
      groom_parents: event.groom_parents,
      bride_parents: event.bride_parents,
      wedding_date: event.wedding_date ?? undefined,
      wedding_time: event.wedding_time,
      ceremony_venue: event.ceremony_venue,
      ceremony_address: event.ceremony_address,
      ceremony_map_url: event.ceremony_map_url,
      reception_venue: event.reception_venue,
      reception_address: event.reception_address,
      reception_map_url: event.reception_map_url,
      video_url: event.video_url,
      view_count: event.view_count,
    },
    template: template
      ? {
          name: template.name,
          group_name: template.group_name,
          css_config: template.css_config,
          layout_config: template.layout_config,
        }
      : null,
    sections: sections
      ? {
          hero_enabled: sections.hero_enabled,
          couple_enabled: sections.couple_enabled,
          event_details_enabled: sections.event_details_enabled,
          gallery_enabled: sections.gallery_enabled,
          video_enabled: sections.video_enabled,
          rsvp_enabled: sections.rsvp_enabled,
          guestbook_enabled: sections.guestbook_enabled,
          love_story_enabled: sections.love_story_enabled ?? false,
          digital_gifts_enabled: sections.digital_gifts_enabled,
          dress_code: sections.dress_code,
          opening_message: sections.opening_message,
          closing_message: sections.closing_message,
          verse_enabled: sections.verse_enabled,
          verse_religion: sections.verse_religion,
          verse_text: sections.verse_text,
          verse_source: sections.verse_source,
          music: music
            ? {
                title: music.title,
                file_url: music.file_url,
                preset: music.preset,
                is_preset: music.is_preset,
              }
            : null,
        }
      : undefined,
    gallery: (gallery ?? []).map((g) => ({
      image_url: g.image_url,
      caption: g.caption,
      sort_order: g.sort_order,
    })),
    love_stories: (loveStories ?? []).map((st) => ({
      id: st.id,
      title: st.title,
      story: st.story,
      year: st.year,
      date: st.date,
      image_url: st.image_url,
      sort_order: st.sort_order,
    })),
    guestbook: [],
    digital_gift: gift
      ? {
          bank_accounts: gift.bank_accounts,
          ewallet: gift.ewallet,
          qris_image_url: gift.qris_image_url,
          gift_message: gift.gift_message,
        }
      : null,
  }
}

/** Render-ready invitation model from live editor state. Pure, no I/O. */
export function buildPreviewModel(
  event: WeddingEvent | null | undefined,
  sections: EventSections | null | undefined,
  gallery: GalleryPhoto[] | null | undefined,
  music: Music | null | undefined,
  gift: DigitalGift | null | undefined,
  loveStories: LoveStory[] | null | undefined,
  template: TemplateSummary | null | undefined
): InvitationModel | null {
  if (!event) return null
  const data = editorStateToPublicEvent(
    event,
    sections,
    gallery,
    music,
    gift,
    loveStories,
    template
  )
  if (!data) return null
  return buildInvitationModel(data, new URLSearchParams(), event.slug)
}
