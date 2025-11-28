// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// LogAuditoriaFilter define os parâmetros de filtro para consultas de logs de auditoria.
// Inclui paginação herdada de PaginationRequest.
type LogAuditoriaFilter struct {
	PaginationRequest                 // Herda Page, Limit, OrderBy, Order
	IDUserAdmin   *int    `form:"id_user_admin"`   // Filtra por ID do usuário administrador
	IDEmpresa     *int    `form:"id_empresa"`      // Filtra por ID da empresa
	DataInicio    *string `form:"data_inicio"`     // Data inicial no formato RFC3339
	DataFim       *string `form:"data_fim"`        // Data final no formato RFC3339
	AcaoRealizada *string `form:"acao_realizada"`  // Filtra por ação realizada
	EnderecoIP    *string `form:"endereco_ip"`     // Filtra por endereço IP
}
