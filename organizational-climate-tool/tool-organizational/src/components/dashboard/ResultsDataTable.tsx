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
  DropdownMenuSeparator
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
  departmentScores?: {
    [department: string]: {
        averageScore: number;
        responseCount: number;
    }
  }
};

// 2. Criando dados de exemplo
export const allMockResults: { [surveyId: string]: SurveyResult[] } = {
    "PESQ-001": [ // Engajamento Q1 2025
        { 
            id: "q1-1", question: "Meu gestor direto me dá feedback construtivo regularmente.", category: "Liderança", averageScore: 4.2, responseCount: 150,
            departmentScores: {
                "Tecnologia": { averageScore: 4.5, responseCount: 40 }, "Recursos Humanos": { averageScore: 4.1, responseCount: 15 }, "Marketing": { averageScore: 4.0, responseCount: 25 }, "Vendas": { averageScore: 4.2, responseCount: 50 }, "Financeiro": { averageScore: 3.9, responseCount: 20 },
            }
        },
        { 
            id: "q1-2", question: "Sinto que tenho um bom equilíbrio entre vida profissional e pessoal.", category: "Bem-estar", averageScore: 3.1, responseCount: 148,
            departmentScores: {
                "Tecnologia": { averageScore: 2.9, responseCount: 38 }, "Recursos Humanos": { averageScore: 3.8, responseCount: 15 }, "Marketing": { averageScore: 3.5, responseCount: 25 }, "Vendas": { averageScore: 3.0, responseCount: 50 }, "Financeiro": { averageScore: 3.2, responseCount: 20 },
            }
        },
        { 
            id: "q1-3", question: "A comunicação entre os departamentos é eficiente.", category: "Comunicação", averageScore: 2.8, responseCount: 145,
            departmentScores: {
                "Tecnologia": { averageScore: 2.5, responseCount: 38 }, "Recursos Humanos": { averageScore: 3.5, responseCount: 14 }, "Marketing": { averageScore: 3.0, responseCount: 23 }, "Vendas": { averageScore: 2.6, responseCount: 50 }, "Financeiro": { averageScore: 2.9, responseCount: 20 },
            }
        },
        { 
            id: "q1-4", question: "Sinto-me reconhecido pelo meu trabalho.", category: "Bem-estar", averageScore: 3.5, responseCount: 151,
            departmentScores: {
                "Tecnologia": { averageScore: 3.3, responseCount: 40 }, "Recursos Humanos": { averageScore: 4.0, responseCount: 15 }, "Marketing": { averageScore: 3.6, responseCount: 25 }, "Vendas": { averageScore: 3.4, responseCount: 51 }, "Financeiro": { averageScore: 3.5, responseCount: 20 },
            }
        },
        { 
            id: "q1-5", question: "As metas da equipe são claras e alcançáveis.", category: "Liderança", averageScore: 4.0, responseCount: 152,
            departmentScores: {
                "Tecnologia": { averageScore: 4.2, responseCount: 40 }, "Recursos Humanos": { averageScore: 4.1, responseCount: 15 }, "Marketing": { averageScore: 3.8, responseCount: 25 }, "Vendas": { averageScore: 4.0, responseCount: 52 }, "Financeiro": { averageScore: 3.9, responseCount: 20 },
            }
        },
    ],
    "PESQ-002": [ // Feedback de Liderança H1
        { 
            id: "q2-1", question: "A liderança da empresa é transparente em suas decisões.", category: "Liderança", averageScore: 3.5, responseCount: 130,
            departmentScores: { "Tecnologia": { averageScore: 3.6, responseCount: 35 }, "Recursos Humanos": { averageScore: 3.8, responseCount: 12 }, "Vendas": { averageScore: 3.4, responseCount: 45 }, "Financeiro": { averageScore: 3.5, responseCount: 18 } }
        },
        { 
            id: "q2-2", question: "Sinto-me à vontade para dar feedback aos meus líderes.", category: "Comunicação", averageScore: 4.0, responseCount: 135,
            departmentScores: { "Tecnologia": { averageScore: 4.2, responseCount: 38 }, "Recursos Humanos": { averageScore: 4.5, responseCount: 13 }, "Vendas": { averageScore: 3.8, responseCount: 46 }, "Financeiro": { averageScore: 4.1, responseCount: 18 } }
        },
        { 
            id: "q2-3", question: "Meu líder direto demonstra preocupação com meu bem-estar.", category: "Liderança", averageScore: 4.1, responseCount: 138,
            departmentScores: { "Tecnologia": { averageScore: 4.3, responseCount: 39 }, "Recursos Humanos": { averageScore: 4.5, responseCount: 13 }, "Vendas": { averageScore: 3.9, responseCount: 47 }, "Financeiro": { averageScore: 4.2, responseCount: 19 } }
        },
        { 
            id: "q2-4", question: "A liderança inspira confiança e motivação na equipe.", category: "Liderança", averageScore: 3.8, responseCount: 132,
            departmentScores: { "Tecnologia": { averageScore: 3.9, responseCount: 36 }, "Recursos Humanos": { averageScore: 4.2, responseCount: 12 }, "Vendas": { averageScore: 3.7, responseCount: 45 }, "Financeiro": { averageScore: 3.8, responseCount: 19 } }
        },
    ],
    "PESQ-003": [ // Pesquisa de Satisfação Anual 2024
        { 
            id: "q3-1", question: "Estou satisfeito com meu pacote de benefícios.", category: "Bem-estar", averageScore: 3.8, responseCount: 170,
            departmentScores: { "Tecnologia": { averageScore: 3.9, responseCount: 45 }, "Recursos Humanos": { averageScore: 4.0, responseCount: 20 }, "Vendas": { averageScore: 3.7, responseCount: 60 }, "Financeiro": { averageScore: 3.8, responseCount: 25 } }
        },
        { 
            id: "q3-2", question: "O ambiente de trabalho é positivo e colaborativo.", category: "Bem-estar", averageScore: 4.1, responseCount: 175,
            departmentScores: { "Tecnologia": { averageScore: 4.2, responseCount: 48 }, "Recursos Humanos": { averageScore: 4.5, responseCount: 20 }, "Vendas": { averageScore: 4.0, responseCount: 62 }, "Financeiro": { averageScore: 4.1, responseCount: 25 } }
        },
        { 
            id: "q3-3", question: "A empresa investe no meu desenvolvimento profissional.", category: "Liderança", averageScore: 3.6, responseCount: 168,
            departmentScores: { "Tecnologia": { averageScore: 3.8, responseCount: 44 }, "Recursos Humanos": { averageScore: 4.0, responseCount: 19 }, "Vendas": { averageScore: 3.4, responseCount: 60 }, "Financeiro": { averageScore: 3.7, responseCount: 25 } }
        },
        { 
            id: "q3-4", question: "A carga de trabalho é gerenciável.", category: "Bem-estar", averageScore: 2.9, responseCount: 178,
            departmentScores: { "Tecnologia": { averageScore: 2.7, responseCount: 49 }, "Recursos Humanos": { averageScore: 3.5, responseCount: 20 }, "Vendas": { averageScore: 2.8, responseCount: 64 }, "Financeiro": { averageScore: 3.1, responseCount: 25 } }
        },
        { 
            id: "q3-5", question: "Sinto que meu trabalho tem um propósito claro.", category: "Bem-estar", averageScore: 4.4, responseCount: 180,
            departmentScores: { "Tecnologia": { averageScore: 4.5, responseCount: 50 }, "Recursos Humanos": { averageScore: 4.6, responseCount: 20 }, "Vendas": { averageScore: 4.3, responseCount: 65 }, "Financeiro": { averageScore: 4.4, responseCount: 25 } }
        },
    ],
    "PESQ-004": [ // Clima Organizacional H2 (em_andamento)
        { 
            id: "q4-1", question: "Os processos internos da empresa são eficientes.", category: "Comunicação", averageScore: 3.2, responseCount: 120,
            departmentScores: { "Tecnologia": { averageScore: 3.0, responseCount: 30 }, "Vendas": { averageScore: 3.3, responseCount: 50 } }
        },
        { 
            id: "q4-2", question: "Há um sentimento de justiça e equidade na empresa.", category: "Bem-estar", averageScore: 3.4, responseCount: 115,
            departmentScores: { "Tecnologia": { averageScore: 3.5, responseCount: 28 }, "Vendas": { averageScore: 3.3, responseCount: 48 } }
        },
        { 
            id: "q4-3", question: "A empresa promove um ambiente de trabalho inclusivo.", category: "Bem-estar", averageScore: 4.0, responseCount: 122,
            departmentScores: { "Tecnologia": { averageScore: 4.1, responseCount: 31 }, "Vendas": { averageScore: 3.9, responseCount: 52 } }
        },
    ],
    "PESQ-005": [ // Onboarding Novos Contratados (em_andamento)
        { 
            id: "q5-1", question: "O processo de integração foi claro e bem estruturado.", category: "Comunicação", averageScore: 4.3, responseCount: 24,
            departmentScores: { "Recursos Humanos": { averageScore: 4.5, responseCount: 24 } }
        },
        { 
            id: "q5-2", question: "Recebi as ferramentas e acessos necessários para começar a trabalhar.", category: "Bem-estar", averageScore: 4.5, responseCount: 25,
            departmentScores: { "Recursos Humanos": { averageScore: 4.5, responseCount: 25 } }
        },
    ],
    "PESQ-008": [ // Segurança Psicológica
        { 
            id: "q8-1", question: "Sinto-me seguro para expressar minhas opiniões, mesmo que discordem da maioria.", category: "Bem-estar", averageScore: 4.5, responseCount: 160,
            departmentScores: { "Tecnologia": { averageScore: 4.6, responseCount: 42 }, "Recursos Humanos": { averageScore: 4.8, responseCount: 18 }, "Marketing": { averageScore: 4.4, responseCount: 30 } }
        },
        { 
            id: "q8-2", question: "Erros são tratados como oportunidades de aprendizado, não como falhas.", category: "Liderança", averageScore: 4.3, responseCount: 162,
            departmentScores: { "Tecnologia": { averageScore: 4.2, responseCount: 43 }, "Recursos Humanos": { averageScore: 4.6, responseCount: 18 }, "Marketing": { averageScore: 4.1, responseCount: 31 } }
        },
        { 
            id: "q8-3", question: "Sinto que posso ser eu mesmo no trabalho, sem medo de julgamentos.", category: "Bem-estar", averageScore: 4.6, responseCount: 165,
            departmentScores: { "Tecnologia": { averageScore: 4.7, responseCount: 44 }, "Recursos Humanos": { averageScore: 4.8, responseCount: 19 }, "Marketing": { averageScore: 4.5, responseCount: 32 } }
        },
    ],
    "PESQ-009": [ // Comunicação Interna (em_andamento)
        { 
            id: "q9-1", question: "As informações importantes sobre a empresa são comunicadas de forma clara.", category: "Comunicação", averageScore: 3.3, responseCount: 90,
            departmentScores: { "Marketing": { averageScore: 3.5, responseCount: 20 }, "Vendas": { averageScore: 3.2, responseCount: 40 } }
        },
        { 
            id: "q9-2", question: "Os canais de comunicação interna (e-mail, chat) são eficazes.", category: "Comunicação", averageScore: 3.7, responseCount: 92,
            departmentScores: { "Marketing": { averageScore: 3.8, responseCount: 21 }, "Vendas": { averageScore: 3.6, responseCount: 41 } }
        },
    ],
    "PESQ-011": [ // Ferramentas de Trabalho
        { 
            id: "q11-1", question: "As ferramentas e softwares que utilizo são adequados para meu trabalho.", category: "Bem-estar", averageScore: 3.7, responseCount: 165,
            departmentScores: { "Tecnologia": { averageScore: 4.2, responseCount: 45 }, "Financeiro": { averageScore: 3.5, responseCount: 25 } }
        },
        { 
            id: "q11-2", question: "Tenho o hardware necessário (computador, monitor) para realizar meu trabalho eficientemente.", category: "Bem-estar", averageScore: 4.1, responseCount: 168,
            departmentScores: { "Tecnologia": { averageScore: 4.5, responseCount: 46 }, "Financeiro": { averageScore: 3.9, responseCount: 26 } }
        },
        { 
            id: "q11-3", question: "O suporte técnico é ágil e resolve meus problemas.", category: "Comunicação", averageScore: 3.9, responseCount: 150,
            departmentScores: { "Tecnologia": { averageScore: 4.0, responseCount: 40 }, "Financeiro": { averageScore: 3.8, responseCount: 22 } }
        },
    ],
    "PESQ-012": [ // e-NPS Semestral (em_andamento)
        { 
            id: "q12-1", question: "Em uma escala de 0 a 10, o quão provável você é de recomendar esta empresa como um bom lugar para trabalhar?", category: "Bem-estar", averageScore: 4.0, responseCount: 105,
            departmentScores: { "Tecnologia": { averageScore: 4.2, responseCount: 30 }, "Vendas": { averageScore: 3.9, responseCount: 50 } }
        },
        { 
            id: "q12-2", question: "Sinto que a empresa se preocupa com a minha saúde mental.", category: "Desmotivação", averageScore: 2.5, responseCount: 108,
            departmentScores: { "Tecnologia": { averageScore: 2.8, responseCount: 32 }, "Vendas": { averageScore: 2.3, responseCount: 51 } }
        },
    ],
    // Pesquisas em rascunho não terão resultados
    "PESQ-006": [],
    "PESQ-007": [],
    "PESQ-010": [],
};

export const mockResults: SurveyResult[] = [];

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
          <Progress
            value={percentage}
            className="w-[60%] [&>div]:bg-[#155dfc]"
          />
          <span className="font-medium">{score.toFixed(1)}</span>
        </div>
      );
    },
  },
];

// 4. O componente da tabela que recebe os dados
export function ResultsDataTable({
  data: tableData,
  isLoading,
}: {
  data: SurveyResult[];
  isLoading?: boolean;
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
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  Carregando resultados...
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
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
