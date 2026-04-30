"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import type { Pesquisa } from "@/lib/types";

interface SurveyQuestionsTabProps {
  survey: Pesquisa;
}

export const SurveyQuestionsTab = ({ survey }: SurveyQuestionsTabProps) => {
  const perguntas = survey.perguntas || [];

  if (perguntas.length === 0) {
    return (
      <p className="text-center py-10 text-muted-foreground">
        Nenhuma pergunta cadastrada para esta pesquisa.
      </p>
    );
  }

  return (
    <div className="overflow-y-auto pr-4 h-full">
      <Accordion type="single" collapsible className="w-full">
        {perguntas.map((q) => (
          <AccordionItem value={String(q.id_pergunta)} key={q.id_pergunta}>
            <AccordionTrigger className="hover:no-underline">
              <div className="flex items-center gap-4 text-left">
                <Badge variant="secondary">{q.tipo_pergunta}</Badge>
                <span>{q.texto_pergunta}</span>
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <div className="p-4 bg-slate-50 rounded-lg">
                <p className="text-sm text-muted-foreground mb-2">
                  Resultados consolidados aparecem no relatório analítico.
                </p>
                <div className="flex items-center gap-2">
                  <Progress value={0} className="flex-1" />
                  <span className="text-xs">Aguardando processamento</span>
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  );
};
