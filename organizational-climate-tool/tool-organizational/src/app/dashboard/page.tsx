import StatCardGrids from "@/components/dashboard/StatCardGrids";
import { ChartBarStacked } from "@/components/dashboard/charts/EngagementChart";
import { ChartRadialShape } from "@/components/dashboard/charts/RadialChart";
import { ChartPieLabel } from "@/components/dashboard/charts/PieChart";
import { ChartBarInteractive } from "@/components/dashboard/charts/BarChartInteractive";
import { DataTable } from "@/components/dashboard/DataTable";

const DashboardPage = () => {
  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Dashboard
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Visão geral da sua organização.
      </p>

      <StatCardGrids />

      <div className="mt-6 lg:col-span-3">
        <ChartBarInteractive />
      </div>

      <div className="mt-6 flex flex-row lg:grid lg:grid-cols-3 gap-4">
        <div className="flex-1">
          <ChartBarStacked />
        </div>
        <div className="flex-1">
          <ChartPieLabel />
        </div>
        <div className="flex-1">
          <ChartRadialShape />
        </div>
      </div>
      <div className="mt-6">
        <div className="bg-background rounded-lg border p-4 h-full">
          <DataTable />
        </div>
      </div>
    </section>
  );
};

export default DashboardPage;
