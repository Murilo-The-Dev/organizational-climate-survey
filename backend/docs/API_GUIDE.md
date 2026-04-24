# Guia da API - Pesquisa de Clima Organizacional

Este guia foi escrito para consumo do frontend e integrações externas.

## Base URLs

- API local: `http://localhost:8080`
- Prefixo versionado: `/api/v1`
- Swagger UI: `/swagger/index.html`
- OpenAPI JSON: `/swagger/doc.json`
- Health check: `/health`

## Autenticação

### Fluxo JWT (rotas protegidas)

1. Faça login em `POST /api/v1/auth/login`.
2. Copie o token retornado em `data.token`.
3. Envie o header abaixo nas rotas protegidas:

```http
Authorization: Bearer <token>
```

### Fluxo de submissão anônima (pesquisa)

1. Gere token de acesso em `POST /api/v1/pesquisas/{pesquisa_id}/token`.
2. Envie respostas em `POST /api/v1/respostas/submit` com `token_acesso`.

## Padrão de resposta

### Sucesso

```json
{
  "success": true,
  "message": "Operação realizada com sucesso",
  "data": {}
}
```

### Erro

```json
{
  "success": false,
  "message": "Validação falhou",
  "error": "campo inválido",
  "code": "BAD_REQUEST"
}
```

## Códigos de erro comuns

- `BAD_REQUEST` (400)
- `UNAUTHORIZED` (401)
- `FORBIDDEN` (403)
- `NOT_FOUND` (404)
- `CONFLICT` (409)
- `VALIDATION_ERROR` (422)
- `RATE_LIMITED` (429)
- `INTERNAL_ERROR` (500)

## Matriz completa de endpoints

### Sistema

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/health` | Pública | Verificação de saúde da API |

### Auth

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/auth/login` | Pública | Login e emissão de JWT |
| POST | `/api/v1/auth/forgot-password` | Pública | Início de recuperação de senha |
| POST | `/api/v1/auth/validate` | Pública | Validação explícita de token |
| POST | `/api/v1/auth/refresh` | Pública | Renovação de token JWT |
| POST | `/api/v1/auth/logout` | Bearer | Logout lógico |
| POST | `/api/v1/auth/change-password` | Bearer | Alteração de senha do usuário autenticado |

### Bootstrap

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/bootstrap` | Pública | Inicialização única do sistema |

### Empresas

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/empresas` | Bearer | Cria empresa |
| GET | `/api/v1/empresas` | Bearer | Lista empresas (limit/offset) |
| GET | `/api/v1/empresas/{id}` | Bearer | Busca empresa por ID |
| PUT | `/api/v1/empresas/{id}` | Bearer | Atualiza empresa |
| DELETE | `/api/v1/empresas/{id}` | Bearer | Remove empresa |
| GET | `/api/v1/empresas/cnpj/{cnpj}` | Bearer | Busca empresa por CNPJ |

### Usuários Administradores

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/usuarios-administradores` | Bearer | Cria usuário administrador |
| GET | `/api/v1/usuarios-administradores/{id}` | Bearer | Busca usuário por ID |
| PUT | `/api/v1/usuarios-administradores/{id}` | Bearer | Atualiza usuário |
| DELETE | `/api/v1/usuarios-administradores/{id}` | Bearer | Inativa usuário |
| PUT | `/api/v1/usuarios-administradores/{id}/password` | Bearer | Atualiza senha |
| PUT | `/api/v1/usuarios-administradores/{id}/status` | Bearer | Atualiza status |
| GET | `/api/v1/usuarios-administradores/email/{email}` | Bearer | Busca usuário por email |
| GET | `/api/v1/empresas/{empresa_id}/usuarios-administradores` | Bearer | Lista usuários por empresa |

### Setores

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/setores` | Bearer | Cria setor |
| GET | `/api/v1/setores/{id}` | Bearer | Busca setor por ID |
| PUT | `/api/v1/setores/{id}` | Bearer | Atualiza setor |
| DELETE | `/api/v1/setores/{id}` | Bearer | Remove setor |
| GET | `/api/v1/empresas/{empresa_id}/setores` | Bearer | Lista setores da empresa |
| GET | `/api/v1/empresas/{empresa_id}/setores/nome/{nome}` | Bearer | Busca setor por nome |

### Pesquisas

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/pesquisas` | Bearer | Cria pesquisa |
| GET | `/api/v1/pesquisas/{id}` | Bearer | Busca pesquisa por ID |
| PUT | `/api/v1/pesquisas/{id}` | Bearer | Atualiza pesquisa |
| DELETE | `/api/v1/pesquisas/{id}` | Bearer | Remove pesquisa |
| PUT | `/api/v1/pesquisas/{id}/status` | Bearer | Atualiza status da pesquisa |
| POST | `/api/v1/pesquisas/{id}/qrcode` | Bearer | Gera QR Code |
| GET | `/api/v1/pesquisas/link/{link}` | Bearer | Busca pesquisa por link |
| GET | `/api/v1/empresas/{empresa_id}/pesquisas` | Bearer | Lista pesquisas da empresa |
| GET | `/api/v1/empresas/{empresa_id}/pesquisas/active` | Bearer | Lista pesquisas ativas |
| GET | `/api/v1/setores/{setor_id}/pesquisas` | Bearer | Lista pesquisas do setor |

### Perguntas

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/perguntas` | Bearer | Cria pergunta |
| POST | `/api/v1/perguntas/batch` | Bearer | Cria perguntas em lote |
| GET | `/api/v1/perguntas/{id}` | Bearer | Busca pergunta por ID |
| PUT | `/api/v1/perguntas/{id}` | Bearer | Atualiza pergunta |
| DELETE | `/api/v1/perguntas/{id}` | Bearer | Remove pergunta |
| PUT | `/api/v1/perguntas/{id}/ordem` | Bearer | Atualiza ordem da pergunta |
| GET | `/api/v1/pesquisas/{pesquisa_id}/perguntas` | Bearer | Lista perguntas por pesquisa |
| PUT | `/api/v1/pesquisas/{pesquisa_id}/perguntas/reorder` | Bearer | Reordena perguntas da pesquisa |
| GET | `/api/v1/pesquisas/{pesquisa_id}/perguntas/with-stats` | Bearer | Lista perguntas com estatísticas |

