"use client";

import { useEffect, useState } from "react";
import StatCard from "./StatCard";
import { DateRange } from "react-day-picker";
import { dashboardService } from "@/lib/services/dashboardService";
import { useAuth } from "@/context/AuthContext";
import { Loader2 } from "lucide-react";

interface StatCardGridsProps {
  dateRange?: DateRange;
}

const StatCardGrids = ({ dateRange }: StatCardGridsProps) => {
  const { user } = useAuth();
  const [metrics, setMetrics] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (user) {
      carregarMetricas();
    }
  }, [user, dateRange]); // Recarrega se o usuário mudar ou se alterar as datas no filtro!

  const carregarMetricas = async () => {
    try {
      setIsLoading(true);
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;

      // Tentamos buscar getMetrics (como pede o backlog).
      // Se o Dev A tiver chamado só de getData, usamos ele como fallback.
      const func =
        dashboardService.getMetrics || (dashboardService as any).getData;
      const data = await func(empresaId);

      setMetrics(data);
    } catch (error) {
      console.error("Erro ao carregar métricas do dashboard:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // 4. Loading state exigido no backlog!
  if (isLoading) {
    return (
      <div className="flex w-full justify-center py-8">
        <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <StatCard
        title="Total de Respostas"
        value={metrics?.totalRespostas?.toString() || "0"}
        iconName="clipboardCheck"
        change=""
      />
      <StatCard
        title="Taxa de Participação"
        value={`${metrics?.taxaParticipacao || 0}%`}
        iconName="users"
        change=""
      />
      <StatCard
        title="Engajamento Médio"
        value={`${metrics?.engajamento || 0}%`} // Fallback pra 0 se não vier
        iconName="smile"
        change=""
      />
      <StatCard
        title="NPS Geral"
        value={metrics?.nps?.toString() || "0"}
        iconName="star"
        change=""
      />
    </div>
  );
};

export default StatCardGrids;
