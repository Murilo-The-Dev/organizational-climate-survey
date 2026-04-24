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

// RespostaHandler gerencia requisições HTTP relacionadas a respostas de pesquisas
type RespostaHandler struct {
	respostaUseCase *usecase.RespostaUseCase
	log             logger.Logger
}

// NewRespostaHandler cria nova instância do handler de respostas
func NewRespostaHandler(respostaUseCase *usecase.RespostaUseCase, log logger.Logger) *RespostaHandler {
	return &RespostaHandler{
		respostaUseCase: respostaUseCase,
		log:             log,
	}
}

// SubmitRespostas processa submissão em lote de respostas de pesquisa
// @Summary Submeter respostas anônimas
// @Description Submete respostas de uma pesquisa usando token de acesso gerado previamente.
// @Tags respostas
// @Accept json
// @Produce json
// @Param body body dto.SubmitRespostasRequest true "Token e respostas"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/respostas/submit [post]
func (h *RespostaHandler) SubmitRespostas(w http.ResponseWriter, r *http.Request) {
	// MODIFICADO: Decodificar struct wrapper com token
	var req dto.SubmitRespostasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar token obrigatório
	if strings.TrimSpace(req.TokenAcesso) == "" {
		h.log.WithContext(r.Context()).Info("Token não fornecido")
		response.WriteError(w, http.StatusBadRequest, "Token obrigatório", "Token de acesso é obrigatório")
		return
	}

	// Validar lista de respostas
	if len(req.Respostas) == 0 {
		h.log.WithContext(r.Context()).Info("Nenhuma resposta enviada")
		response.WriteError(w, http.StatusBadRequest, "Lista vazia", "Pelo menos uma resposta deve ser fornecida")
		return
	}

	// Validar e converter todas as respostas
	respostas := make([]*entity.Resposta, len(req.Respostas))
	for i, respostaReq := range req.Respostas {
		if err := h.validateRespostaCreateRequest(&respostaReq); err != nil {
			h.log.WithContext(r.Context()).Info("Validação falhou na resposta %d: %v", i+1, err)
			response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Erro na resposta %d", i+1), err.Error())
			return
		}
		respostas[i] = respostaReq.ToEntity()
	}

	// MODIFICADO: Passar token para usecase
	if err := h.respostaUseCase.CreateBatch(r.Context(), respostas, req.TokenAcesso); err != nil {
		h.log.WithContext(r.Context()).Error("Erro ao salvar respostas: %v", err)

		// Tratamento de erros específicos
		if strings.Contains(err.Error(), "token inválido") || strings.Contains(err.Error(), "expirado") || strings.Contains(err.Error(), "já utilizado") {
			response.WriteError(w, http.StatusUnauthorized, "Token inválido", err.Error())
			return
		}
		if strings.Contains(err.Error(), "pesquisa não está ativa") {
			response.WriteError(w, http.StatusBadRequest, "Pesquisa inativa", err.Error())
			return
		}
		if strings.Contains(err.Error(), "não pertence à pesquisa") {
			response.WriteError(w, http.StatusBadRequest, "Pergunta inválida", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	h.log.WithContext(r.Context()).Info("Respostas submetidas com sucesso: %d", len(respostas))
	response.WriteSuccess(w, http.StatusCreated, "Respostas submetidas com sucesso", nil)
}

// GetRespostaStats retorna estatísticas agregadas de respostas de uma pesquisa
// @Summary Obter estatísticas de respostas por pesquisa
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/respostas/stats [get]
func (h *RespostaHandler) GetRespostaStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	stats, err := h.respostaUseCase.GetAggregatedByPesquisa(r.Context(), pesquisaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Estatísticas obtidas com sucesso", stats)
}

// GetRespostasByPergunta retorna dados agregados de respostas para pergunta específica
// @Summary Obter respostas agregadas por pergunta
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pergunta_id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{pergunta_id}/respostas/aggregated [get]
func (h *RespostaHandler) GetRespostasByPergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	perguntaID, err := strconv.Atoi(vars["pergunta_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pergunta inválido", "ID deve ser um número inteiro")
		return
	}

	aggregatedData, err := h.respostaUseCase.GetAggregatedByPergunta(r.Context(), perguntaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dados agregados obtidos com sucesso", aggregatedData)
}

// GetRespostasByPesquisa retorna dados agregados de todas as respostas de uma pesquisa
// @Summary Obter respostas agregadas por pesquisa
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/respostas/aggregated [get]
func (h *RespostaHandler) GetRespostasByPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	aggregatedData, err := h.respostaUseCase.GetAggregatedByPesquisa(r.Context(), pesquisaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dados agregados da pesquisa obtidos com sucesso", aggregatedData)
}

// GetRespostasByDateRange retorna respostas de pesquisa filtradas por período
// @Summary Obter respostas por período
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Param start_date query string true "Data inicial (RFC3339)"
// @Param end_date query string true "Data final (RFC3339)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/respostas/by-date [get]
func (h *RespostaHandler) GetRespostasByDateRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if strings.TrimSpace(startDate) == "" {
		response.WriteError(w, http.StatusBadRequest, "Data inicial obrigatória", "Parâmetro start_date é obrigatório")
		return
	}

	if strings.TrimSpace(endDate) == "" {
		response.WriteError(w, http.StatusBadRequest, "Data final obrigatória", "Parâmetro end_date é obrigatório")
		return
	}

	respostas, err := h.respostaUseCase.GetResponsesByDateRange(r.Context(), pesquisaID, startDate, endDate)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		if strings.Contains(err.Error(), "formato de data") {
			response.WriteError(w, http.StatusBadRequest, "Formato de data inválido", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Respostas por período obtidas com sucesso", respostas)
}

// CountRespostasByPesquisa retorna número total de respostas de uma pesquisa
// @Summary Contar respostas por pesquisa
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/respostas/count [get]
func (h *RespostaHandler) CountRespostasByPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	count, err := h.respostaUseCase.CountByPesquisa(r.Context(), pesquisaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	countResponse := map[string]interface{}{
		"pesquisa_id":     pesquisaID,
		"total_respostas": count,
	}

	response.WriteSuccess(w, http.StatusOK, "Contagem de respostas obtida com sucesso", countResponse)
}

// CountRespostasByPergunta retorna número total de respostas de pergunta específica
// @Summary Contar respostas por pergunta
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pergunta_id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{pergunta_id}/respostas/count [get]
func (h *RespostaHandler) CountRespostasByPergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	perguntaID, err := strconv.Atoi(vars["pergunta_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pergunta inválido", "ID deve ser um número inteiro")
		return
	}

	count, err := h.respostaUseCase.CountByPergunta(r.Context(), perguntaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	countResponse := map[string]interface{}{
		"pergunta_id":     perguntaID,
		"total_respostas": count,
	}

	response.WriteSuccess(w, http.StatusOK, "Contagem de respostas da pergunta obtida com sucesso", countResponse)
}

// DeleteRespostasByPesquisa remove todas as respostas de uma pesquisa
// @Summary Remover respostas por pesquisa
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pesquisa_id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/pesquisas/{pesquisa_id}/respostas [delete]
func (h *RespostaHandler) DeleteRespostasByPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	motivo := "Exclusão solicitada pelo administrador"

	if err := h.respostaUseCase.DeleteByPesquisa(r.Context(), pesquisaID, userAdminID, motivo); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Respostas da pesquisa deletadas com sucesso", nil)
}

// DeleteDadosPessoaisBySubmissao anonimiza dados pessoais de uma submissão específica (LGPD).
// @Summary Anonimizar dados pessoais da submissão
// @Description Remove dados pessoais de rastreio de uma submissão, preservando as respostas para análise (LGPD).
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param submissao_id path int true "ID da submissão"
// @Param motivo query string false "Motivo da anonimização"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/submissoes/{submissao_id}/dados-pessoais [delete]
func (h *RespostaHandler) DeleteDadosPessoaisBySubmissao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	submissaoID, err := strconv.Atoi(vars["submissao_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da submissão inválido", "ID deve ser um número inteiro")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	motivo := r.URL.Query().Get("motivo")
	if strings.TrimSpace(motivo) == "" {
		motivo = "Solicitação LGPD"
	}

	if err := h.respostaUseCase.DeletePersonalDataBySubmissao(r.Context(), submissaoID, userAdminID, motivo); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Submissão não encontrada", err.Error())
			return
		}
		if strings.Contains(err.Error(), "obrigatório") || strings.Contains(err.Error(), "inválido") {
			response.WriteError(w, http.StatusBadRequest, "Parâmetros inválidos", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dados pessoais da submissão anonimizados com sucesso", nil)
}

// GetStatsByPergunta retorna estatísticas completas de respostas para pergunta específica
// @Summary Obter estatísticas completas por pergunta
// @Tags respostas
// @Produce json
// @Security BearerAuth
// @Param pergunta_id path int true "ID da pergunta"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/perguntas/{pergunta_id}/respostas/stats [get]
func (h *RespostaHandler) GetStatsByPergunta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	perguntaID, err := strconv.Atoi(vars["pergunta_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pergunta inválido", "ID deve ser um número inteiro")
		return
	}

	stats, err := h.respostaUseCase.GetStatisticsByPergunta(r.Context(), perguntaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pergunta não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Estatísticas da pergunta obtidas com sucesso", stats)
}

// validateRespostaCreateRequest valida campos obrigatórios e regras de negócio para criação
func (h *RespostaHandler) validateRespostaCreateRequest(req *dto.RespostaCreateRequest) error {
	if req.IDPergunta <= 0 {
		return fmt.Errorf("ID da pergunta é obrigatório")
	}
	if strings.TrimSpace(req.ValorResposta) == "" {
		return fmt.Errorf("valor da resposta é obrigatório")
	}
	if len(req.ValorResposta) > 2000 {
		return fmt.Errorf("valor da resposta não pode exceder 2000 caracteres")
	}
	return nil
}

// getUserAdminIDFromContext extrai ID do usuário administrativo do contexto da requisição
func (h *RespostaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP extrai endereço IP do cliente considerando proxies
func (h *RespostaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *RespostaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/respostas/submit", h.SubmitRespostas).Methods("POST")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/stats", h.GetRespostaStats).Methods("GET")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/aggregated", h.GetRespostasByPesquisa).Methods("GET")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/by-date", h.GetRespostasByDateRange).Methods("GET")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/count", h.CountRespostasByPesquisa).Methods("GET")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas", h.DeleteRespostasByPesquisa).Methods("DELETE")
	router.HandleFunc("/submissoes/{submissao_id:[0-9]+}/dados-pessoais", h.DeleteDadosPessoaisBySubmissao).Methods("DELETE")
	router.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/aggregated", h.GetRespostasByPergunta).Methods("GET")
	router.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/count", h.CountRespostasByPergunta).Methods("GET")
	router.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/stats", h.GetStatsByPergunta).Methods("GET")
}
