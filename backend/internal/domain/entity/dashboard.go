// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para dashboards do sistema de pesquisa.
package entity

import "time"

// Dashboard apresenta métricas e análises consolidadas de uma pesquisa
type Dashboard struct {
	ID            int       `json:"id_dashboard" example:"3"`                                         // Identificador único do dashboard
	IDPesquisa    int       `json:"id_pesquisa" example:"10"`                                         // ID da pesquisa associada (relação 1:1)
	Titulo        string    `json:"titulo" example:"Dashboard Clima RH"`                              // Título do painel
	DataCriacao   time.Time `json:"data_criacao" swaggertype:"string" example:"2026-01-05T08:00:00Z"` // Data de criação
	ConfigFiltros *string   `json:"config_filtros" example:"{\"setor\":\"RH\"}"`                      // JSON com filtros aplicáveis

	// Dados agregados calculados dinamicamente
	TotalRespostas   int                    `json:"total_respostas,omitempty" example:"250"`             // Quantidade de respostas coletadas
	TaxaParticipacao float64                `json:"taxa_participacao,omitempty" example:"82.4"`          // Percentual de participação
	Metricas         map[string]interface{} `json:"metricas,omitempty" example:"{\"engajamento\":82.4}"` // Indicadores customizados
}
