// Package response contém structs usadas para enviar dados da API como respostas.
package response

import (
	"time"
	"organizational-climate-survey/backend/internal/domain/entity"
)

// AuditSummaryResponse fornece resumo agregado de eventos de auditoria.
type AuditSummaryResponse struct {
	PeriodoInicio     string          `json:"periodo_inicio"`      // Data inicial do período analisado
	PeriodoFim        string          `json:"periodo_fim"`         // Data final do período analisado
	TotalEventos      int             `json:"total_eventos"`       // Total de eventos registrados
	AcoesPorTipo      map[string]int  `json:"acoes_por_tipo"`      // Contagem de ações agrupadas por tipo
	EventosPorUsuario map[int]int     `json:"eventos_por_usuario"` // Contagem de eventos por ID de usuário
	EventosPorDia     map[string]int  `json:"eventos_por_dia"`     // Contagem de eventos por dia
}

// LogResponse representa um único registro de auditoria.
type LogResponse struct {
	ID            int       `json:"id_log"`          // ID do log
	TimeStamp     time.Time `json:"timestamp"`       // Data e hora do evento
	AcaoRealizada string    `json:"acao_realizada"`  // Ação realizada pelo usuário
	Detalhes      string    `json:"detalhes"`        // Detalhes adicionais do evento
	EnderecoIP    string    `json:"endereco_ip"`     // IP do usuário que realizou a ação
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
