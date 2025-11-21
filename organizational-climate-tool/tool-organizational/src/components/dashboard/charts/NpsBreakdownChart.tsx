"use client";
import { Pie, PieChart } from "recharts";
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
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart";

const chartData = [
  { type: "promotores", count: 118, fill: "var(--color-blue-600)" },
  { type: "neutros", count: 40, fill: "var(--color-blue-500)" },
  { type: "detratores", count: 15, fill: "var(--color-blue-400)" },
];
const chartConfig = {
  count: { label: "Colaboradores" },
  promotores: { label: "Promotores", color: "var(--color-blue-600)" },
  neutros: { label: "Neutros", color: "var(--color-blue-500)" },
  detratores: { label: "Detratores", color: "var(--color-blue-400)" },
} satisfies ChartConfig;

export const NpsBreakdownChart = () => (
  <Card className="h-full grid grid-rows-[auto_1fr]">
    <CardHeader>
      <CardTitle>Detalhamento do e-NPS</CardTitle>
      <CardDescription>
        Distribuição de Promotores, Neutros e Detratores
      </CardDescription>
    </CardHeader>
    <CardContent>
      <ChartContainer
        config={chartConfig}
        className="mx-auto aspect-square h-full"
      >
        <PieChart>
          <ChartTooltip content={<ChartTooltipContent hideLabel />} />
          <Pie
            data={chartData}
            dataKey="count"
            nameKey="type"
            innerRadius="60%"
            outerRadius="80%"
            
          />
          <ChartLegend
            content={<ChartLegendContent nameKey="type" />}
            className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center"
          />
        </PieChart>
      </ChartContainer>
    </CardContent>
  </Card>
);
