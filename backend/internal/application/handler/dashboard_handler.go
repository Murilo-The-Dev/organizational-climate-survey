// Package handler implementa os controladores HTTP do módulo de Dashboards.
// Responsável por processar requisições, validar entradas e orquestrar casos de uso.

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"fmt"
	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/pkg/logger"

	"github.com/gorilla/mux"
)

// DashboardHandler gerencia as rotas e operações HTTP relacionadas a Dashboards
type DashboardHandler struct {
	dashboardUseCase *usecase.DashboardUseCase // Caso de uso principal para operações de dashboard
	log              logger.Logger             // Logger para registrar eventos e erros
}

// NewDashboardHandler instancia um novo handler de dashboard com dependências injetadas
func NewDashboardHandler(dashboardUseCase *usecase.DashboardUseCase, log logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardUseCase: dashboardUseCase,
		log:              log,
	}
}

// CreateDashboard cria um novo dashboard no sistema
func (h *DashboardHandler) CreateDashboard(w http.ResponseWriter, r *http.Request) {
	var req dto.DashboardCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Falha ao decodificar JSON de entrada
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação de campos obrigatórios e regras de negócio
	if err := h.validateDashboardCreateRequest(&req); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	// Conversão do DTO para entidade de domínio
	dashboard := req.ToEntity()
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	// Execução do caso de uso de criação
	if err := h.dashboardUseCase.Create(r.Context(), dashboard, userAdminID, clientIP); err != nil {
		h.log.WithFields(map[string]interface{}{"user_admin_id": userAdminID, "client_ip": clientIP}).Error("Erro ao criar dashboard: %v", err)
		if strings.Contains(err.Error(), "já existe") {
			response.WriteError(w, http.StatusConflict, "Dashboard já existe", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Sucesso na criação
	h.log.WithFields(map[string]interface{}{"dashboard_id": dashboard.ID, "user_admin_id": userAdminID, "client_ip": clientIP}).Info("Dashboard criado com sucesso")
	response.WriteSuccess(w, http.StatusCreated, "Dashboard criado com sucesso", response.ToDashboardResponse(dashboard))
}

// GetDashboard retorna um dashboard específico pelo ID
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Busca o dashboard pelo ID
	dashboard, err := h.dashboardUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dashboard encontrado", response.ToDashboardResponse(dashboard))
}

// GetDashboardByPesquisa retorna o dashboard vinculado a uma pesquisa específica
func (h *DashboardHandler) GetDashboardByPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pesquisaID, err := strconv.Atoi(vars["pesquisa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da pesquisa inválido", "ID deve ser um número inteiro")
		return
	}

	dashboard, err := h.dashboardUseCase.GetByPesquisaID(r.Context(), pesquisaID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado para esta pesquisa", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dashboard da pesquisa encontrado", response.ToDashboardResponse(dashboard))
}

// ListDashboardsByEmpresa lista todos os dashboards associados a uma empresa
func (h *DashboardHandler) ListDashboardsByEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	dashboards, err := h.dashboardUseCase.ListByEmpresa(r.Context(), empresaID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Conversão para DTOs de resposta
	dashboardsResponse := make([]interface{}, len(dashboards))
	for i, dashboard := range dashboards {
		dashboardsResponse[i] = response.ToDashboardResponse(dashboard)
	}

	response.WriteSuccess(w, http.StatusOK, "Dashboards listados com sucesso", dashboardsResponse)
}

// UpdateDashboard atualiza dados de um dashboard existente
func (h *DashboardHandler) UpdateDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req dto.DashboardUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Recupera o dashboard atual
	dashboard, err := h.dashboardUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Aplica modificações
	req.ApplyToEntity(dashboard)

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.dashboardUseCase.Update(r.Context(), dashboard, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dashboard atualizado com sucesso", response.ToDashboardResponse(dashboard))
}

// DeleteDashboard remove um dashboard do sistema
func (h *DashboardHandler) DeleteDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.dashboardUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dashboard deletado com sucesso", nil)
}

// GetDashboardData retorna dados e métricas processadas do dashboard
func (h *DashboardHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	filters := r.URL.Query().Get("filters")

	dashboardData, err := h.dashboardUseCase.GetDashboardData(r.Context(), id, filters)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dados do dashboard obtidos com sucesso", dashboardData)
}

// RefreshDashboard força atualização dos dados e métricas de um dashboard
func (h *DashboardHandler) RefreshDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.dashboardUseCase.RefreshDashboard(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Dashboard atualizado com sucesso", nil)
}

// ExportDashboard exporta o dashboard em PDF, Excel ou CSV
func (h *DashboardHandler) ExportDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "pdf"
	}

	if !h.isValidExportFormat(format) {
		response.WriteError(w, http.StatusBadRequest, "Formato inválido", "Formato deve ser: pdf ou excel")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	exportData, err := h.dashboardUseCase.GenerateReport(r.Context(), id, format, userAdminID, clientIP)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Configuração dos headers de resposta conforme o tipo de exportação
	switch format {
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dashboard_%d.pdf", id))
	case "xlsx":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dashboard_%d.xlsx", id))
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dashboard_%d.csv", id))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(exportData)
}

// GetDashboardMetrics retorna métricas resumidas de um dashboard
func (h *DashboardHandler) GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	metrics, err := h.dashboardUseCase.GetDashboardMetrics(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Dashboard não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Métricas do dashboard obtidas com sucesso", metrics)
}

// validateDashboardCreateRequest valida regras de negócio e obrigatoriedade de campos
func (h *DashboardHandler) validateDashboardCreateRequest(req *dto.DashboardCreateRequest) error {
	if req.IDPesquisa <= 0 {
		return fmt.Errorf("ID da pesquisa é obrigatório")
	}
	if strings.TrimSpace(req.Titulo) == "" {
		return fmt.Errorf("título do dashboard é obrigatório")
	}
	if len(req.Titulo) < 3 {
		return fmt.Errorf("título deve ter pelo menos 3 caracteres")
	}
	if len(req.Titulo) > 255 {
		return fmt.Errorf("título não pode exceder 255 caracteres")
	}
	return nil
}

// isValidExportFormat valida se o formato de exportação é aceito
func (h *DashboardHandler) isValidExportFormat(format string) bool {
	validFormats := []string{"pdf", "xlsx", "csv"}
	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

// getUserAdminIDFromContext obtém o ID do usuário administrativo autenticado
func (h *DashboardHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP identifica o IP real do cliente, considerando cabeçalhos de proxy
func (h *DashboardHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes associa todas as rotas HTTP deste handler
func (h *DashboardHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dashboards", h.CreateDashboard).Methods("POST")
	router.HandleFunc("/dashboards/{id:[0-9]+}", h.GetDashboard).Methods("GET")
	router.HandleFunc("/dashboards/{id:[0-9]+}", h.UpdateDashboard).Methods("PUT")
	router.HandleFunc("/dashboards/{id:[0-9]+}", h.DeleteDashboard).Methods("DELETE")
	router.HandleFunc("/dashboards/{id:[0-9]+}/data", h.GetDashboardData).Methods("GET")
	router.HandleFunc("/dashboards/{id:[0-9]+}/refresh", h.RefreshDashboard).Methods("POST")
	router.HandleFunc("/dashboards/{id:[0-9]+}/export", h.ExportDashboard).Methods("GET")
	router.HandleFunc("/dashboards/{id:[0-9]+}/metrics", h.GetDashboardMetrics).Methods("GET")
	router.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/dashboard", h.GetDashboardByPesquisa).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/dashboards", h.ListDashboardsByEmpresa).Methods("GET")
}
