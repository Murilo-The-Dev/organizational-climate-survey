// Package stats define estruturas de resposta utilizadas para retornar
// estatísticas agregadas de empresas, setores e pesquisas.
// Estes DTOs são voltados para visualizações analíticas em dashboards.

package stats

import "time"

// EstatisticasEmpresaResponse agrega métricas de uma empresa específica,
// incluindo contagem de pesquisas, respostas e taxa média de participação.
type EstatisticasEmpresaResponse struct {
	IDEmpresa             int                        `json:"id_empresa"`
	NomeEmpresa           string                     `json:"nome_empresa"`
	TotalPesquisas        int                        `json:"total_pesquisas"`
	PesquisasAtivas       int                        `json:"pesquisas_ativas"`
	PesquisasConcluidas   int                        `json:"pesquisas_concluidas"`
	TotalRespostas        int                        `json:"total_respostas"`
	TaxaParticipacaoMedia float64                    `json:"taxa_participacao_media"`
	SetorMaisAtivo        string                     `json:"setor_mais_ativo"`
	UltimaPesquisaCriada  *time.Time                 `json:"ultima_pesquisa_criada"`
	EstatisticasPorSetor  []EstatisticasSetorResponse `json:"estatisticas_por_setor,omitempty"`
}

// EstatisticasSetorResponse agrega métricas específicas de um setor dentro
// de uma empresa, incluindo contagem de pesquisas e taxa de participação.
type EstatisticasSetorResponse struct {
	IDSetor          int     `json:"id_setor"`
	NomeSetor        string  `json:"nome_setor"`
	TotalPesquisas   int     `json:"total_pesquisas"`
	TotalRespostas   int     `json:"total_respostas"`
	TaxaParticipacao float64 `json:"taxa_participacao"`
}

// EstatisticasPesquisaResponse agrega métricas detalhadas de uma pesquisa,
// incluindo distribuição de respostas por tipo e médias por pergunta.
type EstatisticasPesquisaResponse struct {
	IDPesquisa        int                    `json:"id_pesquisa"`
	Titulo            string                 `json:"titulo"`
	TotalPerguntas    int                    `json:"total_perguntas"`
	TotalRespostas    int                    `json:"total_respostas"`
	TaxaParticipacao  float64                `json:"taxa_participacao"`
	MediaPorPergunta  float64                `json:"media_por_pergunta"`
	RespostasPorTipo  map[string]int         `json:"respostas_por_tipo"`
	DistribuicaoNotas map[string]int         `json:"distribuicao_notas,omitempty"`
	Tendencias        map[string]interface{} `json:"tendencias,omitempty"`
}

// EstatisticasGeraisResponse fornece uma visão macro do sistema,
// agregando métricas globais como total de empresas, usuários e pesquisas,
// bem como listas das entidades mais ativas.
type EstatisticasGeraisResponse struct {
	TotalEmpresas          int      `json:"total_empresas"`
	TotalUsuarios          int      `json:"total_usuarios"`
	TotalPesquisas         int      `json:"total_pesquisas"`
	TotalRespostas         int      `json:"total_respostas"`
	TaxaParticipacaoGlobal float64  `json:"taxa_participacao_global"`
	PesquisasMaisAtivas    []string `json:"pesquisas_mais_ativas"`
	EmpresasMaisAtivas     []string `json:"empresas_mais_ativas"`
}
