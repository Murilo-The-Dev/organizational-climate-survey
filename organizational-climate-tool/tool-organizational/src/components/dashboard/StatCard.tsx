// src/components/dashboard/StatCard.tsx

"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Users, ClipboardCheck, Smile, Star } from "lucide-react";
import React from "react";

// Este objeto mapeia um nome (string) que vamos passar para o componente do ícone real.
// É uma forma limpa de deixar o componente decidir qual ícone renderizar.
const iconMap = {
  clipboardCheck: ClipboardCheck,
  users: Users,
  smile: Smile,
  star: Star,
};

// Aqui, garantimos que a prop 'iconName' só possa ser uma das chaves do nosso iconMap.
// Isso nos dá autocompletar e segurança de tipo!
type IconName = keyof typeof iconMap;

type StatCardProps = {
  title: string;
  value: string;
  iconName: IconName;
};

const StatCard = ({ title, value, iconName }: StatCardProps) => {
  const Icon = iconMap[iconName];

  if (!Icon) {
    return null;
  }

  return (
    <Card className="hover:translate-y-[-5px] transition-all duration-300">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle>{title}</CardTitle>
        <Icon className="h-6 w-6  text-blue-600" />
      </CardHeader>
      <CardContent>
        <div>{value}</div>
      </CardContent>
    </Card>
  );
};

export default StatCard;
