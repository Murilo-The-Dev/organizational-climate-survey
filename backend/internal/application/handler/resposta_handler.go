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
func (h *RespostaHandler) SubmitRespostas(w http.ResponseWriter, r *http.Request) {
	var reqs []dto.RespostaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}
	
	if len(reqs) == 0 {
		h.log.WithContext(r.Context()).Info("Nenhuma resposta enviada")
		response.WriteError(w, http.StatusBadRequest, "Lista vazia", "Pelo menos uma resposta deve ser fornecida")
		return
	}
	
	// Validar e converter todas as respostas
	respostas := make([]*entity.Resposta, len(reqs))
	for i, req := range reqs {
		if err := h.validateRespostaCreateRequest(&req); err != nil {
			h.log.WithContext(r.Context()).Info("Validação falhou na resposta %d: %v", i+1, err)
			response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Erro na resposta %d", i+1), err.Error())
			return
		}
		respostas[i] = req.ToEntity()
	}
	
	// Executar caso de uso de criação em lote
	if err := h.respostaUseCase.CreateBatch(r.Context(), respostas); err != nil {
		h.log.WithContext(r.Context()).Error("Erro ao salvar respostas: %v", err)
		if strings.Contains(err.Error(), "pesquisa não está ativa") {
			response.WriteError(w, http.StatusBadRequest, "Pesquisa inativa", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}
	
	h.log.WithContext(r.Context()).Info("Respostas submetidas com sucesso: %d", len(respostas))
	response.WriteSuccess(w, http.StatusCreated, "Respostas submetidas com sucesso", nil)
}

// GetRespostaStats retorna estatísticas agregadas de respostas de uma pesquisa
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

// GetStatsByPergunta retorna estatísticas completas de respostas para pergunta específica
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
	router.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/aggregated", h.GetRespostasByPergunta).Methods("GET")
	router.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/count", h.CountRespostasByPergunta).Methods("GET")
	router.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/stats", h.GetStatsByPergunta).Methods("GET")
}