"use client"

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react"
import { setCookie, parseCookies, destroyCookie } from "nookies"
import { apiPost } from "@/lib/api"
import type { UserInfo, LoginResponse, TokenValidationResponse } from "@/lib/types"

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: UserInfo | null;
  token: string | null;
  login: (email: string, senha: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

const COOKIE_TOKEN = 'authToken';
const COOKIE_MAX_AGE = 60 * 60 * 24; // 24 horas

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<UserInfo | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // Ao carregar o app, valida o token salvo no cookie
  useEffect(() => {
    const hydrate = async () => {
      const cookies = parseCookies();
      const storedToken = cookies[COOKIE_TOKEN];

      if (!storedToken) {
        setIsLoading(false);
        return;
      }

      try {
        const validation = await apiPost<TokenValidationResponse>('/auth/validate', {
          token: storedToken,
        });
        if (validation.valid) {
          setToken(storedToken);
          setUser(validation.user);
          setIsAuthenticated(true);
        } else {
          destroyCookie(undefined, COOKIE_TOKEN, { path: '/' });
        }
      } catch {
        destroyCookie(undefined, COOKIE_TOKEN, { path: '/' });
      } finally {
        setIsLoading(false);
      }
    };

    hydrate();
  }, []);

  const login = async (email: string, senha: string): Promise<void> => {
    const data = await apiPost<LoginResponse>('/auth/login', { email, senha });

    setCookie(undefined, COOKIE_TOKEN, data.token, {
      maxAge: COOKIE_MAX_AGE,
      path: '/',
    });

    setToken(data.token);
    setUser(data.user);
    setIsAuthenticated(true);
  };

  const logout = async (): Promise<void> => {
    try {
      await apiPost('/auth/logout', {});
    } catch {
      // Ignora erro no logout — limpa o estado local de qualquer forma
    } finally {
      destroyCookie(undefined, COOKIE_TOKEN, { path: '/' });
      setToken(null);
      setUser(null);
      setIsAuthenticated(false);
    }
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, token, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

