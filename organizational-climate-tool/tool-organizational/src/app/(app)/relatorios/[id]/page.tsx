"use client";

import React, { useMemo } from "react";
import { useParams } from "next/navigation";
import { dadosPesquisas } from "@/components/dashboard/DataTable";
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"; 
import { Badge } from "@/components/ui/badge";
import { Separator } from "@radix-ui/react-dropdown-menu";
import { BarChart, Bar, XAxis, YAxis, ResponsiveContainer, LineChart, Line, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar } from "recharts";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { reportData as mockReportData } from "@/components/dashboard/ReportMockData";
import { allMockResults } from "@/components/dashboard/ResultsDataTable";
import { FileText, Building, BarChart2, TrendingUp, Target, Layers, Shield, Clock, ListChecks, CheckCircle, BrainCircuit, HeartHandshake, UserCheck, Smile, MessageCircle, ShieldCheck, Award, Star, ThumbsUp, LogOut, CalendarX2, AlertTriangle, CheckCircle2, Lightbulb } from "lucide-react";
import { IndividualScoresScatterPlot, mockApiData as allIndividualResponses, RespostaIndividual } from "@/components/dashboard/IndividualScoresScatterPlot";

const brandColor = "#2B7FFF";

const chartConfig = {
  pontuacao: { label: "Pontuação", color: brandColor },
  Liderança: { label: "Liderança", color: brandColor },
  "Bem-estar": { label: "Bem-estar", color: "#579BFF" }, // Tom mais claro do azul da marca
  Comunicação: { label: "Comunicação", color: "#0A4DB2" }, // Tom mais escuro do azul da marca
  score: { label: "Pontuação do Setor", color: brandColor },
} satisfies ChartConfig;

const getRiskColor = (risk: string) => {
  switch (risk) {
    case "Crítico": return "bg-red-600";
    case "Ruim": return "bg-orange-500";
    case "Médio": return "bg-yellow-400";
    case "Bom": return "bg-sky-400";
    case "Excelente": return "bg-green-500";
    default: return "bg-gray-400";
  }
};

const getIndicatorIcon = (name: string) => {
    const iconProps = { className: "h-6 w-6 text-muted-foreground" };
    switch (name) {
        case "Engajamento": return <HeartHandshake {...iconProps} />;
        case "Liderança": return <UserCheck {...iconProps} />;
        case "Bem-estar": return <Smile {...iconProps} />;
        case "Comunicação": return <MessageCircle {...iconProps} />;
        case "Segurança Psicológica": return <ShieldCheck {...iconProps} />;
        case "Crescimento": return <TrendingUp {...iconProps} />;
        case "Reconhecimento": return <Award {...iconProps} />;
        case "Satisfação Geral": return <Star {...iconProps} />;
        case "e-NPS": return <ThumbsUp {...iconProps} />;
        case "Intenção de Rotatividade": return <LogOut {...iconProps} />;
        case "Absenteísmo": return <CalendarX2 {...iconProps} />;
        case "Turnover Voluntário": return <LogOut {...iconProps} />;
        default: return <BarChart2 {...iconProps} />;
    }
};

const getScoreClassification = (score: number): string => {
  if (score < 25) return "Crítico";
  if (score < 45) return "Ruim";
  if (score < 65) return "Médio";
  if (score < 85) return "Bom";
  return "Excelente";
};

const getClassificationBadgeColor = (classification: string) => {
    switch (classification) {
        case "Crítico": return "bg-red-600 hover:bg-red-700";
        case "Ruim": return "bg-orange-500 hover:bg-orange-600";
        case "Médio": return "bg-yellow-400 hover:bg-yellow-500";
        case "Bom": return "bg-sky-400 hover:bg-sky-500";
        case "Excelente": return "bg-green-500 hover:bg-green-600";
        default: return "bg-gray-400";
    }
};

const getRiskIdentificationText = (category: string): string => {
  switch (category) {
    case "Liderança":
      return "Falta de apoio emocional e segurança nas relações hierárquicas.";
    case "Bem-estar":
      return "Desgaste emocional e físico, impactando a saúde e a produtividade.";
    case "Comunicação":
      return "Falhas na comunicação que geram retrabalho, conflitos e desmotivação.";
    case "Engajamento":
      return "Restrição da criatividade e perda de conexão com o trabalho.";
    default:
      return "Risco não especificado identificado para esta categoria.";
  }
};

