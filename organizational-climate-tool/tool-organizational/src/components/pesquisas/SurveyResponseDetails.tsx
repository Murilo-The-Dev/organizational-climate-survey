"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ListChecks } from "lucide-react";

export function SurveyResponseDetails({ surveyId }: { surveyId: string }) {
  // Alexandre: Como não temos um endpoint de "listagem de respostas brutas" por causa do anonimato,
  // esta tela agora serve como placeholder para auditoria de volume por ID.
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <ListChecks className="h-5 w-5 text-green-500" /> Detalhes das
          Respostas
        </CardTitle>
        <CardDescription>
          Visualização de integridade dos dados brutos.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col items-center justify-center py-10 text-center border-2 border-dashed rounded-xl">
          <p className="font-medium">Dados Anonimizados</p>
          <p className="text-sm text-muted-foreground max-w-xs">
            Para garantir o sigilo, as respostas individuais não são exibidas.
            Consulte o Relatório Analítico para ver os agrupamentos.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
