package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"organizational-climate-survey/backend/internal/application/middleware"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	usuarioUseCase      *usecase.UsuarioAdministradorUseCase
	logAuditoriaUseCase *usecase.LogAuditoriaUseCase
	jwtSecret           []byte
}

func NewAuthHandler(
	usuarioUseCase *usecase.UsuarioAdministradorUseCase,
	logAuditoriaUseCase *usecase.LogAuditoriaUseCase,
	jwtSecret string,
) *AuthHandler {
	return &AuthHandler{
		usuarioUseCase:      usuarioUseCase,
		logAuditoriaUseCase: logAuditoriaUseCase,
		jwtSecret:           []byte(jwtSecret),
	}
}

// Login godoc
// @Summary Autenticar usuário administrador
// @Description Realiza login de usuário administrador e retorna token JWT
// @Tags autenticacao
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Credenciais de login"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validateLoginRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	clientIP := h.getClientIP(r)

	// Tentar autenticar
	usuario, err := h.usuarioUseCase.Authenticate(r.Context(), req.Email, req.Senha, clientIP)
	if err != nil {
		if strings.Contains(err.Error(), "credenciais inválidas") || strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusUnauthorized, "Credenciais inválidas", "Email ou senha incorretos")
			return
		}
		if strings.Contains(err.Error(), "inativo") {
			response.WriteError(w, http.StatusUnauthorized, "Usuário inativo", "Conta desativada. Entre em contato com o administrador")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Gerar token JWT
	token, err := h.generateJWT(usuario.ID, usuario.IDEmpresa, usuario.Email)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro ao gerar token", err.Error())
		return
	}

	// Preparar resposta
	loginResponse := response.LoginResponse{
		Token:     token,
		ExpiresIn: 24 * 60 * 60, // 24 horas em segundos
		User: response.UserInfo{
			ID:        usuario.ID,
			Nome:      usuario.NomeAdmin,
			Email:     usuario.Email,
			EmpresaID: usuario.IDEmpresa,
			Status:    usuario.Status,
		},
	}

	response.WriteSuccess(w, http.StatusOK, "Login realizado com sucesso", loginResponse)
}

// RefreshToken godoc
// @Summary Renovar token JWT
// @Description Renova um token JWT válido
// @Tags autenticacao
// @Accept json
// @Produce json
// @Param token body dto.RefreshTokenRequest true "Token para renovação"
// @Success 200 {object} response.RefreshTokenResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	if strings.TrimSpace(req.Token) == "" {
		response.WriteError(w, http.StatusBadRequest, "Token obrigatório", "Token é obrigatório")
		return
	}

	// Validar token atual
	claims, err := h.validateJWT(req.Token)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Token inválido", err.Error())
		return
	}

	// Verificar se usuário ainda está ativo
	usuario, err := h.usuarioUseCase.GetByID(r.Context(), claims.UserID)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Usuário não encontrado", "Usuário não existe ou foi removido")
		return
	}

	if usuario.Status != "Ativo" {
		response.WriteError(w, http.StatusUnauthorized, "Usuário inativo", "Conta desativada")
		return
	}

	// Gerar novo token
	newToken, err := h.generateJWT(usuario.ID, usuario.IDEmpresa, usuario.Email)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro ao gerar token", err.Error())
		return
	}

	refreshResponse := response.RefreshTokenResponse{
		Token:     newToken,
		ExpiresIn: 24 * 60 * 60, // 24 horas em segundos
	}

	response.WriteSuccess(w, http.StatusOK, "Token renovado com sucesso", refreshResponse)
}

// Logout godoc
// @Summary Realizar logout
// @Description Invalida o token atual (implementação básica)
// @Tags autenticacao
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Em uma implementação completa, aqui seria adicionado o token a uma blacklist
	// Por enquanto, apenas retorna sucesso
	
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	// Log de auditoria para logout
	if userAdminID > 0 {
		h.logAuditoriaUseCase.CreateSystemLog(r.Context(), "Logout", fmt.Sprintf("Usuário ID %d realizou logout", userAdminID), clientIP)
	}

	response.WriteSuccess(w, http.StatusOK, "Logout realizado com sucesso", nil)
}

