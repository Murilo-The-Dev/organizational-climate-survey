package response

type PerguntaResponse struct {
	ID               int                    `json:"id_pergunta"`
	TextoPergunta    string                 `json:"texto_pergunta"`
	TipoPergunta     string                 `json:"tipo_pergunta"`
	OrdemExibicao    int                    `json:"ordem_exibicao"`
	OpcoesResposta   *string                `json:"opcoes_resposta"`
	TotalRespostas   int                    `json:"total_respostas,omitempty"`
	Estatisticas     map[string]interface{} `json:"estatisticas,omitempty"`
}