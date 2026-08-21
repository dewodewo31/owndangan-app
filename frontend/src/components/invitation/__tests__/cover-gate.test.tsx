import { act } from "react"
import { createRoot } from "react-dom/client"
import { CoverGate } from "@/components/invitation/cover-gate"
import type { InvitationModel } from "@/templates/types"

const model = {
  slug: "gate-test",
  names: { full: "A & B" },
  parents: {},
  events: {},
  gallery: [],
  loveStories: [],
  guestbook: [],
  sections: {},
} as InvitationModel

function mount(ui: React.ReactElement) {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const root = createRoot(container)
  act(() => {
    root.render(ui)
  })
  return { root, container }
}

function unmount(root: ReturnType<typeof createRoot>, container: HTMLElement) {
  act(() => {
    root.unmount()
  })
  container.remove()
}

describe("CoverGate body scroll lock", () => {
  afterEach(() => {
    document.body.style.overflow = ""
  })

  it("locks body scroll by default", () => {
    const { root, container } = mount(<CoverGate model={model} heroImage="x" primaryColor="#000" />)
    expect(document.body.style.overflow).toBe("hidden")
    unmount(root, container)
    expect(document.body.style.overflow).toBe("")
  })

  it("leaves body scroll untouched when lockBodyScroll is false", () => {
    const { root, container } = mount(
      <CoverGate model={model} heroImage="x" primaryColor="#000" lockBodyScroll={false} />
    )
    expect(document.body.style.overflow).toBe("")
    unmount(root, container)
    expect(document.body.style.overflow).toBe("")
  })
})
