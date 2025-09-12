"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  LayoutDashboard,
  NotebookText,
  BarChart3,
  Settings,
  LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import LogoAtmos from "@/public/images/logoAtmos.svg";
import Image from "next/image";

type SidebarProps = {
  isOpen: boolean;
};

interface NavLink {
  href: string;
  label: string;
  icon: LucideIcon;
}

const navLinks: NavLink[] = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/pesquisas", label: "Pesquisas", icon: NotebookText },
  { href: "/resultados", label: "Resultados", icon: BarChart3 },
  { href: "/configuracoes", label: "Configurações", icon: Settings },
];

const Sidebar = ({ isOpen }: SidebarProps) => {
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        "bg-background rounded-r-4xl rounded-l-4xl text-foreground w-64 h-screen p-4 border-r flex flex-col fixed top-0 left-0 z-40 transition-transform duration-300 ease-in-out",
        !isOpen && "-translate-x-full",
        "md:translate-x-0 md:sticky"
      )}
    >
      <div className="flex justify-center items-center mb-10 px-2">
        <h1 className="text-xl font-bold">
          <Image src={LogoAtmos} alt="Logo Atmos" width={50} height={50} />
        </h1>
      </div>

      <nav className="flex flex-col gap-2">
        {navLinks.map((link) => {
          const isActive = pathname === link.href;
          const Icon = link.icon;

          return (
            <Link href={link.href} key={link.label}>
              <Button
                variant={isActive ? "secondary" : "ghost"}
                className={cn(
                  "w-full justify-start cursor-pointer",
                  !isActive && "hover:bg-blue-600 hover:text-white"
                )}
              >
                <Icon className="h-4 w-4 mr-2" />
                {link.label}
              </Button>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
};

export default Sidebar;
