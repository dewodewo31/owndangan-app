import { TEMPLATES, OCCASION_LABELS, selectTemplate } from "../index"
import type { TemplateDefinition } from "../types"

describe("occasion metadata", () => {
  it("has 11 templates (10 original + corporate)", () => {
    expect(TEMPLATES.length).toBe(11)
  })

  it("every template declares occasions that exist in OCCASION_LABELS", () => {
    for (const t of TEMPLATES) {
      expect(t.occasions.length).toBeGreaterThan(0)
      for (const o of t.occasions) {
        expect(OCCASION_LABELS[o]).toBeTruthy()
      }
    }
  })

  it("corporate template has exactly [corporate, event] occasions", () => {
    const corporate = TEMPLATES.find((t) => t.kind === "corporate")
    expect(corporate).toBeDefined()
    expect(corporate!.occasions).toEqual(["corporate", "event"])
  })
})

describe("selectTemplate", () => {
  it("falls back to a valid definition for unknown names", () => {
    const t = selectTemplate(undefined, "No Such Template")
    expect(t).toBeDefined()
    expect(t.kind).toBeTruthy()
  })

  it("premium group falls back to romantic-elegant", () => {
    const t = selectTemplate("premium")
    expect(t).toBeDefined()
    expect(t.kind).toBe("romantic-elegant")
  })
})
