"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { PlusCircle, Search } from "lucide-react";
import { SurveyCard } from "@/components/pesquisas/SurveyCard";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { CreateSurveyForm } from "@/components/forms/CreateSurveyForm";
import { SurveyDetailsModal } from "@/components/modals/SurveyDetailsModal";
import { SurveyLinkModal } from "@/components/modals/SurveyLinkModal";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuth } from "@/context/AuthContext";
import { pesquisaService } from "@/lib/services/pesquisaService";
import type { Pesquisa, StatusPesquisa } from "@/lib/types";

const PesquisasPage = () => {
  const { user } = useAuth();
  const [pesquisas, setPesquisas] = React.useState<Pesquisa[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [selectedSurvey, setSelectedSurvey] = React.useState<Pesquisa | null>(null);
  const [searchQuery, setSearchQuery] = React.useState("");
  const [statusFilter, setStatusFilter] = React.useState<StatusPesquisa | "todos">("todos");
  const [isLinkModalOpen, setIsLinkModalOpen] = React.useState(false);
  const [selectedSurveyId, setSelectedSurveyId] = React.useState("");

  const fetchPesquisas = React.useCallback(async () => {
    if (!user?.empresa_id) return;
    setIsLoading(true);
    try {
      const data = await pesquisaService.listByEmpresa(user.empresa_id);
      setPesquisas(data);
    } finally {
      setIsLoading(false);
    }
  }, [user?.empresa_id]);

  React.useEffect(() => {
    fetchPesquisas();
  }, [fetchPesquisas]);

  const filtered = pesquisas.filter((p) => {
    const matchesSearch = p.titulo.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus = statusFilter === "todos" || p.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    fetchPesquisas();
  };

  const handleGenerateLink = (surveyId: string) => {
    setSelectedSurveyId(surveyId);
    setIsLinkModalOpen(true);
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Pesquisas</h1>
          <p className="text-muted-foreground mt-2">
            Crie, gerencie e visualize todos os seus formulários.
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-300">
              <PlusCircle className="mr-2 h-4 w-4" />
              Criar Pesquisa
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-5xl">
            <DialogHeader>
              <DialogTitle>Criar Nova Pesquisa</DialogTitle>
              <DialogDescription>
                Preencha as informações abaixo para criar um novo formulário.
              </DialogDescription>
            </DialogHeader>
            <CreateSurveyForm onClose={handleCloseDialog} />
          </DialogContent>
        </Dialog>
      </div>

      <div className="flex items-center gap-4 mb-8">
        <div className="relative w-full md:w-1/3">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Buscar por nome..."
            className="pl-8"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Select
          value={statusFilter}
          onValueChange={(v) => setStatusFilter(v as StatusPesquisa | "todos")}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Status: Todos" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="todos">Todos</SelectItem>
            <SelectItem value="Rascunho">Rascunhos</SelectItem>
            <SelectItem value="Ativa">Ativas</SelectItem>
            <SelectItem value="Concluída">Concluídas</SelectItem>
            <SelectItem value="Arquivada">Arquivadas</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isLoading ? (
          Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-48 w-full rounded-xl" />
          ))
        ) : filtered.length > 0 ? (
          filtered.map((pesquisa) => (
            <SurveyCard
              key={pesquisa.id_pesquisa}
              id={String(pesquisa.id_pesquisa)}
              title={pesquisa.titulo}
              description={pesquisa.descricao}
              tag={pesquisa.status}
              creationDate={new Date(pesquisa.data_criacao).toLocaleDateString("pt-BR")}
              onViewDetails={() => setSelectedSurvey(pesquisa)}
              onGenerateLink={handleGenerateLink}
            />
          ))
        ) : (
          <p className="col-span-3 text-center text-muted-foreground py-10">
            {searchQuery
              ? `Nenhuma pesquisa encontrada com o termo "${searchQuery}".`
              : "Nenhuma pesquisa encontrada."}
          </p>
        )}
      </div>

      <Dialog
        open={!!selectedSurvey}
        onOpenChange={(isOpen) => !isOpen && setSelectedSurvey(null)}
      >
        <DialogContent className="sm:max-w-4xl h-[90vh] flex flex-col">
          {selectedSurvey && <SurveyDetailsModal survey={selectedSurvey as any} />}
        </DialogContent>
      </Dialog>

      {selectedSurveyId && (
        <SurveyLinkModal
          isOpen={isLinkModalOpen}
          onClose={() => setIsLinkModalOpen(false)}
          surveyId={selectedSurveyId}
        />
      )}
    </section>
  );
};

export default PesquisasPage;


type Survey = Pesquisa & {
  id?: string;
  title: string;
  description: string;
  status: string;
  tag: string;
  creationDate: string;
  onViewDetails: () => void;
};

