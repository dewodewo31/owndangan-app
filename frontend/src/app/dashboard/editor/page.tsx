"use client"

import { useEffect, useRef, useState, useCallback } from "react"
import { PenSquare, Eye, Upload, Trash2, Save, Plus, ImageIcon, Music as MusicIcon, Gift, Copy, ExternalLink, Check, Loader2 } from "lucide-react"
import ProtectedRoute from "@/components/dashboard/protected-route"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import InvitationPreview from "@/components/dashboard/invitation-preview"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Textarea } from "@/components/ui/textarea"
import api from "@/lib/api"
import type { WeddingEvent, EventSections, GalleryPhoto, Music, DigitalGift, TemplateSummary } from "@/lib/types"

type Tab = "detail" | "sections" | "gallery" | "music" | "gifts" | "template"

const VERSE_PRESETS: Record<string, { label: string; text: string; source: string }[]> = {
  quran: [
    {
      label: "Ar-Rum 30:21",
      text: "Dan di antara tanda-tanda (kebesaran)-Nya ialah Dia menciptakan pasangan-pasangan untukmu dari jenismu sendiri, agar kamu cenderung dan merasa tenteram kepadanya, dan Dia menjadikan di antaramu rasa kasih dan sayang. Sungguh, pada yang demikian itu benar-benar terdapat tanda-tanda (kebesaran Allah) bagi kaum yang berpikir.",
      source: "Q.S. Ar-Rum: 21",
    },
    {
      label: "An-Nisa 4:1",
      text: "Wahai manusia! Bertakwalah kepada Tuhanmu yang telah menciptakan kamu dari diri yang satu, dan daripadanya Dia menciptakan pasangannya, dan dari keduanya Dia memperkembangbiakkan laki-laki dan perempuan yang banyak.",
      source: "Q.S. An-Nisa: 1",
    },
    {
      label: "Al-Baqarah 2:187",
      text: "Dihalalkan bagi kamu pada malam hari sahur (berjimak) dengan istri-istrimu; mereka adalah pakaian bagimu, dan kamupun pakaian bagi mereka.",
      source: "Q.S. Al-Baqarah: 187",
    },
  ],
  alkitab: [
    {
      label: "1 Korintus 13:4-7",
      text: "Kasih itu sabar; kasih itu murah hati; ia tidak cemburu. Ia tidak memegahkan diri dan tidak sombong. Ia tidak melakukan yang tidak sopan dan tidak mencari keuntungan diri sendiri. Ia tidak pemarah dan tidak menyimpan kesalahan orang lain. Ia tidak bersukacita karena ketidakadilan, tetapi karena kebenaran. Ia menutupi segala sesuatu, percaya segala sesuatu, mengharapkan segala sesuatu, sabar menanggung segala sesuatu.",
      source: "1 Korintus 13:4-7",
    },
    {
      label: "Kejadian 2:24",
      text: "Sebab itu seorang laki-laki akan meninggalkan ayahnya dan ibunya dan bersatu dengan istrinya, sehingga keduanya menjadi satu daging.",
      source: "Kejadian 2:24",
    },
    {
      label: "Pengkhotbah 4:9-10",
      text: "Berdua lebih baik dari pada seorang diri, karena mereka menerima upah yang baik dalam jerih payah mereka. Karena kalau mereka jatuh, yang seorang dapat membangun temannya kembali; tetapi kalau ia seorang diri, padahal ia jatuh, tak ada yang membangun dia kembali.",
      source: "Pengkhotbah 4:9-10",
    },
  ],
}

function classNames(...c: (string | false | undefined)[]) {
  return c.filter(Boolean).join(" ")
}

function debounce<T extends (...a: any[]) => void>(fn: T, ms: number) {
  let t: any
  return (...a: Parameters<T>) => {
    clearTimeout(t)
    t = setTimeout(() => fn(...a), ms)
  }
}

