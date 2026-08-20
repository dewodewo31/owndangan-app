"use client"

import { useState } from "react"
import { Eye, Heart, Church } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ScrollReveal, StaggerContainer, StaggerItem } from "@/components/animation/scroll-reveal"
import { TEMPLATES, templateOccasions, OCCASION_LABELS } from "@/templates"
import { cn } from "@/lib/utils"

export default function TemplateShowcase() {
  const [activeOccasion, setActiveOccasion] = useState("Semua")

  const occasions = ["Semua", ...templateOccasions()]
  const occasionLabel = (o: string) => OCCASION_LABELS[o] ?? o

  const filteredTemplates =
    activeOccasion === "Semua"
      ? TEMPLATES
      : TEMPLATES.filter((t) => t.occasions.includes(activeOccasion))

  return (
    <section id="templates" className="py-24 bg-gradient-to-b from-slate-50 to-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-12 max-w-2xl mx-auto">
            <span className="inline-block px-3 py-1 rounded-full bg-secondary-container border border-secondary/40 text-primary text-sm font-medium mb-4">
              Galeri Template
            </span>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-display font-bold text-foreground tracking-tight">
              Template yang sesuai
              <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                dengan cerita kalian.
              </span>
            </h2>
          </div>
        </ScrollReveal>

        <ScrollReveal delay={0.1}>
          <div className="flex flex-wrap justify-center gap-2 mb-12">
            {occasions.map((occasion) => (
              <button
                key={occasion}
                onClick={() => setActiveOccasion(occasion)}
                className={cn(
                  "px-4 py-2 rounded-full text-sm font-medium transition-all",
                  activeOccasion === occasion
                    ? "bg-gradient-to-r from-plum to-rosegold text-white shadow-elevation-1"
                    : "bg-background border border-border text-muted-foreground hover:text-foreground hover:border-primary/40"
                )}
              >
                {occasionLabel(occasion)}
              </button>
            ))}
          </div>
        </ScrollReveal>

        {filteredTemplates.length === 0 ? (
          <p className="text-center text-muted-foreground py-16">
            Belum ada template untuk kategori ini
          </p>
        ) : (
          <StaggerContainer className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredTemplates.map((template) => (
              <StaggerItem key={template.kind} className="h-full">
                <div className="group relative bg-white rounded-2xl overflow-hidden border border-border shadow-sm transition-all duration-300 hover:-translate-y-1.5 hover:shadow-2xl hover:shadow-primary/10">
                  <div className="aspect-[3/4] bg-gradient-to-br from-primary-container/60 via-secondary-container/50 to-tertiary-container/50 flex items-center justify-center p-6">
                    {template.thumbnail ? (
                      <img
                        src={template.thumbnail}
                        alt={template.name}
                        loading="lazy"
                        className="h-full w-full object-cover rounded-lg"
                      />
                    ) : (
                      <div className="text-center">
                        <div className="w-20 h-20 mx-auto rounded-full bg-gradient-to-br from-primary-container to-secondary-container border border-primary/20 flex items-center justify-center mb-4">
                          <Church className="h-9 w-9 text-primary" />
                        </div>
                        <p className="text-xs text-muted-foreground font-medium tracking-widest uppercase mb-1">
                          Undangan Pernikahan
                        </p>
                        <p className="text-lg font-semibold text-foreground">Andi & Siti</p>
                      </div>
                    )}
                  </div>
                  <div className="p-5">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <h3 className="font-semibold text-foreground">{template.name}</h3>
                        <p className="text-sm text-muted-foreground">{template.category}</p>
                      </div>
                      <span className="inline-flex items-center gap-1 text-xs text-muted-foreground shrink-0">
                        <Heart className="h-3.5 w-3.5 text-rose-400" /> 128
                      </span>
                    </div>
                    <div className="flex flex-wrap gap-1.5 mt-3">
                      {template.occasions.map((occasion) => (
                        <span
                          key={occasion}
                          className="px-2 py-0.5 rounded-full bg-secondary-container/60 border border-secondary/30 text-xs text-primary font-medium"
                        >
                          {occasionLabel(occasion)}
                        </span>
                      ))}
                    </div>
                  </div>
                  <div className="absolute inset-0 bg-gradient-to-t from-slate-900/70 via-slate-900/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <Button className="bg-white text-foreground hover:bg-white/90 border-0">
                      <Eye className="h-4 w-4 mr-2" /> Lihat Template
                    </Button>
                  </div>
                </div>
              </StaggerItem>
            ))}
          </StaggerContainer>
        )}
      </div>
    </section>
  )
}
