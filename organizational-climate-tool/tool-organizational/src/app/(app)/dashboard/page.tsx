"use client"

import StatCardGrids from "@/components/dashboard/StatCardGrids";
import { ChartBarStacked } from "@/components/dashboard/charts/EngagementChart";
import { ChartRadialShape } from "@/components/dashboard/charts/RadialChart";
import { ChartPieLabel } from "@/components/dashboard/charts/PieChart";
import { ChartBarInteractive } from "@/components/dashboard/charts/BarChartInteractive";
import { ChartLineTrends } from "@/components/dashboard/charts/ChartLineTrends";
import { ChartBarComparative } from "@/components/dashboard/charts/ChartBarComparative";
import { DataTable, Pesquisa } from "@/components/dashboard/DataTable";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { useState } from "react";
import { DateRange } from "react-day-picker";
import { ColumnDef } from "@tanstack/react-table";
import { StatusBadge } from "@/components/ui/status-badge";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

// Dados de exemplo (você pode substituir por dados reais da API)
const dadosPesquisas: Pesquisa[] = [
  { id: "PESQ-001", title: "Engajamento Q1 2025", status: "concluido", participantes: 152, dataCriacao: "2025-03-28" },
  { id: "PESQ-002", title: "Feedback de Liderança H1", status: "concluido", participantes: 140, dataCriacao: "2025-06-15" },
  { id: "PESQ-003", title: "Pesquisa de Satisfação Anual 2024", status: "concluido", participantes: 180, dataCriacao: "2024-12-20" },
  { id: "PESQ-004", title: "Clima Organizacional H2", status: "em_andamento", participantes: 125, dataCriacao: "2025-09-01" },
  { id: "PESQ-005", title: "Onboarding Novos Contratados", status: "em_andamento", participantes: 25, dataCriacao: "2025-09-10" },
  { id: "PESQ-006", title: "Avaliação de Benefícios", status: "rascunho", participantes: 0, dataCriacao: "2025-09-18" },
];

// Definição das colunas
const columns: ColumnDef<Pesquisa>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: ({ row }) => <div className="font-medium">{row.getValue("id")}</div>,
  },
  {
    accessorKey: "title",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Título
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue("title")}</div>,
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status") as "concluido" | "em_andamento" | "rascunho";
      return <StatusBadge status={status} />;
    },
  },
  {
    accessorKey: "participantes",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Participantes
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const participantes = row.getValue("participantes") as number;
      return <div className="text-center font-medium">{participantes}</div>;
    },
  },
  {
    accessorKey: "dataCriacao",
    header: "Data de Criação",
    cell: ({ row }) => {
      const data = new Date(row.getValue("dataCriacao"));
      return <div>{data.toLocaleDateString("pt-BR")}</div>;
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const pesquisa = row.original;

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
            <DropdownMenuItem onClick={() => navigator.clipboard.writeText(pesquisa.id)}>
              Copiar ID
            </DropdownMenuItem>
            <DropdownMenuItem>Ver detalhes</DropdownMenuItem>
            <DropdownMenuItem>Editar</DropdownMenuItem>
            <DropdownMenuItem className="text-red-600">Excluir</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];

const DashboardPage = () => {
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined);

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Dashboard
        </h1>
        <DateRangePicker date={dateRange} onSelect={setDateRange} />
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Visão geral da sua organização.
      </p>

      <StatCardGrids dateRange={dateRange} />

      <div className="mt-6 lg:col-span-3">
        <ChartBarInteractive dateRange={dateRange} />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-6">
        <div>
          <ChartBarStacked dateRange={dateRange} />
        </div>
        <div>
          <ChartPieLabel dateRange={dateRange} />
        </div>
        <div>
          <ChartRadialShape dateRange={dateRange} />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
        <ChartLineTrends dateRange={dateRange} />
        <ChartBarComparative dateRange={dateRange} />
      </div>

      <div className="mt-6">
        <div className="bg-background rounded-lg border p-4 h-full">
          <DataTable columns={columns} data={dadosPesquisas} />
        </div>
      </div>
    </section>
  );
};

export default DashboardPage;