
"use client";

import * as React from "react";
import { CardDescription, CardHeader } from "../ui/card";
import { CardContent } from "../ui/card";
import { Card } from "../ui/card";
import { Button } from "../ui/button";
import { Cog, ArrowDown } from "lucide-react";



export function FormsData() {

  return (
    <div className="w-full flex flex-col gap-4">
      <div className="w-full flex flex-row nowrap gap-2 items-center justify-center">
      <h1 className=" text-blue-600 text-2xl font-bold">Formulários<ArrowDown className="w-6 h-6 text-blue-600" /></h1>
      </div>
        <Card className="hover:shadow-lg hover:translate-y-[-5px] transition-all duration-300">
          <CardHeader>
            <CardDescription>Descrição do Formulário</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-row gap-2 items-center justify-between">
            <h1 className="text-lg font-bold">Nome do Formulário</h1>
            <p className="text-muted-foreground">Descrição do Formulário</p>
            <p className="text-muted-foreground">Data de criação: 18/09/2025</p>
              <Button className="bg-blue-600 w-fit text-white hover:bg-blue-500 cursor-pointer"><Cog className="w-4 h-4 text-white" />Ver mais </Button>
              </div>
          </CardContent>
        </Card>

        <Card className="hover:shadow-lg hover:translate-y-[-5px] transition-all duration-300">
          <CardHeader>
            <CardDescription>Descrição do Formulário</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-row gap-2 items-center justify-between">
            <h1 className="text-lg font-bold">Nome do Formulário</h1>
            <p className="text-muted-foreground">Descrição do Formulário</p>
            <p className="text-muted-foreground">Data de criação: 18/09/2025</p>
            <Button className="bg-blue-600 w-fit text-white hover:bg-blue-500 cursor-pointer"><Cog className="w-4 h-4 text-white" />Ver mais</Button>
            </div>
          </CardContent>
        </Card>

        <Card className="hover:shadow-lg hover:translate-y-[-5px] transition-all duration-300">
          <CardHeader>
            <CardDescription>Descrição do Formulário</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-row gap-2 items-center justify-between">
            <h1 className="text-lg font-bold">Nome do Formulário</h1>
            <p className="text-muted-foreground">Descrição do Formulário</p>
            <p className="text-muted-foreground">Data de criação: 18/09/2025</p>
            <Button className="bg-blue-600 w-fit text-white hover:bg-blue-500 cursor-pointer"><Cog className="w-4 h-4 text-white" />Ver mais</Button>
            </div>
          </CardContent>
        </Card>

        <Card className="hover:shadow-lg hover:translate-y-[-5px] transition-all duration-300">
          <CardHeader>
            <CardDescription>Descrição do Formulário</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-row gap-2 items-center justify-between">
            <h1 className="text-lg font-bold">Nome do Formulário</h1>
            <p className="text-muted-foreground">Descrição do Formulário</p>
            <p className="text-muted-foreground">Data de criação: 18/09/2025</p>
            <Button className="bg-blue-600 w-fit text-white hover:bg-blue-500 cursor-pointer"><Cog className="w-4 h-4 text-white" />Ver mais</Button>
            </div>
          </CardContent>
        </Card>

        <Card className="hover:shadow-lg hover:translate-y-[-5px] transition-all duration-300">
          <CardHeader>
            <CardDescription>Descrição do Formulário</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-row gap-2 items-center justify-between">
            <h1 className="text-lg font-bold">Nome do Formulário</h1>
            <p className="text-muted-foreground">Descrição do Formulário</p>
            <p className="text-muted-foreground">Data de criação: 18/09/2025</p>
            <Button className="bg-blue-600 w-fit text-white hover:bg-blue-500 cursor-pointer"><Cog className="w-4 h-4 text-white" />Ver mais</Button>
            </div>
          </CardContent>
        </Card>

        <Card className="hover:shadow-lg hover:translate-y-[-5px] transition-all duration-300">
          <CardHeader>
            <CardDescription>Descrição do Formulário</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-row gap-2 items-center justify-between">
            <h1 className="text-lg font-bold">Nome do Formulário</h1>
            <p className="text-muted-foreground">Descrição do Formulário</p>
            <p className="text-muted-foreground">Data de criação: 18/09/2025</p>
            <Button className="bg-blue-600 w-fit text-white hover:bg-blue-500 cursor-pointer"><Cog className="w-4 h-4 text-white" />Ver mais</Button>
            </div>
          </CardContent>
        </Card>
    </div>
  );
}