import { useEffect, useState, useCallback } from 'react'
import useSWR from 'swr'
import api from './api'
import type { ApiResponse } from './types'

export function useAuth() {
  return useSWR('/auth/me', (url: string) => api.get<ApiResponse<unknown>>(url).then((res) => res.data).catch(() => null))
}

export function usePackages() {
  return useSWR('/packages', (url: string) =>
    api.get<ApiResponse<unknown>>(url).then((res) => res.data)
  )
}

export function useLocalStorage<T>(key: string, initialValue: T) {
  const [storedValue, setStoredValue] = useState<T>(initialValue)

  const updateStoredValue = useCallback(
    (value: T | ((v: T) => T)) => {
      setStoredValue((prev) => {
        const v = typeof value === 'function' ? (value as (v: T) => T)(prev) : value
        if (typeof window !== 'undefined') {
          localStorage.setItem(key, JSON.stringify(v))
        }
        return v
      })
    },
    [key]
  )

  return [storedValue, updateStoredValue] as const
}
