// src/components/reportMockData.ts

export const reportData = {
  resumo: {
    pontuacaoGeral: 82,
    participantes: 152,
    dataAvaliacao: "2025-03-28",
    proximaMedicao: "2025-09-28",
  },
  infoEmpresa: {
    razaoSocial: "Soluções Inovadoras LTDA",
    cnpj: "12.345.678/0001-99",
    tipo: "Tecnologia",
    setores: ["Produção", "ADM", "SGQ"],
    funcionarios: 180,
  },
  indicadores: {
    categorias: [
      { nome: "Engajamento", pontuacao: 88, classificacao: "Bom" },
      { nome: "Liderança", pontuacao: 75, classificacao: "Bom" },
      { nome: "Bem-estar", pontuacao: 65, classificacao: "Médio" },
      { nome: "Comunicação", pontuacao: 55, classificacao: "Médio" },
      { nome: "Segurança Psicológica", pontuacao: 92, classificacao: "Bom" },
      { nome: "Crescimento", pontuacao: 78, classificacao: "Bom" },
      { nome: "Reconhecimento", pontuacao: 45, classificacao: "Ruim" },
    ],
    satisfacao: [
      { nome: "Satisfação Geral", pontuacao: 85, classificacao: "Bom" },
      { nome: "e-NPS", pontuacao: 70, classificacao: "Bom" },
      { nome: "Intenção de Rotatividade", pontuacao: 30, classificacao: "Bom" },
      { nome: "Absenteísmo", pontuacao: 20, classificacao: "Bom" },
    ],
  },
  evolucao: [
    { data: "Q1 2024", Liderança: 70, "Bem-estar": 60, Comunicação: 50 },
    { data: "Q2 2024", Liderança: 72, "Bem-estar": 65, Comunicação: 52 },
    { data: "Q3 2024", Liderança: 68, "Bem-estar": 63, Comunicação: 58 },
    { data: "Q4 2024", Liderança: 78, "Bem-estar": 70, Comunicação: 65 },
    { data: "Q1 2025", Liderança: 75, "Bem-estar": 65, Comunicação: 55 },
  ],
  desempenhoDepartamentos: [
    {
      departamento: "ADM",
      Engajamento: 85,
      Liderança: 70,
      "Bem-estar": 60,
      Comunicação: 50,
      "Segurança Psicológica": 90,
    },
    {
      departamento: "Produção",
      Engajamento: 90,
      Liderança: 80,
      "Bem-estar": 70,
      Comunicação: 60,
      "Segurança Psicológica": 95,
    },
    {
      departamento: "SGQ",
      Engajamento: 80,
      Liderança: 75,
      "Bem-estar": 65,
      Comunicação: 55,
      "Segurança Psicológica": 88,
    },
  ],
  detalhamentoDepartamentos: [
    { id: 1, setor: "Produção", pontuacaoGeral: 88, Engajamento: 90, Liderança: 80, "Bem-estar": 70, Comunicação: 60 },
    { id: 2, setor: "ADM", pontuacaoGeral: 75, Engajamento: 85, Liderança: 70, "Bem-estar": 60, Comunicação: 50 },
    { id: 3, setor: "SGQ", pontuacaoGeral: 82, Engajamento: 80, Liderança: 75, "Bem-estar": 65, Comunicação: 55 },
  ],
  matrizRisco: [
    { categoria: "Reconhecimento", risco: "Crítico" },
    { categoria: "Comunicação", risco: "Ruim" },
    { categoria: "Bem-estar", risco: "Médio" },
    { categoria: "Liderança", risco: "Bom" },
    { categoria: "Crescimento", risco: "Bom" },
    { categoria: "Engajamento", risco: "Excelente" },
    { categoria: "Segurança Psicológica", risco: "Excelente" },
  ],
  mapaArvore: [
    { name: "Segurança Psicológica", size: 92 },
    { name: "Engajamento", size: 88 },
    { name: "Crescimento", size: 78 },
    { name: "Liderança", size: 75 },
    { name: "Bem-estar", size: 65 },
    { name: "Comunicação", size: 55 },
    { name: "Reconhecimento", size: 45 },
  ],
  historicoIntervencoes: [
    { data: "2024-07-15", acao: "Workshop de comunicação não-violenta para líderes.", responsavel: "RH" },
    { data: "2024-10-20", acao: "Revisão do plano de carreira e cargos.", responsavel: "Diretoria" },
    { data: "2025-01-10", acao: "Implementação de happy hour mensal.", responsavel: "RH" },
  ],
  recomendacoes: [
    { area: "Reconhecimento", risco: "Crítico", plano: "Criar um programa de reconhecimento formal (funcionário do mês, bônus por performance)." },
    { area: "Comunicação", risco: "Ruim", plano: "Implementar reuniões de alinhamento interdepartamental semanais." },
    { area: "Bem-estar", risco: "Médio", plano: "Oferecer subsídio para atividades físicas e de saúde mental." },
  ]
};