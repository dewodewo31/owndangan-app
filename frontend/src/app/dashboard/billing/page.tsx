"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { CreditCard, Receipt, ArrowRight, Crown } from "lucide-react"
import ProtectedRoute from "@/components/dashboard/protected-route"
import DashboardLayout from "@/components/dashboard/dashboard-layout"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
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
import { formatCurrency, formatDate } from "@/lib/utils"

interface Subscription {
  id: string
  package: { name: string; code: string; price: number }
  status: string
  start_at: string
  expires_at?: string
}

interface Transaction {
  id: string
  order_id: string
  gross_amount: number
  status: string
  payment_type: string
  created_at: string
}

function getStatusVariant(status: string) {
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

function getStatusLabel(status: string) {
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

export default function BillingPage() {
  const [subscription, setSubscription] = useState<Subscription | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchBilling()
  }, [])

  const fetchBilling = async () => {
    try {
      const [subsRes, txnsRes] = await Promise.all([
        api.get("/subscriptions/current"),
        api.get("/payments/transactions"),
      ])
      setSubscription(subsRes.data?.data)
      setTransactions(txnsRes.data?.data || [])
    } catch (error) {
      console.error("Failed to fetch billing:", error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="space-y-8">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-foreground tracking-tight">Billing</h1>
            <p className="mt-1 text-muted-foreground">
              Kelola langganan dan riwayat transaksi.
            </p>
          </div>

          {loading ? (
            <div className="space-y-6">
              <Skeleton className="h-36 rounded-2xl" />
              <Skeleton className="h-64 rounded-2xl" />
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <Card className="lg:col-span-2">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <CreditCard className="h-5 w-5 text-primary" />
                      Langganan Saat Ini
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    {subscription ? (
                      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                        <div className="flex items-center gap-4">
                          <span className="inline-flex items-center justify-center h-14 w-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-violet-600 text-white shadow-lg shadow-indigo-500/20">
                            <Crown className="h-6 w-6" />
                          </span>
                          <div>
                            <p className="text-xl font-bold text-foreground">
                              Paket {subscription.package?.name || "Free"}
                            </p>
                            <p className="text-sm text-muted-foreground mt-0.5">
                              Mulai {formatDate(subscription.start_at)}
                              {subscription.expires_at && ` · Berakhir ${formatDate(subscription.expires_at)}`}
                            </p>
                          </div>
                        </div>
                        <Badge variant={getStatusVariant(subscription.status)}>
                          {getStatusLabel(subscription.status)}
                        </Badge>
                      </div>
                    ) : (
                      <p className="text-muted-foreground">Belum ada langganan aktif.</p>
                    )}
                  </CardContent>
                </Card>

                <Card className="bg-gradient-to-br from-indigo-600 to-violet-700 text-white border-0">
                  <CardContent className="p-6 flex flex-col justify-between h-full">
                    <div>
                      <h3 className="font-semibold text-lg">Upgrade Paket</h3>
                      <p className="text-sm text-indigo-100 mt-1.5 leading-relaxed">
                        Buka lebih banyak fitur: lebih banyak tamu, custom domain,
                        dan template eksklusif.
                      </p>
                    </div>
                    <Link href="/packages" className="mt-6">
                      <Button className="w-full bg-white text-primary hover:bg-indigo-50 border-0">
                        Lihat Paket
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </Button>
                    </Link>
                  </CardContent>
                </Card>
              </div>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Receipt className="h-5 w-5 text-primary" />
                    Riwayat Transaksi
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {transactions.length === 0 ? (
                    <div className="text-center py-12">
                      <div className="mx-auto w-14 h-14 rounded-full bg-muted flex items-center justify-center mb-4">
                        <Receipt className="h-6 w-6 text-muted-foreground" />
                      </div>
                      <p className="text-muted-foreground">Belum ada transaksi.</p>
                    </div>
                  ) : (
                    <Table>
                      <TableHeader>
                        <TableRow className="hover:bg-transparent">
                          <TableHead>Order ID</TableHead>
                          <TableHead>Jumlah</TableHead>
                          <TableHead>Status</TableHead>
                          <TableHead>Metode</TableHead>
                          <TableHead>Tanggal</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {transactions.map((txn) => (
                          <TableRow key={txn.id}>
                            <TableCell className="font-medium">{txn.order_id}</TableCell>
                            <TableCell>{formatCurrency(txn.gross_amount)}</TableCell>
                            <TableCell>
                              <Badge variant={getStatusVariant(txn.status)}>
                                {getStatusLabel(txn.status)}
                              </Badge>
                            </TableCell>
                            <TableCell className="text-muted-foreground">{txn.payment_type || "-"}</TableCell>
                            <TableCell className="text-muted-foreground">{formatDate(txn.created_at)}</TableCell>
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
