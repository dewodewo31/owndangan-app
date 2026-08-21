"use client"

import { ScrollReveal } from "@/components/animation/scroll-reveal"
import { FaqList } from "./faq-list"

export default function FAQSection() {
  return (
    <section id="faq" className="py-24 bg-background">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-16">
            <span className="inline-block px-3 py-1 rounded-full bg-primary-container border border-primary/20 text-primary text-sm font-medium mb-4">
              FAQ
            </span>
            <h2 className="text-3xl sm:text-4xl font-display font-bold text-foreground tracking-tight">
              Pertanyaan yang
              <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                sering diajukan.
              </span>
            </h2>
          </div>
        </ScrollReveal>

        <FaqList />
      </div>
    </section>
  )
}
