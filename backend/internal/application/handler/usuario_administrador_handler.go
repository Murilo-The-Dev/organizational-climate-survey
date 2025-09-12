package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/entity"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UsuarioAdministradorHandler struct {
	usuarioUseCase *usecase.UsuarioAdministradorUseCase
}

func NewUsuarioAdministradorHandler(usuarioUseCase *usecase.UsuarioAdministradorUseCase) *UsuarioAdministradorHandler {
	return &UsuarioAdministradorHandler{
		usuarioUseCase: usuarioUseCase,
	}
}

// DTOs locais para operações específicas
type PasswordUpdateRequest struct {
	NovaSenha string `json:"nova_senha" binding:"required,min=8,max=128"`
}

type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=Ativo Inativo Pendente"`
}

// CreateUsuarioAdministrador godoc
// @Summary Criar novo usuário administrador
// @Description Cria um novo usuário administrador no sistema
// @Tags usuarios-administradores
// @Accept json
// @Produce json
// @Param usuario body dto.UsuarioAdministradorCreateRequest true "Dados do usuário"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores [post]
func (h *UsuarioAdministradorHandler) CreateUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	var req dto.UsuarioAdministradorCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validateUsuarioCreateRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	// Gerar hash da senha
	senhaHash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro ao processar senha", err.Error())
		return
	}

	usuario := req.ToEntity(string(senhaHash))
	
	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.usuarioUseCase.Create(r.Context(), usuario, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "já cadastrado") || strings.Contains(err.Error(), "já está sendo usado") {
			response.WriteError(w, http.StatusConflict, "Email já em uso", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	usuarioResponse := h.toUsuarioResponse(usuario)
	response.WriteSuccess(w, http.StatusCreated, "Usuário criado com sucesso", usuarioResponse)
}

// GetUsuarioAdministrador godoc
// @Summary Buscar usuário administrador por ID
// @Description Retorna um usuário administrador específico pelo ID
// @Tags usuarios-administradores
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores/{id} [get]
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

// ListUsuariosByEmpresa godoc
// @Summary Listar usuários por empresa
// @Description Retorna lista de usuários administradores de uma empresa
// @Tags usuarios-administradores
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param status query string false "Filtrar por status"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /empresas/{empresa_id}/usuarios-administradores [get]
func (h *UsuarioAdministradorHandler) ListUsuariosByEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	status := r.URL.Query().Get("status")
	
	var usuarios []*entity.UsuarioAdministrador
	
	if status != "" {
		usuarios, err = h.usuarioUseCase.ListByStatus(r.Context(), empresaID, status)
	} else {
		usuarios, err = h.usuarioUseCase.ListByEmpresa(r.Context(), empresaID)
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter para response
	usuariosResponse := make([]interface{}, len(usuarios))
	for i, usuario := range usuarios {
		usuariosResponse[i] = h.toUsuarioResponse(usuario)
	}

	response.WriteSuccess(w, http.StatusOK, "Usuários listados com sucesso", usuariosResponse)
}

// UpdateUsuarioAdministrador godoc
// @Summary Atualizar usuário administrador
// @Description Atualiza dados de um usuário administrador existente
// @Tags usuarios-administradores
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param usuario body dto.UsuarioAdministradorUpdateRequest true "Dados para atualização"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores/{id} [put]
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

	// Aplicar atualizações
	req.ApplyToEntity(usuario)

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

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

// UpdatePassword godoc
// @Summary Atualizar senha do usuário
// @Description Atualiza a senha de um usuário administrador
// @Tags usuarios-administradores
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param password body PasswordUpdateRequest true "Nova senha"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores/{id}/password [put]
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

	// Validação da senha
	if err := h.validatePassword(req.NovaSenha); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Senha inválida", err.Error())
		return
	}

	// Obter informações do usuário autenticado
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

// UpdateStatus godoc
// @Summary Atualizar status do usuário
// @Description Atualiza o status de um usuário administrador
// @Tags usuarios-administradores
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param status body StatusUpdateRequest true "Novo status"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores/{id}/status [put]
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

	// Validação do status
	if !h.isValidStatus(req.Status) {
		response.WriteError(w, http.StatusBadRequest, "Status inválido", "Status deve ser: Ativo, Inativo ou Pendente")
		return
	}

	// Obter informações do usuário autenticado
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

// DeleteUsuarioAdministrador godoc
// @Summary Deletar usuário administrador
// @Description Remove um usuário administrador do sistema
// @Tags usuarios-administradores
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores/{id} [delete]
func (h *UsuarioAdministradorHandler) DeleteUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.usuarioUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Usuário não encontrado", err.Error())
			return
		}
		if strings.Contains(err.Error(), "possui") && strings.Contains(err.Error(), "vinculados") {
			response.WriteError(w, http.StatusConflict, "Usuário possui dependências", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Usuário deletado com sucesso", nil)
}

// GetUsuarioByEmail godoc
// @Summary Buscar usuário por email
// @Description Retorna um usuário administrador específico pelo email
// @Tags usuarios-administradores
// @Produce json
// @Param email path string true "Email do usuário"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /usuarios-administradores/email/{email} [get]
func (h *UsuarioAdministradorHandler) GetUsuarioByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	if strings.TrimSpace(email) == "" {
		response.WriteError(w, http.StatusBadRequest, "Email inválido", "Email é obrigatório")
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

// Métodos auxiliares

func (h *UsuarioAdministradorHandler) validateUsuarioCreateRequest(req *dto.UsuarioAdministradorCreateRequest) error {
	if req.IDEmpresa <= 0 {
		return fmt.Errorf("ID da empresa é obrigatório")
	}
	if strings.TrimSpace(req.NomeAdmin) == "" {
		return fmt.Errorf("nome do administrador é obrigatório")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email é obrigatório")
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("email inválido")
	}
	if err := h.validatePassword(req.Senha); err != nil {
		return err
	}
	if !h.isValidStatus(req.Status) {
		return fmt.Errorf("status inválido")
	}
	return nil
}

func (h *UsuarioAdministradorHandler) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("senha deve ter pelo menos 8 caracteres")
	}
	if len(password) > 100 {
		return fmt.Errorf("senha não pode exceder 100 caracteres")
	}
	return nil
}

func (h *UsuarioAdministradorHandler) isValidStatus(status string) bool {
	validStatuses := []string{"Ativo", "Inativo", "Pendente"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func (h *UsuarioAdministradorHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *UsuarioAdministradorHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// toUsuarioResponse converte entity.UsuarioAdministrador para response.UsuarioAdministradorResponse
func (h *UsuarioAdministradorHandler) toUsuarioResponse(usuario *entity.UsuarioAdministrador) response.UsuarioAdministradorResponse {
	resp := response.UsuarioAdministradorResponse{
		ID:           usuario.ID,
		NomeAdmin:    usuario.NomeAdmin,
		Email:        usuario.Email,
		DataCadastro: usuario.DataCadastro,
		Status:       usuario.Status,
	}

	// Converte empresa se carregada
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

// RegisterRoutes registra as rotas do handler
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