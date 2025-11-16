// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

import "time"

// Resposta representa uma resposta individual dada a uma pergunta de pesquisa
type Resposta struct {
	ID            int       `json:"id_resposta"`    // Identificador único da resposta
	IDPergunta    int       `json:"id_pergunta"`    // ID da pergunta respondida
	IDPesquisa    int       `json:"id_pesquisa"`    // ID da pesquisa associada
	ValorResposta string    `json:"valor_resposta"` // Conteúdo da resposta
	DataResposta  time.Time `json:"data_resposta"`  // Data em que resposta foi iniciada
	DataSubmissao time.Time `json:"data_submissao"` // Data em que resposta foi finalizada

	// Relacionamentos (opcional, para carregamento sob demanda)
	Pergunta *Pergunta `json:"pergunta,omitempty"` // Dados da pergunta
	Pesquisa *Pesquisa `json:"pesquisa,omitempty"` // Dados da pesquisa
}
