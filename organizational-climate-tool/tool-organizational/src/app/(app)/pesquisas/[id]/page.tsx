"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { SurveyResponseDetails } from "@/components/pesquisas/SurveyResponseDetails";
import { SurveyHistoricalTrends } from "@/components/pesquisas/SurveyHistoricalTrends";
import { ExportReportButton } from "@/components/ui/export-report-button";

const mockSurvey = {
  id: "SURV-001",
  title: "Engajamento Q1 2025",
  status: "concluido",
  participants: 152,
  createdAt: "2025-03-28",
};

export default function SurveyDetailsPage({ params }: { params: { id: string } }) {
  const [dateRange, setDateRange] = useState<any | undefined>(undefined);

  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center space-x-4">
          <Link href="/pesquisas">
            <Button variant="outline" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>

          <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
            Resultados: {mockSurvey.title}
          </h1>
        </div>

        <div className="flex items-center space-x-4">
          <ExportReportButton surveyId={params.id} surveyName={mockSurvey.title} />
        </div>
      </div>
      <Card className="mb-6">
        <CardContent className="pt-6">
          <div className="grid grid-cols-3 gap-4 text-sm">
            <p><strong>ID:</strong> {mockSurvey.id}</p>
            <p><strong>Status:</strong> {mockSurvey.status}</p>
            <p><strong>Participantes:</strong> {mockSurvey.participants}</p>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="tendencia">
        <TabsList>
          <TabsTrigger value="tendencia">Ver Tendência Histórica</TabsTrigger>
          <TabsTrigger value="respostas">Ver Detalhes das Respostas</TabsTrigger>
        </TabsList>

        <TabsContent value="tendencia" className="mt-4">
          <SurveyHistoricalTrends surveyId={params.id} dateRange={dateRange} />
        </TabsContent>
        <TabsContent value="respostas" className="mt-4">
          <SurveyResponseDetails surveyId={params.id} />
        </TabsContent>
      </Tabs>
    </section>
  );
}
