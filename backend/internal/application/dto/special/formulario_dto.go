package special

import "organizational-climate-survey/backend/internal/domain/entity"

// Para submissão de formulário completo (múltiplas respostas)
type FormularioRespostaRequest struct {
	IDPesquisa int                     `json:"id_pesquisa" binding:"required,gt=0"`
	Respostas  []RespostaCreateRequest `json:"respostas" binding:"required,dive"`
}

type RespostaCreateRequest struct {
	IDPergunta    int    `json:"id_pergunta" binding:"required,gt=0"`
	ValorResposta string `json:"valor_resposta" binding:"required,max=2000"`
}

func (r *RespostaCreateRequest) ToEntity() *entity.Resposta {
	return &entity.Resposta{
		IDPergunta:    r.IDPergunta,
		ValorResposta: r.ValorResposta,
	}
}

// Para operações em lote com perguntas
type PerguntaBatchCreateRequest struct {
	IDPesquisa int                     `json:"id_pesquisa" binding:"required,gt=0"`
	Perguntas  []PerguntaCreateRequest `json:"perguntas" binding:"required,dive"`
}

type PerguntaCreateRequest struct {
	TextoPergunta  string  `json:"texto_pergunta" binding:"required,min=5,max=500"`
	TipoPergunta   string  `json:"tipo_pergunta" binding:"required,oneof=MultiplaEscolha RespostaAberta EscalaNumerica SimNao"`
	OrdemExibicao  int     `json:"ordem_exibicao" binding:"required,gte=1"`
	OpcoesResposta *string `json:"opcoes_resposta,omitempty"`
}