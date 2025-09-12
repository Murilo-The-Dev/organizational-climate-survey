package response

import "time"

type PesquisaResponse struct {
	ID                   int                            `json:"id_pesquisa"`
	IDEmpresa            int                            `json:"id_empresa"`
	IDSetor              int                            `json:"id_setor"`
	Titulo               string                         `json:"titulo"`
	Descricao            string                         `json:"descricao"`
	DataCriacao          time.Time                      `json:"data_criacao"`
	DataAbertura         *time.Time                     `json:"data_abertura"`
	DataFechamento       *time.Time                     `json:"data_fechamento"`
	Status               string                         `json:"status"`
	LinkAcesso           string                         `json:"link_acesso"`
	QRCodePath           string                         `json:"qrcode_path"`
	Anonimato            bool                           `json:"anonimato"`
	TotalPerguntas       int                            `json:"total_perguntas,omitempty"`
	TotalRespostas       int                            `json:"total_respostas,omitempty"`
	TaxaParticipacao     float64                        `json:"taxa_participacao,omitempty"`
	UsuarioAdministrador *UsuarioAdministradorResponse  `json:"usuario_administrador,omitempty"`
	Setor                *SetorResponse                 `json:"setor,omitempty"`
	Perguntas            []PerguntaResponse             `json:"perguntas,omitempty"`
}