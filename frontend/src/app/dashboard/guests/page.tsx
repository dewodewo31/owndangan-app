"use client"

import { useEffect, useState } from "react"
import { Plus, MessageCircle, Trash2, Users, Download, X, RotateCcw, QrCode, AlertCircle, CheckCircle2 } from "lucide-react"
import ProtectedRoute from "@/components/dashboard/protected-route"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Label } from "@/components/ui/label"
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
import { GuestImportDialog } from "@/components/dashboard/guest-import-dialog"
import api from "@/lib/api"
import { cn } from "@/lib/utils"

interface Guest {
  id: string
  name: string
  phone: string
  category: string
  note: string
  token: string
}

const categoryLabels: Record<string, string> = {
  family: "Keluarga",
  friend: "Teman",
  colleague: "Kolega",
  other: "Lainnya",
}

const categoryColors: Record<string, string> = {
  family: "bg-rose-500/10 text-rose-600",
  friend: "bg-primary-container text-primary",
  colleague: "bg-emerald-500/10 text-emerald-600",
  other: "bg-muted text-muted-foreground",
}

export default function GuestsPage() {
  const [guests, setGuests] = useState<Guest[]>([])
  const [deletedGuests, setDeletedGuests] = useState<Guest[]>([])
  const [loading, setLoading] = useState(true)
  const [deletedLoading, setDeletedLoading] = useState(false)
  const [view, setView] = useState<'active' | 'trash'>('active')
  const [showAddModal, setShowAddModal] = useState(false)
  const [showImportModal, setShowImportModal] = useState(false)
  const [eventId, setEventId] = useState<string | null>(null)
  const [newGuest, setNewGuest] = useState({ name: '', phone: '', category: 'family', note: '' })
  const [checkInToken, setCheckInToken] = useState('')
  const [checkingIn, setCheckingIn] = useState(false)
  const [checkInError, setCheckInError] = useState<string | null>(null)
  const [checkInResult, setCheckInResult] = useState<{ token: string; attendedAt: string } | null>(null)
  const [attendedMap, setAttendedMap] = useState<Record<string, string>>({})

  useEffect(() => {
    fetchGuests()
  }, [])

  const fetchGuests = async () => {
    try {
      const eventsRes = await api.get('/events')
      const events = eventsRes.data?.data || []
      if (events.length > 0) {
        setEventId(events[0].id)
        const guestsRes = await api.get(`/events/${events[0].id}/guests`)
        setGuests(guestsRes.data?.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch guests:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchDeletedGuests = async () => {
    try {
      setDeletedLoading(true)
      if (!eventId) return
      const deletedRes = await api.get(`/events/${eventId}/guests/deleted`)
      setDeletedGuests(deletedRes.data?.data || [])
    } catch (error) {
      console.error('Failed to fetch deleted guests:', error)
    } finally {
      setDeletedLoading(false)
    }
  }

  const handleRestoreGuest = async (guestID: string) => {
    try {
      if (!eventId) return
      await api.post(`/events/${eventId}/guests/${guestID}/restore`)
      alert('Tamu berhasil dipulihkan')
      fetchGuests()
      fetchDeletedGuests()
    } catch (error) {
      alert('Gagal memulihkan tamu')
    }
  }

  const handleAddGuest = async () => {
    try {
      if (!eventId) return
      await api.post(`/events/${eventId}/guests`, newGuest)
      setShowAddModal(false)
      setNewGuest({ name: '', phone: '', category: 'family', note: '' })
      fetchGuests()
    } catch (error) {
      alert('Gagal menambahkan tamu')
    }
  }

  const handleDeleteGuest = async (guestID: string) => {
    if (!confirm('Yakin ingin menghapus tamu ini?')) return
    try {
      if (!eventId) return
      await api.delete(`/events/${eventId}/guests/${guestID}`)
      fetchGuests()
    } catch (error) {
      alert('Gagal menghapus tamu')
    }
  }

  const generateWhatsAppLink = (guest: Guest) => {
    const phone = guest.phone ?? ''
    const message = `Halo ${guest.name}, kami mengundangmu untuk menghadiri pernikahan kami! Lihat undangan di: ${window.location.origin}/${guest.token}`
    return `https://wa.me/${phone.replace(/^0/, '62')}?text=${encodeURIComponent(message)}`
  }

  const checkIn = async (token: string) => {
    const clean = token.trim()
    if (!clean) return
    setCheckingIn(true)
    setCheckInError(null)
    try {
      const baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
      const res = await fetch(`${baseURL}/guests/check-in`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token: clean }),
      })
      if (!res.ok) {
        if (res.status === 403) {
          throw new Error('Check-in hanya tersedia untuk pengguna Pro.')
        }
        const body = await res.json().catch(() => null)
        throw new Error(body?.error?.message || 'Gagal melakukan check-in.')
      }
      const body = await res.json()
      const attendedAt = body?.data?.attended_at
      setCheckInResult({ token: clean, attendedAt })
      const g = guests.find((x) => x.token === clean)
      if (g) setAttendedMap((m) => ({ ...m, [g.id]: attendedAt }))
    } catch (error) {
      setCheckInError(error instanceof Error ? error.message : 'Gagal melakukan check-in.')
    } finally {
      setCheckingIn(false)
    }
  }

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Daftar Tamu</h1>
              <p className="mt-1 text-muted-foreground">
                Kelola dan bagikan undangan ke tamu kamu.
              </p>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setShowImportModal(true)} disabled={!eventId}>
                <Download className="h-4 w-4 mr-2" />
                Impor CSV
              </Button>
              <Button
                onClick={() => setShowAddModal(true)}
                className="bg-gradient-to-r from-primary to-secondary hover:from-primary/90 hover:to-secondary/90 border-0"
              >
                <Plus className="h-4 w-4 mr-2" />
                Tambah Tamu
              </Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-4 sm:p-6">
              <div className="flex flex-col sm:flex-row sm:items-end gap-3">
                <div className="flex-1 space-y-2">
                  <Label htmlFor="checkin-token">Token Tamu</Label>
                  <Input
                    id="checkin-token"
                    type="text"
                    placeholder="Tempel token undangan tamu (mis. hasil pindai QR)"
                    value={checkInToken}
                    onChange={(e) => setCheckInToken(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') checkIn(checkInToken)
                    }}
                  />
                </div>
                <Button
                  onClick={() => checkIn(checkInToken)}
                  loading={checkingIn}
                  disabled={!checkInToken.trim()}
                  className="bg-gradient-to-r from-primary to-secondary hover:from-primary/90 hover:to-secondary/90 border-0"
                >
                  <QrCode className="h-4 w-4 mr-2" />
                  Check-in
                </Button>
              </div>
              {checkInError && (
                <div className="mt-3 flex items-start gap-2 rounded-lg border border-destructive/30 bg-destructive/10 px-3.5 py-2.5 text-sm text-destructive">
                  <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
                  <span>{checkInError}</span>
                </div>
              )}
              {checkInResult && (
                <div className="mt-3 flex items-start gap-2 rounded-lg border border-success/30 bg-success/10 px-3.5 py-2.5 text-sm text-success">
                  <CheckCircle2 className="h-4 w-4 mt-0.5 shrink-0" />
                  <span>
                    Tamu berhasil check-in
                    {checkInResult.attendedAt
                      ? ` · ${new Date(checkInResult.attendedAt).toLocaleString('id-ID')}`
                      : ''}
                  </span>
                </div>
              )}
            </CardContent>
          </Card>

          <div className="flex gap-2">
            <Button
              variant={view === 'active' ? 'primary' : 'outline'}
              onClick={() => setView('active')}
            >
              Tamu Aktif
            </Button>
            <Button
              variant={view === 'trash' ? 'primary' : 'outline'}
              onClick={() => {
                setView('trash')
                if (deletedGuests.length === 0) fetchDeletedGuests()
              }}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Sampah
            </Button>
          </div>

          {view === 'trash' ? (
            deletedLoading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-14 rounded-xl" />
                ))}
              </div>
            ) : (
              <Card>
                <CardContent className="p-0">
                  {deletedGuests.length === 0 ? (
                    <div className="text-center py-16 px-6">
                      <div className="mx-auto w-14 h-14 rounded-full bg-muted flex items-center justify-center mb-4">
                        <Trash2 className="h-6 w-6 text-muted-foreground" />
                      </div>
                      <p className="font-medium text-foreground">Sampah kosong</p>
                      <p className="text-sm text-muted-foreground mt-1">
                        Tamu yang dihapus akan muncul di sini dan bisa dipulihkan.
                      </p>
                    </div>
                  ) : (
                    <Table>
                      <TableHeader>
                        <TableRow className="hover:bg-transparent">
                          <TableHead>Nama</TableHead>
                          <TableHead>No. WhatsApp</TableHead>
                          <TableHead>Kategori</TableHead>
                          <TableHead className="text-right">Aksi</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {deletedGuests.map((guest) => (
                          <TableRow key={guest.id}>
                            <TableCell>
                              <div className="flex items-center gap-3">
                                <span className="inline-flex items-center justify-center h-9 w-9 rounded-full bg-muted text-muted-foreground text-xs font-bold">
                                  {guest.name.slice(0, 2).toUpperCase()}
                                </span>
                                <div>
                                  <p className="font-medium">{guest.name}</p>
                                  {guest.note && (
                                    <p className="text-xs text-muted-foreground">{guest.note}</p>
                                  )}
                                </div>
                              </div>
                            </TableCell>
                            <TableCell className="text-muted-foreground">{guest.phone}</TableCell>
                            <TableCell>
                              <Badge variant="outline" className={cn("border-0", categoryColors[guest.category] || "bg-muted")}>
                                {categoryLabels[guest.category] || guest.category}
                              </Badge>
                            </TableCell>
                            <TableCell>
                              <div className="flex justify-end">
                                <button
                                  onClick={() => handleRestoreGuest(guest.id)}
                                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-primary hover:bg-primary/10 transition-colors"
                                >
                                  <RotateCcw className="h-4 w-4" />
                                  Pulihkan
                                </button>
                              </div>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  )}
                </CardContent>
              </Card>
            )
          ) : (
          loading ? (
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-14 rounded-xl" />
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-0">
                {guests.length === 0 ? (
                  <div className="text-center py-16 px-6">
                    <div className="mx-auto w-14 h-14 rounded-full bg-muted flex items-center justify-center mb-4">
                      <Users className="h-6 w-6 text-muted-foreground" />
                    </div>
                    <p className="font-medium text-foreground">Belum ada tamu</p>
                    <p className="text-sm text-muted-foreground mt-1">
                      Tambahkan tamu pertamamu untuk mulai membagikan undangan.
                    </p>
                    <Button
                      onClick={() => setShowAddModal(true)}
                      className="mt-6 bg-gradient-to-r from-primary to-secondary border-0"
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Tambah Tamu
                    </Button>
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow className="hover:bg-transparent">
                        <TableHead>Nama</TableHead>
                        <TableHead>No. WhatsApp</TableHead>
                        <TableHead>Kategori</TableHead>
                        <TableHead className="text-right">Aksi</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {guests.map((guest) => (
                        <TableRow key={guest.id}>
                          <TableCell>
                            <div className="flex items-center gap-3">
                              <span className="inline-flex items-center justify-center h-9 w-9 rounded-full bg-gradient-to-br from-primary-container to-secondary-container text-primary text-xs font-bold">
                                {guest.name.slice(0, 2).toUpperCase()}
                              </span>
                              <div>
                                <p className="font-medium">{guest.name}</p>
                                {guest.note && (
                                  <p className="text-xs text-muted-foreground">{guest.note}</p>
                                )}
                              </div>
                            </div>
                          </TableCell>
                          <TableCell className="text-muted-foreground">{guest.phone}</TableCell>
                          <TableCell>
                            <Badge variant="outline" className={cn("border-0", categoryColors[guest.category] || "bg-muted")}>
                              {categoryLabels[guest.category] || guest.category}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <div className="flex justify-end gap-2">
                              {attendedMap[guest.id] ? (
                                <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-success">
                                  <CheckCircle2 className="h-4 w-4" />
                                  Hadir
                                </span>
                              ) : (
                                <button
                                  onClick={() => checkIn(guest.token)}
                                  disabled={checkingIn}
                                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-primary hover:bg-primary/10 transition-colors disabled:opacity-50"
                                >
                                  <QrCode className="h-4 w-4" />
                                  Check-in
                                </button>
                              )}
                              {guest.phone ? (
                                <a
                                  href={generateWhatsAppLink(guest)}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-emerald-600 hover:bg-emerald-500/10 transition-colors"
                                >
                                  <MessageCircle className="h-4 w-4" />
                                  WhatsApp
                                </a>
                              ) : null}
                              <button
                                onClick={() => handleDeleteGuest(guest.id)}
                                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-destructive hover:bg-destructive/10 transition-colors"
                              >
                                <Trash2 className="h-4 w-4" />
                                Hapus
                              </button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Add Guest Modal */}
        {showAddModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={() => setShowAddModal(false)} />
            <div className="relative w-full max-w-md bg-card rounded-2xl shadow-2xl border border-border p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-semibold text-foreground">Tambah Tamu</h2>
                <button onClick={() => setShowAddModal(false)} className="p-1.5 hover:bg-accent rounded-lg" aria-label="Close">
                  <X className="h-5 w-5" />
                </button>
              </div>

              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="guest-name">Nama</Label>
                  <Input
                    id="guest-name"
                    type="text"
                    placeholder="Nama tamu"
                    value={newGuest.name}
                    onChange={(e) => setNewGuest({ ...newGuest, name: e.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="guest-phone">No. WhatsApp</Label>
                  <Input
                    id="guest-phone"
                    type="text"
                    placeholder="08xxxxxxxxxx"
                    value={newGuest.phone}
                    onChange={(e) => setNewGuest({ ...newGuest, phone: e.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="guest-category">Kategori</Label>
                  <select
                    id="guest-category"
                    value={newGuest.category}
                    onChange={(e) => setNewGuest({ ...newGuest, category: e.target.value })}
                    className="flex h-11 w-full rounded-lg border border-input bg-background px-3.5 py-2 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  >
                    <option value="family">Keluarga</option>
                    <option value="friend">Teman</option>
                    <option value="colleague">Kolega</option>
                    <option value="other">Lainnya</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="guest-note">Catatan</Label>
                  <Textarea
                    id="guest-note"
                    placeholder="Catatan opsional"
                    value={newGuest.note}
                    onChange={(e) => setNewGuest({ ...newGuest, note: e.target.value })}
                    rows={3}
                  />
                </div>
              </div>

              <div className="flex justify-end gap-2 mt-6">
                <Button variant="ghost" onClick={() => setShowAddModal(false)}>
                  Batal
                </Button>
                <Button
                  onClick={handleAddGuest}
                  className="bg-gradient-to-r from-primary to-secondary hover:from-primary/90 hover:to-secondary/90 border-0"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Tambah
                </Button>
              </div>
            </div>
          </div>
        )}

        {eventId && (
          <GuestImportDialog
            open={showImportModal}
            onOpenChange={setShowImportModal}
            eventId={eventId}
            onImported={fetchGuests}
          />
        )}
      </DashboardLayout>
    </ProtectedRoute>
  )
}
