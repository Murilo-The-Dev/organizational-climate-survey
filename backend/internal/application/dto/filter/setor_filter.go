// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// SetorFilter define filtros para listagem de setores.
type SetorFilter struct {
	PaginationRequest         // Campos de paginação (herança)
	IDEmpresa         *int    `form:"id_empresa" example:"1"` // Filtra pelo ID da empresa
	Nome              *string `form:"nome" example:"RH"`      // Filtra pelo nome do setor
}
