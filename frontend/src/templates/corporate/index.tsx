import type { TemplateDefinition } from "@/templates/types"

export const definition: TemplateDefinition = {
  kind: "corporate",
  name: "Corporate",
  category: "Corporate",
  occasions: ["corporate", "event"],
  description:
    "A sharp, professional invitation for business events — product launches, seminars and company gatherings. Dark slate and indigo with geometric accents.",
  tags: ["Corporate", "Business", "Geometric"],
  thumbnail:
    "https://images.unsplash.com/photo-1540575467063-178a50c2df87?auto=format&fit=crop&w=800&q=60",
  theme: {
    primary: "#818cf8",
    secondary: "#1e293b",
    background: "#0f172a",
    surface: "#1e293b",
    text: "#e2e8f0",
    muted: "#94a3b8",
    accent: "#a5b4fc",
    border: "#334155",
    radius: "0.25rem",
    fontHeading: "'Archivo', sans-serif",
    fontBody: "'Manrope', sans-serif",
    fontAccent: "'Archivo', sans-serif",
    sectionSpacing: "5rem",
    contentWidth: "64rem",
    heroHeight: "88vh",
    animationDuration: "0.6s",
    animationEasing: "cubic-bezier(0.22,1,0.36,1)",
    revealDistance: "32px",
  },
  night: {
    primary: "#6d78e8",
    secondary: "#16213a",
    background: "#0a0f1e",
    surface: "#131c2e",
    text: "#e2e8f0",
    muted: "#94a3b8",
    accent: "#8f9ef0",
    border: "#26334d",
  },
  nav: "bottom-floating",
  animation: { variant: "fade-up", stagger: true, parallax: false },
  decoration: "geometric",
  timelineStyle: "modern",
  sections: [
    { key: "cover", variant: "cinematic" },
    { key: "events", variant: "timeline" },
    { key: "location" },
    { key: "rsvp" },
    { key: "gallery", variant: "grid" },
    { key: "closing" },
  ],
}

export default definition