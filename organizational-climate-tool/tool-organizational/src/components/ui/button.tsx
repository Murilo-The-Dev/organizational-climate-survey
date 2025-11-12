import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:ring-[3px] aria-invalid:border-destructive",
  {
    variants: {
      variant: {
        default:
          "bg-[#155DFC] text-white shadow-xs hover:bg-[#0D4AD9] focus-visible:border-[#155DFC] focus-visible:ring-[#155DFC]/30",
        destructive:
          "bg-destructive text-white shadow-xs hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60 focus-visible:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40",
        outline:
          "border border-[#155DFC]/30 bg-background shadow-xs hover:bg-[#155DFC]/10 hover:text-[#155DFC] hover:border-[#155DFC] dark:bg-input/30 dark:border-[#155DFC]/20 dark:hover:bg-[#155DFC]/20 focus-visible:border-[#155DFC] focus-visible:ring-[#155DFC]/20",
        secondary:
          "bg-[#E6F0FF] text-[#0D4AD9] shadow-xs hover:bg-[#CCE0FF] focus-visible:border-[#155DFC] focus-visible:ring-[#155DFC]/20 dark:bg-[#155DFC]/20 dark:text-[#6BA3FF] dark:hover:bg-[#155DFC]/30",
        ghost:
          "hover:bg-[#155DFC]/10 hover:text-[#155DFC] dark:hover:bg-[#155DFC]/20 focus-visible:ring-[#155DFC]/20",
        link: 
          "text-[#155DFC] underline-offset-4 hover:underline hover:text-[#0D4AD9] focus-visible:ring-[#155DFC]/20",
      },
      size: {
        default: "h-9 px-4 py-2 has-[>svg]:px-3",
        sm: "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5",
        lg: "h-10 rounded-md px-6 has-[>svg]:px-4",
        icon: "size-9",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot : "button"

  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

export { Button, buttonVariants }