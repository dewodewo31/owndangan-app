"use client"

import { useEffect, useState } from "react"
import { Package, Boxes } from "lucide-react"
import AdminProtectedRoute from "@/components/admin/admin-protected-route"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { EmptyState, ErrorState } from "@/components/ui/empty-state"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog"
import api from "@/lib/api"
import { formatCurrency } from "@/lib/utils"
import type { AdminPackage } from "@/lib/types"

export default function AdminPackagesPage() {
  const [packages, setPackages] = useState<AdminPackage[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [selected, setSelected] = useState<AdminPackage | null>(null)

  const fetchPackages = async () => {
    setLoading(true)
    setError(false)
    try {
      const res = await api.get("/admin/packages")
      setPackages(res.data?.data ?? [])
    } catch {
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchPackages()
  }, [])

  return (
    <AdminProtectedRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Paket</h1>
          <p className="mt-1 text-muted-foreground">Daftar paket langganan platform.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5 text-primary" />
              Daftar Paket
            </CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-12 rounded-xl" />
                ))}
              </div>
            ) : error ? (
              <ErrorState retry={fetchPackages} />
            ) : packages.length === 0 ? (
              <EmptyState icon={<Boxes className="h-8 w-8" />} title="Belum ada paket" description="Paket langganan belum ditambahkan." />
            ) : (
              <Table>
                <TableHeader>
                  <TableRow className="hover:bg-transparent">
                    <TableHead>Nama</TableHead>
                    <TableHead>Kode</TableHead>
                    <TableHead>Harga</TableHead>
                    <TableHead>Grup Template</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Aksi</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {packages.map((p) => (
                    <TableRow key={p.id}>
                      <TableCell className="font-medium">{p.name}</TableCell>
                      <TableCell className="text-muted-foreground">{p.code}</TableCell>
                      <TableCell>{formatCurrency(p.price)}</TableCell>
                      <TableCell className="text-muted-foreground capitalize">{p.template_group}</TableCell>
                      <TableCell>
                        <Badge variant={p.is_active ? "success" : "secondary"}>
                          {p.is_active ? "Aktif" : "Nonaktif"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <Button size="sm" variant="outline" onClick={() => setSelected(p)}>
                          Detail
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>
      </div>

      <Dialog open={!!selected} onOpenChange={(open) => !open && setSelected(null)}>
        <DialogContent>
          {selected && (
            <>
              <DialogHeader>
                <DialogTitle className="flex items-center gap-2">
                  <Package className="h-5 w-5 text-primary" />
                  {selected.name}
                </DialogTitle>
              </DialogHeader>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <Badge variant="secondary">{selected.code}</Badge>
                  <Badge variant={selected.is_active ? "success" : "secondary"}>
                    {selected.is_active ? "Aktif" : "Nonaktif"}
                  </Badge>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Harga</p>
                    <p className="mt-1 font-semibold">{formatCurrency(selected.price)}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Grup Template</p>
                    <p className="mt-1 font-medium capitalize">{selected.template_group}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Durasi</p>
                    <p className="mt-1 font-medium">{selected.duration_days ? `${selected.duration_days} hari` : "-"}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Batas Tamu</p>
                    <p className="mt-1 font-medium">{selected.guest_limit ?? "-"}</p>
                  </div>
                </div>

                {selected.features && Object.keys(selected.features).length > 0 && (
                  <div className="rounded-xl border border-border p-4">
                    <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-2">Fitur</p>
                    <ul className="space-y-1.5 text-sm">
                      {Object.entries(selected.features).map(([k, v]) => (
                        <li key={k} className="flex justify-between gap-4">
                          <span className="text-muted-foreground">{k}</span>
                          <span className="font-medium text-right break-all">{String(v)}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>

              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">Tutup</Button>
                </DialogClose>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </AdminProtectedRoute>
  )
}
