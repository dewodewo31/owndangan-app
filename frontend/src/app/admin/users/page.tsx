"use client"

import { useEffect, useState, useCallback } from "react"
import { Search, Users, UserCog, ShieldAlert } from "lucide-react"
import AdminProtectedRoute from "@/components/admin/admin-protected-route"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
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
import { formatDate } from "@/lib/utils"
import type { AdminUser, PaginationMeta } from "@/lib/types"

const PER_PAGE = 10

function statusVariant(status: string) {
  return status === "active" ? "success" : status === "suspended" ? "destructive" : "secondary"
}

function statusLabel(status: string) {
  return status === "active" ? "Aktif" : status === "suspended" ? "Ditangguhkan" : status
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [pagination, setPagination] = useState<PaginationMeta | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState("")
  const [search, setSearch] = useState("")

  const [selected, setSelected] = useState<AdminUser | null>(null)
  const [toggling, setToggling] = useState(false)

  const fetchUsers = useCallback(async () => {
    setLoading(true)
    setError(false)
    try {
      const params = new URLSearchParams({ page: String(page), per_page: String(PER_PAGE) })
      if (statusFilter) params.set("status", statusFilter)
      const res = await api.get(`/admin/users?${params.toString()}`)
      setUsers(res.data?.data ?? [])
      setPagination(res.data?.meta?.pagination ?? null)
    } catch {
      setError(true)
    } finally {
      setLoading(false)
    }
  }, [page, statusFilter])

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  const filtered = users.filter((u) => {
    if (!search) return true
    const q = search.toLowerCase()
    return u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q)
  })

  const toggleStatus = async () => {
    if (!selected) return
    const next = selected.status === "active" ? "suspended" : "active"
    setToggling(true)
    try {
      await api.put(`/admin/users/${selected.id}/status`, { status: next })
      setSelected({ ...selected, status: next })
      fetchUsers()
    } catch {
      setError(true)
    } finally {
      setToggling(false)
    }
  }

  const totalPages = pagination?.total_pages ?? 1

  return (
    <AdminProtectedRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Pengguna</h1>
          <p className="mt-1 text-muted-foreground">Kelola akun pengguna platform.</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5 text-primary" />
              Daftar Pengguna
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-col sm:flex-row gap-3">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Cari nama atau email..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9"
                />
              </div>
              <Select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value)
                  setPage(1)
                }}
                options={[
                  { value: "", label: "Semua Status" },
                  { value: "active", label: "Aktif" },
                  { value: "suspended", label: "Ditangguhkan" },
                ]}
                className="sm:w-48"
              />
            </div>

            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                  <Skeleton key={i} className="h-12 rounded-xl" />
                ))}
              </div>
            ) : error ? (
              <ErrorState retry={fetchUsers} />
            ) : filtered.length === 0 ? (
              <EmptyState icon={<Users className="h-8 w-8" />} title="Tidak ada pengguna" description="Coba ubah filter atau kata kunci pencarian." />
            ) : (
              <>
                <Table>
                  <TableHeader>
                    <TableRow className="hover:bg-transparent">
                      <TableHead>Nama</TableHead>
                      <TableHead>Email</TableHead>
                      <TableHead>Peran</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Bergabung</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filtered.map((u) => (
                      <TableRow key={u.id}>
                        <TableCell className="font-medium">{u.name}</TableCell>
                        <TableCell className="text-muted-foreground">{u.email}</TableCell>
                        <TableCell>
                          <Badge variant={u.role === "admin" ? "default" : "secondary"}>{u.role}</Badge>
                        </TableCell>
                        <TableCell>
                          <Badge variant={statusVariant(u.status)}>{statusLabel(u.status)}</Badge>
                        </TableCell>
                        <TableCell className="text-muted-foreground">{formatDate(u.created_at)}</TableCell>
                        <TableCell className="text-right">
                          <Button size="sm" variant="outline" onClick={() => setSelected(u)}>
                            Detail
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>

                <div className="flex items-center justify-between pt-2">
                  <p className="text-sm text-muted-foreground">
                    Halaman {page} dari {totalPages}
                    {pagination ? ` · ${pagination.total} pengguna` : ""}
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
          {selected && (
            <>
              <DialogHeader>
                <DialogTitle className="flex items-center gap-2">
                  <UserCog className="h-5 w-5 text-primary" />
                  Detail Pengguna
                </DialogTitle>
              </DialogHeader>

              <div className="space-y-4">
                <div className="flex items-center gap-4">
                  <span className="inline-flex items-center justify-center h-14 w-14 rounded-2xl bg-gradient-to-br from-plum to-rosegold text-white text-lg font-bold">
                    {selected.name.split(" ").map((n) => n[0]).slice(0, 2).join("").toUpperCase()}
                  </span>
                  <div className="min-w-0">
                    <p className="text-lg font-bold text-foreground truncate">{selected.name}</p>
                    <p className="text-sm text-muted-foreground truncate">{selected.email}</p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Peran</p>
                    <p className="mt-1 font-medium capitalize">{selected.role}</p>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Status Akun</p>
                    <div className="mt-1">
                      <Badge variant={statusVariant(selected.status)}>{statusLabel(selected.status)}</Badge>
                    </div>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Status Langganan</p>
                    <div className="mt-1">
                      <Badge variant={selected.subscriptions?.[0]?.status === "active" ? "success" : "secondary"}>
                        {selected.subscriptions?.[0]?.status ?? "Tidak ada"}
                      </Badge>
                    </div>
                  </div>
                  <div className="rounded-xl bg-muted/50 p-3">
                    <p className="text-xs text-muted-foreground">Bergabung</p>
                    <p className="mt-1 font-medium">{formatDate(selected.created_at)}</p>
                  </div>
                </div>
              </div>

              <DialogFooter className="gap-2">
                <DialogClose asChild>
                  <Button variant="outline">Tutup</Button>
                </DialogClose>
                <Button
                  variant={selected.status === "active" ? "destructive" : "primary"}
                  loading={toggling}
                  onClick={toggleStatus}
                >
                  {selected.status === "active" ? (
                    <>
                      <ShieldAlert className="h-4 w-4 mr-2" />
                      Tangguhkan
                    </>
                  ) : (
                    "Aktifkan"
                  )}
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </AdminProtectedRoute>
  )
}
