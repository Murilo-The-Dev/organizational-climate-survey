"use client";

import * as React from "react";
import {
  Scatter,
  ScatterChart,
  ResponsiveContainer,
  XAxis,
  YAxis,
} from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { TrendingUp } from "lucide-react";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";

/**
 * Estrutura de dados para uma resposta individual da pesquisa.
 */
export interface RespostaIndividual {
  usuarioId: string;
  perguntaId: string;
  categoria: string;
  departamento?: string;
  pontuacao: number;
  respondidoEm: string; // Data/hora da resposta no formato ISO 8601
  pontuacaoMaxima: number; // Pontuação máxima possível para a pergunta (ex: 5 ou 10)
}

// TODO: Substituir com a busca de dados reais da API (ex: surveyResults).
// Este mock simula os dados que seriam recebidos da sua fonte de dados.
export const mockApiData: RespostaIndividual[] = [
  // Mês 1: Jan/24 - Escala 0-5
  { usuarioId: "U1", perguntaId: "Q1", categoria: "Liderança", pontuacao: 4, respondidoEm: "2024-01-15T10:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U2", perguntaId: "Q1", categoria: "Liderança", pontuacao: 5, respondidoEm: "2024-01-20T14:30:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U3", perguntaId: "Q2", categoria: "Bem-estar", pontuacao: 3, respondidoEm: "2024-01-10T09:15:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U11", perguntaId: "Q3", categoria: "Comunicação", pontuacao: 2, respondidoEm: "2024-01-25T11:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U12", perguntaId: "Q4", categoria: "Engajamento", pontuacao: 5, respondidoEm: "2024-01-28T16:20:00Z", pontuacaoMaxima: 5, departamento: "RH" },
  // Mês 2: Fev/24 - Mix de escalas 0-5 e 0-10 (o gráfico usará a maior: 10)
  { usuarioId: "U4", perguntaId: "Q5", categoria: "Comunicação", pontuacao: 8, respondidoEm: "2024-02-12T11:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U5", perguntaId: "Q6", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2024-02-05T16:20:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U6", perguntaId: "Q7", categoria: "Bem-estar", pontuacao: 4, respondidoEm: "2024-02-22T18:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U7", perguntaId: "Q8", categoria: "Liderança", pontuacao: 2, respondidoEm: "2024-02-25T10:45:00Z", pontuacaoMaxima: 5, departamento: "RH" },
  // Mês 3: Mar/24 - Escala 0-10
  { usuarioId: "U8", perguntaId: "Q6", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2024-03-01T12:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U9", perguntaId: "Q7", categoria: "Bem-estar", pontuacao: 8, respondidoEm: "2024-03-19T14:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U10", perguntaId: "Q8", categoria: "Liderança", pontuacao: 10, respondidoEm: "2024-03-15T09:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },

  { usuarioId: "U13", perguntaId: "Q9", categoria: "Liderança", pontuacao: 8, respondidoEm: "2024-03-20T09:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U14", perguntaId: "Q9", categoria: "Liderança", pontuacao: 7, respondidoEm: "2024-03-20T09:05:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U15", perguntaId: "Q10", categoria: "Bem-estar", pontuacao: 6, respondidoEm: "2024-03-20T09:10:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U16", perguntaId: "Q10", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2024-03-20T09:12:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U17", perguntaId: "Q11", categoria: "Comunicação", pontuacao: 5, respondidoEm: "2024-03-20T10:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U18", perguntaId: "Q11", categoria: "Comunicação", pontuacao: 7, respondidoEm: "2024-03-20T10:05:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U19", perguntaId: "Q12", categoria: "Engajamento", pontuacao: 10, respondidoEm: "2024-03-20T11:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U20", perguntaId: "Q12", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2024-03-20T11:02:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U21", perguntaId: "Q13", categoria: "Liderança", pontuacao: 6, respondidoEm: "2024-03-20T14:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U22", perguntaId: "Q14", categoria: "Bem-estar", pontuacao: 7, respondidoEm: "2024-03-20T15:30:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },

  { usuarioId: "U23", perguntaId: "Q15", categoria: "Liderança", pontuacao: 3, respondidoEm: "2024-04-02T10:00:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },
  { usuarioId: "U24", perguntaId: "Q15", categoria: "Liderança", pontuacao: 4, respondidoEm: "2024-04-03T11:00:00Z", pontuacaoMaxima: 5, departamento: "Marketing" },
  { usuarioId: "U25", perguntaId: "Q16", categoria: "Bem-estar", pontuacao: 5, respondidoEm: "2024-04-05T09:30:00Z", pontuacaoMaxima: 5, departamento: "Financeiro" },
  { usuarioId: "U26", perguntaId: "Q16", categoria: "Bem-estar", pontuacao: 2, respondidoEm: "2024-04-05T14:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U27", perguntaId: "Q17", categoria: "Comunicação", pontuacao: 4, respondidoEm: "2024-04-10T08:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U28", perguntaId: "Q17", categoria: "Comunicação", pontuacao: 3, respondidoEm: "2024-04-11T16:00:00Z", pontuacaoMaxima: 5, departamento: "RH" },
  { usuarioId: "U29", perguntaId: "Q18", categoria: "Engajamento", pontuacao: 5, respondidoEm: "2024-04-15T11:00:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },
  { usuarioId: "U30", perguntaId: "Q18", categoria: "Engajamento", pontuacao: 4, respondidoEm: "2024-04-18T13:20:00Z", pontuacaoMaxima: 5, departamento: "Marketing" },
  { usuarioId: "U31", perguntaId: "Q15", categoria: "Liderança", pontuacao: 5, respondidoEm: "2024-04-20T10:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U32", perguntaId: "Q16", categoria: "Bem-estar", pontuacao: 3, respondidoEm: "2024-04-22T15:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U33", perguntaId: "Q17", categoria: "Comunicação", pontuacao: 2, respondidoEm: "2024-04-25T09:00:00Z", pontuacaoMaxima: 5, departamento: "Financeiro" },
  { usuarioId: "U34", perguntaId: "Q18", categoria: "Engajamento", pontuacao: 4, respondidoEm: "2024-04-28T17:00:00Z", pontuacaoMaxima: 5, departamento: "RH" },
  { usuarioId: "U35", perguntaId: "Q15", categoria: "Liderança", pontuacao: 3, respondidoEm: "2024-04-29T12:00:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },

  { usuarioId: "U36", perguntaId: "Q19", categoria: "Liderança", pontuacao: 8, respondidoEm: "2024-05-01T09:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U37", perguntaId: "Q19", categoria: "Liderança", pontuacao: 9, respondidoEm: "2024-05-02T14:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U38", perguntaId: "Q20", categoria: "Bem-estar", pontuacao: 7, respondidoEm: "2024-05-05T11:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U39", perguntaId: "Q20", categoria: "Bem-estar", pontuacao: 6, respondidoEm: "2024-05-06T10:30:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U40", perguntaId: "Q21", categoria: "Comunicação", pontuacao: 10, respondidoEm: "2024-05-10T15:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U41", perguntaId: "Q21", categoria: "Comunicação", pontuacao: 8, respondidoEm: "2024-05-12T18:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U42", perguntaId: "Q22", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2024-05-15T09:45:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U43", perguntaId: "Q22", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2024-05-18T16:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U44", perguntaId: "Q19", categoria: "Liderança", pontuacao: 6, respondidoEm: "2024-05-20T11:30:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U45", perguntaId: "Q20", categoria: "Bem-estar", pontuacao: 8, respondidoEm: "2024-05-21T10:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U46", perguntaId: "Q21", categoria: "Comunicação", pontuacao: 9, respondidoEm: "2024-05-25T14:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U47", perguntaId: "Q22", categoria: "Engajamento", pontuacao: 10, respondidoEm: "2024-05-28T09:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U48", perguntaId: "Q19", categoria: "Liderança", pontuacao: 7, respondidoEm: "2024-05-30T13:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },

  { usuarioId: "U49", perguntaId: "Q23", categoria: "Liderança", pontuacao: 9, respondidoEm: "2024-06-03T09:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U50", perguntaId: "Q23", categoria: "Liderança", pontuacao: 8, respondidoEm: "2024-06-04T10:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U51", perguntaId: "Q24", categoria: "Bem-estar", pontuacao: 7, respondidoEm: "2024-06-05T11:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U52", perguntaId: "Q24", categoria: "Bem-estar", pontuacao: 6, respondidoEm: "2024-06-06T14:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U53", perguntaId: "Q25", categoria: "Comunicação", pontuacao: 10, respondidoEm: "2024-06-10T16:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U54", perguntaId: "Q25", categoria: "Comunicação", pontuacao: 9, respondidoEm: "2024-06-11T09:30:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U55", perguntaId: "Q26", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2024-06-15T13:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U56", perguntaId: "Q26", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2024-06-18T15:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U57", perguntaId: "Q23", categoria: "Liderança", pontuacao: 10, respondidoEm: "2024-06-20T10:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U58", perguntaId: "Q24", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2024-06-22T11:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U59", perguntaId: "Q25", categoria: "Comunicação", pontuacao: 8, respondidoEm: "2024-06-25T14:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U60", perguntaId: "Q26", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2024-06-28T16:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U61", perguntaId: "Q23", categoria: "Liderança", pontuacao: 9, respondidoEm: "2024-06-29T09:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },

  { usuarioId: "U62", perguntaId: "Q27", categoria: "Bem-estar", pontuacao: 5, respondidoEm: "2024-01-18T10:00:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },
  { usuarioId: "U63", perguntaId: "Q27", categoria: "Bem-estar", pontuacao: 4, respondidoEm: "2024-01-19T11:00:00Z", pontuacaoMaxima: 5, departamento: "Marketing" },
  { usuarioId: "U64", perguntaId: "Q28", categoria: "Liderança", pontuacao: 8, respondidoEm: "2024-02-15T09:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U65", perguntaId: "Q28", categoria: "Liderança", pontuacao: 7, respondidoEm: "2024-02-16T14:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U66", perguntaId: "Q29", categoria: "Comunicação", pontuacao: 9, respondidoEm: "2024-03-22T08:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U67", perguntaId: "Q29", categoria: "Comunicação", pontuacao: 6, respondidoEm: "2024-03-23T16:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U68", perguntaId: "Q30", categoria: "Engajamento", pontuacao: 4, respondidoEm: "2024-04-08T11:00:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },
  { usuarioId: "U69", perguntaId: "Q30", categoria: "Engajamento", pontuacao: 5, respondidoEm: "2024-04-09T13:20:00Z", pontuacaoMaxima: 5, departamento: "Marketing" },
  { usuarioId: "U70", perguntaId: "Q31", categoria: "Bem-estar", pontuacao: 8, respondidoEm: "2024-05-13T10:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U71", perguntaId: "Q31", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2024-05-14T15:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U72", perguntaId: "Q32", categoria: "Liderança", pontuacao: 7, respondidoEm: "2024-06-12T09:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U73", perguntaId: "Q32", categoria: "Liderança", pontuacao: 10, respondidoEm: "2024-06-13T17:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U74", perguntaId: "Q33", categoria: "Comunicação", pontuacao: 3, respondidoEm: "2024-01-22T12:00:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },
  { usuarioId: "U75", perguntaId: "Q33", categoria: "Comunicação", pontuacao: 4, respondidoEm: "2024-01-23T09:00:00Z", pontuacaoMaxima: 5, departamento: "Marketing" },
  { usuarioId: "U76", perguntaId: "Q34", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2024-02-18T14:30:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U77", perguntaId: "Q34", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2024-02-19T10:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U78", perguntaId: "Q35", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2024-03-25T11:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U79", perguntaId: "Q35", categoria: "Bem-estar", pontuacao: 10, respondidoEm: "2024-03-26T15:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U80", perguntaId: "Q36", categoria: "Liderança", pontuacao: 2, respondidoEm: "2024-04-12T09:30:00Z", pontuacaoMaxima: 5, departamento: "Vendas" },
  { usuarioId: "U81", perguntaId: "Q36", categoria: "Liderança", pontuacao: 3, respondidoEm: "2024-04-13T16:00:00Z", pontuacaoMaxima: 5, departamento: "Marketing" },
  { usuarioId: "U82", perguntaId: "Q37", categoria: "Comunicação", pontuacao: 6, respondidoEm: "2024-05-16T11:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U83", perguntaId: "Q37", categoria: "Comunicação", pontuacao: 7, respondidoEm: "2024-05-17T14:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U84", perguntaId: "Q38", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2024-06-19T10:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U85", perguntaId: "Q38", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2024-06-21T12:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U86", perguntaId: "Q39", categoria: "Bem-estar", pontuacao: 4, respondidoEm: "2024-01-05T09:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U87", perguntaId: "Q39", categoria: "Bem-estar", pontuacao: 5, respondidoEm: "2024-01-06T10:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U88", perguntaId: "Q40", categoria: "Liderança", pontuacao: 9, respondidoEm: "2024-02-08T11:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U89", perguntaId: "Q40", categoria: "Liderança", pontuacao: 8, respondidoEm: "2024-02-09T14:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U90", perguntaId: "Q41", categoria: "Comunicação", pontuacao: 7, respondidoEm: "2024-03-11T16:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U91", perguntaId: "Q41", categoria: "Comunicação", pontuacao: 6, respondidoEm: "2024-03-12T09:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U92", perguntaId: "Q42", categoria: "Engajamento", pontuacao: 3, respondidoEm: "2024-04-14T13:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U93", perguntaId: "Q42", categoria: "Engajamento", pontuacao: 4, respondidoEm: "2024-04-16T15:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U94", perguntaId: "Q43", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2024-05-19T10:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U95", perguntaId: "Q43", categoria: "Bem-estar", pontuacao: 10, respondidoEm: "2024-05-22T11:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U96", perguntaId: "Q44", categoria: "Liderança", pontuacao: 8, respondidoEm: "2024-06-23T14:30:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U97", perguntaId: "Q44", categoria: "Liderança", pontuacao: 7, respondidoEm: "2024-06-24T16:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U98", perguntaId: "Q45", categoria: "Comunicação", pontuacao: 5, respondidoEm: "2024-01-29T09:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U99", perguntaId: "Q45", categoria: "Comunicação", pontuacao: 4, respondidoEm: "2024-01-30T10:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U100", perguntaId: "Q46", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2024-02-28T11:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U101", perguntaId: "Q46", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2024-02-29T14:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U102", perguntaId: "Q47", categoria: "Bem-estar", pontuacao: 7, respondidoEm: "2024-03-05T16:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U103", perguntaId: "Q47", categoria: "Bem-estar", pontuacao: 6, respondidoEm: "2024-03-06T09:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U104", perguntaId: "Q48", categoria: "Liderança", pontuacao: 2, respondidoEm: "2024-04-01T13:00:00Z", pontuacaoMaxima: 5, departamento: "Engenharia" },
  { usuarioId: "U105", perguntaId: "Q48", categoria: "Liderança", pontuacao: 3, respondidoEm: "2024-04-04T15:00:00Z", pontuacaoMaxima: 5, departamento: "Produto" },
  { usuarioId: "U106", perguntaId: "Q49", categoria: "Comunicação", pontuacao: 8, respondidoEm: "2024-05-08T10:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U107", perguntaId: "Q49", categoria: "Comunicação", pontuacao: 9, respondidoEm: "2024-05-09T11:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U108", perguntaId: "Q50", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2024-06-08T14:30:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U109", perguntaId: "Q50", categoria: "Engajamento", pontuacao: 6, respondidoEm: "2024-06-09T16:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U110", perguntaId: "Q51", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2025-01-10T10:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U111", perguntaId: "Q51", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2025-01-12T11:30:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U112", perguntaId: "Q52", categoria: "Liderança", pontuacao: 7, respondidoEm: "2025-01-15T09:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U113", perguntaId: "Q52", categoria: "Liderança", pontuacao: 6, respondidoEm: "2025-01-18T14:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U114", perguntaId: "Q53", categoria: "Bem-estar", pontuacao: 5, respondidoEm: "2025-01-20T16:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U115", perguntaId: "Q53", categoria: "Bem-estar", pontuacao: 8, respondidoEm: "2025-01-22T08:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U116", perguntaId: "Q54", categoria: "Comunicação", pontuacao: 9, respondidoEm: "2025-01-25T13:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U117", perguntaId: "Q54", categoria: "Comunicação", pontuacao: 7, respondidoEm: "2025-01-28T10:15:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U118", perguntaId: "Q55", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2025-02-05T09:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U119", perguntaId: "Q55", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2025-02-07T11:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U120", perguntaId: "Q56", categoria: "Liderança", pontuacao: 9, respondidoEm: "2025-02-10T14:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U121", perguntaId: "Q56", categoria: "Liderança", pontuacao: 6, respondidoEm: "2025-02-12T16:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U122", perguntaId: "Q57", categoria: "Bem-estar", pontuacao: 8, respondidoEm: "2025-02-15T10:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U123", perguntaId: "Q57", categoria: "Bem-estar", pontuacao: 7, respondidoEm: "2025-02-18T11:45:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U124", perguntaId: "Q58", categoria: "Comunicação", pontuacao: 6, respondidoEm: "2025-02-20T09:30:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U125", perguntaId: "Q58", categoria: "Comunicação", pontuacao: 9, respondidoEm: "2025-02-22T15:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U126", perguntaId: "Q59", categoria: "Engajamento", pontuacao: 10, respondidoEm: "2025-02-25T12:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U127", perguntaId: "Q59", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2025-02-28T13:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U128", perguntaId: "Q60", categoria: "Liderança", pontuacao: 7, respondidoEm: "2025-03-03T10:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U129", perguntaId: "Q60", categoria: "Liderança", pontuacao: 8, respondidoEm: "2025-03-05T11:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U130", perguntaId: "Q61", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2025-03-08T14:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U131", perguntaId: "Q61", categoria: "Bem-estar", pontuacao: 6, respondidoEm: "2025-03-10T16:00:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U132", perguntaId: "Q62", categoria: "Comunicação", pontuacao: 8, respondidoEm: "2025-03-12T09:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U133", perguntaId: "Q62", categoria: "Comunicação", pontuacao: 7, respondidoEm: "2025-03-15T10:30:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U134", perguntaId: "Q63", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2025-03-18T13:00:00Z", pontuacaoMaxima: 10, departamento: "Engenharia" },
  { usuarioId: "U135", perguntaId: "Q63", categoria: "Engajamento", pontuacao: 10, respondidoEm: "2025-03-20T15:00:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U136", perguntaId: "Q64", categoria: "Liderança", pontuacao: 8, respondidoEm: "2025-03-22T11:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U137", perguntaId: "Q64", categoria: "Liderança", pontuacao: 7, respondidoEm: "2025-03-25T14:30:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U138", perguntaId: "Q65", categoria: "Bem-estar", pontuacao: 6, respondidoEm: "2025-03-28T09:00:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U139", perguntaId: "Q65", categoria: "Bem-estar", pontuacao: 9, respondidoEm: "2025-03-30T10:00:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U140", perguntaId: "Q66", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2025-03-25T09:00:00Z", pontuacaoMaxima: 10, departamento: "Tecnologia" },
  { usuarioId: "U141", perguntaId: "Q66", categoria: "Engajamento", pontuacao: 7, respondidoEm: "2025-03-25T09:05:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U142", perguntaId: "Q67", categoria: "Liderança", pontuacao: 9, respondidoEm: "2025-03-25T09:10:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
  { usuarioId: "U143", perguntaId: "Q67", categoria: "Liderança", pontuacao: 6, respondidoEm: "2025-03-25T09:15:00Z", pontuacaoMaxima: 10, departamento: "Tecnologia" },
  { usuarioId: "U144", perguntaId: "Q68", categoria: "Bem-estar", pontuacao: 5, respondidoEm: "2025-03-25T10:00:00Z", pontuacaoMaxima: 10, departamento: "RH" },
  { usuarioId: "U145", perguntaId: "Q68", categoria: "Bem-estar", pontuacao: 8, respondidoEm: "2025-03-25T10:05:00Z", pontuacaoMaxima: 10, departamento: "Vendas" },
  { usuarioId: "U146", perguntaId: "Q69", categoria: "Comunicação", pontuacao: 10, respondidoEm: "2025-03-25T11:20:00Z", pontuacaoMaxima: 10, departamento: "Produto" },
  { usuarioId: "U147", perguntaId: "Q69", categoria: "Comunicação", pontuacao: 7, respondidoEm: "2025-03-25T11:25:00Z", pontuacaoMaxima: 10, departamento: "Financeiro" },
  { usuarioId: "U148", perguntaId: "Q70", categoria: "Engajamento", pontuacao: 9, respondidoEm: "2025-03-25T14:00:00Z", pontuacaoMaxima: 10, departamento: "Tecnologia" },
  { usuarioId: "U149", perguntaId: "Q70", categoria: "Engajamento", pontuacao: 8, respondidoEm: "2025-03-25T14:05:00Z", pontuacaoMaxima: 10, departamento: "Marketing" },
];

const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="p-2 text-sm bg-background border rounded-lg shadow-lg">
          <p className="font-bold">
            {`Respondido em: ${new Date(data.respondidoEm).toLocaleDateString(
              "pt-BR",
              { day: "2-digit", month: "short", year: "numeric" }
            )}`}
          </p>
          <p style={{ color: payload[0].color }}>
            {`${data.categoria}: ${data.pontuacao.toFixed(0)}`}
          </p>
          {data.departamento && (
            <p className="text-muted-foreground">{`Depto: ${data.departamento}`}</p>
          )}
        </div>
      );
    }
    return null;
};

export function IndividualScoresScatterPlot() {
  // Em uma aplicação real, os dados viriam de uma prop ou de um hook de busca de dados.
  const [unifyPoints, setUnifyPoints] = React.useState(false);
  const data = mockApiData;

  if (!data || data.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp /> Dispersão de Pontuações Individuais
          </CardTitle>
          <CardDescription>
            Visualização da dispersão das respostas individuais ao longo do tempo, por categoria.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p>Não há dados de respostas para exibir no momento.</p>
        </CardContent>
      </Card>
    );
  }

  // Detecta a pontuação máxima para definir o domínio do eixo Y dinamicamente.
  const yAxisMax = Math.max(...data.map(item => item.pontuacaoMaxima), 5);
  const yAxisTicks = Array.from({ length: yAxisMax + 1 }, (_, i) => i);

  const { periodos, dataWithJitter } = React.useMemo(() => {
    const getMonthYear = (dateString: string) => {
      const date = new Date(dateString);
      return date.toLocaleDateString('pt-BR', { month: 'short', year: '2-digit' }).replace('.', '').replace(' de ', '/');
    };

    const meses: { [key: string]: number } = { 'jan': 1, 'fev': 2, 'mar': 3, 'abr': 4, 'mai': 5, 'jun': 6, 'jul': 7, 'ago': 8, 'set': 9, 'out': 10, 'nov': 11, 'dez': 12 };

    const currentPeriodos = [...new Set(data.map(item => getMonthYear(item.respondidoEm)))].sort((a, b) => {
          const [m1Str, y1] = a.toLowerCase().split('/');
          const [m2Str, y2] = b.toLowerCase().split('/');
          const dateA = new Date(parseInt(`20${y1}`), meses[m1Str] - 1);
          const dateB = new Date(parseInt(`20${y2}`), meses[m2Str] - 1);
          return dateA.getTime() - dateB.getTime();
        });

    const periodoMap = new Map(currentPeriodos.map((p, i) => [p, i]));

    const processedData = data.map(item => {
      const periodo = getMonthYear(item.respondidoEm);
      const jitter = unifyPoints ? 0 : (Math.random() - 0.5) * 0.7;
      return {
        ...item,
        periodo,
        x: (periodoMap.get(periodo) ?? 0) + jitter,
        y: item.pontuacao,
      };
    });

    return { periodos: currentPeriodos, dataWithJitter: processedData };
  }, [data, unifyPoints]);

  // Agrupa os dados por categoria para renderizar um <Scatter> para cada uma.
  const groupedData = dataWithJitter.reduce((acc, item) => {
    const key = item.categoria;
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(item);
    return acc;
  }, {} as Record<string, typeof dataWithJitter>);

  // Configuração de cores para cada categoria.
  const categoryColors = ["#155DFC", "#16A34A", "#DB2777", "#9333EA"];
  const chartConfig = Object.keys(groupedData).reduce((acc, key, index) => {
    acc[key] = {
      label: key,
      color: categoryColors[index % categoryColors.length],
    };
    return acc;
  }, {} as ChartConfig);

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <div>
            <CardTitle className="flex items-center gap-2">
                <TrendingUp /> Dispersão de Pontuações Individuais
            </CardTitle>
            <CardDescription>
                Visualização da dispersão das respostas individuais ao longo do tempo, por categoria.
            </CardDescription>
          </div>
          <div className="flex items-center space-x-2">
            <Switch id="unify-points" onCheckedChange={setUnifyPoints} />
            <Label htmlFor="unify-points">Unificar Pontos</Label>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="min-h-[450px] max-h-[600px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <ScatterChart margin={{ top: 20, right: 40, bottom: 40, left: 10 }}>
              <XAxis
                type="number"
                dataKey="x"
                name="Período"
                domain={[-0.5, periodos.length - 0.5]}
                ticks={Array.from({ length: periodos.length }, (_, i) => i)}
                tickFormatter={(tick) => periodos[tick] || ""}
                angle={periodos.length > 15 ? -45 : -35}
                textAnchor="end"
                height={periodos.length > 15 ? 80 : 60}
                interval={0}
                tickLine={true}
                axisLine={false}
              />
              <YAxis
                type="number"
                dataKey="y"
                name="Pontuação"
                domain={[0, yAxisMax]}
                ticks={yAxisTicks}
                allowDecimals={false}
                tickLine={false}
                axisLine={false}
              />
              <ChartTooltip
                cursor={{ strokeDasharray: '3 3' }}
                content={<CustomTooltip />}
              />
              <ChartLegend content={<ChartLegendContent />} />
              {Object.entries(groupedData).map(([categoria, points]) => (
                <Scatter
                  key={categoria}
                  name={categoria}
                  data={points}
                  fill={`var(--color-${categoria})`}
                  shape="circle"
                  opacity={0.7}
                />
              ))}
            </ScatterChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}