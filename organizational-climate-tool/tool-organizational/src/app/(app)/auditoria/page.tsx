"use client";

import { useEffect, useState } from "react";
import { DataTable } from "@/components/dashboard/DataTable";
import { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge"; // Mudamos para o Badge comum
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";
import { auditoriaService } from "@/lib/services/auditoriaService";
import { LogAuditoria } from "@/lib/types"; // Importando o tipo real!

// Atualizamos as colunas para bater com os nomes em português que costumam vir do backend
const columns: ColumnDef<LogAuditoria>[] = [
  {
    accessorKey: "data_hora",
    header: "Data/Hora",
    cell: ({ row }) => {
      const valor = row.getValue("data_hora") as string;
      // Formata a data se ela existir, ou mostra traço
      if (!valor) return "-";
      return new Date(valor).toLocaleString("pt-BR");
    },
  },
  { accessorKey: "id_usuario", header: "ID Usuário" },
  { accessorKey: "acao", header: "Ação" },
  { accessorKey: "entidade", header: "Entidade" },
  { accessorKey: "id_entidade", header: "ID Entidade" },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status") as string;

      // Lógica simples para colorir o Badge comum baseado no texto do status
      let variant: "default" | "destructive" | "outline" | "secondary" =
        "default";
      if (status?.toLowerCase() === "erro" || status?.toLowerCase() === "falha")
        variant = "destructive";
      if (status?.toLowerCase() === "info") variant = "secondary";

      return <Badge variant={variant}>{status || "N/A"}</Badge>;
    },
  },
  {
    id: "actions",
    header: "Ações",
    cell: ({ row }) => {
      const log = row.original;
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Abrir menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Ações</DropdownMenuLabel>
            {/* O ID real que vem do banco geralmente é id_log ou id */}
            <DropdownMenuItem
              onClick={() =>
                navigator.clipboard.writeText(
                  String((log as any).id || (log as any).id_log),
                )
              }
            >
              Copiar ID do Log
            </DropdownMenuItem>
            <DropdownMenuItem>Ver Detalhes</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];

export default function AuditoriaPage() {
  const { user } = useAuth();
  const [logs, setLogs] = useState<LogAuditoria[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (user) {
      carregarLogs();
    }
  }, [user]);

  const carregarLogs = async () => {
    try {
      setIsLoading(true);
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;

      const data = await auditoriaService.listByEmpresa(empresaId);

      // Previne erros se a API retornar vazio
      setLogs(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error("Erro ao carregar auditoria:", error);
      toast.error("Não foi possível carregar os logs de auditoria.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Painel de Auditoria
        </h1>
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Visualize e filtre os logs de atividades do sistema.
      </p>

      <div className="bg-background rounded-lg border p-4 h-full min-h-[400px]">
        {isLoading ? (
          <div className="flex h-full w-full items-center justify-center pt-20">
            <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
          </div>
        ) : (
          <DataTable columns={columns} data={logs} />
        )}
      </div>
    </section>
  );
}
