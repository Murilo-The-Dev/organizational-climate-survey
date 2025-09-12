package auth

type LoginRequest struct {
    Email string `json:"email" binding:"required,email"`
    Senha string `json:"senha" binding:"required"`
}

type RefreshTokenRequest struct {
    Token string `json:"token" binding:"required"`
}

type ValidateTokenRequest struct {
    Token string `json:"token" binding:"required"`
}

type ChangePasswordRequest struct {
    SenhaAtual string `json:"senha_atual" binding:"required"`
    NovaSenha  string `json:"nova_senha" binding:"required,min=8,max=128"`
}

type ForgotPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}