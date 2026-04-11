"use client";
import React, { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {  ListFilter, Filter, ArrowDownToLine } from "lucide-react";
import {
  ResultsDataTable,
  allMockResults,
  SurveyResult,
} from "@/components/dashboard/ResultsDataTable";
import { ExportReportButton } from "@/components/ui/export-report-button";

import { dadosPesquisas, Pesquisa } from "@/components/dashboard/DataTable";

const departamentos = [
  "Tecnologia",
  "Recursos Humanos",
  "Marketing",
  "Vendas",
  "Financeiro",
];

const ResultadosPage = () => {
  const [selectedSurvey, setSelectedSurvey] = useState<string>("");
  const [selectedDepartment, setSelectedDepartment] = useState<string>("todos");
  const [displayResults, setDisplayResults] = useState<SurveyResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // Função centralizada que busca e filtra os dados
  const applyFilters = (surveyId: string, department: string) => {
    if (!surveyId) {
      setDisplayResults([]);
      return;
    }

    setIsLoading(true);
    // Simula uma chamada de API para buscar os resultados
    setTimeout(() => {
      let results = allMockResults[surveyId] || [];

      // Se um departamento específico for selecionado (diferente de "todos"),
      // simulamos a filtragem alterando os dados para dar um feedback visual.
      if (department !== "todos") {
        // Em um cenário real, você faria uma nova busca na API com o filtro
        // de departamento ou filtraria um conjunto de dados que já contém essa informação.
        results = results.map((result) => {
          // Usa o nome do departamento para criar um "hash" simples e previsível
          // para modificar os dados de forma consistente na simulação.
          const hash = department
            .split("")
            .reduce((acc, char) => acc + char.charCodeAt(0), 0);

          // Modifica a pontuação e o número de respostas de forma simulada
          const scoreVariation = 1 + ((hash % 10) - 4.5) / 50; // Variação de até ~ +/- 9%
          const responseCountVariation = 0.2 + ((hash % 10) / 15); // Respostas entre 20% e ~86% do total

          const newScore = result.averageScore * scoreVariation;
          const newResponseCount = Math.round(
            result.responseCount * responseCountVariation,
          );

          return {
            ...result,
            averageScore: parseFloat(
              Math.max(1, Math.min(5, newScore)).toFixed(1),
            ),
            responseCount: newResponseCount,
          };
        });
      }

      setDisplayResults(results);
      setIsLoading(false);
    }, 300); // Simula um delay de rede
  };

  const handleSurveyChange = (surveyId: string) => {
    setSelectedSurvey(surveyId);
    applyFilters(surveyId, selectedDepartment);
  };

  const handleDepartmentChange = (department: string) => {
    setSelectedDepartment(department);
    // Aplica o filtro apenas se uma pesquisa já estiver selecionada
    if (selectedSurvey) {
      applyFilters(selectedSurvey, department);
    }
  };

  const surveyForReport = dadosPesquisas.find((p) => p.id === selectedSurvey);

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
        <ExportReportButton surveyId={selectedSurvey} surveyName={surveyForReport?.title ?? "Resultados Detalhados"} />
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ListFilter className="h-5 w-5" />
            Filtros de Análise
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4 md:flex-row">
            <Select onValueChange={handleSurveyChange} value={selectedSurvey}>
              <SelectTrigger>
                <SelectValue placeholder="Selecione a pesquisa" />
              </SelectTrigger>
              <SelectContent>
                {dadosPesquisas.map((pesquisa: Pesquisa) => (
                  <SelectItem key={pesquisa.id} value={pesquisa.id}>
                    {pesquisa.title}
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
                {departamentos.map((depto) => (
                  <SelectItem key={depto} value={depto.toLowerCase()}>
                    {depto}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      <ResultsDataTable data={displayResults} isLoading={isLoading} />
    </section>
  );
};

export default ResultadosPage;
