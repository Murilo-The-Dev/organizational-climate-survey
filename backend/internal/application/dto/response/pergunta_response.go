// Package response contém structs usadas para enviar dados da API como respostas.
// PerguntaResponse representa a estrutura de resposta de uma pergunta dentro de uma pesquisa.
package response

// PerguntaResponse retorna informações detalhadas sobre uma pergunta específica.
type PerguntaResponse struct {
	ID             int                    `json:"id_pergunta"`              // ID único da pergunta
	TextoPergunta  string                 `json:"texto_pergunta"`           // Texto da pergunta
	TipoPergunta   string                 `json:"tipo_pergunta"`            // Tipo da pergunta (MultiplaEscolha, RespostaAberta, EscalaNumerica, SimNao)
	OrdemExibicao  int                    `json:"ordem_exibicao"`           // Posição da pergunta na pesquisa
	OpcoesResposta *string                `json:"opcoes_resposta"`          // Opções de resposta, se aplicável (para múltipla escolha)
	TotalRespostas int                    `json:"total_respostas,omitempty"`// Total de respostas recebidas, opcional
	Estatisticas   map[string]interface{} `json:"estatisticas,omitempty"`   // Estatísticas agregadas da pergunta, opcional
}
