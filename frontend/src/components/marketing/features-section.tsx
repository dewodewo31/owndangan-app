"use client"

import { ScrollReveal, StaggerContainer, StaggerItem } from "@/components/animation/scroll-reveal"
import { Card, CardContent } from "@/components/ui/card"
import { features } from "@/data/marketing"

export default function FeaturesSection() {
  return (
    <section id="features" className="py-24 bg-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-16 max-w-2xl mx-auto">
            <span className="inline-block px-3 py-1 rounded-full bg-primary-container border border-primary/20 text-primary text-sm font-medium mb-4">
              Fitur Lengkap
            </span>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-display font-bold text-foreground tracking-tight">
              Semua yang kamu butuhkan
              <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                untuk hari istimewa.
              </span>
            </h2>
            <p className="mt-4 text-muted-foreground text-lg">
              Dari template hingga analitik RSVP — semua dalam satu platform yang mudah digunakan.
            </p>
          </div>
        </ScrollReveal>

        <StaggerContainer className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          {features.map((feature) => (
            <StaggerItem key={feature.id} className="h-full">
              <Card className="h-full transition-all duration-300 hover:-translate-y-1.5 hover:shadow-xl hover:shadow-primary/10 hover:border-primary/40 group">
                <CardContent className="p-6">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-container to-secondary-container border border-primary/10 flex items-center justify-center mb-5 transition-colors group-hover:from-primary group-hover:to-secondary">
                    <feature.icon className="h-6 w-6 text-primary transition-transform group-hover:scale-110" />
                  </div>
                  <h3 className="font-semibold text-foreground mb-2">{feature.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">{feature.description}</p>
                </CardContent>
              </Card>
            </StaggerItem>
          ))}
        </StaggerContainer>
      </div>
    </section>
  )
}
