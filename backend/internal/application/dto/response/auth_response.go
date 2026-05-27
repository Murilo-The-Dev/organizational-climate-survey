// Package response contém structs de resposta relacionadas a autenticação e tokens.
package response

import "time"

// LoginResponse representa a resposta enviada após login bem-sucedido.
type LoginResponse struct {
	Token     string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT ou token de autenticação
	ExpiresIn int64    `json:"expires_in" example:"86400"`                              // Tempo de expiração em segundos
	User      UserInfo `json:"user" example:"{\"id\":1,\"nome\":\"Joao Silva\"}"`       // Informações básicas do usuário logado
}

// UserInfo mantém dados do usuário.
type UserInfo struct {
	ID        int    `json:"id" example:"1"`                   // ID do usuário
	Nome      string `json:"nome" example:"Joao Silva"`        // Nome completo
	Email     string `json:"email" example:"joao@empresa.com"` // E-mail
	EmpresaID int    `json:"empresa_id" example:"1"`           // Empresa associada
	Status    string `json:"status" example:"Ativo"`           // Status do usuário (Ativo/Inativo)
}

// RefreshTokenResponse representa a resposta ao renovar um token.
type RefreshTokenResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Novo token de autenticação
	ExpiresIn int64  `json:"expires_in" example:"86400"`                              // Expiração em segundos
}

// TokenValidationResponse é retornado ao validar um token.
type TokenValidationResponse struct {
	Valid     bool      `json:"valid" example:"true"`                                           // Indica se o token é válido
	User      UserInfo  `json:"user" example:"{\"id\":1,\"nome\":\"Joao Silva\"}"`              // Dados do usuário associado ao token
	ExpiresAt time.Time `json:"expires_at" swaggertype:"string" example:"2026-01-31T23:59:59Z"` // Data e hora de expiração do token
}
