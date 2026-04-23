import Link from "next/link";
import { CheckCircle2 } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function AgradecimentoPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-100 p-4">
      <Card className="w-full max-w-md text-center">
        <CardHeader>
          <div className="flex justify-center mb-4">
            <CheckCircle2 className="h-16 w-16 text-green-500" />
          </div>
          <CardTitle className="text-2xl font-bold">Resposta enviada!</CardTitle>
          <CardDescription className="mt-2 text-gray-600">
            Obrigado pela sua participação. Suas respostas foram registradas com sucesso e contribuirão para a melhoria do ambiente de trabalho.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button asChild variant="outline" className="mt-2">
            <Link href="/">Voltar ao início</Link>
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
