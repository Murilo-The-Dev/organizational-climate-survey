"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Download,
  Loader2,
  FileSpreadsheet,
  FileText,
  TableProperties,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";

interface ExportReportButtonProps {
  surveyId: string;
}

export function ExportReportButton({ surveyId }: ExportReportButtonProps) {
  const [isExporting, setIsExporting] = useState(false);
  const { user } = useAuth();

  const handleExport = async (format: "csv" | "xlsx" | "pdf") => {
    if (!surveyId) {
      toast.error("Nenhuma pesquisa selecionada para exportação.");
      return;
    }

    try {
      setIsExporting(true);
      toast.info(`Iniciando download em ${format.toUpperCase()}...`);

      // Alexandre: Estou batendo direto no endpoint com fetch em vez do axios do api.ts
      // porque lidar com responseType 'blob' no Axios centralizado costuma dar dor de cabeça.
      const token = (user as any)?.token;
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

      const response = await fetch(
        `${baseUrl}/dashboards/${surveyId}/export?format=${format}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        },
      );

      if (!response.ok) {
        throw new Error(
          `Falha ao exportar relatório: Status ${response.status}`,
        );
      }

      // Transforma a resposta binária em um Blob (arquivo)
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);

      // Cria um link <a href="..."> invisível na tela para forçar o download no navegador
      const a = document.createElement("a");
      a.href = url;
      a.download = `relatorio-pesquisa-${surveyId}.${format}`;
      document.body.appendChild(a);
      a.click();

      // Limpa a memória
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      toast.success("Download concluído com sucesso!");
    } catch (error) {
      console.error("Erro na exportação:", error);
      toast.error(
        "Ocorreu um erro ao gerar o arquivo. Verifique com o administrador.",
      );
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          className="gap-2 border-blue-200 hover:bg-blue-50 hover:text-blue-700 transition-colors"
          disabled={isExporting || !surveyId}
        >
          {isExporting ? (
            <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
          ) : (
            <Download className="h-4 w-4" />
          )}
          {isExporting ? "Gerando Arquivo..." : "Exportar Resultados"}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuItem
          onClick={() => handleExport("csv")}
          className="cursor-pointer gap-2 py-3"
        >
          <TableProperties className="h-4 w-4 text-green-600" />
          <span>Exportar como CSV</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => handleExport("xlsx")}
          className="cursor-pointer gap-2 py-3"
        >
          <FileSpreadsheet className="h-4 w-4 text-emerald-600" />
          <span>Exportar como Excel (.xlsx)</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => handleExport("pdf")}
          className="cursor-pointer gap-2 py-3"
        >
          <FileText className="h-4 w-4 text-red-500" />
          <span>Exportar como PDF</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
