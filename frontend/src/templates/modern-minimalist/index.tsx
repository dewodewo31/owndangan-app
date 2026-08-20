import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "modern-minimalist",
  name: "Modern Minimalist",
  category: "Modern",
  occasions: ["pernikahan"],
  description:
    "Clean whitespace, large editorial typography, asymmetric sections and a restrained palette. Built for the modern urban wedding.",
  tags: ["Modern", "Minimal", "Editorial"],
  thumbnail:
    "https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#1f2937",
    secondary: "#111827",
    background: "#ffffff",
    surface: "#f7f7f5",
    text: "#1f2937",
    muted: "#6b7280",
    accent: "#9ca3af",
    border: "#e5e7eb",
    radius: "0.5rem",
    fontHeading: "'Playfair Display', serif",
    fontBody: "'Inter', sans-serif",
    fontAccent: "'Playfair Display', serif",
    sectionSpacing: "5rem",
    contentWidth: "64rem",
    heroHeight: "92vh",
    animationDuration: "0.7s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "36px",
  },
  nav: "bottom-floating",
  animation: { variant: "fade-up", stagger: true, parallax: false },
  decoration: "none",
  sections: [
    { key: "cover", variant: "minimal" },
    { key: "quote" },
    { key: "couple", variant: "portrait" },
    { key: "events", variant: "cards" },
    { key: "gallery", variant: "grid" },
    { key: "rsvp" },
    { key: "gift" },
    { key: "guestbook" },
    { key: "closing" },
  ],
}

export default definition
