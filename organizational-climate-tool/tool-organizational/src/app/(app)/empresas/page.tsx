"use client";

import { useEffect, useState } from "react";
import { DataTable } from "@/components/dashboard/DataTable";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { toast } from "sonner";
import { Loader2, MoreHorizontal } from "lucide-react";
import { empresaService } from "@/lib/services/empresaService"; 
import { Empresa } from "@/lib/types";

export default function EmpresasPage() {
  const [empresas, setEmpresas] = useState<Empresa[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const columns = [
    { accessorKey: "name", header: "Nome da Empresa" },
    { accessorKey: "cnpj", header: "CNPJ" },
    { accessorKey: "status", header: "Status" },
    {
      id: "actions",
      header: "Ações",
      cell: ({ row }: any) => (
        <Button variant="ghost" className="h-8 w-8 p-0">
          <span className="sr-only">Abrir menu</span>
          <MoreHorizontal className="h-4 w-4" />
        </Button>
      ),
    },
  ];

  useEffect(() => {
    carregarEmpresas();
  }, []);

  const carregarEmpresas = async () => {
    try {
      setIsLoading(true);
      const data = await empresaService.list(); 
      setEmpresas(data);
    } catch (error) {
      console.error("Erro ao carregar empresas:", error);
      toast.error("Não foi possível carregar a lista de empresas.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Empresas
        </h1>
        <Link href="/empresas/nova">
          <Button>Adicionar Nova Empresa</Button>
        </Link>
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie as empresas cadastradas no sistema.
      </p>
      
      <div className="bg-background rounded-lg border p-4 h-full min-h-[400px]">
        {isLoading ? (
          <div className="flex h-full w-full items-center justify-center pt-20">
            <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
          </div>
        ) : (
          <DataTable columns={columns} data={empresas} />
        )}
      </div>
    </section>
  );
}