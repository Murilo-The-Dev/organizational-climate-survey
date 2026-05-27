#!/bin/bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"
EMAIL="${EMAIL:-admin@test.com}"
SENHA="${SENHA:-senha123}"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

print_section() {
  echo -e "\n${GREEN}=== $1 ===${NC}"
}

request() {
  local method="$1"
  local path="$2"
  local data="${3:-}"
  local auth="${4:-}"

  echo "-> ${method} ${BASE_URL}${path}"
  if [[ -n "$auth" ]]; then
    if [[ -n "$data" ]]; then
      curl -sS -X "$method" "${BASE_URL}${path}" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${auth}" \
        -d "$data" | jq .
    else
      curl -sS -X "$method" "${BASE_URL}${path}" \
        -H "Authorization: Bearer ${auth}" | jq .
    fi
  else
    if [[ -n "$data" ]]; then
      curl -sS -X "$method" "${BASE_URL}${path}" \
        -H "Content-Type: application/json" \
        -d "$data" | jq .
    else
      curl -sS -X "$method" "${BASE_URL}${path}" | jq .
    fi
  fi
}

print_section "Health"
curl -sS "${BASE_URL%/api/v1}/health" | jq .

print_section "Bootstrap (first admin)"
request POST "/bootstrap" '{"empresa":{"nome_fantasia":"Empresa Teste","razao_social":"Empresa Teste LTDA","cnpj":"12.345.678/0001-99"},"usuario_administrador":{"nome_admin":"Admin Teste","email":"'"${EMAIL}"'","senha":"'"${SENHA}"'"}}'

print_section "Auth"
LOGIN_RESPONSE=$(curl -sS -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"'"${EMAIL}"'","senha":"'"${SENHA}"'"}')

echo "$LOGIN_RESPONSE" | jq .
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // empty')

if [[ -z "$TOKEN" ]]; then
  echo -e "${RED}Nao foi possivel obter token no login.${NC}"
  exit 1
fi

echo "Token obtido: ${TOKEN:0:20}..."

request POST "/auth/validate" '{"token":"'"${TOKEN}"'"}'
request POST "/auth/refresh" '{"token":"'"${TOKEN}"'"}'
request POST "/auth/change-password" '{"senha_atual":"'"${SENHA}"'","nova_senha":"senha12345"}' "$TOKEN"
request POST "/auth/logout" '' "$TOKEN"
request POST "/auth/forgot-password" '{"email":"'"${EMAIL}"'"}'

print_section "Empresas"
request GET "/empresas" '' "$TOKEN"
request GET "/empresas/1" '' "$TOKEN"
request GET "/empresas/cnpj/12.345.678/0001-99" '' "$TOKEN"
request POST "/empresas" '{"nome_fantasia":"Nova Empresa","razao_social":"Nova Empresa LTDA","cnpj":"98.765.432/0001-10"}' "$TOKEN"
request PUT "/empresas/1" '{"nome_fantasia":"Empresa Atualizada","razao_social":"Empresa Atualizada LTDA","cnpj":"12.345.678/0001-99"}' "$TOKEN"
request DELETE "/empresas/2" '' "$TOKEN"

print_section "Usuarios Administradores"
request GET "/usuarios-administradores" '' "$TOKEN"
request GET "/usuarios-administradores/1" '' "$TOKEN"
request POST "/usuarios-administradores" '{"id_empresa":1,"nome_admin":"Novo Admin","email":"novo.admin@test.com","senha_hash":"senha12345"}' "$TOKEN"
request PUT "/usuarios-administradores/1" '{"id_empresa":1,"nome_admin":"Admin Editado","email":"admin.editado@test.com","status":"Ativo"}' "$TOKEN"
request PUT "/usuarios-administradores/1/password" '{"nova_senha":"senha12345"}' "$TOKEN"
request PUT "/usuarios-administradores/1/status" '{"status":"Ativo"}' "$TOKEN"
request DELETE "/usuarios-administradores/2" '' "$TOKEN"

print_section "Setores"
request GET "/setores" '' "$TOKEN"
request GET "/setores/1" '' "$TOKEN"
request POST "/setores" '{"id_empresa":1,"nome_setor":"RH","descricao":"Recursos Humanos"}' "$TOKEN"
request PUT "/setores/1" '{"id_empresa":1,"nome_setor":"RH Atualizado","descricao":"RH atualizado"}' "$TOKEN"
request DELETE "/setores/2" '' "$TOKEN"

print_section "Pesquisas"
request GET "/pesquisas" '' "$TOKEN"
request GET "/pesquisas/1" '' "$TOKEN"
request POST "/pesquisas" '{"id_empresa":1,"id_setor":1,"titulo":"Pesquisa Q1","descricao":"Pesquisa de clima"}' "$TOKEN"
request PUT "/pesquisas/1" '{"id_empresa":1,"id_setor":1,"titulo":"Pesquisa Q1 Atualizada","descricao":"Pesquisa atualizada","status":"Rascunho"}' "$TOKEN"
request DELETE "/pesquisas/2" '' "$TOKEN"

print_section "Perguntas"
request GET "/perguntas/1" '' "$TOKEN"
request POST "/perguntas" '{"id_pesquisa":1,"texto_pergunta":"Voce esta satisfeito?","tipo_pergunta":"SimNao","ordem_exibicao":1}' "$TOKEN"
request PUT "/perguntas/1" '{"id_pesquisa":1,"texto_pergunta":"Voce esta satisfeito com o ambiente?","tipo_pergunta":"SimNao","ordem_exibicao":1}' "$TOKEN"
request DELETE "/perguntas/2" '' "$TOKEN"

print_section "Submissao anonima"
TOKEN_RESPONSE=$(curl -sS -X POST "${BASE_URL}/pesquisas/1/token" -H "Content-Type: application/json" -d '{}')
echo "$TOKEN_RESPONSE" | jq .
SURVEY_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.data.token_acesso // empty')

if [[ -n "$SURVEY_TOKEN" ]]; then
  request POST "/respostas/submit" '{"token_acesso":"'"${SURVEY_TOKEN}"'","respostas":[{"id_pergunta":1,"valor_resposta":"Sim"}]}'
fi

print_section "Respostas e estatisticas"
request GET "/pesquisas/1/respostas/stats" '' "$TOKEN"
request GET "/pesquisas/1/respostas/aggregated" '' "$TOKEN"
request GET "/pesquisas/1/respostas/by-date?start_date=2025-01-01&end_date=2025-12-31" '' "$TOKEN"
request GET "/pesquisas/1/respostas/count" '' "$TOKEN"
request DELETE "/pesquisas/1/respostas" '' "$TOKEN"
request GET "/perguntas/1/respostas/aggregated" '' "$TOKEN"
request GET "/perguntas/1/respostas/count" '' "$TOKEN"
request GET "/perguntas/1/respostas/stats" '' "$TOKEN"

print_section "Dashboard e auditoria"
request GET "/dashboards/1" '' "$TOKEN"
request GET "/logs-auditoria" '' "$TOKEN"
request GET "/logs-auditoria/1" '' "$TOKEN"
request GET "/pesquisas/1/submissions/stats" '' "$TOKEN"

echo -e "\n${GREEN}Script finalizado. Todos os endpoints registrados foram exercitados.${NC}"
