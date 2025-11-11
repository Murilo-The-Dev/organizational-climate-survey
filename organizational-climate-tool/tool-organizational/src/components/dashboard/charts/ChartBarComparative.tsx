import * as React from "react";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis, Legend } from "recharts";
import { DateRange } from "react-day-picker";
import { isWithinInterval, parseISO } from "date-fns";

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

const allChartData = [
  { category: "Liderança", q1: 60, q2: 75, q3: 70 },
  { category: "Comunicação", q1: 70, q2: 65, q3: 80 },
  { category: "Colaboração", q1: 80, q2: 85, q3: 75 },
  { category: "Reconhecimento", q1: 50, q2: 60, q3: 65 },
  { category: "Desenvolvimento", q1: 75, q2: 80, q3: 85 },
];

const chartConfig = {
  q1: {
    label: "Q1",
    color: "hsl(var(--chart-1))",
  },
  q2: {
    label: "Q2",
    color: "hsl(var(--chart-2))",
  },
  q3: {
    label: "Q3",
    color: "hsl(var(--chart-3))",
  },
} satisfies ChartConfig;

interface ChartBarComparativeProps {
  dateRange?: DateRange;
}

export function ChartBarComparative({ dateRange }: ChartBarComparativeProps) {
  // Para este gráfico comparativo, o dateRange pode ser usado para selecionar quais trimestres/períodos comparar
  // Por simplicidade, vamos manter os dados mockados fixos por enquanto, mas a prop está disponível.
  console.log("Date range for ChartBarComparative:", dateRange);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Comparativo de Desempenho por Categoria</CardTitle>
        <CardDescription>Comparação de métricas chave ao longo de diferentes trimestres.</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart
            accessibilityLayer
            data={allChartData}
            margin={{
              left: 0,
              right: 5,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="category"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
            />
            <YAxis tickLine={false} axisLine={false} tickMargin={8} />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent />}
            />
            <Legend />
            <Bar dataKey="q1" fill="#2B7FFF" />
            <Bar dataKey="q2" fill="#5790e5" />
            <Bar dataKey="q3" fill="#0A4DB2" />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

