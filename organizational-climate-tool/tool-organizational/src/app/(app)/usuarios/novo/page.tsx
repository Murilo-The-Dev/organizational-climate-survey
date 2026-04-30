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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { usuarioService } from "@/lib/services/usuarioService";

// Validação idêntica às regras do backend
const usuarioSchema = z.object({
  nome: z.string().min(1, { message: "O Nome é obrigatório." }),
  email: z.string().email({ message: "E-mail inválido." }),
  senha: z
    .string()
    .min(8, { message: "A senha deve ter pelo menos 8 caracteres." }),
  role: z.enum(["admin", "viewer"]),
});

type UsuarioFormInputs = z.infer<typeof usuarioSchema>;

export default function NovoUsuarioPage() {
  const router = useRouter();
  const { user } = useAuth();

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<UsuarioFormInputs>({
    resolver: zodResolver(usuarioSchema),
    defaultValues: {
      role: "admin", // Valor padrão
    },
  });

  const onSubmit = async (data: UsuarioFormInputs) => {
    try {
      const empresaId = (user as any)?.empresaId
        ? Number((user as any).empresaId)
        : 1;

      await usuarioService.create({
        nome: data.nome,
        email: data.email,
        senha: data.senha,
        role: data.role,
        id_empresa: empresaId,
      } as any);

      toast.success("Usuário cadastrado com sucesso!");
      router.push("/usuarios");
    } catch (error) {
      console.error("Erro ao cadastrar usuário:", error);
      toast.error("Erro ao cadastrar usuário. O e-mail já pode estar em uso.");
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Novo Usuário
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Cadastre um novo usuário para acessar o painel.
      </p>

      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Dados do Usuário</CardTitle>
          <CardDescription>
            Defina as credenciais e o nível de acesso.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="nome">Nome Completo</Label>
              <Input
                id="nome"
                placeholder="Ex: João Silva"
                {...register("nome")}
              />
              {errors.nome && (
                <p className="text-red-500 text-sm">{errors.nome.message}</p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="email">E-mail</Label>
              <Input
                id="email"
                type="email"
                placeholder="joao@empresa.com"
                {...register("email")}
              />
              {errors.email && (
                <p className="text-red-500 text-sm">{errors.email.message}</p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="senha">Senha Temporária</Label>
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
              <Label htmlFor="role">Nível de Acesso</Label>
              <Controller
                name="role"
                control={control}
                render={({ field }) => (
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Selecione o papel" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="admin">Administrador</SelectItem>
                      <SelectItem value="viewer">
                        Visualizador (Sem edição)
                      </SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.role && (
                <p className="text-red-500 text-sm">{errors.role.message}</p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full mt-4 bg-blue-600 hover:bg-blue-700 text-white"
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
