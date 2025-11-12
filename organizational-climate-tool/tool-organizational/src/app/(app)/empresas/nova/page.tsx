import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

const companySchema = z.object({
  name: z.string().min(1, { message: "Nome da empresa é obrigatório." }),
  cnpj: z.string().regex(/^\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}$/, { message: "CNPJ inválido. Use o formato XX.XXX.XXX/XXXX-XX." }),
});

type CompanyFormInputs = z.infer<typeof companySchema>;

export default function NovaEmpresaPage() {
  const router = useRouter();
  const { register, handleSubmit, formState: { errors, isSubmitting }, reset } = useForm<CompanyFormInputs>({
    resolver: zodResolver(companySchema),
  });

  const onSubmit = async (data: CompanyFormInputs) => {
    try {
      // Simular chamada de API para cadastrar empresa
      await new Promise(resolve => setTimeout(resolve, 1500));
      console.log("Nova empresa cadastrada:", data);
      toast.success("Empresa cadastrada com sucesso!");
      reset(); // Limpa o formulário
      router.push("/empresas"); // Redireciona para a lista de empresas
    } catch (error) {
      console.error("Erro ao cadastrar empresa:", error);
      toast.error("Erro ao cadastrar empresa. Tente novamente.");
    }
  };

  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
        Nova Empresa
      </h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Preencha os dados para cadastrar uma nova empresa.
      </p>

      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Dados da Empresa</CardTitle>
          <CardDescription>Insira as informações da nova empresa.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid gap-2">
              <Label htmlFor="name">Nome da Empresa</Label>
              <Input id="name" {...register("name")} />
              {errors.name && <p className="text-red-500 text-sm">{errors.name.message}</p>}
            </div>
            <div className="grid gap-2">
              <Label htmlFor="cnpj">CNPJ</Label>
              <Input id="cnpj" placeholder="XX.XXX.XXX/XXXX-XX" {...register("cnpj")} />
              {errors.cnpj && <p className="text-red-500 text-sm">{errors.cnpj.message}</p>}
            </div>
            <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white" disabled={isSubmitting}>
              {isSubmitting ? "Cadastrando..." : "Cadastrar Empresa"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </section>
  );
}

