"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import {
  Heart,
  Users,
  MessageSquare,
  MessagesSquare,
  PenSquare,
  CreditCard,
  ArrowRight,
  Sparkles,
} from "lucide-react"
import ProtectedRoute from "@/components/dashboard/protected-route"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Badge } from "@/components/ui/badge"
import api from "@/lib/api"
import { cn } from "@/lib/utils"

interface DashboardStats {
  eventCount: number
  guestCount: number
  rsvpCount: number
  pendingMessages: number
  subscription: string
  paymentStatus: string
}

const quickActions = [
  { href: "/dashboard/editor", label: "Edit Undangan", icon: PenSquare, color: "from-indigo-500 to-violet-500" },
  { href: "/dashboard/guests", label: "Kelola Tamu", icon: Users, color: "from-emerald-500 to-teal-500" },
  { href: "/dashboard/billing", label: "Upgrade Paket", icon: CreditCard, color: "from-fuchsia-500 to-pink-500" },
]

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      const [eventsRes, subsRes] = await Promise.all([
        api.get("/events"),
        api.get("/subscriptions/current"),
      ])

      setStats({
        eventCount: eventsRes.data?.data?.length || 0,
        guestCount: 0,
        rsvpCount: 0,
        pendingMessages: 0,
        subscription: subsRes.data?.data?.package?.name || "Free",
        paymentStatus: subsRes.data?.data?.status || "active",
      })
    } catch (error) {
      console.error("Failed to fetch stats:", error)
    } finally {
      setLoading(false)
    }
  }

  const statCards = [
    { title: "Undangan", value: stats?.eventCount || 0, icon: Heart, color: "text-rose-500 bg-rose-500/10" },
    { title: "Tamu", value: stats?.guestCount || 0, icon: Users, color: "text-indigo-500 bg-indigo-500/10" },
    { title: "RSVP", value: stats?.rsvpCount || 0, icon: MessageSquare, color: "text-emerald-500 bg-emerald-500/10" },
    { title: "Pesan Masuk", value: stats?.pendingMessages || 0, icon: MessagesSquare, color: "text-violet-500 bg-violet-500/10" },
  ]

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Dashboard</h1>
              <p className="mt-1 text-muted-foreground">
                Kelola undangan, tamu, dan RSVP dalam satu tempat.
              </p>
            </div>
            <Link href="/dashboard/editor">
              <Button className="bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 border-0">
                <Sparkles className="h-4 w-4 mr-2" />
                Buat Undangan Baru
              </Button>
            </Link>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
              {[1, 2, 3, 4].map((i) => (
                <Card key={i}>
                  <CardContent className="p-6">
                    <Skeleton className="h-4 w-1/2 mb-3" />
                    <Skeleton className="h-8 w-1/4" />
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
              {statCards.map((stat) => (
                <Card
                  key={stat.title}
                  className="transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-indigo-500/5"
                >
                  <CardContent className="p-6">
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="text-sm text-muted-foreground">{stat.title}</p>
                        <p className="mt-2 text-3xl font-bold text-foreground">{stat.value}</p>
                      </div>
                      <span className={cn("inline-flex items-center justify-center h-11 w-11 rounded-xl", stat.color)}>
                        <stat.icon className="h-5 w-5" />
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Subscription */}
            <Card className="lg:col-span-1">
              <CardContent className="p-6">
                <h2 className="font-semibold text-foreground mb-4">Langganan</h2>
                <div className="rounded-2xl bg-gradient-to-br from-indigo-600 to-violet-600 p-5 text-white">
                  <p className="text-xs text-indigo-200 font-medium">Paket Aktif</p>
                  <div className="mt-1 flex items-center gap-2">
                    <p className="text-2xl font-bold">{stats?.subscription || "Free"}</p>
                  </div>
                  <div className="mt-3 inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-white/15 text-xs font-medium">
                    <span className="h-1.5 w-1.5 rounded-full bg-emerald-400 animate-pulse" />
                    {stats?.paymentStatus || "active"}
                  </div>
                </div>
                <Link href="/dashboard/billing" className="mt-4 block">
                  <Button variant="outline" className="w-full">
                    Kelola Langganan
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card className="lg:col-span-2">
              <CardContent className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="font-semibold text-foreground">Aksi Cepat</h2>
                  <Link href="/dashboard/editor" className="text-xs text-primary hover:underline">
                    Lihat semua
                  </Link>
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  {quickActions.map((action) => (
                    <Link
                      key={action.href}
                      href={action.href}
                      className="group flex items-center gap-3 rounded-xl border border-border p-4 transition-all hover:-translate-y-0.5 hover:shadow-md hover:border-primary/40"
                    >
                      <span
                        className={cn(
                          "inline-flex items-center justify-center h-10 w-10 rounded-lg bg-gradient-to-br text-white shrink-0",
                          action.color
                        )}
                      >
                        <action.icon className="h-5 w-5" />
                      </span>
                      <span className="text-sm font-medium group-hover:text-primary transition-colors">
                        {action.label}
                      </span>
                    </Link>
                  ))}
                </div>

                <div className="mt-6 rounded-xl bg-muted/60 border border-border p-4">
                  <div className="flex items-start gap-3">
                    <Badge variant="secondary" className="shrink-0 mt-0.5">
                      Tips
                    </Badge>
                    <p className="text-sm text-muted-foreground">
                      Bagikan link undangan ke tamu melalui WhatsApp agar konfirmasi RSVP
                      masuk secara otomatis.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}
