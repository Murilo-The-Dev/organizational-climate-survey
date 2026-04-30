import axios, { AxiosError } from "axios";
import { parseCookies } from "nookies";
import { toast } from "sonner";

// Tipo base que espelha a estrutura padrão de resposta da API Go
export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data: T;
  error?: string;
};

const api = axios.create({
  baseURL:
    process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1",
});

// Interceptor de request: injeta o token JWT em todas as requisições
api.interceptors.request.use((config) => {
  const { authToken: token } = parseCookies();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor de response: trata 401 e erros genéricos
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse<unknown>>) => {
    if (error.response?.status === 401) {
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
      return Promise.reject(error);
    }

    const message =
      error.response?.data?.message ||
      error.response?.data?.error ||
      "Ocorreu um erro inesperado. Tente novamente.";

    if (typeof window !== "undefined") {
      toast.error(message);
    }

    return Promise.reject(error);
  },
);

// Helpers tipados que unwrapam `data` automaticamente
export async function apiGet<T>(
  url: string,
  params?: Record<string, unknown>,
): Promise<T> {
  const response = await api.get<ApiResponse<T>>(url, { params });
  return response.data.data;
}

export async function apiPost<T>(url: string, body?: unknown): Promise<T> {
  const response = await api.post<ApiResponse<T>>(url, body);
  return response.data.data;
}

export async function apiPut<T>(url: string, body?: unknown): Promise<T> {
  const response = await api.put<ApiResponse<T>>(url, body);
  return response.data.data;
}

export async function apiDelete<T = void>(url: string): Promise<T> {
  const response = await api.delete<ApiResponse<T>>(url);
  return response.data.data;
}

export default api;
