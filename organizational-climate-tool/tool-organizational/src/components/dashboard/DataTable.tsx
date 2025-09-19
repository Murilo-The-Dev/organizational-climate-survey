"use client";

import * as React from "react";
import {
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from "@tanstack/react-table";
import { ArrowUpDown, ChevronDown, MoreHorizontal } from "lucide-react";
import { StatusBadge } from "@/components/ui/status-badge";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const dadosPesquisas: Pesquisa[] = [
  { id: "PESQ-001", 
    titulo: "Engajamento Q1 2025", 
    status: "concluido", participantes: 152, 
    dataCriacao: "2025-03-28" 
  },
  { id: "PESQ-002", 
    titulo: "Feedback de Liderança H1", 
    status: "concluido", 
    participantes: 140, 
    dataCriacao: "2025-06-15" 
  },
  { id: "PESQ-003", 
    titulo: "Pesquisa de Satisfação Anual 2024", 
    status: "concluido", 
    participantes: 180, 
    dataCriacao: "2024-12-20"
   },
  { id: "PESQ-004", 
    titulo: "Clima Organizacional H2", 
    status: "em_andamento", 
    participantes: 125, 
    dataCriacao: "2025-09-01" 
  },
  { id: "PESQ-005", 
    titulo: "Onboarding Novos Contratados", 
    status: "em_andamento", 
    participantes: 25, 
    dataCriacao: "2025-09-10" 
  },
  { id: "PESQ-006", 
    titulo: "Avaliação de Benefícios", 
    status: "rascunho", 
    participantes: 0, 
    dataCriacao: "2025-09-18" 
  },
  { id: "PESQ-007", 
    titulo: "Engajamento Q2 2025", 
    status: "rascunho", 
    participantes: 0, 
    dataCriacao: "2025-09-15" 
  },
  { id: "PESQ-008", 
    titulo: "Segurança Psicológica", 
    status: "concluido", 
    participantes: 165, 
    dataCriacao: "2025-01-30" 
  },
  { id: "PESQ-009", 
    titulo: "Comunicação Interna", 
    status: "em_andamento", 
    participantes: 95, 
    dataCriacao: "2025-08-22" 
  },
  { id: "PESQ-010", 
    titulo: "Planejamento Estratégico 2026", 
    status: "rascunho",
     participantes: 0, 
     dataCriacao: "2025-09-19" 
    },
  { id: "PESQ-011", 
    titulo: "Ferramentas de Trabalho", 
    status: "concluido", 
    participantes: 170, 
    dataCriacao: "2025-05-10" 
  },
  { id: "PESQ-012", 
    titulo: "e-NPS Semestral", 
    status: "em_andamento", 
    participantes: 110, 
    dataCriacao: "2025-09-05" 
  },
  { id: "PESQ-003", 
    titulo: "Engajamento Q1 2025", 
    status: "concluido", participantes: 152, 
    dataCriacao: "2025-03-28" 
  },
  { id: "PESQ-004", 
    titulo: "Feedback de Liderança H1", 
    status: "concluido", 
    participantes: 140, 
    dataCriacao: "2025-06-15" 
  },
  { id: "PESQ-005", 
    titulo: "Pesquisa de Satisfação Anual 2024", 
    status: "concluido", 
    participantes: 180, 
    dataCriacao: "2024-12-20"
   },
  { id: "PESQ-006", 
    titulo: "Clima Organizacional H2", 
    status: "em_andamento", 
    participantes: 125, 
    dataCriacao: "2025-09-01" 
  },
  { id: "PESQ-007", 
    titulo: "Onboarding Novos Contratados", 
    status: "em_andamento", 
    participantes: 25, 
    dataCriacao: "2025-09-10" 
  },
  { id: "PESQ-008", 
    titulo: "Avaliação de Benefícios", 
    status: "rascunho", 
    participantes: 0, 
    dataCriacao: "2025-09-18" 
  },
  { id: "PESQ-009", 
    titulo: "Engajamento Q2 2025", 
    status: "concluido", 
    participantes: 0, 
    dataCriacao: "2025-09-15" 
  },
  { id: "PESQ-010", 
    titulo: "Segurança Psicológica", 
    status: "concluido", 
    participantes: 165, 
    dataCriacao: "2025-01-30" 
  },
];

export type Pesquisa = {
  id: string;
  titulo: string;
  status: "concluido" | "em_andamento" | "rascunho";
  participantes: number;
  dataCriacao: string;
};

export const columns: ColumnDef<Pesquisa>[] = [
  // Coluna de seleção (pode manter a mesma)
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
  },
  // Nova coluna 'Título da Pesquisa'
  {
    accessorKey: "titulo",
    header: "Título",
  },
  // Nova coluna 'Status' com o Badge
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      // Adapte o StatusBadge para os novos status
      return <StatusBadge status={row.getValue("status")} />;
    },
  },
  // Nova coluna 'Participantes'
  {
    accessorKey: "participantes",
    header: "Participantes",
  },
  // Nova coluna 'Data de Criação'
  {
    accessorKey: "dataCriacao",
    header: "Data de Criação",
  },
  // Coluna de Ações
  {
    id: "actions",
    cell: ({ row }) => {
      const pesquisa = row.original;
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuLabel>Ações</DropdownMenuLabel>
            <DropdownMenuItem>Ver Resultados</DropdownMenuItem>
            <DropdownMenuItem>Editar Pesquisa</DropdownMenuItem>
            <DropdownMenuItem className="text-red-600">
              Excluir
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];

export function DataTable() {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState({});

  const table = useReactTable({
    data: dadosPesquisas,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
  });

  return (
    <div className="w-full">
      <div className="flex items-center py-4">
        <Input
          placeholder="Filtrar titulos..."
          value={(table.getColumn("titulo")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("titulo")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="ml-auto">
              Colunas <ChevronDown />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {table
              .getAllColumns()
              .filter((column) => column.getCanHide())
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className="capitalize"
                    checked={column.getIsVisible()}
                    onCheckedChange={(value) =>
                      column.toggleVisibility(!!value)
                    }
                  >
                    {column.id}
                  </DropdownMenuCheckboxItem>
                );
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="overflow-hidden rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  Nenhum resultado.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-end space-x-2 py-4">
        <div className="text-muted-foreground flex-1 text-sm">
          {table.getFilteredSelectedRowModel().rows.length} of{" "}
          {table.getFilteredRowModel().rows.length} row(s) selected
        </div>
        <div className="space-x-2">
        <Button
            className="cursor-pointer"
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Anterior
          </Button>
          <Button
            className="cursor-pointer"
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Próximo
          </Button>
        </div>
      </div>
    </div>
  );
}
