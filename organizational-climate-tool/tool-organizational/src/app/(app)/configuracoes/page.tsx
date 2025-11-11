import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';

export default function ConfiguracoesPage() {
  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Configurações
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie as configurações da sua conta e do sistema.
      </p>

      <div className="space-y-8">
        <Card>
          <CardHeader>
            <CardTitle>Perfil</CardTitle>
            <CardDescription>Atualize as informações do seu perfil.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="nome">Nome</Label>
              <Input id="nome" defaultValue="Eduarda" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="email">E-mail</Label>
              <Input id="email" type="email" defaultValue="eduarda@example.com" />
            </div>
            <Button>Salvar alterações</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Segurança</CardTitle>
            <CardDescription>Altere sua senha e outras configurações de segurança.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="senhaAtual">Senha Atual</Label>
              <Input id="senhaAtual" type="password" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="novaSenha">Nova Senha</Label>
              <Input id="novaSenha" type="password" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="confirmarSenha">Confirmar Nova Senha</Label>
              <Input id="confirmarSenha" type="password" />
            </div>
            <Button>Alterar Senha</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Notificações</CardTitle>
            <CardDescription>Gerencie suas preferências de notificação.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Adicionar componentes de checkbox ou switch para notificações */}
            <p>Em breve...</p>
            <Button>Salvar preferências</Button>
          </CardContent>
        </Card>
      </div>
    </section>
  );
}

