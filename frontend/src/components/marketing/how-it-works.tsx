"use client"

import { ScrollReveal, StaggerContainer, StaggerItem } from "@/components/animation/scroll-reveal"
import { howItWorks } from "@/data/marketing"
import { cn } from "@/lib/utils"

export default function HowItWorks() {
  return (
    <section id="how-it-works" className="py-24 bg-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-16 max-w-2xl mx-auto">
            <span className="inline-block px-3 py-1 rounded-full bg-indigo-50 border border-indigo-200 text-primary text-sm font-medium mb-4">
              Cara Kerja
            </span>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-foreground tracking-tight">
              Mulai dari nol, selesai
              <span className="block bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent">
                dalam 3 langkah mudah.
              </span>
            </h2>
          </div>
        </ScrollReveal>

        <StaggerContainer className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {howItWorks.map((step, index) => (
            <StaggerItem key={step.step} className="h-full">
              <div className="relative h-full">
                {index < howItWorks.length - 1 && (
                  <div className="hidden md:block absolute top-10 left-[60%] w-[80%] h-[2px] bg-gradient-to-r from-indigo-300 to-transparent" />
                )}
                <div className="relative bg-white rounded-2xl border border-border p-8 h-full text-center shadow-sm transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/10 hover:border-primary/40">
                  <div
                    className={cn(
                      "w-20 h-20 mx-auto rounded-2xl flex items-center justify-center mb-6 relative",
                      "bg-gradient-to-br from-indigo-500 to-violet-600 shadow-lg shadow-indigo-500/25"
                    )}
                  >
                    <span className="text-2xl font-bold text-white">{step.step}</span>
                    <span className="absolute -top-2 -right-2 h-5 w-5 rounded-full bg-white border-2 border-indigo-500" />
                  </div>
                  <h3 className="font-semibold text-lg text-foreground mb-3">{step.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">{step.description}</p>
                </div>
              </div>
            </StaggerItem>
          ))}
        </StaggerContainer>
      </div>
    </section>
  )
}
