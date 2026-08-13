import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "japanese-zen",
  name: "Japanese Zen",
  category: "Minimal",
  description:
    "Enso-inspired minimalism with sumi-ink typography, generous whitespace and thin organic borders. A calm, sophisticated canvas for the intimate wedding.",
  tags: ["Zen", "Minimal", "Calm"],
  thumbnail:
    "https://images.unsplash.com/photo-1520854229188-3da5971b5e0c?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#3a3a3a",
    secondary: "#2b2b2b",
    background: "#fbfaf7",
    surface: "#f1efe9",
    text: "#2b2b2b",
    muted: "#8a857c",
    accent: "#9c8f7a",
    border: "#e2ddd3",
    radius: "0.25rem",
    fontHeading: "'Shippori Mincho', serif",
    fontBody: "'Zen Maru Gothic', sans-serif",
    fontAccent: "'Shippori Mincho', serif",
    sectionSpacing: "6rem",
    contentWidth: "60rem",
    heroHeight: "90vh",
    animationDuration: "0.9s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "28px",
  },
  nav: "floating-menu",
  animation: { variant: "fade", stagger: true, parallax: false },
  decoration: "zen",
  sections: [
    { key: "cover", variant: "minimal" },
    { key: "quote" },
    { key: "couple", variant: "portrait" },
    { key: "events", variant: "side" },
    { key: "gallery", variant: "columns" },
    { key: "location" },
    { key: "rsvp" },
    { key: "gift" },
    { key: "closing" },
  ],
}

export default definition
