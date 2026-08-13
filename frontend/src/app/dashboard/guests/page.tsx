"use client"

import { useEffect, useState } from "react"
import { Plus, MessageCircle, Trash2, Users, Download, X } from "lucide-react"
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
  friend: "bg-indigo-500/10 text-indigo-600",
  colleague: "bg-emerald-500/10 text-emerald-600",
  other: "bg-muted text-muted-foreground",
}

export default function GuestsPage() {
  const [guests, setGuests] = useState<Guest[]>([])
  const [loading, setLoading] = useState(true)
  const [showAddModal, setShowAddModal] = useState(false)
  const [newGuest, setNewGuest] = useState({ name: '', phone: '', category: 'family', note: '' })

  useEffect(() => {
    fetchGuests()
  }, [])

  const fetchGuests = async () => {
    try {
      const eventsRes = await api.get('/events')
      const events = eventsRes.data?.data || []
      if (events.length > 0) {
        const guestsRes = await api.get(`/events/${events[0].id}/guests`)
        setGuests(guestsRes.data?.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch guests:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAddGuest = async () => {
    try {
      const eventsRes = await api.get('/events')
      const events = eventsRes.data?.data || []
      if (events.length > 0) {
        await api.post(`/events/${events[0].id}/guests`, newGuest)
        setShowAddModal(false)
        setNewGuest({ name: '', phone: '', category: 'family', note: '' })
        fetchGuests()
      }
    } catch (error) {
      alert('Gagal menambahkan tamu')
    }
  }

  const handleDeleteGuest = async (guestID: string) => {
    if (!confirm('Yakin ingin menghapus tamu ini?')) return
    try {
      const eventsRes = await api.get('/events')
      const events = eventsRes.data?.data || []
      if (events.length > 0) {
        await api.delete(`/events/${events[0].id}/guests/${guestID}`)
        fetchGuests()
      }
    } catch (error) {
      alert('Gagal menghapus tamu')
    }
  }

  const generateWhatsAppLink = (guest: Guest) => {
    const message = `Halo ${guest.name}, kami mengundangmu untuk menghadiri pernikahan kami! Lihat undangan di: ${window.location.origin}/${guest.token}`
    return `https://wa.me/${guest.phone.replace(/^0/, '62')}?text=${encodeURIComponent(message)}`
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
              <Button variant="outline">
                <Download className="h-4 w-4 mr-2" />
                Impor CSV
              </Button>
              <Button
                onClick={() => setShowAddModal(true)}
                className="bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 border-0"
              >
                <Plus className="h-4 w-4 mr-2" />
                Tambah Tamu
              </Button>
            </div>
          </div>

          {loading ? (
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
                      className="mt-6 bg-gradient-to-r from-indigo-600 to-violet-600 border-0"
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
                              <span className="inline-flex items-center justify-center h-9 w-9 rounded-full bg-gradient-to-br from-indigo-100 to-violet-100 text-primary text-xs font-bold">
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
                              <a
                                href={generateWhatsAppLink(guest)}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-emerald-600 hover:bg-emerald-500/10 transition-colors"
                              >
                                <MessageCircle className="h-4 w-4" />
                                WhatsApp
                              </a>
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
          )}
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
                  className="bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 border-0"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Tambah
                </Button>
              </div>
            </div>
          </div>
        )}
      </DashboardLayout>
    </ProtectedRoute>
  )
}
