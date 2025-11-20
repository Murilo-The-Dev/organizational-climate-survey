import { DataTable } from "@/components/dashboard/DataTable";
import { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger, DropdownMenuSeparator } from "@/components/ui/dropdown-menu";
import { Badge, badgeVariants } from "@/components/ui/badge";
import type { VariantProps } from "class-variance-authority";

export type AuditLog = {
  id: string;
  timestamp: string;
  user: string;
  action: string;
  entity: string;
  entityId: string;
  status: "success" | "failed" | "info";
};

const mockAuditLogs: AuditLog[] = [
  { id: "LOG-001", timestamp: "2025-10-12 10:00:00", user: "admin@example.com", action: "LOGIN", entity: "Auth", entityId: "N/A", status: "success" },
  { id: "LOG-002", timestamp: "2025-10-12 10:05:15", user: "user1@example.com", action: "CREATE_SURVEY", entity: "Survey", entityId: "SURV-007", status: "success" },
  { id: "LOG-003", timestamp: "2025-10-12 10:10:30", user: "admin@example.com", action: "UPDATE_COMPANY", entity: "Company", entityId: "COMP-001", status: "success" },
  { id: "LOG-004", timestamp: "2025-10-12 10:15:00", user: "user2@example.com", action: "LOGIN", entity: "Auth", entityId: "N/A", status: "failed" },
  { id: "LOG-005", timestamp: "2025-10-12 10:20:45", user: "admin@example.com", action: "DELETE_USER", entity: "User", entityId: "USER-003", status: "info" },
  { id: "LOG-006", timestamp: "2025-10-12 10:25:00", user: "user1@example.com", action: "RESPOND_SURVEY", entity: "Survey", entityId: "SURV-001", status: "success" },
  { id: "LOG-007", timestamp: "2025-10-12 10:30:00", user: "admin@example.com", action: "VIEW_REPORT", entity: "Report", entityId: "REP-001", status: "success" },
  { id: "LOG-008", timestamp: "2025-10-12 10:35:00", user: "user3@example.com", action: "CREATE_COMPANY", entity: "Company", entityId: "COMP-005", status: "failed" },
];

const columns: ColumnDef<AuditLog>[] = [
  { accessorKey: "timestamp", header: "Data/Hora" },
  { accessorKey: "user", header: "Usuário" },
  { accessorKey: "action", header: "Ação" },
  { accessorKey: "entity", header: "Entidade" },
  { accessorKey: "entityId", header: "ID da Entidade" },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status");
      let badgeVariant: VariantProps<typeof badgeVariants>["variant"];
      let badgeText = "";
      switch (status) {
        case "success":
          badgeVariant = "success";
          badgeText = "Sucesso";
          break;
        case "failed":
          badgeVariant = "destructive";
          badgeText = "Falha";
          break;
        case "info":
          badgeVariant = "default";
          badgeText = "Info";
          break;
        default:
          badgeText = String(status);
          badgeVariant = "secondary";
      }
      return <Badge variant={badgeVariant}>{badgeText}</Badge>;
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
            <DropdownMenuLabel className="text-center">Ações</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => navigator.clipboard.writeText(log.id)}>
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

      <div className="bg-background rounded-lg border p-4 h-full">
        <DataTable columns={columns} data={mockAuditLogs} />
      </div>
    </section>
  );
}
