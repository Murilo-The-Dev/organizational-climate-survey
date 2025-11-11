"use client";

import { DataTable } from "@/components/dashboard/DataTable";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { MoreHorizontal } from "lucide-react";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { ColumnDef } from "@tanstack/react-table";
import { StatusBadge } from "@/components/ui/status-badge";
import { SurveyLinkModal } from "@/components/modals/SurveyLinkModal";
import { useState } from "react";

export type Survey = {
  id: string;
  title: string;
  status: "completed" | "in_progress" | "draft";
  participants: number;
  creationDate: string;
};

const mockSurveys: Survey[] = [
  { id: "SURV-001", title: "Employee Engagement Q1", status: "completed", participants: 152, creationDate: "2025-03-28" },
  { id: "SURV-002", title: "Leadership Feedback H1", status: "completed", participants: 140, creationDate: "2025-06-15" },
  { id: "SURV-003", title: "Annual Satisfaction 2024", status: "completed", participants: 180, creationDate: "2024-12-20" },
  { id: "SURV-004", title: "Organizational Climate H2", status: "in_progress", participants: 125, creationDate: "2025-09-01" },
  { id: "SURV-005", title: "New Hire Onboarding", status: "in_progress", participants: 25, creationDate: "2025-09-10" },
  { id: "SURV-006", title: "Benefits Evaluation", status: "draft", participants: 0, creationDate: "2025-09-18" },
];

export default function PesquisasPage() {
  const [isLinkModalOpen, setIsLinkModalOpen] = useState(false);
  const [selectedSurveyId, setSelectedSurveyId] = useState("");

  const handleGenerateLink = (surveyId: string) => {
    setSelectedSurveyId(surveyId);
    setIsLinkModalOpen(true);
  };

  const columns: ColumnDef<Survey>[] = [
    { accessorKey: "title", header: "Título da Pesquisa" },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status");
        let badgeStatus: "success" | "warning" | "info" | "default" = "default";
        let badgeText = "";

        switch (status) {
          case "completed":
            badgeStatus = "success";
            badgeText = "Concluído";
            break;
          case "in_progress":
            badgeStatus = "warning";
            badgeText = "Em Andamento";
            break;
          case "draft":
            badgeStatus = "info";
            badgeText = "Rascunho";
            break;
          default:
            badgeText = String(status);
        }
        return <StatusBadge status={badgeStatus}>{badgeText}</StatusBadge>;
      },
    },
    { accessorKey: "participants", header: "Participantes" },
    { accessorKey: "creationDate", header: "Data de Criação" },
    {
      id: "actions",
      header: "Ações",
      cell: ({ row }) => {
        const survey = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Abrir menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Ações</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => handleGenerateLink(survey.id)}>
                Gerar Link
              </DropdownMenuItem>
              <DropdownMenuItem>Ver Detalhes</DropdownMenuItem>
              <DropdownMenuItem>Editar Pesquisa</DropdownMenuItem>
              <DropdownMenuItem className="text-red-600">Excluir</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Pesquisas
        </h1>
        <Link href="/pesquisas/nova">
          <Button>Criar Nova Pesquisa</Button>
        </Link>
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie todas as pesquisas criadas no sistema.
      </p>

      <div className="bg-background rounded-lg border p-4 h-full">
        <DataTable columns={columns} data={mockSurveys} />
      </div>

      {selectedSurveyId && (
        <SurveyLinkModal
          isOpen={isLinkModalOpen}
          onClose={() => setIsLinkModalOpen(false)}
          surveyId={selectedSurveyId}
        />
      )}
    </section>
  );
}

