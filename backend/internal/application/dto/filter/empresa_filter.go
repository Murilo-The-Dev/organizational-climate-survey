// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// EmpresaFilter define os parâmetros de filtro para consultas de empresas.
// Inclui paginação herdada de PaginationRequest.
type EmpresaFilter struct {
	PaginationRequest                 // Herda Page, Limit, OrderBy, Order
	NomeFantasia *string `form:"nome_fantasia"` // Filtra pelo nome fantasia da empresa
	RazaoSocial  *string `form:"razao_social"`  // Filtra pela razão social da empresa
	CNPJ         *string `form:"cnpj"`          // Filtra pelo CNPJ
	DataInicio   *string `form:"data_inicio"`   // Data inicial de criação no formato RFC3339
	DataFim      *string `form:"data_fim"`      // Data final de criação no formato RFC3339
}