// Package handler implementa os controladores HTTP da aplicação.
// Processa requisições, valida entrada e coordena a execução de casos de uso.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/pkg/logger"
	"organizational-climate-survey/backend/pkg/validator"

	"github.com/gorilla/mux"
)

// UsuarioAdministradorHandler gerencia requisições HTTP relacionadas a usuários administrativos
type UsuarioAdministradorHandler struct {
	usuarioUseCase *usecase.UsuarioAdministradorUseCase
	log            logger.Logger
	validator      *validator.Validator
}

// NewUsuarioAdministradorHandler cria nova instância do handler de usuários administrativos
func NewUsuarioAdministradorHandler(usuarioUseCase *usecase.UsuarioAdministradorUseCase, log logger.Logger, val *validator.Validator) *UsuarioAdministradorHandler {
	return &UsuarioAdministradorHandler{
		usuarioUseCase: usuarioUseCase,
		log:            log,
		validator:      val,
	}
}

// PasswordUpdateRequest representa requisição de atualização de senha
type PasswordUpdateRequest struct {
	NovaSenha string `json:"nova_senha"` // Nova senha a ser definida
}

// StatusUpdateRequest representa requisição de atualização de status
type StatusUpdateRequest struct {
	Status string `json:"status"` // Novo status do usuário
}

