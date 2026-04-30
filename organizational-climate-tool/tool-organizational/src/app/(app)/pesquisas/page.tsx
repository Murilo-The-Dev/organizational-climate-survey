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
import { PlusCircle, Search, Loader2 } from "lucide-react";
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
import { useAuth } from "@/context/AuthContext";
import { pesquisaService } from "@/lib/services/pesquisaService";
import type { Pesquisa, StatusPesquisa } from "@/lib/types";

const PesquisasPage = () => {
  const { user } = useAuth();
  const [pesquisas, setPesquisas] = React.useState<Pesquisa[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [selectedSurvey, setSelectedSurvey] = React.useState<Pesquisa | null>(
    null,
  );
  const [searchQuery, setSearchQuery] = React.useState("");
  const [statusFilter, setStatusFilter] = React.useState<
    StatusPesquisa | "todos"
  >("todos");

  // Estados para o modal de QR Code/Link
  const [isLinkModalOpen, setIsLinkModalOpen] = React.useState(false);
  const [selectedSurveyId, setSelectedSurveyId] = React.useState("");

  const fetchPesquisas = React.useCallback(async () => {
    // Garantimos que pegamos o ID da empresa corretamente
    const empresaId = (user as any)?.empresaId
      ? Number((user as any).empresaId)
      : 1;

    setIsLoading(true);
    try {
      const data = await pesquisaService.listByEmpresa(empresaId);
      setPesquisas(data);
    } catch (error) {
      console.error("Erro ao buscar pesquisas:", error);
    } finally {
      setIsLoading(false);
    }
  }, [user]);

  React.useEffect(() => {
    if (user) fetchPesquisas();
  }, [fetchPesquisas, user]);

  const filtered = pesquisas.filter((p) => {
    const matchesSearch = p.titulo
      .toLowerCase()
      .includes(searchQuery.toLowerCase());
    const matchesStatus = statusFilter === "todos" || p.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    fetchPesquisas(); // Recarrega a lista após criar uma nova
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

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 min-h-[300px]">
        {isLoading ? (
          <div className="col-span-full flex justify-center items-center">
            <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
          </div>
        ) : filtered.length > 0 ? (
          filtered.map((pesquisa) => (
            <SurveyCard
              key={pesquisa.id_pesquisa}
              id={String(pesquisa.id_pesquisa)}
              title={pesquisa.titulo}
              description={pesquisa.descricao}
              tag={pesquisa.status}
              creationDate={new Date(pesquisa.data_criacao).toLocaleDateString(
                "pt-BR",
              )}
              onViewDetails={() => setSelectedSurvey(pesquisa)}
              onGenerateLink={handleGenerateLink}
            />
          ))
        ) : (
          <p className="col-span-full text-center text-muted-foreground py-10">
            {searchQuery
              ? `Nenhuma pesquisa encontrada com o termo "${searchQuery}".`
              : "Nenhuma pesquisa encontrada."}
          </p>
        )}
      </div>

      <DialogContent className="sm:max-w-4xl h-[90vh] flex flex-col">
        {selectedSurvey && (
          <SurveyDetailsModal
            survey={selectedSurvey as any}
            description={selectedSurvey.descricao}
            tag={selectedSurvey.status}
            creationDate={new Date(
              selectedSurvey.data_criacao,
            ).toLocaleDateString("pt-BR")}
          />
        )}
      </DialogContent>

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
