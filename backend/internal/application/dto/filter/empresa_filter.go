// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// EmpresaFilter define os parâmetros de filtro para consultas de empresas.
// Inclui paginação herdada de PaginationRequest.
type EmpresaFilter struct {
	PaginationRequest         // Herda Page, Limit, OrderBy, Order
	NomeFantasia      *string `form:"nome_fantasia" example:"Acme"`               // Filtra pelo nome fantasia da empresa
	RazaoSocial       *string `form:"razao_social" example:"Acme LTDA"`           // Filtra pela razão social da empresa
	CNPJ              *string `form:"cnpj" example:"12.345.678/0001-90"`          // Filtra pelo CNPJ
	DataInicio        *string `form:"data_inicio" example:"2026-01-01T00:00:00Z"` // Data inicial de criação no formato RFC3339
	DataFim           *string `form:"data_fim" example:"2026-01-31T23:59:59Z"`    // Data final de criação no formato RFC3339
}
