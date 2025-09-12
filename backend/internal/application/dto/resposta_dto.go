package dto

import "organizational-climate-survey/backend/internal/domain/entity"

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