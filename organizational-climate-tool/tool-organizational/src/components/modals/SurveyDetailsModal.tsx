"use client";

import {
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { type Pesquisa } from "@/components/dashboard/DataTable";
import { SurveyOverviewTab } from "../pesquisas/SurveyOverviewTab";
import { SurveyQuestionsTab } from "../pesquisas/SurveyQuestionsTab";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ExternalDataForm } from "@/components/forms/ExternalDataForm";

type SurveyDetailsModalProps = {
  description: string;
  tag: string;
  creationDate: string;
  survey: Pesquisa;
};

export const SurveyDetailsModal = ({ survey }: SurveyDetailsModalProps) => {
  if (!survey) return null;

  return (
    <>
      <DialogHeader>
        <DialogTitle>{survey.title}</DialogTitle>
        <DialogDescription>
          Análise detalhada e gerenciamento da pesquisa.
        </DialogDescription>
      </DialogHeader>

      <Tabs defaultValue="overview" className="h-full w-full mt-4 overflow-y-auto">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="overview">Visão Geral</TabsTrigger>
          <TabsTrigger value="questions">Perguntas e Respostas</TabsTrigger>
          <TabsTrigger value="external-data">Indicadores de RH</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="py-4">
          <SurveyOverviewTab survey={survey} />
        </TabsContent>

        <TabsContent value="questions" className="py-4">
          <SurveyQuestionsTab survey={survey} />
        </TabsContent>

        <TabsContent value="external-data" className="py-4">
          <Card>
            <CardHeader>
              <CardTitle>Dados Externos de RH</CardTitle>
              <CardDescription>
                Insira indicadores como absenteísmo e turnover para cruzar com os resultados da pesquisa.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ExternalDataForm surveyId={survey.id!} />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </>
  );
};
