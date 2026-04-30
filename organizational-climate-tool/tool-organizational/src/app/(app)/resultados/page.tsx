"use client";

import React, { useState, useEffect } from "react";
import { ListFilter } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  ResultsDataTable,
  SurveyResult,
} from "@/components/dashboard/ResultsDataTable";
import { ExportReportButton } from "@/components/ui/export-report-button";
import { pesquisaService } from "@/lib/services/pesquisaService";
import { dashboardService } from "@/lib/services/dashboardService";
import { useAuth } from "@/context/AuthContext";
import { toast } from "sonner";

export default function ResultadosPage() {
  const { user } = useAuth();
  const [pesquisas, setPesquisas] = useState<any[]>([]);
  const [selectedSurvey, setSelectedSurvey] = useState<string>("");
  const [selectedDepartment, setSelectedDepartment] = useState<string>("todos");
  const [displayResults, setDisplayResults] = useState<SurveyResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // 1. Carrega as pesquisas reais da empresa para preencher o Select
  useEffect(() => {
    if (user) {
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;
      pesquisaService
        .listByEmpresa(empresaId)
        .then(setPesquisas)
        .catch(console.error);
    }
  }, [user]);

  // 2. Busca os dados da pesquisa selecionada direto da API
  const applyFilters = async (surveyId: string, department: string) => {
    if (!surveyId) {
      setDisplayResults([]);
      return;
    }

    setIsLoading(true);

    try {
      const func =
        dashboardService.getMetrics || (dashboardService as any).getData;
      const data = await func(Number(surveyId));
      const metricas = data?.metricas_por_pergunta || [];

      // 3. Traduz os dados do backend para o formato que a tabela espera
      let results = metricas.map((m: any, index: number) => ({
        id: String(m.id_pergunta || index),
        question: m.texto_pergunta || "Métrica",
        category: "Geral" as any,
        averageScore: m.media ? Number(m.media.toFixed(1)) : 0,
        responseCount: m.total_respostas || 0,
      }));

      setDisplayResults(results);
    } catch (error) {
      console.error("Erro ao buscar resultados:", error);
      toast.error("Não foi possível carregar os resultados.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSurveyChange = (surveyId: string) => {
    setSelectedSurvey(surveyId);
    applyFilters(surveyId, selectedDepartment);
  };

  const handleDepartmentChange = (department: string) => {
    setSelectedDepartment(department);
    if (selectedSurvey) applyFilters(selectedSurvey, department);
  };

  const surveyForReport = pesquisas.find(
    (p) => String(p.id_pesquisa) === selectedSurvey,
  );

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Resultados Detalhados
          </h1>
          <p className="text-muted-foreground mt-2">
            Filtre e analise as respostas de cada pesquisa em detalhes.
          </p>
        </div>
        <ExportReportButton surveyId={selectedSurvey} />
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ListFilter className="h-5 w-5" /> Filtros de Análise
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4 md:flex-row">
            <Select onValueChange={handleSurveyChange} value={selectedSurvey}>
              <SelectTrigger>
                <SelectValue placeholder="Selecione a pesquisa" />
              </SelectTrigger>
              <SelectContent>
                {pesquisas.map((p: any) => (
                  <SelectItem key={p.id_pesquisa} value={String(p.id_pesquisa)}>
                    {p.titulo}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              onValueChange={handleDepartmentChange}
              value={selectedDepartment}
            >
              <SelectTrigger>
                <SelectValue placeholder="Todos os departamentos" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="todos">Todos os departamentos</SelectItem>
                <SelectItem value="ti">Tecnologia</SelectItem>
                <SelectItem value="rh">Recursos Humanos</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      <ResultsDataTable data={displayResults} isLoading={isLoading} />
    </section>
  );
}
