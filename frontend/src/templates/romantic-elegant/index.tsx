import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "romantic-elegant",
  name: "Romantic Elegant",
  category: "Classic",
  occasions: ["pernikahan"],
  description:
    "Soft ivory palette with gold accents, romantic serif typography and graceful curved cards. A luxurious love-story feel built to sweep your loved ones away.",
  tags: ["Romantic", "Elegant", "Luxury"],
  thumbnail:
    "https://images.unsplash.com/photo-1465495976272-a759b7a8e978?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#9a7b4f",
    secondary: "#f2e9d8",
    background: "#fdfaf5",
    surface: "#f7efe1",
    text: "#3d3330",
    muted: "#8a7a6e",
    accent: "#c9a96a",
    border: "#ecd9b8",
    radius: "1.25rem",
    fontHeading: "'Cormorant Garamond', serif",
    fontBody: "'Lora', serif",
    fontAccent: "'Great Vibes', cursive",
    sectionSpacing: "6rem",
    contentWidth: "60rem",
    heroHeight: "100vh",
    animationDuration: "0.9s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "40px",
  },
  nav: "floating-menu",
  animation: { variant: "fade-up", stagger: true, parallax: false },
  decoration: "floral",
  sections: [
    { key: "cover", variant: "centered" },
    { key: "quote" },
    { key: "couple", variant: "portrait" },
    { key: "events", variant: "cards" },
    { key: "gallery", variant: "horizontal" },
    { key: "rsvp" },
    { key: "gift" },
    { key: "guestbook" },
    { key: "closing" },
  ],
}

export default definition
