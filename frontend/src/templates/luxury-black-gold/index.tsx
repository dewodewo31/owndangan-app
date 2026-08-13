import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "luxury-black-gold",
  name: "Luxury Black & Gold",
  category: "Luxury",
  description:
    "Cinematic drama: near-black backdrop, gold accents, elegant serif typography and large imagery. Built for a high-end premium celebration.",
  tags: ["Luxury", "Black & Gold", "Premium"],
  thumbnail:
    "https://images.unsplash.com/photo-1465495976272-a759b7a8e978?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#d4af37",
    secondary: "#bfa04a",
    background: "#0e0e10",
    surface: "#1a1a1d",
    text: "#f5f0e6",
    muted: "#a9a29a",
    accent: "#d4af37",
    border: "#2c2c30",
    radius: "0.25rem",
    fontHeading: "'Playfair Display', serif",
    fontBody: "'Jost', sans-serif",
    fontAccent: "'Pinyon Script', cursive",
    sectionSpacing: "7rem",
    contentWidth: "68rem",
    heroHeight: "100vh",
    animationDuration: "0.9s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "48px",
  },
  nav: "side",
  animation: { variant: "scale", stagger: true, parallax: true },
  decoration: "gold",
  sections: [
    { key: "cover", variant: "cinematic" },
    { key: "quote" },
    { key: "couple", variant: "stacked" },
    { key: "events", variant: "cards" },
    { key: "gallery", variant: "horizontal" },
    { key: "location" },
    { key: "rsvp" },
    { key: "gift" },
    { key: "guestbook" },
    { key: "closing" },
  ],
}

export default definition
