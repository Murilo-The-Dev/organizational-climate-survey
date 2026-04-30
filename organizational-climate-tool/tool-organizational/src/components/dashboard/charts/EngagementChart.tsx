"use client";

import { useEffect, useState } from "react";
import { TrendingUp, Loader2 } from "lucide-react";
import { Bar, BarChart, CartesianGrid, XAxis } from "recharts";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { dashboardService } from "@/lib/services/dashboardService";
import { useAuth } from "@/context/AuthContext";

const chartConfig = {
  positivo: {
    label: "Positivo (Notas 4-5)",
    color: "var(--color-blue-400)",
  },
  negativo: {
    label: "Atenção (Notas 1-3)",
    color: "var(--color-blue-600)",
  },
} satisfies ChartConfig;

export function ChartBarStacked() {
  const { user } = useAuth();
  const [chartData, setChartData] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (user) carregarDados();
  }, [user]);

  const carregarDados = async () => {
    try {
      setIsLoading(true);
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;

      const func =
        dashboardService.getData || (dashboardService as any).getMetrics;
      const data = await func(empresaId);

      const metricas = data?.metricas_por_pergunta || [];

      const formatado = metricas.slice(0, 5).map((m: any) => {
        let pos = 0;
        let neg = 0;

        if (m.distribuicao) {
          Object.entries(m.distribuicao).forEach(([nota, qtd]) => {
            const numNota = Number(nota);
            const quantidade = Number(qtd);

            if (numNota >= 4) pos += quantidade;
            else if (numNota > 0) neg += quantidade;
            else {
              if (
                ["promotor", "bom", "excelente", "sim"].some((palavra) =>
                  nota.toLowerCase().includes(palavra),
                )
              ) {
                pos += quantidade;
              } else {
                neg += quantidade;
              }
            }
          });
        } else if (m.media) {
          pos = Math.round((m.media / 5) * (m.total_respostas || 10));
          neg = (m.total_respostas || 10) - pos;
        }

        return {
          categoria: m.texto_pergunta
            ? m.texto_pergunta.substring(0, 12) + "..."
            : "Métrica",
          positivo: pos,
          negativo: neg,
        };
      });

      setChartData(
        formatado.length > 0
          ? formatado
          : [{ categoria: "Sem dados", positivo: 0, negativo: 0 }],
      );
    } catch (error) {
      console.error("Erro ao carregar ChartBarStacked:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="flex flex-col h-full">
      <CardHeader>
        <CardTitle>Engajamento por Pergunta</CardTitle>
        <CardDescription>Sentimento geral (Top 5 perguntas)</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 min-h-[250px] flex items-center justify-center">
        {isLoading ? (
          <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
        ) : (
          <ChartContainer
            config={chartConfig}
            className="w-full aspect-auto h-[250px]"
          >
            <BarChart accessibilityLayer data={chartData}>
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="categoria"
                tickLine={false}
                tickMargin={10}
                axisLine={false}
              />
              <ChartTooltip content={<ChartTooltipContent hideLabel />} />
              <ChartLegend content={<ChartLegendContent />} />
              <Bar
                dataKey="positivo"
                stackId="a"
                fill="var(--color-positivo)"
                radius={[0, 0, 4, 4]}
              />
              <Bar
                dataKey="negativo"
                stackId="a"
                fill="var(--color-negativo)"
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm mt-auto">
        <div className="flex gap-2 leading-none font-medium">
          Monitoramento contínuo <TrendingUp className="h-4 w-4" />
        </div>
        <div className="text-muted-foreground leading-none">
          Baseado nas respostas validadas do sistema.
        </div>
      </CardFooter>
    </Card>
  );
}
