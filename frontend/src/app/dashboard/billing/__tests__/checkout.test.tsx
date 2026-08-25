import { act } from "react"
import { createRoot, type Root } from "react-dom/client"
import BillingPage from "@/app/dashboard/billing/page"

jest.mock("next/link", () => ({
  __esModule: true,
  default: ({ children, href, ...props }: any) => <a href={href} {...props}>{children}</a>,
}))

// Render the page without the auth gate / chrome so we can drive the checkout
// directly and assert what gets passed to Midtrans Snap.
jest.mock("@/components/dashboard/protected-route", () => ({
  __esModule: true,
  default: ({ children }: any) => children,
}))
jest.mock("@/components/dashboard/dashboard-layout", () => ({
  __esModule: true,
  default: ({ children }: any) => children,
}))

jest.mock("@/lib/api", () => ({
  __esModule: true,
  default: { get: jest.fn(), post: jest.fn() },
}))

;(globalThis as any).IntersectionObserver = class {
  observe() {}
  unobserve() {}
  disconnect() {}
  takeRecords() {
    return []
  }
}
;(globalThis as any).IS_REACT_ACT_ENVIRONMENT = true

const SNAP_TOKEN = "uuid-snap-token-abc123"
const ORDER_ID = "INV-should-never-be-passed-as-token"

async function flush(ms = 80) {
  await act(async () => {
    await new Promise((r) => setTimeout(r, ms))
  })
}

describe("billing checkout token contract", () => {
  beforeEach(() => {
    const api: any = require("@/lib/api").default
    api.get.mockReset()
    api.post.mockReset()
    api.get.mockImplementation((url: string) => {
      if (url.includes("/subscriptions/current")) return Promise.resolve({ data: { data: null } })
      if (url.includes("/payments/transactions")) return Promise.resolve({ data: { data: [] } })
      if (url.includes("/packages"))
        return Promise.resolve({
          data: { data: [{ id: "pkg-1", name: "Basic", code: "basic", price: 99000, is_active: true }] },
        })
      return Promise.resolve({ data: { data: null } })
    })
    api.post.mockImplementation((url: string) => {
      if (url.includes("/payments/snap"))
        return Promise.resolve({ data: { success: true, data: { snap_token: SNAP_TOKEN, order_id: ORDER_ID } } })
      return Promise.resolve({ data: { data: null } })
    })
    ;(window as any).snap = { pay: jest.fn() }
  })

  it("passes snap_token (not order_id) to snap.pay", async () => {
    const container = document.createElement("div")
    document.body.appendChild(container)
    const root: Root = createRoot(container)
    act(() => {
      root.render(<BillingPage />)
    })
    await flush()

    const btn = Array.from(container.querySelectorAll("button")).find((b) =>
      (b.textContent || "").includes("Berlangganan")
    )
    expect(btn).toBeTruthy()

    await act(async () => {
      ;(btn as HTMLButtonElement).click()
      await new Promise((r) => setTimeout(r, 60))
    })

    const pay = (window as any).snap.pay as jest.Mock
    expect(pay).toHaveBeenCalledTimes(1)
    expect(pay.mock.calls[0][0]).toBe(SNAP_TOKEN)
    expect(pay.mock.calls[0][0]).not.toBe(ORDER_ID)

    act(() => {
      root.unmount()
    })
    container.remove()
  })
})
