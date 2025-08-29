package repository

import (
    "context"
    "organizational-climate-survey/backend/internal/domain/entity"
)

// DashboardRepository gerencia operações relacionadas aos dashboards
type DashboardRepository interface {
    Create(ctx context.Context, dashboard *entity.Dashboard) error
    GetByID(ctx context.Context, id int) (*entity.Dashboard, error)
    GetByPesquisaID(ctx context.Context, pesquisaID int) (*entity.Dashboard, error) // Relação 1:1
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Dashboard, error)
    Update(ctx context.Context, dashboard *entity.Dashboard) error
    Delete(ctx context.Context, id int) error
}

// EmpresaRepository gerencia operações relacionadas às empresas
type EmpresaRepository interface {
    Create(ctx context.Context, empresa *entity.Empresa) error
    GetByID(ctx context.Context, id int) (*entity.Empresa, error)
    GetByCNPJ(ctx context.Context, cnpj string) (*entity.Empresa, error)
    List(ctx context.Context, limit, offset int) ([]*entity.Empresa, error)
    Update(ctx context.Context, empresa *entity.Empresa) error
    Delete(ctx context.Context, id int) error
}

// LogAuditoriaRepository gerencia logs de auditoria
type LogAuditoriaRepository interface {
    Create(ctx context.Context, log *entity.LogAuditoria) error
    GetByID(ctx context.Context, id int) (*entity.LogAuditoria, error)
    ListByEmpresa(ctx context.Context, empresaID int, limit, offset int) ([]*entity.LogAuditoria, error)
    ListByUsuarioAdmin(ctx context.Context, userAdminID int, limit, offset int) ([]*entity.LogAuditoria, error)
    ListByDateRange(ctx context.Context, empresaID int, startDate, endDate string) ([]*entity.LogAuditoria, error)
}

// PesquisaRepository gerencia operações relacionadas às pesquisas
type PesquisaRepository interface {
    Create(ctx context.Context, pesquisa *entity.Pesquisa) error
    GetByID(ctx context.Context, id int) (*entity.Pesquisa, error)
    GetByLinkAcesso(ctx context.Context, link string) (*entity.Pesquisa, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Pesquisa, error)
    ListBySetor(ctx context.Context, setorID int) ([]*entity.Pesquisa, error)
    ListByStatus(ctx context.Context, empresaID int, status string) ([]*entity.Pesquisa, error)
    ListActive(ctx context.Context, empresaID int) ([]*entity.Pesquisa, error)
    Update(ctx context.Context, pesquisa *entity.Pesquisa) error
    UpdateStatus(ctx context.Context, id int, status string) error
    Delete(ctx context.Context, id int) error
}

// PerguntaRepository gerencia operações relacionadas às perguntas
type PerguntaRepository interface {
    Create(ctx context.Context, pergunta *entity.Pergunta) error
    CreateBatch(ctx context.Context, perguntas []*entity.Pergunta) error // Para criar múltiplas perguntas
    GetByID(ctx context.Context, id int) (*entity.Pergunta, error)
    ListByPesquisa(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error)
    Update(ctx context.Context, pergunta *entity.Pergunta) error
    Delete(ctx context.Context, id int) error
    UpdateOrdem(ctx context.Context, perguntaID int, novaOrdem int) error
}

// RespostaRepository gerencia operações relacionadas às respostas (com foco em anonimato)
type RespostaRepository interface {
    CreateBatch(ctx context.Context, respostas []*entity.Resposta) error
    CountByPesquisa(ctx context.Context, pesquisaID int) (int, error)
    CountByPergunta(ctx context.Context, perguntaID int) (int, error)
    GetAggregatedByPergunta(ctx context.Context, perguntaID int) (map[string]int, error)
    GetAggregatedByPesquisa(ctx context.Context, pesquisaID int) (map[int]map[string]int, error)
    GetResponsesByDateRange(ctx context.Context, pesquisaID int, startDate, endDate string) ([]*entity.Resposta, error)
    DeleteByPesquisa(ctx context.Context, pesquisaID int) error // Para limpeza após análise
}

// SetorRepository gerencia operações relacionadas aos setores
type SetorRepository interface {
    Create(ctx context.Context, setor *entity.Setor) error
    GetByID(ctx context.Context, id int) (*entity.Setor, error)
    GetByNome(ctx context.Context, empresaID int, nome string) (*entity.Setor, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Setor, error)
    Update(ctx context.Context, setor *entity.Setor) error
    Delete(ctx context.Context, id int) error
}

// UsuarioAdministradorRepository gerencia operações relacionadas aos usuários administradores
type UsuarioAdministradorRepository interface {
    Create(ctx context.Context, usuario *entity.UsuarioAdministrador) error
    GetByID(ctx context.Context, id int) (*entity.UsuarioAdministrador, error)
    GetByEmail(ctx context.Context, email string) (*entity.UsuarioAdministrador, error) // Para login
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.UsuarioAdministrador, error)
    ListByStatus(ctx context.Context, empresaID int, status string) ([]*entity.UsuarioAdministrador, error)
    Update(ctx context.Context, usuario *entity.UsuarioAdministrador) error
    UpdatePassword(ctx context.Context, id int, senhaHash string) error
    UpdateStatus(ctx context.Context, id int, status string) error
    Delete(ctx context.Context, id int) error
}

// Interfaces para operações mais complexas que podem envolver múltiplas entidades

// AnalyticsRepository para operações de análise de dados
type AnalyticsRepository interface {
    GetPesquisaMetrics(ctx context.Context, pesquisaID int) (map[string]interface{}, error)
    GetComparisonData(ctx context.Context, pesquisaIDs []int) (map[string]interface{}, error)
    GetSetorComparison(ctx context.Context, empresaID int, pesquisaID int) (map[string]interface{}, error)
}