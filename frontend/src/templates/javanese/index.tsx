import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "javanese",
  name: "Javanese Traditional",
  category: "Traditional",
  description:
    "Heritage Javanese elegance with batik-inspired motifs, warm maroon and earthy cream tones, and classical serif typography for a refined traditional Indonesian wedding.",
  tags: ["Traditional", "Javanese", "Heritage"],
  thumbnail:
    "https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#6b1f2a",
    secondary: "#4a1220",
    background: "#fbf3e6",
    surface: "#f4e4c9",
    text: "#3a2a1a",
    muted: "#8a6d4f",
    accent: "#b8860b",
    border: "#dcc8a0",
    radius: "0.375rem",
    fontHeading: "'DM Serif Display', serif",
    fontBody: "'Inter', sans-serif",
    fontAccent: "'Pinyon Script', cursive",
    sectionSpacing: "5.5rem",
    contentWidth: "62rem",
    heroHeight: "92vh",
    animationDuration: "0.8s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "30px",
  },
  nav: "decorative-bottom",
  animation: { variant: "fade", stagger: true, parallax: false },
  decoration: "batik",
  sections: [
    { key: "cover", variant: "framed" },
    { key: "quote" },
    { key: "couple", variant: "portrait" },
    { key: "parents" },
    { key: "events", variant: "cards" },
    { key: "gallery", variant: "masonry" },
    { key: "gift" },
    { key: "guestbook" },
    { key: "closing" },
  ],
}

export default definition
