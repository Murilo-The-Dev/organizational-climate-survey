import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "../../styles/globals.css";

// Definições de fontes (mantidas)
const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

const interMono = Inter({
  variable: "--font-inter-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Atmos - Login",
  description: "Página de Login da plataforma Atmos",
};

export default function LoginLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    // Adicionado suppressHydrationWarning aqui
    <html lang="pt-BR" suppressHydrationWarning> 
      {/* Adicionado suppressHydrationWarning aqui */}
      <body className="antialiased bg-zinc-50" suppressHydrationWarning>
        {children}
      </body>
    </html>
  );
}