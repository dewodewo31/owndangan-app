import { act } from "react"
import { createRoot, type Root } from "react-dom/client"
import AnalyticsWidget from "../analytics-widget"
import api from "@/lib/api"

;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

jest.mock("@/lib/api", () => ({
  __esModule: true,
  default: { get: jest.fn() },
}))

jest.mock("next/link", () => ({
  __esModule: true,
  default: ({ children, ...props }: { children: React.ReactNode; [key: string]: unknown }) => (
    <a {...props}>{children}</a>
  ),
}))

const mockGet = api.get as jest.Mock

const analyticsPayload = {
  views: 120,
  unique_views: 87,
  whatsapp_clicks: 34,
  map_clicks: 21,
  phone_clicks: 9,
  rsvp_count: 15,
}

const UPSELL_TEXT = "Aktifkan paket berbayar untuk melihat analitik undangan"

const roots: Root[] = []

function renderWidget(eventId = "evt_001") {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const root = createRoot(container)
  roots.push(root)
  act(() => {
    root.render(<AnalyticsWidget eventId={eventId} />)
  })
  return container
}

afterEach(() => {
  roots.forEach((root) => act(() => root.unmount()))
  roots.length = 0
  jest.clearAllMocks()
})

describe("AnalyticsWidget", () => {
  it("renders the 6 analytics numbers for an entitled user", async () => {
    mockGet.mockImplementation((url: string) => {
      if (url === "/subscriptions/current") {
        return Promise.resolve({ data: { data: { package: { code: "premium" } } } })
      }
      if (url === "/events/evt_001/analytics") {
        return Promise.resolve({ data: { data: analyticsPayload } })
      }
      return Promise.reject(new Error(`unexpected url: ${url}`))
    })

    const container = renderWidget()
    await act(async () => {})

    expect(container.innerHTML).toContain("120")
    expect(container.innerHTML).toContain("87")
    expect(container.innerHTML).toContain("34")
    expect(container.innerHTML).toContain("21")
    expect(container.innerHTML).toContain("9")
    expect(container.innerHTML).toContain("15")
    expect(container.innerHTML).toContain("Views")
    expect(container.innerHTML).toContain("Pengunjung Unik")
    expect(container.innerHTML).toContain("Klik WhatsApp")
    expect(container.innerHTML).toContain("Klik Peta")
    expect(container.innerHTML).toContain("Klik Telepon")
    expect(container.innerHTML).toContain("Konfirmasi RSVP")
    expect(mockGet).toHaveBeenCalledWith("/events/evt_001/analytics")
  })

  it("shows the upsell and does NOT fetch analytics for a free user", async () => {
    mockGet.mockImplementation((url: string) => {
      if (url === "/subscriptions/current") {
        return Promise.resolve({ data: { data: { package: { code: "free" } } } })
      }
      return Promise.reject(new Error(`unexpected url: ${url}`))
    })

    const container = renderWidget()
    await act(async () => {})

    expect(container.innerHTML).toContain(UPSELL_TEXT)
    expect(mockGet).not.toHaveBeenCalledWith("/events/evt_001/analytics")
  })

  it("shows the upsell without crashing when the analytics GET returns 403", async () => {
    mockGet.mockImplementation((url: string) => {
      if (url === "/subscriptions/current") {
        return Promise.resolve({ data: { data: { package: { code: "premium" } } } })
      }
      if (url === "/events/evt_001/analytics") {
        return Promise.reject({ response: { status: 403 } })
      }
      return Promise.reject(new Error(`unexpected url: ${url}`))
    })

    const container = renderWidget()
    await act(async () => {})

    expect(container.innerHTML).toContain(UPSELL_TEXT)
    expect(container.innerHTML).not.toContain("120")
  })
})