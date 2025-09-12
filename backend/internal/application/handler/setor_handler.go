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
)

type SetorHandler struct {
	setorUseCase *usecase.SetorUseCase
}

func NewSetorHandler(setorUseCase *usecase.SetorUseCase) *SetorHandler {
	return &SetorHandler{
		setorUseCase: setorUseCase,
	}
}

// CreateSetor godoc
// @Summary Criar novo setor
// @Description Cria um novo setor para uma empresa
// @Tags setores
// @Accept json
// @Produce json
// @Param setor body dto.SetorCreateRequest true "Dados do setor"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /setores [post]
func (h *SetorHandler) CreateSetor(w http.ResponseWriter, r *http.Request) {
	var req dto.SetorCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validateSetorCreateRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	setor := req.ToEntity()
	
	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.setorUseCase.Create(r.Context(), setor, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "já existe") {
			response.WriteError(w, http.StatusConflict, "Setor já existe", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	setorResponse := h.toSetorResponse(setor)
	response.WriteSuccess(w, http.StatusCreated, "Setor criado com sucesso", setorResponse)
}

// GetSetor godoc
// @Summary Buscar setor por ID
// @Description Retorna um setor específico pelo ID
// @Tags setores
// @Produce json
// @Param id path int true "ID do setor"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /setores/{id} [get]
func (h *SetorHandler) GetSetor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	setor, err := h.setorUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Setor não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	setorResponse := h.toSetorResponse(setor)
	response.WriteSuccess(w, http.StatusOK, "Setor encontrado", setorResponse)
}

// ListSetoresByEmpresa godoc
// @Summary Listar setores por empresa
// @Description Retorna lista de setores de uma empresa específica
// @Tags setores
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /empresas/{empresa_id}/setores [get]
func (h *SetorHandler) ListSetoresByEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	setores, err := h.setorUseCase.ListByEmpresa(r.Context(), empresaID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter para response
	setoresResponse := make([]interface{}, len(setores))
	for i, setor := range setores {
		setoresResponse[i] = h.toSetorResponse(setor)
	}

	response.WriteSuccess(w, http.StatusOK, "Setores listados com sucesso", setoresResponse)
}

// UpdateSetor godoc
// @Summary Atualizar setor
// @Description Atualiza dados de um setor existente
// @Tags setores
// @Accept json
// @Produce json
// @Param id path int true "ID do setor"
// @Param setor body dto.SetorUpdateRequest true "Dados para atualização"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /setores/{id} [put]
func (h *SetorHandler) UpdateSetor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req dto.SetorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Buscar setor existente
	setor, err := h.setorUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Setor não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Aplicar atualizações
	req.ApplyToEntity(setor)

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.setorUseCase.Update(r.Context(), setor, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "já existe") {
			response.WriteError(w, http.StatusConflict, "Nome do setor já em uso", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	setorResponse := h.toSetorResponse(setor)
	response.WriteSuccess(w, http.StatusOK, "Setor atualizado com sucesso", setorResponse)
}

// DeleteSetor godoc
// @Summary Deletar setor
// @Description Remove um setor do sistema
// @Tags setores
// @Produce json
// @Param id path int true "ID do setor"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /setores/{id} [delete]
func (h *SetorHandler) DeleteSetor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.setorUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Setor não encontrado", err.Error())
			return
		}
		if strings.Contains(err.Error(), "possui") && strings.Contains(err.Error(), "vinculados") {
			response.WriteError(w, http.StatusConflict, "Setor possui dependências", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Setor deletado com sucesso", nil)
}

// GetSetorByNome godoc
// @Summary Buscar setor por nome
// @Description Retorna um setor específico pelo nome dentro de uma empresa
// @Tags setores
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param nome path string true "Nome do setor"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /empresas/{empresa_id}/setores/nome/{nome} [get]
func (h *SetorHandler) GetSetorByNome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	nome := vars["nome"]
	if strings.TrimSpace(nome) == "" {
		response.WriteError(w, http.StatusBadRequest, "Nome inválido", "Nome do setor é obrigatório")
		return
	}

	setor, err := h.setorUseCase.GetByNome(r.Context(), empresaID, nome)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Setor não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	setorResponse := h.toSetorResponse(setor)
	response.WriteSuccess(w, http.StatusOK, "Setor encontrado", setorResponse)
}

// Métodos auxiliares

func (h *SetorHandler) validateSetorCreateRequest(req *dto.SetorCreateRequest) error {
	if req.IDEmpresa <= 0 {
		return fmt.Errorf("ID da empresa é obrigatório")
	}
	if strings.TrimSpace(req.NomeSetor) == "" {
		return fmt.Errorf("nome do setor é obrigatório")
	}
	if len(req.NomeSetor) < 2 {
		return fmt.Errorf("nome do setor deve ter pelo menos 2 caracteres")
	}
	if len(req.NomeSetor) > 255 {
		return fmt.Errorf("nome do setor não pode exceder 255 caracteres")
	}
	return nil
}

func (h *SetorHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *SetorHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// toSetorResponse converte entity.Setor para response.SetorResponse
func (h *SetorHandler) toSetorResponse(setor *entity.Setor) response.SetorResponse {
	resp := response.SetorResponse{
		ID:        setor.ID,
		NomeSetor: setor.NomeSetor,
		Descricao: setor.Descricao,
	}

	// Converte empresa se carregada
	if setor.Empresa != nil {
		empresaResp := &response.EmpresaResponse{
			ID:           setor.Empresa.ID,
			NomeFantasia: setor.Empresa.NomeFantasia,
			RazaoSocial:  setor.Empresa.RazaoSocial,
			CNPJ:         setor.Empresa.CNPJ,
			DataCadastro: setor.Empresa.DataCadastro,
		}
		resp.Empresa = empresaResp
	}

	return resp
}

// RegisterRoutes registra as rotas do handler
func (h *SetorHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/setores", h.CreateSetor).Methods("POST")
	router.HandleFunc("/setores/{id:[0-9]+}", h.GetSetor).Methods("GET")
	router.HandleFunc("/setores/{id:[0-9]+}", h.UpdateSetor).Methods("PUT")
	router.HandleFunc("/setores/{id:[0-9]+}", h.DeleteSetor).Methods("DELETE")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/setores", h.ListSetoresByEmpresa).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/setores/nome/{nome}", h.GetSetorByNome).Methods("GET")
}