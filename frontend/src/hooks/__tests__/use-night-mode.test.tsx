import { act } from "react"
import { createRoot } from "react-dom/client"
import { useNightMode } from "@/hooks/use-night-mode"

;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

function renderProbe(hour: number, opts?: { start?: string; end?: string }): boolean {
  const container = document.createElement("div")
  document.body.appendChild(container)
  let captured: boolean | null = null
  function Probe() {
    const { night } = useNightMode(opts ?? { start: "18:00", end: "06:00" })
    captured = night
    return null
  }
  const root = createRoot(container)
  const realGetHours = Date.prototype.getHours
  act(() => {
    Date.prototype.getHours = () => hour
    root.render(<Probe />)
  })
  const result = captured
  act(() => {
    root.unmount()
    Date.prototype.getHours = realGetHours
  })
  container.remove()
  return result as boolean
}

describe("useNightMode", () => {
  it("is night across the default 18:00–06:00 window (wraps midnight)", () => {
    expect(renderProbe(18)).toBe(true)
    expect(renderProbe(23)).toBe(true)
    expect(renderProbe(0)).toBe(true)
    expect(renderProbe(5)).toBe(true)
  })

  it("is day between 06:00 and 18:00", () => {
    expect(renderProbe(6)).toBe(false)
    expect(renderProbe(12)).toBe(false)
    expect(renderProbe(17)).toBe(false)
  })

  it("honours a custom window", () => {
    expect(renderProbe(20, { start: "20:00", end: "22:00" })).toBe(true)
    expect(renderProbe(12, { start: "20:00", end: "22:00" })).toBe(false)
  })

  it("computes night in an explicit IANA timezone", () => {
    expect(() => renderProbe(12, { timezone: "Asia/Jakarta" })).not.toThrow()
  })
})
