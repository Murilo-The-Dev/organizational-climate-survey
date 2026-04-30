"use client";

import React, { useEffect, useState, useMemo } from "react";
import { useParams } from "next/navigation";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  ResponsiveContainer,
  LineChart,
  Line,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
} from "recharts";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart";
import {
  FileText,
  Building,
  BarChart2,
  TrendingUp,
  Target,
  Layers,
  Shield,
  Clock,
  ListChecks,
  CheckCircle2,
  AlertTriangle,
  Lightbulb,
  Loader2,
} from "lucide-react";
import { pesquisaService } from "@/lib/services/pesquisaService";
import { dashboardService } from "@/lib/services/dashboardService";
import { useAuth } from "@/context/AuthContext";

const brandColor = "#2B7FFF";

// Alexandre: Mantive a tipagem de cores do gráfico, se adicionar novos setores no backend, avisa pra eu colocar mais cores aqui
const departmentColors = [
  brandColor,
  "#16A34A",
  "#DB2777",
  "#9333EA",
  "#15D9A8",
];

const getScoreClassification = (score: number): string => {
  if (score < 25) return "Crítico";
  if (score < 45) return "Ruim";
  if (score < 65) return "Médio";
  if (score < 85) return "Bom";
  return "Excelente";
};

const getClassificationBadgeColor = (classification: string) => {
  switch (classification) {
    case "Crítico":
      return "bg-red-600 hover:bg-red-700";
    case "Ruim":
      return "bg-orange-500 hover:bg-orange-600";
    case "Médio":
      return "bg-yellow-400 hover:bg-yellow-500";
    case "Bom":
      return "bg-sky-400 hover:bg-sky-500";
    case "Excelente":
      return "bg-green-500 hover:bg-green-600";
    default:
      return "bg-gray-400";
  }
};

const getRiskColor = (risk: string) => {
  switch (risk) {
    case "Crítico":
      return "bg-red-600";
    case "Ruim":
      return "bg-orange-500";
    case "Médio":
      return "bg-yellow-400";
    case "Bom":
      return "bg-sky-400";
    case "Excelente":
      return "bg-green-500";
    default:
      return "bg-gray-400";
  }
};

