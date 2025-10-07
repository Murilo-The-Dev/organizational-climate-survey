import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import Link from 'next/link';

export default function LoginPage() {
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
              <form className="space-y-6">
                <div className="grid gap-2">
                  <Label htmlFor="email">E-mail</Label>
                  <Input id="email" type="email" placeholder="m@example.com" required />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="password">Senha</Label>
                  <Input id="password" type="password" required />
                </div>
                <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white">
                  Continuar
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
          {/* Conteúdo opcional para o lado direito, como uma imagem ou ilustração */}
        </div>
      </div>
    </div>
  );
}

