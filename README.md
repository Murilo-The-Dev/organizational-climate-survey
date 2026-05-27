# 🏢 Sistema de Pesquisa de Clima Organizacional

<div align="center">

![Status](https://img.shields.io/badge/Status-Em%20Desenvolvimento-yellow)
![Versão](https://img.shields.io/badge/Versão-1.0-blue)
![Licença](https://img.shields.io/badge/Licença-Acadêmico-green)
![Go](https://img.shields.io/badge/Go-1.23.5-00ADD8?logo=go)
![Next.js](https://img.shields.io/badge/Next.js-15.5.2-black?logo=next.js)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?logo=postgresql)

**Sistema completo para coleta, análise e gestão de pesquisas de clima organizacional com garantia de anonimato e conformidade com a LGPD**

[📋 Sobre](#-sobre-o-projeto) • [🏗️ Arquitetura](#️-arquitetura-técnica) • [🚀 Instalação](#-instalação-e-configuração) • [📊 Roadmap](#-roadmap)

</div>

---

## 📖 Sobre o Projeto

Este sistema foi desenvolvido como **projeto de extensão curricular** pelos alunos:

| Nome | RA | Responsabilidade |
|------|----|-----------------|
| **Murilo do Amaral Christofoletti** | 8204209 | Backend |
| **Alexandre Ricardo Calore** | 8205280 | Frontend |
| **Geovanni Adrian de Oliveira Muniz** | 8203566 | Database |
| **Guilherme Rodrigues da Conceição** | 8183961 | Frontend |

### 🎯 Objetivo

Oferecer uma solução abrangente para empresas realizarem pesquisas internas de clima organizacional, priorizando:

- **🔒 Anonimato completo** dos respondentes via hash de IP/fingerprint
- **📊 Análises segmentadas** por departamentos e setores
- **📈 Comparações históricas** de resultados
- **🔄 Agendamento** de pesquisas recorrentes
- **⚖️ Conformidade** com LGPD e regulamentações brasileiras

---

## ✨ Funcionalidades Principais

### 👨‍💼 Para Administradores (Ferramenta)
- ✅ **Criação e gestão de pesquisas** com perguntas customizáveis
- ✅ **Gestão de empresas** e setores organizacionais
- ✅ **Gestão de usuários** administradores por empresa
- ✅ **Geração de links e QR Codes** para compartilhamento
- ✅ **Dashboards interativos** com múltiplos tipos de gráfico
- ✅ **Exportação de relatórios** em PDF
- ✅ **Logs de auditoria** completos de ações administrativas
- ✅ **Visualização de resultados** por pesquisa e setor

### 👥 Para Respondentes (Colaboradores)
- ✅ **Acesso anônimo** via link ou QR Code, sem necessidade de cadastro
- ✅ **Interface responsiva** para qualquer dispositivo
- ✅ **Proteção contra múltiplas submissões** via hash anônimo
- ✅ **Múltiplos tipos de pergunta** suportados

### 🌐 Landing Page
- ✅ **Página institucional** apresentando o sistema
- ✅ **Design responsivo** com Tailwind CSS

---

## 🏗️ Arquitetura Técnica

### Stack Tecnológica

```
Backend              Frontend (Ferramenta)   Frontend (Landing Page)   Database
┌──────────────────┐ ┌────────────────────┐  ┌──────────────────────┐  ┌──────────────┐
│  Go 1.23.5       │ │  Next.js 15.5.2    │  │  Next.js 15.5.2      │  │ PostgreSQL   │
│  gorilla/mux     │ │  React 19          │  │  React 19            │  │              │
│  golang-jwt/v5   │ │  TypeScript 5      │  │  TypeScript 5        │  │ Migrações    │
│  lib/pq          │ │  Tailwind CSS v4   │  │  Tailwind CSS v4     │  │ SQL puras    │
│  godotenv        │ │  Radix UI          │  │                      │  │ (psql)       │
│  google/uuid     │ │  shadcn/ui         │  │                      │  │              │
│  bcrypt          │ │  Recharts          │  │                      │  │              │
│  SQL puro        │ │  React Hook Form   │  │                      │  │              │
│                  │ │  Zod               │  │                      │  │              │
│                  │ │  TanStack Table    │  │                      │  │              │
│                  │ │  Axios + nookies   │  │                      │  │              │
│                  │ │  jsPDF + html2canvas│  │                      │  │              │
│                  │ │  react-qr-code     │  │                      │  │              │
│                  │ │  Sonner (toasts)   │  │                      │  │              │
│                  │ │  next-themes       │  │                      │  │              │
│                  │ │  lucide-react      │  │                      │  │              │
└──────────────────┘ └────────────────────┘  └──────────────────────┘  └──────────────┘
```

### Arquitetura do Backend — Clean Architecture

O backend segue os princípios da **Clean Architecture** com separação clara de responsabilidades:

```
┌──────────────────────────────────────────────────┐
│                  HTTP Request                    │
└────────────────────────┬─────────────────────────┘
                         ▼
┌──────────────────────────────────────────────────┐
│         Middleware (Auth JWT, CORS, Logger)       │
│              internal/application/middleware      │
└────────────────────────┬─────────────────────────┘
                         ▼
┌──────────────────────────────────────────────────┐
│            Handler + DTO Validation               │
│              internal/application/handler         │
│              internal/application/dto             │
└────────────────────────┬─────────────────────────┘
                         ▼
┌──────────────────────────────────────────────────┐
│              Use Case (regras de negócio)         │
│              internal/domain/usecase              │
└────────────────────────┬─────────────────────────┘
                         ▼
┌──────────────────────────────────────────────────┐
│      Repository Interface  ←→  Entity             │
│      internal/domain/repository                   │
│      internal/domain/entity                       │
└────────────────────────┬─────────────────────────┘
                         ▼
┌──────────────────────────────────────────────────┐
│        Repository Implementation (SQL puro)       │
│        internal/infrastructure/postgres           │
└────────────────────────┬─────────────────────────┘
                         ▼
┌──────────────────────────────────────────────────┐
│                  PostgreSQL                       │
└──────────────────────────────────────────────────┘
```

### Entidades do Domínio

| Entidade | Descrição |
|----------|-----------|
| `Empresa` | Empresa cadastrada no sistema |
| `Setor` | Setor/departamento dentro de uma empresa |
| `UsuarioAdministrador` | Admin vinculado a uma empresa |
| `Pesquisa` | Pesquisa de clima com período e status |
| `Pergunta` | Pergunta dentro de uma pesquisa |
| `SubmissaoPesquisa` | Registro anônimo de submissão (anti-duplicate) |
| `Resposta` | Resposta individual a uma pergunta |
| `LogAuditoria` | Log de ações administrativas |
| `Dashboard` | Agregações e métricas calculadas |

---

## 📁 Estrutura do Repositório

```
organizational-climate-survey/
│
├── backend/                              # API REST em Go
│   ├── cmd/api/
│   │   └── main.go                       # Entry point
│   ├── config/
│   │   └── config.go                     # Carregamento de variáveis de ambiente
│   ├── internal/
│   │   ├── application/
│   │   │   ├── dto/                      # Data Transfer Objects
│   │   │   ├── handler/                  # HTTP handlers (um por entidade)
│   │   │   └── middleware/               # Auth JWT, CORS, Logger
│   │   ├── domain/
│   │   │   ├── entity/                   # Entidades de negócio
│   │   │   ├── repository/               # Interfaces de repositório
│   │   │   └── usecase/                  # Casos de uso
│   │   └── infrastructure/
│   │       ├── auth/                     # JWT e bcrypt
│   │       ├── http/                     # Servidor HTTP e roteamento (gorilla/mux)
│   │       └── postgres/                 # Implementações SQL puras dos repositórios
│   ├── migrations/                       # Arquivos SQL versionados
│   │   ├── 000_setup_migrations.sql      # Tabela de controle de migrações
│   │   ├── 001_initial_state.sql         # Schema principal
│   │   ├── 002_adding_procedures_and_triggers.sql
│   │   ├── 003_adding_views_for_reports.sql
│   │   ├── 004_add_submissao_pesquisa.sql
│   │   ├── 005_alter_resposta_add_submissao.sql
│   │   ├── 006_feat_add_submissao_signals.sql
│   │   └── 007_fix_resposta.sql
│   ├── pkg/
│   │   ├── crypto/                       # Serviço de criptografia e hash
│   │   ├── logger/                       # Logger estruturado
│   │   └── validator/                    # Utilitários de validação
│   ├── migrate.sh                        # Script de gerenciamento de migrações
│   ├── .env                              # Variáveis de ambiente (não versionado)
│   ├── go.mod
│   └── go.sum
│
├── organizational-climate-tool/          # Frontend — Ferramenta Administrativa
│   └── tool-organizational/
│       └── src/
│           ├── app/
│           │   ├── (app)/                # Rotas protegidas (requer autenticação)
│           │   │   ├── dashboard/        # Dashboard principal
│           │   │   ├── pesquisas/        # Listagem e gestão de pesquisas
│           │   │   ├── empresas/         # Gestão de empresas
│           │   │   ├── usuarios/         # Gestão de usuários administradores
│           │   │   ├── resultados/       # Resultados e analytics
│           │   │   ├── auditoria/        # Log de auditoria
│           │   │   └── configuracoes/    # Configurações do sistema
│           │   ├── login/                # Página pública de login
│           │   └── pesquisas/[id]/       # Página pública de resposta à pesquisa
│           │       └── responder/
│           ├── components/
│           │   ├── dashboard/            # StatCard, DataTable, ResultsDataTable
│           │   │   └── charts/           # BarChart, LineChart, PieChart, RadialChart,
│           │   │                         # StackedChart, EngagementChart, NpsBreakdown, etc.
│           │   ├── forms/                # Formulários reutilizáveis
│           │   ├── layout/               # Componentes de layout (sidebar, header)
│           │   ├── modals/               # SurveyDetailsModal, SurveyLinkModal (QR Code)
│           │   ├── pesquisas/            # SurveyCard, SurveyOverviewTab, SurveyQuestionsTab
│           │   └── ui/                   # Componentes base (shadcn/ui)
│           ├── context/
│           │   └── AuthContext.tsx       # Contexto de autenticação (JWT via cookies)
│           └── lib/
│               ├── api.ts                # Cliente Axios com interceptor de token
│               └── utils.ts
│
└── organizational-climate-lp/           # Frontend — Landing Page
    └── lp-organizational/
        └── src/
            └── app/                      # Página institucional (Next.js App Router)
```

---

## 🚀 Instalação e Configuração

### Pré-requisitos

- **Go** 1.23+ ([Download](https://golang.org/dl/))
- **Node.js** 18+ e **npm** ([Download](https://nodejs.org/))
- **PostgreSQL** 14+ ([Download](https://www.postgresql.org/download/))
- **Git** ([Download](https://git-scm.com/))
- **psql** (cliente PostgreSQL, incluso na instalação padrão do PostgreSQL)

### Clone o Repositório

```bash
git clone https://github.com/Murilo-The-Dev/organizational-climate-survey.git
cd organizational-climate-survey
```

---

### 🗄️ 1. Configuração do Banco de Dados (PostgreSQL)

```bash
# Acessar o PostgreSQL
psql -U postgres

# Criar o banco de dados
CREATE DATABASE organizational_climate;
\q

# Configurar o schema de controle de migrações
cd backend
psql -U postgres -d organizational_climate -f migrations/000_setup_migrations.sql

# Aplicar todas as migrações pendentes via script
chmod +x migrate.sh
./migrate.sh aplicar
```

#### Comandos do `migrate.sh`

```bash
./migrate.sh aplicar            # Aplica migrações pendentes
./migrate.sh status             # Exibe o status de cada migration
./migrate.sh nova 'descrição'   # Cria um novo arquivo de migration
```

---

### 🔧 2. Backend (Go)

```bash
cd backend

# Instalar dependências
go mod download

# Configurar variáveis de ambiente
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

**Variáveis de ambiente (`.env`):**

```env
# Aplicação
APP_NAME=organizational-climate-survey
APP_PORT=8080
APP_ENV=development

# Banco de dados (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=sua_senha_aqui
DB_NAME=organizational_climate
DB_SSLMODE=disable

# Autenticação JWT
JWT_SECRET=seu_jwt_secret_muito_seguro_aqui

# Segurança — salt para hash anônimo de submissões
HASH_SALT=seu_hash_salt_aqui

# Logs
LOG_LEVEL=debug
```

```bash
# Iniciar o servidor em modo desenvolvimento
go run cmd/api/main.go

# Build para produção
go build -o bin/api cmd/api/main.go
./bin/api
```

O servidor estará disponível em:
- **API Base URL:** `http://localhost:8080/api/v1`
- **Health Check:** `http://localhost:8080/health`

---

### 🎨 3. Frontend — Ferramenta Administrativa

```bash
cd organizational-climate-tool/tool-organizational

# Instalar dependências
npm install

# Configurar variável de ambiente
echo "NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1" > .env.local

# Iniciar em modo desenvolvimento
npm run dev
```

A ferramenta estará disponível em `http://localhost:3000`.

**Scripts disponíveis:**

```bash
npm run dev      # Servidor de desenvolvimento
npm run build    # Build de produção
npm run start    # Inicia o build de produção
npm run lint     # Verificação de lint
```

---

### 🌐 4. Frontend — Landing Page

```bash
cd organizational-climate-lp/lp-organizational

# Instalar dependências
npm install

# Iniciar em modo desenvolvimento
npm run dev
```

A landing page estará disponível em `http://localhost:3001` (ou outra porta disponível).

---

## 🔌 API — Principais Endpoints

Todos os endpoints protegidos requerem o header `Authorization: Bearer <token>`.

<details>
<summary><strong>🔐 Autenticação</strong></summary>

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@empresa.com",
  "password": "senha123"
}
```

```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

```http
POST /api/v1/bootstrap
Content-Type: application/json
# Cria empresa e primeiro administrador (setup inicial)
```

</details>

<details>
<summary><strong>🏢 Empresas e Setores</strong></summary>

```http
GET    /api/v1/empresas
POST   /api/v1/empresas
GET    /api/v1/empresas/{id}
PUT    /api/v1/empresas/{id}
DELETE /api/v1/empresas/{id}

GET    /api/v1/setores
POST   /api/v1/setores
GET    /api/v1/setores/{id}
PUT    /api/v1/setores/{id}
DELETE /api/v1/setores/{id}
```

</details>

<details>
<summary><strong>👤 Usuários Administradores</strong></summary>

```http
GET    /api/v1/usuarios
POST   /api/v1/usuarios
GET    /api/v1/usuarios/{id}
PUT    /api/v1/usuarios/{id}
DELETE /api/v1/usuarios/{id}
```

</details>

<details>
<summary><strong>📋 Pesquisas e Perguntas</strong></summary>

```http
GET    /api/v1/pesquisas
POST   /api/v1/pesquisas
GET    /api/v1/pesquisas/{id}
PUT    /api/v1/pesquisas/{id}
DELETE /api/v1/pesquisas/{id}

GET    /api/v1/pesquisas/{id}/perguntas
POST   /api/v1/perguntas
PUT    /api/v1/perguntas/{id}
DELETE /api/v1/perguntas/{id}
```

</details>

<details>
<summary><strong>📝 Submissão Anônima de Respostas</strong></summary>

```http
# Rotas públicas — não requerem autenticação
GET  /api/v1/public/pesquisas/{id}
POST /api/v1/public/pesquisas/{id}/submissoes
POST /api/v1/public/pesquisas/{id}/respostas
```

> A proteção contra múltiplas submissões é feita via hash anônimo do IP + fingerprint do respondente, sem armazenar dados identificáveis.

</details>

<details>
<summary><strong>📊 Dashboard e Auditoria</strong></summary>

```http
GET /api/v1/dashboard
GET /api/v1/dashboard/pesquisas/{id}

GET /api/v1/logs
GET /api/v1/logs/{id}
```

</details>

---

## 🛡️ Segurança e Conformidade

### 🔒 Medidas de Segurança

- **JWT Authentication** com `golang-jwt/v5` para sessões de administradores
- **Bcrypt** para hash de senhas (`golang.org/x/crypto`)
- **Hash anônimo** de IP/fingerprint para submissões (sem armazenar dados pessoais)
- **CORS** configurado via middleware
- **SQL Injection** prevenida via prepared statements (SQL puro com `database/sql`)
- **Logs de auditoria** completos de todas as operações administrativas

### ⚖️ Conformidade LGPD

- **Anonimização completa** — respostas não são vinculadas a identificadores pessoais
- **Minimização de dados** — apenas dados estritamente necessários são coletados
- **Hash unidirecional** — impossível reverter IP/fingerprint ao dado original
- **Direito de exclusão** — dados podem ser removidos pelo administrador
- **Transparência** — respondentes informados sobre o tratamento de dados

---

## 📊 Roadmap

#### ✅ Fase 1 — Setup e Infraestrutura
- [x] Repositório Git configurado
- [x] Arquitetura Clean Architecture definida
- [x] Schema do banco de dados e migrações
- [x] Script de migrações (`migrate.sh`)

#### ✅ Fase 2 — Autenticação e Gestão
- [x] Sistema de login JWT
- [x] Middleware de autenticação
- [x] CRUD de empresas e setores
- [x] CRUD de usuários administradores
- [x] Bootstrap (setup inicial do sistema)

#### ✅ Fase 3 — Pesquisas
- [x] Criação e edição de pesquisas
- [x] Sistema de perguntas dinâmico
- [x] Geração de links e QR Codes para compartilhamento
- [x] Interface pública de resposta anônima
- [x] Proteção contra múltiplas submissões

#### 🔄 Fase 4 — Analytics e Relatórios
- [x] Dashboard com gráficos interativos (Bar, Line, Pie, Radial, Stacked, NPS)
- [x] Exportação de relatórios em PDF (jsPDF + html2canvas)
- [ ] Análises históricas comparativas
- [ ] Segmentação avançada por setor

#### ⏳ Fase 5 — Testes e Deploy
- [ ] Testes unitários (backend e frontend)
- [ ] CI/CD
- [ ] Deploy em produção
- [ ] Documentação final da API

---

## 📝 Padrões de Código

### Commits (Conventional Commits)

```bash
feat(pesquisas): adiciona endpoint de criação de pesquisa
fix(dashboard): corrige cálculo de percentuais por setor
docs(readme): atualiza stack tecnológica
test(auth): adiciona testes para middleware JWT
refactor(repos): migra queries para prepared statements
```

### Code Review Checklist

- [ ] Código limpo e bem comentado
- [ ] DTOs validados na camada de application
- [ ] Erros tratados adequadamente
- [ ] Logs de auditoria adicionados onde necessário
- [ ] Segurança e anonimato preservados

---

## 📞 Equipe

| Membro | GitHub |
|--------|--------|
| **Murilo Christofoletti** | [@Murilo-The-Dev](https://github.com/Murilo-The-Dev) |
| **Geovanni Muniz** | [@geovanniz](https://github.com/geovanniz) |
| **Guilherme Conceição** | [@rodriguesdev-ui](https://github.com/rodriguesdev-ui) |
| **Alexandre Calore** | [@AlexandreCalore](https://github.com/AlexandreCalore) |

### 🐛 Reportar Bugs

Encontrou um problema? [Abra uma issue](https://github.com/Murilo-The-Dev/organizational-climate-survey/issues) com:
- Descrição detalhada do problema
- Passos para reproduzir
- Screenshots (se aplicável)
- Ambiente (OS, browser, versões)

---

## 🏆 Tecnologias Utilizadas

### Backend
- [Go](https://golang.org/) — Linguagem principal
- [gorilla/mux](https://github.com/gorilla/mux) — Roteamento HTTP
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) — Autenticação JWT
- [lib/pq](https://github.com/lib/pq) — Driver PostgreSQL
- [google/uuid](https://github.com/google/uuid) — Geração de UUIDs
- [godotenv](https://github.com/joho/godotenv) — Carregamento de `.env`
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) — Bcrypt

### Frontend (Ferramenta Administrativa)
- [Next.js 15](https://nextjs.org/) — Framework React
- [React 19](https://react.dev/) — Biblioteca UI
- [TypeScript 5](https://www.typescriptlang.org/) — Tipagem estática
- [Tailwind CSS v4](https://tailwindcss.com/) — Estilização
- [Radix UI](https://www.radix-ui.com/) — Componentes acessíveis
- [shadcn/ui](https://ui.shadcn.com/) — Sistema de componentes
- [Recharts](https://recharts.org/) — Gráficos
- [TanStack Table](https://tanstack.com/table) — Tabelas avançadas
- [React Hook Form](https://react-hook-form.com/) + [Zod](https://zod.dev/) — Formulários e validação
- [Axios](https://axios-http.com/) + [nookies](https://github.com/matipan/nookies) — HTTP e cookies
- [jsPDF](https://github.com/parallax/jsPDF) + [html2canvas](https://html2canvas.hertzen.com/) — Exportação PDF
- [react-qr-code](https://github.com/rosskhanas/react-qr-code) — Geração de QR Codes
- [Sonner](https://sonner.emilkowal.ski/) — Notificações toast
- [next-themes](https://github.com/pacocoursey/next-themes) — Tema claro/escuro
- [lucide-react](https://lucide.dev/) — Ícones
- [date-fns](https://date-fns.org/) — Manipulação de datas

### Frontend (Landing Page)
- [Next.js 15](https://nextjs.org/) — Framework React
- [Tailwind CSS v4](https://tailwindcss.com/) — Estilização

### Banco de Dados
- [PostgreSQL](https://www.postgresql.org/) — SGBD relacional

---

## 🎓 Instituição

Projeto desenvolvido como **Extensão Curricular** do **Centro Universitário Claretiano** com foco em aplicação prática de conhecimentos acadêmicos em cenário real.

---

## 📄 Licença

Este projeto é um **trabalho acadêmico** desenvolvido para fins educacionais.

Para uso comercial ou adaptações, entre em contato com a equipe de desenvolvimento.

---

<div align="center">

**⭐ Se este projeto foi útil, considere dar uma estrela!**

---

Feito com ❤️ pela equipe de Extensão Curricular — Claretiano

</div>
