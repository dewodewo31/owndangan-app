"use client"

const SNAP_SANDBOX_URL = "https://app.sandbox.midtrans.com/snap/snap.js"
const SNAP_PRODUCTION_URL = "https://app.midtrans.com/snap/snap.js"

export interface SnapResult {
  order_id: string
  transaction_status: string
  transaction_id?: string
  status_code?: string
  status_message?: string
  payment_type?: string
  [key: string]: unknown
}

export interface SnapCallbacks {
  onSuccess?: (result: SnapResult) => void
  onPending?: (result: SnapResult) => void
  onError?: (result: SnapResult) => void
  onClose?: () => void
}

declare global {
  interface Window {
    snap?: {
      pay: (
        token: string,
        callbacks?: {
          onSuccess?: (result: SnapResult) => void
          onPending?: (result: SnapResult) => void
          onError?: (result: SnapResult) => void
          onClose?: () => void
        }
      ) => void
    }
  }
}

function getSnapScriptUrl(): string {
  const isProduction = process.env.NEXT_PUBLIC_MIDTRANS_IS_PRODUCTION === "true"
  return isProduction ? SNAP_PRODUCTION_URL : SNAP_SANDBOX_URL
}

let snapScriptPromise: Promise<void> | null = null

export function loadSnapScript(): Promise<void> {
  if (typeof window === "undefined") {
    return Promise.reject(new Error("Snap.js can only be loaded in the browser"))
  }
  if (window.snap) {
    return Promise.resolve()
  }
  if (snapScriptPromise) {
    return snapScriptPromise
  }

  snapScriptPromise = new Promise<void>((resolve, reject) => {
    const clientKey = process.env.NEXT_PUBLIC_MIDTRANS_CLIENT_KEY || ""
    const script = document.createElement("script")
    script.src = getSnapScriptUrl()
    script.setAttribute("data-client-key", clientKey)
    script.async = true
    script.onload = () => {
      if (window.snap) {
        resolve()
      } else {
        reject(new Error("Midtrans Snap failed to initialize"))
      }
    }
    script.onerror = () => {
      snapScriptPromise = null
      reject(new Error("Failed to load Midtrans Snap.js"))
    }
    document.body.appendChild(script)
  })

  return snapScriptPromise
}

export async function payWithSnap(
  token: string,
  callbacks: SnapCallbacks = {}
): Promise<void> {
  await loadSnapScript()
  if (!window.snap) {
    throw new Error("Midtrans Snap not available")
  }
  const snap = window.snap
  return new Promise<void>((resolve) => {
    snap.pay(token, {
      onSuccess: (result: SnapResult) => {
        callbacks.onSuccess?.(result)
        resolve()
      },
      onPending: (result: SnapResult) => {
        callbacks.onPending?.(result)
        resolve()
      },
      onError: (result: SnapResult) => {
        callbacks.onError?.(result)
        resolve()
      },
      onClose: () => {
        callbacks.onClose?.()
        resolve()
      },
    })
  })
}
