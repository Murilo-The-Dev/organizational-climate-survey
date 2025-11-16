// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para dashboards do sistema de pesquisa.
package entity

import "time"

// Dashboard apresenta métricas e análises consolidadas de uma pesquisa
type Dashboard struct {
    ID            int       `json:"id_dashboard"`     // Identificador único do dashboard
    IDPesquisa    int       `json:"id_pesquisa"`      // ID da pesquisa associada (relação 1:1)
    Titulo        string    `json:"titulo"`           // Título do painel
    DataCriacao   time.Time `json:"data_criacao"`     // Data de criação
    ConfigFiltros *string   `json:"config_filtros"`   // JSON com filtros aplicáveis
    
    // Dados agregados calculados dinamicamente
    TotalRespostas   int                    `json:"total_respostas,omitempty"`   // Quantidade de respostas coletadas
    TaxaParticipacao float64                `json:"taxa_participacao,omitempty"` // Percentual de participação
    Metricas         map[string]interface{} `json:"metricas,omitempty"`          // Indicadores customizados
}