// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// PesquisaFilter define filtros para listagem de pesquisas.
type PesquisaFilter struct {
	PaginationRequest           // Campos de paginação (herança)
	IDEmpresa      *int    `form:"id_empresa"`       // Filtra pelo ID da empresa
	IDSetor        *int    `form:"id_setor"`         // Filtra pelo ID do setor
	IDUserAdmin    *int    `form:"id_user_admin"`    // Filtra pelo ID do usuário administrador
	Status         *string `form:"status" binding:"omitempty,oneof=Rascunho Ativa Concluída Arquivada"` // Filtra pelo status da pesquisa
	DataInicio     *string `form:"data_inicio"`      // Filtra a partir da data inicial (RFC3339)
	DataFim        *string `form:"data_fim"`         // Filtra até a data final (RFC3339)
	Titulo         *string `form:"titulo"`           // Filtra pelo título da pesquisa
	Anonimato      *bool   `form:"anonimato"`        // Filtra se a pesquisa é anônima
	ComPerguntas   *bool   `form:"com_perguntas"`    // Filtra pesquisas que possuem perguntas
	ComRespostas   *bool   `form:"com_respostas"`    // Filtra pesquisas que possuem respostas
}
