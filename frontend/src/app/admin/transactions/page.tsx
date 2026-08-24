"use client"

import { useEffect, useState, useCallback } from "react"
import { Receipt, Search } from "lucide-react"
import AdminProtectedRoute from "@/components/admin/admin-protected-route"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Select } from "@/components/ui/select"
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
import { formatCurrency, formatDate } from "@/lib/utils"
import type { AdminTransaction, PaginationMeta } from "@/lib/types"

const PER_PAGE = 10

function statusVariant(status: string) {
  switch (status) {
    case "settlement":
    case "paid":
    case "active":
      return "success"
    case "pending":
      return "warning"
    case "expire":
    case "cancel":
    case "deny":
    case "expired":
      return "destructive"
    default:
      return "secondary"
  }
}

function statusLabel(status: string) {
  switch (status) {
    case "settlement":
      return "Berhasil"
    case "pending":
      return "Menunggu"
    case "expire":
    case "expired":
      return "Kedaluwarsa"
    case "cancel":
    case "deny":
      return "Dibatalkan"
    default:
      return status
  }
}

export default function AdminTransactionsPage() {
  const [txns, setTxns] = useState<AdminTransaction[]>([])
  const [pagination, setPagination] = useState<PaginationMeta | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState("")
  const [packageFilter, setPackageFilter] = useState("")

  const [selected, setSelected] = useState<AdminTransaction | null>(null)
  const [detail, setDetail] = useState<AdminTransaction | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)

  const fetchTxns = useCallback(async () => {
    setLoading(true)
    setError(false)
    try {
      const params = new URLSearchParams({ page: String(page), per_page: String(PER_PAGE) })
      if (statusFilter) params.set("status", statusFilter)
      const res = await api.get(`/admin/transactions?${params.toString()}`)
      setTxns(res.data?.data ?? [])
      setPagination(res.data?.meta?.pagination ?? null)
    } catch {
      setError(true)
    } finally {
      setLoading(false)
    }
  }, [page, statusFilter])

  useEffect(() => {
    fetchTxns()
  }, [fetchTxns])

  const filtered = txns.filter((t) => {
    if (!packageFilter) return true
    return t.package_id === packageFilter
  })

  const openDetail = async (txn: AdminTransaction) => {
    setSelected(txn)
    setDetailLoading(true)
    setDetail(null)
    try {
      const res = await api.get(`/admin/transactions/${txn.id}`)
      setDetail(res.data?.data ?? txn)
    } catch {
      setDetail(txn)
    } finally {
      setDetailLoading(false)
    }
  }

  const packageOptions = Array.from(new Set(txns.map((t) => t.package_id)))
    .map((id) => {
      const t = txns.find((x) => x.package_id === id)
      return { value: id, label: t?.package?.name || t?.package?.code || id.slice(0, 8) }
    })

  const totalPages = pagination?.total_pages ?? 1

  return (
    <AdminProtectedRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Transaksi</h1>
          <p className="mt-1 text-muted-foreground">Pantau pembayaran dan status Midtrans.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Receipt className="h-5 w-5 text-primary" />
              Daftar Transaksi
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-col sm:flex-row gap-3">
              <Select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value)
                  setPage(1)
                }}
                options={[
                  { value: "", label: "Semua Status" },
                  { value: "settlement", label: "Berhasil" },
                  { value: "pending", label: "Menunggu" },
                  { value: "expire", label: "Kedaluwarsa" },
                  { value: "cancel", label: "Dibatalkan" },
                  { value: "deny", label: "Ditolak" },
                ]}
                className="sm:w-48"
              />
              <Select
                value={packageFilter}
                onChange={(e) => setPackageFilter(e.target.value)}
                options={[{ value: "", label: "Semua Paket" }, ...packageOptions]}
                className="sm:w-56"
              />
            </div>

            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                  <Skeleton key={i} className="h-12 rounded-xl" />
                ))}
              </div>
            ) : error ? (
              <ErrorState retry={fetchTxns} />
            ) : filtered.length === 0 ? (
              <EmptyState icon={<Receipt className="h-8 w-8" />} title="Tidak ada transaksi" description="Coba ubah filter status atau paket." />
            ) : (
              <>
                <Table>
                  <TableHeader>
                    <TableRow className="hover:bg-transparent">
                      <TableHead>Order ID</TableHead>
                      <TableHead>Pengguna</TableHead>
                      <TableHead>Paket</TableHead>
                      <TableHead>Jumlah</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Tanggal</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filtered.map((t) => (
                      <TableRow
                        key={t.id}
                        className="cursor-pointer hover:bg-accent/60"
                        onClick={() => openDetail(t)}
                      >
                        <TableCell className="font-medium truncate max-w-[140px]">{t.order_id}</TableCell>
                        <TableCell className="text-muted-foreground truncate max-w-[140px]">{t.user?.name ?? "-"}</TableCell>
                        <TableCell className="text-muted-foreground">{t.package?.name ?? "-"}</TableCell>
                        <TableCell>{formatCurrency(t.gross_amount)}</TableCell>
                        <TableCell>
                          <Badge variant={statusVariant(t.status)}>{statusLabel(t.status)}</Badge>
                        </TableCell>
                        <TableCell className="text-muted-foreground">{formatDate(t.created_at)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>

                <div className="flex items-center justify-between pt-2">
                  <p className="text-sm text-muted-foreground">
                    Halaman {page} dari {totalPages}
                    {pagination ? ` · ${pagination.total} transaksi` : ""}
                  </p>
                  <div className="flex gap-2">
                    <Button size="sm" variant="outline" disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))}>
                      Sebelumnya
                    </Button>
                    <Button size="sm" variant="outline" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)}>
                      Berikutnya
                    </Button>
                  </div>
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </div>

      <Dialog open={!!selected} onOpenChange={(open) => !open && setSelected(null)}>
        <DialogContent>
          {detailLoading ? (
            <div className="py-8">
              <Skeleton className="h-40 rounded-xl" />
            </div>
          ) : detail ? (
            <>
              <DialogHeader>
                <DialogTitle className="flex items-center gap-2">
                  <Receipt className="h-5 w-5 text-primary" />
                  Detail Transaksi
                </DialogTitle>
              </DialogHeader>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <p className="font-mono text-sm text-muted-foreground break-all">{detail.order_id}</p>
                  <Badge variant={statusVariant(detail.status)}>{statusLabel(detail.status)}</Badge>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Jumlah</p>
                    <p className="mt-1 font-semibold">{formatCurrency(detail.gross_amount)}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Metode</p>
                    <p className="mt-1 font-medium capitalize">{detail.payment_type || "-"}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Waktu Pembayaran</p>
                    <p className="mt-1 font-medium">{detail.transaction_time ? formatDate(detail.transaction_time) : "-"}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Waktu Settlement</p>
                    <p className="mt-1 font-medium">{detail.settlement_time ? formatDate(detail.settlement_time) : "-"}</p>
                  </div>
                </div>

                <div className="rounded-xl border border-border p-4">
                  <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-2">Informasi Pengguna</p>
                  <p className="font-medium">{detail.user?.name ?? "-"}</p>
                  <p className="text-sm text-muted-foreground">{detail.user?.email ?? "-"}</p>
                  <p className="text-sm text-muted-foreground mt-1">Paket: {detail.package?.name ?? "-"}</p>
                </div>
              </div>

              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">Tutup</Button>
                </DialogClose>
              </DialogFooter>
            </>
          ) : null}
        </DialogContent>
      </Dialog>
    </AdminProtectedRoute>
  )
}
