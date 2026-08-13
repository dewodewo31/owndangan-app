"use client"

import { ScrollReveal } from "@/components/animation/scroll-reveal"
import { trustIndicators } from "@/data/marketing"

export default function TrustBar() {
  return (
    <section className="py-12 bg-background border-y border-border">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <p className="text-center text-sm font-medium text-muted-foreground mb-8">
            Dipercaya ribuan pasangan untuk momen istimewa mereka
          </p>
          <div className="flex flex-wrap justify-center gap-x-12 gap-y-4">
            {trustIndicators.map((item) => (
              <span
                key={item}
                className="flex items-center gap-2 text-sm font-semibold text-foreground/70"
              >
                <svg className="h-4 w-4 text-primary" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2l2.4 7.2H22l-6 4.8 2.4 7.2-6.4-4.6-6.4 4.6 2.4-7.2-6-4.8h7.6z" />
                </svg>
                {item}
              </span>
            ))}
          </div>
        </ScrollReveal>
      </div>
    </section>
  )
}
