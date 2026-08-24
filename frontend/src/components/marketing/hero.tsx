"use client"

import { Button } from "@/components/ui/button"
import { ScrollReveal } from "@/components/animation/scroll-reveal"
import { Sparkles, Users, MessageSquare, CheckCircle2, Gem } from "lucide-react"
import Link from "next/link"

const heroStats = [
  { value: "10K+", label: "Undangan dibuat" },
  { value: "50K+", label: "Tamu terkelola" },
  { value: "4.9/5", label: "Rating pengguna" },
]

export default function Hero() {
  return (
    <section className="relative min-h-screen flex items-center overflow-hidden pt-24 pb-16">
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary-container/50 via-background to-secondary-container/50" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,rgba(99,102,241,0.12),transparent_50%)]" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_left,rgba(139,92,246,0.10),transparent_50%)]" />
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-primary/40 to-transparent" />

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <div className="grid lg:grid-cols-2 gap-16 items-center">
          {/* Left Content */}
          <div className="text-center lg:text-left">
            <ScrollReveal>
              <span className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-white/70 backdrop-blur border border-primary/20 text-primary text-sm font-medium mb-6 shadow-sm">
                <Sparkles className="h-3.5 w-3.5" />
                Platform Undangan Digital Modern
              </span>
            </ScrollReveal>

            <ScrollReveal delay={0.1}>
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-display font-bold text-foreground leading-[1.1] tracking-tight">
                Undangan Pernikahan
                <span className="block bg-gradient-to-r from-primary via-primary/80 to-secondary bg-clip-text text-transparent">
                  Digital yang Elegan.
                </span>
              </h1>
            </ScrollReveal>

            <ScrollReveal delay={0.2}>
              <p className="mt-6 text-lg text-muted-foreground max-w-xl mx-auto lg:mx-0 leading-relaxed">
                Buat undangan digital yang indah, kelola tamu secara otomatis,
                terima RSVP real-time, dan bagikan momen spesial melalui WhatsApp
                dalam hitungan menit.
              </p>
            </ScrollReveal>

            <ScrollReveal delay={0.3}>
              <div className="mt-8 flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                <Link href="/register">
                  <Button size="lg" className="w-full sm:w-auto bg-gradient-to-r from-primary to-secondary hover:from-primary/90 hover:to-secondary/90 border-0">
                    Buat Undangan Gratis
                  </Button>
                </Link>
                <Link href="#templates">
                  <Button size="lg" variant="outline" className="w-full sm:w-auto bg-background/60">
                    Lihat Template
                  </Button>
                </Link>
              </div>
            </ScrollReveal>

            <ScrollReveal delay={0.4}>
              <div className="mt-8 flex flex-wrap gap-x-6 gap-y-2 justify-center lg:justify-start text-sm text-muted-foreground">
                <span className="flex items-center gap-1.5">
                  <CheckCircle2 className="h-4 w-4 text-success" /> Mobile Friendly
                </span>
                <span className="flex items-center gap-1.5">
                  <CheckCircle2 className="h-4 w-4 text-success" /> RSVP Management
                </span>
                <span className="flex items-center gap-1.5">
                  <CheckCircle2 className="h-4 w-4 text-success" /> WhatsApp Sharing
                </span>
              </div>
            </ScrollReveal>

            <ScrollReveal delay={0.5}>
              <div className="mt-10 grid grid-cols-3 gap-6 max-w-sm mx-auto lg:mx-0 border-t border-border pt-8">
                {heroStats.map((stat) => (
                  <div key={stat.label}>
                    <p className="text-2xl font-bold text-foreground">{stat.value}</p>
                    <p className="text-xs text-muted-foreground mt-0.5">{stat.label}</p>
                  </div>
                ))}
              </div>
            </ScrollReveal>
          </div>

          {/* Right Visual */}
          <ScrollReveal direction="left" delay={0.2}>
            <div className="relative mx-auto max-w-md">
              {/* Glow behind mockup */}
              <div className="absolute -inset-6 bg-gradient-to-br from-primary/20 via-secondary/15 to-tertiary/20 rounded-[2.5rem] blur-2xl" />

              {/* Main Preview Card */}
              <div className="relative bg-white/90 backdrop-blur rounded-3xl shadow-elevation-3 p-4 border border-border">
                <div className="flex items-center justify-between px-2 py-2 mb-3">
                  <div className="flex gap-1.5">
                    <span className="h-3 w-3 rounded-full bg-red-400" />
                    <span className="h-3 w-3 rounded-full bg-yellow-400" />
                    <span className="h-3 w-3 rounded-full bg-green-400" />
                  </div>
                  <span className="text-xs text-muted-foreground font-medium">owndangan.app/andi-siti</span>
                </div>
                <div className="aspect-[4/5] rounded-2xl bg-gradient-to-br from-primary-container via-secondary-container to-tertiary-container flex items-center justify-center overflow-hidden relative">
                  <div className="text-center p-8">
                    <div className="w-20 h-20 mx-auto rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center mb-5 shadow-lg shadow-primary/30">
                      <Gem className="h-9 w-9 text-white" />
                    </div>
                    <p className="text-xs text-muted-foreground font-medium tracking-widest uppercase mb-2">
                      Undangan Pernikahan
                    </p>
                    <h3 className="font-bold text-2xl text-foreground">Andi & Siti</h3>
                    <p className="text-sm text-muted-foreground mt-3">15 Agustus 2026</p>
                    <p className="text-xs text-muted-foreground mt-1">Grand Ballroom, Jakarta</p>
                  </div>
                </div>
              </div>

              {/* Floating Cards */}
              <div className="absolute -top-5 -right-4 bg-white rounded-2xl shadow-xl border border-border p-3.5 animate-float">
                <div className="flex items-center gap-2.5">
                  <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-success/10">
                    <CheckCircle2 className="h-4 w-4 text-success" />
                  </span>
                  <div>
                    <p className="text-xs font-semibold">Published</p>
                    <p className="text-[10px] text-muted-foreground">Live di internet</p>
                  </div>
                </div>
              </div>

              <div className="absolute -bottom-5 -left-4 bg-white rounded-2xl shadow-xl border border-border p-3.5 animate-float" style={{ animationDelay: "1.2s" }}>
                <div className="flex items-center gap-2.5">
                  <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-primary-container">
                    <Users className="h-4 w-4 text-primary" />
                  </span>
                  <div>
                    <p className="text-xs font-semibold">127 Tamu</p>
                    <p className="text-[10px] text-muted-foreground">Sudah dikelola</p>
                  </div>
                </div>
              </div>

              <div className="absolute top-1/2 -right-6 bg-white rounded-2xl shadow-xl border border-border p-3.5 animate-float" style={{ animationDelay: "2.4s" }}>
                <div className="flex items-center gap-2.5">
                  <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-secondary-container">
                    <MessageSquare className="h-4 w-4 text-primary" />
                  </span>
                  <div>
                    <p className="text-xs font-semibold">86 RSVP Hadir</p>
                    <p className="text-[10px] text-muted-foreground">Konfirmasi real-time</p>
                  </div>
                </div>
              </div>
            </div>
          </ScrollReveal>
        </div>
      </div>
    </section>
  )
}