const mockSurveys = [
  {
    id: "SURV-001",
    title: "Engajamento Trimestral Q3",
    description:
      "Pesquisa para medir o nível de engajamento e satisfação dos colaboradores neste trimestre.",
    tag: "Engajamento",
    creationDate: "15/08/2025",
  },
  {
    id: "SURV-002",
    title: "Feedback de Liderança H2",
    description:
      "Avaliação 360º dos líderes e gestores da organização para o segundo semestre.",
    tag: "Liderança",
    creationDate: "01/09/2025",
  },
  {
    id: "SURV-003",
    title: "Pesquisa de Benefícios 2025",
    description:
      "Coleta de feedback sobre o pacote de benefícios atual e sugestões de melhorias.",
    tag: "RH",
    creationDate: "18/09/2025",
  },
  {
    id: "SURV-004",
    title: "Pesquisa de Cultura Organizacional 2025",
    description:
      "Coleta de feedback sobre a cultura organizacional atual e sugestões de melhorias.",
    tag: "Cultura Organizacional",
    creationDate: "22/10/2025",
  },
  {
    id: "SURV-005",
    title: "Pesquisa de Satisfação do Colaborador 2025",
    description:
      "Coleta de feedback sobre a satisfação do colaborador atual e sugestões de melhorias.",
    tag: "Satisfação do Colaborador",
    creationDate: "22/11/2025",
  },
  {
    id: "SURV-006",
    title: "Pesquisa de Satisfação do Cliente 2025",
    description:
      "Coleta de feedback sobre a satisfação do cliente atual e sugestões de melhorias.",
    tag: "Satisfação do Cliente",
    creationDate: "22/11/2025",
  },
  {
    id: "SURV-007",
    title: "Análise de Trabalhadores 2025",
    description:
      "Coleta de feedback sobre a satisfação do trabalhador atual e sugestões de melhorias.",
    tag: "Análise de Trabalhadores",
    creationDate: "22/11/2025",
  },
  {
    id: "SURV-008",
    title: "Pesquisa de Satisfação do Trabalhador 2025",
    description:
      "Coleta de feedback sobre a satisfação do trabalhador atual e sugestões de melhorias.",
    tag: "Satisfação do Trabalhador",
    creationDate: "22/12/2025",
  },
];

const PesquisasPage = () => {
  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [selectedSurvey, setSelectedSurvey] = React.useState<Survey | null>(null);
  const [searchQuery, setSearchQuery] = React.useState("");

  // Estados para o modal de QR Code
  const [isLinkModalOpen, setIsLinkModalOpen] = React.useState(false);
  const [selectedSurveyId, setSelectedSurveyId] = React.useState("");

  const filteredSurveys = mockSurveys.filter((survey) =>
    survey.title.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
  };

  // Função para gerar link/QR Code
  const handleGenerateLink = (surveyId: string) => {
    setSelectedSurveyId(surveyId);
    setIsLinkModalOpen(true);
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Pesquisas</h1>
          <p className="text-muted-foreground mt-2">
            Crie, gerencie e visualize todos os seus formulários.
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-300">
              <PlusCircle className="mr-2 h-4 w-4" />
              Criar Pesquisa
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-5xl">
            <DialogHeader>
              <DialogTitle>Criar Nova Pesquisa</DialogTitle>
              <DialogDescription>
                Preencha as informações abaixo para criar um novo formulário.
              </DialogDescription>
            </DialogHeader>
            <CreateSurveyForm onClose={handleCloseDialog} />
          </DialogContent>
        </Dialog>
      </div>

      <div className="flex items-center gap-4 mb-8">
        <div className="relative w-full md:w-1/3">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Buscar por nome..."
            className="pl-8"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Select>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Status: Todos" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="todos">Todos</SelectItem>
            <SelectItem value="rascunho">Rascunhos</SelectItem>
            <SelectItem value="ativo">Ativas</SelectItem>
            <SelectItem value="concluido">Concluídas</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredSurveys.length > 0 ? (
          filteredSurveys.map((survey: any) => (
            <SurveyCard
              key={survey.id || survey.title}
              id={survey.id || survey.title}
              title={survey.title}
              description={survey.description}
              tag={survey.tag}
              creationDate={survey.creationDate}
              onViewDetails={() => setSelectedSurvey(survey)}
              onGenerateLink={handleGenerateLink}
            />
          ))
        ) : (
          <p className="col-span-3 text-center text-muted-foreground py-10">
            Nenhuma pesquisa encontrada com o termo "{searchQuery}".
          </p>
        )}
      </div>

      {/* Modal de detalhes da pesquisa */}
      <Dialog
        open={!!selectedSurvey}
        onOpenChange={(isOpen) => !isOpen && setSelectedSurvey(null)}
      >
        <DialogContent className="sm:max-w-4xl h-[90vh] flex flex-col">
          {selectedSurvey && <SurveyDetailsModal survey={selectedSurvey} />}
        </DialogContent>
      </Dialog>

      {/* Modal de QR Code */}
      {selectedSurveyId && (
        <SurveyLinkModal
          isOpen={isLinkModalOpen}
          onClose={() => setIsLinkModalOpen(false)}
          surveyId={selectedSurveyId}
        />
      )}
    </section>
  );
};

export default PesquisasPage;