// ValidateToken godoc
// @Summary Validar token JWT
// @Description Valida se um token JWT é válido
// @Tags autenticacao
// @Accept json
// @Produce json
// @Param token body dto.ValidateTokenRequest true "Token para validação"
// @Success 200 {object} response.TokenValidationResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/validate [post]
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var req ValidateTokenRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	if strings.TrimSpace(req.Token) == "" {
		response.WriteError(w, http.StatusBadRequest, "Token obrigatório", "Token é obrigatório")
		return
	}

	// Validar token
	claims, err := h.validateJWT(req.Token)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Token inválido", err.Error())
		return
	}

	// Verificar se usuário ainda existe e está ativo
	usuario, err := h.usuarioUseCase.GetByID(r.Context(), claims.UserID)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Usuário não encontrado", "Usuário não existe ou foi removido")
		return
	}

	validationResponse := response.TokenValidationResponse{
		Valid: usuario.Status == "Ativo",
		User: response.UserInfo{
			ID:        usuario.ID,
			Nome:      usuario.NomeAdmin,
			Email:     usuario.Email,
			EmpresaID: usuario.IDEmpresa,
			Status:    usuario.Status,
		},
		ExpiresAt: claims.ExpiresAt.Time,
	}

	response.WriteSuccess(w, http.StatusOK, "Token validado", validationResponse)
}

// ChangePassword godoc
// @Summary Alterar senha
// @Description Permite ao usuário alterar sua própria senha
// @Tags autenticacao
// @Accept json
// @Produce json
// @Param password body dto.ChangePasswordRequest true "Dados para alteração de senha"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validateChangePasswordRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	if userAdminID == 0 {
		response.WriteError(w, http.StatusUnauthorized, "Não autorizado", "Token inválido ou expirado")
		return
	}

	clientIP := h.getClientIP(r)

	// Verificar senha atual
	usuario, err := h.usuarioUseCase.GetByID(r.Context(), userAdminID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Validar senha atual
	if !h.usuarioUseCase.ValidatePassword(req.SenhaAtual, usuario.SenhaHash) {
		response.WriteError(w, http.StatusUnauthorized, "Senha atual incorreta", "A senha atual fornecida está incorreta")
		return
	}

	// Atualizar senha
	if err := h.usuarioUseCase.UpdatePassword(r.Context(), userAdminID, req.NovaSenha, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Senha alterada com sucesso", nil)
}

// ForgotPassword godoc
// @Summary Solicitar recuperação de senha
// @Description Inicia processo de recuperação de senha por email
// @Tags autenticacao
// @Accept json
// @Produce json
// @Param email body dto.ForgotPasswordRequest true "Email para recuperação"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		response.WriteError(w, http.StatusBadRequest, "Email obrigatório", "Email é obrigatório")
		return
	}

	clientIP := h.getClientIP(r)

	// Processar solicitação de recuperação
	if err := h.usuarioUseCase.RequestPasswordReset(r.Context(), req.Email, clientIP); err != nil {
		// Por segurança, sempre retorna sucesso mesmo se email não existir
		response.WriteSuccess(w, http.StatusOK, "Se o email existir, instruções foram enviadas", nil)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Se o email existir, instruções foram enviadas", nil)
}

// Métodos auxiliares

func (h *AuthHandler) generateJWT(userID, empresaID int, email string) (string, error) {
	claims := middleware.JWTClaims{
		UserID:    userID,
		EmpresaID: empresaID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "organizational-climate-survey",
			Subject:   fmt.Sprintf("user_%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

func (h *AuthHandler) validateJWT(tokenString string) (*middleware.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return h.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*middleware.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

func (h *AuthHandler) validateLoginRequest(req *LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email é obrigatório")
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("email inválido")
	}
	if strings.TrimSpace(req.Senha) == "" {
		return fmt.Errorf("senha é obrigatória")
	}
	return nil
}

func (h *AuthHandler) validateChangePasswordRequest(req *ChangePasswordRequest) error {
	if strings.TrimSpace(req.SenhaAtual) == "" {
		return fmt.Errorf("senha atual é obrigatória")
	}
	if strings.TrimSpace(req.NovaSenha) == "" {
		return fmt.Errorf("nova senha é obrigatória")
	}
	if len(req.NovaSenha) < 8 {
		return fmt.Errorf("nova senha deve ter pelo menos 8 caracteres")
	}
	if req.SenhaAtual == req.NovaSenha {
		return fmt.Errorf("nova senha deve ser diferente da atual")
	}
	return nil
}

func (h *AuthHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *AuthHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra as rotas do handler
func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	// Rotas públicas (sem autenticação)
	router.HandleFunc("/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/auth/forgot-password", h.ForgotPassword).Methods("POST")
	router.HandleFunc("/auth/validate", h.ValidateToken).Methods("POST")
	
	// Rotas que requerem autenticação
	router.HandleFunc("/auth/refresh", h.RefreshToken).Methods("POST")
	router.HandleFunc("/auth/logout", h.Logout).Methods("POST")
	router.HandleFunc("/auth/change-password", h.ChangePassword).Methods("POST")
}