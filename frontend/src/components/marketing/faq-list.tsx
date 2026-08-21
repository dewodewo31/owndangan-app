"use client"

import { useState } from "react"
import { ChevronDown } from "lucide-react"
import { ScrollReveal } from "@/components/animation/scroll-reveal"
import { faqItems } from "@/data/marketing"
import { cn } from "@/lib/utils"

export function FaqList() {
  const [openIndex, setOpenIndex] = useState<number | null>(0)

  return (
    <div className="space-y-3">
      {faqItems.map((item, index) => (
        <ScrollReveal key={index} delay={index * 0.04}>
          <div
            className={cn(
              "border rounded-xl overflow-hidden transition-all duration-300 bg-white",
              openIndex === index
                ? "border-primary/40 shadow-lg shadow-primary/5"
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
  )
}
