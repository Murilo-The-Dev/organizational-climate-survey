"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
import { usuarioService } from "@/lib/services/usuarioService";
import { useAuth } from "@/context/AuthContext";

// 1. Mudamos de 'nome' para 'nome_admin'
const userSchema = z.object({
  nome_admin: z.string().min(1, { message: "Nome é obrigatório." }),
  email: z.string().email({ message: "E-mail inválido." }),
  senha: z
    .string()
    .min(8, { message: "A senha deve ter no mínimo 8 caracteres." }),
  role: z.string().min(1, { message: "Função obrigatória." }),
});

type UserFormInputs = z.infer<typeof userSchema>;

export default function NovoUsuarioPage() {
  const router = useRouter();
  const { user } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<UserFormInputs>({
    resolver: zodResolver(userSchema),
  });

  const onSubmit = async (data: UserFormInputs) => {
    try {
      // 2. Mudamos a variável para 'id_empresa' pra bater certinho com o que o Service quer
      const id_empresa = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;
      const payload = { ...data, id_empresa };

      await usuarioService.create(payload as any);

      toast.success("Usuário cadastrado com sucesso!");
      reset();
      router.push("/usuarios");
    } catch (error) {
      console.error("Erro ao cadastrar usuário:", error);
      toast.error(
        "Erro ao cadastrar usuário. Verifique os dados e tente novamente.",
      );
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Novo Usuário
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Preencha os dados para cadastrar um novo usuário administrador.
      </p>

      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Dados do Usuário</CardTitle>
          <CardDescription>
            Insira as informações do novo usuário.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid gap-2">
              {/* 3. Atualizamos o HTML e o register para usar 'nome_admin' */}
              <Label htmlFor="nome_admin">Nome</Label>
              <Input
                id="nome_admin"
                placeholder="Ex: Maria Silva"
                {...register("nome_admin")}
              />
              {errors.nome_admin && (
                <p className="text-red-500 text-sm">
                  {errors.nome_admin.message}
                </p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="email">E-mail</Label>
              <Input
                id="email"
                type="email"
                placeholder="maria@empresa.com"
                {...register("email")}
              />
              {errors.email && (
                <p className="text-red-500 text-sm">{errors.email.message}</p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="senha">Senha</Label>
              <Input
                id="senha"
                type="password"
                placeholder="Mínimo de 8 caracteres"
                {...register("senha")}
              />
              {errors.senha && (
                <p className="text-red-500 text-sm">{errors.senha.message}</p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="role">Função (Role)</Label>
              <Input
                id="role"
                placeholder="Ex: admin ou viewer"
                {...register("role")}
              />
              {errors.role && (
                <p className="text-red-500 text-sm">{errors.role.message}</p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full bg-blue-600 hover:bg-blue-700 text-white"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Cadastrando..." : "Cadastrar Usuário"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </section>
  );
}
