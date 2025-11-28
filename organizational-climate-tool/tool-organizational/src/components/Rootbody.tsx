// RootBody.tsx
"use client";

export default function RootBody({
  children,
  inter,
  interMono,
}: {
  children: React.ReactNode;
  inter: any;
  interMono: any;
}) {
  return (
    <body className={`${inter.variable} ${interMono.variable} antialiased bg-zinc-50`}>
      {children}
    </body>
  );
}
