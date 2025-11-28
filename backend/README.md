# Organizational Climate Survey - Backend

ATMOS

Sistema de pesquisa de clima organizacional com arquitetura limpa e separação clara de responsabilidades.

## 📌 Funcionalidades Principais
- Autenticação de administradores via login e senha (JWT).
- Cadastro de empresas, setores e usuários administradores.
- Criação, edição e agendamento de pesquisas.
- Coleta de respostas anônimas.
- Dashboards analíticos e exportação de relatórios.
- Logs de auditoria para ações administrativas.

---

## 🚀 Tecnologias Utilizadas
- **Go** (linguagem principal).
- **Postgres** (banco de dados relacional).
- **godotenv** (carregar variáveis de ambiente).
- **jwt-go** (autenticação com JSON Web Tokens).

------

## 📁 Estrutura do Projeto

organizational-climate-survey/backend/
├── cmd/api/              # ✓ Entry point
├── config/               # ✓ Configurações
├── internal/
│   ├── application/      # ✓ Camada de aplicação
│   │   ├── dto/          # ✓ Data Transfer Objects
│   │   ├── handler/      # ✓ HTTP handlers
│   │   └── middleware/   # ✓ Middlewares específicos
│   ├── domain/           # ✓ Regras de negócio
│   │   ├── entity/       # ✓ Entidades
│   │   ├── repository/   # ✓ Interfaces
│   │   └── usecase/      # ✓ Casos de uso
│   └── infrastructure/   # ✓ Implementações externas
│       ├── auth/         # ✓ JWT, hash
│       ├── http/         # ✓ Servidor HTTP
│       └── postgres/     # ✓ Implementações repository
├── migrations/           # ✓ SQL migrations
├── pkg/
├── .env                              # Variáveis de ambiente
├── go.mod                            # Dependências do módulo Go
├── go.sum                            # Checksums das dependências
└── README.md                         # Documentação do projeto

## 🏗️ Arquitetura

### Clean Architecture

O projeto segue os princípios da Clean Architecture com separação clara de responsabilidades:

**1. Domain Layer (internal/domain/)**
- Contém as regras de negócio fundamentais
- Independente de frameworks e implementações externas
- Entities: Representação das entidades de negócio
- Repository Interfaces: Contratos para acesso a dados
- Use Cases: Orquestração de lógica de negócio

**2. Application Layer (internal/application/)**
- Camada de adaptação entre HTTP e domínio
- DTOs: Transformação de dados entre camadas
- Handlers: Processamento de requisições HTTP
- Middlewares: Interceptação de requisições

**3. Infrastructure Layer (internal/infrastructure/)**
- Implementações concretas de detalhes técnicos
- Database: Conexões e transações
- Auth: JWT, bcrypt, tokens
- HTTP: Servidor e configuração de rotas
- Postgres: Implementações SQL dos repositórios

**4. Package Layer (pkg/)**
- Utilitários reutilizáveis e independentes
- Validações, logging, helpers

### Fluxo de Requisição

HTTP Request
↓
Middleware (Auth, CORS, Logger)
↓
Handler (application/handler)
↓
DTO Validation
↓
Use Case (domain/usecase)
↓
Entity Business Logic (domain/entity)
↓
Repository Interface (domain/repository)
↓
Repository Implementation (infrastructure/postgres)
↓
Database

## 🔐 Segurança

- **Autenticação:** JWT com refresh tokens
- **Passwords:** Bcrypt com custo configurável
- **Validação:** Validação robusta de entrada com validator package
- **Auditoria:** Logs detalhados de todas as operações sensíveis
- **CORS:** Configuração restritiva para APIs

## 📊 Logging

Sistema de logging estruturado com:
- Níveis configuráveis (DEBUG, INFO, WARN, ERROR, FATAL)
- Context propagation para request tracing
- Fields injection para dados estruturados
- Caller information para debugging

## 🗄️ Banco de Dados

- PostgreSQL como SGBD principal
- Migrações versionadas com up/down
- Transações gerenciadas na camada de infrastructure
- Connection pooling configurável

## 🚀 Executando o Projeto
```bash
# Instalar dependências
go mod download

# Executar migrações
make migrate-up

# Iniciar servidor
go run cmd/api/main.go


📦 Dependências Principais

gorilla/mux - Roteamento HTTP
jackc/pgx/v5 - Driver PostgreSQL
golang-jwt/jwt/v5 - Autenticação JWT
golang.org/x/crypto - Bcrypt para senhas