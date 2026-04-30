"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Download, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";

interface ExportButtonProps {
  surveyId?: string;
  filename?: string;
}

export function ExportButton({
  surveyId,
  filename = "dados",
}: ExportButtonProps) {
  const [isExporting, setIsExporting] = useState(false);
  const { user } = useAuth();

  const handleExportCsv = async () => {
    if (!surveyId) {
      toast.info("Selecione uma pesquisa para exportar os dados.");
      return;
    }

    try {
      setIsExporting(true);
      const token = (user as any)?.token;
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

      const response = await fetch(
        `${baseUrl}/dashboards/${surveyId}/export?format=csv`,
        {
          method: "GET",
          headers: { Authorization: `Bearer ${token}` },
        },
      );

      if (!response.ok) throw new Error("Falha na API");

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `${filename}.csv`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      toast.success("Dados exportados com sucesso!");
    } catch (error) {
      console.error(error);
      toast.error("Erro ao exportar dados do servidor.");
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <Button
      variant="outline"
      size="sm"
      className="h-8 gap-1"
      onClick={handleExportCsv}
      disabled={isExporting || !surveyId}
    >
      {isExporting ? (
        <Loader2 className="h-3.5 w-3.5 animate-spin" />
      ) : (
        <Download className="h-3.5 w-3.5" />
      )}
      <span className="sr-only sm:not-sr-only">Exportar CSV</span>
    </Button>
  );
}
