"use client"

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react"
import { setCookie, parseCookies, destroyCookie } from "nookies"

interface AuthContextType {
  // define os tipos aqui
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)


export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<{ email: string } | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    const cookies = parseCookies();
    const storedToken = cookies.authToken;
    const storedUser = cookies.authUser ? JSON.parse(cookies.authUser) : null;

    if (storedToken && storedUser) {
      setToken(storedToken);
      setUser(storedUser);
      setIsAuthenticated(true);
    }
    setIsLoading(false);
  }, []);

  const login = (email: string, newToken: string) => {
    setCookie(undefined, 'authToken', newToken, { maxAge: 60 * 60 * 24 * 7, path: '/' }); // 7 days
    setCookie(undefined, 'authUser', JSON.stringify({ email }), { maxAge: 60 * 60 * 24 * 7, path: '/' });
    setToken(newToken);
    setUser({ email });
    setIsAuthenticated(true);
  };

  const logout = () => {
    destroyCookie(undefined, 'authToken', { path: '/' });
    destroyCookie(undefined, 'authUser', { path: '/' });
    setToken(null);
    setUser(null);
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, token, login, logout, isLoading }}>
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

