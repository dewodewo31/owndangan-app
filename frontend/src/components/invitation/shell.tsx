"use client"

import { themeVars, themeFontsHref } from "@/templates/theme"
import { Reveal } from "@/templates/anim"
import { Nav } from "./nav"
import { Divider } from "./decorations"
import { MusicBar } from "./MusicBar"
import { ShareButton } from "./share"
import * as S from "./sections"
import type { InvitationModel, SectionSpec, TemplateDefinition } from "@/templates/types"

const SECTION_MAP: Record<string, (p: S.SectionProps) => JSX.Element | null> = {
  cover: S.Cover,
  quote: S.Quote,
  couple: S.Couple,
  parents: S.Parents,
  countdown: S.Countdown,
  events: S.Events,
  "add-to-calendar": S.AddToCalendar,
  "love-story": S.LoveStory,
  video: S.Video,
  gallery: S.Gallery,
  location: S.Location,
  rsvp: S.RSVP,
  gift: S.Gift,
  guestbook: S.Guestbook,
  closing: S.Closing,
}

export function TemplateShell({
  definition,
  model,
}: {
  definition: TemplateDefinition
  model: InvitationModel
}) {
  const { theme, nav, animation, sections, decoration } = definition

  return (
    <div
      style={themeVars(theme)}
      className="invitation-root relative min-h-screen"
    >
      <link rel="stylesheet" href={themeFontsHref(theme)} />

      <Nav style={nav} sections={sections} />

      {model.music?.file_url && model.sections && (
        <MusicBar music={model.music} theme={theme} placement={definition.musicPlacement} />
      )}

      <main>
        {sections.map((spec: SectionSpec, i: number) => {
          const Comp = SECTION_MAP[spec.key]
          if (!Comp) return null
          return (
            <div key={spec.key + i} id={spec.key} style={{ scrollMarginTop: "64px" }}>
              {i > 0 && spec.key !== "cover" && decoration !== "none" && (
                <Divider style={decoration} theme={theme} />
              )}
              <Reveal variant={animation.variant} as="div">
                <Comp model={model} theme={theme} spec={spec} />
              </Reveal>
            </div>
          )
        })}
      </main>

      <div className="px-6 pb-12 pt-4">
        <ShareButton slug={model.slug} guestName={model.guestName} theme={theme} eventId={model.eventId} />
      </div>
    </div>
  )
}
