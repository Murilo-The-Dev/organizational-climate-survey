
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

const mockSurveys = [
  {
    title: "Engajamento Trimestral Q3",
    description:
      "Pesquisa para medir o nível de engajamento e satisfação dos colaboradores neste trimestre.",
    tag: "Engajamento",
    creationDate: "15/08/2025",
  },
  {
    title: "Feedback de Liderança H2",
    description:
      "Avaliação 360º dos líderes e gestores da organização para o segundo semestre.",
    tag: "Liderança",
    creationDate: "01/09/2025",
  },
  {
    title: "Pesquisa de Benefícios 2025",
    description:
      "Coleta de feedback sobre o pacote de benefícios atual e sugestões de melhorias.",
    tag: "RH",
    creationDate: "18/09/2025",
  },
  {
    title: "Pesquisa de Cultura Organizacional 2025",
    description:
      "Coleta de feedback sobre a cultura organizacional atual e sugestões de melhorias.",
    tag: "Cultura Organizacional",
    creationDate: "22/10/2025",
  },
  {
    title: "Pesquisa de Satisfação do Colaborador 2025",
    description:
      "Coleta de feedback sobre a satisfação do colaborador atual e sugestões de melhorias.",
    tag: "Satisfação do Colaborador",
    creationDate: "22/11/2025",
  },
  {
    title: "Pesquisa de Satisfação do Cliente 2025",
    description:
      "Coleta de feedback sobre a satisfação do cliente atual e sugestões de melhorias.",
    tag: "Satisfação do Cliente",
    creationDate: "22/11/2025",
  },
  {
    title: "Análise de Trabalhadores 2025",
    description:
      "Coleta de feedback sobre a satisfação do trabalhador atual e sugestões de melhorias.",
    tag: "Análise de Trabalhadores",
    creationDate: "22/11/2025",
  },
  {
    title: "Pesquisa de Satisfação do Trabalhador 2025",
    description:
      "Coleta de feedback sobre a satisfação do trabalhador atual e sugestões de melhorias.",
    tag: "Satisfação do Trabalhador",
    creationDate: "22/12/2025",
  },
  
];

const PesquisasPage = () => {
  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Pesquisas</h1>
          <p className="text-muted-foreground mt-2">
            Crie, gerencie e visualize todos os seus formulários.
          </p>
        </div>
        <Button className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-300">
          <PlusCircle className="mr-2 h-4 w-4" />
          Criar Pesquisa
        </Button>
      </div>

      <div className="flex items-center gap-4 mb-8">
        <div className="relative w-full md:w-1/3">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Buscar por nome..." className="pl-8" />
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
        {mockSurveys.map((survey) => (
          <SurveyCard
            key={survey.title}
            title={survey.title}
            description={survey.description}
            tag={survey.tag}
            creationDate={survey.creationDate}
          />
        ))}
      </div>
    </section>
  );
};

export default PesquisasPage;
