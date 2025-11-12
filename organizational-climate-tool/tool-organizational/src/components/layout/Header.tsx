import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";

type HeaderProps = {
  onMenuClick: () => void;
};

const Header = ({ onMenuClick }: HeaderProps) => {
  return (
    <header className="sticky top-0 z-10 w-full h-16 bg-background/95 backdrop-blur-sm border-b">
      <div className="container flex items-center h-full max-w-7xl mx-auto px-4">
        <div className="md:hidden">
          <Button onClick={onMenuClick} variant="ghost" size="icon">
            <Menu className="h-6 w-6" />
            <span className="sr-only">Abrir Menu</span>
          </Button>
        </div>
      </div>
    </header>
  );
};

export default Header;
