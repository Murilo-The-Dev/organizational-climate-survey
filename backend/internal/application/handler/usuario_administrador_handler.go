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
	"golang.org/x/crypto/bcrypt"
)

type UsuarioAdministradorHandler struct {
	usuarioUseCase *usecase.UsuarioAdministradorUseCase
	log            logger.Logger
	validator      *validator.Validator
}

func NewUsuarioAdministradorHandler(usuarioUseCase *usecase.UsuarioAdministradorUseCase, log logger.Logger, val *validator.Validator) *UsuarioAdministradorHandler {
	return &UsuarioAdministradorHandler{
		usuarioUseCase: usuarioUseCase,
		log:            log,
		validator:      val,
	}
}

type PasswordUpdateRequest struct {
	NovaSenha string `json:"nova_senha"`
}

type StatusUpdateRequest struct {
	Status string `json:"status"`
}

func (h *UsuarioAdministradorHandler) CreateUsuarioAdministrador(w http.ResponseWriter, r *http.Request) {
	var req dto.UsuarioAdministradorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	if err := h.validator.IsEmail(req.Email); err != nil {
		h.log.WithContext(r.Context()).Info("Email inválido: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Email inválido", err.Error())
		return
	}

	if err := h.validator.IsPasswordStrong(req.Senha); err != nil {
		h.log.WithContext(r.Context()).Info("Senha fraca: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Senha não atende requisitos", err.Error())
		return
	}

	if err := h.validateUsuarioCreateRequest(&req); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	senhaHash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		h.log.WithContext(r.Context()).Error("Erro ao gerar hash: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Erro ao processar senha", err.Error())
		return
	}

	usuario := req.ToEntity(string(senhaHash))
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

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

	usuariosResponse := make([]interface{}, len(usuarios))
	for i, usuario := range usuarios {
		usuariosResponse[i] = h.toUsuarioResponse(usuario)
	}

	response.WriteSuccess(w, http.StatusOK, "Usuários listados com sucesso", usuariosResponse)
}

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

	if req.Email != nil && *req.Email != "" {
		if err := h.validator.IsEmail(*req.Email); err != nil {
			h.log.WithContext(r.Context()).Info("Email inválido: %v", err)
			response.WriteError(w, http.StatusBadRequest, "Email inválido", err.Error())
			return
		}
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

	req.ApplyToEntity(usuario)

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
		if strings.Contains(err.Error(), "possui") && strings.Contains(err.Error(), "vinculados") {
			response.WriteError(w, http.StatusConflict, "Usuário possui dependências", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Usuário deletado com sucesso", nil)
}

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

func (h *UsuarioAdministradorHandler) toUsuarioResponse(usuario *entity.UsuarioAdministrador) response.UsuarioAdministradorResponse {
	resp := response.UsuarioAdministradorResponse{
		ID:           usuario.ID,
		NomeAdmin:    usuario.NomeAdmin,
		Email:        usuario.Email,
		DataCadastro: usuario.DataCadastro,
		Status:       usuario.Status,
	}

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