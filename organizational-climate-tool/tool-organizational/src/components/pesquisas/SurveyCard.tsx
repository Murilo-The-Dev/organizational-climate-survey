import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Eye, QrCode, MoreHorizontal } from "lucide-react";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger, DropdownMenuSeparator } from "@/components/ui/dropdown-menu";

type SurveyCardProps = {
  id: string;
  title: string;
  description: string;
  tag: string;
  creationDate: string;
  onViewDetails: () => void;
  onGenerateLink: (id: string) => void;
};

export const SurveyCard = ({
  id,
  title,
  description,
  tag,
  creationDate,
  onViewDetails,
  onGenerateLink,
}: SurveyCardProps) => {
  return (
    <Card className="hover:shadow-lg hover:border-blue-600 hover:translate-y-[-5px] transition-all duration-500">
      <CardHeader>
        <Badge variant="secondary" className="w-fit">
          {tag}
        </Badge>
        <div className="flex justify-between items-start pt-2">
          <CardTitle>{title}</CardTitle>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0 -mt-1 text-muted-foreground">
                <span className="sr-only">Abrir menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel className="text-center">Ações</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Editar Pesquisa</DropdownMenuItem>
              <DropdownMenuItem className="text-red-600">Excluir</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <p className="text-sm text-muted-foreground line-clamp-2">
          {description}
        </p>
      </CardContent>
      <CardFooter className="flex justify-between items-center text-sm text-muted-foreground">
        <span>Criado em: {creationDate}</span>
        <div className="flex items-center gap-2">
          <Button
            size="icon"
            className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-500"
            onClick={() => onGenerateLink(id)}
          >
            <QrCode className="h-5 w-5" />
          </Button>
          <Button
            className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-500"
            onClick={onViewDetails}
          >
            <Eye className="mr-2 h-4 w-4" /> Ver mais
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
};