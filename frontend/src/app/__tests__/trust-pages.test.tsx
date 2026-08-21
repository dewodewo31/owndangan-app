import { act } from "react"
import { createRoot, type Root } from "react-dom/client"

;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

class IntersectionObserverStub {
  observe() {}
  unobserve() {}
  disconnect() {}
  takeRecords() {
    return []
  }
}
;(globalThis as unknown as { IntersectionObserver: unknown }).IntersectionObserver =
  IntersectionObserverStub

jest.mock("next/link", () => ({
  __esModule: true,
  default: ({ children, href, ...props }: { children: React.ReactNode; href: string; [key: string]: unknown }) => (
    <a href={href} {...props}>{children}</a>
  ),
}))

import FAQPage from "@/app/faq/page"
import CaraOrderPage from "@/app/cara-order/page"
import TestimoniPage from "@/app/testimoni/page"
import Navbar from "@/components/marketing/navbar"
import { Footer } from "@/components/marketing/footer"

function render(element: React.ReactElement): { root: Root; container: HTMLElement } {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const root = createRoot(container)
  act(() => {
    root.render(element)
  })
  return { root, container }
}

function cleanup(root: Root, container: HTMLElement) {
  act(() => {
    root.unmount()
  })
  container.remove()
}

describe("trust pages", () => {
  it("FAQ page renders the shared FaqList content", () => {
    const { root, container } = render(<FAQPage />)
    expect(container.textContent).toContain("Pertanyaan yang")
    expect(container.textContent).toContain("Apa itu Owndangan?")
    expect(container.querySelectorAll("button[aria-expanded]").length).toBeGreaterThan(0)
    cleanup(root, container)
  })

  it("Cara Order page shows four steps", () => {
    const { root, container } = render(<CaraOrderPage />)
    const steps = Array.from(container.querySelectorAll("ol li"))
    expect(steps).toHaveLength(4)
    expect(container.textContent).toContain("Pilih Template")
    expect(container.textContent).toContain("Isi Detail Pernikahan")
    expect(container.textContent).toContain("Pilih Paket & Bayar")
    expect(container.textContent).toContain("Bagikan Undangan")
    cleanup(root, container)
  })

  it("Testimoni page shows placeholder and Contoh badges with no fabricated identities", () => {
    const { root, container } = render(<TestimoniPage />)
    expect(container.textContent).toContain(
      "Kumpulan testimoni pelanggan akan segera hadir. Hubungi kami untuk berbagi pengalaman Anda."
    )
    const badges = Array.from(container.querySelectorAll("span")).filter((s) => s.textContent === "Contoh")
    expect(badges.length).toBe(3)

    const quotes = Array.from(container.querySelectorAll("blockquote")).map((q) => q.textContent ?? "")
    expect(quotes).toHaveLength(3)
    for (const quote of quotes) {
      expect(quote).not.toMatch(/[A-Z][a-z]+\s+[A-Z][a-z]+/)
    }
    cleanup(root, container)
  })

  it("Navbar includes trust-page links", () => {
    const { root, container } = render(<Navbar />)
    for (const href of ["/faq", "/cara-order", "/testimoni"]) {
      expect(container.querySelector(`a[href="${href}"]`)).not.toBeNull()
    }
    cleanup(root, container)
  })

  it("Footer includes trust-page links", () => {
    const { root, container } = render(<Footer />)
    for (const href of ["/faq", "/cara-order", "/testimoni"]) {
      expect(container.querySelector(`a[href="${href}"]`)).not.toBeNull()
    }
    cleanup(root, container)
  })
})
