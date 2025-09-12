package response

import (
    "time"
    "organizational-climate-survey/backend/internal/domain/entity"
)

type DashboardResponse struct {
    ID               int                    `json:"id_dashboard"`
    IDPesquisa       int                    `json:"id_pesquisa"`
    Titulo           string                 `json:"titulo"`
    DataCriacao      time.Time              `json:"data_criacao"`
    ConfigFiltros    *string                `json:"config_filtros,omitempty"`
    TotalRespostas   int                    `json:"total_respostas,omitempty"`
    TaxaParticipacao float64                `json:"taxa_participacao,omitempty"`
    Metricas         map[string]interface{} `json:"metricas,omitempty"`
}

type DashboardDataResponse struct {
    Dashboard DashboardResponse `json:"dashboard"`
    Data      interface{}       `json:"data"`
    Metrics   interface{}       `json:"metrics"`
    UpdatedAt time.Time         `json:"updated_at"`
}

type ExportResponse struct {
    Format    string    `json:"format"`
    Filename  string    `json:"filename"`
    URL       string    `json:"url,omitempty"`
    Size      int64     `json:"size,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

type MetricsResponse struct {
    TotalRespostas     int                    `json:"total_respostas"`
    DataUltimaResp     *time.Time             `json:"data_ultima_resposta,omitempty"`
    TaxaParticipacao   float64                `json:"taxa_participacao"`
    MetricasPorTipo    map[string]interface{} `json:"metricas_por_tipo"`
    ResumoEstatistico  map[string]interface{} `json:"resumo_estatistico"`
}

func ToDashboardResponse(dashboard *entity.Dashboard) DashboardResponse {
    return DashboardResponse{
        ID:               dashboard.ID,
        IDPesquisa:       dashboard.IDPesquisa,
        Titulo:           dashboard.Titulo,
        DataCriacao:      dashboard.DataCriacao,
        ConfigFiltros:    dashboard.ConfigFiltros, // Já é ponteiro
        TotalRespostas:   dashboard.TotalRespostas,
        TaxaParticipacao: dashboard.TaxaParticipacao,
        Metricas:         dashboard.Metricas,
    }
}