"use client"

import { useEffect, useState } from "react"
import {
  Users,
  CreditCard,
  Mail,
  Receipt,
  Wallet,
  Package,
  LayoutTemplate,
  ArrowUpRight,
} from "lucide-react"
import AdminProtectedRoute from "@/components/admin/admin-protected-route"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { EmptyState } from "@/components/ui/empty-state"
import api from "@/lib/api"
import { formatCurrency, formatDate } from "@/lib/utils"
import type { AdminAnalytics, AdminTransaction, AdminUser } from "@/lib/types"

function txnStatusVariant(status: string) {
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

function txnStatusLabel(status: string) {
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

function userStatusVariant(status: string) {
  return status === "active" ? "success" : status === "suspended" ? "destructive" : "secondary"
}

export default function AdminOverviewPage() {
  const [analytics, setAnalytics] = useState<AdminAnalytics | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  useEffect(() => {
    fetchAnalytics()
  }, [])

  const fetchAnalytics = async () => {
    try {
      const res = await api.get("/admin/analytics")
      setAnalytics(res.data?.data ?? null)
    } catch {
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  const statCards = [
    { title: "Total Pengguna", value: analytics?.total_users ?? 0, icon: Users, iconClass: "bg-primary-container text-primary" },
    { title: "Langganan Aktif", value: analytics?.active_subscriptions ?? 0, icon: CreditCard, iconClass: "bg-secondary-container text-secondary" },
    { title: "Total Undangan", value: analytics?.total_invitations ?? 0, icon: Mail, iconClass: "bg-tertiary-container text-tertiary" },
    { title: "Total Transaksi", value: analytics?.total_transactions ?? 0, icon: Receipt, iconClass: "bg-surface-container-high text-on-surface-variant" },
    { title: "Total Pendapatan", value: formatCurrency(analytics?.total_revenue ?? 0), icon: Wallet, iconClass: "bg-primary-container text-primary", isCurrency: true },
    { title: "Paket Aktif", value: analytics?.active_packages ?? 0, icon: Package, iconClass: "bg-secondary-container text-secondary" },
    { title: "Template Aktif", value: analytics?.active_templates ?? 0, icon: LayoutTemplate, iconClass: "bg-tertiary-container text-tertiary" },
  ]

  const recentTxns: AdminTransaction[] = analytics?.recent_transactions ?? []
  const recentUsers: AdminUser[] = analytics?.recent_users ?? []

  return (
    <AdminProtectedRoute>
      <div className="space-y-8">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Overview</h1>
            <p className="mt-1 text-muted-foreground">Ringkasan platform Owndangan.</p>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
              {[1, 2, 3, 4, 5, 6, 7].map((i) => (
                <Card key={i}>
                  <CardContent className="p-6">
                    <Skeleton className="h-4 w-1/2 mb-3" />
                    <Skeleton className="h-8 w-1/3" />
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : error ? (
            <EmptyState title="Gagal memuat data" description="Terjadi kesalahan saat mengambil analytics." action={{ label: "Coba Lagi", onClick: fetchAnalytics }} />
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
              {statCards.map((stat) => (
                <Card key={stat.title} className="transition-all duration-300 hover:-translate-y-1 hover:shadow-elevation-2">
                  <CardContent className="p-6">
                    <div className="flex items-start justify-between">
                      <div className="min-w-0">
                        <p className="text-sm text-muted-foreground">{stat.title}</p>
                        <p className="mt-2 text-2xl font-bold text-foreground font-inter tabular-nums truncate">{stat.value}</p>
                      </div>
                      <span className={`inline-flex items-center justify-center h-11 w-11 rounded-xl shrink-0 ${stat.iconClass}`}>
                        <stat.icon className="h-5 w-5" />
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Receipt className="h-5 w-5 text-primary" />
                  Transaksi Terbaru
                </CardTitle>
              </CardHeader>
              <CardContent>
                {loading ? (
                  <Skeleton className="h-40 rounded-xl" />
                ) : recentTxns.length === 0 ? (
                  <EmptyState title="Belum ada data" description="Transaksi terbaru akan muncul di sini." />
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow className="hover:bg-transparent">
                        <TableHead>Order ID</TableHead>
                        <TableHead>Jumlah</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {recentTxns.slice(0, 5).map((txn) => (
                        <TableRow key={txn.id}>
                          <TableCell className="font-medium truncate max-w-[140px]">{txn.order_id}</TableCell>
                          <TableCell>{formatCurrency(txn.gross_amount)}</TableCell>
                          <TableCell>
                            <Badge variant={txnStatusVariant(txn.status)}>{txnStatusLabel(txn.status)}</Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Users className="h-5 w-5 text-primary" />
                  Pengguna Terbaru
                </CardTitle>
              </CardHeader>
              <CardContent>
                {loading ? (
                  <Skeleton className="h-40 rounded-xl" />
                ) : recentUsers.length === 0 ? (
                  <EmptyState title="Belum ada data" description="Pengguna terbaru akan muncul di sini." />
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow className="hover:bg-transparent">
                        <TableHead>Nama</TableHead>
                        <TableHead>Email</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {recentUsers.slice(0, 5).map((u) => (
                        <TableRow key={u.id}>
                          <TableCell className="font-medium truncate max-w-[140px]">{u.name}</TableCell>
                          <TableCell className="text-muted-foreground truncate max-w-[160px]">{u.email}</TableCell>
                          <TableCell>
                            <Badge variant={userStatusVariant(u.status)}>{u.status}</Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </div>

          <p className="text-xs text-muted-foreground flex items-center gap-1">
            Periode analytics: {analytics?.period_start ? formatDate(analytics.period_start) : "-"} – {analytics?.period_end ? formatDate(analytics.period_end) : "-"}
            <ArrowUpRight className="h-3 w-3" />
          </p>
        </div>
    </AdminProtectedRoute>
  )
}
