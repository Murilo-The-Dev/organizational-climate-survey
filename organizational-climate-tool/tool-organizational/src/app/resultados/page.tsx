
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
import { ListFilter } from "lucide-react";
import { FormsData } from "@/components/dashboard/FormsData";


const ResultadosPage = () => {
  return (
    <section className="container mx-auto px-4 mt-10">
      <h1 className="bg-blue-500 text-white p-2 rounded-lg w-fit text-3xl font-bold tracking-tight">Resultados Detalhados</h1>
      <p className="text-muted-foreground mt-2 mb-6">
        Filtre e analise as respostas de cada pesquisa em detalhes.
      </p>
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ListFilter className="h-5 w-5 text-blue-700" />
            Filtros de Análise
          </CardTitle>
          <CardDescription>
            Selecione uma pesquisa e refine os dados para encontrar insights.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <div className="flex flex-col gap-2">
              <label htmlFor="pesquisa" className="text-sm font-medium">
                Pesquisa
              </label>
              <Select>
                <SelectTrigger id="pesquisa">
                  <SelectValue placeholder="Selecione a pesquisa" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="q1-2025">Engajamento Q1 2025</SelectItem>
                  <SelectItem value="q4-2024">Liderança Q4 2024</SelectItem>
                  <SelectItem value="e-nps-2024">e-NPS Anual 2024</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex flex-col gap-2">
              <label htmlFor="departamento" className="text-sm font-medium">
                Departamento
              </label>
              <Select>
                <SelectTrigger id="departamento">
                  <SelectValue placeholder="Todos os departamentos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="tecnologia">Tecnologia</SelectItem>
                  <SelectItem value="marketing">Marketing</SelectItem>
                  <SelectItem value="rh">Recursos Humanos</SelectItem>
                </SelectContent>
              </Select>
            </div>

            
            <Button className="w-fit self-end bg-blue-600 text-white hover:bg-blue-500 cursor-pointer">
              Aplicar Filtros
            </Button>
          </div>
        </CardContent>
      </Card>
      
      <div className="mt-6">
            <FormsData />
      </div>
    </section>
  );
};

export default ResultadosPage;