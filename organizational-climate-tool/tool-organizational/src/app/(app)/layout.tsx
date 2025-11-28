import DashboardLayout from "@/components/dashboard/DashboardLayout";
import React from "react";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <DashboardLayout>{children}</DashboardLayout>
  );
}

