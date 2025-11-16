// Package response contém structs de resposta relacionadas a autenticação e tokens.
package response

import "time"

// LoginResponse representa a resposta enviada após login bem-sucedido.
type LoginResponse struct {
    Token     string   `json:"token"`       // JWT ou token de autenticação
    ExpiresIn int64    `json:"expires_in"`  // Tempo de expiração em segundos
    User      UserInfo `json:"user"`        // Informações básicas do usuário logado
}

// UserInfo mantém dados do usuário.
type UserInfo struct {
    ID        int    `json:"id"`          // ID do usuário
    Nome      string `json:"nome"`        // Nome completo
    Email     string `json:"email"`       // E-mail
    EmpresaID int    `json:"empresa_id"`  // Empresa associada
    Status    string `json:"status"`      // Status do usuário (Ativo/Inativo)
}

// RefreshTokenResponse representa a resposta ao renovar um token.
type RefreshTokenResponse struct {
    Token     string `json:"token"`       // Novo token de autenticação
    ExpiresIn int64  `json:"expires_in"`  // Expiração em segundos
}

// TokenValidationResponse é retornado ao validar um token.
type TokenValidationResponse struct {
    Valid     bool      `json:"valid"`      // Indica se o token é válido
    User      UserInfo  `json:"user"`       // Dados do usuário associado ao token
    ExpiresAt time.Time `json:"expires_at"` // Data e hora de expiração do token
}
