// Package response contém structs usadas para enviar dados da API como respostas.
package response

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"time"
)

// AuditSummaryResponse fornece resumo agregado de eventos de auditoria.
type AuditSummaryResponse struct {
	PeriodoInicio     string         `json:"periodo_inicio" example:"2026-01-01T00:00:00Z"`                      // Data inicial do período analisado
	PeriodoFim        string         `json:"periodo_fim" example:"2026-01-31T23:59:59Z"`                         // Data final do período analisado
	TotalEventos      int            `json:"total_eventos" example:"120"`                                        // Total de eventos registrados
	AcoesPorTipo      map[string]int `json:"acoes_por_tipo" example:"{\"Login\":45,\"Criacao de pesquisa\":12}"` // Contagem de ações agrupadas por tipo
	EventosPorUsuario map[int]int    `json:"eventos_por_usuario" example:"{\"1\":30,\"2\":25}"`                  // Contagem de eventos por ID de usuário
	EventosPorDia     map[string]int `json:"eventos_por_dia" example:"{\"2026-01-10\":8,\"2026-01-11\":11}"`     // Contagem de eventos por dia
}

// LogResponse representa um único registro de auditoria.
type LogResponse struct {
	ID            int       `json:"id_log" example:"55"`                                           // ID do log
	TimeStamp     time.Time `json:"timestamp" swaggertype:"string" example:"2026-01-10T09:00:00Z"` // Data e hora do evento
	AcaoRealizada string    `json:"acao_realizada" example:"Login"`                                // Ação realizada pelo usuário
	Detalhes      string    `json:"detalhes" example:"Usuário autenticado com sucesso"`            // Detalhes adicionais do evento
	EnderecoIP    string    `json:"endereco_ip" example:"192.168.1.10"`                            // IP do usuário que realizou a ação
}

// ToLogResponse converte a entidade de domínio LogAuditoria em LogResponse
func ToLogResponse(log *entity.LogAuditoria) LogResponse {
	return LogResponse{
		ID:            log.ID,
		TimeStamp:     log.TimeStamp,
		AcaoRealizada: log.AcaoRealizada,
		Detalhes:      log.Detalhes,
		EnderecoIP:    log.EnderecoIP,
	}
}
