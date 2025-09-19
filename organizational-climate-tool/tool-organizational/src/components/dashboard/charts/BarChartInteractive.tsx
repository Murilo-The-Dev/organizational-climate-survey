"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, XAxis } from "recharts";

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

export const description = "um gráfico de barras interativo";

const chartData = [
  { date: "2024-04-01", engajamento: 82, satisfacao: 78 },
  { date: "2024-04-02", engajamento: 97, satisfacao: 180 },
  { date: "2024-04-03", engajamento: 167, satisfacao: 120 },
  { date: "2024-04-04", engajamento: 242, satisfacao: 260 },
  { date: "2024-04-05", engajamento: 373, satisfacao: 290 },
  { date: "2024-04-06", engajamento: 301, satisfacao: 340 },
  { date: "2024-04-07", engajamento: 245, satisfacao: 180 },
  { date: "2024-04-08", engajamento: 409, satisfacao: 320 },
  { date: "2024-04-09", engajamento: 59, satisfacao: 110 },
  { date: "2024-04-10", engajamento: 261, satisfacao: 190 },
  { date: "2024-04-11", engajamento: 327, satisfacao: 350 },
  { date: "2024-04-12", engajamento: 292, satisfacao: 210 },
  { date: "2024-04-13", engajamento: 342, satisfacao: 380 },
  { date: "2024-04-14", engajamento: 137, satisfacao: 220 },
  { date: "2024-04-15", engajamento: 120, satisfacao: 170 },
  { date: "2024-04-16", engajamento: 138, satisfacao: 190 },
  { date: "2024-04-17", engajamento: 446, satisfacao: 360 },
  { date: "2024-04-18", engajamento: 364, satisfacao: 410 },
  { date: "2024-04-19", engajamento: 243, satisfacao: 180 },
  { date: "2024-04-20", engajamento: 89, satisfacao: 150 },
  { date: "2024-04-21", engajamento: 137, satisfacao: 200 },
  { date: "2024-04-22", engajamento: 224, satisfacao: 170 },
  { date: "2024-04-23", engajamento: 138, satisfacao: 230 },
  { date: "2024-04-24", engajamento: 387, satisfacao: 290 },
  { date: "2024-04-25", engajamento: 215, satisfacao: 250 },
  { date: "2024-04-26", engajamento: 75, satisfacao: 130 },
  { date: "2024-04-27", engajamento: 383, satisfacao: 420 },
  { date: "2024-04-28", engajamento: 122, satisfacao: 180 },
  { date: "2024-04-29", engajamento: 315, satisfacao: 240 },
  { date: "2024-04-30", engajamento: 454, satisfacao: 380 },
  { date: "2024-05-01", engajamento: 165, satisfacao: 220 },
  { date: "2024-05-02", engajamento: 293, satisfacao: 310 },
  { date: "2024-05-03", engajamento: 247, satisfacao: 190 },
  { date: "2024-05-04", engajamento: 385, satisfacao: 420 },
  { date: "2024-05-05", engajamento: 481, satisfacao: 390 },
  { date: "2024-05-06", engajamento: 498, satisfacao: 520 },
  { date: "2024-05-07", engajamento: 388, satisfacao: 300 },
  { date: "2024-05-08", engajamento: 149, satisfacao: 210 },
  { date: "2024-05-09", engajamento: 227, satisfacao: 180 },
  { date: "2024-05-10", engajamento: 293, satisfacao: 330 },
  { date: "2024-05-11", engajamento: 335, satisfacao: 270 },
  { date: "2024-05-12", engajamento: 197, satisfacao: 240 },
  { date: "2024-05-13", engajamento: 197, satisfacao: 160 },
  { date: "2024-05-14", engajamento: 448, satisfacao: 490 },
  { date: "2024-05-15", engajamento: 473, satisfacao: 380 },
  { date: "2024-05-16", engajamento: 338, satisfacao: 400 },
  { date: "2024-05-17", engajamento: 499, satisfacao: 420 },
  { date: "2024-05-18", engajamento: 315, satisfacao: 350 },
  { date: "2024-05-19", engajamento: 235, satisfacao: 180 },
  { date: "2024-05-20", engajamento: 177, satisfacao: 230 },
  { date: "2024-05-21", engajamento: 82, satisfacao: 140 },
  { date: "2024-05-22", engajamento: 81, satisfacao: 120 },
  { date: "2024-05-23", engajamento: 252, satisfacao: 290 },
  { date: "2024-05-24", engajamento: 294, satisfacao: 220 },
  { date: "2024-05-25", engajamento: 201, satisfacao: 250 },
  { date: "2024-05-26", engajamento: 213, satisfacao: 170 },
  { date: "2024-05-27", engajamento: 420, satisfacao: 460 },
  { date: "2024-05-28", engajamento: 233, satisfacao: 190 },
  { date: "2024-05-29", engajamento: 78, satisfacao: 130 },
  { date: "2024-05-30", engajamento: 340, satisfacao: 280 },
  { date: "2024-05-31", engajamento: 178, satisfacao: 230 },
  { date: "2024-06-01", engajamento: 178, satisfacao: 200 },
  { date: "2024-06-02", engajamento: 470, satisfacao: 410 },
  { date: "2024-06-03", engajamento: 103, satisfacao: 160 },
  { date: "2024-06-04", engajamento: 439, satisfacao: 380 },
  { date: "2024-06-05", engajamento: 88, satisfacao: 140 },
  { date: "2024-06-06", engajamento: 294, satisfacao: 250 },
  { date: "2024-06-07", engajamento: 323, satisfacao: 370 },
  { date: "2024-06-08", engajamento: 385, satisfacao: 320 },
  { date: "2024-06-09", engajamento: 438, satisfacao: 480 },
  { date: "2024-06-10", engajamento: 155, satisfacao: 200 },
  { date: "2024-06-11", engajamento: 92, satisfacao: 150 },
  { date: "2024-06-12", engajamento: 492, satisfacao: 420 },
  { date: "2024-06-13", engajamento: 81, satisfacao: 130 },
  { date: "2024-06-14", engajamento: 426, satisfacao: 380 },
  { date: "2024-06-15", engajamento: 307, satisfacao: 350 },
  { date: "2024-06-16", engajamento: 371, satisfacao: 310 },
  { date: "2024-06-17", engajamento: 475, satisfacao: 520 },
  { date: "2024-06-18", engajamento: 107, satisfacao: 170 },
  { date: "2024-06-19", engajamento: 341, satisfacao: 290 },
  { date: "2024-06-20", engajamento: 408, satisfacao: 450 },
  { date: "2024-06-21", engajamento: 169, satisfacao: 210 },
  { date: "2024-06-22", engajamento: 317, satisfacao: 270 },
  { date: "2024-06-23", engajamento: 480, satisfacao: 530 },
  { date: "2024-06-24", engajamento: 132, satisfacao: 180 },
  { date: "2024-06-25", engajamento: 141, satisfacao: 190 },
  { date: "2024-06-26", engajamento: 434, satisfacao: 380 },
  { date: "2024-06-27", engajamento: 448, satisfacao: 490 },
  { date: "2024-06-28", engajamento: 149, satisfacao: 200 },
  { date: "2024-06-29", engajamento: 103, satisfacao: 160 },
  { date: "2024-06-30", engajamento: 446, satisfacao: 400 },
];

