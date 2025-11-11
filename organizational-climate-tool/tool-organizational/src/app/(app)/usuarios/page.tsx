import { DataTable } from "@/components/dashboard/DataTable";
import { Button } from "@/components/ui/button";
import Link from "next/link";

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

const mockUsers: User[] = [
  { id: "1", name: "Alexandre Calore", email: "alexandre@example.com", role: "Administrador" },
  { id: "2", name: "Guilherme Conceição", email: "guilherme@example.com", role: "Editor" },
  { id: "3", name: "Usuário Teste", email: "teste@example.com", role: "Visualizador" },
];

const columns = [
  { accessorKey: "name", header: "Nome" },
  { accessorKey: "email", header: "E-mail" },
  { accessorKey: "role", header: "Função" },
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

export default function UsuariosPage() {
  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex justify-between items-center mb-6">
        <h1 className="w-fit text-3xl font-bold tracking-tight bg-blue-500 text-white p-2 rounded-lg">
          Usuários
        </h1>
        <Link href="/usuarios/novo">
          <Button>Adicionar Novo Usuário</Button>
        </Link>
      </div>
      <p className="text-muted-foreground mt-2 mb-6">
        Gerencie os usuários administradores do sistema.
      </p>

      <div className="bg-background rounded-lg border p-4 h-full">
        <DataTable columns={columns} data={mockUsers} />
      </div>
    </section>
  );
}

