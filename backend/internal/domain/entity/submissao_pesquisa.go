package entity

import "time"

// SubmissaoPesquisa representa uma tentativa única e anônima de responder uma pesquisa
type SubmissaoPesquisa struct {
	ID              int        `json:"id_submissao" example:"300"`                                         // Identificador único
	IDPesquisa      int        `json:"id_pesquisa" example:"10"`                                           // Pesquisa sendo respondida
	TokenAcesso     string     `json:"token_acesso" example:"tok_abc123def456"`                            // Token único de acesso
	IPHash          string     `json:"ip_hash" example:"sha256:abc..."`                                    // Hash do IP (anti-spam, não identificação)
	FingerprintHash string     `json:"fingerprint_hash" example:"sha256:def..."`                           // Hash do fingerprint do browser
	UserAgentHash   string     `json:"user_agent_hash" example:"sha256:ghi..."`                            // Hash do User-Agent
	AcceptLangHash  string     `json:"accept_lang_hash" example:"sha256:jkl..."`                           // Hash do Accept-Language
	FingerprintComp string     `json:"fingerprint_comp" example:"fp_comp_01"`                              // Fingerprint composto para heurística antifraude
	Status          string     `json:"status" example:"completa"`                                          // pendente, completa, expirada
	DataCriacao     time.Time  `json:"data_criacao" swaggertype:"string" example:"2026-01-10T11:30:00Z"`   // Quando token foi gerado
	DataExpiracao   time.Time  `json:"data_expiracao" swaggertype:"string" example:"2026-01-10T12:00:00Z"` // Quando token expira
	DataConclusao   *time.Time `json:"data_conclusao" swaggertype:"string" example:"2026-01-10T11:45:00Z"` // Quando foi finalizada

	// Relacionamentos
	Pesquisa  *Pesquisa  `json:"pesquisa,omitempty" example:"{\"id_pesquisa\":10,\"titulo\":\"Pesquisa de Clima 2026\"}"`
	Respostas []Resposta `json:"respostas,omitempty" example:"[{\"id_resposta\":1001,\"valor_resposta\":\"5\"}]"`
}
