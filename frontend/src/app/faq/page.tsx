import type { Metadata } from "next"
import Link from "next/link"
import { ArrowLeft } from "lucide-react"
import { FaqList } from "@/components/marketing/faq-list"
import { Button } from "@/components/ui/button"

export const metadata: Metadata = {
  title: "FAQ — Owndangan",
  description: "Pertanyaan yang sering diajukan tentang Owndangan, platform undangan pernikahan digital.",
}

export default function FAQPage() {
  return (
    <main className="min-h-screen bg-background">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center mb-12">
          <span className="inline-block px-3 py-1 rounded-full bg-primary-container border border-primary/20 text-primary text-sm font-medium mb-4">
            FAQ
          </span>
          <h1 className="text-3xl sm:text-4xl font-display font-bold text-foreground tracking-tight">
            Pertanyaan yang
            <span className="block bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              sering diajukan.
            </span>
          </h1>
        </div>

        <FaqList />

        <div className="mt-12 text-center">
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
