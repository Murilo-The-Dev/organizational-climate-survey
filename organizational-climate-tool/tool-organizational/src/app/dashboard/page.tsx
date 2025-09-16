import StatCardGrids from "@/components/dashboard/StatCardGrids";
import { EngagementChart } from "@/components/dashboard/charts/EngagementChart";
import { ChartPieLabel } from "@/components/dashboard/charts/PieChart";
import { DataTable } from "@/components/dashboard/DataTable";

const DashboardPage = () => {
  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Visão geral da sua organização.
      </p>

      <StatCardGrids />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
        <div className="lg:col-span-4">
          <EngagementChart />
        </div>
          <div className="">
              <ChartPieLabel />
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
