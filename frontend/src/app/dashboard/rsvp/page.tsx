"use client"

import { useEffect, useState } from "react"
import { MessageSquare, CheckCircle2, XCircle, Users, CalendarCheck } from "lucide-react"
import ProtectedRoute from "@/components/dashboard/protected-route"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import api from "@/lib/api"
import { cn } from "@/lib/utils"

interface RSVPRecap {
  total_responded: number
  attending: number
  not_attending: number
  total_guest_count: number
}

interface RSVP {
  id: string
  guest_id: string
  attendance: string
  guest_count: number
  message: string
  submitted_at: string
}

export default function RSVPPage() {
  const [recap, setRecap] = useState<RSVPRecap | null>(null)
  const [rsvps, setRsvps] = useState<RSVP[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchRSVP()
  }, [])

  const fetchRSVP = async () => {
    try {
      const eventsRes = await api.get('/events')
      const events = eventsRes.data?.data || []
      if (events.length > 0) {
        const eventID = events[0].id
        const [recapRes, rsvpsRes] = await Promise.all([
          api.get(`/rsvp/${eventID}/recap`),
          api.get(`/rsvp/${eventID}`),
        ])
        setRecap(recapRes.data?.data)
        setRsvps(rsvpsRes.data?.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch RSVP:', error)
    } finally {
      setLoading(false)
    }
  }

  const recapCards = [
    { title: "Total Respon", value: recap?.total_responded || 0, icon: CalendarCheck, color: "text-indigo-500 bg-indigo-500/10" },
    { title: "Hadir", value: recap?.attending || 0, icon: CheckCircle2, color: "text-emerald-500 bg-emerald-500/10" },
    { title: "Tidak Hadir", value: recap?.not_attending || 0, icon: XCircle, color: "text-rose-500 bg-rose-500/10" },
    { title: "Total Tamu", value: recap?.total_guest_count || 0, icon: Users, color: "text-violet-500 bg-violet-500/10" },
  ]

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">RSVP</h1>
            <p className="mt-1 text-muted-foreground">
              Pantau konfirmasi kehadiran tamu secara real-time.
            </p>
          </div>

          {loading ? (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-4 gap-5">
                {[1, 2, 3, 4].map((i) => (
                  <Skeleton key={i} className="h-28 rounded-2xl" />
                ))}
              </div>
              <Skeleton className="h-64 rounded-2xl" />
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
                {recapCards.map((stat) => (
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

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <MessageSquare className="h-5 w-5 text-primary" />
                    Daftar Respon Tamu
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {rsvps.length === 0 ? (
                    <div className="text-center py-12">
                      <div className="mx-auto w-14 h-14 rounded-full bg-muted flex items-center justify-center mb-4">
                        <MessageSquare className="h-6 w-6 text-muted-foreground" />
                      </div>
                      <p className="text-muted-foreground">Belum ada respon RSVP.</p>
                    </div>
                  ) : (
                    <Table>
                      <TableHeader>
                        <TableRow className="hover:bg-transparent">
                          <TableHead>Tamu</TableHead>
                          <TableHead>Kehadiran</TableHead>
                          <TableHead>Jumlah</TableHead>
                          <TableHead>Pesan</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {rsvps.map((rsvp) => (
                          <TableRow key={rsvp.id}>
                            <TableCell className="font-medium">{rsvp.guest_id}</TableCell>
                            <TableCell>
                              <Badge variant={rsvp.attendance === 'attending' ? 'success' : 'destructive'}>
                                {rsvp.attendance === 'attending' ? 'Hadir' : 'Tidak Hadir'}
                              </Badge>
                            </TableCell>
                            <TableCell>{rsvp.guest_count}</TableCell>
                            <TableCell className="max-w-xs truncate text-muted-foreground">
                              {rsvp.message || "-"}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  )}
                </CardContent>
              </Card>
            </>
          )}
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}
