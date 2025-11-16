// src/components/dashboard/charts/ScoreDistributionChart.tsx
"use client";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";
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
  { score: 1, responses: 60 },
  { score: 2, responses: 25 },
  { score: 3, responses: 40 },
  { score: 4, responses: 120 },
  { score: 5, responses: 98 },
];
const chartConfig = {
  responses: { label: "Respostas", color: "var(--color-blue-600)" },
} satisfies ChartConfig;

export const ScoreDistributionChart = () => (
  <Card className="h-full grid grid-rows-[auto_1fr] overflow-hidden">
    <CardHeader>
      <CardTitle>Distribuição de Notas</CardTitle>
      <CardDescription>
        Contagem de respostas para cada nota (1-5)
      </CardDescription>
    </CardHeader>
    <CardContent className="h-full w-full p-4 overflow-x-auto overflow-y-hidden self-center">
      <ChartContainer config={chartConfig} className="h-full min-w-[200px] w-full max-h-96">
        <BarChart 
          accessibilityLayer 
          data={chartData}
          barSize={25}
          margin={{ top: 5, right: 20, bottom: 20, left: 0 }}
        >
          <CartesianGrid vertical={false} />
          <XAxis
            dataKey="score"
            tickLine={false}
            tickMargin={10}
            axisLine={false}
          />
          <YAxis tickLine={false} axisLine={false} tickMargin={8} />
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