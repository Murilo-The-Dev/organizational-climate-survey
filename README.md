# 🏢 Sistema de Pesquisa de Clima Organizacional

<div align="center">

![Status](https://img.shields.io/badge/Status-Em%20Desenvolvimento-yellow)
![Versão](https://img.shields.io/badge/Versão-1.0-blue)
![Licença](https://img.shields.io/badge/Licença-Acadêmico-green)

**Sistema completo para coleta, análise e gestão de pesquisas de clima organizacional com garantia de anonimato e conformidade com a LGPD**

</div>

---

📖 Sobre o Projeto
Este sistema foi desenvolvido como projeto de extensão curricular pelos alunos:

Murilo do Amaral Christofoletti (8204209) - Backend
Alexandre Ricardo Calore (8205280) - Frontend
Geovanni Adrian de Oliveira Muniz (8203566) - Database
Guilherme Rodrigues da Conceição (8183961) - Frontend

🎯 Objetivo
Oferecer uma solução abrangente para empresas realizarem pesquisas internas de clima organizacional, priorizando:

🔒 Anonimato completo dos respondentes
📊 Análises segmentadas por departamentos
📈 Comparações históricas de resultados
🔄 Automação de processos recorrentes
⚖️ Conformidade com LGPD e regulamentações


✨ Funcionalidades Principais
👨‍💼 Para Administradores (RH)

✅ Criação de pesquisas com formulários customizáveis
✅ Gestão de empresas e setores organizacionais
✅ Geração automática de links e QR Codes
✅ Agendamento de pesquisas recorrentes
✅ Dashboards interativos com métricas em tempo real
✅ Exportação de relatórios (Excel, PDF, CSV)
✅ Auditoria completa de ações no sistema

👥 Para Respondentes (Colaboradores)

✅ Acesso anônimo via link ou QR Code
✅ Interface responsiva para qualquer dispositivo
✅ Múltiplos tipos de pergunta (múltipla escolha, escala Likert, texto livre)
✅ Proteção contra múltiplas submissões
✅ Experiência intuitiva sem necessidade de cadastro

🏢 Para Empresas

✅ Análise segmentada por setores e equipes
✅ Comparações históricas de indicadores
✅ Insights acionáveis para tomada de decisão
✅ Conformidade total com LGPD

### Arquitetura de Sistema

```
Stack Tecnológica
Frontend        Backend         Database       
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ 
│  Next.js 14 │ │   Go 1.21+  │ │PostgreSQL 15│
│ TypeScript  │ │ Gorilla Mux │ │    pgx/v5   │ 
│  Tailwind   │ │     JWT     │ │  Migrations │ 
│   Recharts  │ │   Swagger   │ │   Indexes   │ 
└─────────────┘ └─────────────┘ └─────────────┘ 

```

## 🛡️ Segurança e Conformidade

### 🔒 Medidas de Segurança Implementadas
- **JWT Authentication** com refresh tokens
- **Rate Limiting** em endpoints sensíveis  
- **CORS** configurado adequadamente
- **SQL Injection** prevenção via ORM/prepared statements
- **XSS Protection** com sanitização de inputs
- **HTTPS** obrigatório em produção
- **Logs de auditoria** completos

### ⚖️ Conformidade LGPD
- **Anonimização completa** de respostas
- **Minimização de dados** - apenas dados necessários
- **Transparência** - usuários informados sobre tratamento
- **Direito de exclusão** - dados podem ser removidos
- **Logs de auditoria** para rastreabilidade
- **Criptografia** em dados sensíveis

---

### 📝 Padrões de Código

#### Commits
```bash
feat(auth): adiciona endpoint de refresh token
fix(dashboard): corrige cálculo de percentuais
docs(api): atualiza documentação dos endpoints
test(survey): adiciona testes para criação de pesquisa
```

---

## 📞 Suporte e Comunidade

### 🐛 Reportar Bugs
Encontrou um problema? [Abra uma issue](https://github.com/Murilo-The-Dev/sistema-clima-organizacional/issues) com:
- Descrição detalhada do problema
- Passos para reproduzir
- Screenshots (se aplicável)
- Ambiente (OS, browser, versões)

### 💡 Sugerir Melhorias
Tem uma ideia? [Abra uma feature request](https://github.com/Murilo-The-Dev/sistema-clima-organizacional/issues/new?template=feature_request.md)

### 📧 Contato da Equipe
- **Murilo Christofoletti** - [@murilo_christofoletti](https://github.com/Murilo-The-Dev)
- **Geovanni Muniz** - [@geovanni_adri](https://github.com/geovanniz) 
- **Guilherme Conceição** - [@rodriguesg.dev](https://github.com/rodriguesdev-ui)
- **Alexandre Calore** - [@alexandre_calore1](https://github.com/AlexandreCalore)

---

## 🏆 Reconhecimentos

### 📚 Tecnologias Utilizadas
- [Next.js](https://nextjs.org/) - Framework React
- [Golang](https://golang.org/) - Linguagem backend
- [Postgres](https://www.postgresql.org) - Banco de dados
- [Tailwind CSS](https://tailwindcss.com/) - Framework CSS
- [Recharts](https://recharts.org/) - Biblioteca de gráficos
- [JWT](https://jwt.io/) - Autenticação
- [Docker](https://docker.com/) - Containerização


### 🎓 Instituição
Projeto desenvolvido como **Extensão Curricular** do Centro Universitário Claretiano com foco em aplicação prática de conhecimentos acadêmicos em cenário real.

---

## 📄 Licença

Este projeto é um **trabalho acadêmico** desenvolvido para fins educacionais. 

Para uso comercial ou adaptações, entre em contato com a equipe de desenvolvimento.

---

<div align="center">

**⭐ Se este projeto foi útil, considere dar uma estrela!**

**🤝 Contribuições são sempre bem-vindas!**

---

Feito com ❤️ pela equipe de Extensão Curricular

</div>