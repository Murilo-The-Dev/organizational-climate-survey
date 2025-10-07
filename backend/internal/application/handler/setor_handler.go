// Package handler implementa os controladores HTTP da aplicação.
// Processa requisições, valida entrada e coordena a execução de casos de uso.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/pkg/logger"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// SetorHandler gerencia requisições HTTP relacionadas a setores organizacionais
type SetorHandler struct {
	setorUseCase *usecase.SetorUseCase
	log          logger.Logger
}

// NewSetorHandler cria nova instância do handler de setores
func NewSetorHandler(setorUseCase *usecase.SetorUseCase, log logger.Logger) *SetorHandler {
	return &SetorHandler{
		setorUseCase: setorUseCase,
		log:          log,
	}
}

// CreateSetor cria novo setor organizacional no sistema
func (h *SetorHandler) CreateSetor(w http.ResponseWriter, r *http.Request) {
	var req dto.SetorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}
	
	// Validar campos obrigatórios e regras de negócio
	if err := h.validateSetorCreateRequest(&req); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}
	
	// Converter DTO para entidade de domínio
	setor := req.ToEntity()
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)
	
	// Executar caso de uso de criação
	if err := h.setorUseCase.Create(r.Context(), setor, userAdminID, clientIP); err != nil {
		h.log.WithFields(map[string]interface{}{"user_admin_id": userAdminID, "client_ip": clientIP}).Error("Erro ao criar setor: %v", err)
		if strings.Contains(err.Error(), "já existe") {
			response.WriteError(w, http.StatusConflict, "Setor já existe", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}
	
	h.log.WithFields(map[string]interface{}{"setor_id": setor.ID, "user_admin_id": userAdminID}).Info("Setor criado com sucesso")
	response.WriteSuccess(w, http.StatusCreated, "Setor criado com sucesso", h.toSetorResponse(setor))
}

// GetSetor busca setor organizacional por ID
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

// ListSetoresByEmpresa lista todos os setores de uma empresa específica
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

	// Converter entidades para DTOs de resposta
	setoresResponse := make([]interface{}, len(setores))
	for i, setor := range setores {
		setoresResponse[i] = h.toSetorResponse(setor)
	}

	response.WriteSuccess(w, http.StatusOK, "Setores listados com sucesso", setoresResponse)
}

// UpdateSetor atualiza dados de setor organizacional existente
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

	// Aplicar alterações parciais à entidade
	req.ApplyToEntity(setor)

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	// Executar atualização
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

// DeleteSetor remove setor organizacional do sistema
func (h *SetorHandler) DeleteSetor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

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

// GetSetorByNome busca setor organizacional por nome dentro de uma empresa
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

// validateSetorCreateRequest valida campos obrigatórios e regras de negócio para criação
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

// getUserAdminIDFromContext extrai ID do usuário administrativo do contexto da requisição
func (h *SetorHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP extrai endereço IP do cliente considerando proxies
func (h *SetorHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// toSetorResponse converte entidade de domínio para DTO de resposta
func (h *SetorHandler) toSetorResponse(setor *entity.Setor) response.SetorResponse {
	resp := response.SetorResponse{
		ID:        setor.ID,
		NomeSetor: setor.NomeSetor,
		Descricao: setor.Descricao,
	}

	// Incluir dados da empresa se carregada
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

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *SetorHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/setores", h.CreateSetor).Methods("POST")
	router.HandleFunc("/setores/{id:[0-9]+}", h.GetSetor).Methods("GET")
	router.HandleFunc("/setores/{id:[0-9]+}", h.UpdateSetor).Methods("PUT")
	router.HandleFunc("/setores/{id:[0-9]+}", h.DeleteSetor).Methods("DELETE")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/setores", h.ListSetoresByEmpresa).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/setores/nome/{nome}", h.GetSetorByNome).Methods("GET")
}