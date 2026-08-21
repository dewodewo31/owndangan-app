"use client"

import { useState } from "react"
import Link from "next/link"
import { AlertCircle, Mail, Lock, User, Phone, ArrowRight, Check } from "lucide-react"
import { useAuth } from "@/providers/auth-context"
import AuthLayout from "@/components/auth/auth-layout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

const passwordChecks = [
  { test: (v: string) => v.length >= 8, label: "Minimal 8 karakter" },
  { test: (v: string) => /[A-Z]/.test(v), label: "1 huruf kapital" },
  { test: (v: string) => /[0-9]/.test(v), label: "1 angka" },
]

export default function RegisterPage() {
  const { register } = useAuth()
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [phone, setPhone] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const passwordScore = passwordChecks.filter((c) => c.test(password)).length

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)
    try {
      await register({ name, email, password, phone })
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Registration failed")
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout side="register">
      <div className="text-center lg:text-left">
        <h1 className="text-3xl font-bold text-foreground tracking-tight">Buat akun gratis</h1>
        <p className="mt-2 text-muted-foreground">
          Mulai buat undangan digital impianmu dalam hitungan menit.
        </p>
      </div>

      <div className="mt-8">
        {error && (
          <div className="mb-4 flex items-start gap-2.5 p-3.5 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-sm">
            <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="name">Nama Lengkap</Label>
            <div className="relative">
              <User className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                placeholder="Nama kamu"
                className="pl-10"
                autoComplete="name"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <div className="relative">
              <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="you@example.com"
                className="pl-10"
                autoComplete="email"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="phone">No. WhatsApp (opsional)</Label>
            <div className="relative">
              <Phone className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="phone"
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="6281234567890"
                className="pl-10"
                autoComplete="tel"
              />
            </div>
            <p className="text-xs text-muted-foreground">
              Format: 62xxxxxxxxxx (diawali 62)
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Kata Sandi</Label>
            <div className="relative">
              <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={8}
                placeholder="Minimal 8 karakter"
                className="pl-10"
                autoComplete="new-password"
              />
            </div>

            {password.length > 0 && (
              <div className="space-y-1.5 pt-1">
                <div className="flex gap-1.5">
                  {[0, 1, 2].map((i) => (
                    <span
                      key={i}
                      className={`h-1 flex-1 rounded-full transition-colors ${
                        i < passwordScore
                          ? "bg-success"
                          : "bg-muted"
                      }`}
                    />
                  ))}
                </div>
                <ul className="grid grid-cols-1 gap-1">
                  {passwordChecks.map((check) => (
                    <li
                      key={check.label}
                      className={`flex items-center gap-1.5 text-xs transition-colors ${
                        check.test(password)
                          ? "text-success"
                          : "text-muted-foreground"
                      }`}
                    >
                      <Check className="h-3 w-3" />
                      {check.label}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          <Button
            type="submit"
            className="w-full h-11 bg-gradient-to-r from-primary to-secondary hover:from-primary/90 hover:to-secondary/90 border-0"
            loading={loading}
          >
            Buat Akun Gratis
            {!loading && <ArrowRight className="ml-2 h-4 w-4" />}
          </Button>

          <p className="text-xs text-muted-foreground text-center">
            Dengan mendaftar, kamu menyetujui{" "}
            <Link href="#" className="text-primary hover:underline">Syarat & Ketentuan</Link> dan{" "}
            <Link href="#" className="text-primary hover:underline">Kebijakan Privasi</Link>.
          </p>
        </form>

        <p className="mt-6 text-center text-sm text-muted-foreground">
          Sudah punya akun?{" "}
          <Link href="/login" className="text-primary font-medium hover:underline">
            Masuk
          </Link>
        </p>
      </div>
    </AuthLayout>
  )
}
