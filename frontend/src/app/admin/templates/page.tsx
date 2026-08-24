"use client"

import { useEffect, useState } from "react"
import { LayoutTemplate, ImageOff, Power, PowerOff } from "lucide-react"
import AdminProtectedRoute from "@/components/admin/admin-protected-route"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { EmptyState, ErrorState } from "@/components/ui/empty-state"
import api from "@/lib/api"
import type { AdminTemplate } from "@/lib/types"

export default function AdminTemplatesPage() {
  const [templates, setTemplates] = useState<AdminTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [togglingId, setTogglingId] = useState<string | null>(null)

  const fetchTemplates = async () => {
    setLoading(true)
    setError(false)
    try {
      const res = await api.get("/admin/templates")
      setTemplates(res.data?.data ?? [])
    } catch {
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTemplates()
  }, [])

  const toggle = async (tpl: AdminTemplate) => {
    setTogglingId(tpl.id)
    try {
      const res = await api.put(`/admin/templates/${tpl.id}`, { is_active: !tpl.is_active })
      const updated: AdminTemplate = res.data?.data ?? { ...tpl, is_active: !tpl.is_active }
      setTemplates((prev) => prev.map((t) => (t.id === tpl.id ? updated : t)))
    } catch {
      setError(true)
    } finally {
      setTogglingId(null)
    }
  }

  return (
    <AdminProtectedRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Template</h1>
          <p className="mt-1 text-muted-foreground">Kelola template undangan dan status ketersediaannya.</p>
        </div>

        {loading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <Skeleton key={i} className="h-64 rounded-2xl" />
            ))}
          </div>
        ) : error ? (
          <ErrorState retry={fetchTemplates} />
        ) : templates.length === 0 ? (
          <Card>
            <CardContent>
              <EmptyState icon={<LayoutTemplate className="h-8 w-8" />} title="Belum ada template" description="Template undangan belum ditambahkan." />
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
            {templates.map((tpl) => (
              <Card key={tpl.id} className="overflow-hidden transition-all duration-300 hover:-translate-y-1 hover:shadow-elevation-2">
                <div className="relative aspect-[4/3] bg-muted flex items-center justify-center">
                  {tpl.thumbnail_url ? (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img src={tpl.thumbnail_url} alt={tpl.name} className="h-full w-full object-cover" />
                  ) : (
                    <ImageOff className="h-8 w-8 text-muted-foreground" />
                  )}
                  <span className="absolute top-3 right-3">
                    <Badge variant={tpl.is_active ? "success" : "secondary"}>
                      {tpl.is_active ? "Aktif" : "Nonaktif"}
                    </Badge>
                  </span>
                </div>
                <CardContent className="p-5 space-y-3">
                  <div>
                    <p className="font-semibold text-foreground">{tpl.name}</p>
                    <p className="text-sm text-muted-foreground capitalize">{tpl.group_name}</p>
                  </div>
                  <Button
                    variant={tpl.is_active ? "outline" : "primary"}
                    className="w-full"
                    loading={togglingId === tpl.id}
                    onClick={() => toggle(tpl)}
                  >
                    {tpl.is_active ? (
                      <>
                        <PowerOff className="h-4 w-4 mr-2" />
                        Nonaktifkan
                      </>
                    ) : (
                      <>
                        <Power className="h-4 w-4 mr-2" />
                        Aktifkan
                      </>
                    )}
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </AdminProtectedRoute>
  )
}
