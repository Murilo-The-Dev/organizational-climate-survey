// Package handler implementa os controladores HTTP da aplicação.
// Processa requisições, valida entrada e coordena a execução de casos de uso.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/pkg/logger"
	"organizational-climate-survey/backend/pkg/validator"

	"github.com/gorilla/mux"
)

// EmpresaHandler gerencia requisições HTTP relacionadas a empresas
type EmpresaHandler struct {
	empresaUseCase *usecase.EmpresaUseCase
	log            logger.Logger
	validator      *validator.Validator
}

// NewEmpresaHandler cria nova instância do handler de empresas
func NewEmpresaHandler(empresaUseCase *usecase.EmpresaUseCase, log logger.Logger, val *validator.Validator) *EmpresaHandler {
	return &EmpresaHandler{
		empresaUseCase: empresaUseCase,
		log:            log,
		validator:      val,
	}
}

// CreateEmpresa cria nova empresa no sistema
func (h *EmpresaHandler) CreateEmpresa(w http.ResponseWriter, r *http.Request) {
	var req dto.EmpresaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	if strings.TrimSpace(req.NomeFantasia) == "" {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", "nome fantasia é obrigatório")
		return
	}
	if strings.TrimSpace(req.RazaoSocial) == "" {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", "razão social é obrigatória")
		return
	}
	if err := h.validator.IsCNPJ(req.CNPJ); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	empresa := req.ToEntity()
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.empresaUseCase.Create(r.Context(), empresa, userAdminID, clientIP); err != nil {
		h.log.WithFields(map[string]interface{}{"user_admin_id": userAdminID, "client_ip": clientIP}).Error("Erro ao criar empresa: %v", err)
		if strings.Contains(err.Error(), "já cadastrada") {
			response.WriteError(w, http.StatusConflict, "Empresa já existe", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	h.log.WithFields(map[string]interface{}{"empresa_id": empresa.ID, "user_admin_id": userAdminID, "client_ip": clientIP}).Info("Empresa criada com sucesso")
	response.WriteSuccess(w, http.StatusCreated, "Empresa criada com sucesso", response.ToEmpresaResponse(empresa))
}

// GetEmpresa busca empresa por ID
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

// ListEmpresas lista todas as empresas com paginação
func (h *EmpresaHandler) ListEmpresas(w http.ResponseWriter, r *http.Request) {
    limit, offset := h.getPaginationParams(r)

    empresas, err := h.empresaUseCase.List(r.Context(), limit, offset)
    if err != nil {
        response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
        return
    }

    // Converter entidades para DTOs de resposta
    empresasResponse := make([]interface{}, len(empresas))
    for i, empresa := range empresas {
        empresasResponse[i] = response.ToEmpresaResponse(empresa)
    }

    // Calcular metadados de paginação
    pagination := response.PaginationInfo{
        Page:       (offset / limit) + 1,
        Limit:      limit,
        Total:      len(empresasResponse),
        TotalPages: (len(empresasResponse) + limit - 1) / limit,
    }

    response.WritePaginated(w, http.StatusOK, "Empresas listadas com sucesso", empresasResponse, pagination)
}

// UpdateEmpresa atualiza dados de empresa existente
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

	// Aplicar alterações parciais à entidade
	req.ApplyToEntity(empresa)

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

// DeleteEmpresa remove empresa do sistema
func (h *EmpresaHandler) DeleteEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

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

// GetEmpresaByCNPJ busca empresa por CNPJ
func (h *EmpresaHandler) GetEmpresaByCNPJ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cnpj := vars["cnpj"]
	
	// ADICIONAR: Decodificar URL
	cnpj, err := url.QueryUnescape(cnpj)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "CNPJ inválido", "Formato de CNPJ inválido")
		return
	}

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

// validateEmpresaCreateRequest valida campos obrigatórios e regras de negócio para criação
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

// getPaginationParams extrai parâmetros de paginação da query string com valores padrão
func (h *EmpresaHandler) getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset = 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}

// getUserAdminIDFromContext extrai ID do usuário administrativo do contexto da requisição
func (h *EmpresaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP extrai endereço IP do cliente considerando proxies
func (h *EmpresaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *EmpresaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/empresas", h.CreateEmpresa).Methods("POST")
	router.HandleFunc("/empresas", h.ListEmpresas).Methods("GET")
	router.HandleFunc("/empresas/{id:[0-9]+}", h.GetEmpresa).Methods("GET")
	router.HandleFunc("/empresas/{id:[0-9]+}", h.UpdateEmpresa).Methods("PUT")
	router.HandleFunc("/empresas/{id:[0-9]+}", h.DeleteEmpresa).Methods("DELETE")
	router.HandleFunc("/empresas/cnpj/{cnpj:.+}", h.GetEmpresaByCNPJ).Methods("GET")
}