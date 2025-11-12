"use client";

import { TrendingUp } from "lucide-react";
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

export const description = "A stacked bar chart with a legend";

const chartData = [
  { month: "January", positivo: 186, negativo: 80 },
  { month: "February", positivo: 305, negativo: 200 },
  { month: "March", positivo: 237, negativo: 120 },
  { month: "April", positivo: 73, negativo: 190 },
  { month: "May", positivo: 209, negativo: 130 },
  { month: "June", positivo: 214, negativo: 140 },
];

const chartConfig = {
  positivo: {
    label: "Positivo",
    color: "var(--color-blue-400)",
  },
  negativo: {
    label: "Negativo",
    color: "var(--color-blue-600)",
  },
} satisfies ChartConfig;

export function ChartBarStacked() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Gráfico de Barras - Positivo + Negativo</CardTitle>
        <CardDescription>Janeiro - Junho 2024</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart accessibilityLayer data={chartData}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="month"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              tickFormatter={(value) => value.slice(0, 3)}
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
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 leading-none font-medium">
          Aumento de 5.2% este mês <TrendingUp className="h-4 w-4" />
        </div>
        <div className="text-muted-foreground leading-none">
          Mostrando o total de Respostas para os últimos 6 meses
        </div>
      </CardFooter>
    </Card>
  );
}
