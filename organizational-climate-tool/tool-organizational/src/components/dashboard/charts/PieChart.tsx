"use client"

import { TrendingUp } from "lucide-react"
import { Pie, PieChart } from "recharts"

import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart"

export const description = "A pie chart with a label"

const chartData = [
    { departamento: "Tecnologia da Informação", visitors: 275, fill: "var(--color-blue-100)" },
    { departamento: "Produção", visitors: 200, fill: "var(--color-blue-200)" },
    { departamento: "Recursos Humanos", visitors: 187, fill: "var(--color-blue-300)" },
    { departamento: "Vendas", visitors: 173, fill: "var(--color-blue-400)" },
    { departamento: "Engenharia", visitors: 90, fill: "var(--color-blue-500)" },
]

const chartConfig = {
    visitors: {
        label: "Departamentos",
    },
    chrome: {
        label: "Tecnologia da Informação",
        color: "var(--color-blue-100)",
    },
    safari: {
        label: "Produção",
        color: "var(--color-blue-200)",
    },
    firefox: {
        label: "Recursos Humanos",
        color: "var(color-blue-300)",
    },
    edge: {
        label: "Vendas",
        color: "var(color-blue-400)",
    },
    other: {
        label: "Engenharia",
        color: "var(--color-blue-500)",
    },
} satisfies ChartConfig

export function ChartPieLabel() {
    return (
        <Card className="flex flex-col">
            <CardHeader className="items-center pb-0">
                <CardTitle>Distribuição de Departamentos</CardTitle>
                <CardDescription>Janeiro - Junho 2024</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 pb-0">
                <ChartContainer
                    config={chartConfig}
                    className="[&_.recharts-pie-label-text]:fill-foreground mx-auto aspect-square max-h-[250px] pb-0"
                >
                    <PieChart>
                        <ChartTooltip content={<ChartTooltipContent hideLabel />} />
                        <Pie data={chartData} dataKey="visitors" label nameKey="departamento" />
                    </PieChart>
                </ChartContainer>
            </CardContent>
            <CardFooter className="flex-col gap-2 text-sm">
                <div className="flex items-center gap-2 leading-none font-medium">
                    Aumento de 5.2% este mês <TrendingUp className="h-4 w-4" />
                </div>
                <div className="text-muted-foreground leading-none">
                    Mostrando o total dos Departamentos que responderam a pesquisa
                </div>
            </CardFooter>
        </Card>
    )
}
