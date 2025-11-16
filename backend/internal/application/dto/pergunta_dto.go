// Package dto contém estruturas de transferência de dados (Data Transfer Objects)
// utilizadas para comunicação entre as camadas de entrada e o domínio.
// Este arquivo define os DTOs de criação e atualização de perguntas associadas a uma pesquisa.

package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

// PerguntaCreateRequest representa os dados necessários para criar uma nova pergunta
// vinculada a uma pesquisa específica.
type PerguntaCreateRequest struct {
	IDPesquisa     int     `json:"id_pesquisa" binding:"required,gt=0"`                                                      // Identificador da pesquisa associada (obrigatório)
	TextoPergunta  string  `json:"texto_pergunta" binding:"required,min=5,max=500"`                                          // Enunciado da pergunta (obrigatório)
	TipoPergunta   string  `json:"tipo_pergunta" binding:"required,oneof=MultiplaEscolha RespostaAberta EscalaNumerica SimNao"` // Tipo da pergunta, restringido a opções válidas
	OrdemExibicao  int     `json:"ordem_exibicao" binding:"required,gte=1"`                                                  // Posição de exibição da pergunta (obrigatório)
	OpcoesResposta *string `json:"opcoes_resposta,omitempty"`                                                                // Opções disponíveis para múltipla escolha ou escala (opcional)
}

// PerguntaUpdateRequest representa os campos permitidos para atualização parcial
// de uma pergunta existente.
type PerguntaUpdateRequest struct {
	TextoPergunta  *string `json:"texto_pergunta,omitempty" binding:"omitempty,min=5,max=500"`                                          // Novo texto da pergunta (opcional)
	TipoPergunta   *string `json:"tipo_pergunta,omitempty" binding:"omitempty,oneof=MultiplaEscolha RespostaAberta EscalaNumerica SimNao"` // Novo tipo da pergunta (opcional)
	OrdemExibicao  *int    `json:"ordem_exibicao,omitempty" binding:"omitempty,gte=1"`                                                  // Nova ordem de exibição (opcional)
	OpcoesResposta *string `json:"opcoes_resposta,omitempty"`                                                                           // Novas opções de resposta (opcional)
}

// ToEntity converte a requisição de criação em uma entidade de domínio Pergunta,
// sanitizando entradas textuais.
func (r *PerguntaCreateRequest) ToEntity() *entity.Pergunta {
	return &entity.Pergunta{
		IDPesquisa:     r.IDPesquisa,
		TextoPergunta:  strings.TrimSpace(r.TextoPergunta),
		TipoPergunta:   r.TipoPergunta,
		OrdemExibicao:  r.OrdemExibicao,
		OpcoesResposta: r.OpcoesResposta,
	}
}

// ApplyToEntity aplica os campos fornecidos na requisição de atualização
// sobre uma instância existente da entidade Pergunta.
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
