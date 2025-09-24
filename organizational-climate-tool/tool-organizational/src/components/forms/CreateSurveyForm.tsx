// src/components/forms/CreateSurveyForm.tsx

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function CreateSurveyForm() {
  return (
    <div className="grid gap-6 py-4">
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="title" className="text-right">
          Título
        </Label>
        <Input
          id="title"
          placeholder="Ex: Pesquisa de Engajamento Q3"
          className="col-span-3"
        />
      </div>
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="description" className="text-right">
          Descrição
        </Label>
        <Textarea
          id="description"
          placeholder="Descreva o objetivo desta pesquisa."
          className="col-span-3"
        />
      </div>
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="tag" className="text-right">
          Categoria
        </Label>
        <Select>
          <SelectTrigger id="tag" className="col-span-3">
            <SelectValue placeholder="Selecione uma categoria" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="engajamento">Engajamento</SelectItem>
            <SelectItem value="lideranca">Liderança</SelectItem>
            <SelectItem value="bem-estar">Bem-estar</SelectItem>
            <SelectItem value="rh">RH</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
