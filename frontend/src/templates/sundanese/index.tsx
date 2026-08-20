import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "sundanese",
  name: "Sundanese Traditional",
  category: "Traditional",
  occasions: ["pernikahan"],
  description:
    "Sundanese-inspired West Java elegance with a fresh leaf-green and warm natural palette, botanical accents and organic minimal composition.",
  tags: ["Sundanese", "Natural", "Elegant"],
  thumbnail:
    "https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#2e7d52",
    secondary: "#1f5c3d",
    background: "#f4f1e8",
    surface: "#ece7d8",
    text: "#2b3a2e",
    muted: "#7a857a",
    accent: "#c9a14a",
    border: "#dcd6c4",
    radius: "1.25rem",
    fontHeading: "'Petrona', serif",
    fontBody: "'Inter', sans-serif",
    fontAccent: "'Caveat', cursive",
    sectionSpacing: "5.5rem",
    contentWidth: "66rem",
    heroHeight: "90vh",
    animationDuration: "0.8s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "32px",
  },
  nav: "floating-menu",
  animation: { variant: "fade-up", stagger: true, parallax: false },
  decoration: "botanical",
  sections: [
    { key: "cover", variant: "centered" },
    { key: "quote" },
    { key: "couple", variant: "editorial" },
    { key: "parents" },
    { key: "events", variant: "side" },
    { key: "gallery", variant: "columns" },
    { key: "location" },
    { key: "rsvp" },
    { key: "gift" },
    { key: "guestbook" },
    { key: "closing" },
  ],
}

export default definition
