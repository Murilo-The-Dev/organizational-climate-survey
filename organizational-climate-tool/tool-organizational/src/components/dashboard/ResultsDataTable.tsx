// src/components/dashboard/ResultsDataTable.tsx

"use client";

import * as React from "react";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { MoreHorizontal } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

// 1. Definindo a estrutura dos nossos dados de resultado
export type SurveyResult = {
  id: string;
  question: string;
  category: "Liderança" | "Bem-estar" | "Comunicação" | "Desmotivação";
  averageScore: number; // Uma nota de 0 a 5
  responseCount: number;
};

// 2. Criando dados de exemplo
export const data: SurveyResult[] = [
  {
    id: "q1",
    question: "Meu gestor direto me dá feedback construtivo regularmente.",
    category: "Liderança",
    averageScore: 4.2,
    responseCount: 85,
  },
  {
    id: "q2",
    question:
      "Sinto que tenho um bom equilíbrio entre vida profissional e pessoal.",
    category: "Bem-estar",
    averageScore: 3.1,
    responseCount: 88,
  },
  {
    id: "q3",
    question: "A comunicação entre os departamentos é eficiente.",
    category: "Comunicação",
    averageScore: 2.8,
    responseCount: 82,
  },
  {
    id: "q4",
    question: "Tenho clareza sobre as oportunidades de crescimento na empresa.",
    category: "Liderança",
    averageScore: 3.9,
    responseCount: 86,
  },
  {
    id: "q5",
    question:
      "Não tenho clareza sobre as oportunidades de crescimento na empresa.",
    category: "Desmotivação",
    averageScore: 1.9,
    responseCount: 10,
  },
];

// 3. Definindo as colunas da nossa tabela
export const columns: ColumnDef<SurveyResult>[] = [
  {
    accessorKey: "question",
    header: "Pergunta",
    cell: ({ row }) => (
      <div>
        <div className="font-medium">{row.original.question}</div>
        <Badge variant="outline" className="mt-1">
          {row.original.category}
        </Badge>
      </div>
    ),
  },
  {
    accessorKey: "responseCount",
    header: () => <div className="text-center">Nº de Respostas</div>,
    cell: ({ row }) => (
      <div className="text-center">{row.getValue("responseCount")}</div>
    ),
  },
  {
    accessorKey: "averageScore",
    header: "Média da Pontuação",
    cell: ({ row }) => {
      const score = row.getValue("averageScore") as number;
      const percentage = (score / 5) * 100;
      return (
        <div className="flex items-center gap-2">
          <Progress value={percentage} className="w-[60%]" />
          <span className="font-medium">{score.toFixed(1)}</span>
        </div>
      );
    },
  },
  {
    id: "actions",
    cell: ({ row }) => (
      <div className="text-right">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Abrir menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Ações</DropdownMenuLabel>
            <DropdownMenuItem>Ver detalhes das respostas</DropdownMenuItem>
            <DropdownMenuItem>Ver tendência histórica</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    ),
  },
];

// 4. O componente da tabela que recebe os dados
export function ResultsDataTable({
  data: tableData,
}: {
  data: SurveyResult[];
}) {
  const table = useReactTable({
    data: tableData,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Resultados</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
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
                  Nenhum resultado para os filtros selecionados.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
