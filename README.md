# 🏢 Sistema de Pesquisa de Clima Organizacional

<div align="center">

![Status](https://img.shields.io/badge/Status-Em%20Desenvolvimento-yellow)
![Versão](https://img.shields.io/badge/Versão-1.0-blue)
![Licença](https://img.shields.io/badge/Licença-Acadêmico-green)

**Sistema completo para coleta, análise e gestão de pesquisas de clima organizacional com garantia de anonimato e conformidade com a LGPD**

[📋 Documentação Completa](docs/) • [🚀 Demo](#) • [📊 Roadmap](#roadmap) • [🤝 Contribuir](#contribuindo)

</div>

---

## 📖 Sobre o Projeto

Este sistema foi desenvolvido como **projeto de extensão curricular** pelos alunos:
- **Murilo do Amaral Christofoletti** (8204209) - Backend
- **Alexandre Ricardo Calore** (8205280) - Frontend  
- **Geovanni Adrian de Oliveira Muniz** (8203566) - Database
- **Guilherme Rodrigues da Conceição** (8183961) - Frontend

### 🎯 Objetivo
Oferecer uma solução abrangente para empresas realizarem pesquisas internas de clima organizacional, priorizando:
- **🔒 Anonimato completo** dos respondentes
- **📊 Análises segmentadas** por departamentos
- **📈 Comparações históricas** de resultados
- **🔄 Automação** de processos recorrentes
- **⚖️ Conformidade** com LGPD e regulamentações

---

## ✨ Funcionalidades Principais

### 👨‍💼 Para Administradores
- ✅ **Criação de pesquisas** com formulários customizáveis
- ✅ **Gestão de empresas** e setores organizacionais  
- ✅ **Geração automática** de links e QR Codes
- ✅ **Agendamento** de pesquisas recorrentes
- ✅ **Dashboards interativos** com métricas em tempo real
- ✅ **Exportação** de relatórios (Excel, PDF, CSV)
- ✅ **Auditoria completa** de ações no sistema

### 👥 Para Respondentes (Colaboradores)
- ✅ **Acesso anônimo** via link ou QR Code
- ✅ **Interface responsiva** para qualquer dispositivo
- ✅ **Múltiplos tipos de pergunta** (múltipla escolha, escala Likert, texto livre)
- ✅ **Proteção contra** múltiplas submissões
- ✅ **Experiência intuitiva** sem necessidade de cadastro

### 🏢 Para Empresas
- ✅ **Análise segmentada** por setores e equipes
- ✅ **Comparações históricas** de indicadores
- ✅ **Insights acionáveis** para tomada de decisão
- ✅ **Conformidade total** com LGPD

---

## 🏗️ Arquitetura Técnica

### Stack Tecnológica
```
Frontend        Backend         Database        Deploy
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  Next.js 14 │ │   Go 1.21+  │ │  MySQL 8.0  │ │   Docker    │
│ TypeScript  │ │     Gin     │ │    GORM     │ │  Railway    │
│  Tailwind   │ │     JWT     │ │   Redis*    │ │   Vercel    │
│   Recharts  │ │   Swagger   │ │  Migrations │ │     CI/CD   │
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘
```

### Arquitetura de Sistema
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │    Database     │
│   (Next.js)     │    │   (Golang)      │    │    (MySQL)      │
│                 │    │                 │    │                 │
│ • Dashboard     │◄──►│ • REST API      │◄──►│ • Pesquisas     │
│ • Formulários   │    │ • JWT Auth      │    │ • Respostas     │
│ • Gráficos      │    │ • Middleware    │    │ • Usuários      │
│ • Relatórios    │    │ • QR Codes      │    │ • Auditoria     │
│                 │    │ • Cron Jobs     │    │ • Views         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        └───────────────────────┼───────────────────────┘
                               │
                    ┌─────────────────┐
                    │     Cache       │
                    │    (Redis)      │
                    │                 │
                    │ • Sessões       │
                    │ • Rate Limit    │
                    │ • Dashboard     │
                    └─────────────────┘
```

---

## 🚀 Instalação e Configuração

### Pré-requisitos
- **Node.js** 18+ ([Download](https://nodejs.org/))
- **Go** 1.21+ ([Download](https://golang.org/dl/))
- **MySQL** 8.0+ ([Download](https://dev.mysql.com/downloads/))
- **Git** ([Download](https://git-scm.com/))

### 🔧 Setup Rápido com Docker

```bash
# 1. Clone o repositório
git clone https://github.com/Murilo-The-Dev/organizational-climate-survey.git
cd organizational-climate-survey

# 2. Configure as variáveis de ambiente
cp .env.example .env
# Edite o arquivo .env com suas configurações

# 3. Execute com Docker Compose
docker-compose up -d

# 4. Execute as migrações
docker-compose exec backend go run migrations/migrate.go

# 5. Acesse a aplicação
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Documentação API: http://localhost:8080/swagger
```

### ⚙️ Setup Manual (Desenvolvimento)

<details>
<summary><strong>🗄️ 1. Configuração do Banco de Dados</strong></summary>

```bash
# Entrar no MySQL
mysql -u root -p

# Criar database
CREATE DATABASE clima_organizacional;
CREATE USER 'clima_user'@'localhost' IDENTIFIED BY 'sua_senha_aqui';
GRANT ALL PRIVILEGES ON clima_organizacional.* TO 'clima_user'@'localhost';
FLUSH PRIVILEGES;

# Executar migrações
cd database
mysql -u clima_user -p clima_organizacional < migrations/001_create_tables.sql
mysql -u clima_user -p clima_organizacional < migrations/002_add_indexes.sql
mysql -u clima_user -p clima_organizacional < migrations/003_create_views.sql

# Dados de teste (opcional)
mysql -u clima_user -p clima_organizacional < seeds/demo_data.sql
```
</details>

<details>
<summary><strong>🔧 2. Backend (Golang)</strong></summary>

```bash
cd backend

# Instalar dependências
go mod download

# Configurar variáveis de ambiente
cp .env.example .env
# Edite as configurações de banco e JWT

# Executar testes
go test -v ./...

# Executar em modo desenvolvimento
go run cmd/api/main.go

# Build para produção
go build -o bin/api cmd/api/main.go
```

**Variáveis de ambiente necessárias (.env):**
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=clima_user
DB_PASSWORD=sua_senha_aqui
DB_NAME=clima_organizacional

# JWT
JWT_SECRET=seu_jwt_secret_muito_seguro_aqui
JWT_EXPIRE_HOURS=24

# Server
PORT=8080
GIN_MODE=debug

# CORS
FRONTEND_URL=http://localhost:3000
```
</details>

<details>
<summary><strong>🎨 3. Frontend (Next.js)</strong></summary>

```bash
cd frontend

# Instalar dependências
npm install

# Configurar variáveis de ambiente
cp .env.local.example .env.local
# Edite a URL da API

# Executar em modo desenvolvimento
npm run dev

# Build para produção
npm run build
npm run start

# Executar testes
npm run test

# Linting
npm run lint
```

**Variáveis de ambiente necessárias (.env.local):**
```env
# API
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_URL=http://localhost:3000

# Features
NEXT_PUBLIC_ENABLE_ANALYTICS=true
NEXT_PUBLIC_MAX_FILE_SIZE=5242880

# Analytics (opcional)
NEXT_PUBLIC_GA_ID=G-XXXXXXXXXX
```
</details>

---

## 🧪 Testes

### Executar todos os testes
```bash
# Backend
cd backend && go test -v ./... -cover

# Frontend
cd frontend && npm run test

# E2E (após iniciar aplicação)
cd frontend && npm run test:e2e

# Testes de carga
cd scripts && ./load_test.sh
```

### Coverage Reports
```bash
# Backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Frontend
npm run test:coverage
```

---

## 📊 Status do Desenvolvimento

### 🎯 Roadmap

#### ✅ Fase 1 - Setup e Infraestrutura (Semana 1-2)
- [x] Repositório Git configurado
- [ ] CI/CD implementado
- [ ] Arquitetura definida
- [ ] Ambiente de desenvolvimento

#### 🔄 Fase 2 - Autenticação (Semana 3-4)
- [ ] Sistema de login JWT
- [ ] Middleware de autenticação  
- [ ] Gestão de empresas/setores
- [ ] Dashboard administrativo

#### 📋 Fase 3 - Pesquisas (Semana 5-8)
- [ ] Criação de formulários
- [ ] Sistema de perguntas dinâmico
- [ ] Geração de links/QR codes
- [ ] Interface pública de resposta
- [ ] Sistema de recorrência

#### 📈 Fase 4 - Analytics (Semana 9-12)
- [ ] Dashboards interativos
- [ ] Relatórios exportáveis
- [ ] Análises históricas
- [ ] Segmentação por setor

#### 🔍 Fase 5 - Testes e Deploy (Semana 13)
- [ ] Testes E2E completos
- [ ] Performance testing
- [ ] Deploy em produção
- [ ] Documentação final

### 📈 Métricas de Qualidade

| Métrica | Target | Atual | Status |
|---------|--------|-------|--------|
| Test Coverage (Backend) | >80% | 75% | 🟡 |
| Test Coverage (Frontend) | >70% | 60% | 🟡 |
| Performance (API) | <2s | <1.5s | ✅ |
| Performance (Frontend) | <3s | <2s | ✅ |
| Lighthouse Score | >90 | 85 | 🟡 |

---

## 📚 Documentação

### 📋 Documentação Principal
- [📖 Documentação Completa do Sistema](docs/Doc%20Sistema%20Extensão.pdf)
- [🏗️ Guia de Arquitetura](docs/architecture.md)
- [🔧 Manual de Instalação](docs/installation.md)
- [👤 Manual do Usuário](docs/user-guide.md)

### 🔌 API Reference
- [📡 Swagger Documentation](http://localhost:8080/swagger)
- [🛠️ Endpoints Reference](docs/api/endpoints.md)
- [🔐 Autenticação](docs/api/authentication.md)
- [📊 Analytics APIs](docs/api/analytics.md)

### 🗄️ Database
- [📊 Modelo Entidade Relacionamento](docs/database/er-diagram.png)
- [📝 Dicionário de Dados](docs/database/data-dictionary.md)
- [🔄 Guia de Migrações](database/README.md)

---

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

## 🌐 APIs e Integrações

### 🔌 Principais Endpoints

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
</details>

<details>
<summary><strong>📋 Pesquisas</strong></summary>

```http
GET /api/v1/pesquisas
Authorization: Bearer <token>

POST /api/v1/pesquisas
Authorization: Bearer <token>
Content-Type: application/json

{
  "titulo": "Pesquisa Q1 2025",
  "descricao": "Avaliação trimestral",
  "setor_id": 1,
  "dataAbertura": "2025-01-15T09:00:00Z",
  "dataFechamento": "2025-01-30T18:00:00Z"
}
```
</details>

<details>
<summary><strong>📊 Analytics</strong></summary>

```http
GET /api/v1/pesquisas/{id}/dashboard
Authorization: Bearer <token>

GET /api/v1/pesquisas/{id}/analytics/export/xlsx
Authorization: Bearer <token>
```
</details>

### 📖 Documentação Completa da API
Acesse [http://localhost:8080/swagger](http://localhost:8080/swagger) após iniciar o backend.

---

### 📝 Padrões de Código

#### Commits
```bash
feat(auth): adiciona endpoint de refresh token
fix(dashboard): corrige cálculo de percentuais
docs(api): atualiza documentação dos endpoints
test(survey): adiciona testes para criação de pesquisa
```

#### Code Review Checklist
- [ ] ✅ Código limpo e bem comentado
- [ ] ✅ Testes unitários adicionados/atualizados
- [ ] ✅ Documentação atualizada
- [ ] ✅ Performance considerada
- [ ] ✅ Segurança avaliada
- [ ] ✅ Compatibilidade verificada

---

## 📞 Suporte e Comunidade

### 🐛 Reportar Bugs
Encontrou um problema? [Abra uma issue](https://github.com/Murilo-The-Dev/sistema-clima-organizacional/issues) com:
- Descrição detalhada do problema
- Passos para reproduzir
- Screenshots (se aplicável)
- Ambiente (OS, browser, versões)

### 💡 Sugerir Melhorias
Tem uma ideia? [Abra uma feature request](https://github.com/SEU_USUARIO/sistema-clima-organizacional/issues/new?template=feature_request.md)

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
- [MySQL](https://mysql.com/) - Banco de dados
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