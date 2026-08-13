"use client"

import { useState } from "react"
import { Eye, Heart } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ScrollReveal, StaggerContainer, StaggerItem } from "@/components/animation/scroll-reveal"
import { templates, templateCategories } from "@/data/marketing"
import { cn } from "@/lib/utils"

export default function TemplateShowcase() {
  const [activeCategory, setActiveCategory] = useState("All")

  const filteredTemplates =
    activeCategory === "All"
      ? templates
      : templates.filter((t) => t.style === activeCategory)

  return (
    <section id="templates" className="py-24 bg-gradient-to-b from-slate-50 to-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-12 max-w-2xl mx-auto">
            <span className="inline-block px-3 py-1 rounded-full bg-violet-50 border border-violet-200 text-primary text-sm font-medium mb-4">
              Galeri Template
            </span>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-foreground tracking-tight">
              Template yang sesuai
              <span className="block bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent">
                dengan cerita kalian.
              </span>
            </h2>
          </div>
        </ScrollReveal>

        <ScrollReveal delay={0.1}>
          <div className="flex flex-wrap justify-center gap-2 mb-12">
            {templateCategories.map((category) => (
              <button
                key={category}
                onClick={() => setActiveCategory(category)}
                className={cn(
                  "px-4 py-2 rounded-full text-sm font-medium transition-all",
                  activeCategory === category
                    ? "bg-gradient-to-r from-indigo-600 to-violet-600 text-white shadow-md shadow-indigo-500/25"
                    : "bg-background border border-border text-muted-foreground hover:text-foreground hover:border-primary/40"
                )}
              >
                {category}
              </button>
            ))}
          </div>
        </ScrollReveal>

        <StaggerContainer className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredTemplates.map((template) => (
            <StaggerItem key={template.id} className="h-full">
              <div className="group relative bg-white rounded-2xl overflow-hidden border border-border shadow-sm transition-all duration-300 hover:-translate-y-1.5 hover:shadow-2xl hover:shadow-indigo-500/10">
                <div className="aspect-[3/4] bg-gradient-to-br from-indigo-50 via-violet-50 to-fuchsia-50 flex items-center justify-center p-6">
                  <div className="text-center">
                    <div className="w-20 h-20 mx-auto rounded-full bg-gradient-to-br from-indigo-100 to-violet-100 border border-indigo-200 flex items-center justify-center mb-4">
                      <span className="text-3xl">💒</span>
                    </div>
                    <p className="text-xs text-muted-foreground font-medium tracking-widest uppercase mb-1">
                      Undangan Pernikahan
                    </p>
                    <p className="text-lg font-semibold text-foreground">Andi & Siti</p>
                  </div>
                </div>
                <div className="p-5 flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold text-foreground">{template.name}</h3>
                    <p className="text-sm text-muted-foreground">{template.style}</p>
                  </div>
                  <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
                    <Heart className="h-3.5 w-3.5 text-rose-400" /> 128
                  </span>
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
      </div>
    </section>
  )
}
