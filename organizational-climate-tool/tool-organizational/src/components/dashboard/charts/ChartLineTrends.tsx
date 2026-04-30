"use client";

import * as React from "react";
import {
  CartesianGrid,
  Line,
  LineChart as RechartsLineChart,
  XAxis,
  YAxis,
} from "recharts";
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
  engajamento: { label: "Média (0-5)", color: "var(--color-blue-600)" },
  satisfacao: { label: "Total Respostas", color: "var(--color-blue-400)" },
} satisfies ChartConfig;

export function ChartLineTrends({ dateRange }: { dateRange?: DateRange }) {
  const { user } = useAuth();
  const [chartData, setChartData] = React.useState<any[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);

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
      const formatado = metricas.map((m: any, i: number) => ({
        name: m.texto_pergunta
          ? m.texto_pergunta.substring(0, 10) + "..."
          : `Q${i + 1}`,
        engajamento: m.media ? Number(m.media.toFixed(1)) : 0,
        satisfacao: m.total_respostas || 0,
      }));

      setChartData(formatado);
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="flex flex-col h-full">
      <CardHeader>
        <CardTitle>Tendências por Métrica</CardTitle>
        <CardDescription>
          Média e volume ao longo do questionário
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-1 min-h-[250px] flex items-center justify-center">
        {isLoading ? (
          <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
        ) : (
          <ChartContainer config={chartConfig} className="w-full h-[250px]">
            <RechartsLineChart
              data={chartData}
              margin={{ left: 12, right: 12 }}
            >
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="name"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
              />
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent hideLabel />}
              />
              <Line
                dataKey="engajamento"
                type="monotone"
                stroke="var(--color-blue-600)"
                strokeWidth={2}
                dot={false}
              />
              <Line
                dataKey="satisfacao"
                type="monotone"
                stroke="var(--color-blue-400)"
                strokeWidth={2}
                dot={false}
              />
            </RechartsLineChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