const chartConfig = {
  views: {
    label: "Visualizações de Página",
  },
  engajamento: {
    label: "Engajamento",
    color: "var(--color-blue-400)",
  },
  satisfacao: {
    label: "Satisfação",
    color: "var(--color-blue-600)",
  },
} satisfies ChartConfig;

export function ChartBarInteractive() {
  const [activeChart, setActiveChart] =
    React.useState<keyof typeof chartConfig>("engajamento");

  const total = React.useMemo(
    () => ({
      engajamento: chartData.reduce((acc, curr) => acc + curr.engajamento, 0),
      satisfacao: chartData.reduce((acc, curr) => acc + curr.satisfacao, 0),
    }),
    []
  );

  return (
    <Card className="py-0 w-full">
      <CardHeader className="flex flex-col items-stretch border-b !p-0 sm:flex-row">
        <div className="flex flex-1 flex-col justify-center gap-1 px-6 pt-4 pb-3 sm:!py-0">
          <CardTitle>Gráfico de Barras - Interativo</CardTitle>
          <CardDescription>
            Mostrando o total de visitantes para os últimos 3 meses
          </CardDescription>
        </div>
        <div className="flex">
          {["engajamento", "satisfacao"].map((key) => {
            const chart = key as keyof typeof chartConfig;
            return (
              <button
                key={chart}
                data-active={activeChart === chart}
                className="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6"
                onClick={() => setActiveChart(chart)}
              >
                <span className="text-muted-foreground text-xs">
                  {chartConfig[chart].label}
                </span>
                <span className="text-lg leading-none font-bold sm:text-3xl">
                  {total[key as keyof typeof total].toLocaleString()}
                </span>
              </button>
            );
          })}
        </div>
      </CardHeader>
      <CardContent className="px-2 sm:p-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          <BarChart
            accessibilityLayer
            data={chartData}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={(value) => {
                const date = new Date(value);
                return date.toLocaleDateString("en-US", {
                  month: "short",
                  day: "numeric",
                });
              }}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  className="w-[150px]"
                  nameKey="views"
                  labelFormatter={(value) => {
                    return new Date(value).toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                      year: "numeric",
                    });
                  }}
                />
              }
            />
            <Bar dataKey={activeChart} fill={`var(--color-${activeChart})`} />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
