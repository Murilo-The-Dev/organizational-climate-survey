"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, XAxis } from "recharts";
import { DateRange } from "react-day-picker";
import { Loader2 } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
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

const chartConfig = {
  media: {
    label: "Média de Notas",
    color: "var(--color-blue-600)",
  },
  respostas: {
    label: "Volume de Respostas",
    color: "var(--color-blue-400)",
  },
} satisfies ChartConfig;

export function ChartBarInteractive({ dateRange }: { dateRange?: DateRange }) {
  const { user } = useAuth();
  const [chartData, setChartData] = React.useState<any[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [activeChart, setActiveChart] =
    React.useState<keyof typeof chartConfig>("media");

  React.useEffect(() => {
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

      const formatado = metricas.map((m: any, index: number) => ({
        perguntaCurta: m.texto_pergunta
          ? m.texto_pergunta.substring(0, 15) + "..."
          : `Q${index + 1}`,
        perguntaCompleta: m.texto_pergunta || "Métrica desconhecida",
        media: m.media ? Number(m.media.toFixed(1)) : 0,
        respostas: m.total_respostas || 0,
      }));

      setChartData(
        formatado.length > 0
          ? formatado
          : [
              {
                perguntaCurta: "Sem dados",
                perguntaCompleta: "Nenhum dado disponível",
                media: 0,
                respostas: 0,
              },
            ],
      );
    } catch (error) {
      console.error("Erro ao carregar ChartBarInteractive:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Calcula os totais dinâmicos para os botões do topo do gráfico
  const totalMedia = React.useMemo(() => {
    if (chartData.length === 0 || chartData[0].perguntaCurta === "Sem dados")
      return "0.0";
    const sum = chartData.reduce((acc, curr) => acc + curr.media, 0);
    return (sum / chartData.length).toFixed(1);
  }, [chartData]);

  const totalRespostas = React.useMemo(() => {
    if (chartData.length === 0 || chartData[0].perguntaCurta === "Sem dados")
      return "0";
    return chartData.reduce((acc, curr) => acc + curr.respostas, 0);
  }, [chartData]);

  return (
    <Card className="py-0 w-full">
      <CardHeader className="flex flex-col items-stretch border-b !p-0 sm:flex-row">
        <div className="flex flex-1 flex-col justify-center gap-1 px-6 pt-4 pb-3 sm:!py-0">
          <CardTitle>Desempenho Detalhado por Métrica</CardTitle>
          <CardDescription>
            Alterne entre as abas para ver a média de notas ou o volume de
            participação
          </CardDescription>
        </div>
        <div className="flex">
          <button
            data-active={activeChart === "media"}
            className="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6 cursor-pointer hover:bg-muted/30 transition-colors"
            onClick={() => setActiveChart("media")}
          >
            <span className="text-muted-foreground text-xs">
              Média Geral (0-5)
            </span>
            <span className="text-lg leading-none font-bold sm:text-3xl">
              {totalMedia}
            </span>
          </button>
          <button
            data-active={activeChart === "respostas"}
            className="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6 cursor-pointer hover:bg-muted/30 transition-colors"
            onClick={() => setActiveChart("respostas")}
          >
            <span className="text-muted-foreground text-xs">
              Total de Avaliações
            </span>
            <span className="text-lg leading-none font-bold sm:text-3xl">
              {totalRespostas}
            </span>
          </button>
        </div>
      </CardHeader>
      <CardContent className="px-2 sm:p-6 min-h-[250px] flex items-center justify-center">
        {isLoading ? (
          <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
        ) : (
          <ChartContainer
            config={chartConfig}
            className="aspect-auto h-[250px] w-full"
          >
            <BarChart
              accessibilityLayer
              data={chartData}
              margin={{ left: 12, right: 12 }}
            >
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="perguntaCurta"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    className="w-[200px]"
                    nameKey={activeChart}
                    labelFormatter={(value, payload) => {
                      // Mostra a pergunta inteira no tooltip quando passar o mouse!
                      return payload?.[0]?.payload?.perguntaCompleta || value;
                    }}
                  />
                }
              />
              <Bar
                dataKey={activeChart}
                fill={`var(--color-${activeChart})`}
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
