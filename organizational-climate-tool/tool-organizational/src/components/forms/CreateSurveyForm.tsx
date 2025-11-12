import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { PlusCircle, Trash2 } from "lucide-react";
import { useState } from "react";

const questionSchema = z.object({
  type: z.enum(["text", "single-choice", "multi-choice", "likert"]),
  questionText: z.string().min(1, { message: "O texto da pergunta é obrigatório." }),
  options: z.array(z.string().min(1, { message: "A opção não pode ser vazia." })).optional(),
});

const createSurveySchema = z.object({
  title: z.string().min(1, { message: "O título da pesquisa é obrigatório." }),
  description: z.string().min(1, { message: "A descrição da pesquisa é obrigatória." }),
  tag: z.string().min(1, { message: "A categoria da pesquisa é obrigatória." }),
  questions: z.array(questionSchema).min(1, { message: "A pesquisa deve ter pelo menos uma pergunta." }),
});

type CreateSurveyFormInputs = z.infer<typeof createSurveySchema>;

export function CreateSurveyForm() {
  const { register, handleSubmit, control, formState: { errors, isSubmitting }, reset } = useForm<CreateSurveyFormInputs>({
    resolver: zodResolver(createSurveySchema),
    defaultValues: {
      questions: [{ type: "text", questionText: "", options: [] }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: "questions",
  });

  const [selectedQuestionType, setSelectedQuestionType] = useState<z.infer<typeof questionSchema>["type"]>("text");

  const onSubmit = async (data: CreateSurveyFormInputs) => {
    try {
      // Simular chamada de API para criar pesquisa
      await new Promise(resolve => setTimeout(resolve, 1500));
      console.log("Nova pesquisa criada:", data);
      toast.success("Pesquisa criada com sucesso!");
      reset(); // Limpa o formulário
      // router.push("/pesquisas"); // Redireciona para a lista de pesquisas
    } catch (error) {
      console.error("Erro ao criar pesquisa:", error);
      toast.error("Erro ao criar pesquisa. Tente novamente.");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="grid gap-6 py-4">
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="title" className="text-right">Título</Label>
        <Input id="title" placeholder="Ex: Pesquisa de Engajamento Q3" className="col-span-3" {...register("title")} />
        {errors.title && <p className="col-span-4 text-right text-red-500 text-sm">{errors.title.message}</p>}
      </div>
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="description" className="text-right">Descrição</Label>
        <Textarea id="description" placeholder="Descreva o objetivo desta pesquisa." className="col-span-3" {...register("description")} />
        {errors.description && <p className="col-span-4 text-right text-red-500 text-sm">{errors.description.message}</p>}
      </div>
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="tag" className="text-right">Categoria</Label>
        <Select onValueChange={(value) => register("tag").onChange({ target: { value } })}>
          <SelectTrigger id="tag" className="col-span-3">
            <SelectValue placeholder="Selecione uma categoria" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="engajamento">Engajamento</SelectItem>
            <SelectItem value="lideranca">Liderança</SelectItem>
            <SelectItem value="bem-estar">Bem-estar</SelectItem>
            <SelectItem value="rh">RH</SelectItem>
          </SelectContent>
        </Select>
        {errors.tag && <p className="col-span-4 text-right text-red-500 text-sm">{errors.tag.message}</p>}
      </div>

      <div className="mt-6">
        <h3 className="text-lg font-semibold mb-4">Perguntas da Pesquisa</h3>
        {fields.map((field, index) => (
          <Card key={field.id} className="mb-4 p-4">
            <div className="flex justify-between items-center mb-4">
              <h4 className="font-medium">Pergunta {index + 1}</h4>
              <Button type="button" variant="destructive" size="sm" onClick={() => remove(index)}>
                <Trash2 className="h-4 w-4 mr-2" /> Remover
              </Button>
            </div>
            <div className="grid gap-2 mb-4">
              <Label htmlFor={`questions.${index}.questionText`}>Texto da Pergunta</Label>
              <Input id={`questions.${index}.questionText`} {...register(`questions.${index}.questionText`)} />
              {errors.questions?.[index]?.questionText && <p className="text-red-500 text-sm">{errors.questions[index]?.questionText?.message}</p>}
            </div>
            <div className="grid gap-2 mb-4">
              <Label htmlFor={`questions.${index}.type`}>Tipo de Pergunta</Label>
              <Select
                onValueChange={(value: z.infer<typeof questionSchema>["type"]) => {
                  register(`questions.${index}.type`).onChange({ target: { value } });
                  setSelectedQuestionType(value);
                }}
                defaultValue={field.type}
              >
                <SelectTrigger id={`questions.${index}.type`}>
                  <SelectValue placeholder="Selecione o tipo" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="text">Texto Livre</SelectItem>
                  <SelectItem value="single-choice">Múltipla Escolha (Única)</SelectItem>
                  <SelectItem value="multi-choice">Múltipla Escolha (Múltipla)</SelectItem>
                  <SelectItem value="likert">Escala Likert</SelectItem>
                </SelectContent>
              </Select>
              {errors.questions?.[index]?.type && <p className="text-red-500 text-sm">{errors.questions[index]?.type?.message}</p>}
            </div>

            {(field.type === "single-choice" || field.type === "multi-choice") && (
              <div className="grid gap-2">
                <Label>Opções</Label>
                {field.options?.map((option, optIndex) => (
                  <div key={optIndex} className="flex items-center gap-2">
                    <Input
                      placeholder="Nova opção"
                      {...register(`questions.${index}.options.${optIndex}`)}
                    />
                    <Button type="button" variant="ghost" size="icon" onClick={() => {
                      const currentOptions = control._fields[index]?.options?.value || [];
                      const newOptions = currentOptions.filter((_: any, i: number) => i !== optIndex);
                      control._fields[index].options.value = newOptions;
                      // Forçar re-render para atualizar o array de opções
                      // Isso é um workaround, idealmente o useFieldArray deveria gerenciar isso diretamente
                      // setForceUpdate(prev => !prev);
                    }}>
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    const currentOptions = control._fields[index]?.options?.value || [];
                    const newOptions = [...currentOptions, ""];
                    control._fields[index].options.value = newOptions;
                    // setForceUpdate(prev => !prev);
                  }}
                >
                  <PlusCircle className="h-4 w-4 mr-2" /> Adicionar Opção
                </Button>
              </div>
            )}
          </Card>
        ))}
        {errors.questions && <p className="text-red-500 text-sm">{errors.questions.message}</p>}
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => append({ type: "text", questionText: "", options: [] })}
          className="mt-4"
        >
          <PlusCircle className="h-4 w-4 mr-2" /> Adicionar Pergunta
        </Button>
      </div>

      <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white" disabled={isSubmitting}>
        {isSubmitting ? "Criando Pesquisa..." : "Criar Pesquisa"}
      </Button>
    </form>
  );
}

