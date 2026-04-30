"use client";

import { useEffect, useState } from "react";
import { DataTable } from "@/components/dashboard/DataTable";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { toast } from "sonner";
import { Loader2, MoreHorizontal } from "lucide-react";
import { setorService } from "@/lib/services/setorService";
import { useAuth } from "@/context/AuthContext";
import type { Setor } from "@/lib/types";

export default function SetoresPage() {
  const { user } = useAuth();
  const [setores, setSetores] = useState<Setor[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Colunas da tabela
  const columns = [
    { accessorKey: "nome_setor", header: "Nome do Setor" },
    { accessorKey: "descricao", header: "Descrição" },
    {
      id: "actions",
      header: "Ações",
      cell: () => (
        <Button variant="ghost" className="h-8 w-8 p-0">
          <span className="sr-only">Abrir menu</span>
          <MoreHorizontal className="h-4 w-4" />
        </Button>
      ),
    },
  ];

  useEffect(() => {
    if (user) {
      carregarSetores();
    }
  }, [user]);

  const carregarSetores = async () => {
    try {
      setIsLoading(true);
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;
      const data = await setorService.listByEmpresa(empresaId);
      setSetores(data);
    } catch (error) {
      console.error("Erro ao carregar setores:", error);
      toast.error("Não foi possível carregar a lista de setores.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Setores
        </h1>
        <Link href="/setores/nova">
          <Button>Adicionar Novo Setor</Button>
        </Link>
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie os setores cadastrados na sua empresa.
      </p>

      <div className="bg-background rounded-lg border p-4 h-full min-h-[400px]">
        {isLoading ? (
          <div className="flex h-full w-full items-center justify-center pt-20">
            <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
          </div>
        ) : (
          <DataTable columns={columns as any} data={setores} />
        )}
      </div>
    </section>
  );
}
