package repository

import (
    "context"
	"organizational-climate-survey/backend/internal/domain/entity"
)

type DashboardRepository interface {
    Create(ctx context.Context, dashboard *entity.Dashboard) error
    GetByID(ctx context.Context, id int) (*entity.Dashboard, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Dashboard, error)
}

type EmpresaRepository interface {
    Create(ctx context.Context, empresa *entity.Empresa) error
    GetByID(ctx context.Context, id int) (*entity.Empresa, error)
    GetByCNPJ(ctx context.Context, cnpj string) (*entity.Empresa, error)
    Update(ctx context.Context, empresa *entity.Empresa) error
    Delete(ctx context.Context, id int) error
}

type LogAuditoriaRepository interface {
    Create(ctx context.Context, log *entity.LogAuditoria) error
    GetByID(ctx context.Context, id int) (*entity.LogAuditoria, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.LogAuditoria, error)
}

type PesquisaRepository interface {
    Create(ctx context.Context, pesquisa *entity.Pesquisa) error
    GetByID(ctx context.Context, id int) (*entity.Pesquisa, error)
    GetByLinkAcesso(ctx context.Context, link string) (*entity.Pesquisa, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Pesquisa, error)
    UpdateStatus(ctx context.Context, id int, status string) error
}

type PerguntaRepository interface {
    Create(ctx context.Context, pergunta *entity.Pergunta) error
    GetByID(ctx context.Context, id int) (*entity.Pergunta, error)
    ListByPesquisa(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error)
    Update(ctx context.Context, pergunta *entity.Pergunta) error
    Delete(ctx context.Context, id int) error
}

type RespostaRepository interface {
    CreateBatch(ctx context.Context, respostas []*entity.Resposta) error
    CountByPesquisa(ctx context.Context, pesquisaID int) (int, error)
    GetAggregatedByPergunta(ctx context.Context, perguntaID int) (map[string]int, error)
}

type SetorRepository interface {
    Create(ctx context.Context, setor *entity.Setor) error
    GetByID(ctx context.Context, id int) (*entity.Setor, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Setor, error)
    Update(ctx context.Context, setor *entity.Setor) error
    Delete(ctx context.Context, id int) error
}

type UsuarioAdministradorRepository interface {
    Create(ctx context.Context, usuario *entity.UsuarioAdministrador) error
    GetByID(ctx context.Context, id int) (*entity.UsuarioAdministrador, error)
    ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.UsuarioAdministrador, error)
    Update(ctx context.Context, usuario *entity.UsuarioAdministrador) error
    Delete(ctx context.Context, id int) error
}