package response

import "time"

type LoginResponse struct {
    Token     string   `json:"token"`
    ExpiresIn int64    `json:"expires_in"`
    User      UserInfo `json:"user"`
}

type UserInfo struct {
    ID        int    `json:"id"`
    Nome      string `json:"nome"`
    Email     string `json:"email"`
    EmpresaID int    `json:"empresa_id"`
    Status    string `json:"status"`
}

type RefreshTokenResponse struct {
    Token     string `json:"token"`
    ExpiresIn int64  `json:"expires_in"`
}

type TokenValidationResponse struct {
    Valid     bool      `json:"valid"`
    User      UserInfo  `json:"user"`
    ExpiresAt time.Time `json:"expires_at"`
}