"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import {
  Heart,
  Users,
  MessageSquare,
  Eye,
  PenSquare,
  ArrowRight,
  Sparkles,
  ExternalLink,
  Wand2,
  CalendarDays,
  Sun,
} from "lucide-react"
import ProtectedRoute from "@/components/dashboard/protected-route"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import AnalyticsWidget from "@/components/dashboard/analytics-widget"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Badge } from "@/components/ui/badge"
import { useAuth } from "@/providers/auth-context"
import api from "@/lib/api"
import { cn, formatDate } from "@/lib/utils"
import type { WeddingEvent } from "@/lib/types"

interface DashboardStats {
  eventCount: number
  guestCount: number
  rsvpCount: number
  viewCount: number
  subscription: string
}

function greeting(): string {
  const h = new Date().getHours()
  if (h < 12) return "Selamat pagi"
  if (h < 15) return "Selamat siang"
  if (h < 18) return "Selamat sore"
  return "Selamat malam"
}

function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return "baru saja"
  if (mins < 60) return `${mins} menit lalu`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} jam lalu`
  return `${Math.floor(hours / 24)} hari lalu`
}

export default function DashboardPage() {
  const { user } = useAuth()
  const [events, setEvents] = useState<WeddingEvent[]>([])
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
      const list: WeddingEvent[] = eventsRes.data?.data || []
      setEvents(list)
      setStats({
        eventCount: list.length,
        guestCount: list.reduce((sum, e) => sum + (e.guest_count || 0), 0),
        rsvpCount: list.reduce((sum, e) => sum + (e.rsvp_count || 0), 0),
        viewCount: list.reduce((sum, e) => sum + (e.view_count || 0), 0),
        subscription: subsRes.data?.data?.package?.name || "",
      })
    } catch (error) {
      console.error("Failed to fetch stats:", error)
    } finally {
      setLoading(false)
    }
  }

  const firstName = user?.name?.split(" ")[0] || "Teman"

  const statCards = [
    { title: "Undangan", value: stats?.eventCount || 0, icon: Heart, iconClass: "bg-primary-container text-primary" },
    { title: "Tamu", value: stats?.guestCount || 0, icon: Users, iconClass: "bg-secondary-container text-secondary" },
    { title: "RSVP", value: stats?.rsvpCount || 0, icon: MessageSquare, iconClass: "bg-tertiary-container text-tertiary" },
    { title: "Dilihat", value: stats?.viewCount || 0, icon: Eye, iconClass: "bg-surface-container-high text-on-surface-variant" },
  ]

  const activity = [...events]
    .sort((a, b) => new Date(b.updated_at || "").getTime() - new Date(a.updated_at || "").getTime())
    .slice(0, 5)
    .filter((e) => e.updated_at)
    .map((e) => ({
      event: e,
      label: e.status === "published" ? "diperbarui" : "masih disusun",
      time: relativeTime(e.updated_at || ""),
    }))

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          {/* Greeting */}
          <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
            <div>
              <p className="text-sm text-on-surface-variant">{new Date().toLocaleDateString("id-ID", { weekday: "long", day: "numeric", month: "long", year: "numeric" })}</p>
              <h1 className="mt-1 text-2xl sm:text-3xl font-bold text-foreground tracking-tight">
                {greeting()}, <span className="text-primary">{firstName}</span> <Sun className="ml-1 inline h-6 w-6 text-amber-400 align-middle" aria-hidden />
              </h1>
              <p className="mt-1 text-muted-foreground">
                Ini rangkuman undangan dan aktivitas pernikahanmu.
              </p>
            </div>
            <Link href="/dashboard/editor">
              <Button className="bg-gradient-to-r from-plum to-rosegold hover:from-plum/90 hover:to-rosegold/90 border-0 text-white shadow-elevation-1">
                <Wand2 className="h-4 w-4 mr-2" />
                Buat Undangan Baru
              </Button>
            </Link>
          </div>

          {/* Stat cards */}
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
                <Card key={stat.title} className="transition-all duration-300 hover:-translate-y-1 hover:shadow-elevation-2">
                  <CardContent className="p-6">
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="text-sm text-muted-foreground">{stat.title}</p>
                        <p className="mt-2 text-3xl font-bold text-foreground font-inter tabular-nums">{stat.value}</p>
                      </div>
                      <span className={cn("inline-flex items-center justify-center h-11 w-11 rounded-xl", stat.iconClass)}>
                        <stat.icon className="h-5 w-5" />
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
            {/* Invitations */}
            <div className="lg:col-span-2 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="font-semibold text-foreground">Undangan Kamu</h2>
                <Link href="/dashboard/editor" className="text-xs text-primary hover:underline">
                  Lihat semua
                </Link>
              </div>

              {loading ? (
                <div className="space-y-4">
                  {[1, 2].map((i) => (
                    <Skeleton key={i} className="h-28 rounded-2xl" />
                  ))}
                </div>
              ) : events.length === 0 ? (
                <Card>
                  <CardContent className="flex flex-col items-center text-center py-14">
                    <div className="w-14 h-14 rounded-full bg-primary-container flex items-center justify-center mb-4">
                      <Heart className="h-6 w-6 text-primary" />
                    </div>
                    <p className="font-medium text-foreground">Belum ada undangan</p>
                    <p className="text-sm text-muted-foreground mt-1 max-w-sm">
                      Mulai undangan pertamamu dan bagikan ke orang-orang tersayang.
                    </p>
                    <Link href="/dashboard/editor" className="mt-6">
                      <Button>
                        <Sparkles className="h-4 w-4 mr-2" />
                        Buat Undangan
                      </Button>
                    </Link>
                  </CardContent>
                </Card>
              ) : (
                events.map((event) => (
                  <Card key={event.id} className="transition-all duration-300 hover:-translate-y-0.5 hover:shadow-elevation-2">
                    <CardContent className="p-5 sm:p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                      <div className="flex items-center gap-4 min-w-0">
                        <span className="inline-flex items-center justify-center h-12 w-12 rounded-xl bg-gradient-to-br from-primary-container to-secondary-container text-primary shrink-0">
                          <Heart className="h-5 w-5" />
                        </span>
                        <div className="min-w-0">
                          <div className="flex items-center gap-2 flex-wrap">
                            <h3 className="font-semibold text-foreground truncate">{event.couple_name || event.title}</h3>
                            <Badge variant={event.status === "published" ? "success" : "secondary"}>
                              {event.status === "published" ? "Diterbitkan" : "Draft"}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground mt-1 flex items-center gap-1.5">
                            {event.wedding_date ? (
                              <>
                                <CalendarDays className="h-3.5 w-3.5" />
                                {formatDate(event.wedding_date)}
                              </>
                            ) : (
                              "Tanggal belum diatur"
                            )}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2 shrink-0">
                        <Link href="/dashboard/editor">
                          <Button size="sm" variant="outline">
                            <PenSquare className="h-3.5 w-3.5 mr-1.5" />
                            Edit
                          </Button>
                        </Link>
                        {event.status === "published" && event.slug ? (
                          <Link href={`/${event.slug}`} target="_blank" rel="noreferrer">
                            <Button size="sm" variant="ghost">
                              <ExternalLink className="h-3.5 w-3.5 mr-1.5" />
                              Pratinjau
                            </Button>
                          </Link>
                        ) : (
                          <Link href="/dashboard/editor">
                            <Button size="sm" variant="ghost">
                              <ArrowRight className="h-3.5 w-3.5 mr-1.5" />
                              Publikasikan
                            </Button>
                          </Link>
                        )}
                      </div>
                    </CardContent>
                    <div className="px-5 sm:px-6 pb-5 sm:pb-6">
                      <AnalyticsWidget eventId={event.id} />
                    </div>
                  </Card>
                ))
              )}
            </div>

            {/* Activity + subscription */}
            <div className="space-y-6">
              <Card>
                <CardContent className="p-6">
                  <h2 className="font-semibold text-foreground mb-4">Aktivitas Terbaru</h2>
                  {activity.length === 0 ? (
                    <p className="text-sm text-muted-foreground">Belum ada aktivitas.</p>
                  ) : (
                    <ul className="space-y-4">
                      {activity.map((item) => (
                        <li key={item.event.id} className="flex items-start gap-3">
                          <span className="mt-1 h-2 w-2 rounded-full bg-primary shrink-0" />
                          <div className="min-w-0">
                            <p className="text-sm text-foreground">
                              Undangan <span className="font-medium">“{item.event.couple_name || item.event.title}”</span> {item.label}
                            </p>
                            <p className="text-xs text-muted-foreground mt-0.5">{item.time}</p>
                          </div>
                        </li>
                      ))}
                    </ul>
                  )}
                </CardContent>
              </Card>

              <Card className="bg-gradient-to-br from-plum to-rosegold text-white border-0">
                <CardContent className="p-6">
                  <p className="text-xs font-medium text-white/70">Paket Aktif</p>
                  <p className="mt-1 text-2xl font-bold">{stats?.subscription || "Free"}</p>
                  <p className="mt-2 text-sm text-white/80 leading-relaxed">
                    Kelola langganan, fitur, dan metode pembayaran dari satu tempat.
                  </p>
                  <Link href="/dashboard/billing" className="mt-4 block">
                    <Button className="w-full bg-white text-plum hover:bg-white/90 border-0">
                      Kelola Langganan
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </Link>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}