"use client";

import { useEffect, useState } from "react";
import { Loader2 } from "lucide-react";
import {
  Label,
  PolarGrid,
  PolarRadiusAxis,
  RadialBar,
  RadialBarChart,
} from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ChartConfig, ChartContainer } from "@/components/ui/chart";
import { dashboardService } from "@/lib/services/dashboardService";
import { useAuth } from "@/context/AuthContext";

const chartConfig = {
  valor: {
    label: "Participação",
  },
  participacao: {
    label: "Participação",
    color: "var(--color-blue-600)",
  },
} satisfies ChartConfig;

export function ChartRadialShape() {
  const { user } = useAuth();
  const [taxa, setTaxa] = useState(0);
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
        dashboardService.getMetrics || (dashboardService as any).getData;
      const data = await func(empresaId);

      setTaxa(data?.taxa_participacao || 0);
    } catch (error) {
      console.error("Erro ao carregar dados do gráfico radial:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const chartData = [
    { metrica: "participacao", valor: taxa, fill: "var(--color-blue-600)" },
  ];

  return (
    <Card className="flex flex-col h-full">
      <CardHeader className="items-center pb-0">
        <CardTitle>Gráfico De Forma - Participação</CardTitle>
        <CardDescription>Taxa de adesão atual</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0 flex items-center justify-center min-h-[250px]">
        {isLoading ? (
          <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
        ) : (
          <ChartContainer
            config={chartConfig}
            className="mx-auto aspect-square max-h-[250px] w-full"
          >
            <RadialBarChart
              data={chartData}
              endAngle={100}
              innerRadius={80}
              outerRadius={140}
            >
              <PolarGrid
                gridType="circle"
                radialLines={false}
                stroke="none"
                className="first:fill-muted last:fill-background"
                polarRadius={[86, 74]}
              />
              <RadialBar dataKey="valor" background />
              <PolarRadiusAxis tick={false} tickLine={false} axisLine={false}>
                <Label
                  content={({ viewBox }) => {
                    if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                      return (
                        <text
                          x={viewBox.cx}
                          y={viewBox.cy}
                          textAnchor="middle"
                          dominantBaseline="middle"
                        >
                          <tspan
                            x={viewBox.cx}
                            y={viewBox.cy}
                            className="fill-foreground text-4xl font-bold"
                          >
                            {chartData[0].valor}%
                          </tspan>
                          <tspan
                            x={viewBox.cx}
                            y={(viewBox.cy || 0) + 24}
                            className="fill-muted-foreground"
                          >
                            Participação
                          </tspan>
                        </text>
                      );
                    }
                  }}
                />
              </PolarRadiusAxis>
            </RadialBarChart>
          </ChartContainer>
        )}
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm mt-4 text-center">
        <div className="text-muted-foreground leading-none">
          Lendo métricas oficiais da empresa
        </div>
      </CardFooter>
    </Card>
  );
}
