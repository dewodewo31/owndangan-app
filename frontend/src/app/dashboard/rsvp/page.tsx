"use client"

import { useEffect, useState } from "react"
import { MessageSquare, CheckCircle2, XCircle, Users, CalendarCheck, Download, HelpCircle, FileSpreadsheet, AlertCircle } from "lucide-react"
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
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface RSVPRecap {
  total_responded: number
  attending: number
  not_attending: number
  maybe: number
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
  const [eventID, setEventID] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [exporting, setExporting] = useState(false)
  const [exportError, setExportError] = useState<string | null>(null)

  useEffect(() => {
    fetchRSVP()
  }, [])

  const fetchRSVP = async () => {
    try {
      const eventsRes = await api.get('/events')
      const events = eventsRes.data?.data || []
      if (events.length > 0) {
        const id = events[0].id
        setEventID(id)
        const [recapRes, rsvpsRes] = await Promise.all([
          api.get(`/rsvp/${id}/recap`),
          api.get(`/rsvp/${id}`),
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

  const downloadExport = async (format: "csv" | "xlsx") => {
    if (!eventID) return
    setExporting(true)
    setExportError(null)
    try {
      const token = typeof window !== "undefined" ? localStorage.getItem("access_token") : null
      const baseURL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"
      const res = await fetch(`${baseURL}/events/${eventID}/rsvp/export?format=${format}`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      })
      if (!res.ok) {
        if (res.status === 403) {
          throw new Error("Ekspor XLSX hanya untuk pengguna Pro. Upgrade paket untuk mengunduh XLSX.")
        }
        throw new Error("Gagal mengekspor data RSVP")
      }
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement("a")
      a.href = url
      a.download = `rsvp-${eventID}.${format}`
      a.click()
      URL.revokeObjectURL(url)
    } catch (error) {
      console.error("Failed to export RSVP:", error)
      setExportError(error instanceof Error ? error.message : "Gagal mengekspor data RSVP")
    } finally {
      setExporting(false)
    }
  }

  const recapCards = [
    { title: "Total Respon", value: recap?.total_responded || 0, icon: CalendarCheck, color: "text-primary bg-primary-container" },
    { title: "Hadir", value: recap?.attending || 0, icon: CheckCircle2, color: "text-success bg-success/10" },
    { title: "Tidak Hadir", value: recap?.not_attending || 0, icon: XCircle, color: "text-secondary bg-secondary-container" },
    { title: "Ragu", value: recap?.maybe || 0, icon: HelpCircle, color: "text-warning bg-warning/10" },
    { title: "Total Tamu", value: recap?.total_guest_count || 0, icon: Users, color: "text-primary bg-secondary-container" },
  ]

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">RSVP</h1>
              <p className="mt-1 text-muted-foreground">
                Pantau konfirmasi kehadiran tamu secara real-time.
              </p>
            </div>
            <div className="flex flex-col items-start sm:items-end gap-2">
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => downloadExport("csv")}
                  disabled={!eventID || exporting}
                >
                  <Download className="h-4 w-4 mr-2" /> Ekspor CSV
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => downloadExport("xlsx")}
                  disabled={!eventID || exporting}
                >
                  <FileSpreadsheet className="h-4 w-4 mr-2" /> Ekspor XLSX
                  <Badge variant="secondary" className="ml-2 border-0">Pro</Badge>
                </Button>
              </div>
              {exportError && (
                <p className="flex items-center gap-1.5 text-xs text-destructive max-w-xs text-left">
                  <AlertCircle className="h-3.5 w-3.5 shrink-0" />
                  {exportError}
                </p>
              )}
            </div>
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
                    className="transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-primary/5"
                  >
                    <CardContent className="p-6">
                      <div className="flex items-start justify-between">
                        <div>
                          <p className="text-sm text-muted-foreground">{stat.title}</p>
                          <p className="mt-2 text-3xl font-bold text-foreground font-inter tabular-nums">{stat.value}</p>
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
                              <Badge
                                variant={rsvp.attendance === 'attending' ? 'success' : rsvp.attendance === 'maybe' ? 'secondary' : 'destructive'}
                              >
                                {rsvp.attendance === 'attending' ? 'Hadir' : rsvp.attendance === 'maybe' ? 'Ragu' : 'Tidak Hadir'}
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
