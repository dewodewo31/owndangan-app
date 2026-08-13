"use client"

import { themeVars, GOOGLE_FONTS_HREF } from "@/templates/theme"
import { Reveal } from "@/templates/anim"
import { Nav } from "./nav"
import { Divider } from "./decorations"
import { MusicBar } from "./MusicBar"
import * as S from "./sections"
import type { InvitationModel, SectionSpec, TemplateDefinition } from "@/templates/types"

const SECTION_MAP: Record<string, (p: S.SectionProps) => JSX.Element | null> = {
  cover: S.Cover,
  quote: S.Quote,
  couple: S.Couple,
  parents: S.Parents,
  countdown: S.Countdown,
  events: S.Events,
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
      <link rel="stylesheet" href={GOOGLE_FONTS_HREF} />

      <Nav style={nav} sections={sections} />

      {model.music?.file_url && model.sections && (
        <MusicBar music={model.music} theme={theme} />
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
    </div>
  )
}
