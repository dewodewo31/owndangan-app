"use client"

import { useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  Users,
  Receipt,
  Package,
  LayoutTemplate,
  LogOut,
  Menu,
  X,
  ChevronDown,
  ShieldCheck,
} from "lucide-react"
import { useAuth } from "@/providers/auth-context"
import { cn } from "@/lib/utils"

const navItems = [
  { href: "/admin", label: "Overview", icon: LayoutDashboard, exact: true },
  { href: "/admin/users", label: "Pengguna", icon: Users },
  { href: "/admin/transactions", label: "Transaksi", icon: Receipt },
  { href: "/admin/packages", label: "Paket", icon: Package },
  { href: "/admin/templates", label: "Template", icon: LayoutTemplate },
]

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth()
  const pathname = usePathname()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)

  const getInitials = () => {
    const name = user?.name || "A"
    return name
      .split(" ")
      .map((n) => n[0])
      .slice(0, 2)
      .join("")
      .toUpperCase()
  }

  const SidebarContent = (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-5 h-16 border-b border-border">
        <Link href="/admin" className="flex items-center gap-2.5">
          <span className="inline-flex items-center justify-center h-8 w-8 rounded-lg bg-gradient-to-br from-plum to-rosegold text-white">
            <ShieldCheck className="h-4 w-4" />
          </span>
          <span className="font-bold text-foreground">Owndangan</span>
          <span className="text-xs font-semibold text-white bg-gradient-to-r from-plum to-rosegold px-1.5 py-0.5 rounded">
            Admin
          </span>
        </Link>
        <button className="lg:hidden p-1.5" onClick={() => setSidebarOpen(false)} aria-label="Tutup menu">
          <X className="h-5 w-5" />
        </button>
      </div>

      <nav className="flex-1 p-4 space-y-1">
        <p className="px-3 py-2 text-xs font-medium text-muted-foreground uppercase tracking-wider">
          Manajemen
        </p>
        {navItems.map((item) => {
          const isActive = item.exact
            ? pathname === item.href
            : pathname === item.href || pathname.startsWith(item.href)
          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={() => setSidebarOpen(false)}
              className={cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all",
                isActive
                  ? "bg-gradient-to-r from-plum to-rosegold text-white shadow-elevation-1"
                  : "text-muted-foreground hover:bg-accent hover:text-foreground"
              )}
            >
              <item.icon className="h-5 w-5" />
              {item.label}
            </Link>
          )
        })}
      </nav>

      <div className="p-4 border-t border-border">
        <Link
          href="/dashboard"
          className="flex items-center gap-2 px-3 py-2.5 rounded-xl text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground transition-all"
        >
          <LayoutDashboard className="h-5 w-5" />
          Ke Dashboard Pengguna
        </Link>
      </div>
    </div>
  )

  return (
    <div className="min-h-screen bg-muted/40 lg:flex">
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-50 w-64 bg-background border-r border-border transition-transform duration-300 lg:translate-x-0 lg:static lg:block lg:shrink-0",
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        {SidebarContent}
      </aside>

      <div className="flex-1 flex flex-col min-w-0 min-h-screen">
        <header className="sticky top-0 z-30 bg-background/80 backdrop-blur-xl border-b border-border h-16 flex items-center px-4 sm:px-6">
          <button
            className="lg:hidden p-2 text-foreground"
            onClick={() => setSidebarOpen(true)}
            aria-label="Buka menu"
          >
            <Menu className="h-5 w-5" />
          </button>

          <div className="flex-1" />

          <div className="relative">
            <button
              onClick={() => setUserMenuOpen(!userMenuOpen)}
              className="flex items-center gap-2.5 px-2 py-1.5 rounded-xl hover:bg-accent transition-colors"
            >
              <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-gradient-to-br from-plum to-rosegold text-white text-xs font-bold">
                {getInitials()}
              </span>
              <span className="hidden sm:block text-sm font-medium">{user?.name}</span>
              <ChevronDown className="hidden sm:block h-4 w-4 text-muted-foreground" />
            </button>

            {userMenuOpen && (
              <>
                <div className="fixed inset-0 z-40" onClick={() => setUserMenuOpen(false)} />
                <div className="absolute right-0 mt-2 w-56 rounded-xl bg-card border border-border shadow-xl shadow-black/5 z-50 overflow-hidden">
                  <div className="px-4 py-3 border-b border-border">
                    <p className="text-sm font-semibold truncate">{user?.name}</p>
                    <p className="text-xs text-muted-foreground truncate">{user?.email}</p>
                    <p className="mt-1 text-xs font-medium text-primary">Administrator</p>
                  </div>
                  <button
                    onClick={logout}
                    className="w-full flex items-center gap-2.5 px-4 py-3 text-sm text-destructive hover:bg-destructive/5 transition-colors"
                  >
                    <LogOut className="h-4 w-4" />
                    Keluar
                  </button>
                </div>
              </>
            )}
          </div>
        </header>

        <main className="flex-1 p-4 sm:p-6 lg:p-8">
          <div className="max-w-7xl mx-auto">{children}</div>
        </main>
      </div>
    </div>
  )
}
