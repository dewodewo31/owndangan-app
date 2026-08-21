import type { Metadata } from "next"
import Link from "next/link"
import { ArrowLeft, BadgeCheck, MessageCircle, Quote } from "lucide-react"
import { contactWhatsApp } from "@/data/marketing"
import { Button } from "@/components/ui/button"

export const metadata: Metadata = {
  title: "Testimoni — Owndangan",
  description: "Testimoni pelanggan Owndangan.",
}

const sampleQuotes = [
  "Ucapan pelanggan akan tampil di sini.",
  "Setiap testimoni akan ditinjau sebelum dipublikasikan.",
  "Bagikan pengalaman Anda bersama Owndangan.",
]

export default function TestimoniPage() {
  return (
    <main className="min-h-screen bg-background">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center mb-10">
          <span className="inline-block px-3 py-1 rounded-full bg-primary-container border border-primary/20 text-primary text-sm font-medium mb-4">
            Testimoni
          </span>
          <h1 className="text-3xl sm:text-4xl font-display font-bold text-foreground tracking-tight">
            Kata mereka tentang
            <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              Owndangan.
            </span>
          </h1>
          <p className="mt-5 text-muted-foreground leading-relaxed max-w-xl mx-auto">
            Kumpulan testimoni pelanggan akan segera hadir. Hubungi kami untuk
            berbagi pengalaman Anda.
          </p>
          <div className="mt-7 flex flex-col sm:flex-row items-center justify-center gap-4">
            <a href={contactWhatsApp} target="_blank" rel="noopener noreferrer">
              <Button className="bg-[#25D366] hover:bg-[#1fb857] text-white">
                <MessageCircle className="mr-2 h-4 w-4" />
                WA Kami
              </Button>
            </a>
            <Link href="/">
              <Button variant="outline">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Kembali ke Beranda
              </Button>
            </Link>
          </div>
        </div>

        <section aria-label="Contoh tampilan testimoni" className="mt-8">
          <div className="flex items-center gap-2 mb-5">
            <BadgeCheck className="h-4 w-4 text-primary" />
            <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
              Contoh tampilan
            </h2>
          </div>
          <div className="grid gap-5 sm:grid-cols-3">
            {sampleQuotes.map((quote) => (
              <figure
                key={quote}
                className="relative rounded-2xl border border-dashed border-border bg-white p-6 opacity-80"
              >
                <span className="absolute top-4 right-4 rounded-full bg-secondary-container px-2.5 py-0.5 text-xs font-medium text-secondary-foreground">
                  Contoh
                </span>
                <Quote className="h-5 w-5 text-primary/40 mb-3" />
                <blockquote className="text-sm text-muted-foreground leading-relaxed">
                  {quote}
                </blockquote>
              </figure>
            ))}
          </div>
        </section>
      </div>
    </main>
  )
}
