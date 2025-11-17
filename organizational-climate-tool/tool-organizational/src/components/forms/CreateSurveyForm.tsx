// CreateSurveyForm.tsx
"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { PlusCircle, Trash2 } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Card } from "@/components/ui/card";

// Definição do Schema para as Opções de Múltipla Escolha
const optionSchema = z.object({
  text: z.string().min(1, "A opção não pode ser vazia."),
});

// Definição do Schema para as Perguntas
const questionSchema = z.object({
  text: z.string().min(1, "O texto da pergunta é obrigatório."),
  type: z.enum(["text", "radio", "checkbox", "scale"], {
    required_error: "O tipo de pergunta é obrigatório.",
  }),
  options: z.array(optionSchema).optional(), // Opcional para perguntas de texto
});

// Definição do Schema Principal
const surveySchema = z.object({
  title: z.string().min(3, "O título deve ter no mínimo 3 caracteres."),
  description: z.string().optional(),
  questions: z.array(questionSchema).min(1, "A pesquisa deve ter pelo menos 1 pergunta."),
});

type SurveyFormData = z.infer<typeof surveySchema>;

interface CreateSurveyFormProps {
  onClose?: () => void;
}

export function CreateSurveyForm({ onClose }: CreateSurveyFormProps) {
  const form = useForm<SurveyFormData>({
    resolver: zodResolver(surveySchema),
    defaultValues: {
      title: "",
      description: "",
      questions: [{ text: "", type: "text", options: [] }],
    },
  });

  // USE useFieldArray CORRETAMENTE
  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "questions",
  });

  const onSubmit = async (data: SurveyFormData) => {
    try {
      // Simulação de chamada de API
      await new Promise(resolve => setTimeout(resolve, 1500));
      console.log("Nova Pesquisa:", data);
      toast.success("Pesquisa criada com sucesso!");
      if (onClose) onClose(); // Fecha o modal após o sucesso
    } catch (error) {
      toast.error("Erro ao criar pesquisa. Tente novamente.");
    }
  };

  const addQuestion = () => {
    append({ text: "", type: "text", options: [] });
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 max-h-[70vh] overflow-y-auto p-2">
      {/* Título e Descrição */}
      <div className="grid gap-2">
        <Label htmlFor="title">Título da Pesquisa</Label>
        <Input id="title" {...form.register("title")} />
        {form.formState.errors.title && (
          <p className="text-red-500 text-sm">{form.formState.errors.title.message}</p>
        )}
      </div>
      <div className="grid gap-2">
        <Label htmlFor="description">Descrição (Opcional)</Label>
        <Textarea id="description" {...form.register("description")} />
      </div>

      {/* Seção de Perguntas */}
      <h3 className="text-xl font-semibold mt-6">Perguntas</h3>
      <div className="space-y-4">
        {fields.map((field, index) => (
          <Card key={field.id} className="p-4 border-l-4 border-blue-500">
            <div className="flex justify-between items-start mb-3">
              <h4 className="font-medium">Pergunta #{index + 1}</h4>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => remove(index)}
                className="text-red-500 hover:text-red-700"
              >
                <Trash2 className="w-4 h-4" />
              </Button>
            </div>

            <div className="grid gap-2 mb-3">
              <Label htmlFor={`questions.${index}.text`}>Texto da Pergunta</Label>
              <Input id={`questions.${index}.text`} {...form.register(`questions.${index}.text`)} />
              {form.formState.errors.questions?.[index]?.text && (
                <p className="text-red-500 text-sm">{form.formState.errors.questions[index].text.message}</p>
              )}
            </div>

            <div className="grid gap-2 mb-3">
              <Label>Tipo de Resposta</Label>
              <Select
                onValueChange={(value) => form.setValue(`questions.${index}.type`, value as any)}
                defaultValue={field.type}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o tipo" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="text">Texto Livre</SelectItem>
                  <SelectItem value="radio">Múltipla Escolha (Única)</SelectItem>
                  <SelectItem value="checkbox">Múltipla Escolha (Múltipla)</SelectItem>
                  <SelectItem value="scale">Escala (1 a 5)</SelectItem>
                </SelectContent>
              </Select>
              {form.formState.errors.questions?.[index]?.type && (
                <p className="text-red-500 text-sm">{form.formState.errors.questions[index].type.message}</p>
              )}
            </div>

            {/* Lógica para Opções (Apenas para Múltipla Escolha/Escala) */}
            {(field.type === "radio" || field.type === "checkbox") && (
              <div className="mt-4 p-3 border rounded-md">
                <h5 className="font-medium mb-2">Opções de Resposta</h5>
                <p className="text-sm text-muted-foreground">Lógica de opções aninhadas seria implementada aqui.</p>
              </div>
            )}
          </Card>
        ))}
      </div>

      <Button type="button" variant="outline" onClick={addQuestion} className="w-full">
        <PlusCircle className="w-4 h-4 mr-2" /> Adicionar Pergunta
      </Button>

      {/* Botão de Submissão */}
      <div className="flex justify-end pt-4 border-t">
        <Button type="submit" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? "Criando..." : "Criar Pesquisa"}
        </Button>
      </div>
    </form>
  );
}

