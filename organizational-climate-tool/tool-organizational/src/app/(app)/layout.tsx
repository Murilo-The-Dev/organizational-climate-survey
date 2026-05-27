import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { AuthGuard } from "@/components/AuthGuard";
import React from "react";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <AuthGuard>
      <DashboardLayout>{children}</DashboardLayout>
    </AuthGuard>
  );
}

