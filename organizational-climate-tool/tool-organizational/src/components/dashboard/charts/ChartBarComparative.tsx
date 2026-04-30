"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";
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
  media: { label: "Média (0-5)", color: "var(--color-blue-600)" },
} satisfies ChartConfig;

export function ChartBarComparative({ dateRange }: { dateRange?: DateRange }) {
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
        category: m.texto_pergunta
          ? m.texto_pergunta.substring(0, 15) + "..."
          : `Q${i + 1}`,
        media: m.media ? Number(m.media.toFixed(1)) : 0,
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
        <CardTitle>Comparativo Geral</CardTitle>
        <CardDescription>Comparação das médias entre métricas.</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 min-h-[250px] flex items-center justify-center">
        {isLoading ? (
          <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
        ) : (
          <ChartContainer config={chartConfig} className="w-full h-[250px]">
            <BarChart data={chartData} margin={{ left: 0, right: 5 }}>
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="category"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
              />
              <YAxis tickLine={false} axisLine={false} tickMargin={8} />
              <ChartTooltip cursor={false} content={<ChartTooltipContent />} />
              <Bar
                dataKey="media"
                fill="var(--color-blue-600)"
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
