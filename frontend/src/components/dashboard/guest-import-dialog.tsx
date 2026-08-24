"use client"

import { useState, useRef } from "react"
import { Upload, AlertCircle, CheckCircle2, FileSpreadsheet, RefreshCw } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Select } from "@/components/ui/select"
import { Label } from "@/components/ui/label"
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

export interface ImportPreviewRow {
  index: number
  name: string
  email: string
  phone: string
  category: string
  status: "valid" | "duplicate" | "invalid"
  errors: string[]
}

export interface ImportPreviewResponse {
  columns: string[]
  rows: ImportPreviewRow[]
  summary: { total: number; valid: number; duplicate: number; invalid: number }
}

export interface ImportConfirmResponse {
  total: number
  imported: number
  duplicates: number
  errors: { index: number; errors: string[] }[]
}

const MAX_SIZE = 5 * 1024 * 1024
const MAPPING_FIELDS = ["name", "email", "phone", "category"] as const
type MappingField = (typeof MAPPING_FIELDS)[number]
const MAPPING_LABELS: Record<MappingField, string> = {
  name: "Nama",
  email: "Email",
  phone: "No. WhatsApp",
  category: "Kategori",
}

interface ApiError {
  code: string
  message: string
}

interface ApiErrorResponse {
  response?: { data?: { error?: ApiError } }
}

function extractError(err: unknown): ApiError {
  const data = (err as ApiErrorResponse)?.response?.data
  return data?.error ?? { code: "UNKNOWN", message: "Terjadi kesalahan. Coba lagi." }
}

interface GuestImportDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  eventId: string
  onImported: () => void
}

