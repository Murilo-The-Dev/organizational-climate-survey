// Package response contém structs usadas para enviar dados da API como respostas.
// PesquisaResponse representa a estrutura de resposta detalhada de uma pesquisa.
package response

import "time"

// PesquisaResponse retorna informações completas sobre uma pesquisa, incluindo
// dados de setor, administrador, perguntas e métricas de participação.
type PesquisaResponse struct {
	ID                   int                           `json:"id_pesquisa"`                       // ID único da pesquisa
	IDEmpresa            int                           `json:"id_empresa"`                         // ID da empresa associada
	IDSetor              int                           `json:"id_setor"`                           // ID do setor associado
	Titulo               string                        `json:"titulo"`                             // Título da pesquisa
	Descricao            string                        `json:"descricao"`                          // Descrição da pesquisa
	DataCriacao          time.Time                     `json:"data_criacao"`                       // Data de criação da pesquisa
	DataAbertura         *time.Time                    `json:"data_abertura"`                      // Data de abertura da pesquisa, opcional
	DataFechamento       *time.Time                    `json:"data_fechamento"`                    // Data de fechamento da pesquisa, opcional
	Status               string                        `json:"status"`                             // Status da pesquisa (Rascunho, Ativa, Concluída, Arquivada)
	LinkAcesso           string                        `json:"link_acesso"`                        // Link de acesso à pesquisa
	QRCodePath           string                        `json:"qrcode_path"`                        // Caminho para QR Code da pesquisa
	Anonimato            bool                          `json:"anonimato"`                          // Indica se a pesquisa é anônima
	TotalPerguntas       int                           `json:"total_perguntas,omitempty"`          // Número total de perguntas, opcional
	TotalRespostas       int                           `json:"total_respostas,omitempty"`          // Número total de respostas, opcional
	TaxaParticipacao     float64                       `json:"taxa_participacao,omitempty"`        // Taxa média de participação, opcional
	UsuarioAdministrador *UsuarioAdministradorResponse `json:"usuario_administrador,omitempty"`   // Informações do administrador da pesquisa, opcional
	Setor                *SetorResponse                `json:"setor,omitempty"`                    // Informações do setor da pesquisa, opcional
	Perguntas            []PerguntaResponse            `json:"perguntas,omitempty"`                // Lista de perguntas da pesquisa, opcional
}
