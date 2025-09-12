package response

import (
	"time"
	"organizational-climate-survey/backend/internal/domain/entity"
)

// AuditSummaryResponse resposta para resumo de auditoria
type AuditSummaryResponse struct {
	PeriodoInicio     string                 `json:"periodo_inicio"`
	PeriodoFim        string                 `json:"periodo_fim"`
	TotalEventos      int                    `json:"total_eventos"`
	AcoesPorTipo      map[string]int         `json:"acoes_por_tipo"`
	EventosPorUsuario map[int]int           `json:"eventos_por_usuario"`
	EventosPorDia     map[string]int        `json:"eventos_por_dia"`
}

// LogResponse resposta simplificada para log
type LogResponse struct {
	ID            int       `json:"id_log"`
	TimeStamp     time.Time `json:"timestamp"`
	AcaoRealizada string    `json:"acao_realizada"`
	Detalhes      string    `json:"detalhes"`
	EnderecoIP    string    `json:"endereco_ip"`
}

func ToLogResponse(log *entity.LogAuditoria) LogResponse {
	return LogResponse{
		ID:            log.ID,
		TimeStamp:     log.TimeStamp,
		AcaoRealizada: log.AcaoRealizada,
		Detalhes:      log.Detalhes,
		EnderecoIP:    log.EnderecoIP,
	}
}