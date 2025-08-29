package entity

import "time"

type Dashboard struct {
    ID           int       `json:"id_dashboard"`
    IDPesquisa   int       `json:"id_pesquisa"` // Relação 1:1 com pesquisa
    Titulo       string    `json:"titulo"`
    DataCriacao  time.Time `json:"data_criacao"`
    ConfigFiltros *string   `json:"config_filtros"` // JSON com configurações
    
    // Dados agregados - calculados em runtime
    TotalRespostas   int                    `json:"total_respostas,omitempty"`
    TaxaParticipacao float64                `json:"taxa_participacao,omitempty"`
    Metricas         map[string]interface{} `json:"metricas,omitempty"`
}