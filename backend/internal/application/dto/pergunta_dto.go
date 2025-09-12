package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

type PerguntaCreateRequest struct {
	IDPesquisa     int     `json:"id_pesquisa" binding:"required,gt=0"`
	TextoPergunta  string  `json:"texto_pergunta" binding:"required,min=5,max=500"`
	TipoPergunta   string  `json:"tipo_pergunta" binding:"required,oneof=MultiplaEscolha RespostaAberta EscalaNumerica SimNao"`
	OrdemExibicao  int     `json:"ordem_exibicao" binding:"required,gte=1"`
	OpcoesResposta *string `json:"opcoes_resposta,omitempty"`
}

type PerguntaUpdateRequest struct {
	TextoPergunta  *string `json:"texto_pergunta,omitempty" binding:"omitempty,min=5,max=500"`
	TipoPergunta   *string `json:"tipo_pergunta,omitempty" binding:"omitempty,oneof=MultiplaEscolha RespostaAberta EscalaNumerica SimNao"`
	OrdemExibicao  *int    `json:"ordem_exibicao,omitempty" binding:"omitempty,gte=1"`
	OpcoesResposta *string `json:"opcoes_resposta,omitempty"`
}

func (r *PerguntaCreateRequest) ToEntity() *entity.Pergunta {
	return &entity.Pergunta{
		IDPesquisa:     r.IDPesquisa,
		TextoPergunta:  strings.TrimSpace(r.TextoPergunta),
		TipoPergunta:   r.TipoPergunta,
		OrdemExibicao:  r.OrdemExibicao,
		OpcoesResposta: r.OpcoesResposta,
	}
}

func (r *PerguntaUpdateRequest) ApplyToEntity(pergunta *entity.Pergunta) {
	if r.TextoPergunta != nil {
		pergunta.TextoPergunta = strings.TrimSpace(*r.TextoPergunta)
	}
	if r.TipoPergunta != nil {
		pergunta.TipoPergunta = *r.TipoPergunta
	}
	if r.OrdemExibicao != nil {
		pergunta.OrdemExibicao = *r.OrdemExibicao
	}
	if r.OpcoesResposta != nil {
		pergunta.OpcoesResposta = r.OpcoesResposta
	}
}