export default function EditorPage() {
  const [events, setEvents] = useState<WeddingEvent[]>([])
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [event, setEvent] = useState<WeddingEvent | null>(null)
  const [sections, setSections] = useState<EventSections | null>(null)
  const [gallery, setGallery] = useState<GalleryPhoto[]>([])
  const [music, setMusic] = useState<Music | null>(null)
  const [presets, setPresets] = useState<Music[]>([])
  const [gift, setGift] = useState<DigitalGift | null>(null)
  const [templates, setTemplates] = useState<TemplateSummary[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [activeTab, setActiveTab] = useState<Tab>("detail")
  const [showCreate, setShowCreate] = useState(false)
  const [newTitle, setNewTitle] = useState("")

  const tabs: { id: Tab; label: string; icon: any }[] = [
    { id: "detail", label: "Detail Undangan", icon: PenSquare },
    { id: "sections", label: "Bagian", icon: Copy },
    { id: "gallery", label: "Galeri", icon: ImageIcon },
    { id: "music", label: "Musik", icon: MusicIcon },
    { id: "gifts", label: "Hadiah Digital", icon: Gift },
    { id: "template", label: "Templat", icon: ImageIcon },
  ]

  useEffect(() => {
    fetchEvents()
  }, [])

  useEffect(() => {
    if (selectedId) {
      fetchConfig(selectedId)
    }
  }, [selectedId])

  async function fetchEvents() {
    setLoading(true)
    try {
      const res = await api.get("/events")
      const list: WeddingEvent[] = res.data?.data || []
      setEvents(list)
      if (list.length > 0 && !selectedId) setSelectedId(list[0].id)
      if (list.length > 0 && !event) setEvent(list[0])
    } catch (e) {
      console.error("fetch events failed", e)
    } finally {
      setLoading(false)
    }
  }

  async function fetchConfig(id: string) {
    const safe = async <T,>(p: Promise<{ data: T }>, fallback: T) => {
      try {
        const r = await p
        return r.data as T
      } catch {
        return fallback
      }
    }
    try {
      const [ev, sec, gal, mus, gif, tmpl, presets] = await Promise.all([
        safe(api.get(`/events/${id}`).then((r) => r.data), null as any),
        safe(api.get(`/events/${id}/sections`).then((r) => r.data), null as any),
        safe(api.get(`/events/${id}/gallery`).then((r) => r.data), [] as GalleryPhoto[]),
        safe(api.get(`/events/${id}/music`).then((r) => r.data), null as unknown as Music | null),
        safe(api.get(`/events/${id}/digital-gifts`).then((r) => r.data), null as unknown as DigitalGift | null),
        safe(api.get("/templates").then((r) => r.data), [] as TemplateSummary[]),
        safe(api.get(`/events/${id}/music/presets`).then((r) => r.data), [] as Music[]),
      ])
      if (ev) setEvent(ev as WeddingEvent)
      setSections(sec as EventSections)
      setGallery(gal as GalleryPhoto[] || [])
      setMusic(mus as Music | null)
      setGift(gif as DigitalGift | null)
      setTemplates(tmpl as TemplateSummary[])
      setPresets(presets as Music[])
    } catch (e) {
      console.error("fetch config failed", e)
    }
  }

  const saveField = useRef(
    debounce(async (patch: Record<string, unknown>) => {
      if (!event) return
      setSaving(true)
      try {
        const res = await api.put(`/events/${event.id}`, patch)
        setEvent(res.data.data as WeddingEvent)
      } catch (e) {
        console.error("save failed", e)
        alert("Gagal menyimpan perubahan")
      } finally {
        setSaving(false)
      }
    }, 600)
  ).current

  function updateField(key: keyof WeddingEvent, value: unknown) {
    if (!event) return
    const next = { ...event, [key]: value }
    setEvent(next)
    saveField({ [key]: value })
  }

  async function createEvent() {
    if (!newTitle.trim()) return
    try {
      const res = await api.post("/events", { title: newTitle, couple_name: newTitle })
      const created = res.data.data as WeddingEvent
      setEvents([created, ...events])
      setSelectedId(created.id)
      setShowCreate(false)
      setNewTitle("")
    } catch (e: any) {
      alert(e?.message || "Gagal membuat undangan")
    }
  }

  async function saveSections() {
    if (!event || !sections) return
    setSaving(true)
    try {
      await api.put(`/events/${event.id}/sections`, bodyFromSections(sections))
    } catch (e) {
      alert("Gagal menyimpan bagian")
    } finally {
      setSaving(false)
    }
  }

  async function toggleSection(key: keyof EventSections, value: boolean) {
    if (!sections) return
    const next = { ...sections, [key]: value }
    setSections(next)
    setSaving(true)
    try {
      await api.put(`/events/${event!.id}/sections`, bodyFromSections(next))
    } catch {
      alert("Gagal menyimpan")
    } finally {
      setSaving(false)
    }
  }

  function bodyFromSections(s: EventSections) {
    return {
      hero_enabled: s.hero_enabled,
      couple_enabled: s.couple_enabled,
      event_details_enabled: s.event_details_enabled,
      gallery_enabled: s.gallery_enabled,
      video_enabled: s.video_enabled,
      music_id: s.music_id,
      rsvp_enabled: s.rsvp_enabled,
      guestbook_enabled: s.guestbook_enabled,
      digital_gifts_enabled: s.digital_gifts_enabled,
      dress_code: s.dress_code,
      closing_message: s.closing_message,
      opening_message: s.opening_message,
      verse_enabled: s.verse_enabled,
      verse_religion: s.verse_religion,
      verse_text: s.verse_text,
      verse_source: s.verse_source,
    }
  }

  function setSectionField(key: keyof EventSections, value: unknown) {
    if (!sections) return
    setSections({ ...sections, [key]: value })
  }

  async function saveSectionText() {
    await saveSections()
  }

  async function uploadGallery(e: React.ChangeEvent<HTMLInputElement>) {
    if (!event || !e.target.files || !e.target.files[0]) return
    const file = e.target.files[0]
    const form = new FormData()
    form.append("file", file)
    form.append("caption", "")
    try {
      const res = await api.post(`/events/${event.id}/gallery/upload`, form)
      const photo = res.data.data as GalleryPhoto
      setGallery((g) => [...g, photo])
    } catch (err: any) {
      alert(err?.message || "Upload gagal")
    }
    e.target.value = ""
  }

  async function deleteGallery(id: string) {
    if (!event) return
    try {
      await api.delete(`/events/${event.id}/gallery/${id}`)
      setGallery((g) => g.filter((p) => p.id !== id))
    } catch {
      alert("Hapus gagal")
    }
  }

  async function reorderGallery(reordered: GalleryPhoto[]) {
    if (!event) return
    setGallery(reordered)
    try {
      await api.put(`/events/${event.id}/gallery/reorder`, {
        photos: reordered.map((p, i) => ({ id: p.id, sort_order: i })),
      })
    } catch {
      alert("Urutan gagal disimpan")
    }
  }

  async function assignPreset(presetId: string) {
    if (!event) return
    try {
      await api.post(`/events/${event.id}/music/presets`, { preset_id: presetId })
    } catch (e: any) {
      alert(e?.message || "Gagal memilih musik")
    }
  }

  async function uploadMusic(e: React.ChangeEvent<HTMLInputElement>) {
    if (!event || !e.target.files || !e.target.files[0]) return
    const file = e.target.files[0]
    const form = new FormData()
    form.append("file", file)
    try {
      const res = await api.post(`/events/${event.id}/music/upload`, form)
      setMusic(res.data.data as Music)
    } catch (e: any) {
      alert(e?.message || "Upload musik gagal")
    }
    e.target.value = ""
  }

  async function saveGift() {
    if (!event || !gift) return
    setSaving(true)
    try {
      await api.put(`/events/${event.id}/digital-gifts`, {
        bank_accounts: gift.bank_accounts,
        ewallet: gift.ewallet,
        qris_image_url: gift.qris_image_url,
        gift_message: gift.gift_message,
      })
    } catch {
      alert("Gagal menyimpan hadiah")
    } finally {
      setSaving(false)
    }
  }

  async function assignTemplate(id: string) {
    if (!event) return
    try {
      const res = await api.put(`/events/${event.id}/template`, { template_id: id })
      const data = res.data.data as WeddingEvent
      setEvent(data)
    } catch (e: any) {
      alert(e?.message || "Gagal menerapkan templat")
    }
  }

  async function publish() {
    if (!event) return
    if (!window.confirm("Terbitkan undangan ini? Undangan akan dapat dilihat publik.")) return
    try {
      const res = await api.post(`/events/${event.id}/publish`)
      const data = res.data.data as WeddingEvent
      setEvent(data)
      alert("Undangan berhasil dipublikasikan!")
    } catch (e: any) {
      alert(e?.message || "Gagal mempublikasikan")
    }
  }

  function previewURL() {
    if (!event?.slug) return null
    return `/${event.slug}`
  }

  const hasUnsaved = saving
  const canPublish = event?.status === "draft"

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Editor Undangan</h1>
              <p className="mt-1 text-muted-foreground">Sesuaikan setiap bagian undanganmu.</p>
            </div>
            <div className="flex items-center gap-3">
              {event?.status === "published" && (
                <Badge variant="success" className="self-center">
                  <span className="h-1.5 w-1.5 rounded-full bg-current mr-1.5 animate-pulse" />
                  Terbit
                </Badge>
              )}
              {hasUnsaved && (
                <span className="text-xs text-muted-foreground flex items-center gap-1">
                  <Loader2 className="h-3 w-3 animate-spin" /> Menyimpan…
                </span>
              )}
              {event && canPublish && (
                <Button onClick={publish} className="bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 border-0">
                  <PenSquare className="h-4 w-4 mr-2" /> Terbitkan
                </Button>
              )}
              {previewURL() && event?.status === "published" && (
                <Button variant="outline" size="sm" asChild>
                  <a href={previewURL()!} target="_blank" rel="noreferrer">
                    <Eye className="h-4 w-4 mr-2" /> Pratinjau
                  </a>
                </Button>
              )}
            </div>
          </div>

          {!showCreate && (
            <div className="flex items-center gap-3">
              <Button variant={events.length === 0 ? "primary" : "outline"} onClick={() => setShowCreate(true)}>
                <Plus className="h-4 w-4 mr-2" /> Undangan Baru
              </Button>
            </div>
          )}

          {showCreate && (
            <Card>
              <CardHeader>
                <CardTitle>Buat Undangan Baru</CardTitle>
              </CardHeader>
              <CardContent className="flex items-end gap-3">
                <div className="flex-1 space-y-1">
                  <Label htmlFor="new-title">Judul Undangan</Label>
                  <Input id="new-title" placeholder="Mis. The Wedding of Adi & Aisyah" value={newTitle} onChange={(e) => setNewTitle(e.target.value)} />
                </div>
                <Button onClick={createEvent}>Buat</Button>
                <Button variant="ghost" onClick={() => setShowCreate(false)}>Batal</Button>
              </CardContent>
            </Card>
          )}

          {loading && !event ? (
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
              <Skeleton className="lg:col-span-1 h-96 rounded-2xl" />
              <Skeleton className="lg:col-span-3 h-96 rounded-2xl" />
            </div>
          ) : !event ? (
            <Card>
              <CardContent className="py-16 text-center">
                <PenSquare className="mx-auto h-10 w-10 text-muted-foreground" />
                <p className="mt-3 font-medium text-foreground">Belum ada undangan</p>
                <p className="text-sm text-muted-foreground">Buat undangan pertamamu untuk mulai mengedit.</p>
              </CardContent>
            </Card>
          ) : (
              <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
                <div className="lg:col-span-2">
                  <Card>
                    <CardContent className="p-4">
                      <p className="px-3 py-2 text-xs font-medium text-muted-foreground uppercase tracking-wider">Bagian Undangan</p>
                      <nav className="space-y-1">
                        {tabs.map((tab) => (
                          <button
                            key={tab.id}
                            onClick={() => setActiveTab(tab.id)}
                            className={classNames(
                              "w-full flex items-center gap-3 px-3 py-2.5 text-sm rounded-xl transition-all",
                              activeTab === tab.id
                                ? "bg-gradient-to-r from-indigo-600 to-violet-600 text-white shadow-md shadow-indigo-500/20"
                                : "text-muted-foreground hover:bg-accent hover:text-foreground"
                            )}
                          >
                            <tab.icon className="h-4 w-4" />
                            {tab.label}
                          </button>
                        ))}
                      </nav>
                    </CardContent>
                  </Card>
                </div>

                <div className="lg:col-span-5">
                  <Card>
                    <CardHeader>
                      <div className="flex items-center justify-between">
                        <CardTitle>{tabs.find((t) => t.id === activeTab)?.label}</CardTitle>
                        <Badge variant="secondary">{event.status}</Badge>
                      </div>
                    </CardHeader>
                    <CardContent>
                      {activeTab === "detail" && renderDetail()}
                      {activeTab === "sections" && renderSections()}
                      {activeTab === "gallery" && renderGallery()}
                      {activeTab === "music" && renderMusic()}
                      {activeTab === "gifts" && renderGifts()}
                      {activeTab === "template" && renderTemplate()}
                    </CardContent>
                  </Card>
                </div>

                <div className="lg:col-span-5 lg:sticky lg:top-6">
                  <InvitationPreview
                    event={event}
                    sections={sections}
                    gallery={gallery}
                    music={music}
                    gift={gift}
                    template={templates.find((t) => t.id === event?.template_id) || null}
                  />
                </div>
              </div>
          )}
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )

  function renderDetail() {
    if (!event) return null
    const input = (label: string, key: keyof WeddingEvent, opts?: any) => (
      <div className="space-y-1.5">
        <Label>{label}</Label>
        <Input
          defaultValue={event[key] as any}
          onBlur={(e) => updateField(key, opts?.date ? e.target.value : e.target.value)}
          type={opts?.type}
          placeholder={label}
        />
      </div>
    )
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
          {input("Judul Undangan", "title")}
          {input("Slug", "slug")}
          {input("Nama Mempelai Pria", "groom_name")}
          {input("Nama Mempelai Wanita", "bride_name")}
          {input("Nama Pasangan", "couple_name")}
          {input("Orang Tua Mempelai Pria", "groom_parents")}
          {input("Orang Tua Mempelai Wanita", "bride_parents")}
          {input("Tanggal Pernikahan", "wedding_date", { type: "date" })}
          {input("Waktu Pernikahan", "wedding_time")}
          {input("Lokasi Akad / Tempat Nikah", "ceremony_venue")}
        </div>
        <div className="space-y-1.5">
          <Label>Alamat Akad</Label>
          <Textarea defaultValue={event.ceremony_address || ""} onBlur={(e) => updateField("ceremony_address", e.target.value)} />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
          {input("Lokasi Resepsi", "reception_venue")}
          {input("Map URL Lokasi Akad", "ceremony_map_url")}
          {input("Map URL Resepsi", "reception_map_url")}
        </div>
        <div className="space-y-1.5">
          <Label>Alamat Resepsi</Label>
          <Textarea defaultValue={event.reception_address || ""} onBlur={(e) => updateField("reception_address", e.target.value)} />
        </div>
        <div className="flex justify-end">
          <Button size="sm" disabled={saving} onClick={() => {}}>
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4 mr-2" />}
            Simpan
          </Button>
        </div>
      </div>
    )
  }

  function renderSections() {
    if (!sections) return <Skeleton className="h-64 w-full" />
    const toggle = (label: string, key: keyof EventSections) => (
      <label className="flex items-center justify-between py-2">
        <span className="text-sm">{label}</span>
        <input
          type="checkbox"
          checked={!!(sections as any)[key]}
          onChange={(e) => toggleSection(key, e.target.checked)}
          className="h-4 w-4 rounded border-input accent-primary"
        />
      </label>
    )
    return (
      <div className="space-y-6">
        <div className="space-y-1">
          <h3 className="font-medium">Aktifkan / Nonaktifkan Bagian</h3>
          {toggle("Sorotan Hero", "hero_enabled")}
          {toggle("Profil Pasangan", "couple_enabled")}
          {toggle("Detail Acara", "event_details_enabled")}
          {toggle("Galeri", "gallery_enabled")}
          {toggle("Video", "video_enabled")}
          {toggle("RSVP", "rsvp_enabled")}
          {toggle("Buku Tamu", "guestbook_enabled")}
          {toggle("Hadiah Digital", "digital_gifts_enabled")}
        </div>

        <div className="space-y-1.5">
          <Label>Pesan Pembuka</Label>
          <Textarea value={sections.opening_message || ""} onChange={(e) => setSectionField("opening_message", e.target.value)} onBlur={saveSectionText} />
        </div>
        <div className="space-y-1.5">
          <Label>Pesan Penutup</Label>
          <Textarea value={sections.closing_message || ""} onChange={(e) => setSectionField("closing_message", e.target.value)} onBlur={saveSectionText} />
        </div>
        <div className="space-y-1.5">
          <Label>Dress Code</Label>
          <Input value={sections.dress_code || ""} onChange={(e) => setSectionField("dress_code", e.target.value)} onBlur={saveSectionText} placeholder="Mis. Black Tie" />
        </div>

        <div className="space-y-3 rounded-lg border p-3">
          {toggle("Ayat Suci (Al-Quran / Alkitab)", "verse_enabled")}
          {sections.verse_enabled && (
            <div className="space-y-3 pt-1">
              <div className="space-y-1.5">
                <Label>Sumber Ayat</Label>
                <select
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={sections.verse_religion || "quran"}
                  onChange={(e) => setSectionField("verse_religion", e.target.value)}
                  onBlur={saveSectionText}
                >
                  <option value="quran">Al-Quran</option>
                  <option value="alkitab">Alkitab</option>
                </select>
              </div>
              <div className="space-y-1.5">
                <Label>Ayat Pilihan</Label>
                <select
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value=""
                  onChange={(e) => {
                    const p = (VERSE_PRESETS[sections.verse_religion || "quran"] || []).find(
                      (x) => x.label === e.target.value
                    )
                    if (p) {
                      setSectionField("verse_text", p.text)
                      setSectionField("verse_source", p.source)
                      saveSectionText()
                    }
                  }}
                >
                  <option value="">Pilih ayat preset…</option>
                  {(VERSE_PRESETS[sections.verse_religion || "quran"] || []).map((p) => (
                    <option key={p.label} value={p.label}>
                      {p.label}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1.5">
                <Label>Teks Ayat</Label>
                <Textarea
                  rows={4}
                  value={sections.verse_text || ""}
                  onChange={(e) => setSectionField("verse_text", e.target.value)}
                  onBlur={saveSectionText}
                />
              </div>
              <div className="space-y-1.5">
                <Label>Sumber / Referensi</Label>
                <Input
                  value={sections.verse_source || ""}
                  onChange={(e) => setSectionField("verse_source", e.target.value)}
                  onBlur={saveSectionText}
                  placeholder="Mis. Q.S. Ar-Rum: 21"
                />
              </div>
            </div>
          )}
        </div>

        <div className="space-y-1">
          <h3 className="font-medium">Musik Latar</h3>
          <p className="text-sm text-muted-foreground">Musik pilihan akan diputar otomatis di undangan.</p>
          {presets.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {presets.map((p) => (
                <button
                  key={p.id}
                  type="button"
                  onClick={() => assignPreset(p.id)}
                  className={classNames(
                    "text-left px-3 py-2 rounded-lg border text-sm transition",
                    music?.id === p.id ? "border-primary bg-primary/10" : "border-input hover:bg-accent"
                  )}
                >
                  {p.title}
                </button>
              ))}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">Tidak ada lagu preset.</p>
          )}
          {music && !music.is_preset && music.file_url && (
            <audio controls src={music.file_url} className="mt-2" />
          )}
          {music && music.is_preset && music.preset && (
            <p className="text-sm text-muted-foreground">Preset aktif: {music.preset}</p>
          )}
        </div>
      </div>
    )
  }

  function renderGallery() {
    if (!event) return null
    const reorder = (from: number, to: number) => {
      const next = [...gallery]
      const [moved] = next.splice(from, 1)
      next.splice(to, 0, moved)
      reorderGallery(next)
    }
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="font-medium">Foto Galeri</h3>
          <Label className="cursor-pointer inline-flex items-center gap-2">
            <Upload className="h-4 w-4" /> Unggah
            <input type="file" accept="image/*" onChange={uploadGallery} className="hidden" />
          </Label>
        </div>
        {gallery.length === 0 ? (
          <div className="text-center py-12 border-2 border-dashed border-input rounded-xl">
            <ImageIcon className="mx-auto h-8 w-8 text-muted-foreground" />
            <p className="mt-2 text-sm text-muted-foreground">Belum ada foto. Unggah foto pertama.</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
            {gallery.map((p, i) => (
              <div key={p.id} className="relative group">
                <img src={p.image_url} alt={p.caption || "gallery"} className="aspect-square w-full object-cover rounded-lg border" />
                <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg flex items-center justify-center gap-2">
                  <button type="button" onClick={() => deleteGallery(p.id)} className="p-1 rounded bg-destructive text-destructive-foreground">
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
                <button
                  type="button"
                  draggable
                  onDragEnd={() => {}}
                  className="absolute bottom-1 left-1 right-1 text-left"
                  title="Seret untuk urutkan (drag & drop)"
                >
                  <MoveHandle photo={p} index={i} onMove={reorder} />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    )
  }

  function renderMusic() {
    if (!event) return null
    return (
      <div className="space-y-4">
        <h3 className="font-medium">Musik Latar</h3>
        <p className="text-sm text-muted-foreground">Pilih lagu dari preset atau unggah lagu milikmu.</p>
        <div className="flex items-center gap-3 pt-2">
          <Label className="cursor-pointer">
            <Upload className="h-4 w-4 mr-2 inline" /> Unggah MP3
            <input type="file" accept="audio/*" onChange={uploadMusic} className="hidden" />
          </Label>
        </div>
        {music ? (
          <div className="flex items-center gap-3 p-3 border rounded-lg">
            {music.is_preset ? <MusicIcon className="h-5 w-5 text-primary" /> : <MusicIcon className="h-5 w-5" />}
            <div className="flex-1">
              <p className="font-medium">{music.title}</p>
              {music.preset && <p className="text-xs text-muted-foreground">{music.preset}</p>}
            </div>
            {music.file_url && <audio controls src={music.file_url} />}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">Tidak ada musik yang dipilih.</p>
        )}
      </div>
    )
  }

  function renderGifts() {
    if (!gift) return <Skeleton className="h-48 w-full" />
    const bankFields = (gift.bank_accounts && gift.bank_accounts.length > 0 ? gift.bank_accounts : [{}]) as Array<Record<string, unknown>>
    return (
      <div className="space-y-6">
        <h3 className="font-medium">Informasi Hadiah</h3>
        <div className="space-y-1.5">
          <Label>Pesan Hadiah</Label>
          <Textarea
            value={gift.gift_message || ""}
            onChange={(e) => setGift({ ...gift, gift_message: e.target.value })}
            onBlur={saveGift}
            placeholder="Terima kasih atas kehadiran..."
          />
        </div>
        {bankFields.map((b, i) => (
          <div key={i} className="grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
            <div className="space-y-1"><Label>Bank</Label><Input value={String(b.bank || "")} onChange={(e) => {
              const copy = [...bankFields]; copy[i] = { ...copy[i], bank: e.target.value }; setGift({ ...gift, bank_accounts: copy })
            }} /></div>
            <div className="space-y-1"><Label>Nomor Rekening</Label><Input value={String(b.account || "")} onChange={(e) => {
              const copy = [...bankFields]; copy[i] = { ...copy[i], account: e.target.value }; setGift({ ...gift, bank_accounts: copy })
            }} /></div>
            <div className="space-y-1"><Label>Nama Pemilik</Label><Input value={String(b.name || "")} onChange={(e) => {
              const copy = [...bankFields]; copy[i] = { ...copy[i], name: e.target.value }; setGift({ ...gift, bank_accounts: copy })
            }} /></div>
          </div>
        ))}
        <div className="flex justify-between pt-2">
          <Button variant="outline" size="sm" onClick={() => {
            const copy = [...(gift.bank_accounts || []), {}]
            setGift({ ...gift, bank_accounts: copy })
          }}>Tambah Akun</Button>
          <Button size="sm" disabled={saving} onClick={saveGift}>{saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4 mr-2" />} Simpan Hadiah</Button>
        </div>
      </div>
    )
  }

  function renderTemplate() {
    return (
      <div className="space-y-4">
        <h3 className="font-medium">Pilih Templat</h3>
        <p className="text-sm text-muted-foreground">Tampilan undangan publik akan mengikuti templat yang dipilih.</p>
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
          {templates.map((t) => (
            <button
              key={t.id}
              type="button"
              onClick={() => assignTemplate(t.id)}
              className={classNames(
                "text-left rounded-xl border p-3 transition",
                event?.template_id === t.id ? "border-primary ring-2 ring-primary" : "border-input hover:border-primary"
              )}
            >
              {t.thumbnail_url ? <img src={t.thumbnail_url} alt={t.name} className="aspect-video w-full object-cover rounded" /> : <div className="aspect-video w-full bg-muted rounded flex items-center justify-center"><Copy className="h-6 w-6 text-muted-foreground" /></div>}
              <p className="mt-2 font-medium text-sm">{t.name}</p>
              <Badge variant="secondary" className="mt-1">{t.group_name}</Badge>
            </button>
          ))}
        </div>
      </div>
    )
  }
}

function MoveHandle({ photo, index, onMove }: { photo: GalleryPhoto; index: number; onMove: (from: number, to: number) => void }) {
  return (
    <button
      type="button"
      draggable
      onDragStart={(e) => e.dataTransfer.setData("text/plain", String(index))}
      onDragOver={(e) => e.preventDefault()}
      onDrop={(e) => {
        const from = Number(e.dataTransfer.getData("text/plain"))
        if (from !== index) onMove(from, index)
      }}
      className="flex items-center gap-1 text-xs text-muted-foreground"
    >
      <span>↕</span> #{index + 1}
    </button>
  )
}
