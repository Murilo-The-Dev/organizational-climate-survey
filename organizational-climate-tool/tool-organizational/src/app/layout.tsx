import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "../styles/globals.css";
import DashboardLayout from "@/components/dashboard/DashboardLayout";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

const interMono = Inter({
  variable: "--font-inter-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Atmos",
  description: "Monitorando o clima organizacional da sua empresa",
  icons: {
    icon: "./public/images/logoAtmos.svg",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body
        className={`${inter.variable} ${interMono.variable} antialiased bg-zinc-50`}
      >
      <DashboardLayout>{children}</DashboardLayout>
      </body>
    </html>
  );
}
