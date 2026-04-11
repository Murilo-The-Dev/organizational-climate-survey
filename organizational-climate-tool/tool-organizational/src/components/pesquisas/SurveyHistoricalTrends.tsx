// src/components/pesquisas/SurveyHistoricalTrends.tsx
"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@radix-ui/react-dropdown-menu";
import { TrendingUp } from "lucide-react";

interface SurveyHistoricalTrendsProps {
  surveyId: string;
  dateRange: any; // O tipo real seria DateRange do react-day-picker
}

export function SurveyHistoricalTrends({ surveyId, dateRange }: SurveyHistoricalTrendsProps) {
  // Lógica para buscar dados históricos baseados em surveyId e dateRange
  const startDate = dateRange?.from ? dateRange.from.toLocaleDateString('pt-BR') : 'Início';
  const endDate = dateRange?.to ? dateRange.to.toLocaleDateString('pt-BR') : 'Fim';

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <TrendingUp className="h-5 w-5 text-blue-500" />
          Tendência Histórica da Pesquisa
        </CardTitle>
        <CardDescription>
          Análise da evolução dos resultados ao longo do tempo.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground">
          **Pesquisa ID:** {surveyId} | **Intervalo Selecionado:** {startDate} até {endDate}
        </p>
        <Separator />
        
        {/* Placeholder para o Gráfico de Tendência */}
        <div className="h-96 w-full flex items-center justify-center bg-gray-50 border rounded-lg">
          <p className="text-lg text-gray-500 italic">
            Placeholder: Gráfico de Linha com a Tendência Histórica (Ex: Média de Engajamento por Mês)
          </p>
        </div>

        {/* Placeholder para Insights */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4">
            <div className="p-4 border rounded-lg bg-green-50/50">
                <h4 className="font-semibold text-green-700">Ponto Alto</h4>
                <p className="text-sm">O score de "Satisfação com a Liderança" aumentou 15% no último trimestre.</p>
            </div>
            <div className="p-4 border rounded-lg bg-red-50/50">
                <h4 className="font-semibold text-red-700">Ponto de Atenção</h4>
                <p className="text-sm">O score de "Equilíbrio entre Vida Pessoal e Profissional" caiu 8% no último mês.</p>
            </div>
            <div className="p-4 border rounded-lg bg-blue-50/50">
                <h4 className="font-semibold text-blue-700">Comparativo</h4>
                <p className="text-sm">O resultado atual está 5% acima da média anual.</p>
            </div>
        </div>
      </CardContent>
    </Card>
  );
}
