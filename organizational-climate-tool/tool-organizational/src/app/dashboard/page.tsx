import StatCardGrids from "@/components/dashboard/StatCardGrids";
import { EngagementChart } from "@/components/dashboard/EngagementChart";
import { DataTable } from "@/components/dashboard/DataTable";

const DashboardPage = () => {
  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Visão geral da sua organização.
      </p>

      <StatCardGrids />

      <div className="flex flex-col gap-6 mt-6">
        <div className="lg:col-span-4">
          <EngagementChart />
        </div>
        <div className="lg:col-span-3">
          <div className="bg-background rounded-lg border p-4 h-full">
            <DataTable />
          </div>
        </div>
      </div>
    </section>
  );
};

export default DashboardPage;
