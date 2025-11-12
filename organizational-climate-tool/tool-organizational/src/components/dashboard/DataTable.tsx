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
import { ExportButton } from "@/components/ui/export-button";

const dadosPesquisas: Pesquisa[] = [
  { id: "PESQ-001", 
    title: "Engajamento Q1 2025", 
    status: "concluido", participantes: 152, 
    dataCriacao: "2025-03-28" 
  },
  { id: "PESQ-002", 
    title: "Feedback de Liderança H1", 
    status: "concluido", 
    participantes: 140, 
    dataCriacao: "2025-06-15" 
  },
  { id: "PESQ-003", 
    title: "Pesquisa de Satisfação Anual 2024", 
    status: "concluido", 
    participantes: 180, 
    dataCriacao: "2024-12-20"
   },
  { id: "PESQ-004", 
    title: "Clima Organizacional H2", 
    status: "em_andamento", 
    participantes: 125, 
    dataCriacao: "2025-09-01" 
  },
  { id: "PESQ-005", 
    title: "Onboarding Novos Contratados", 
    status: "em_andamento", 
    participantes: 25, 
    dataCriacao: "2025-09-10" 
  },
  { id: "PESQ-006", 
    title: "Avaliação de Benefícios", 
    status: "rascunho", 
    participantes: 0, 
    dataCriacao: "2025-09-18" 
  },
  { id: "PESQ-007", 
    title: "Engajamento Q2 2025", 
    status: "rascunho", 
    participantes: 0, 
    dataCriacao: "2025-09-15" 
  },
  { id: "PESQ-008", 
    title: "Segurança Psicológica", 
    status: "concluido", 
    participantes: 165, 
    dataCriacao: "2025-01-30" 
  },
  { id: "PESQ-009", 
    title: "Comunicação Interna", 
    status: "em_andamento", 
    participantes: 95, 
    dataCriacao: "2025-08-22" 
  },
  { id: "PESQ-010", 
    title: "Planejamento Estratégico 2026", 
    status: "rascunho",
     participantes: 0, 
     dataCriacao: "2025-09-19" 
    },
  { id: "PESQ-011", 
    title: "Ferramentas de Trabalho", 
    status: "concluido", 
    participantes: 170, 
    dataCriacao: "2025-05-10" 
  },
  { id: "PESQ-012", 
    title: "e-NPS Semestral", 
    status: "em_andamento", 
    participantes: 110, 
    dataCriacao: "2025-09-05" 
  },
  { id: "PESQ-003", 
    title: "Engajamento Q1 2025", 
    status: "concluido", participantes: 152, 
    dataCriacao: "2025-03-28" 
  },
  { id: "PESQ-004", 
    title: "Feedback de Liderança H1", 
    status: "concluido", 
    participantes: 140, 
    dataCriacao: "2025-06-15" 
  },
  { id: "PESQ-005", 
    title: "Pesquisa de Satisfação Anual 2024", 
    status: "concluido", 
    participantes: 180, 
    dataCriacao: "2024-12-20"
   },
  { id: "PESQ-006", 
    title: "Clima Organizacional H2", 
    status: "em_andamento", 
    participantes: 125, 
    dataCriacao: "2025-09-01" 
  },
  { id: "PESQ-007", 
    title: "Onboarding Novos Contratados", 
    status: "em_andamento", 
    participantes: 25, 
    dataCriacao: "2025-09-10" 
  },
  { id: "PESQ-008", 
    title: "Avaliação de Benefícios", 
    status: "rascunho", 
    participantes: 0, 
    dataCriacao: "2025-09-18" 
  },
  { id: "PESQ-009", 
    title: "Engajamento Q2 2025", 
    status: "concluido", 
    participantes: 0, 
    dataCriacao: "2025-09-15" 
  },
  { id: "PESQ-010", 
    title: "Segurança Psicológica", 
    status: "concluido", 
    participantes: 165, 
    dataCriacao: "2025-01-30" 
  },
];

export type Pesquisa = {
  id: string;
  title: string;
  status: "concluido" | "em_andamento" | "rascunho";
  participantes: number;
  dataCriacao: string;
};

interface DataTableProps<TData extends object> {
  columns: ColumnDef<TData>[];
  data: TData[];
}

export function DataTable<TData extends object>({ columns, data }: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState({});

  const table = useReactTable({
    data,
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
            placeholder="Filtrar títulos..."
            value={table.getColumn("title") ? (table.getColumn("title")!.getFilterValue() as string) : ""}
            onChange={(event) => table.getColumn("title")?.setFilterValue(event.target.value)}
            className="max-w-sm"
        />
        <div className="ml-auto flex items-center gap-2">
          <ExportButton data={table.getFilteredRowModel().rows.map(row => row.original)} filename="dados_tabela" />
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

