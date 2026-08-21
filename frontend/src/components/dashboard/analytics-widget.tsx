"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { Eye, Users, MessageSquare, MapPin, Phone, CalendarCheck2 } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import api from "@/lib/api"

interface Analytics {
  views: number
  unique_views: number
  whatsapp_clicks: number
  map_clicks: number
  phone_clicks: number
  rsvp_count: number
}

type WidgetState = "checking" | "upsell" | "ready"

const UPSELL_TEXT = "Aktifkan paket berbayar untuk melihat analitik undangan"

export default function AnalyticsWidget({ eventId }: { eventId: string }) {
  const [analytics, setAnalytics] = useState<Analytics | null>(null)
  const [state, setState] = useState<WidgetState>("checking")

  useEffect(() => {
    let cancelled = false
    const load = async () => {
      try {
        // Entitlement gate: analytics is a paid feature. The auth context does
        // not expose the subscription, so check it here before fetching.
        const subsRes = await api.get("/subscriptions/current")
        const pkg = subsRes.data?.data?.package
        if (!pkg || pkg.code === "free") {
          if (!cancelled) setState("upsell")
          return
        }
        const res = await api.get(`/events/${eventId}/analytics`)
        if (!cancelled) {
          setAnalytics(res.data?.data)
          setState("ready")
        }
      } catch {
        // 403/404 (free user, other user's event, endpoint missing) or any
        // network failure: hide the widget behind the upsell note, never crash.
        if (!cancelled) setState("upsell")
      }
    }
    load()
    return () => {
      cancelled = true
    }
  }, [eventId])

  if (state === "checking") return null

  if (state !== "ready" || !analytics) {
    return (
      <div className="mt-4 rounded-xl border border-dashed bg-muted/40 px-4 py-3 text-sm text-muted-foreground">
        {UPSELL_TEXT}.{" "}
        <Link href="/dashboard/billing" className="text-primary font-medium hover:underline">
          Lihat paket
        </Link>
      </div>
    )
  }

  const stats = [
    { label: "Views", value: analytics.views, icon: Eye },
    { label: "Pengunjung Unik", value: analytics.unique_views, icon: Users },
    { label: "Klik WhatsApp", value: analytics.whatsapp_clicks, icon: MessageSquare },
    { label: "Klik Peta", value: analytics.map_clicks, icon: MapPin },
    { label: "Klik Telepon", value: analytics.phone_clicks, icon: Phone },
    { label: "Konfirmasi RSVP", value: analytics.rsvp_count, icon: CalendarCheck2 },
  ]

  return (
    <div className="mt-4 grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-6 gap-3">
      {stats.map((stat) => (
        <Card key={stat.label} className="bg-muted/40">
          <CardContent className="p-4">
            <stat.icon className="h-4 w-4 text-primary" aria-hidden />
            <p className="mt-2 text-xl font-bold text-foreground font-inter tabular-nums">{stat.value}</p>
            <p className="text-xs text-muted-foreground">{stat.label}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}