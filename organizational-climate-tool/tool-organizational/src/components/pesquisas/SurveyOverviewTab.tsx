"use client";

import StatCard from "@/components/dashboard/StatCard";
import { NpsBreakdownChart } from "@/components/dashboard/charts/NpsBreakdownChart";
import { ScoreDistributionChart } from "@/components/dashboard/charts/ScoreDistributionChart";
import { ChartLineDefault } from "@/components/dashboard/charts/LineChart";
import { ChartStacked } from "@/components/dashboard/charts/StackedChart";

export const SurveyOverviewTab = () => {
  return (
    <div className="flex flex-col gap-6 overflow-y-auto pr-4">
      <div className="grid gap-4 md:grid-cols-3">
        <StatCard
          title="Taxa de Participação"
          value="87%"
          iconName="users"
          change="152 de 175 convidados"
        />
        <StatCard
          title="Média Geral"
          value="4.1 / 5.0"
          iconName="smile"
          change="+0.2 vs. última pesquisa"
        />
        <StatCard
          title="e-NPS"
          value="+45"
          iconName="star"
          change="Excelente"
        />
      </div>

      <div className="grid grid-cols-2 md:grid-cols-2 gap-6">
        <ScoreDistributionChart />
        <NpsBreakdownChart />
      </div>
      <div className="grid grid-cols-2 md:grid-cols-2 gap-6">
        <ChartLineDefault />
        <ChartStacked />
      </div>
    </div>
  );
};
