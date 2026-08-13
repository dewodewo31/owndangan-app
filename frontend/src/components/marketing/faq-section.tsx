"use client"

import { useState } from "react"
import { ChevronDown } from "lucide-react"
import { ScrollReveal } from "@/components/animation/scroll-reveal"
import { faqItems } from "@/data/marketing"
import { cn } from "@/lib/utils"

export default function FAQSection() {
  const [openIndex, setOpenIndex] = useState<number | null>(0)

  return (
    <section id="faq" className="py-24 bg-background">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-16">
            <span className="inline-block px-3 py-1 rounded-full bg-indigo-50 border border-indigo-200 text-primary text-sm font-medium mb-4">
              FAQ
            </span>
            <h2 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
              Pertanyaan yang
              <span className="block bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent">
                sering diajukan.
              </span>
            </h2>
          </div>
        </ScrollReveal>

        <div className="space-y-3">
          {faqItems.map((item, index) => (
            <ScrollReveal key={index} delay={index * 0.04}>
              <div
                className={cn(
                  "border rounded-xl overflow-hidden transition-all duration-300 bg-white",
                  openIndex === index
                    ? "border-primary/40 shadow-lg shadow-indigo-500/5"
                    : "border-border hover:border-primary/30"
                )}
              >
                <button
                  className="w-full px-6 py-4 text-left flex items-center justify-between gap-4 hover:bg-accent/40 transition-colors"
                  onClick={() => setOpenIndex(openIndex === index ? null : index)}
                  aria-expanded={openIndex === index}
                >
                  <span className="font-medium text-foreground">{item.question}</span>
                  <ChevronDown
                    className={cn(
                      "h-4 w-4 text-primary shrink-0 transition-transform duration-300",
                      openIndex === index && "rotate-180"
                    )}
                  />
                </button>
                <div
                  className={cn(
                    "overflow-hidden transition-all duration-300",
                    openIndex === index ? "max-h-40" : "max-h-0"
                  )}
                >
                  <p className="px-6 pb-4 text-muted-foreground leading-relaxed text-sm">
                    {item.answer}
                  </p>
                </div>
              </div>
            </ScrollReveal>
          ))}
        </div>
      </div>
    </section>
  )
}
