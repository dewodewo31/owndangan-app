"use client"

import { ScrollReveal } from "@/components/animation/scroll-reveal"
import { Button } from "@/components/ui/button"
import { Sparkles } from "lucide-react"
import Link from "next/link"

export default function FinalCTA() {
  return (
    <section className="py-24 bg-background">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-primary via-primary/80 to-secondary p-10 sm:p-16 text-center shadow-elevation-3">
            {/* Decorative circles */}
            <div className="absolute -top-24 -left-24 h-64 w-64 rounded-full bg-white/10 blur-2xl" />
            <div className="absolute -bottom-24 -right-24 h-64 w-64 rounded-full bg-white/10 blur-2xl" />
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-40 w-40 rounded-full border border-white/20" />

            <div className="relative">
              <span className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-white/15 text-white text-sm font-medium mb-6 backdrop-blur">
                <Sparkles className="h-4 w-4" /> Gratis untuk mulai
              </span>
              <h2 className="text-3xl sm:text-4xl lg:text-5xl font-display font-bold text-white tracking-tight leading-tight">
                Siap membuat undangan
                <span className="block text-primary-foreground/80">untuk hari istimewamu?</span>
              </h2>
              <p className="mt-6 text-lg text-primary-foreground/80 max-w-2xl mx-auto">
                Buat undangan digital yang indah, mudah dibagikan, dan mudah dikelola.
                Mulai sekarang — tanpa kartu kredit.
              </p>
              <div className="mt-10 flex flex-col sm:flex-row gap-4 justify-center">
                <Link href="/register">
                  <Button size="lg" className="w-full sm:w-auto bg-white text-primary hover:bg-primary-container border-0 shadow-lg">
                    Buat Undangan Gratis
                  </Button>
                </Link>
                <Link href="#templates">
                  <Button size="lg" variant="outline" className="w-full sm:w-auto bg-transparent border-white/40 text-white hover:bg-white/10">
                    Lihat Template
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </ScrollReveal>
      </div>
    </section>
  )
}
