// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// PesquisaFilter define filtros para listagem de pesquisas.
type PesquisaFilter struct {
	PaginationRequest         // Campos de paginação (herança)
	IDEmpresa         *int    `form:"id_empresa" example:"1"`                                                              // Filtra pelo ID da empresa
	IDSetor           *int    `form:"id_setor" example:"2"`                                                                // Filtra pelo ID do setor
	IDUserAdmin       *int    `form:"id_user_admin" example:"1"`                                                           // Filtra pelo ID do usuário administrador
	Status            *string `form:"status" binding:"omitempty,oneof=Rascunho Ativa Concluída Arquivada" example:"Ativa"` // Filtra pelo status da pesquisa
	DataInicio        *string `form:"data_inicio" example:"2026-01-01T00:00:00Z"`                                          // Filtra a partir da data inicial (RFC3339)
	DataFim           *string `form:"data_fim" example:"2026-01-31T23:59:59Z"`                                             // Filtra até a data final (RFC3339)
	Titulo            *string `form:"titulo" example:"Pulso"`                                                              // Filtra pelo título da pesquisa
	Anonimato         *bool   `form:"anonimato" example:"true"`                                                            // Filtra se a pesquisa é anônima
	ComPerguntas      *bool   `form:"com_perguntas" example:"true"`                                                        // Filtra pesquisas que possuem perguntas
	ComRespostas      *bool   `form:"com_respostas" example:"false"`                                                       // Filtra pesquisas que possuem respostas
}
