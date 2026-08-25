"use client"

import { useCallback, useEffect, useState } from "react"
import Link from "next/link"
import { CreditCard, Receipt, ArrowRight, Crown, Sparkles, Loader2, AlertCircle } from "lucide-react"
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
import { payWithSnap } from "@/lib/midtrans"
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

interface Pkg {
  id: string
  name: string
  code: string
  price: number
  duration_days?: number | null
  guest_limit?: number | null
  features?: Record<string, unknown> | unknown[] | null
  is_active: boolean
}

const POLL_INTERVAL = 3000
const POLL_MAX = 10

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

function featureList(features?: Record<string, unknown> | unknown[] | null): string[] {
  if (!features) return []
  if (Array.isArray(features)) {
    return features.filter((f) => typeof f === "string").map((f) => f as string)
  }
  return Object.values(features).filter((v) => typeof v === "string").map((v) => v as string)
}

export default function BillingPage() {
  const [subscription, setSubscription] = useState<Subscription | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [packages, setPackages] = useState<Pkg[]>([])
  const [loading, setLoading] = useState(true)
  const [checkingOut, setCheckingOut] = useState<string | null>(null)
  const [checkoutError, setCheckoutError] = useState<string | null>(null)

  useEffect(() => {
    fetchBilling()
    fetchPackages()
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

  const fetchPackages = async () => {
    try {
      const res = await api.get("/packages")
      setPackages(res.data?.data || [])
    } catch (error) {
      console.error("Failed to fetch packages:", error)
    }
  }

  // After a Snap payment resolves, the webhook lands asynchronously; poll the
  // backend so the UI reflects settlement/activation without a manual refresh.
  const pollStatus = useCallback(async (orderId: string) => {
    for (let i = 0; i < POLL_MAX; i++) {
      await new Promise((r) => setTimeout(r, POLL_INTERVAL))
      try {
        const [subRes, txnsRes] = await Promise.all([
          api.get("/subscriptions/current"),
          api.get("/payments/transactions"),
        ])
        const sub = subRes.data?.data
        const txns: Transaction[] = txnsRes.data?.data || []
        setSubscription(sub || null)
        setTransactions(txns)
        if (sub && sub.status === "active") return
        const myTxn = txns.find((t) => t.order_id === orderId)
        if (
          myTxn &&
          (myTxn.status === "settlement" ||
            myTxn.status === "expire" ||
            myTxn.status === "deny" ||
            myTxn.status === "cancel")
        ) {
          return
        }
      } catch {
        // keep polling through transient errors
      }
    }
  }, [])

  const handleCheckout = async (pkg: Pkg) => {
    if (checkingOut) return
    setCheckingOut(pkg.id)
    setCheckoutError(null)
    try {
      const res = await api.post("/payments/snap", { package_id: pkg.id })
      const token = res.data?.data?.snap_token
      if (!token) throw new Error("No snap token returned")
      await payWithSnap(token, {
        onSuccess: (result) => pollStatus(result.order_id),
        onPending: (result) => pollStatus(result.order_id),
        onClose: () => fetchBilling(),
      })
    } catch (error: any) {
      const msg =
        error?.response?.data?.error?.message ||
        error?.message ||
        "Gagal memproses pembayaran. Coba lagi."
      setCheckoutError(msg)
      console.error("Checkout failed:", error)
    } finally {
      setCheckingOut(null)
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
                          <span className="inline-flex items-center justify-center h-14 w-14 rounded-2xl bg-gradient-to-br from-plum to-rosegold text-white shadow-elevation-1">
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

                <Card className="bg-gradient-to-br from-plum to-rosegold text-white border-0">
                  <CardContent className="p-6 flex flex-col justify-between h-full">
                    <div>
                      <h3 className="font-semibold text-lg">Upgrade Paket</h3>
                      <p className="text-sm text-primary-foreground/80 mt-1.5 leading-relaxed">
                        Buka lebih banyak fitur: lebih banyak tamu, custom domain,
                        dan template eksklusif.
                      </p>
                    </div>
                    <Link href="#pilih-paket" className="mt-6">
                      <Button className="w-full bg-white text-primary hover:bg-primary-container border-0">
                        Lihat Paket
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </Button>
                    </Link>
                  </CardContent>
                </Card>
              </div>

              <Card id="pilih-paket">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Sparkles className="h-5 w-5 text-primary" />
                    Pilih Paket
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {checkoutError && (
                    <div className="mb-4 flex items-start gap-2.5 rounded-xl border border-destructive/20 bg-destructive/10 p-3.5 text-sm text-destructive">
                      <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
                      <span>{checkoutError}</span>
                    </div>
                  )}
                  {packages.length === 0 ? (
                    <p className="text-muted-foreground">Tidak ada paket tersedia.</p>
                  ) : (
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                      {packages.map((pkg) => {
                        const features = featureList(pkg.features)
                        const isCurrent =
                          subscription?.status === "active" &&
                          subscription?.package?.code === pkg.code
                        return (
                          <div
                            key={pkg.id}
                            className="rounded-2xl border border-border p-5 flex flex-col"
                          >
                            <p className="font-bold text-lg text-foreground">{pkg.name}</p>
                            <p className="text-2xl font-extrabold mt-2 text-foreground">
                              {formatCurrency(pkg.price)}
                            </p>
                            <p className="text-sm text-muted-foreground mt-1">
                              {pkg.duration_days ? `${pkg.duration_days} hari` : "Lifetime"}
                              {pkg.guest_limit ? ` · ${pkg.guest_limit} tamu` : ""}
                            </p>
                            {features.length > 0 && (
                              <ul className="mt-3 space-y-1 text-sm text-muted-foreground">
                                {features.map((f, i) => (
                                  <li key={i}>• {f}</li>
                                ))}
                              </ul>
                            )}
                            <Button
                              className="mt-5 w-full"
                              disabled={checkingOut !== null}
                              variant={isCurrent ? "outline" : "primary"}
                              onClick={() => handleCheckout(pkg)}
                            >
                              {checkingOut === pkg.id ? (
                                <>
                                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                  Memproses…
                                </>
                              ) : isCurrent ? (
                                "Paket Aktif"
                              ) : (
                                "Berlangganan"
                              )}
                            </Button>
                          </div>
                        )
                      })}
                    </div>
                  )}
                </CardContent>
              </Card>

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
