// Package response contém structs usadas para enviar dados da API como respostas.
package response

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"time"
)

// DashboardResponse representa as informações principais de um dashboard.
type DashboardResponse struct {
	ID               int                    `json:"id_dashboard" example:"3"`                                         // ID do dashboard
	IDPesquisa       int                    `json:"id_pesquisa" example:"10"`                                         // ID da pesquisa associada
	Titulo           string                 `json:"titulo" example:"Dashboard Clima RH"`                              // Título do dashboard
	DataCriacao      time.Time              `json:"data_criacao" swaggertype:"string" example:"2026-01-05T08:00:00Z"` // Data de criação do dashboard
	ConfigFiltros    *string                `json:"config_filtros,omitempty" example:"{\"setor\":\"RH\"}"`            // Configuração de filtros, opcional
	TotalRespostas   int                    `json:"total_respostas,omitempty" example:"250"`                          // Total de respostas coletadas
	TaxaParticipacao float64                `json:"taxa_participacao,omitempty" example:"82.4"`                       // Taxa de participação média
	Metricas         map[string]interface{} `json:"metricas,omitempty" example:"{\"engajamento\":82.4}"`              // Métricas adicionais, opcionais
}

// DashboardDataResponse inclui o dashboard com dados e métricas adicionais.
type DashboardDataResponse struct {
	Dashboard DashboardResponse `json:"dashboard" example:"{\"id_dashboard\":3}"`                       // Dashboard principal
	Data      interface{}       `json:"data"`                                                           // Dados detalhados do dashboard
	Metrics   interface{}       `json:"metrics"`                                                        // Métricas adicionais
	UpdatedAt time.Time         `json:"updated_at" swaggertype:"string" example:"2026-01-10T09:30:00Z"` // Data da última atualização
}

// ExportResponse representa os dados de um arquivo exportado do dashboard.
type ExportResponse struct {
	Format    string    `json:"format" example:"pdf"`                                                  // Formato do arquivo (ex: CSV, PDF)
	Filename  string    `json:"filename" example:"dashboard_3.pdf"`                                    // Nome do arquivo
	URL       string    `json:"url,omitempty" example:"https://api.exemplo.com/files/dashboard_3.pdf"` // URL para download, opcional
	Size      int64     `json:"size,omitempty" example:"102400"`                                       // Tamanho do arquivo, opcional
	CreatedAt time.Time `json:"created_at" swaggertype:"string" example:"2026-01-10T09:45:00Z"`        // Data de criação do arquivo
}

// MetricsResponse agrega métricas de respostas e estatísticas do dashboard.
type MetricsResponse struct {
	TotalRespostas    int                    `json:"total_respostas" example:"300"`                                                      // Total de respostas
	DataUltimaResp    *time.Time             `json:"data_ultima_resposta,omitempty" swaggertype:"string" example:"2026-01-10T09:40:00Z"` // Data da última resposta, opcional
	TaxaParticipacao  float64                `json:"taxa_participacao" example:"79.2"`                                                   // Taxa de participação média
	MetricasPorTipo   map[string]interface{} `json:"metricas_por_tipo" example:"{\"EscalaNumerica\":4.2}"`                               // Métricas detalhadas por tipo
	ResumoEstatistico map[string]interface{} `json:"resumo_estatistico" example:"{\"media_geral\":4.1}"`                                 // Resumo estatístico geral
}

// ToDashboardResponse converte uma entidade Dashboard para DashboardResponse.
func ToDashboardResponse(dashboard *entity.Dashboard) DashboardResponse {
	return DashboardResponse{
		ID:               dashboard.ID,
		IDPesquisa:       dashboard.IDPesquisa,
		Titulo:           dashboard.Titulo,
		DataCriacao:      dashboard.DataCriacao,
		ConfigFiltros:    dashboard.ConfigFiltros, // já é ponteiro
		TotalRespostas:   dashboard.TotalRespostas,
		TaxaParticipacao: dashboard.TaxaParticipacao,
		Metricas:         dashboard.Metricas,
	}
}
