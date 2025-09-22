// src/components/dashboard/SurveyCard.tsx

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
import { Eye } from "lucide-react";

type SurveyCardProps = {
  title: string;
  description: string;
  tag: string;
  creationDate: string;
};

export const SurveyCard = ({
  title,
  description,
  tag,
  creationDate,
}: SurveyCardProps) => {
  return (
    <Card className="hover:shadow-lg hover:border-blue-600 hover:translate-y-[-5px] transition-all duration-500">
      <CardHeader>
        <Badge variant="secondary" className="w-fit">
          {tag}
        </Badge>
        <CardTitle className="pt-2">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground line-clamp-2">
          {description}
        </p>
      </CardContent>
      <CardFooter className="flex justify-between items-center text-sm text-muted-foreground">
        <span>Criado em: {creationDate}</span>
        <Button className="cursor-pointer bg-blue-600 text-white hover:bg-blue-500 hover:text-white transition-all duration-500">
          <Eye className="mr-2 h-4 w-4" /> Ver mais
        </Button>
      </CardFooter>
    </Card>
  );
};
