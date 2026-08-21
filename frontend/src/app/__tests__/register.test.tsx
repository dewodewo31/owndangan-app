import { act } from "react"
import { createRoot } from "react-dom/client"

;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

const mockPush = jest.fn()
jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}))

import { AuthProvider, useAuth } from "@/providers/auth-context"

function Probe() {
  const { register } = useAuth()
  return (
    <button onClick={() => register({ name: "A", email: "a@b.c", password: "secret123" })}>
      go
    </button>
  )
}

async function renderAndRegister() {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const root = createRoot(container)
  await act(async () => {
    root.render(
      <AuthProvider>
        <Probe />
      </AuthProvider>
    )
  })
  await act(async () => {
    container.querySelector("button")?.click()
  })
  return { root, container }
}

describe("register auto-login flow", () => {
  beforeEach(() => {
    localStorage.clear()
    mockPush.mockClear()
  })

  it("stores tokens and redirects to /dashboard on success", async () => {
    ;(globalThis as unknown as { fetch: unknown }).fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        data: { access_token: "tok-123", refresh_token: "ref-456" },
      }),
    })

    const { root, container } = await renderAndRegister()

    expect(localStorage.getItem("access_token")).toBe("tok-123")
    expect(localStorage.getItem("refresh_token")).toBe("ref-456")
    expect(mockPush).toHaveBeenCalledWith("/dashboard")

    act(() => root.unmount())
    container.remove()
  })

  it("redirects to /login without storing a session when tokens are missing", async () => {
    ;(globalThis as unknown as { fetch: unknown }).fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ data: {} }),
    })

    const { root, container } = await renderAndRegister()

    expect(mockPush).toHaveBeenCalledWith("/login?notice=sesi")
    expect(localStorage.getItem("access_token")).toBeNull()

    act(() => root.unmount())
    container.remove()
  })
})
