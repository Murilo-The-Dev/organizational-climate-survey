// Package dto contém estruturas de transferência de dados (Data Transfer Objects)
// usadas para comunicação entre a camada de entrada (handlers) e o domínio.
// Este arquivo define o DTO para criação de respostas em pesquisas.

package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"time"
)

// RespostaCreateRequest representa os dados necessários para registrar
// uma nova resposta associada a uma pergunta de pesquisa.
type RespostaCreateRequest struct {
	IDPergunta    int    `json:"id_pergunta"`    // Identificador da pergunta respondida
	ValorResposta string `json:"valor_resposta"` // Conteúdo textual ou valor da resposta
}

// ToEntity converte o DTO em uma entidade de domínio Resposta,
// pronta para persistência no repositório.
func (r *RespostaCreateRequest) ToEntity() *entity.Resposta {
	return &entity.Resposta{
		IDPergunta:    r.IDPergunta,
		ValorResposta: r.ValorResposta,
		DataSubmissao: time.Now(), // Adicionar timestamp
		// IDSubmissao será setado no UseCase após validar token
	}
}