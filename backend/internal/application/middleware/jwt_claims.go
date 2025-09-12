package middleware

import (
    "github.com/golang-jwt/jwt/v5"
)

// JWTClaims estrutura das claims do JWT
type JWTClaims struct {
    UserID    int    `json:"user_id"`
    EmpresaID int    `json:"empresa_id"`
    Email     string `json:"email"`
    jwt.RegisteredClaims
}