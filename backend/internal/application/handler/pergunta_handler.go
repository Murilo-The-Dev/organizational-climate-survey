package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/application/dto/response"

	"github.com/gorilla/mux"
)

type PerguntaHandler struct {
	perguntaUseCase *usecase.PerguntaUseCase
}

func NewPerguntaHandler(perguntaUseCase *usecase.PerguntaUseCase) *PerguntaHandler {
	return &PerguntaHandler{
		perguntaUseCase: perguntaUseCase,
	}
}

// CreatePergunta godoc
// @Summary Criar nova pergunta
// @Description Cria uma nova pergunta para uma pesquisa
// @Tags perguntas
// @Accept json
// @Produce json
// @Param pergunta body dto.PerguntaCreateRequest true "Dados da pergunta"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /perguntas [post]
func (h *PerguntaHandler) CreatePergunta(w http.ResponseWriter, r *http.Request) {
	var req dto.PerguntaCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validatePerguntaCreateRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	pergunta := req.ToEntity()
	
	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.Create(r.Context(), pergunta, userAdminID, clientIP); err != nil {
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

	response.WriteSuccess(w, http.StatusCreated, "Pergunta criada com sucesso", perguntaResponse)
}

// CreatePerguntasBatch godoc
// @Summary Criar múltiplas perguntas
// @Description Cria múltiplas perguntas para uma pesquisa em uma única operação
// @Tags perguntas
// @Accept json
// @Produce json
// @Param perguntas body []dto.PerguntaCreateRequest true "Lista de perguntas"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /perguntas/batch [post]
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

	// Validar todas as perguntas
	perguntas := make([]*entity.Pergunta, len(reqs))
	for i, req := range reqs {
		if err := h.validatePerguntaCreateRequest(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Validação falhou na pergunta %d", i+1), err.Error())
			return
		}
		perguntas[i] = req.ToEntity()
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.CreateBatch(r.Context(), perguntas, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter para response
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

// GetPergunta godoc
// @Summary Buscar pergunta por ID
// @Description Retorna uma pergunta específica pelo ID
// @Tags perguntas
// @Produce json
// @Param id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /perguntas/{id} [get]
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

// ListPerguntasByPesquisa godoc
// @Summary Listar perguntas por pesquisa
// @Description Retorna lista ordenada de perguntas de uma pesquisa específica
// @Tags perguntas
// @Produce json
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{pesquisa_id}/perguntas [get]
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

	// Converter para response
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

// UpdatePergunta godoc
// @Summary Atualizar pergunta
// @Description Atualiza dados de uma pergunta existente
// @Tags perguntas
// @Accept json
// @Produce json
// @Param id path int true "ID da pergunta"
// @Param pergunta body dto.PerguntaUpdateRequest true "Dados para atualização"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /perguntas/{id} [put]
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

	// Aplicar atualizações
	req.ApplyToEntity(pergunta)

	// Obter informações do usuário autenticado
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

// UpdateOrdemPergunta godoc
// @Summary Atualizar ordem da pergunta
// @Description Atualiza a ordem de exibição de uma pergunta
// @Tags perguntas
// @Accept json
// @Produce json
// @Param id path int true "ID da pergunta"
// @Param ordem body struct{NovaOrdem int `json:"nova_ordem"`} true "Nova ordem"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /perguntas/{id}/ordem [put]
func (h *PerguntaHandler) UpdateOrdemPergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req struct {
		NovaOrdem int `json:"nova_ordem"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação da ordem
	if req.NovaOrdem <= 0 {
		response.WriteError(w, http.StatusBadRequest, "Ordem inválida", "Ordem deve ser maior que zero")
		return
	}

	// Obter informações do usuário autenticado
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

// DeletePergunta godoc
// @Summary Deletar pergunta
// @Description Remove uma pergunta do sistema
// @Tags perguntas
// @Produce json
// @Param id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /perguntas/{id} [delete]
func (h *PerguntaHandler) DeletePergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Obter informações do usuário autenticado
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

// ReorderPerguntas godoc
// @Summary Reordenar perguntas
// @Description Reordena todas as perguntas de uma pesquisa
// @Tags perguntas
// @Accept json
// @Produce json
// @Param pesquisa_id path int true "ID da pesquisa"
// @Param ordem body struct{PerguntaIDs []int `json:"pergunta_ids"`} true "Nova ordem das perguntas"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{pesquisa_id}/perguntas/reorder [put]
func (h *PerguntaHandler) ReorderPerguntas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	var req struct {
		PerguntaIDs []int `json:"pergunta_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if len(req.PerguntaIDs) == 0 {
		response.WriteError(w, http.StatusBadRequest, "Lista vazia", "Lista de IDs das perguntas é obrigatória")
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.perguntaUseCase.ReorderPerguntas(r.Context(), pesquisaID, req.PerguntaIDs, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Perguntas reordenadas com sucesso", nil)
}

// GetPerguntasWithStats godoc
// @Summary Listar perguntas com estatísticas
// @Description Retorna perguntas de uma pesquisa com estatísticas de respostas
// @Tags perguntas
// @Produce json
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{pesquisa_id}/perguntas/with-stats [get]
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

// Métodos auxiliares

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

func (h *PerguntaHandler) isValidTipoPergunta(tipo string) bool {
	validTypes := []string{"MultiplaEscolha", "RespostaAberta", "EscalaNumerica", "SimNao"}
	for _, validType := range validTypes {
		if tipo == validType {
			return true
		}
	}
	return false
}

func (h *PerguntaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *PerguntaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra as rotas do handler
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