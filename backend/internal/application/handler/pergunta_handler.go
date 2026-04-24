// Package handler implementa os controladores HTTP da aplicação.
// Processa requisições, valida entrada e coordena a execução de casos de uso.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/pkg/logger"

	"github.com/gorilla/mux"
)

// PerguntaHandler gerencia requisições HTTP relacionadas a perguntas de pesquisas
type PerguntaHandler struct {
	perguntaUseCase *usecase.PerguntaUseCase
	log             logger.Logger
}

// PerguntaOrdemUpdateRequest define payload para atualização da ordem de uma pergunta.
type PerguntaOrdemUpdateRequest struct {
	NovaOrdem int `json:"nova_ordem" example:"2"`
}

// PerguntaReorderRequest define payload para reordenação em lote das perguntas de uma pesquisa.
type PerguntaReorderRequest struct {
	PerguntaIDs []int `json:"pergunta_ids" example:"12,13,14"`
}

// NewPerguntaHandler cria nova instância do handler de perguntas
func NewPerguntaHandler(perguntaUseCase *usecase.PerguntaUseCase, log logger.Logger) *PerguntaHandler {
	return &PerguntaHandler{
		perguntaUseCase: perguntaUseCase,
		log:             log,
	}
}

// CreatePergunta cria nova pergunta no sistema
// @Summary Criar pergunta
// @Tags perguntas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.PerguntaCreateRequest true "Dados da pergunta"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas [post]
func (h *PerguntaHandler) CreatePergunta(w http.ResponseWriter, r *http.Request) {
	var req dto.PerguntaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}
	if err := h.validatePerguntaCreateRequest(&req); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}
	pergunta := req.ToEntity()
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)
	if err := h.perguntaUseCase.Create(r.Context(), pergunta, userAdminID, clientIP); err != nil {
		h.log.WithFields(map[string]interface{}{"user_admin_id": userAdminID, "client_ip": clientIP}).Error("Erro ao criar pergunta: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	perguntaResponse := &response.PerguntaResponse{
		ID:             pergunta.ID,
		TextoPergunta:  pergunta.TextoPergunta,
		TipoPergunta:   pergunta.TipoPergunta,
		OrdemExibicao:  pergunta.OrdemExibicao,
		OpcoesResposta: pergunta.OpcoesResposta,
	}

	h.log.WithFields(map[string]interface{}{"pergunta_id": pergunta.ID, "user_admin_id": userAdminID}).Info("Pergunta criada com sucesso")
	response.WriteSuccess(w, http.StatusCreated, "Pergunta criada com sucesso", perguntaResponse)
}

