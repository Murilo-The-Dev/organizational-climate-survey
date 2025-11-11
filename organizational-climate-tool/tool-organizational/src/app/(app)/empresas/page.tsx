import { DataTable } from "@/components/dashboard/DataTable";
import { Button } from "@/components/ui/button";
import Link from "next/link";

interface Company {
  id: string;
  name: string;
  cnpj: string;
  status: string;
}

const mockCompanies: Company[] = [
  { id: "1", name: "Empresa A", cnpj: "11.111.111/0001-11", status: "Ativo" },
  { id: "2", name: "Empresa B", cnpj: "22.222.222/0001-22", status: "Inativo" },
  { id: "3", name: "Empresa C", cnpj: "33.333.333/0001-33", status: "Ativo" },
];

const columns = [
  { accessorKey: "name", header: "Nome da Empresa" },
  { accessorKey: "cnpj", header: "CNPJ" },
  { accessorKey: "status", header: "Status" },
  {
    id: "actions",
    header: "Ações",
    cell: ({ row }: any) => (
      <Button variant="ghost" className="h-8 w-8 p-0">
        <span className="sr-only">Abrir menu</span>
        {/* Ícone de menu ou ação */}
        ...
      </Button>
    ),
  },
];

export default function EmpresasPage() {
  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Empresas
        </h1>
        <Link href="/empresas/nova">
          <Button>Adicionar Nova Empresa</Button>
        </Link>
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie as empresas cadastradas no sistema.
      </p>

      <div className="bg-background rounded-lg border p-4 h-full">
        <DataTable columns={columns} data={mockCompanies} />
      </div>
    </section>
  );
}

