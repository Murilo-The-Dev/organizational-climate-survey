
import { Button } from "@/components/ui/button";
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
import {  ListFilter, Filter, ArrowDownToLine } from "lucide-react";
import { ResultsDataTable } from "@/components/dashboard/ResultsDataTable";

import { data as mockResults } from "@/components/dashboard/ResultsDataTable";

const ResultadosPage = () => {
  return (
    <section className="container mx-auto px-4 mt-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Resultados Detalhados
          </h1>
          <p className="text-muted-foreground mt-2">
            Filtre e analise as respostas de cada pesquisa em detalhes.
          </p>
        </div>
        <Button className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-300"> <ArrowDownToLine className="h-5 w-5" /> Exportar Relatório</Button>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ListFilter className="h-5 w-5" />
            Filtros de Análise
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Select>
              <SelectTrigger>
                <SelectValue placeholder="Selecione a pesquisa" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="q1-2025">Engajamento Q1 2025</SelectItem>
              </SelectContent>
            </Select>
            <Select>
              <SelectTrigger>
                <SelectValue placeholder="Todos os departamentos" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="tecnologia">Tecnologia</SelectItem>
              </SelectContent>
            </Select>
            <Button className="md:w-fit cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-300"> <Filter className="h-5 w-5" /> Aplicar Filtros</Button>
          </div>
        </CardContent>
      </Card>

      <ResultsDataTable data={mockResults} />
    </section>
  );
};

export default ResultadosPage;
