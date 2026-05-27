// Package response contém structs usadas para enviar dados da API como respostas.
// PerguntaResponse representa a estrutura de resposta de uma pergunta dentro de uma pesquisa.
package response

// PerguntaResponse retorna informações detalhadas sobre uma pergunta específica.
type PerguntaResponse struct {
	ID             int                    `json:"id_pergunta" example:"15"`                                          // ID único da pergunta
	TextoPergunta  string                 `json:"texto_pergunta" example:"Como você avalia o ambiente de trabalho?"` // Texto da pergunta
	TipoPergunta   string                 `json:"tipo_pergunta" example:"EscalaNumerica"`                            // Tipo da pergunta (MultiplaEscolha, RespostaAberta, EscalaNumerica, SimNao)
	OrdemExibicao  int                    `json:"ordem_exibicao" example:"1"`                                        // Posição da pergunta na pesquisa
	OpcoesResposta *string                `json:"opcoes_resposta" example:"1,2,3,4,5"`                               // Opções de resposta, se aplicável (para múltipla escolha)
	TotalRespostas int                    `json:"total_respostas,omitempty" example:"250"`                           // Total de respostas recebidas, opcional
	Estatisticas   map[string]interface{} `json:"estatisticas,omitempty" example:"{\"media\":4.2,\"desvio\":0.7}"`   // Estatísticas agregadas da pergunta, opcional
}
