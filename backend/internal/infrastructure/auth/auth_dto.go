// Package auth implementa DTOs e handlers para autenticação e autorização.
// Fornece estruturas para requisições de login, tokens e gerenciamento de senha.
package auth

// LoginRequest representa os dados necessários para autenticação
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"` // Email do usuário
	Senha string `json:"senha" binding:"required"`       // Senha em texto plano
}

// RefreshTokenRequest representa pedido de renovação de token
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"` // Token JWT atual
}

// ValidateTokenRequest representa pedido de validação de token
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"` // Token JWT a ser validado
}

// ChangePasswordRequest representa pedido de alteração de senha
type ChangePasswordRequest struct {
	SenhaAtual string `json:"senha_atual" binding:"required"`              // Senha atual
	NovaSenha  string `json:"nova_senha" binding:"required,min=8,max=128"` // Nova senha
}

// ForgotPasswordRequest representa pedido de recuperação de senha
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"` // Email para recuperação
}
