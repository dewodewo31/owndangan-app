"use client"

import Link from "next/link"
import { Heart } from "lucide-react"

const productLinks = [
  { href: "#features", label: "Fitur" },
  { href: "#templates", label: "Template" },
  { href: "#pricing", label: "Harga" },
  { href: "/cara-order", label: "Cara Order" },
  { href: "/testimoni", label: "Testimoni" },
  { href: "/faq", label: "FAQ" },
]

const supportLinks = [
  { href: "#", label: "Bantuan" },
  { href: "#", label: "Kontak" },
  { href: "#", label: "Pusat Bantuan" },
]

const legalLinks = [
  { href: "#", label: "Privacy Policy" },
  { href: "#", label: "Terms of Service" },
]

export function Footer() {
  return (
    <footer className="bg-slate-950 text-slate-300">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-10">
          <div className="col-span-2">
            <Link href="/" className="flex items-center gap-2.5">
              <span className="inline-flex items-center justify-center h-9 w-9 rounded-xl bg-gradient-to-br from-plum to-rosegold text-white shadow-md shadow-primary/30">
                <svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
              </span>
              <span className="text-lg font-bold text-white">Owndangan</span>
            </Link>
            <p className="mt-4 text-sm text-slate-400 max-w-xs leading-relaxed">
              Platform undangan pernikahan digital modern. Buat, kelola, dan bagikan
              undangan dengan mudah dan elegan.
            </p>
            <div className="mt-6 flex gap-3">
              {["Instagram", "TikTok", "WhatsApp"].map((social) => (
                <a
                  key={social}
                  href="#"
                  className="inline-flex items-center justify-center h-9 px-3.5 rounded-lg bg-slate-800/60 hover:bg-slate-700 text-xs font-medium text-slate-300 hover:text-white transition-colors"
                >
                  {social}
                </a>
              ))}
            </div>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-4">Produk</h4>
            <ul className="space-y-2.5 text-sm">
              {productLinks.map((link) => (
                <li key={link.label}>
                  <a href={link.href} className="hover:text-white transition-colors">
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-4">Dukungan</h4>
            <ul className="space-y-2.5 text-sm">
              {supportLinks.map((link) => (
                <li key={link.label}>
                  <a href={link.href} className="hover:text-white transition-colors">
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-4">Legal</h4>
            <ul className="space-y-2.5 text-sm">
              {legalLinks.map((link) => (
                <li key={link.label}>
                  <a href={link.href} className="hover:text-white transition-colors">
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="mt-12 pt-8 border-t border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-slate-500">
          <p>© 2026 Owndangan. All rights reserved.</p>
          <p>
            Dibuat dengan <Heart className="inline h-4 w-4 fill-rose-400 text-rose-400 align-[-2px]" aria-hidden /> untuk pasangan Indonesia
          </p>
        </div>
      </div>
    </footer>
  )
}
