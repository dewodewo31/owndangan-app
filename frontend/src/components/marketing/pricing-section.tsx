"use client"

import { Check, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { ScrollReveal, StaggerContainer, StaggerItem } from "@/components/animation/scroll-reveal"
import { pricingPlans } from "@/data/marketing"
import { formatCurrency } from "@/lib/utils"
import { cn } from "@/lib/utils"
import Link from "next/link"

export default function PricingSection() {
  return (
    <section id="pricing" className="py-24 bg-gradient-to-b from-slate-50 to-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ScrollReveal>
          <div className="text-center mb-16 max-w-2xl mx-auto">
            <span className="inline-block px-3 py-1 rounded-full bg-secondary-container border border-secondary/40 text-primary text-sm font-medium mb-4">
              Harga
            </span>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-display font-bold text-foreground tracking-tight">
              Paket yang sesuai
              <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                dengan kebutuhan kalian.
              </span>
            </h2>
            <p className="mt-4 text-muted-foreground text-lg">
              Mulai gratis, upgrade kapan saja. Tanpa komitmen.
            </p>
          </div>
        </ScrollReveal>

        <StaggerContainer className="grid grid-cols-1 md:grid-cols-3 gap-8 items-stretch">
          {pricingPlans.map((plan) => (
            <StaggerItem key={plan.id} className="h-full">
              <Card
                className={cn(
                  "relative h-full flex flex-col transition-all duration-300 hover:-translate-y-1.5",
                  plan.popular
                    ? "border-transparent bg-gradient-to-b from-plum to-rosegold text-white shadow-elevation-3 md:scale-[1.05]"
                    : "hover:shadow-xl hover:shadow-primary/10 hover:border-primary/40"
                )}
              >
                {plan.popular && (
                  <div className="absolute -top-3.5 left-1/2 -translate-x-1/2">
                    <span className="inline-flex items-center gap-1 px-4 py-1.5 bg-white text-primary text-xs font-bold rounded-full shadow-lg">
                      <Sparkles className="h-3.5 w-3.5" /> Paling Populer
                    </span>
                  </div>
                )}
                <CardContent className="p-8 flex flex-col h-full">
                  <h3 className={cn("text-lg font-semibold", plan.popular ? "text-white" : "text-foreground")}>
                    {plan.name}
                  </h3>
                  <p className={cn("text-sm mt-1", plan.popular ? "text-primary-foreground/80" : "text-muted-foreground")}>
                    {plan.description}
                  </p>

                  <div className="mt-6 mb-8">
                    <span className={cn("text-4xl font-bold tracking-tight", plan.popular ? "text-white" : "text-foreground")}>
                      {formatCurrency(plan.price)}
                    </span>
                    <span className={cn("text-sm", plan.popular ? "text-primary-foreground/70" : "text-muted-foreground")}>
                      /{plan.duration}
                    </span>
                  </div>

                  <ul className="space-y-3 mb-8 flex-1">
                    {plan.features.map((feature) => (
                      <li key={feature} className="flex items-start gap-2.5 text-sm">
                        <span
                          className={cn(
                            "inline-flex items-center justify-center h-5 w-5 rounded-full shrink-0 mt-0.5",
                            plan.popular ? "bg-white/20" : "bg-success/10"
                          )}
                        >
                          <Check className={cn("h-3 w-3", plan.popular ? "text-white" : "text-success")} />
                        </span>
                        <span className={plan.popular ? "text-primary-foreground" : "text-foreground"}>{feature}</span>
                      </li>
                    ))}
                  </ul>

                  <Link href="/register" className="mt-auto">
                    <Button
                      className={cn(
                        "w-full border-0",
                        plan.popular
                          ? "bg-white text-primary hover:bg-primary-container"
                          : "bg-gradient-to-r from-primary to-secondary hover:from-primary/90 hover:to-secondary/90 text-white"
                      )}
                    >
                      {plan.cta}
                    </Button>
                  </Link>
                </CardContent>
              </Card>
            </StaggerItem>
          ))}
        </StaggerContainer>
      </div>
    </section>
  )
}
