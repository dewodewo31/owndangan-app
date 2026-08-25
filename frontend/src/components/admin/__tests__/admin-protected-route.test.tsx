import { act } from "react"
import { createRoot, type Root } from "react-dom/client"

;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

const authState = { user: null as any, loading: false, isAuthenticated: false, logout: jest.fn() }
jest.mock("@/providers/auth-context", () => ({ useAuth: () => authState }))

const mockPush = jest.fn()
jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}))

import AdminProtectedRoute from "../admin-protected-route"

const roots: Root[] = []

function render() {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const root = createRoot(container)
  roots.push(root)
  act(() => {
    root.render(<AdminProtectedRoute>SECRET</AdminProtectedRoute>)
  })
  return container
}

afterEach(() => {
  roots.forEach((r) => act(() => r.unmount()))
  roots.length = 0
  authState.user = null
  authState.loading = false
  authState.isAuthenticated = false
  mockPush.mockClear()
})

describe("AdminProtectedRoute", () => {
  it("renders children for an admin", () => {
    authState.user = { name: "A", email: "a@b.c", role: "admin" }
    authState.isAuthenticated = true
    const container = render()
    expect(container.innerHTML).toContain("SECRET")
    expect(mockPush).not.toHaveBeenCalled()
  })

  it("redirects a non-admin user to /dashboard", () => {
    authState.user = { name: "U", email: "u@b.c", role: "user" }
    authState.isAuthenticated = true
    const container = render()
    expect(mockPush).toHaveBeenCalledWith("/dashboard")
    expect(container.innerHTML).not.toContain("SECRET")
  })

  it("redirects an unauthenticated user to /login", () => {
    authState.user = null
    authState.isAuthenticated = false
    const container = render()
    expect(mockPush).toHaveBeenCalledWith("/login")
    expect(container.innerHTML).not.toContain("SECRET")
  })
})
