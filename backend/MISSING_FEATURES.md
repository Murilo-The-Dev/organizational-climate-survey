# Missing Backend Features (Audit)

This document tracks roadmap status for backend after the latest implementation cycle.

## Implemented
- JWT authentication and refresh endpoint
- Company and sector management endpoints
- Survey creation and question types (MultiplaEscolha, EscalaNumerica, RespostaAberta, SimNao)
- Anonymous submission token flow
- Audit logging endpoints
- Stronger anti-fraud duplicate prevention for submissions
  - Heuristic: 2-of-3 matching signals in 24h window (IP hash, User-Agent hash, Accept-Language hash)
  - API behavior: duplicate attempt returns HTTP 409 with message "Você já respondeu esta pesquisa."
- LGPD per-submission anonymization endpoint
  - Endpoint: DELETE /submissoes/{submissao_id}/dados-pessoais
  - Behavior: anonymizes personal tracking fields while preserving analytical responses
- QR code generation for survey links
  - Endpoint: POST /pesquisas/{id}/qrcode
  - Behavior: generates PNG, stores file, persists qrcode_path, returns qr_code_base64/survey_url/qr_payload
- Recurring survey scheduler in runtime
  - Cron-based cycle for materializing recurring surveys and cloning questions/dashboard
- Historical comparison endpoint
  - Endpoint: GET /pesquisas/{pesquisa_id}/comparativo
  - Behavior: historical score series with variation from previous survey
- Segmentation by sector in dashboard listing/comparison
  - Query support: setor_id in dashboard listing and historical comparison
- Export endpoints (CSV, XLSX, PDF)
  - Endpoint: GET /dashboards/{id}/export
  - Endpoint: GET /pesquisas/{pesquisa_id}/export

## Remaining Gaps
- Additional business-level KPIs and advanced analytics models can still be expanded (non-blocking enhancement).

## Migrations Applied
- Added: migrations/006_strengthen_submissao_signals.sql
  - Columns: user_agent_hash, accept_language_hash, fingerprint_composto
  - Indexes for anti-fraud signal queries
