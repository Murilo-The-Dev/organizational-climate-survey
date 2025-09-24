import {
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { type Pesquisa } from "@/components/dashboard/DataTable"; // Supondo que seu tipo e dados estão aqui
import { SurveyOverviewTab } from "../pesquisas/SurveyOverviewTab";

type SurveyDetailsModalProps = {
  survey: Pesquisa;
};

export const SurveyDetailsModal = ({ survey }: SurveyDetailsModalProps) => {
  return (
    <>
      <DialogHeader>
        <DialogTitle>{survey.titulo}</DialogTitle>
        <DialogDescription>
          Análise detalhada e gerenciamento da pesquisa.
        </DialogDescription>
      </DialogHeader>

      <Tabs
        defaultValue="overview"
        className="h-full w-full mt-4 overflow-y-auto"
      >
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="overview">Visão Geral</TabsTrigger>
          <TabsTrigger value="questions">Perguntas e Respostas</TabsTrigger>
        </TabsList>
        <TabsContent value="overview" className="py-4">
          <SurveyOverviewTab />
        </TabsContent>
        <TabsContent value="questions" className="py-4">
          <p className="h-full w-full overflow-y-auto">
            Aqui ficará a lista de perguntas e as respostas dos colaboradores.
          </p>
        </TabsContent>
      </Tabs>
    </>
  );
};
