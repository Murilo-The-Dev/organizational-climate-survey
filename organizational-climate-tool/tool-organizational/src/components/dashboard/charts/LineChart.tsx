"use client";

import { TrendingUp } from "lucide-react";
import { CartesianGrid, Line, LineChart, XAxis } from "recharts";

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

export const description = "Gráfico de curva de respostas";

const chartData = [
  { month: "January", responses: 186, mobile: 80 },
  { month: "February", responses: 305, mobile: 200 },
  { month: "March", responses: 237, mobile: 120 },
  { month: "April", responses: 73, mobile: 190 },
  { month: "May", responses: 209, mobile: 130 },
  { month: "June", responses: 214, mobile: 140 },
];

const chartConfig = {
  responses: {
    label: "Respostas",
    color: "var(--color-blue-600)",
  },
} satisfies ChartConfig;

export function ChartLineDefault() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Curva de Respostas</CardTitle>
        <CardDescription>Janeiro - Junho 2024</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <LineChart
            accessibilityLayer
            data={chartData}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="month"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => value.slice(0, 3)}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Line
              dataKey="responses"
              type="natural"
              stroke="var(--color-blue-600)"
              strokeWidth={2}
              dot={{
                fill: "var(--color-blue-600)",
              }}
              activeDot={{
                r: 6,
              }}
            />
          </LineChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 leading-none font-medium">
          Comparativo de respostas dos últimos 6 meses.
          <TrendingUp className="h-4 w-4" />
        </div>
        <div className="text-muted-foreground leading-none">
          Curva de respostas dos últimos 6 meses.
        </div>
      </CardFooter>
    </Card>
  );
}
