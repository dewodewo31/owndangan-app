import { act } from "react"
import { createRoot, type Root } from "react-dom/client"

;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

const authState = { user: null as any, loading: false, isAuthenticated: false, logout: jest.fn() }
jest.mock("@/providers/auth-context", () => ({ useAuth: () => authState }))

const pathnameState = { value: "/dashboard" }
jest.mock("next/navigation", () => ({
  usePathname: () => pathnameState.value,
}))

jest.mock("next/link", () => ({
  __esModule: true,
  default: ({ children, ...props }: { children: React.ReactNode; [key: string]: unknown }) => (
    <a {...props}>{children}</a>
  ),
}))

import DashboardLayout from "../dashboard-layout"

const roots: Root[] = []

function render() {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const root = createRoot(container)
  roots.push(root)
  act(() => {
    root.render(<DashboardLayout>CONTENT</DashboardLayout>)
  })
  return container
}

afterEach(() => {
  roots.forEach((r) => act(() => r.unmount()))
  roots.length = 0
  authState.user = null
  pathnameState.value = "/dashboard"
  jest.clearAllMocks()
})

describe("DashboardLayout admin nav item", () => {
  it("does NOT show Admin Dashboard for a non-admin user", () => {
    authState.user = { name: "U", email: "u@b.c", role: "user" }
    const container = render()
    expect(container.querySelector('a[href="/admin"]')).toBeNull()
    expect(container.innerHTML).not.toContain("Admin Dashboard")
  })

  it("does NOT show Admin Dashboard when unauthenticated", () => {
    authState.user = null
    const container = render()
    expect(container.querySelector('a[href="/admin"]')).toBeNull()
  })

  it("shows Admin Dashboard linking to /admin for an admin", () => {
    authState.user = { name: "A", email: "a@b.c", role: "admin" }
    const container = render()
    const link = container.querySelector('a[href="/admin"]') as HTMLAnchorElement
    expect(link).not.toBeNull()
    expect(container.innerHTML).toContain("Admin Dashboard")
  })

  it("marks Admin Dashboard active when inside /admin", () => {
    authState.user = { name: "A", email: "a@b.c", role: "admin" }
    pathnameState.value = "/admin"
    const container = render()
    const link = container.querySelector('a[href="/admin"]') as HTMLAnchorElement
    expect(link.className).toContain("shadow-elevation-1")
  })

  it("does NOT mark Admin Dashboard active on the user dashboard", () => {
    authState.user = { name: "A", email: "a@b.c", role: "admin" }
    pathnameState.value = "/dashboard"
    const container = render()
    const link = container.querySelector('a[href="/admin"]') as HTMLAnchorElement
    expect(link.className).not.toContain("shadow-elevation-1")
  })
})
