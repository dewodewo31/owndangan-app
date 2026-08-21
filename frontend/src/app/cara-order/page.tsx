import type { Metadata } from "next"
import Link from "next/link"
import { ArrowLeft, MessageCircle } from "lucide-react"
import { howItWorks, contactWhatsApp } from "@/data/marketing"
import { Button } from "@/components/ui/button"

export const metadata: Metadata = {
  title: "Cara Order — Owndangan",
  description: "Empat langkah mudah membuat undangan pernikahan digital di Owndangan.",
}

const orderSteps = [
  ...howItWorks.slice(0, 2),
  {
    step: "03",
    title: "Pilih Paket & Bayar",
    description: "Selesaikan pembayaran dengan aman melalui Midtrans",
  },
  ...howItWorks.slice(2).map((s) => ({ ...s, step: "04" })),
]

export default function CaraOrderPage() {
  return (
    <main className="min-h-screen bg-background">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center mb-14">
          <span className="inline-block px-3 py-1 rounded-full bg-primary-container border border-primary/20 text-primary text-sm font-medium mb-4">
            Cara Order
          </span>
          <h1 className="text-3xl sm:text-4xl font-display font-bold text-foreground tracking-tight">
            Empat langkah menuju
            <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              undangan impian Anda.
            </span>
          </h1>
        </div>

        <ol className="space-y-5">
          {orderSteps.map((step) => (
            <li
              key={step.step}
              className="flex items-start gap-5 rounded-2xl border border-border bg-white p-6 shadow-sm"
            >
              <span className="shrink-0 inline-flex items-center justify-center h-11 w-11 rounded-xl bg-gradient-to-br from-primary to-secondary text-white font-display font-bold">
                {step.step}
              </span>
              <div>
                <h2 className="font-semibold text-lg text-foreground">{step.title}</h2>
                <p className="text-muted-foreground text-sm mt-1 leading-relaxed">{step.description}</p>
              </div>
            </li>
          ))}
        </ol>

        <div className="mt-12 flex flex-col sm:flex-row items-center justify-center gap-4">
          <a href={contactWhatsApp} target="_blank" rel="noopener noreferrer">
            <Button className="bg-[#25D366] hover:bg-[#1fb857] text-white">
              <MessageCircle className="mr-2 h-4 w-4" />
              Konsultasi via WhatsApp
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
    </main>
  )
}