### Respostas e Submissões

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/pesquisas/{pesquisa_id}/token` | Pública | Gera token para responder pesquisa |
| POST | `/api/v1/respostas/submit` | Token de pesquisa | Submete respostas anônimas |
| GET | `/api/v1/pesquisas/{pesquisa_id}/submissions/stats` | Bearer | Estatísticas de submissões |
| GET | `/api/v1/pesquisas/{pesquisa_id}/respostas/stats` | Bearer | Estatísticas de respostas |
| GET | `/api/v1/pesquisas/{pesquisa_id}/respostas/aggregated` | Bearer | Respostas agregadas da pesquisa |
| GET | `/api/v1/pesquisas/{pesquisa_id}/respostas/by-date` | Bearer | Respostas por período |
| GET | `/api/v1/pesquisas/{pesquisa_id}/respostas/count` | Bearer | Contagem de respostas por pesquisa |
| DELETE | `/api/v1/pesquisas/{pesquisa_id}/respostas` | Bearer | Remove respostas por pesquisa |
| DELETE | `/api/v1/submissoes/{submissao_id}/dados-pessoais` | Bearer | Anonimiza dados pessoais (LGPD) |
| GET | `/api/v1/perguntas/{pergunta_id}/respostas/aggregated` | Bearer | Respostas agregadas por pergunta |
| GET | `/api/v1/perguntas/{pergunta_id}/respostas/count` | Bearer | Contagem por pergunta |
| GET | `/api/v1/perguntas/{pergunta_id}/respostas/stats` | Bearer | Estatísticas por pergunta |

### Dashboards

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/dashboards` | Bearer | Cria dashboard |
| GET | `/api/v1/dashboards/{id}` | Bearer | Busca dashboard por ID |
| PUT | `/api/v1/dashboards/{id}` | Bearer | Atualiza dashboard |
| DELETE | `/api/v1/dashboards/{id}` | Bearer | Remove dashboard |
| GET | `/api/v1/dashboards/{id}/data` | Bearer | Dados analíticos do dashboard |
| POST | `/api/v1/dashboards/{id}/refresh` | Bearer | Recalcula dashboard |
| GET | `/api/v1/dashboards/{id}/export` | Bearer | Exporta dashboard (pdf/xlsx/csv) |
| GET | `/api/v1/dashboards/{id}/metrics` | Bearer | Métricas resumidas |
| GET | `/api/v1/pesquisas/{pesquisa_id}/comparativo` | Bearer | Comparativo histórico |
| GET | `/api/v1/pesquisas/{pesquisa_id}/export` | Bearer | Exporta análise da pesquisa |
| GET | `/api/v1/pesquisas/{pesquisa_id}/dashboard` | Bearer | Busca dashboard por pesquisa |
| GET | `/api/v1/empresas/{empresa_id}/dashboards` | Bearer | Lista dashboards por empresa |

### Logs de Auditoria (admin)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/v1/logs-auditoria` | Bearer (admin) | Cria log |
| GET | `/api/v1/logs-auditoria/{id}` | Bearer (admin) | Busca log por ID |
| POST | `/api/v1/logs-auditoria/clean` | Bearer (admin) | Limpa logs antigos |
| GET | `/api/v1/empresas/{empresa_id}/logs-auditoria` | Bearer (admin) | Lista logs por empresa |
| GET | `/api/v1/empresas/{empresa_id}/logs-auditoria/by-date` | Bearer (admin) | Lista logs por período |
| GET | `/api/v1/empresas/{empresa_id}/logs-auditoria/by-action` | Bearer (admin) | Lista logs por ação |
| GET | `/api/v1/empresas/{empresa_id}/logs-auditoria/summary` | Bearer (admin) | Resumo de auditoria |
| GET | `/api/v1/empresas/{empresa_id}/logs-auditoria/export` | Bearer (admin) | Exporta logs |
| GET | `/api/v1/usuarios-administradores/{user_admin_id}/logs-auditoria` | Bearer (admin) | Lista logs por usuário |

## Observações importantes

- O endpoint `POST /api/v1/auth/validate` é `POST` (não `GET`).
- O endpoint `GET /api/v1/pesquisas/link/{link}` está registrado em rotas autenticadas no backend atual.
- A submissão anônima depende do token emitido por `POST /api/v1/pesquisas/{pesquisa_id}/token`.
