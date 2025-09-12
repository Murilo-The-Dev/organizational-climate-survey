// src/components/dashboard/StatCardsGrid.tsx

"use client"; // Marcamos este container como um Client Component

import StatCard from "./StatCard";

// Documentação:
// Este componente serve como um "container" no lado do cliente para
// agrupar e renderizar todos os StatCards. Isso resolve o problema de um
// Server Component (a página) renderizar diretamente múltiplos Client Components (os cards).

const StatCardGrids = () => {
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
