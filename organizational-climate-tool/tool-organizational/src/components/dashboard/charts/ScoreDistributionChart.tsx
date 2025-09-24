// src/components/dashboard/charts/ScoreDistributionChart.tsx
"use client";
import { Bar, BarChart, CartesianGrid, XAxis } from "recharts";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

const chartData = [
  { score: 1, responses: 10 },
  { score: 2, responses: 25 },
  { score: 3, responses: 40 },
  { score: 4, responses: 120 },
  { score: 5, responses: 98 },
];
const chartConfig = {
  responses: { label: "Respostas", color: "var(--color-blue-600)" },
} satisfies ChartConfig;

export const ScoreDistributionChart = () => (
  <Card className="h-full grid grid-rows-[auto_1fr]">
    <CardHeader>
      <CardTitle>Distribuição de Notas</CardTitle>
      <CardDescription>
        Contagem de respostas para cada nota (1-5)
      </CardDescription>
    </CardHeader>
    <CardContent>
      <ChartContainer config={chartConfig} className="h-full w-full">
        <BarChart accessibilityLayer data={chartData}>
          <CartesianGrid vertical={false} />
          <XAxis
            dataKey="score"
            tickLine={false}
            tickMargin={10}
            axisLine={false}
          />
          <ChartTooltip
            cursor={false}
            content={<ChartTooltipContent hideLabel />}
          />
          <Bar dataKey="responses" fill="var(--color-blue-600)" radius={8} />
        </BarChart>
      </ChartContainer>
    </CardContent>
  </Card>
);
