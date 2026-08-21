import { TEMPLATES } from "../index"
import { themeVars } from "../theme"

const COLOR_KEYS = [
  "primary",
  "secondary",
  "background",
  "surface",
  "text",
  "muted",
  "accent",
  "border",
] as const

describe("night theme tokens", () => {
  it("every template defines a night palette with all color keys", () => {
    expect(TEMPLATES.length).toBe(11)
    for (const def of TEMPLATES) {
      expect(def.night).toBeDefined()
      for (const key of COLOR_KEYS) {
        expect(def.night?.[key]).toBeDefined()
      }
    }
  })

  it("themeVars(def.theme, def.night) emits the night values for the 8 color vars", () => {
    for (const def of TEMPLATES) {
      const vars = themeVars(def.theme, def.night)
      for (const key of COLOR_KEYS) {
        expect(vars[`--t-${key}`]).toBe(def.night![key])
      }
    }
  })

  it("themeVars(def.theme) matches themeVars(def.theme, false) and keeps day values", () => {
    for (const def of TEMPLATES) {
      const day = themeVars(def.theme)
      const explicit = themeVars(def.theme, false)
      expect(day).toEqual(explicit)
      for (const key of COLOR_KEYS) {
        expect(day[`--t-${key}`]).toBe(def.theme[key])
      }
    }
  })
})