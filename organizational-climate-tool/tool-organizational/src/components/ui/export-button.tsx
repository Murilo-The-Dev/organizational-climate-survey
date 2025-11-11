import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Download } from "lucide-react";
import { toast } from "sonner";

interface ExportButtonProps {
  data: any[];
  filename?: string;
  headers?: string[];
}

export function ExportButton({ data, filename = "export", headers }: ExportButtonProps) {
  const convertToCsv = (data: any[], headers?: string[]) => {
    if (!data || data.length === 0) return "";

    const csvRows = [];
    const actualHeaders = headers || Object.keys(data[0]);
    csvRows.push(actualHeaders.join(","));

    for (const row of data) {
      const values = actualHeaders.map(header => {
        const escaped = ('' + row[header]).replace(/"/g, '""');
        return `"${escaped}"`;
      });
      csvRows.push(values.join(","));
    }
    return csvRows.join("\n");
  };

  const handleExportCsv = () => {
    if (data.length === 0) {
      toast.info("Não há dados para exportar.");
      return;
    }
    const csvString = convertToCsv(data, headers);
    const blob = new Blob([csvString], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.setAttribute("download", `${filename}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    toast.success("Dados exportados com sucesso para CSV!");
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="h-8 gap-1">
          <Download className="h-3.5 w-3.5" />
          <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
            Exportar
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={handleExportCsv}>
          Exportar para CSV
        </DropdownMenuItem>
        {/* Futuras opções de exportação, como PDF, podem ser adicionadas aqui */}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

