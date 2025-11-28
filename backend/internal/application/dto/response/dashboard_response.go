// Package response contém structs usadas para enviar dados da API como respostas.
package response

import (
    "time"
    "organizational-climate-survey/backend/internal/domain/entity"
)

// DashboardResponse representa as informações principais de um dashboard.
type DashboardResponse struct {
    ID               int                    `json:"id_dashboard"`            // ID do dashboard
    IDPesquisa       int                    `json:"id_pesquisa"`             // ID da pesquisa associada
    Titulo           string                 `json:"titulo"`                  // Título do dashboard
    DataCriacao      time.Time              `json:"data_criacao"`            // Data de criação do dashboard
    ConfigFiltros    *string                `json:"config_filtros,omitempty"`// Configuração de filtros, opcional
    TotalRespostas   int                    `json:"total_respostas,omitempty"`// Total de respostas coletadas
    TaxaParticipacao float64                `json:"taxa_participacao,omitempty"`// Taxa de participação média
    Metricas         map[string]interface{} `json:"metricas,omitempty"`      // Métricas adicionais, opcionais
}

// DashboardDataResponse inclui o dashboard com dados e métricas adicionais.
type DashboardDataResponse struct {
    Dashboard DashboardResponse `json:"dashboard"` // Dashboard principal
    Data      interface{}       `json:"data"`      // Dados detalhados do dashboard
    Metrics   interface{}       `json:"metrics"`   // Métricas adicionais
    UpdatedAt time.Time         `json:"updated_at"`// Data da última atualização
}

// ExportResponse representa os dados de um arquivo exportado do dashboard.
type ExportResponse struct {
    Format    string    `json:"format"`                // Formato do arquivo (ex: CSV, PDF)
    Filename  string    `json:"filename"`              // Nome do arquivo
    URL       string    `json:"url,omitempty"`         // URL para download, opcional
    Size      int64     `json:"size,omitempty"`        // Tamanho do arquivo, opcional
    CreatedAt time.Time `json:"created_at"`            // Data de criação do arquivo
}

// MetricsResponse agrega métricas de respostas e estatísticas do dashboard.
type MetricsResponse struct {
    TotalRespostas    int                    `json:"total_respostas"`             // Total de respostas
    DataUltimaResp    *time.Time             `json:"data_ultima_resposta,omitempty"` // Data da última resposta, opcional
    TaxaParticipacao  float64                `json:"taxa_participacao"`           // Taxa de participação média
    MetricasPorTipo   map[string]interface{} `json:"metricas_por_tipo"`           // Métricas detalhadas por tipo
    ResumoEstatistico map[string]interface{} `json:"resumo_estatistico"`          // Resumo estatístico geral
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
