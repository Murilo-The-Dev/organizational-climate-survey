'use client';

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
//import { Copy, QrCode } from "lucide-react";
import QRCode from "react-qr-code";
import { toast } from "sonner";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import { ArrowUpDown, MoreHorizontal, Copy } from "lucide-react";

interface SurveyLinkModalProps {
  isOpen: boolean;
  onClose: () => void;
  surveyId: string;
}

export function SurveyLinkModal({ isOpen, onClose, surveyId }: SurveyLinkModalProps) {
  const surveyLink = `${window.location.origin}/pesquisas/${surveyId}/responder`;

  const copyToClipboard = () => {
    navigator.clipboard.writeText(surveyLink);
    toast.success("Link da pesquisa copiado para a área de transferência!");
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Link da Pesquisa</DialogTitle>
          <DialogDescription>
            Compartilhe este link para que os participantes possam responder à pesquisa.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="survey-link">Link</Label>
            <div className="flex space-x-2">
              <Input id="survey-link" value={surveyLink} readOnly />
              <Button type="button" size="sm" onClick={copyToClipboard}>
                <Copy className="h-4 w-4" />
                <span className="sr-only">Copiar</span>
              </Button>
            </div>
          </div>
          <div className="flex flex-col items-center justify-center p-4 border rounded-md">
            <Label className="mb-2">QR Code</Label>
            <QRCode value={surveyLink} size={180} level="H" />
            <p className="text-sm text-muted-foreground mt-2">Escaneie para responder</p>
          </div>
        </div>
        <DialogFooter>
          <Button type="button" variant="secondary" onClick={onClose}>Fechar</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