// CreatePerguntasBatch cria múltiplas perguntas em uma única operação
// @Summary Criar perguntas em lote
// @Tags perguntas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body []dto.PerguntaCreateRequest true "Lista de perguntas"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/batch [post]
func (h *PerguntaHandler) CreatePerguntasBatch(w http.ResponseWriter, r *http.Request) {
	var reqs []dto.PerguntaCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	if len(reqs) == 0 {
		response.WriteError(w, http.StatusBadRequest, "Lista vazia", "Pelo menos uma pergunta deve ser fornecida")
		return
	}

	if len(reqs) > 50 {
		response.WriteError(w, http.StatusBadRequest, "Muitas perguntas", "Máximo de 50 perguntas por operação")
		return
	}

	// Validar e converter todas as perguntas
	perguntas := make([]*entity.Pergunta, len(reqs))
	for i, req := range reqs {
		if err := h.validatePerguntaCreateRequest(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Validação falhou na pergunta %d", i+1), err.Error())
			return
		}
		perguntas[i] = req.ToEntity()
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.CreateBatch(r.Context(), perguntas, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	perguntasResponse := make([]response.PerguntaResponse, len(perguntas))
	for i, pergunta := range perguntas {
		perguntasResponse[i] = response.PerguntaResponse{
			ID:             pergunta.ID,
			TextoPergunta:  pergunta.TextoPergunta,
			TipoPergunta:   pergunta.TipoPergunta,
			OrdemExibicao:  pergunta.OrdemExibicao,
			OpcoesResposta: pergunta.OpcoesResposta,
		}
	}

	response.WriteSuccess(w, http.StatusCreated, "Perguntas criadas com sucesso", perguntasResponse)
}

// GetPergunta busca pergunta por ID
// @Summary Buscar pergunta por ID
// @Tags perguntas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{id} [get]
func (h *PerguntaHandler) GetPergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	pergunta, err := h.perguntaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	perguntaResponse := &response.PerguntaResponse{
		ID:             pergunta.ID,
		TextoPergunta:  pergunta.TextoPergunta,
		TipoPergunta:   pergunta.TipoPergunta,
		OrdemExibicao:  pergunta.OrdemExibicao,
		OpcoesResposta: pergunta.OpcoesResposta,
	}

	response.WriteSuccess(w, http.StatusOK, "Pergunta encontrada", perguntaResponse)
}

// ListPerguntasByPesquisa lista perguntas de pesquisa ordenadas por exibição
// @Summary Listar perguntas por pesquisa
// @Tags perguntas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/perguntas [get]
func (h *PerguntaHandler) ListPerguntasByPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	perguntas, err := h.perguntaUseCase.ListByPesquisa(r.Context(), pesquisaID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	perguntasResponse := make([]response.PerguntaResponse, len(perguntas))
	for i, pergunta := range perguntas {
		perguntasResponse[i] = response.PerguntaResponse{
			ID:             pergunta.ID,
			TextoPergunta:  pergunta.TextoPergunta,
			TipoPergunta:   pergunta.TipoPergunta,
			OrdemExibicao:  pergunta.OrdemExibicao,
			OpcoesResposta: pergunta.OpcoesResposta,
		}
	}

	response.WriteSuccess(w, http.StatusOK, "Perguntas listadas com sucesso", perguntasResponse)
}

// UpdatePergunta atualiza dados de pergunta existente
// @Summary Atualizar pergunta
// @Tags perguntas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da pergunta"
// @Param body body dto.PerguntaUpdateRequest true "Dados para atualização"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{id} [put]
func (h *PerguntaHandler) UpdatePergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req dto.PerguntaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Buscar pergunta existente
	pergunta, err := h.perguntaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Aplicar alterações parciais à entidade
	req.ApplyToEntity(pergunta)

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.Update(r.Context(), pergunta, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	perguntaResponse := &response.PerguntaResponse{
		ID:             pergunta.ID,
		TextoPergunta:  pergunta.TextoPergunta,
		TipoPergunta:   pergunta.TipoPergunta,
		OrdemExibicao:  pergunta.OrdemExibicao,
		OpcoesResposta: pergunta.OpcoesResposta,
	}

	response.WriteSuccess(w, http.StatusOK, "Pergunta atualizada com sucesso", perguntaResponse)
}

// UpdateOrdemPergunta atualiza ordem de exibição de pergunta específica
// @Summary Atualizar ordem de exibição da pergunta
// @Tags perguntas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da pergunta"
// @Param body body PerguntaOrdemUpdateRequest true "Nova ordem"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{id}/ordem [put]
func (h *PerguntaHandler) UpdateOrdemPergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req PerguntaOrdemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar ordem fornecida
	if req.NovaOrdem <= 0 {
		response.WriteError(w, http.StatusBadRequest, "Ordem inválida", "Ordem deve ser maior que zero")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.UpdateOrdem(r.Context(), id, req.NovaOrdem, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Ordem da pergunta atualizada com sucesso", nil)
}

// DeletePergunta remove pergunta do sistema
// @Summary Remover pergunta
// @Tags perguntas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{id} [delete]
func (h *PerguntaHandler) DeletePergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		if strings.Contains(err.Error(), "possui") && strings.Contains(err.Error(), "respostas") {
			response.WriteError(w, http.StatusConflict, "Pergunta possui respostas", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Pergunta deletada com sucesso", nil)
}

// ReorderPerguntas reordena todas as perguntas de uma pesquisa
// @Summary Reordenar perguntas da pesquisa
// @Tags perguntas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Param body body PerguntaReorderRequest true "IDs de perguntas na ordem desejada"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/perguntas/reorder [put]
func (h *PerguntaHandler) ReorderPerguntas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	var req PerguntaReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar lista de IDs
	if len(req.PerguntaIDs) == 0 {
		response.WriteError(w, http.StatusBadRequest, "Lista vazia", "Lista de IDs das perguntas é obrigatória")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.ReorderPerguntas(r.Context(), pesquisaID, req.PerguntaIDs, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Perguntas reordenadas com sucesso", nil)
}

// GetPerguntasWithStats lista perguntas de pesquisa com estatísticas de respostas
// @Summary Listar perguntas com estatísticas
// @Tags perguntas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/perguntas/with-stats [get]
func (h *PerguntaHandler) GetPerguntasWithStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	perguntasWithStats, err := h.perguntaUseCase.GetPerguntasWithStats(r.Context(), pesquisaID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Perguntas com estatísticas listadas com sucesso", perguntasWithStats)
}

// validatePerguntaCreateRequest valida campos obrigatórios e regras de negócio para criação
func (h *PerguntaHandler) validatePerguntaCreateRequest(req *dto.PerguntaCreateRequest) error {
	if req.IDPesquisa <= 0 {
		return fmt.Errorf("ID da pesquisa é obrigatório")
	}
	if strings.TrimSpace(req.TextoPergunta) == "" {
		return fmt.Errorf("texto da pergunta é obrigatório")
	}
	if len(req.TextoPergunta) < 5 {
		return fmt.Errorf("texto da pergunta deve ter pelo menos 5 caracteres")
	}
	if len(req.TextoPergunta) > 500 {
		return fmt.Errorf("texto da pergunta não pode exceder 500 caracteres")
	}
	if !h.isValidTipoPergunta(req.TipoPergunta) {
		return fmt.Errorf("tipo de pergunta inválido")
	}
	if req.OrdemExibicao <= 0 {
		return fmt.Errorf("ordem de exibição deve ser maior que zero")
	}
	return nil
}

// isValidTipoPergunta verifica se tipo de pergunta fornecido é válido
func (h *PerguntaHandler) isValidTipoPergunta(tipo string) bool {
	validTypes := []string{"MultiplaEscolha", "RespostaAberta", "EscalaNumerica", "SimNao"}
	for _, validType := range validTypes {
		if tipo == validType {
			return true
		}
	}
	return false
}

// getUserAdminIDFromContext extrai ID do usuário administrativo do contexto da requisição
func (h *PerguntaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP extrai endereço IP do cliente considerando proxies
func (h *PerguntaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *PerguntaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/perguntas", h.CreatePergunta).Methods("POST")
	router.HandleFunc("/perguntas/batch", h.CreatePerguntasBatch).Methods("POST")
	router.HandleFunc("/perguntas/{id:[0-9]+}", h.GetPergunta).Methods("GET")
	router.HandleFunc("/perguntas/{id:[0-9]+}", h.UpdatePergunta).Methods("PUT")
	router.HandleFunc("/perguntas/{id:[0-9]+}", h.DeletePergunta).Methods("DELETE")
	router.HandleFunc("/perguntas/{id:[0-9]+}/ordem", h.UpdateOrdemPergunta).Methods("PUT")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/perguntas", h.ListPerguntasByPesquisa).Methods("GET")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/perguntas/reorder", h.ReorderPerguntas).Methods("PUT")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/perguntas/with-stats", h.GetPerguntasWithStats).Methods("GET")
}
