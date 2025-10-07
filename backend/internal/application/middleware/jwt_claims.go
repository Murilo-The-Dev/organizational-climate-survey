// Package middleware fornece componentes intermediários para processamento de requisições.
// Implementa autenticação, autorização e validações de segurança da aplicação.
package middleware

import (
    "github.com/golang-jwt/jwt/v5"
)

// JWTClaims define os dados de autenticação contidos no token JWT
type JWTClaims struct {
    UserID    int    `json:"user_id"`    // ID do usuário autenticado
    EmpresaID int    `json:"empresa_id"` // ID da empresa vinculada
    Email     string `json:"email"`      // Email do usuário
    jwt.RegisteredClaims                 // Claims padrão JWT (exp, iat, iss)
}