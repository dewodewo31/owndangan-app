import { TEMPLATES } from "@/templates"
import type { ThemeTokens } from "@/templates/types"

function hexToRgb(hex: string): [number, number, number] {
  let h = hex.trim().replace("#", "")
  if (h.length === 3) h = h.split("").map((c) => c + c).join("")
  const n = parseInt(h.slice(0, 6), 16)
  return [(n >> 16) & 255, (n >> 8) & 255, n & 255]
}

function luminance(hex: string): number {
  const [r, g, b] = hexToRgb(hex).map((c) => {
    const s = c / 255
    return s <= 0.03928 ? s / 12.92 : Math.pow((s + 0.055) / 1.055, 2.4)
  })
  return 0.2126 * r + 0.7152 * g + 0.0722 * b
}

export function contrast(a: string, b: string): number {
  const l1 = luminance(a)
  const l2 = luminance(b)
  return (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05)
}

function check(pairs: Array<[string, string, number]>, label: string, name: string, mode: string, failures: string[]) {
  for (const [bg, fg, min] of pairs) {
    if (!bg || !fg) continue
    const ratio = contrast(bg, fg)
    if (ratio < min) {
      failures.push(`${name} [${mode}] ${fg} on ${bg}: ${ratio.toFixed(2)} < ${min} (${label})`)
    }
  }
}

describe("template palette contrast (WCAG)", () => {
  it("keeps text readable on background/surface in day AND night", () => {
    const failures: string[] = []
    for (const t of TEMPLATES) {
      const modes: Array<[string, ThemeTokens]> = [["day", t.theme]]
      if (t.night) modes.push(["night", { ...t.theme, ...t.night }])
      for (const [mode, tok] of modes) {
        check(
          [
            [tok.background, tok.text, 4.5],
            [tok.surface, tok.text, 4.5],
            [tok.background, tok.muted, 3.5],
          ],
          "body-text", t.name, mode, failures
        )
      }
    }
    if (failures.length) console.log("CONTRAST FAILURES:\n" + failures.join("\n"))
    expect(failures).toEqual([])
  })
})
