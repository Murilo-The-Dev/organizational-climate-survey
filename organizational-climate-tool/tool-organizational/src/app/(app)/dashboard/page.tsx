"use client"

import StatCardGrids from "@/components/dashboard/StatCardGrids";
import { ChartBarStacked } from "@/components/dashboard/charts/EngagementChart";
import { ChartRadialShape } from "@/components/dashboard/charts/RadialChart";
import { ChartPieLabel } from "@/components/dashboard/charts/PieChart";
import { ChartBarInteractive } from "@/components/dashboard/charts/BarChartInteractive";
import { ChartLineTrends } from "@/components/dashboard/charts/ChartLineTrends";
import { ChartBarComparative } from "@/components/dashboard/charts/ChartBarComparative";
import { DataTable, Pesquisa, columns, dadosPesquisas } from "@/components/dashboard/DataTable";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { useState } from "react";
import { DateRange } from "react-day-picker";
import { ExternalDataModal } from "@/components/dashboard/ExternalDataModal";

const DashboardPage = () => {
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedSurvey, setSelectedSurvey] = useState<Pesquisa | null>(null);

  const handleOpenModal = (survey: Pesquisa) => {
    setSelectedSurvey(survey);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedSurvey(null);
  };

  const handleSaveData = (surveyId: string, data: { absenteismo: number; turnover: number }) => {
    console.log("Salvando dados para a pesquisa:", surveyId, data);
    // IMPORTANTE: Aqui você faria a chamada para sua API para salvar os dados no banco.
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Dashboard
        </h1>
        <DateRangePicker date={dateRange} onSelect={setDateRange} />
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Visão geral da sua organização.
      </p>

      <StatCardGrids dateRange={dateRange} />

      <div className="mt-6 lg:col-span-3">
        <ChartBarInteractive dateRange={dateRange} />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-6">
        <div>
          <ChartBarStacked />
        </div>
        <div>
          <ChartPieLabel />
        </div>
        <div>
          <ChartRadialShape />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
        <ChartLineTrends dateRange={dateRange} />
        <ChartBarComparative dateRange={dateRange} />
      </div>

      <div className="mt-6">
        <div className="bg-background rounded-lg border p-4 h-full">
          <DataTable 
            columns={columns} 
            data={dadosPesquisas} 
            meta={{
              onOpenExternalDataModal: handleOpenModal,
            }}
          />
        </div>
      </div>

      {selectedSurvey && (
        <ExternalDataModal
          isOpen={isModalOpen}
          onClose={handleCloseModal}
          surveyId={selectedSurvey.id}
          surveyTitle={selectedSurvey.title}
          onSave={handleSaveData}
        />
      )}
    </section>
  );
};

export default DashboardPage;