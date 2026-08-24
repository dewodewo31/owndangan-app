"use client"

import Link from "next/link"
import { Heart, Sparkles, Users, MessageSquare, Star } from "lucide-react"
import { cn } from "@/lib/utils"

interface AuthLayoutProps {
  children: React.ReactNode
  side?: "login" | "register"
  className?: string
}

const highlights = [
  {
    icon: Heart,
    title: "Undangan elegan",
    description: "Template modern yang mencerminkan cerita kalian",
  },
  {
    icon: Users,
    title: "Kelola tamu mudah",
    description: "Impor dan kelola daftar tamu dalam satu tempat",
  },
  {
    icon: MessageSquare,
    title: "RSVP real-time",
    description: "Terima konfirmasi kehadiran secara otomatis",
  },
]

export default function AuthLayout({ children, side = "login", className }: AuthLayoutProps) {
  return (
    <div className="min-h-screen grid lg:grid-cols-2 bg-background">
      {/* Left: Brand panel */}
      <div className="hidden lg:flex flex-col justify-between bg-gradient-to-br from-plum via-rosegold/70 to-rosegold p-12 text-white relative overflow-hidden">
        <div className="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-72 w-72 rounded-full border border-white/20" />

        <div className="relative">
          <Link href="/" className="inline-flex items-center gap-2.5">
            <span className="inline-flex items-center justify-center h-10 w-10 rounded-xl bg-white/15 backdrop-blur">
              <svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
            </span>
            <span className="text-xl font-bold">Owndangan</span>
          </Link>
        </div>

        <div className="relative max-w-md">
          <span className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/15 text-sm font-medium mb-6 backdrop-blur">
            <Sparkles className="h-4 w-4" />
            {side === "register" ? "Mulai gratis, tanpa kartu kredit" : "Selamat datang kembali"}
          </span>
          <h2 className="text-4xl font-bold leading-tight tracking-tight">
            {side === "register" ? (
              <>
                Mulai buat undangan
                <span className="block text-primary-foreground/70">impian kalian.</span>
              </>
            ) : (
              <>
                Undangan digital
                <span className="block text-primary-foreground/70">yang tak terlupakan.</span>
              </>
            )}
          </h2>
          <div className="mt-10 space-y-6">
            {highlights.map((item) => (
              <div key={item.title} className="flex items-start gap-4">
                <span className="inline-flex items-center justify-center h-10 w-10 rounded-xl bg-white/15 backdrop-blur shrink-0">
                  <item.icon className="h-5 w-5" />
                </span>
                <div>
                  <p className="font-semibold">{item.title}</p>
                  <p className="text-sm text-primary-foreground/80 mt-0.5">{item.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="relative">
          <div className="flex items-center gap-2 text-sm text-primary-foreground/80">
            <span className="flex gap-1">
              {Array.from({ length: 5 }).map((_, i) => (
                <Star key={i} className="h-4 w-4 fill-current" aria-hidden />
              ))}
            </span>
            <span>Dipercaya ribuan pasangan di Indonesia</span>
          </div>
        </div>
      </div>

      {/* Right: Form */}
      <div className={cn("flex items-center justify-center px-4 sm:px-6 lg:px-12 py-12", className)}>
        <div className="w-full max-w-md">
          <Link href="/" className="lg:hidden inline-flex items-center gap-2 mb-8">
            <span className="inline-flex items-center justify-center h-9 w-9 rounded-xl bg-gradient-to-br from-plum to-rosegold text-white">
              <svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
            </span>
            <span className="text-lg font-bold">Owndangan</span>
          </Link>
          {children}
        </div>
      </div>
    </div>
  )
}
