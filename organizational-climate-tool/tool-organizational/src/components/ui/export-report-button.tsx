"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowDownToLine } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogTrigger,
  DialogClose,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { dadosPesquisas, Pesquisa } from "@/components/dashboard/DataTable";

interface ExportReportButtonProps {
  surveyId: string;
}

export function ExportReportButton({ surveyId: initialSurveyId }: ExportReportButtonProps) {
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [selectedSurvey, setSelectedSurvey] = useState<string>(initialSurveyId || "");

  // Sincroniza o estado interno com a prop quando o modal é aberto ou a prop muda.
  // Isso garante que a pesquisa selecionada na página principal seja pré-selecionada no modal.
  React.useEffect(() => {
    if (isOpen) {
      setSelectedSurvey(initialSurveyId || "");
    }
  }, [isOpen, initialSurveyId]);

  const handleExport = () => {
    // Usa o ID da pesquisa selecionada para construir a URL dinamicamente
    const surveyId = selectedSurvey || "RESULTADOS_GERAIS"; // Fallback para um valor geral
    const url = `/relatorios/${surveyId}`;
    router.push(url);
    setIsOpen(false); // Fecha o modal após a navegação
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button>
          <ArrowDownToLine className="mr-2 h-4 w-4" />
          Exportar Relatório
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Filtros para Relatório de Pesquisa</DialogTitle>
          <DialogDescription>
            Selecione a pesquisa para a qual deseja gerar um relatório detalhado.
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="pesquisa">Pesquisa</Label>
            <Select value={selectedSurvey} onValueChange={setSelectedSurvey}>
              <SelectTrigger id="pesquisa">
                <SelectValue placeholder="Selecione uma pesquisa" />
              </SelectTrigger>
              <SelectContent>
                {dadosPesquisas.map((pesquisa: Pesquisa) => (
                  <SelectItem key={pesquisa.id} value={pesquisa.id}>
                    {pesquisa.title}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Cancelar</Button>
          </DialogClose>
          <Button onClick={handleExport} disabled={!selectedSurvey}>Exportar</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}