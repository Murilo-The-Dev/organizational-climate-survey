// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para perguntas do sistema de pesquisa.
package entity

// Pergunta representa uma questão individual dentro de uma pesquisa de clima
type Pergunta struct {
	ID             int     `json:"id_pergunta" example:"15"`                                          // Identificador único da pergunta
	IDPesquisa     int     `json:"id_pesquisa" example:"10"`                                          // ID da pesquisa associada
	TextoPergunta  string  `json:"texto_pergunta" example:"Como você avalia o ambiente de trabalho?"` // Texto exibido ao respondente
	TipoPergunta   string  `json:"tipo_pergunta" example:"EscalaNumerica"`                            // Tipo de resposta esperada
	OrdemExibicao  int     `json:"ordem_exibicao" example:"1"`                                        // Sequência de apresentação
	OpcoesResposta *string `json:"opcoes_resposta" example:"1,2,3,4,5"`                               // JSON com opções para múltipla escolha

	// Relacionamento com respostas (carregamento opcional)
	Respostas []Resposta `json:"respostas,omitempty" example:"[{\"id_resposta\":1001,\"valor_resposta\":\"5\"}]"` // Respostas coletadas
}