// CreateUsuarioAdministrador cria novo usuário administrativo no sistema
func (h *UsuarioAdministradorHandler) CreateUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	var req dto.UsuarioAdministradorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar formato do email
	if err := h.validator.IsEmail(req.Email); err != nil {
		h.log.WithContext(r.Context()).Info("Email inválido: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Email inválido", err.Error())
		return
	}

	// Validar força da senha
	if err := h.validator.IsPasswordStrong(req.Senha); err != nil {
		h.log.WithContext(r.Context()).Info("Senha fraca: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Senha não atende requisitos", err.Error())
		return
	}

	// Validar campos obrigatórios e regras de negócio
	if err := h.validateUsuarioCreateRequest(&req); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	// Converter DTO para entidade de domínio
	usuario := req.ToEntity(req.Senha)
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	usuario.SenhaHash = req.Senha

	// Executar caso de uso de criação
	if err := h.usuarioUseCase.Create(r.Context(), usuario, userAdminID, clientIP); err != nil {
		h.log.WithFields(map[string]interface{}{"user_admin_id": userAdminID}).Error("Erro ao criar usuario: %v", err)
		if strings.Contains(err.Error(), "já cadastrado") {
			response.WriteError(w, http.StatusConflict, "Email já em uso", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	h.log.WithFields(map[string]interface{}{"usuario_id": usuario.ID, "user_admin_id": userAdminID}).Info("Usuário administrador criado com sucesso")
	response.WriteSuccess(w, http.StatusCreated, "Usuário criado com sucesso", h.toUsuarioResponse(usuario))
}

// GetUsuarioAdministrador busca usuário administrativo por ID
func (h *UsuarioAdministradorHandler) GetUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	usuario, err := h.usuarioUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	usuarioResponse := h.toUsuarioResponse(usuario)
	response.WriteSuccess(w, http.StatusOK, "Usuário encontrado", usuarioResponse)
}

// ListUsuariosByEmpresa lista usuários administrativos de uma empresa com filtro opcional de status
func (h *UsuarioAdministradorHandler) ListUsuariosByEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	status := r.URL.Query().Get("status")
	
	var usuarios []*entity.UsuarioAdministrador
	
	// Aplicar filtro de status se fornecido
	if status != "" {
		usuarios, err = h.usuarioUseCase.ListByStatus(r.Context(), empresaID, status)
	} else {
		usuarios, err = h.usuarioUseCase.ListByEmpresa(r.Context(), empresaID)
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	usuariosResponse := make([]interface{}, len(usuarios))
	for i, usuario := range usuarios {
		usuariosResponse[i] = h.toUsuarioResponse(usuario)
	}

	response.WriteSuccess(w, http.StatusOK, "Usuários listados com sucesso", usuariosResponse)
}

// UpdateUsuarioAdministrador atualiza dados de usuário administrativo existente
func (h *UsuarioAdministradorHandler) UpdateUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req dto.UsuarioAdministradorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar email se fornecido na atualização
	if req.Email != nil && *req.Email != "" {
		if err := h.validator.IsEmail(*req.Email); err != nil {
			h.log.WithContext(r.Context()).Info("Email inválido: %v", err)
			response.WriteError(w, http.StatusBadRequest, "Email inválido", err.Error())
			return
		}
	}

	// Buscar usuário existente
	usuario, err := h.usuarioUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Aplicar alterações parciais à entidade
	req.ApplyToEntity(usuario)

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	// Executar atualização
	if err := h.usuarioUseCase.Update(r.Context(), usuario, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "já está sendo usado") {
			response.WriteError(w, http.StatusConflict, "Email já em uso", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	usuarioResponse := h.toUsuarioResponse(usuario)
	response.WriteSuccess(w, http.StatusOK, "Usuário atualizado com sucesso", usuarioResponse)
}

// UpdatePassword atualiza senha de usuário administrativo
func (h *UsuarioAdministradorHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req PasswordUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar força da nova senha
	if err := h.validator.IsPasswordStrong(req.NovaSenha); err != nil {
		h.log.WithContext(r.Context()).Info("Senha fraca: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Senha não atende requisitos", err.Error())
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.usuarioUseCase.UpdatePassword(r.Context(), id, req.NovaSenha, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Senha atualizada com sucesso", nil)
}

// UpdateStatus atualiza status de usuário administrativo
func (h *UsuarioAdministradorHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar status contra valores permitidos
	validStatuses := []string{"Ativo", "Inativo", "Pendente"}
	if err := h.validator.IsValidStatus(req.Status, validStatuses); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Status inválido", err.Error())
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.usuarioUseCase.UpdateStatus(r.Context(), id, req.Status, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Status atualizado com sucesso", nil)
}

// DeleteUsuarioAdministrador inativa usuário administrativo (soft delete)
func (h *UsuarioAdministradorHandler) DeleteUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.usuarioUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		if strings.Contains(err.Error(), "já está inativo") {
			response.WriteError(w, http.StatusConflict, "Usuário já inativo", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Usuário inativado com sucesso", nil)
}

// GetUsuarioByEmail busca usuário administrativo por endereço de email
func (h *UsuarioAdministradorHandler) GetUsuarioByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	if err := h.validator.IsEmail(email); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Email inválido", err.Error())
		return
	}

	usuario, err := h.usuarioUseCase.GetByEmail(r.Context(), email)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	usuarioResponse := h.toUsuarioResponse(usuario)
	response.WriteSuccess(w, http.StatusOK, "Usuário encontrado", usuarioResponse)
}

// validateUsuarioCreateRequest valida campos obrigatórios e regras de negócio para criação
func (h *UsuarioAdministradorHandler) validateUsuarioCreateRequest(req *dto.UsuarioAdministradorCreateRequest) error {
	if req.IDEmpresa <= 0 {
		return validator.ValidationError{Field: "id_empresa", Message: "obrigatório"}
	}
	if strings.TrimSpace(req.NomeAdmin) == "" {
		return validator.ValidationError{Field: "nome_admin", Message: "obrigatório"}
	}
	validStatuses := []string{"Ativo", "Inativo", "Pendente"}
	return h.validator.IsValidStatus(req.Status, validStatuses)
}

// getUserAdminIDFromContext extrai ID do usuário administrativo do contexto da requisição
func (h *UsuarioAdministradorHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP extrai endereço IP do cliente considerando proxies
func (h *UsuarioAdministradorHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// toUsuarioResponse converte entidade de domínio para DTO de resposta
func (h *UsuarioAdministradorHandler) toUsuarioResponse(usuario *entity.UsuarioAdministrador) response.UsuarioAdministradorResponse {
	resp := response.UsuarioAdministradorResponse{
		ID:           usuario.ID,
		NomeAdmin:    usuario.NomeAdmin,
		Email:        usuario.Email,
		DataCadastro: usuario.DataCadastro,
		Status:       usuario.Status,
	}

	// Incluir dados da empresa se carregada
	if usuario.Empresa != nil {
		empresaResp := &response.EmpresaResponse{
			ID:           usuario.Empresa.ID,
			NomeFantasia: usuario.Empresa.NomeFantasia,
			RazaoSocial:  usuario.Empresa.RazaoSocial,
			CNPJ:         usuario.Empresa.CNPJ,
			DataCadastro: usuario.Empresa.DataCadastro,
		}
		resp.Empresa = empresaResp
	}

	return resp
}

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *UsuarioAdministradorHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/usuarios-administradores", h.CreateUsuarioAdministrador).Methods("POST")
	router.HandleFunc("/usuarios-administradores/{id:[0-9]+}", h.GetUsuarioAdministrador).Methods("GET")
	router.HandleFunc("/usuarios-administradores/{id:[0-9]+}", h.UpdateUsuarioAdministrador).Methods("PUT")
	router.HandleFunc("/usuarios-administradores/{id:[0-9]+}", h.DeleteUsuarioAdministrador).Methods("DELETE")
	router.HandleFunc("/usuarios-administradores/{id:[0-9]+}/password", h.UpdatePassword).Methods("PUT")
	router.HandleFunc("/usuarios-administradores/{id:[0-9]+}/status", h.UpdateStatus).Methods("PUT")
	router.HandleFunc("/usuarios-administradores/email/{email}", h.GetUsuarioByEmail).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/usuarios-administradores", h.ListUsuariosByEmpresa).Methods("GET")
}