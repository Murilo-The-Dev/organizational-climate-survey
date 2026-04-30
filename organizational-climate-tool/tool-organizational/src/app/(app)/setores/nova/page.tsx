"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { setorService } from "@/lib/services/setorService";
import { useAuth } from "@/context/AuthContext";

const setorSchema = z.object({
  nome_setor: z.string().min(1, { message: "O Nome do Setor é obrigatório." }),
  descricao: z.string().optional(),
});

type SetorFormInputs = z.infer<typeof setorSchema>;

export default function NovoSetorPage() {
  const router = useRouter();
  const { user } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<SetorFormInputs>({
    resolver: zodResolver(setorSchema),
  });

  const onSubmit = async (data: SetorFormInputs) => {
    try {
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;

      await setorService.create({
        nome_setor: data.nome_setor,
        descricao: data.descricao || "",
        id_empresa: empresaId,
      });

      toast.success("Setor cadastrado com sucesso!");
      reset();
      router.push("/setores");
    } catch (error) {
      console.error("Erro ao cadastrar setor:", error);
      toast.error("Erro ao cadastrar setor. Tente novamente.");
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Novo Setor
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Preencha os dados para cadastrar um novo setor.
      </p>

      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Dados do Setor</CardTitle>
          <CardDescription>
            Insira as informações do novo setor da empresa.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid gap-2">
              <Label htmlFor="nome_setor">Nome do Setor</Label>
              <Input
                id="nome_setor"
                placeholder="Ex: Recursos Humanos"
                {...register("nome_setor")}
              />
              {errors.nome_setor && (
                <p className="text-red-500 text-sm">
                  {errors.nome_setor.message}
                </p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="descricao">Descrição (Opcional)</Label>
              <Textarea
                id="descricao"
                placeholder="Breve descrição do setor..."
                {...register("descricao")}
              />
              {errors.descricao && (
                <p className="text-red-500 text-sm">
                  {errors.descricao.message}
                </p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full bg-blue-600 hover:bg-blue-700 text-white"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Cadastrando..." : "Cadastrar Setor"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </section>
  );
}
