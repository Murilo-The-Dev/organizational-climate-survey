package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/application/dto/response"

	"github.com/gorilla/mux"
)

type EmpresaHandler struct {
	empresaUseCase *usecase.EmpresaUseCase
}

func NewEmpresaHandler(empresaUseCase *usecase.EmpresaUseCase) *EmpresaHandler {
	return &EmpresaHandler{
		empresaUseCase: empresaUseCase,
	}
}

// CreateEmpresa godoc
// @Summary Criar nova empresa
// @Description Cria uma nova empresa no sistema
// @Tags empresas
// @Accept json
// @Produce json
// @Param empresa body dto.EmpresaCreateRequest true "Dados da empresa"
// @Success 201 {object} response.EmpresaResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas [post]
func (h *EmpresaHandler) CreateEmpresa(w http.ResponseWriter, r *http.Request) {
	var req dto.EmpresaCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validateEmpresaCreateRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	empresa := req.ToEntity()
	
	// Obter informações do usuário autenticado (middleware deve definir)
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.empresaUseCase.Create(r.Context(), empresa, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "já cadastrada") {
			response.WriteError(w, http.StatusConflict, "Empresa já existe", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Empresa criada com sucesso", response.ToEmpresaResponse(empresa))
}

// GetEmpresa godoc
// @Summary Buscar empresa por ID
// @Description Retorna uma empresa específica pelo ID
// @Tags empresas
// @Produce json
// @Param id path int true "ID da empresa"
// @Success 200 {object} response.EmpresaResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{id} [get]
func (h *EmpresaHandler) GetEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	empresa, err := h.empresaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Empresa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Empresa encontrada", response.ToEmpresaResponse(empresa))
}

// ListEmpresas godoc
// @Summary Listar empresas
// @Description Retorna lista paginada de empresas
// @Tags empresas
// @Produce json
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas [get]
func (h *EmpresaHandler) ListEmpresas(w http.ResponseWriter, r *http.Request) {
    limit, offset := h.getPaginationParams(r)

    empresas, err := h.empresaUseCase.List(r.Context(), limit, offset)
    if err != nil {
        response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
        return
    }

    // Converter para response
    empresasResponse := make([]interface{}, len(empresas))
    for i, empresa := range empresas {
        empresasResponse[i] = response.ToEmpresaResponse(empresa)
    }

    // Criar informações de paginação
    pagination := response.PaginationInfo{
        Page:       (offset / limit) + 1,
        Limit:      limit,
        Total:      len(empresasResponse), // Em produção, fazer COUNT separado
        TotalPages: (len(empresasResponse) + limit - 1) / limit,
    }

    response.WritePaginated(w, http.StatusOK, "Empresas listadas com sucesso", empresasResponse, pagination)
}

// UpdateEmpresa godoc
// @Summary Atualizar empresa
// @Description Atualiza dados de uma empresa existente
// @Tags empresas
// @Accept json
// @Produce json
// @Param id path int true "ID da empresa"
// @Param empresa body dto.EmpresaUpdateRequest true "Dados para atualização"
// @Success 200 {object} response.EmpresaResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{id} [put]
func (h *EmpresaHandler) UpdateEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req dto.EmpresaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Buscar empresa existente
	empresa, err := h.empresaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Empresa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Aplicar atualizações
	req.ApplyToEntity(empresa)

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.empresaUseCase.Update(r.Context(), empresa, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "já está sendo usado") {
			response.WriteError(w, http.StatusConflict, "CNPJ já em uso", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Empresa atualizada com sucesso", response.ToEmpresaResponse(empresa))
}

// DeleteEmpresa godoc
// @Summary Deletar empresa
// @Description Remove uma empresa do sistema
// @Tags empresas
// @Produce json
// @Param id path int true "ID da empresa"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{id} [delete]
func (h *EmpresaHandler) DeleteEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.empresaUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Empresa não encontrada", err.Error())
			return
		}
		if strings.Contains(err.Error(), "possui") && strings.Contains(err.Error(), "vinculados") {
			response.WriteError(w, http.StatusConflict, "Empresa possui dependências", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Empresa deletada com sucesso", nil)
}

// GetEmpresaByCNPJ godoc
// @Summary Buscar empresa por CNPJ
// @Description Retorna uma empresa específica pelo CNPJ
// @Tags empresas
// @Produce json
// @Param cnpj path string true "CNPJ da empresa"
// @Success 200 {object} response.EmpresaResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/cnpj/{cnpj} [get]
func (h *EmpresaHandler) GetEmpresaByCNPJ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cnpj := vars["cnpj"]

	if strings.TrimSpace(cnpj) == "" {
		response.WriteError(w, http.StatusBadRequest, "CNPJ inválido", "CNPJ é obrigatório")
		return
	}

	empresa, err := h.empresaUseCase.GetByCNPJ(r.Context(), cnpj)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Empresa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Empresa encontrada", response.ToEmpresaResponse(empresa))
}

// Métodos auxiliares

func (h *EmpresaHandler) validateEmpresaCreateRequest(req *dto.EmpresaCreateRequest) error {
	if strings.TrimSpace(req.NomeFantasia) == "" {
		return fmt.Errorf("nome fantasia é obrigatório")
	}
	if strings.TrimSpace(req.RazaoSocial) == "" {
		return fmt.Errorf("razão social é obrigatória")
	}
	if strings.TrimSpace(req.CNPJ) == "" {
		return fmt.Errorf("CNPJ é obrigatório")
	}
	return nil
}

func (h *EmpresaHandler) getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 20 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset = 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}

func (h *EmpresaHandler) getUserAdminIDFromContext(r *http.Request) int {
	// Este método deve ser implementado com base no middleware de autenticação
	// Por enquanto, retorna 0 (sistema)
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *EmpresaHandler) getClientIP(r *http.Request) string {
	// Verifica headers de proxy
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra as rotas do handler
func (h *EmpresaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/empresas", h.CreateEmpresa).Methods("POST")
	router.HandleFunc("/empresas", h.ListEmpresas).Methods("GET")
	router.HandleFunc("/empresas/{id:[0-9]+}", h.GetEmpresa).Methods("GET")
	router.HandleFunc("/empresas/{id:[0-9]+}", h.UpdateEmpresa).Methods("PUT")
	router.HandleFunc("/empresas/{id:[0-9]+}", h.DeleteEmpresa).Methods("DELETE")
	router.HandleFunc("/empresas/cnpj/{cnpj}", h.GetEmpresaByCNPJ).Methods("GET")
}