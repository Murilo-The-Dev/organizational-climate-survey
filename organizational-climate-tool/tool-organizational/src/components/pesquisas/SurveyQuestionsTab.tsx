// src/components/dashboard/tabs/SurveyQuestionsTab.tsx

"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";

type QuestionResult = {
  id: string;
  text: string;
  type: "escala" | "multipla-escolha" | "texto-aberto";
  results: any;
  category: string;
};

const mockQuestions: QuestionResult[] = [
  {
    id: "q1",
    text: "Em uma escala de 1 a 5, quão satisfeito você está com a liderança?",
    type: "escala",
    category: "Liderança",
    results: { "1": 10, "2": 25, "3": 40, "4": 120, "5": 98 },
  },
  {
    id: "q2",
    text: "Qual benefício você mais valoriza na empresa?",
    type: "multipla-escolha",
    category: "Bem-estar",
    results: {
      total: 293,
      options: [
        { label: "Plano de Saúde", count: 150 },
        { label: "Vale Alimentação", count: 80 },
        { label: "Horário Flexível", count: 63 },
      ],
    },
  },
  {
    id: "q3",
    text: "Que sugestões você daria para melhorar a comunicação interna?",
    type: "texto-aberto",
    category: "Comunicação",
    results: [
      "Ter reuniões semanais mais curtas e objetivas.",
      "Criar um canal de feedback anônimo mais divulgado.",
      "Melhorar a transparência nas decisões da diretoria.",
      "A newsletter interna poderia ser mais focada em conquistas das equipes.",
    ],
  },
  {
    id: "q4",
    text: "Sinto que tenho autonomia para tomar decisões no meu trabalho.",
    type: "escala",
    category: "Cultura Organizacional",
    results: { "1": 5, "2": 15, "3": 55, "4": 110, "5": 108 },
  },
  {
    id: "q5",
    text: "Com que frequência você tem reuniões 1-on-1 com seu gestor?",
    type: "multipla-escolha",
    category: "Liderança",
    results: {
      total: 293,
      options: [
        { label: "Semanalmente", count: 120 },
        { label: "Quinzenalmente", count: 95 },
        { label: "Mensalmente", count: 58 },
        { label: "Raramente ou nunca", count: 20 },
      ],
    },
  },
  {
    id: "q6",
    text: "Descreva um momento em que você se sentiu orgulhoso(a) de trabalhar aqui.",
    type: "texto-aberto",
    category: "Bem-estar",
    results: [
      "Quando lançamos o projeto X e recebemos elogios do cliente.",
      "No último evento de confraternização da empresa.",
      "Quando meu gestor reconheceu meu esforço publicamente.",
    ],
  },
];

const AnswerDetails = ({ question }: { question: QuestionResult }) => {
  switch (question.type) {
    case "escala":
      const totalResponses = Object.values(question.results).reduce(
        (sum: number, count: any) => sum + count,
        0
      );
      return (
        <div className="flex flex-col gap-2">
          {Object.entries(question.results).map(
            ([score, count]: [string, any]) => (
              <div key={score} className="flex items-center gap-4">
                <span className="text-sm font-medium w-12">Nota {score}</span>
                <Progress
                  value={(count / totalResponses) * 100}
                  className="flex-1"
                />
                <span className="text-sm text-muted-foreground">
                  {count} resp.
                </span>
              </div>
            )
          )}
        </div>
      );
    case "multipla-escolha":
      return (
        <div className="flex flex-col gap-2">
          {question.results.options.map((opt: any) => (
            <div key={opt.label} className="flex items-center gap-4">
              <span className="text-sm font-medium w-32 truncate">
                {opt.label}
              </span>
              <Progress
                value={(opt.count / question.results.total) * 100}
                className="flex-1"
              />
              <span className="text-sm text-muted-foreground">
                {opt.count} resp.
              </span>
            </div>
          ))}
        </div>
      );
    case "texto-aberto":
      return (
        <div className="flex flex-col gap-3">
          {question.results.map((text: string, index: number) => (
            <Card key={index}>
              <CardContent className="p-4 text-sm text-muted-foreground">
                "{text}"
              </CardContent>
            </Card>
          ))}
        </div>
      );
    default:
      return null;
  }
};

export const SurveyQuestionsTab = () => {
  return (
    <div className="overflow-y-auto pr-4 h-full">
      <Accordion type="single" collapsible className="w-full">
        {mockQuestions.map((question) => (
          <AccordionItem value={question.id} key={question.id}>
            <AccordionTrigger>
              <div className="flex items-center gap-4 text-left">
                <Badge>{question.category || question.type}</Badge>
                <span>{question.text}</span>
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <AnswerDetails question={question} />
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  );
};
