"use client";

import { useEffect } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";
import { usuarioService } from "@/lib/services/usuarioService";
import { authService } from "@/lib/services/authService";

// Validação para o formulário de Perfil
const profileSchema = z.object({
  nome: z.string().min(1, "O nome é obrigatório."),
  email: z.string().email("E-mail inválido."),
});

// Validação para o formulário de Troca de Senha
const passwordSchema = z
  .object({
    senhaAtual: z.string().min(1, "A senha atual é obrigatória."),
    novaSenha: z
      .string()
      .min(8, "A nova senha deve ter no mínimo 8 caracteres."),
    confirmarSenha: z.string().min(1, "Confirme sua nova senha."),
  })
  .refine((data) => data.novaSenha === data.confirmarSenha, {
    message: "As senhas não coincidem.",
    path: ["confirmarSenha"],
  });

type ProfileForm = z.infer<typeof profileSchema>;
type PasswordForm = z.infer<typeof passwordSchema>;

export default function ConfiguracoesPage() {
  const { user } = useAuth();

  // Setup do formulário de Perfil
  const {
    register: registerProfile,
    handleSubmit: handleSubmitProfile,
    reset: resetProfile,
    formState: { errors: errorsProfile, isSubmitting: isSubmittingProfile },
  } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
  });

  // Setup do formulário de Senha
  const {
    register: registerPassword,
    handleSubmit: handleSubmitPassword,
    reset: resetPassword,
    formState: { errors: errorsPassword, isSubmitting: isSubmittingPassword },
  } = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
  });

  // Preenche os dados atuais do usuário assim que a tela carrega
  useEffect(() => {
    if (user) {
      resetProfile({
        nome: user.nome,
        email: user.email,
      });
    }
  }, [user, resetProfile]);

  const onUpdateProfile = async (data: ProfileForm) => {
    try {
      if (!user?.id) return;
      // Chamada para a API (usamos 'as any' para ignorar o type strict da senha aqui)
      await usuarioService.update(user.id, {
        nome_admin: data.nome,
        email: data.email,
      } as any);

      toast.success(
        "Perfil atualizado! Você verá as alterações no próximo login.",
      );
    } catch (error) {
      console.error(error);
      toast.error(
        "Erro ao atualizar perfil. Verifique os dados e tente novamente.",
      );
    }
  };

  const onUpdatePassword = async (data: PasswordForm) => {
    try {
      await authService.changePassword(data.senhaAtual, data.novaSenha);
      toast.success("Senha alterada com sucesso!");
      resetPassword(); // Limpa os campos após o sucesso
    } catch (error) {
      console.error(error);
      toast.error(
        "Erro ao alterar senha. Verifique se a sua senha atual está correta.",
      );
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Configurações
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie as configurações da sua conta e do sistema.
      </p>

      <div className="space-y-8 max-w-3xl">
        {/* Bloco 1: Perfil */}
        <Card>
          <CardHeader>
            <CardTitle>Perfil</CardTitle>
            <CardDescription>
              Atualize as informações do seu perfil de administrador.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              onSubmit={handleSubmitProfile(onUpdateProfile)}
              className="space-y-4"
            >
              <div className="grid gap-2">
                <Label htmlFor="nome">Nome Completo</Label>
                <Input id="nome" {...registerProfile("nome")} />
                {errorsProfile.nome && (
                  <p className="text-sm text-red-500">
                    {errorsProfile.nome.message}
                  </p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="email">E-mail corporativo</Label>
                <Input id="email" type="email" {...registerProfile("email")} />
                {errorsProfile.email && (
                  <p className="text-sm text-red-500">
                    {errorsProfile.email.message}
                  </p>
                )}
              </div>
              <Button
                type="submit"
                disabled={isSubmittingProfile}
                className="bg-blue-600 hover:bg-blue-700 text-white"
              >
                {isSubmittingProfile ? "Salvando..." : "Salvar alterações"}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Bloco 2: Segurança (Senha) */}
        <Card>
          <CardHeader>
            <CardTitle>Segurança</CardTitle>
            <CardDescription>
              Altere sua senha de acesso periodicamente para manter a conta
              segura.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              onSubmit={handleSubmitPassword(onUpdatePassword)}
              className="space-y-4"
            >
              <div className="grid gap-2">
                <Label htmlFor="senhaAtual">Senha Atual</Label>
                <Input
                  id="senhaAtual"
                  type="password"
                  {...registerPassword("senhaAtual")}
                />
                {errorsPassword.senhaAtual && (
                  <p className="text-sm text-red-500">
                    {errorsPassword.senhaAtual.message}
                  </p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="novaSenha">Nova Senha</Label>
                <Input
                  id="novaSenha"
                  type="password"
                  placeholder="Mínimo 8 caracteres"
                  {...registerPassword("novaSenha")}
                />
                {errorsPassword.novaSenha && (
                  <p className="text-sm text-red-500">
                    {errorsPassword.novaSenha.message}
                  </p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="confirmarSenha">Confirmar Nova Senha</Label>
                <Input
                  id="confirmarSenha"
                  type="password"
                  {...registerPassword("confirmarSenha")}
                />
                {errorsPassword.confirmarSenha && (
                  <p className="text-sm text-red-500">
                    {errorsPassword.confirmarSenha.message}
                  </p>
                )}
              </div>
              <Button
                type="submit"
                disabled={isSubmittingPassword}
                className="bg-blue-600 hover:bg-blue-700 text-white"
              >
                {isSubmittingPassword ? "Alterando..." : "Alterar Senha"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </section>
  );
}
