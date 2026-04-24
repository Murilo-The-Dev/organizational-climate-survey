// Package response define objetos de transferência de dados para respostas HTTP.
package response

// GenerateTokenResponse representa resposta com token gerado
type GenerateTokenResponse struct {
	TokenAcesso string `json:"token_acesso" example:"tok_abc123def456"`   // Token único para submissão
	ExpiresAt   string `json:"expires_at" example:"2026-01-10T12:00:00Z"` // Data/hora de expiração (ISO 8601)
	ExpiresIn   int    `json:"expires_in" example:"1800"`                 // Tempo até expirar em segundos
}

// SubmissionStatsResponse representa estatísticas de submissões de uma pesquisa
type SubmissionStatsResponse struct {
	TotalSubmissoes     int     `json:"total_submissoes" example:"120"`     // Total de submissões (completas + pendentes)
	Completas           int     `json:"completas" example:"100"`            // Submissões finalizadas
	Pendentes           int     `json:"pendentes" example:"15"`             // Submissões aguardando conclusão
	Expiradas           int     `json:"expiradas" example:"5"`              // Submissões que expiraram
	TaxaConclusao       float64 `json:"taxa_conclusao" example:"83.3"`      // Percentual de conclusão
	ParticipantesUnicos int     `json:"participantes_unicos" example:"100"` // Respondentes únicos (= completas)
}
