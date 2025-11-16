// src/components/ui/status-badge.tsx

import { cn } from "@/lib/utils";
import { type Pesquisa } from "../dashboard/DataTable"; // Ajuste o caminho se o tipo 'Pesquisa' estiver em outro lugar

// Mapeia cada status para um texto e uma classe de cor do Tailwind
const statusConfig = {
  completed: {
    label: "Concluído",
    color: "bg-green-500",
  },
  in_progress: {
    label: "Em Andamento",
    color: "bg-yellow-500", 
  },
  draft: {
    label: "Rascunho",
    color: "bg-gray-500",
  },
};

// Define que a prop 'status' deve ser um dos status válidos do nosso tipo 'Pesquisa'
type StatusBadgeProps = {
  status: Pesquisa["status"];
  className?: string;
};

export const StatusBadge = ({ status, className }: StatusBadgeProps) => {
  // Pega a configuração correta ou uma padrão caso o status seja inesperado
  const config = statusConfig[status] || { label: "Desconhecido", color: "bg-gray-400" };

  return (
    <div
      className={cn(
        "inline-flex items-center gap-2 rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        className
      )}
    >
      <span className={cn("h-2 w-2 rounded-full", config.color)} />
      {config.label}
    </div>
  );
};