const getCategoryAnalysisText = (category: string, score: number): string => {
  const baseAnalysis = (performance: string, details: string) => `A atuação em ${category} é ${performance}, ${details}.`;

  if (score >= 85) { // Excelente
    switch (category) {
      case "Liderança": return baseAnalysis("excelente", "com uma liderança inspiradora que é referência em motivar e desenvolver equipes");
      case "Bem-estar": return baseAnalysis("excelente", "promovendo um ambiente de trabalho excepcionalmente saudável e equilibrado, que é um grande diferencial");
      case "Comunicação": return baseAnalysis("excelente", "com fluxos de informação claros e eficientes que fortalecem a colaboração e a confiança em todos os níveis");
      case "Engajamento": return baseAnalysis("excelente", "com equipes altamente conectadas e motivadas, que demonstram um forte senso de propósito e pertencimento");
      default: return baseAnalysis("excelente", "com práticas consolidadas que servem de referência e contribuem significativamente para um ambiente de trabalho positivo");
    }
  }
  if (score >= 65) { // Bom
    switch (category) {
      case "Liderança": return baseAnalysis("boa", "demonstrando uma gestão eficaz que apoia e orienta as equipes de forma consistente");
      case "Bem-estar": return baseAnalysis("boa", "com iniciativas que promovem um equilíbrio saudável entre vida profissional e pessoal");
      case "Comunicação": return baseAnalysis("boa", "garantindo que as informações importantes sejam compartilhadas de forma clara e acessível");
      case "Engajamento": return baseAnalysis("boa", "com colaboradores que se sentem valorizados e conectados com os objetivos da empresa");
      default: return baseAnalysis("boa", "demonstrando um desempenho consistente e resultados positivos que devem ser mantidos e reforçados");
    }
  }
  if (score >= 45) { // Médio
    switch (category) {
      case "Liderança": return baseAnalysis("média", "sinalizando que, embora haja pontos positivos, há oportunidades para fortalecer o apoio e a direção das lideranças");
      case "Bem-estar": return baseAnalysis("média", "indicando a necessidade de maior atenção ao equilíbrio e à saúde mental dos colaboradores para evitar o desgaste");
      case "Comunicação": return baseAnalysis("média", "apontando para ruídos ou falhas que podem estar afetando o alinhamento e a colaboração entre as equipes");
      case "Engajamento": return baseAnalysis("média", "sugerindo que parte dos colaboradores pode não se sentir totalmente conectada ou motivada, necessitando de ações de incentivo");
      default: return baseAnalysis("média", "sinalizando pontos de atenção que precisam ser aprimorados para garantir maior consistência e evitar a degradação do clima");
    }
  }
  if (score >= 25) { // Ruim
    switch (category) {
      case "Liderança": return baseAnalysis("ruim", "evidenciando uma necessidade urgente de desenvolvimento de líderes para evitar desmotivação e desalinhamento nas equipes");
      case "Bem-estar": return baseAnalysis("ruim", "com sinais claros de esgotamento e insatisfação que podem comprometer a saúde dos colaboradores e a produtividade");
      case "Comunicação": return baseAnalysis("ruim", "com falhas significativas que estão gerando conflitos, retrabalho e um clima de desconfiança");
      case "Engajamento": return baseAnalysis("ruim", "com um baixo nível de conexão dos colaboradores, o que representa um risco para a retenção de talentos e para os resultados");
      default: return baseAnalysis("ruim", "evidenciando vulnerabilidades que comprometem a eficiência e o bem-estar, demandando ações corretivas imediatas");
    }
  }
  // Crítico
  switch (category) {
    case "Liderança": return baseAnalysis("crítica", "com falhas graves na gestão que exigem intervenção imediata para reverter um cenário de alto risco para a cultura e resultados");
    case "Bem-estar": return baseAnalysis("crítica", "com um ambiente de trabalho que apresenta riscos psicossociais severos, necessitando de ações emergenciais para proteger os colaboradores");
    case "Comunicação": return baseAnalysis("crítica", "com um colapso nos fluxos de informação que está minando a confiança e a capacidade de execução da empresa");
    case "Engajamento": return baseAnalysis("crítica", "com um nível alarmante de desmotivação e desconexão, indicando problemas profundos que precisam ser endereçados com urgência");
    default: return baseAnalysis("crítica", "com falhas graves que exigem intervenção urgente para mitigar riscos psicossociais severos");
  }
};

