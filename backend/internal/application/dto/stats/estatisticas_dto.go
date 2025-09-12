package stats

import "time"

type EstatisticasEmpresaResponse struct {
	IDEmpresa               int                      `json:"id_empresa"`
	NomeEmpresa             string                   `json:"nome_empresa"`
	TotalPesquisas          int                      `json:"total_pesquisas"`
	PesquisasAtivas         int                      `json:"pesquisas_ativas"`
	PesquisasConcluidas     int                      `json:"pesquisas_concluidas"`
	TotalRespostas          int                      `json:"total_respostas"`
	TaxaParticipacaoMedia   float64                  `json:"taxa_participacao_media"`
	SetorMaisAtivo          string                   `json:"setor_mais_ativo"`
	UltimaPesquisaCriada    *time.Time               `json:"ultima_pesquisa_criada"`
	EstatisticasPorSetor    []EstatisticasSetorResponse `json:"estatisticas_por_setor,omitempty"`
}

type EstatisticasSetorResponse struct {
	IDSetor            int     `json:"id_setor"`
	NomeSetor          string  `json:"nome_setor"`
	TotalPesquisas     int     `json:"total_pesquisas"`
	TotalRespostas     int     `json:"total_respostas"`
	TaxaParticipacao   float64 `json:"taxa_participacao"`
}

type EstatisticasPesquisaResponse struct {
	IDPesquisa         int                    `json:"id_pesquisa"`
	Titulo             string                 `json:"titulo"`
	TotalPerguntas     int                    `json:"total_perguntas"`
	TotalRespostas     int                    `json:"total_respostas"`
	TaxaParticipacao   float64                `json:"taxa_participacao"`
	MediaPorPergunta   float64                `json:"media_por_pergunta"`
	RespostasPorTipo   map[string]int         `json:"respostas_por_tipo"`
	DistribuicaoNotas  map[string]int         `json:"distribuicao_notas,omitempty"`
	Tendencias         map[string]interface{} `json:"tendencias,omitempty"`
}

type EstatisticasGeraisResponse struct {
	TotalEmpresas           int     `json:"total_empresas"`
	TotalUsuarios           int     `json:"total_usuarios"`
	TotalPesquisas          int     `json:"total_pesquisas"`
	TotalRespostas          int     `json:"total_respostas"`
	TaxaParticipacaoGlobal  float64 `json:"taxa_participacao_global"`
	PesquisasMaisAtivas     []string `json:"pesquisas_mais_ativas"`
	EmpresasMaisAtivas      []string `json:"empresas_mais_ativas"`
}