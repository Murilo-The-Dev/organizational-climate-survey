"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Suspense } from 'react';
import * as z from 'zod';
import { useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useRouter, useSearchParams } from 'next/navigation';
import { AxiosError } from 'axios';
import type { ApiResponse } from '@/lib/api';

const loginSchema = z.object({
  email: z.string().email({ message: 'E-mail inválido.' }),
  password: z.string().min(8, { message: 'A senha deve ter no mínimo 8 caracteres.' }),
});

type LoginFormInputs = z.infer<typeof loginSchema>;

function LoginForm() {
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginFormInputs>({
    resolver: zodResolver(loginSchema),
  });
  const [loginError, setLoginError] = useState<string | null>(null);
  const { login } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();

  const onSubmit = async (data: LoginFormInputs) => {
    setLoginError(null);
    try {
      await login(data.email, data.password);
      const redirect = searchParams.get('redirect') || '/dashboard';
      router.push(redirect);
    } catch (err) {
      const axiosError = err as AxiosError<ApiResponse<unknown>>;
      const message = axiosError.response?.data?.message || axiosError.response?.data?.error;
      setLoginError(message || 'Ocorreu um erro ao tentar fazer login. Tente novamente.');
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-blue-600 p-4">
      <div className="relative flex w-full max-w-6xl h-[600px] rounded-3xl overflow-hidden shadow-2xl bg-white">
        {/* Lado esquerdo: Formulário de Login */}
        <div className="flex-1 flex items-center justify-center p-8">
          <Card className="w-full max-w-md border-none shadow-none">
            <CardHeader className="text-center">
              <CardTitle className="text-2xl font-bold">Acesse nossa plataforma!</CardTitle>
              <CardDescription className="text-gray-600 mt-2">
                Digite suas credenciais de acesso para entrar no sistema. Utilize o e-mail corporativo e a senha cadastrada.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                <div className="grid gap-2">
                  <Label htmlFor="email">E-mail</Label>
                  <Input id="email" type="email" placeholder="m@example.com" {...register('email')} />
                  {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="password">Senha</Label>
                  <Input id="password" type="password" {...register('password')} />
                  {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
                </div>
                {loginError && <p className="text-red-500 text-sm text-center">{loginError}</p>}
                <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white" disabled={isSubmitting}>
                  {isSubmitting ? 'Carregando...' : 'Continuar'}
                </Button>
              </form>
              <div className="mt-6 text-center text-sm">
                Ainda não tem uma conta?{' '}
                <Link href="#" className="underline text-blue-600 hover:text-blue-700">
                  Cadastre-se!
                </Link>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Lado direito: Gradiente de fundo */}
        <div className="flex-1 bg-gradient-to-br from-blue-500 to-blue-800 hidden md:flex items-center justify-center">
        </div>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense>
      <LoginForm />
    </Suspense>
  );
}

