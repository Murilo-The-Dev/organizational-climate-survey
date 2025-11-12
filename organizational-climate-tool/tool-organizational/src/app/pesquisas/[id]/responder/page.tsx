"use client"

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Slider } from "@/components/ui/slider";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { use } from "react";

// Esquemas de validação para os tipos de pergunta
const answerSchema = z.object({
  questionId: z.string(),
  answer: z.any(), // A resposta pode ser string, array, number, etc.
});

const publicSurveyResponseSchema = z.object({
  surveyId: z.string(),
  answers: z.array(answerSchema),
});

type PublicSurveyResponseInputs = z.infer<typeof publicSurveyResponseSchema>;

// Mock de dados de pesquisa para demonstração
const mockSurvey = {
  id: "SURV-001",
  title: "Pesquisa de Engajamento Trimestral Q3",
  description: "Sua opinião é muito importante para nós!",
  questions: [
    { id: "Q1", type: "text", questionText: "Qual sua principal sugestão para melhorar o ambiente de trabalho?" },
    { id: "Q2", type: "single-choice", questionText: "Você se sente valorizado(a) na empresa?", options: ["Sim", "Não", "Às vezes"] },
    { id: "Q3", type: "multi-choice", questionText: "Quais benefícios você considera mais importantes? (Selecione todos que se aplicam)", options: ["Plano de Saúde", "Vale Refeição", "Flexibilidade de Horário", "Desenvolvimento Profissional"] },
    { id: "Q4", type: "likert", questionText: "Concordo totalmente que a comunicação interna é eficaz.", options: ["Discordo Totalmente", "Discordo", "Neutro", "Concordo", "Concordo Totalmente"] },
  ],
};

export default function PublicSurveyResponsePage({ params }: { params: Promise<{ id: string }> }) {
  const router = useRouter();
  const { id } = use(params); // Desembrulha a Promise

  const { handleSubmit, register, setValue, watch, formState: { isSubmitting } } = useForm<PublicSurveyResponseInputs>({
    resolver: zodResolver(publicSurveyResponseSchema),
    defaultValues: {
      surveyId: id,
      answers: mockSurvey.questions.map(q => ({ questionId: q.id, answer: undefined }))
    }
  });

  const onSubmit = async (data: PublicSurveyResponseInputs) => {
    try {
      // Simular envio da resposta da pesquisa
      await new Promise(resolve => setTimeout(resolve, 2000));
      console.log("Resposta da pesquisa enviada:", data);
      toast.success("Sua resposta foi enviada com sucesso! Obrigado por participar.");
      router.push("/pesquisas/agradecimento"); // Redireciona para uma página de agradecimento
    } catch (error) {
      console.error("Erro ao enviar resposta:", error);
      toast.error("Ocorreu um erro ao enviar sua resposta. Por favor, tente novamente.");
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-100 p-4">
      <Card className="w-full max-w-3xl">
        <CardHeader className="text-center">
          <CardTitle className="text-3xl font-bold">{mockSurvey.title}</CardTitle>
          <CardDescription className="text-gray-600 mt-2">
            {mockSurvey.description}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
            {mockSurvey.questions.map((question, index) => (
              <div key={question.id} className="space-y-4 border-b pb-6 last:border-b-0 last:pb-0">
                <Label className="text-lg font-semibold">
                  {index + 1}. {question.questionText}
                </Label>
                {question.type === "text" && (
                  <Textarea
                    placeholder="Digite sua resposta aqui..."
                    {...register(`answers.${index}.answer`)}
                  />
                )}
                {question.type === "single-choice" && question.options && (
                  <RadioGroup
                    onValueChange={(value) => setValue(`answers.${index}.answer`, value)}
                    defaultValue={watch(`answers.${index}.answer`)}
                  >
                    {question.options.map((option, optIndex) => (
                      <div key={optIndex} className="flex items-center space-x-2">
                        <RadioGroupItem value={option} id={`${question.id}-${optIndex}`} />
                        <Label htmlFor={`${question.id}-${optIndex}`}>{option}</Label>
                      </div>
                    ))}
                  </RadioGroup>
                )}
                {question.type === "multi-choice" && question.options && (
                  <div className="space-y-2">
                    {question.options.map((option, optIndex) => (
                      <div key={optIndex} className="flex items-center space-x-2">
                        <Checkbox
                          id={`${question.id}-${optIndex}`}
                          onCheckedChange={(checked) => {
                            const currentAnswers = watch(`answers.${index}.answer`) || [];
                            const newAnswers = checked
                              ? [...currentAnswers, option]
                              : currentAnswers.filter((ans: string) => ans !== option);
                            setValue(`answers.${index}.answer`, newAnswers);
                          }}
                        />
                        <Label htmlFor={`${question.id}-${optIndex}`}>{option}</Label>
                      </div>
                    ))}
                  </div>
                )}
                {question.type === "likert" && question.options && (
                  <div className="space-y-4">
                    <Slider
                      defaultValue={[Math.floor(question.options.length / 2)]}
                      max={question.options.length - 1}
                      step={1}
                      onValueChange={(value) => setValue(`answers.${index}.answer`, question.options[value[0]])}
                      className="w-[60%]"
                    />
                    <div className="flex justify-between text-sm text-muted-foreground">
                      {question.options.map((option, optIndex) => (
                        <span key={optIndex}>{option}</span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
            <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white" disabled={isSubmitting}>
              {isSubmitting ? "Enviando Resposta..." : "Enviar Resposta"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}