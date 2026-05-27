// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

import "time"

// Resposta representa uma resposta individual dada a uma pergunta de pesquisa
type Resposta struct {
	ID            int       `json:"id_resposta" example:"1001"`                                         // Identificador único da resposta
	IDPergunta    int       `json:"id_pergunta" example:"15"`                                           // ID da pergunta respondida
	ValorResposta string    `json:"valor_resposta" example:"5"`                                         // Conteúdo da resposta
	DataSubmissao time.Time `json:"data_submissao" swaggertype:"string" example:"2026-01-11T14:30:00Z"` // Data em que resposta foi finalizada
	IDSubmissao   int       `json:"id_submissao" example:"300"`                                         // NOVO - vincula ao respondente anônimo

	// Relacionamentos (opcional, para carregamento sob demanda)
	Pergunta  *Pergunta          `json:"pergunta,omitempty" example:"{\"id_pergunta\":15,\"texto_pergunta\":\"Como você avalia o ambiente?\"}"` // Dados da pergunta
	Pesquisa  *Pesquisa          `json:"pesquisa,omitempty" example:"{\"id_pesquisa\":10,\"titulo\":\"Pesquisa de Clima 2026\"}"`               // Dados da pesquisa
	Submissao *SubmissaoPesquisa `json:"submissao,omitempty" example:"{\"id_submissao\":300,\"status\":\"completa\"}"`
}