const RelatorioPage = () => {
  const params = useParams();
  const surveyId = params.id as string;
  const survey = dadosPesquisas.find((p) => p.id === surveyId);
  const surveyTitle = survey ? survey.title : "Relatório Geral";

  // Converte a pontuação para uma escala de 0-100.
  // TODO: A pontuação máxima (5) está hardcoded. O ideal é que a fonte de dados forneça a pontuação máxima para cada categoria.
  const scaleScore = (score: number) => Math.round((score / 5) * 100);

  const currentSurveyId = useMemo(() => `survey-${surveyId.split('-')[1]}`, [surveyId]);
  const currentSurveyResponses = useMemo(() => allIndividualResponses.filter(res => res.pesquisaId === currentSurveyId), [currentSurveyId]);

  // Calcula os indicadores de categoria a partir das respostas individuais para maior precisão
  const indicadoresCategorias = useMemo(() => {
    const scoresFromIndividualResponses = currentSurveyResponses.reduce((acc, res) => {
      if (!acc[res.categoria]) {
        acc[res.categoria] = { totalScore: 0, count: 0, maxScore: res.pontuacaoMaxima };
      }
      acc[res.categoria].totalScore += res.pontuacao;
      acc[res.categoria].count += 1;
      return acc;
    }, {} as { [category: string]: { totalScore: number; count: number; maxScore: number } });

    return Object.entries(scoresFromIndividualResponses).map(([category, data]) => {
      const average = data.count > 0 ? data.totalScore / data.count : 0;
      const scaledScore = Math.round((average / data.maxScore) * 100);
      return { nome: category, pontuacao: scaledScore };
    });
  }, [currentSurveyResponses]);

  // Mock de uma "fonte de dados externa". Em um app real, isso viria de um banco de dados.
  // O usuário inseriria esses dados através do botão na tela de listagem de pesquisas.
  const externalIndicatorsData: { [surveyId: string]: { [indicatorName: string]: { value: number, isInverted: boolean } } } = {
    'dps-001': { // Dados inseridos apenas para a pesquisa 'dps-001'
      'Absenteísmo': { value: 5, isInverted: true }, // 5% - valor alto é ruim
      'Turnover Voluntário': { value: 12, isInverted: true } // 12% - valor alto é ruim
    },
    // A pesquisa 'dps-002' não tem dados aqui, então esses indicadores não aparecerão para ela.
  };

  // --- Início do Cálculo de Indicadores Especiais (e-NPS, Rotatividade, etc.) ---
  const specialIndicators = useMemo(() => {
    const indicators: { nome: string; pontuacao: number }[] = [];
  // 1. Satisfação Geral: Calculada como a média de todas as categorias da pesquisa.
  if (indicadoresCategorias.length > 0) {
    const pontuacaoGeralSatisfacao = Math.round(indicadoresCategorias.reduce((acc, cat) => acc + cat.pontuacao, 0) / indicadoresCategorias.length);
    specialIndicators.push({
        nome: "Satisfação Geral",
        pontuacao: pontuacaoGeralSatisfacao,
    });
  }

  // 2. e-NPS (Employee Net Promoter Score) - Aparece se houver a pergunta na pesquisa.
  const eNPSResponses = currentSurveyResponses.filter(res => res.categoria === 'e-NPS');
  if (eNPSResponses.length > 0) {
      const totalRespondents = eNPSResponses.length;
      const maxScoreForENPS = eNPSResponses[0].pontuacaoMaxima > 0 ? eNPSResponses[0].pontuacaoMaxima : 10;
      const promoters = eNPSResponses.filter(r => (r.pontuacao / maxScoreForENPS) * 10 >= 9).length;
      const detractors = eNPSResponses.filter(r => (r.pontuacao / maxScoreForENPS) * 10 <= 6).length;
      const promoterPercentage = (promoters / totalRespondents) * 100;
      const detractorPercentage = (detractors / totalRespondents) * 100;
      const eNPSScore = Math.round(promoterPercentage - detractorPercentage);
      const normalizedENPS = Math.round((eNPSScore + 100) / 2);
      specialIndicators.push({ nome: "e-NPS", pontuacao: normalizedENPS });
  }

  // 3. Intenção de Rotatividade - Aparece se houver a pergunta na pesquisa.
  const turnoverResponses = currentSurveyResponses.filter(res => res.categoria === 'Intenção de Rotatividade');
  if (turnoverResponses.length > 0) {
      const totalScore = turnoverResponses.reduce((acc, res) => acc + res.pontuacao, 0);
      const averageScore = totalScore / turnoverResponses.length;
      const maxPossibleScore = turnoverResponses[0].pontuacaoMaxima;
      const invertedAverage = maxPossibleScore - averageScore;
      const turnoverIndicatorScore = Math.round((invertedAverage / maxPossibleScore) * 100);
      specialIndicators.push({ nome: "Intenção de Rotatividade", pontuacao: turnoverIndicatorScore });
  }

  // 4. Indicadores Externos (Absenteísmo, etc.) - Aparecem se forem inseridos manualmente.
  const surveyExternalData = externalIndicatorsData[surveyId] || {};
  Object.entries(surveyExternalData).forEach(([name, data]) => {
    let score = data.isInverted ? 100 - data.value : data.value;
    indicators.push({ nome: name, pontuacao: Math.round(score) });
  });
    return indicators;
  }, [indicadoresCategorias, currentSurveyResponses, surveyId]);

  // --- Fim do Cálculo de Indicadores Especiais ---

  // --- Início da Refatoração para o Gráfico de Evolução ---

  // 1. Processar os dados brutos para calcular a evolução mensal por categoria
  const getMonthYear = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', { month: 'short', year: '2-digit' }).replace('.', '').replace(' de ', '/');
  };

  const evolucaoData = useMemo(() => {
  // Agrupa os resultados de todas as pesquisas do mesmo TIPO (ex: "Diagnóstico")
  // para mostrar uma linha de tendência real entre as medições.
  const surveyType = survey?.type;
  const relatedSurveyIds = dadosPesquisas
    .filter(p => p.type === surveyType && p.id) // Garante que o tipo e o ID existam
    .map(p => `survey-${p.id.split('-')[1]}`);

  const allRelatedResponses = allIndividualResponses.filter(res =>
    relatedSurveyIds.includes(res.pesquisaId)
  );

  // Agrupa todas as respostas relacionadas por mês/ano, que representam uma medição.
  const responsesByMonth = allRelatedResponses.reduce((acc, response) => {
    const monthYear = getMonthYear(response.respondidoEm);
    if (!acc[monthYear]) {
      acc[monthYear] = [];
    }
    acc[monthYear].push(response);
    return acc;
  }, {} as { [month: string]: RespostaIndividual[] });

  const meses: { [key: string]: number } = { 'jan': 1, 'fev': 2, 'mar': 3, 'abr': 4, 'mai': 5, 'jun': 6, 'jul': 7, 'ago': 8, 'set': 9, 'out': 10, 'nov': 11, 'dez': 12 };

  // Ordena as medições por data para garantir a ordem cronológica (Medição 1, 2, 3...).
  const sortedMeasurements = Object.entries(responsesByMonth).sort(([monthA], [monthB]) => {
    const [m1Str, y1] = monthA.toLowerCase().split('/');
    const [m2Str, y2] = monthB.toLowerCase().split('/');
    const dateA = new Date(parseInt(`20${y1}`), meses[m1Str] - 1);
    const dateB = new Date(parseInt(`20${y2}`), meses[m2Str] - 1);
    return dateA.getTime() - dateB.getTime();
  });

  // Mapeia os dados ordenados para calcular as médias e formatar o rótulo do eixo X.
    let data = sortedMeasurements.map(([monthYear, responses], index) => {
    const categoriesInMonth: { [category: string]: { totalScore: number, count: number } } = {};

    responses.forEach(res => {
      if (!categoriesInMonth[res.categoria]) {
        categoriesInMonth[res.categoria] = { totalScore: 0, count: 0 };
      }
      categoriesInMonth[res.categoria].totalScore += res.pontuacao;
      categoriesInMonth[res.categoria].count += 1;
    });

    const monthlyAverages: { [key: string]: string | number } = {
      data: `Medição ${index + 1} (${monthYear})` // Rótulo com contagem e data
    };
    Object.entries(categoriesInMonth).forEach(([category, data]) => {
      // A pontuação máxima por categoria deve ser consistente. Pega do primeiro registro encontrado.
      const maxScore = responses.find(r => r.categoria === category)?.pontuacaoMaxima || 5; // Fallback para 5
      const average = data.count > 0 ? data.totalScore / data.count : 0;
      monthlyAverages[category] = Math.round((average / maxScore) * 100);
    });
    return monthlyAverages;
  });

  // Se não houver dados históricos (nenhuma medição encontrada), usa os dados da pesquisa ATUAL
  // como o "ponto zero" (Medição 1). Isso garante que o gráfico não fique vazio.
    if (data.length === 0 && indicadoresCategorias.length > 0 && survey?.dataCriacao) {
    const currentSurveyDataPoint: { [key: string]: string | number } = {
        data: `Medição 1 (${getMonthYear(survey.dataCriacao)})`
    };
    indicadoresCategorias.forEach(cat => {
        currentSurveyDataPoint[cat.nome] = cat.pontuacao;
    });
      data.push(currentSurveyDataPoint);
    }
    return data;
  }, [survey, allIndividualResponses, dadosPesquisas, indicadoresCategorias]);

  // --- Fim da Refatoração ---

  // Carrega os resultados "reais" da pesquisa com base no ID
  const surveyResults = allMockResults[surveyId] || [];

  // Extrai e calcula os dados reais por departamento para o gráfico de radar
  const { desempenhoDepartamentos, detalhamentoDepartamentos, surveyCategories } = useMemo(() => {
    const surveyCategoriesSet = new Set(surveyResults.map(r => r.category));

  // Agrega as pontuações por departamento e categoria a partir dos dados reais da pesquisa
  const scoresByDeptAndCat = surveyResults.reduce((acc, result) => {
    if (!result.departmentScores) return acc;

    Object.entries(result.departmentScores).forEach(([deptName, deptScore]) => {
        if (!acc[deptName]) {
            acc[deptName] = {};
        }
        if (!acc[deptName][result.category]) {
            acc[deptName][result.category] = [];
        }
        acc[deptName][result.category].push(deptScore.averageScore);
    });

    return acc;
  }, {} as { [dept: string]: { [cat: string]: number[] } });

  // Calcula a média para cada categoria dentro de cada departamento
    const desempenho = Object.entries(scoresByDeptAndCat).map(([deptName, categories]) => {
      const deptScores: { [category: string]: number } = {};
      surveyCategoriesSet.forEach(catName => {
          const scores = categories[catName] || [];
          const average = scores.length > 0 ? scores.reduce((a, b) => a + b, 0) / scores.length : 0;
          deptScores[catName] = scaleScore(average);
      });

      return { departamento: deptName, ...deptScores };
  });

    const detalhamento = desempenho.map((deptData, index) => {
      const scores = Object.values(deptData).filter(v => typeof v === 'number') as number[];
      const pontuacaoGeral = scores.length > 0 ? Math.round(scores.reduce((a, b) => a + b, 0) / scores.length) : 0;
  
      const row: {[key: string]: any} = {
          id: `dept-${index}`,
          setor: deptData.departamento,
          pontuacaoGeral: pontuacaoGeral,
      };
      surveyCategoriesSet.forEach(cat => {
        row[cat] = (deptData[cat] as number) || 0;
      });
      return row;
  });
    return { desempenhoDepartamentos: desempenho, detalhamentoDepartamentos: detalhamento, surveyCategories: [...surveyCategoriesSet] };
  }, [surveyResults]);

  
  const setoresDinamicos = desempenhoDepartamentos.map(d => d.departamento);

  // Monta o objeto de dados do relatório dinamicamente
  const reportData = {
    ...mockReportData,
    infoEmpresa: {
      ...mockReportData.infoEmpresa,
      setores: setoresDinamicos,
      funcionarios: survey?.participantes ?? 0,
    },
    resumo: {
      ...mockReportData.resumo,
      pontuacaoGeral: indicadoresCategorias.length > 0 ? Math.round(indicadoresCategorias.reduce((acc, cat) => acc + cat.pontuacao, 0) / indicadoresCategorias.length) : 0,
      participantes: survey?.participantes ?? 0,
      dataAvaliacao: survey?.dataCriacao ? new Date(survey.dataCriacao).toLocaleDateString("pt-BR") : mockReportData.resumo.dataAvaliacao,
    },
    evolucao: evolucaoData, // Utiliza os dados de evolução processados
    indicadores: {
      categorias: indicadoresCategorias,
      especiais: specialIndicators,
    },
    desempenhoDepartamentos,
    detalhamentoDepartamentos,
    matrizRisco: indicadoresCategorias.map((cat) => ({
        categoria: cat.nome,
        risco: getScoreClassification(cat.pontuacao),
    })),
    mapaArvore: indicadoresCategorias.map((cat) => ({
        name: cat.nome,
        size: cat.pontuacao,
    })),
  };

  const departmentColors = [
    brandColor,
    "#16A34A",
    "#DB2777",
    "#9333EA",
    "#15D9A8",
  ];

  // Criar configuração de gráfico dinâmica para não mutar o objeto original
  const { dynamicChartConfig, allEvolutionCategories } = useMemo(() => {
    const newConfig = { ...chartConfig };
    const categories = evolucaoData.reduce((acc, monthData) => {
      Object.keys(monthData).forEach(key => {
        if (key !== 'data' && !acc.includes(key)) {
          acc.push(key);
        }
      });
      return acc;
    }, [] as string[]);

    categories.forEach((category, index) => {
      if (!newConfig[category]) {
        newConfig[category] = { label: category, color: departmentColors[index % departmentColors.length] };
      }
    });
    return { dynamicChartConfig: newConfig, allEvolutionCategories: categories };
  }, [evolucaoData]);

  // Dados para seções dinâmicas
  const pontuacaoGeral = reportData.resumo.pontuacaoGeral;
  const sortedCategories = [...indicadoresCategorias].sort((a, b) => b.pontuacao - a.pontuacao);
  const pontosFortes = sortedCategories.slice(0, 3).filter(c => c.pontuacao >= 70);
  const areasDeMelhoria = sortedCategories.slice(-3).reverse().filter(c => c.pontuacao < 70);
  const piorCategoria = areasDeMelhoria[0];

  const proximaDataMedicao = survey?.dataCriacao
    ? new Date(new Date(survey.dataCriacao).setMonth(new Date(survey.dataCriacao).getMonth() + 6)).toLocaleDateString("pt-BR")
    : "N/A";

  const recomendacoesDinamicas = areasDeMelhoria.map(cat => ({
      area: cat.nome,
      risco: getScoreClassification(cat.pontuacao),
      plano: mockReportData.recomendacoes.find(r => r.area === cat.nome)?.plano || "Definir e comunicar um plano de ação específico para esta área."
  }));

  const analisePorSetor = desempenhoDepartamentos.map(dept => {
    const allCategories = Object.entries(dept).filter(([key]) => key !== 'departamento');

    const pontosFortes = allCategories
      .map(([cat, score]) => ({ nome: cat, pontuacao: score as number }))
      .filter(c => c.pontuacao >= 65) // Bom ou Excelente
      .sort((a, b) => b.pontuacao - a.pontuacao);

    const areasDeMelhoria = allCategories
      .map(([cat, score]) => ({ nome: cat, pontuacao: score as number }))
      .filter(c => c.pontuacao < 65) // Médio, Ruim, Crítico
      .sort((a, b) => a.pontuacao - b.pontuacao);

    return {
      setor: dept.departamento,
      pontosFortes,
      areasDeMelhoria,
    };
  });

  const tableCategories = [...indicadoresCategorias]
    .sort((a, b) => a.pontuacao - b.pontuacao)
    .slice(0, 2)
    .map(c => c.nome);

  // Configuração do Gráfico de Radar
  const radarChartConfig = { ...chartConfig };
  reportData.desempenhoDepartamentos.forEach((dept, index) => {
    radarChartConfig[dept.departamento] = {
      label: dept.departamento,
      color: departmentColors[index % departmentColors.length],
    };
  });

  const radarCategories = Object.keys(reportData.desempenhoDepartamentos[0] || {}).filter(
    (key) => key !== "departamento"
  );

  if (surveyResults.length === 0) {
    return (
      <section className="container mx-auto px-4 py-10">
        <Card>
          <CardHeader><CardTitle>Relatório Indisponível</CardTitle></CardHeader>
          <CardContent><p>Não há dados de resultados para a pesquisa selecionada.</p></CardContent>
        </Card>
      </section>
    );
  }

  const multiRadarData = radarCategories.map((category) => {
    const categoryData: { [key: string]: string | number } = { category };
    reportData.desempenhoDepartamentos.forEach((dept) => {
      categoryData[dept.departamento] = dept[category as keyof typeof dept];
    });
    return categoryData;
  });

  return (
    <section className="container mx-auto px-4 py-10 bg-gray-50/50">
      <header className="mb-8 print:hidden">
        <h1 className="text-4xl font-bold tracking-tight text-gray-800">{surveyTitle}</h1>
        <p className="text-muted-foreground mt-2 text-lg">Relatório de Análise Psicossocial e de Clima Organizacional</p>
      </header>

      <div className="flex flex-col gap-6">
        {/* Resumo Executivo */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2"><FileText /> Resumo Executivo</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="flex flex-col items-center justify-center p-4 rounded-lg bg-slate-100">
              <span className="text-3xl font-bold" style={{ color: brandColor }}>{reportData.resumo.pontuacaoGeral}%</span>
              <span className="text-sm text-muted-foreground">Pontuação Geral</span>
            </div>
            <div className="flex flex-col items-center justify-center p-4 rounded-lg">
              <span className="text-3xl font-bold">{reportData.resumo.participantes}</span>
              <span className="text-sm text-muted-foreground">Participantes</span>
            </div>
            <div className="flex flex-col items-center justify-center p-4 rounded-lg">
              <span className="text-xl font-semibold">{reportData.resumo.dataAvaliacao}</span>
              <span className="text-sm text-muted-foreground">Data da Avaliação</span>
            </div>
            <div className="flex flex-col items-center justify-center p-4 rounded-lg">
              <span className="text-xl font-semibold">{reportData.resumo.proximaMedicao}</span>
              <span className="text-sm text-muted-foreground">Próxima Medição</span>
            </div>
          </CardContent>
        </Card>

        {/* Informações da Empresa */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2"><Building /> Informações da Empresa</CardTitle>
          </CardHeader>
          <CardContent className="text-sm space-y-2">
            <p><strong>Razão Social:</strong> {reportData.infoEmpresa.razaoSocial}</p>
            <p><strong>CNPJ:</strong> {reportData.infoEmpresa.cnpj}</p>
            <p><strong>Nº de Funcionários:</strong> {reportData.infoEmpresa.funcionarios}</p>
            <p><strong>Setores Avaliados:</strong> {reportData.infoEmpresa.setores.join(", ")}</p>
          </CardContent>
        </Card>

        {/* Resumo Executivo Textual */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText /> Resumo Executivo
            </CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground space-y-2">
            <p>
              Análise abrangente do ambiente psicossocial da empresa {reportData.infoEmpresa.razaoSocial}, que conta com {reportData.infoEmpresa.funcionarios} funcionários distribuídos em {reportData.infoEmpresa.setores.length} setores: {reportData.infoEmpresa.setores.join(', ')}.
            </p>
            <p>
              A pontuação geral obtida foi de {(reportData.resumo.pontuacaoGeral / 10).toFixed(1)}/10, baseada na avaliação de todos os setores da empresa.
            </p>
            <p>
              Data da última solicitação: {reportData.resumo.dataAvaliacao}
            </p>
            <br />
            <p>
              Próxima data sugerida para medição: {proximaDataMedicao}
            </p>
          </CardContent>
        </Card>

        {/* Indicadores Principais */}
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2"><BarChart2 /> Indicadores Principais</CardTitle>
                <CardDescription>Visão geral do desempenho em cada categoria e indicador de satisfação.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                    {[...reportData.indicadores.categorias, ...reportData.indicadores.especiais].map((indicator) => {
                        const classification = getScoreClassification(indicator.pontuacao);
                        return (
                            <Card key={indicator.nome} className="flex flex-col">
                                <CardHeader className="flex flex-row items-center justify-between pb-2">
                                    <CardTitle className="text-base font-medium">{indicator.nome}</CardTitle>
                                    {getIndicatorIcon(indicator.nome)}
                                </CardHeader>
                                <CardContent className="flex-grow">
                                    <div className="text-4xl font-bold" style={{ color: brandColor }}>{indicator.pontuacao}%</div>
                                </CardContent>
                                <CardFooter>
                                    <Badge className={`${getClassificationBadgeColor(classification)} text-white`}>
                                        {classification}
                                    </Badge>
                                </CardFooter>
                            </Card>
                        );
                    })}
                </div>
            </CardContent>
        </Card>

        {/* Comparação de Categorias */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2"><Layers /> Comparação de Categorias</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="min-h-[300px] w-full">
              <ResponsiveContainer width="100%" height={300}>
                <BarChart 
                  data={reportData.indicadores.categorias} 
                  layout="horizontal" 
                  margin={{ left: 30 }}
                  barSize={40}
                  barCategoryGap="5%"
                >
                  <YAxis 
                    type="number" 
                    dataKey="pontuacao" 
                    domain={[0, 100]}
                    tickLine={true}
                    axisLine={false}
                  />
                  <XAxis 
                    width={60} 
                    type="category" 
                    dataKey="nome" 
                    tickLine={false}
                    axisLine={false} 
                  />                  
                  <ChartTooltip 
                    cursor={false}
                    content={<ChartTooltipContent valueFormatter={(value) => `${value}%`} />} 
                  />
                  <Bar dataKey="pontuacao" fill="var(--color-pontuacao)" radius={8} />
                </BarChart>
              </ResponsiveContainer>
            </ChartContainer>
          </CardContent>
        </Card>

        {/* Análise por Setor (Radar) */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Target /> Análise por Setor
            </CardTitle>
            <CardDescription>
              Comparativo de desempenho entre os diferentes setores da empresa em um único gráfico.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={radarChartConfig} className="mx-auto aspect-square max-h-[600px]">
              <ResponsiveContainer width="100%" height="100%">
                <RadarChart data={multiRadarData} margin={{ top: 80, right: 120, bottom: 80, left: 120 }}>
                  <ChartTooltip cursor={false} content={<ChartTooltipContent valueFormatter={(value) => `${value}%`} />} />
                  <ChartLegend content={<ChartLegendContent />} />
                  {reportData.desempenhoDepartamentos.map((dept) => (
                    <Radar 
                      key={dept.departamento} 
                      name={dept.departamento} 
                      dataKey={dept.departamento} 
                      fill={(radarChartConfig as any)[dept.departamento]?.color}
                      fillOpacity={0.2} 
                      stroke={radarChartConfig[dept.departamento]?.color}
                      strokeWidth={2}
                    />
                  ))}
                  <PolarGrid />
                  <PolarAngleAxis dataKey="category" tick={{ fontSize: 12 }} />
                  <PolarRadiusAxis angle={30} domain={[0, 100]} />
                </RadarChart>
              </ResponsiveContainer>
            </ChartContainer>
          </CardContent>
        </Card>

        {/* Evolução das Pontuações */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2"><TrendingUp /> Evolução das Pontuações (Análise de Tendência)</CardTitle>
          </CardHeader>
          <CardContent>
            {reportData.evolucao.length > 1 ? (
              <ChartContainer config={dynamicChartConfig} className="min-h-[300px] w-full">
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={reportData.evolucao}>
                    <XAxis 
                      dataKey="data" 
                      tickLine={true}
                      axisLine={false}
                    />
                    <YAxis 
                      domain={[0, 100]} 
                      tickLine={true}
                      axisLine={false}
                    />
                    <ChartTooltip content={<ChartTooltipContent valueFormatter={(value) => `${value}%`} />} />
                    <ChartLegend content={<ChartLegendContent />} />
                    
                    {/* 3. Renderizar as linhas do gráfico dinamicamente */}
                    {allEvolutionCategories.map(category => (
                      <Line 
                        key={category}
                        type="monotone" 
                        dataKey={category} 
                        stroke={`var(--color-${category})`}
                        strokeWidth={2} 
                        dot={true} 
                      />
                    ))}
                  </LineChart>
                </ResponsiveContainer>
              </ChartContainer>
            ) : (
              <p className="text-muted-foreground text-center">O gráfico de tendência estará disponível a partir da segunda medição.<br />Ainda não há histórico de evolução para esta pesquisa.</p>
            )}
          </CardContent>
        </Card>

        {/* Gráfico de Dispersão de Pontuações Individuais */}
        <IndividualScoresScatterPlot />

        {/* Detalhamento por Departamento */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2"><Layers /> Detalhamento por Departamento</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Setor</TableHead>
                  <TableHead className="text-right">Pontuação Geral</TableHead>
                  {tableCategories.map(cat => (
                      <TableHead key={cat} className="text-right">{cat}</TableHead>
                  ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                {reportData.detalhamentoDepartamentos.map((depto) => (
                  <TableRow key={depto.id}>
                    <TableCell className="font-medium">{depto.setor}</TableCell>
                    <TableCell className={`text-right font-bold`}>{depto.pontuacaoGeral}%</TableCell>
                    {tableCategories.map(cat => (
                        <TableCell key={cat} className="text-right">{depto[cat]}%</TableCell>
                    ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        {/* Matriz de Risco */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2"><Shield /> Matriz de Risco</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {reportData.matrizRisco.map(item => (
                <div key={item.categoria} className="flex items-center justify-between text-sm p-2 rounded-md bg-gray-50">
                    <span>{item.categoria}</span>
                    <Badge className={`${getRiskColor(item.risco)} text-white`}>{item.risco}</Badge>
                </div>
            ))}
          </CardContent>
        </Card>

        {/* Mapa de Árvore */}
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2"><BrainCircuit /> Mapa de Fatores Psicossociais</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="w-full h-48 flex flex-wrap gap-1 content-start p-2 border rounded-lg bg-slate-50">
                    {reportData.mapaArvore.map(item => (
                        <div key={item.name} style={{ width: `${item.size - 20}%`, height: `${item.size - 20}%`, backgroundColor: brandColor, opacity: item.size / 100 }}
                            className="flex items-center justify-center rounded-md text-white text-xs font-bold p-1 leading-tight text-center">
                            {item.name}
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>

        {/* Histórico de Intervenções */}
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2"><Clock /> Histórico de Intervenções</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="relative pl-6">
                    <div className="absolute left-0 top-0 h-full w-0.5 bg-gray-200"></div>
                    {reportData.historicoIntervencoes.map((item, index) => (
                        <div key={index} className="relative mb-6">
                            <div className="absolute -left-[34px] top-1.5 h-4 w-4 rounded-full" style={{backgroundColor: brandColor}}></div>
                            <p className="font-semibold">{item.acao}</p>
                            <p className="text-sm text-muted-foreground">{item.data} - {item.responsavel}</p>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>

        {/* Análise e Recomendações */}
        <div className="space-y-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight text-slate-900 flex items-center justify-center gap-3">
              <ListChecks style={{ color: brandColor }} /> Análise e Recomendações
            </h2>
            <p className="text-muted-foreground mt-2 max-w-2xl mx-auto">
              Análise detalhada por setor, destacando pontos fortes e áreas que necessitam de atenção para um ambiente de trabalho mais saudável e produtivo.
            </p>
          </div>

          {analisePorSetor.map(setorData => (
            <Card key={setorData.setor} className="overflow-hidden border-2 border-slate-100 shadow-lg">
              <CardHeader className="bg-slate-100/80 border-b">
                <CardTitle className="text-2xl text-slate-800">{setorData.setor}</CardTitle>
              </CardHeader>
              <CardContent className="p-6 grid grid-cols-1 lg:grid-cols-2 gap-x-8 gap-y-6">
                {/* Coluna de Pontos Fortes */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2 text-green-600 border-b pb-2">
                    <CheckCircle2 className="text-green-600" /> Pontos Fortes
                  </h3>
                  {setorData.pontosFortes.length > 0 ? (
                    setorData.pontosFortes.map(cat => {
                      const classification = getScoreClassification(cat.pontuacao);
                      return (
                        <div key={`forte-${cat.nome}`} className="text-sm p-3 rounded-lg bg-slate-50 border border-slate-200">
                          <p className="font-semibold text-slate-800">{cat.nome} <span className="font-normal text-slate-500">- {cat.pontuacao}% ({classification})</span></p>
                          <p className="text-muted-foreground mt-1">{getCategoryAnalysisText(cat.nome, cat.pontuacao)}</p>
                        </div>
                      );
                    })
                  ) : <p className="text-sm text-muted-foreground p-3 bg-slate-50 rounded-lg">Nenhum ponto forte destacado para este setor.</p>}
                </div>

                {/* Coluna de Áreas de Melhoria */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2 text-orange-500 border-b pb-2">
                    <AlertTriangle className="text-orange-500" /> Áreas de Melhoria
                  </h3>
                  {setorData.areasDeMelhoria.length > 0 ? (
                    setorData.areasDeMelhoria.map(cat => {
                      const classification = getScoreClassification(cat.pontuacao);
                      const recomendacao = recomendacoesDinamicas.find(r => r.area === cat.nome);
                      return (
                        <Card key={`melhoria-${cat.nome}`} className="bg-orange-50/60 border-orange-200 shadow-sm hover:shadow-md transition-shadow">
                          <CardHeader className="pb-3">
                            <CardTitle className="text-base">{cat.nome} - {cat.pontuacao}%</CardTitle>
                            <CardDescription>
                              <Badge className={`${getClassificationBadgeColor(classification)} text-white shadow-sm`}>{classification}</Badge>
                            </CardDescription>
                          </CardHeader>
                          <CardContent className="text-sm space-y-3">
                            <p className="text-muted-foreground">{getCategoryAnalysisText(cat.nome, cat.pontuacao)}</p>
                            <p className="font-semibold text-orange-800">Risco Identificado: <span className="font-normal">{getRiskIdentificationText(cat.nome)}</span></p>
                            {recomendacao && (
                              <>
                                <Separator className="my-2" />
                                <div className="space-y-1">
                                  <h4 className="font-semibold flex items-center gap-1.5" style={{ color: brandColor }}>
                                    <Lightbulb size={16} style={{ color: brandColor }} /> Plano de Ação Sugerido
                                  </h4>
                                  <p className="text-muted-foreground pl-5">{recomendacao.plano}</p>
                                </div>
                              </>
                            )}
                          </CardContent>
                        </Card>
                      );
                    })
                  ) : <p className="text-sm text-muted-foreground p-3 bg-slate-50 rounded-lg">Nenhuma área de melhoria crítica identificada para este setor.</p>}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Recomendações / Planos de Ação - A tabela foi integrada na seção de análise acima, esta pode ser removida ou mantida como um resumo. Por ora, será removida para evitar redundância. */}
        {/* <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2"><CheckCircle /> Planos de Ação Recomendados</CardTitle>
            </CardHeader>
            <CardContent>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Área de Melhoria</TableHead>
                            <TableHead>Nível de Risco</TableHead>
                            <TableHead>Plano de Ação Sugerido</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {recomendacoesDinamicas.map(rec => (
                            <TableRow key={rec.area}>
                                <TableCell className="font-medium">{rec.area}</TableCell>
                                <TableCell><Badge className={`${getRiskColor(rec.risco)} text-white`}>{rec.risco}</Badge></TableCell>
                                <TableCell>{rec.plano}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card> */}

        {/* Conclusão */}
        <Card>
            <CardHeader>
                <CardTitle>Conclusão do Relatório</CardTitle>
            </CardHeader>
            <CardContent>
                <p className="text-sm text-muted-foreground">Os dados aqui apresentados refletem a coleta realizada na data da avaliação e servem como um instrumento para a tomada de decisão estratégica. Recomenda-se o acompanhamento dos planos de ação propostos e a realização de uma nova medição na data sugerida para avaliar a evolução dos indicadores.</p>
            </CardContent>
            <CardFooter className="flex justify-between">
                <div className="text-sm">
                    <strong>Data do Relatório:</strong> {new Date().toLocaleDateString('pt-BR')}
                </div>
                <div className="w-1/3">
                    <Separator className="mb-2"/>
                    <p className="text-center text-sm font-semibold">Responsável pela Análise</p>
                </div>
            </CardFooter>
        </Card>

      </div>
    </section>
  );
};

export default RelatorioPage;