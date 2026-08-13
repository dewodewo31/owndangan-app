import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "rustic-bohemian",
  name: "Rustic Bohemian",
  category: "Rustic",
  description:
    "Warm cream textures, earthy terracotta and sage, organic botanical separators and a handmade boho spirit. Built for outdoor and garden weddings.",
  tags: ["Rustic", "Bohemian", "Natural"],
  thumbnail:
    "https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#a0522d",
    secondary: "#6b3f23",
    background: "#faf3e7",
    surface: "#f4ead9",
    text: "#3d2f23",
    muted: "#8a7663",
    accent: "#8a9a5b",
    border: "#e2d3b8",
    radius: "1rem",
    fontHeading: "'Jost', sans-serif",
    fontBody: "'Jost', sans-serif",
    fontAccent: "'Pinyon Script', cursive",
    sectionSpacing: "5.5rem",
    contentWidth: "64rem",
    heroHeight: "88vh",
    animationDuration: "0.7s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "32px",
  },
  nav: "decorative-bottom",
  animation: { variant: "fade-up", stagger: true, parallax: false },
  decoration: "botanical",
  sections: [
    { key: "cover", variant: "split" },
    { key: "quote" },
    { key: "couple", variant: "stacked" },
    { key: "countdown" },
    { key: "events", variant: "timeline" },
    { key: "gallery", variant: "masonry" },
    { key: "location" },
    { key: "rsvp" },
    { key: "gift" },
    { key: "guestbook" },
    { key: "closing" },
  ],
}

export default definition
