// src/components/modals/ReportFilterModal.tsx
"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { X, Filter, Send } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DateRangePicker } from "@/components/ui/date-range-picker"; // Componente para seleção de intervalo de datas
import { DateRange } from "react-day-picker";

interface ReportFilterModalProps {
  isOpen: boolean;
  onClose: () => void;
  surveyId: string;
  surveyName: string;
}

interface ReportFilters {
  dateRange: DateRange | undefined;
  company: string;
  sector: string;
  category: string; // Novo filtro
  question: string;
  comparisonSurveyId: string;
}

export function ReportFilterModal({ isOpen, onClose, surveyId, surveyName }: ReportFilterModalProps) {
  const router = useRouter();
  const [filters, setFilters] = useState<ReportFilters>({
    dateRange: undefined,
    company: "Todas",
    sector: "Todos",
    category: "Todas", // Novo filtro
    question: "Todas",
    comparisonSurveyId: "Nenhuma",
  });

  const handleGenerateReport = () => {
    // Constrói a string de query params, convertendo as datas para string no formato ISO
    const params = {
      ...filters,
      surveyId,
      startDate: filters.dateRange?.from ? filters.dateRange.from.toISOString().split('T')[0] : '',
      endDate: filters.dateRange?.to ? filters.dateRange.to.toISOString().split('T')[0] : '',
    };

    // Remove chaves com valor "Nenhuma" ou "Todas" para query params mais limpos
    const cleanParams = Object.fromEntries(
      Object.entries(params).filter(([_, v]) => v !== "Nenhuma" && v !== "Todas" && v !== "")
    );

    const queryParams = new URLSearchParams(cleanParams as Record<string, string>).toString();

    // Navega para a nova página de relatório interativo
    router.push(`/relatorios/${surveyId}?${queryParams}`);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[90vw] h-[90vh] p-0 flex flex-col">
        <DialogHeader className="p-6 border-b flex flex-row items-center justify-between">
          <div className="flex items-center gap-3">
            <Filter className="h-6 w-6 text-primary" />
            <DialogTitle className="text-2xl font-bold">
              Filtros para Relatório de Pesquisa: {surveyName}
            </DialogTitle>
          </div>
          <Button variant="ghost" size="icon" onClick={onClose}>
            <X className="h-6 w-6" />
          </Button>
        </DialogHeader>

        <div className="flex-grow overflow-y-auto p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Filtro de Intervalo de Datas */}
            <div className="space-y-2 col-span-1 md:col-span-2 lg:col-span-1">
              <Label htmlFor="dateRange">Intervalo de Datas</Label>
              <DateRangePicker
                date={filters.dateRange}
                onSelect={(dateRange) => setFilters(prev => ({ ...prev, dateRange }))}
              />
            </div>

            {/* Filtro de Empresa */}
            <div className="space-y-2">
              <Label htmlFor="company">Empresa</Label>
              <Select onValueChange={(value) => setFilters(prev => ({ ...prev, company: value }))} defaultValue={filters.company}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a Empresa" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Todas">Todas as Empresas</SelectItem>
                  <SelectItem value="Empresa A">Empresa A</SelectItem>
                  <SelectItem value="Empresa B">Empresa B</SelectItem>
                </SelectContent>
              </Select>
            </div>

	            {/* Filtro de Setor */}
	            <div className="space-y-2">
	              <Label htmlFor="sector">Setor</Label>
	              <Select onValueChange={(value) => setFilters(prev => ({ ...prev, sector: value }))} defaultValue={filters.sector}>
	                <SelectTrigger>
	                  <SelectValue placeholder="Selecione o Setor" />
	                </SelectTrigger>
	                <SelectContent>
	                  <SelectItem value="Todos">Todos os Setores</SelectItem>
	                  <SelectItem value="TI">TI</SelectItem>
	                  <SelectItem value="RH">RH</SelectItem>
	                </SelectContent>
	              </Select>
	            </div>
	
	            {/* Filtro de Categoria */}
	            <div className="space-y-2">
	              <Label htmlFor="category">Categoria</Label>
	              <Select onValueChange={(value) => setFilters(prev => ({ ...prev, category: value }))} defaultValue={filters.category}>
	                <SelectTrigger>
	                  <SelectValue placeholder="Selecione a Categoria" />
	                </SelectTrigger>
	                <SelectContent>
	                  <SelectItem value="Todas">Todas as Categorias</SelectItem>
	                  <SelectItem value="Lideranca">Liderança</SelectItem>
	                  <SelectItem value="Comunicacao">Comunicação</SelectItem>
	                  <SelectItem value="Reconhecimento">Reconhecimento</SelectItem>
	                </SelectContent>
	              </Select>
	            </div>
	
	            {/* Filtro de Pergunta */}
            <div className="space-y-2">
              <Label htmlFor="question">Pergunta</Label>
              <Select onValueChange={(value) => setFilters(prev => ({ ...prev, question: value }))} defaultValue={filters.question}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a Pergunta" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Todas">Todas as Perguntas</SelectItem>
                  <SelectItem value="Q1">Pergunta 1</SelectItem>
                  <SelectItem value="Q2">Pergunta 2</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Comparação entre Pesquisas */}
            <div className="space-y-2">
              <Label htmlFor="comparison">Comparar com</Label>
              <Select onValueChange={(value) => setFilters(prev => ({ ...prev, comparisonSurveyId: value }))} defaultValue={filters.comparisonSurveyId}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione outra Pesquisa" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Nenhuma">Nenhuma Comparação</SelectItem>
                  <SelectItem value="survey-xyz">Pesquisa Anterior (XYZ)</SelectItem>
                  <SelectItem value="survey-abc">Pesquisa Global (ABC)</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        <div className="p-6 border-t flex justify-end gap-3">
          <Button variant="outline" onClick={onClose}>
            Cancelar
          </Button>
          <Button onClick={handleGenerateReport}>
            <Send className="h-4 w-4 mr-2" /> Gerar Relatório
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}