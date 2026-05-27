// Package stats define estruturas de resposta utilizadas para retornar
// estatísticas agregadas de empresas, setores e pesquisas.
// Estes DTOs são voltados para visualizações analíticas em dashboards.

package stats

import "time"

// EstatisticasEmpresaResponse agrega métricas de uma empresa específica,
// incluindo contagem de pesquisas, respostas e taxa média de participação.
type EstatisticasEmpresaResponse struct {
	IDEmpresa             int                         `json:"id_empresa" example:"1"`
	NomeEmpresa           string                      `json:"nome_empresa" example:"Acme"`
	TotalPesquisas        int                         `json:"total_pesquisas" example:"8"`
	PesquisasAtivas       int                         `json:"pesquisas_ativas" example:"2"`
	PesquisasConcluidas   int                         `json:"pesquisas_concluidas" example:"6"`
	TotalRespostas        int                         `json:"total_respostas" example:"450"`
	TaxaParticipacaoMedia float64                     `json:"taxa_participacao_media" example:"78.5"`
	SetorMaisAtivo        string                      `json:"setor_mais_ativo" example:"RH"`
	UltimaPesquisaCriada  *time.Time                  `json:"ultima_pesquisa_criada" swaggertype:"string" example:"2026-01-15T10:00:00Z"`
	EstatisticasPorSetor  []EstatisticasSetorResponse `json:"estatisticas_por_setor,omitempty" example:"[{\"id_setor\":2,\"nome_setor\":\"RH\",\"total_pesquisas\":3}]"`
}

// EstatisticasSetorResponse agrega métricas específicas de um setor dentro
// de uma empresa, incluindo contagem de pesquisas e taxa de participação.
type EstatisticasSetorResponse struct {
	IDSetor          int     `json:"id_setor" example:"2"`
	NomeSetor        string  `json:"nome_setor" example:"RH"`
	TotalPesquisas   int     `json:"total_pesquisas" example:"3"`
	TotalRespostas   int     `json:"total_respostas" example:"180"`
	TaxaParticipacao float64 `json:"taxa_participacao" example:"81.3"`
}

// EstatisticasPesquisaResponse agrega métricas detalhadas de uma pesquisa,
// incluindo distribuição de respostas por tipo e médias por pergunta.
type EstatisticasPesquisaResponse struct {
	IDPesquisa        int                    `json:"id_pesquisa" example:"10"`
	Titulo            string                 `json:"titulo" example:"Pesquisa de Clima 2026"`
	TotalPerguntas    int                    `json:"total_perguntas" example:"20"`
	TotalRespostas    int                    `json:"total_respostas" example:"300"`
	TaxaParticipacao  float64                `json:"taxa_participacao" example:"75.0"`
	MediaPorPergunta  float64                `json:"media_por_pergunta" example:"4.1"`
	RespostasPorTipo  map[string]int         `json:"respostas_por_tipo" example:"{\"EscalaNumerica\":250,\"SimNao\":180}"`
	DistribuicaoNotas map[string]int         `json:"distribuicao_notas,omitempty" example:"{\"1\":5,\"2\":20,\"3\":80,\"4\":120,\"5\":75}"`
	Tendencias        map[string]interface{} `json:"tendencias,omitempty" example:"{\"variacao_media\":0.3}"`
}

// EstatisticasGeraisResponse fornece uma visão macro do sistema,
// agregando métricas globais como total de empresas, usuários e pesquisas,
// bem como listas das entidades mais ativas.
type EstatisticasGeraisResponse struct {
	TotalEmpresas          int      `json:"total_empresas" example:"12"`
	TotalUsuarios          int      `json:"total_usuarios" example:"48"`
	TotalPesquisas         int      `json:"total_pesquisas" example:"40"`
	TotalRespostas         int      `json:"total_respostas" example:"3200"`
	TaxaParticipacaoGlobal float64  `json:"taxa_participacao_global" example:"73.8"`
	PesquisasMaisAtivas    []string `json:"pesquisas_mais_ativas" example:"Pesquisa de Clima 2026,Pulso Semanal"`
	EmpresasMaisAtivas     []string `json:"empresas_mais_ativas" example:"Acme,Globex"`
}