export function GuestImportDialog({ open, onOpenChange, eventId, onImported }: GuestImportDialogProps) {
  const [step, setStep] = useState<"upload" | "preview" | "done">("upload")
  const [file, setFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<ImportPreviewResponse | null>(null)
  const [mapping, setMapping] = useState<Record<MappingField, string>>({
    name: "",
    email: "",
    phone: "",
    category: "",
  })
  const [loadingPreview, setLoadingPreview] = useState(false)
  const [loadingImport, setLoadingImport] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [confirmResult, setConfirmResult] = useState<ImportConfirmResponse | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const reset = () => {
    setStep("upload")
    setFile(null)
    setPreview(null)
    setMapping({ name: "", email: "", phone: "", category: "" })
    setError(null)
    setConfirmResult(null)
    if (fileInputRef.current) fileInputRef.current.value = ""
  }

  const handleClose = (next: boolean) => {
    if (!next) reset()
    onOpenChange(next)
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setError(null)
    const f = e.target.files?.[0]
    if (!f) return
    if (!f.name.toLowerCase().endsWith(".csv")) {
      setError("Hanya file CSV yang didukung (.xlsx tidak tersedia).")
      return
    }
    if (f.size > MAX_SIZE) {
      setError("Ukuran file melebihi batas 5MB.")
      return
    }
    setFile(f)
  }

  const runPreview = async (f: File, m: Record<MappingField, string>) => {
    setLoadingPreview(true)
    setError(null)
    try {
      const form = new FormData()
      form.append("file", f)
      const active = MAPPING_FIELDS.filter((k) => m[k]).reduce<Record<string, string>>(
        (acc, k) => {
          acc[k] = m[k]
          return acc
        },
        {}
      )
      if (Object.keys(active).length > 0) {
        form.append("mapping", JSON.stringify(active))
      }
      const res = await api.post<{ data: ImportPreviewResponse }>(
        `/events/${eventId}/guests/import`,
        form,
        { headers: { "Content-Type": "multipart/form-data" } }
      )
      setPreview(res.data.data)
      setStep("preview")
    } catch (err) {
      const e = extractError(err)
      setError(e.message || "Gagal memproses file.")
    } finally {
      setLoadingPreview(false)
    }
  }

  const handleUpload = () => {
    if (!file) return
    void runPreview(file, mapping)
  }

  const handleMappingChange = (field: MappingField, value: string) => {
    const next = { ...mapping, [field]: value }
    setMapping(next)
    if (file) void runPreview(file, next)
  }

  const validRows = preview?.rows.filter((r) => r.status === "valid") ?? []

  const handleImport = async () => {
    if (validRows.length === 0) return
    setLoadingImport(true)
    setError(null)
    try {
      const res = await api.post<{ data: ImportConfirmResponse }>(
        `/events/${eventId}/guests/import/confirm`,
        {
          rows: validRows.map((r) => ({
            name: r.name,
            email: r.email,
            phone: r.phone,
            category: r.category,
          })),
        }
      )
      setConfirmResult(res.data.data)
      setStep("done")
      onImported()
    } catch (err) {
      const e = extractError(err)
      setError(e.message || "Gagal mengimpor tamu.")
    } finally {
      setLoadingImport(false)
    }
  }

  const statusBadge = (status: ImportPreviewRow["status"]) => {
    if (status === "valid")
      return <Badge variant="success">Valid</Badge>
    if (status === "duplicate")
      return <Badge variant="warning">Duplikat</Badge>
    return <Badge variant="destructive">Invalid</Badge>
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Impor Tamu dari CSV</DialogTitle>
          <DialogDescription>
            Unggah file CSV berisi daftar tamu. Kolom akan dideteksi otomatis.
          </DialogDescription>
        </DialogHeader>

        {error && (
          <div className="flex items-start gap-2 rounded-lg border border-destructive/30 bg-destructive/10 px-3.5 py-2.5 text-sm text-destructive">
            <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        {step === "upload" && (
          <div className="space-y-4">
            <label
              className={cn(
                "flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed border-input px-6 py-10 text-center transition-colors hover:border-primary/50 hover:bg-accent/40 cursor-pointer",
                file && "border-primary/50 bg-accent/40"
              )}
            >
              <input
                ref={fileInputRef}
                type="file"
                accept=".csv"
                className="hidden"
                onChange={handleFileChange}
              />
              {file ? (
                <>
                  <FileSpreadsheet className="h-8 w-8 text-primary" />
                  <div>
                    <p className="font-medium text-foreground">{file.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {(file.size / 1024).toFixed(1)} KB · Klik untuk ganti file
                    </p>
                  </div>
                </>
              ) : (
                <>
                  <Upload className="h-8 w-8 text-muted-foreground" />
                  <div>
                    <p className="font-medium text-foreground">Pilih file CSV</p>
                    <p className="text-xs text-muted-foreground">Maksimal 5MB · hanya .csv</p>
                  </div>
                </>
              )}
            </label>
          </div>
        )}

        {step === "preview" && preview && (
          <div className="space-y-4">
            <div className="grid grid-cols-4 gap-2">
              <SummaryStat label="Total" value={preview.summary.total} />
              <SummaryStat label="Valid" value={preview.summary.valid} tone="success" />
              <SummaryStat label="Duplikat" value={preview.summary.duplicate} tone="warning" />
              <SummaryStat label="Invalid" value={preview.summary.invalid} tone="destructive" />
            </div>

            <details className="rounded-lg border border-border bg-muted/30 px-4 py-3">
              <summary className="cursor-pointer text-sm font-medium text-foreground">
                Atur pemetaan kolom
              </summary>
              <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
                {MAPPING_FIELDS.map((field) => (
                  <div key={field} className="space-y-1.5">
                    <Label htmlFor={`map-${field}`}>{MAPPING_LABELS[field]}</Label>
                    <Select
                      id={`map-${field}`}
                      value={mapping[field]}
                      onChange={(e) => handleMappingChange(field, e.target.value)}
                      options={[
                        { value: "", label: "Otomatis" },
                        ...preview.columns.map((c) => ({ value: c, label: c })),
                      ]}
                    />
                  </div>
                ))}
              </div>
            </details>

            <div className="rounded-lg border border-border">
              <Table>
                <TableHeader>
                  <TableRow className="hover:bg-transparent">
                    <TableHead className="w-16">#</TableHead>
                    <TableHead>Nama</TableHead>
                    <TableHead>Kontak</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {preview.rows.map((row) => (
                    <TableRow
                      key={row.index}
                      className={cn(
                        row.status === "invalid" && "bg-destructive/5",
                        row.status === "duplicate" && "bg-warning/5"
                      )}
                    >
                      <TableCell className="text-muted-foreground">{row.index + 1}</TableCell>
                      <TableCell>
                        <p className="font-medium text-foreground">{row.name || "—"}</p>
                        {row.status === "invalid" && row.errors.length > 0 && (
                          <p className="text-xs text-destructive mt-0.5">{row.errors.join(", ")}</p>
                        )}
                      </TableCell>
                      <TableCell className="text-muted-foreground text-xs">
                        {[row.email, row.phone].filter(Boolean).join(" · ") || "—"}
                      </TableCell>
                      <TableCell>{statusBadge(row.status)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </div>
        )}

        {step === "done" && confirmResult && (
          <div className="space-y-4">
            <div className="flex items-center gap-3 rounded-xl bg-success/10 px-4 py-3">
              <CheckCircle2 className="h-6 w-6 text-success" />
              <div>
                <p className="font-medium text-foreground">Impor selesai</p>
                <p className="text-sm text-muted-foreground">
                  {confirmResult.imported} tamu berhasil diimpor
                  {confirmResult.duplicates > 0 && `, ${confirmResult.duplicates} duplikat dilewati`}.
                </p>
              </div>
            </div>
            <div className="grid grid-cols-3 gap-2">
              <SummaryStat label="Total" value={confirmResult.total} />
              <SummaryStat label="Diimpor" value={confirmResult.imported} tone="success" />
              <SummaryStat label="Duplikat" value={confirmResult.duplicates} tone="warning" />
            </div>
            {confirmResult.errors.length > 0 && (
              <div className="rounded-lg border border-destructive/30 bg-destructive/10 px-3.5 py-2.5 text-sm text-destructive">
                {confirmResult.errors.length} baris gagal:{" "}
                {confirmResult.errors.map((e) => `baris ${e.index + 1}`).join(", ")}
              </div>
            )}
          </div>
        )}

        <DialogFooter className="gap-2">
          {step === "upload" && (
            <>
              <Button variant="ghost" onClick={() => handleClose(false)}>
                Batal
              </Button>
              <Button onClick={handleUpload} loading={loadingPreview} disabled={!file}>
                Pratinjau
              </Button>
            </>
          )}
          {step === "preview" && (
            <>
              <Button
                variant="ghost"
                onClick={() => {
                  setStep("upload")
                  setPreview(null)
                }}
                disabled={loadingPreview || loadingImport}
              >
                <RefreshCw className="h-4 w-4 mr-2" />
                Pilih File Lain
              </Button>
              <Button
                onClick={handleImport}
                loading={loadingImport}
                disabled={validRows.length === 0}
              >
                Impor {validRows.length > 0 ? `(${validRows.length})` : ""}
              </Button>
            </>
          )}
          {step === "done" && (
            <Button onClick={() => handleClose(false)}>Selesai</Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function SummaryStat({
  label,
  value,
  tone,
}: {
  label: string
  value: number
  tone?: "success" | "warning" | "destructive"
}) {
  const toneClass =
    tone === "success"
      ? "text-success"
      : tone === "warning"
        ? "text-warning"
        : tone === "destructive"
          ? "text-destructive"
          : "text-foreground"
  return (
    <div className="rounded-lg border border-border bg-muted/30 px-3 py-2 text-center">
      <p className={cn("text-xl font-bold", toneClass)}>{value}</p>
      <p className="text-xs text-muted-foreground">{label}</p>
    </div>
  )
}
