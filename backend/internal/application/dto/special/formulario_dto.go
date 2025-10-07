// Package special contém structs usadas para operações especiais da API,
// principalmente para submissão de respostas completas de formulários e criação
// em lote de perguntas associadas a uma pesquisa.
package special

import "organizational-climate-survey/backend/internal/domain/entity"

// FormularioRespostaRequest é usado para enviar todas as respostas de uma pesquisa
// de uma só vez, garantindo consistência no envio.
type FormularioRespostaRequest struct {
	IDPesquisa int                     `json:"id_pesquisa" binding:"required,gt=0"` // ID da pesquisa à qual as respostas pertencem
	Respostas  []RespostaCreateRequest `json:"respostas" binding:"required,dive"`  // Lista de respostas a cada pergunta da pesquisa
}

// RespostaCreateRequest representa o envio de uma única resposta de pergunta.
type RespostaCreateRequest struct {
	IDPergunta    int    `json:"id_pergunta" binding:"required,gt=0"`  // ID da pergunta
	ValorResposta string `json:"valor_resposta" binding:"required,max=2000"` // Valor fornecido para a pergunta
}

// ToEntity converte a requisição para a entidade de domínio correspondente.
func (r *RespostaCreateRequest) ToEntity() *entity.Resposta {
	return &entity.Resposta{
		IDPergunta:    r.IDPergunta,
		ValorResposta: r.ValorResposta,
	}
}

// PerguntaBatchCreateRequest permite criar múltiplas perguntas em uma única requisição.
type PerguntaBatchCreateRequest struct {
	IDPesquisa int                     `json:"id_pesquisa" binding:"required,gt=0"` // ID da pesquisa
	Perguntas  []PerguntaCreateRequest `json:"perguntas" binding:"required,dive"`   // Lista de perguntas a serem criadas
}

// PerguntaCreateRequest define os dados necessários para criar uma única pergunta.
type PerguntaCreateRequest struct {
	TextoPergunta  string  `json:"texto_pergunta" binding:"required,min=5,max=500"` // Texto da pergunta
	TipoPergunta   string  `json:"tipo_pergunta" binding:"required,oneof=MultiplaEscolha RespostaAberta EscalaNumerica SimNao"` // Tipo da pergunta
	OrdemExibicao  int     `json:"ordem_exibicao" binding:"required,gte=1"` // Ordem de exibição na pesquisa
	OpcoesResposta *string `json:"opcoes_resposta,omitempty"`               // Opções para perguntas de múltipla escolha
}
