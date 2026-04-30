"use client";

import { useEffect, useState } from "react";
import { TrendingUp, Loader2 } from "lucide-react";
import { Pie, PieChart } from "recharts";
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
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { dashboardService } from "@/lib/services/dashboardService";
import { useAuth } from "@/context/AuthContext";

export function ChartPieLabel() {
  const { user } = useAuth();
  const [chartData, setChartData] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (user) {
      carregarDados();
    }
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
      const metricaComDistribuicao = metricas.find((m: any) => m.distribuicao);

      if (metricaComDistribuicao && metricaComDistribuicao.distribuicao) {
        const cores = [
          "var(--color-blue-100)",
          "var(--color-blue-300)",
          "var(--color-blue-500)",
          "var(--color-blue-700)",
          "#0f172a",
        ];

        const formatado = Object.entries(
          metricaComDistribuicao.distribuicao,
        ).map(([chave, valor], index) => ({
          nome: chave,
          valor: Number(valor),
          fill: cores[index % cores.length],
        }));

        setChartData(formatado);
      }
    } catch (error) {
      console.error("Erro ao carregar dados do gráfico:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const chartConfig = {
    valor: { label: "Respostas" },
  } satisfies ChartConfig;

  return (
    <Card className="flex flex-col h-full">
      <CardHeader className="items-center pb-0">
        <CardTitle>Distribuição de Respostas</CardTitle>
        <CardDescription>Métrica com maior volume de dados</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0 flex items-center justify-center min-h-[250px]">
        {isLoading ? (
          <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
        ) : chartData.length > 0 ? (
          <ChartContainer
            config={chartConfig}
            className="[&_.recharts-pie-label-text]:fill-foreground mx-auto aspect-square max-h-[250px] pb-0 w-full"
          >
            <PieChart>
              <ChartTooltip content={<ChartTooltipContent hideLabel />} />
              <Pie data={chartData} dataKey="valor" nameKey="nome" label />
            </PieChart>
          </ChartContainer>
        ) : (
          <p className="text-sm text-muted-foreground">
            Nenhum dado disponível.
          </p>
        )}
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm mt-4">
        <div className="flex items-center gap-2 leading-none font-medium">
          Atualizado dinamicamente <TrendingUp className="h-4 w-4" />
        </div>
        <div className="text-muted-foreground leading-none text-center">
          Lendo informações direto da base de dados.
        </div>
      </CardFooter>
    </Card>
  );
}
