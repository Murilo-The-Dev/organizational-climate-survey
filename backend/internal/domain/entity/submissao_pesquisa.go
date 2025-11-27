package entity

import "time"

// SubmissaoPesquisa representa uma tentativa única e anônima de responder uma pesquisa
type SubmissaoPesquisa struct {
    ID              int        `json:"id_submissao"`      // Identificador único
    IDPesquisa      int        `json:"id_pesquisa"`       // Pesquisa sendo respondida
    TokenAcesso     string     `json:"token_acesso"`      // Token único de acesso
    IPHash          string     `json:"ip_hash"`           // Hash do IP (anti-spam, não identificação)
    FingerprintHash string     `json:"fingerprint_hash"`  // Hash do fingerprint do browser
    Status          string     `json:"status"`            // pendente, completa, expirada
    DataCriacao     time.Time  `json:"data_criacao"`      // Quando token foi gerado
    DataExpiracao   time.Time  `json:"data_expiracao"`    // Quando token expira
    DataConclusao   *time.Time `json:"data_conclusao"`    // Quando foi finalizada
    
    // Relacionamentos
    Pesquisa  *Pesquisa   `json:"pesquisa,omitempty"`
    Respostas []Resposta  `json:"respostas,omitempty"`
}