import * as React from "react";
import { CartesianGrid, Line, LineChart as RechartsLineChart, XAxis, YAxis } from "recharts";
import { DateRange } from "react-day-picker";
import { isWithinInterval, parseISO, format } from "date-fns";

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
  { date: "2024-01-01", engajamento: 120, satisfacao: 100 },
  { date: "2024-02-01", engajamento: 150, satisfacao: 130 },
  { date: "2024-03-01", engajamento: 130, satisfacao: 110 },
  { date: "2024-04-01", engajamento: 180, satisfacao: 160 },
  { date: "2024-05-01", engajamento: 200, satisfacao: 180 },
  { date: "2024-06-01", engajamento: 170, satisfacao: 150 },
  { date: "2024-07-01", engajamento: 220, satisfacao: 200 },
  { date: "2024-08-01", engajamento: 250, satisfacao: 230 },
  { date: "2024-09-01", engajamento: 230, satisfacao: 210 },
  { date: "2024-10-01", engajamento: 280, satisfacao: 260 },
  { date: "2024-11-01", engajamento: 300, satisfacao: 280 },
  { date: "2024-12-01", engajamento: 270, satisfacao: 250 },
];

const chartConfig = {
  engajamento: {
    label: "Engajamento",
    color: "hsl(var(--chart-1))",
  },
  satisfacao: {
    label: "Satisfação",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig;

interface ChartLineTrendsProps {
  dateRange?: DateRange;
}

export function ChartLineTrends({ dateRange }: ChartLineTrendsProps) {
  const filteredChartData = React.useMemo(() => {
    if (!dateRange?.from) {
      return allChartData;
    }
    const startDate = dateRange.from;
    const endDate = dateRange.to || new Date();

    return allChartData.filter(item => {
      const itemDate = parseISO(item.date);
      return isWithinInterval(itemDate, { start: startDate, end: endDate });
    });
  }, [dateRange]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Engajamento e Satisfação ao Longo do Tempo</CardTitle>
        <CardDescription>Métricas mensais de engajamento e satisfação.</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <RechartsLineChart
            accessibilityLayer
            data={filteredChartData}
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
              tickFormatter={(value) => format(parseISO(value), "MMM yy")}
            />
            <YAxis tickLine={false} axisLine={false} tickMargin={8} />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Line
              dataKey="engajamento"
              type="monotone"
              stroke="#2B7FFF"
              strokeWidth={2}
              dot={false}
            />
            <Line
              dataKey="satisfacao"
              type="monotone"
              stroke="#2B7FFF"
              strokeWidth={2}
              dot={false}
            />
          </RechartsLineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

