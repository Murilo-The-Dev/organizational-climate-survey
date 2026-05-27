// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// LogAuditoriaFilter define os parâmetros de filtro para consultas de logs de auditoria.
// Inclui paginação herdada de PaginationRequest.
type LogAuditoriaFilter struct {
	PaginationRequest         // Herda Page, Limit, OrderBy, Order
	IDUserAdmin       *int    `form:"id_user_admin" example:"1"`                  // Filtra por ID do usuário administrador
	IDEmpresa         *int    `form:"id_empresa" example:"1"`                     // Filtra por ID da empresa
	DataInicio        *string `form:"data_inicio" example:"2026-01-01T00:00:00Z"` // Data inicial no formato RFC3339
	DataFim           *string `form:"data_fim" example:"2026-01-31T23:59:59Z"`    // Data final no formato RFC3339
	AcaoRealizada     *string `form:"acao_realizada" example:"Login"`             // Filtra por ação realizada
	EnderecoIP        *string `form:"endereco_ip" example:"192.168.1.10"`         // Filtra por endereço IP
}
