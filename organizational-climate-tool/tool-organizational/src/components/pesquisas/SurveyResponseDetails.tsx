// src/components/pesquisas/SurveyResponseDetails.tsx
"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ListChecks } from "lucide-react";

interface SurveyResponseDetailsProps {
  surveyId: string;
}

// Dados mockados para a tabela
const mockResponses = [
    { id: 1, setor: 'TI', data: '2025-04-01', q1: '5', q2: '4', q3: '5' },
    { id: 2, setor: 'RH', data: '2025-04-02', q1: '3', q2: '5', q3: '4' },
    { id: 3, setor: 'Vendas', data: '2025-04-03', q1: '4', q2: '3', q3: '3' },
];

export function SurveyResponseDetails({ surveyId }: SurveyResponseDetailsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <ListChecks className="h-5 w-5 text-green-500" />
          Detalhes das Respostas
        </CardTitle>
        <CardDescription>
          Visualização das respostas individuais e dados brutos da pesquisa.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground mb-4">
          **Pesquisa ID:** {surveyId} | Exibindo as últimas {mockResponses.length} respostas.
        </p>
        
        {/* Tabela de Detalhes das Respostas */}
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[100px]">ID Resposta</TableHead>
                <TableHead>Setor</TableHead>
                <TableHead>Data</TableHead>
                <TableHead>Q1 (Liderança)</TableHead>
                <TableHead>Q2 (Ambiente)</TableHead>
                <TableHead>Q3 (Remuneração)</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockResponses.map((response) => (
                <TableRow key={response.id}>
                  <TableCell className="font-medium">{response.id}</TableCell>
                  <TableCell>{response.setor}</TableCell>
                  <TableCell>{response.data}</TableCell>
                  <TableCell>{response.q1}</TableCell>
                  <TableCell>{response.q2}</TableCell>
                  <TableCell>{response.q3}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
        
        <div className="mt-4 text-center text-sm text-gray-500">
            <p>Use o filtro de data na aba "Tendência Histórica" para refinar os resultados.</p>
        </div>
      </CardContent>
    </Card>
  );
}
