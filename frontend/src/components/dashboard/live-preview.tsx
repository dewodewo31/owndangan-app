"use client"

import type { CSSProperties } from "react"
import { TemplateShell } from "@/components/invitation/shell"
import { CoverGate } from "@/components/invitation/cover-gate"
import { selectTemplate } from "@/templates"
import { themeVars } from "@/templates/theme"
import { useNightMode } from "@/hooks/use-night-mode"
import { buildPreviewModel } from "@/lib/editor-preview"
import type {
  DigitalGift,
  EventSections,
  GalleryPhoto,
  LoveStory,
  Music,
  TemplateSummary,
  WeddingEvent,
} from "@/lib/types"

const FALLBACK_HERO_IMAGE =
  "https://images.unsplash.com/photo-1519046904744-6fd9aeda9e0e?auto=format&fit=crop&w=1600&q=60"

type LivePreviewProps = {
  event: WeddingEvent | null | undefined
  sections: EventSections | null | undefined
  gallery: GalleryPhoto[] | null | undefined
  music: Music | null | undefined
  gift: DigitalGift | null | undefined
  loveStories?: LoveStory[] | null | undefined
  template: TemplateSummary | null | undefined
}

/**
 * Renders the REAL public invitation renderer (same composition as
 * `[slug]/page.tsx`) inside a virtual phone frame, driven purely by editor
 * state. No API calls — `buildPreviewModel` is a pure adapter over the props.
 */
export function LivePreview({
  event,
  sections,
  gallery,
  music,
  gift,
  loveStories,
  template,
}: LivePreviewProps) {
  const model = buildPreviewModel(
    event,
    sections,
    gallery,
    music,
    gift,
    loveStories,
    template
  )
  const { night } = useNightMode()

  if (!event || !template || !model) {
    return (
      <div className="flex min-h-[540px] w-full items-center justify-center">
        <div className="mx-auto w-full max-w-[375px]">
          <p className="mb-3 text-xs font-medium uppercase tracking-wider text-muted-foreground">
            Pratinjau Undangan
          </p>
          <div className="rounded-3xl border border-dashed border-input px-10 py-16 text-center text-sm text-muted-foreground shadow-sm">
            Pratinjau akan muncul setelah undangan dibuat
          </div>
        </div>
      </div>
    )
  }

  const definition = selectTemplate(template.group_name, template.name)
  const css = (template.css_config ?? {}) as Record<string, unknown>
  const primaryColor = (css.primary_color as string) || "#b22234"
  const heroImage =
    (css.hero_image as string) || model.gallery[0]?.image_url || FALLBACK_HERO_IMAGE
  const vars = themeVars(definition.theme, night ? definition.night : undefined) as CSSProperties

  return (
    <div className="flex w-full justify-center">
      <div className="w-full max-w-[375px]">
        <p className="mb-3 text-xs font-medium uppercase tracking-wider text-muted-foreground">
          Pratinjau Undangan
        </p>

        {/* Subtle backdrop behind the phone */}
        <div className="rounded-3xl bg-gradient-to-b from-neutral-200 to-neutral-300 p-3 shadow-inner">
          {/* Phone frame — `transform` makes it the containing block for
              CoverGate's `position:fixed`, so the gate covers only the phone. */}
          <div
            className="relative overflow-hidden rounded-[2.5rem] border border-neutral-400/60 bg-neutral-900 shadow-2xl"
            style={{ width: "100%", height: "620px", transform: "translateZ(0)" }}
          >
            {/* Screen */}
            <div
              className="relative h-full w-full overflow-y-auto bg-white"
              style={vars}
            >
              <TemplateShell definition={definition} model={model} />
            </div>

            {/* Notch / status bar */}
            <div className="pointer-events-none absolute inset-x-0 top-0 z-[70] flex justify-center pt-2.5">
              <div className="h-6 w-28 rounded-full bg-neutral-900 shadow" aria-hidden="true" />
            </div>

            {/* CoverGate — same composition as the published page, confined
                inside the transformed phone frame. */}
            <CoverGate
              model={model}
              heroImage={heroImage}
              primaryColor={primaryColor}
              lockBodyScroll={false}
            />
          </div>
        </div>
      </div>
    </div>
  )
}