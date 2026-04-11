"use client"

import { useState, useCallback } from 'react'
import { toast } from 'sonner'
import { AxiosError } from 'axios'
import type { ApiResponse } from '@/lib/api'

interface UseApiMutationOptions<TData> {
  onSuccess?: (data: TData) => void
  onError?: (error: string) => void
  successMessage?: string
}

interface UseApiMutationReturn<TData, TVariables> {
  mutate: (variables: TVariables) => Promise<TData | undefined>
  isLoading: boolean
  error: string | null
  reset: () => void
}

/**
 * Hook para encapsular chamadas à API com loading state, tratamento de erros
 * e integração com sonner. Funciona com react-hook-form via `error`.
 *
 * Erros 400 retornam a mensagem do backend para exibição no formulário.
 * Erros 500+ exibem toast genérico.
 */
export function useApiMutation<TData = unknown, TVariables = unknown>(
  mutationFn: (variables: TVariables) => Promise<TData>,
  options?: UseApiMutationOptions<TData>
): UseApiMutationReturn<TData, TVariables> {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const reset = useCallback(() => {
    setError(null)
  }, [])

  const mutate = useCallback(
    async (variables: TVariables): Promise<TData | undefined> => {
      setIsLoading(true)
      setError(null)

      try {
        const data = await mutationFn(variables)

        if (options?.successMessage) {
          toast.success(options.successMessage)
        }

        options?.onSuccess?.(data)
        return data
      } catch (err) {
        const axiosError = err as AxiosError<ApiResponse<unknown>>
        const status = axiosError.response?.status

        const message =
          axiosError.response?.data?.message ||
          axiosError.response?.data?.error ||
          'Ocorreu um erro inesperado.'

        // Erros de validação (400): popula o campo de erro no formulário
        // Erros de servidor (500+): já exibem toast via interceptor do api.ts
        if (status && status < 500) {
          setError(message)
        }

        options?.onError?.(message)
        return undefined
      } finally {
        setIsLoading(false)
      }
    },
    [mutationFn, options]
  )

  return { mutate, isLoading, error, reset }
}
