import StatCard from "./StatCard";
import { DateRange } from "react-day-picker";

interface StatCardGridsProps {
  dateRange?: DateRange;
}

const StatCardGrids = ({ dateRange }: StatCardGridsProps) => {
  // Aqui você pode usar dateRange para buscar dados filtrados
  // Por enquanto, os dados são mockados
  console.log("Date range for StatCardGrids:", dateRange);

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <StatCard
        title="Total de Respostas"
        value="1,234"
        iconName="clipboardCheck"
      />
      <StatCard title="Participantes Ativos" value="87" iconName="users" />
      <StatCard title="Engajamento Médio" value="30%" iconName="smile" />
      <StatCard title="NPS Geral" value="+42" iconName="star" />
    </div>
  );
};

export default StatCardGrids;

