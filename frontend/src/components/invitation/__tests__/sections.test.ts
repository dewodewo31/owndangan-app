import { buildIcs } from "../sections"

describe("buildIcs", () => {
  const base = {
    label: "Akad",
    coupleName: "A & B",
    date: "2026-12-25",
    time: "09:00:00",
    end: "11:00:00",
    venue: "Gedung",
    address: "Jakarta",
  }

  it("emits TZID and X-WR-TIMEZONE when a zone is provided", () => {
    const ics = buildIcs({ ...base, timeZone: "Asia/Jakarta" })
    expect(ics).toContain("DTSTART;TZID=Asia/Jakarta:20261225T090000")
    expect(ics).toContain("DTEND;TZID=Asia/Jakarta:20261225T110000")
    expect(ics).toContain("X-WR-TIMEZONE:Asia/Jakarta")
  })

  it("keeps floating times when no zone is provided", () => {
    const ics = buildIcs(base)
    expect(ics).toContain("DTSTART:20261225T090000")
    expect(ics).toContain("DTEND:20261225T110000")
    expect(ics).not.toContain("TZID")
    expect(ics).not.toContain("X-WR-TIMEZONE")
  })

  it("skips DTSTART/DTEND gracefully when there is no date", () => {
    const ics = buildIcs({ ...base, date: undefined })
    expect(ics).not.toContain("DTSTART")
    expect(ics).not.toContain("DTEND")
    expect(ics).toContain("BEGIN:VCALENDAR")
    expect(ics).toContain("END:VCALENDAR")
  })
})