export default function RelatorioPage() {
  const params = useParams();
  const { user } = useAuth();

  // Tratamento do ID da pesquisa vindo da URL
  const rawId = params.id as string;
  const surveyId = rawId.replace(/\D/g, ""); // Extrai apenas números caso venha "PESQ-001"

  const [isLoading, setIsLoading] = useState(true);
  const [reportData, setReportData] = useState<any>(null);

  useEffect(() => {
    if (user && surveyId) {
      carregarRelatorioCompleto();
    }
  }, [user, surveyId]);

  const carregarRelatorioCompleto = async () => {
    try {
      setIsLoading(true);

      // Estou fazendo as chamadas em paralelo pra otimizar o carregamento da tela
      const [pesquisa, metricasData] = await Promise.all([
        pesquisaService.getById(Number(surveyId)).catch(() => null),
        dashboardService.getData(Number(surveyId)).catch(() => null), // Fallback se der erro
      ]);

      const metricas = metricasData?.metricas_por_pergunta || [];
      // Fiz um parser temporário no frontend pegando pelas métricas das perguntas. Quando a API evoluir, a gente limpa isso.
      const indicadoresFormatados = metricas.map((m: any, index: number) => ({
        nome: m.texto_pergunta
          ? m.texto_pergunta.substring(0, 20)
          : `Categoria ${index + 1}`,
        pontuacao: m.media ? Math.round((m.media / 5) * 100) : 0,
        risco: getScoreClassification(
          m.media ? Math.round((m.media / 5) * 100) : 0,
        ),
      }));

      const pontuacaoGeral =
        indicadoresFormatados.length > 0
          ? Math.round(
              indicadoresFormatados.reduce(
                (acc: number, cur: any) => acc + cur.pontuacao,
                0,
              ) / indicadoresFormatados.length,
            )
          : 0;

      // Montando o objeto final que os gráficos esperam
      setReportData({
        resumo: {
          pontuacaoGeral,
          participantes: metricasData?.total_respostas || 0,
          dataAvaliacao: pesquisa?.data_criacao
            ? new Date(pesquisa.data_criacao).toLocaleDateString("pt-BR")
            : new Date().toLocaleDateString("pt-BR"),
          proximaMedicao: "A definir",
        },
        infoEmpresa: {
          razaoSocial: (user as any)?.nome || "Empresa Logada",
          setores: ["Geral"], // TODO(Alexandre): Trazer os setores da API
          funcionarios: metricasData?.total_respostas || 0,
        },
        indicadores: indicadoresFormatados,
        // Mock temporário para os gráficos que exigem histórico e departamentos até o backend ter essas rotas prontas
        evolucao: [
          {
            data: "Atual",
            ...indicadoresFormatados.reduce(
              (acc: any, cur: any) => ({ ...acc, [cur.nome]: cur.pontuacao }),
              {},
            ),
          },
        ],
        desempenhoDepartamentos: [
          {
            departamento: "Geral",
            ...indicadoresFormatados.reduce(
              (acc: any, cur: any) => ({ ...acc, [cur.nome]: cur.pontuacao }),
              {},
            ),
          },
        ],
      });
    } catch (error) {
      console.error("Erro ao carregar o relatório:", error);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-gray-50/50">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="h-12 w-12 text-blue-600 animate-spin" />
          <p className="text-muted-foreground font-medium">
            Processando dados analíticos...
          </p>
        </div>
      </div>
    );
  }

  if (!reportData || reportData.indicadores.length === 0) {
    return (
      <section className="container mx-auto px-4 py-10">
        <Card>
          <CardHeader>
            <CardTitle>Relatório Indisponível</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Não há dados de resultados processados para a pesquisa selecionada
              no momento.
            </p>
          </CardContent>
        </Card>
      </section>
    );
  }

  // Configuração dinâmica para o Recharts não quebrar
  const chartConfig = reportData.indicadores.reduce(
    (acc: any, indicator: any, index: number) => {
      acc[indicator.nome] = {
        label: indicator.nome,
        color: departmentColors[index % departmentColors.length],
      };
      return acc;
    },
    {},
  );

  const sortedCategories = [...reportData.indicadores].sort(
    (a, b) => b.pontuacao - a.pontuacao,
  );
  const pontosFortes = sortedCategories
    .slice(0, 3)
    .filter((c) => c.pontuacao >= 70);
  const areasDeMelhoria = sortedCategories
    .slice(-3)
    .reverse()
    .filter((c) => c.pontuacao < 70);

  return (
    <section className="container mx-auto px-4 py-10 bg-gray-50/50">
      <header className="mb-8 print:hidden">
        <h1 className="text-4xl font-bold tracking-tight text-gray-800">
          Relatório Analítico
        </h1>
        <p className="text-muted-foreground mt-2 text-lg">
          Análise de Clima e Indicadores da Pesquisa
        </p>
      </header>

      <div className="flex flex-col gap-6">
        {/* Resumo Executivo */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText /> Resumo Executivo
            </CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="flex flex-col items-center justify-center p-4 rounded-lg bg-slate-100">
              <span
                className="text-3xl font-bold"
                style={{ color: brandColor }}
              >
                {reportData.resumo.pontuacaoGeral}%
              </span>
              <span className="text-sm text-muted-foreground">
                Pontuação Geral
              </span>
            </div>
            <div className="flex flex-col items-center justify-center p-4 rounded-lg">
              <span className="text-3xl font-bold">
                {reportData.resumo.participantes}
              </span>
              <span className="text-sm text-muted-foreground">
                Participantes
              </span>
            </div>
            <div className="flex flex-col items-center justify-center p-4 rounded-lg">
              <span className="text-xl font-semibold">
                {reportData.resumo.dataAvaliacao}
              </span>
              <span className="text-sm text-muted-foreground">
                Data da Avaliação
              </span>
            </div>
            <div className="flex flex-col items-center justify-center p-4 rounded-lg">
              <span className="text-xl font-semibold">
                {reportData.resumo.proximaMedicao}
              </span>
              <span className="text-sm text-muted-foreground">
                Próxima Medição
              </span>
            </div>
          </CardContent>
        </Card>

        {/* Comparação de Categorias */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Layers /> Desempenho por Categoria
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={chartConfig}
              className="min-h-[300px] w-full"
            >
              <ResponsiveContainer width="100%" height={300}>
                <BarChart
                  data={reportData.indicadores}
                  layout="horizontal"
                  margin={{ left: 0 }}
                  barSize={40}
                >
                  <YAxis
                    type="number"
                    domain={[0, 100]}
                    tickLine={false}
                    axisLine={false}
                  />
                  <XAxis
                    type="category"
                    dataKey="nome"
                    tickLine={false}
                    axisLine={false}
                    tickFormatter={(val) => val.substring(0, 10) + "..."}
                  />
                  <ChartTooltip
                    cursor={false}
                    content={
                      <ChartTooltipContent
                        valueFormatter={(value) => `${value}%`}
                      />
                    }
                  />
                  <Bar
                    dataKey="pontuacao"
                    fill={brandColor}
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </ChartContainer>
          </CardContent>
        </Card>

        {/* Análise Textual Rápida */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card className="border-green-200 bg-green-50/30">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-green-700">
                <CheckCircle2 /> Pontos Fortes
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {pontosFortes.length > 0 ? (
                pontosFortes.map((cat, i) => (
                  <div
                    key={i}
                    className="flex justify-between items-center border-b border-green-100 pb-2"
                  >
                    <span className="font-medium">{cat.nome}</span>
                    <Badge className="bg-green-600 text-white">
                      {cat.pontuacao}%
                    </Badge>
                  </div>
                ))
              ) : (
                <p className="text-sm text-muted-foreground">
                  Ainda não há dados suficientes com notas altas.
                </p>
              )}
            </CardContent>
          </Card>

          <Card className="border-orange-200 bg-orange-50/30">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-orange-700">
                <AlertTriangle /> Atenção Necessária
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {areasDeMelhoria.length > 0 ? (
                areasDeMelhoria.map((cat, i) => (
                  <div
                    key={i}
                    className="flex justify-between items-center border-b border-orange-100 pb-2"
                  >
                    <span className="font-medium">{cat.nome}</span>
                    <Badge className={getRiskColor(cat.risco) + " text-white"}>
                      {cat.pontuacao}%
                    </Badge>
                  </div>
                ))
              ) : (
                <p className="text-sm text-muted-foreground">
                  Nenhuma área crítica identificada nesta pesquisa.
                </p>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  );
}
