// Package response define objetos de transferência de dados para respostas HTTP.
package response

// GenerateTokenResponse representa resposta com token gerado
type GenerateTokenResponse struct {
	TokenAcesso string `json:"token_acesso"` // Token único para submissão
	ExpiresAt   string `json:"expires_at"`   // Data/hora de expiração (ISO 8601)
	ExpiresIn   int    `json:"expires_in"`   // Tempo até expirar em segundos
}

// SubmissionStatsResponse representa estatísticas de submissões de uma pesquisa
type SubmissionStatsResponse struct {
	TotalSubmissoes    int     `json:"total_submissoes"`     // Total de submissões (completas + pendentes)
	Completas          int     `json:"completas"`            // Submissões finalizadas
	Pendentes          int     `json:"pendentes"`            // Submissões aguardando conclusão
	Expiradas          int     `json:"expiradas"`            // Submissões que expiraram
	TaxaConclusao      float64 `json:"taxa_conclusao"`       // Percentual de conclusão
	ParticipantesUnicos int    `json:"participantes_unicos"` // Respondentes únicos (= completas)
}