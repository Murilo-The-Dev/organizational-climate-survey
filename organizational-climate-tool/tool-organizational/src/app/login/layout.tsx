import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "../../styles/globals.css";

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
    <html lang="pt-BR">
      <body className={`${inter.variable} ${interMono.variable} antialiased bg-zinc-50`}>
        {children}
      </body>
    </html>
  );
}
