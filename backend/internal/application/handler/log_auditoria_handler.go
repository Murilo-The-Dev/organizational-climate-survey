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
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/pkg/logger"
	"github.com/gorilla/mux"
)

// LogAuditoriaHandler gerencia requisições HTTP relacionadas a logs de auditoria
type LogAuditoriaHandler struct {
	logAuditoriaUseCase *usecase.LogAuditoriaUseCase
	log                 logger.Logger
}

// NewLogAuditoriaHandler cria nova instância do handler de logs de auditoria
func NewLogAuditoriaHandler(logAuditoriaUseCase *usecase.LogAuditoriaUseCase, log logger.Logger) *LogAuditoriaHandler {
	return &LogAuditoriaHandler{
		logAuditoriaUseCase: logAuditoriaUseCase,
		log:                 log,
	}
}

// CreateLogAuditoria cria novo registro de auditoria no sistema
func (h *LogAuditoriaHandler) CreateLogAuditoria(w http.ResponseWriter, r *http.Request) {
	var req dto.LogAuditoriaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithContext(r.Context()).Warn("Decode erro: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}
	if err := h.validateLogCreateRequest(&req); err != nil {
		h.log.WithContext(r.Context()).Info("Validação falhou: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}
	logEntity := req.ToEntity()
	if err := h.logAuditoriaUseCase.Create(r.Context(), logEntity); err != nil {
		h.log.WithContext(r.Context()).Error("Erro ao criar log auditoria: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}
	h.log.WithFields(map[string]interface{}{"log_id": logEntity.ID}).Info("Log auditoria criado com sucesso")
	response.WriteSuccess(w, http.StatusCreated, "Log criado com sucesso", response.ToLogResponse(logEntity))
}

// GetLogAuditoria busca registro de auditoria por ID
func (h *LogAuditoriaHandler) GetLogAuditoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	log, err := h.logAuditoriaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			response.WriteError(w, http.StatusNotFound, "Log não encontrado", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Log encontrado", response.ToLogResponse(log))
}

// ListLogsByEmpresa lista logs de auditoria de empresa com paginação
func (h *LogAuditoriaHandler) ListLogsByEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	limit, offset := h.getPaginationParams(r)

	logs, err := h.logAuditoriaUseCase.ListByEmpresa(r.Context(), empresaID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	logsResponse := make([]response.LogResponse, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.LogResponse{
			ID:            log.ID,
			TimeStamp:     log.TimeStamp,
			AcaoRealizada: log.AcaoRealizada,
			Detalhes:      log.Detalhes,
			EnderecoIP:    log.EnderecoIP,
		}
	}

	// Calcular metadados de paginação
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	total := len(logsResponse)
	totalPages := 1
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	pagination := response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	paginatedResponse := response.NewPaginatedResponse(logsResponse, pagination, "Logs listados com sucesso")

	response.WriteSuccess(w, http.StatusOK, "Logs listados com sucesso", paginatedResponse)
}

// ListLogsByUsuario lista logs de auditoria de usuário específico com paginação
func (h *LogAuditoriaHandler) ListLogsByUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userAdminID, err := strconv.Atoi(vars["user_admin_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID do usuário inválido", "ID deve ser um número inteiro")
		return
	}

	limit, offset := h.getPaginationParams(r)

	logs, err := h.logAuditoriaUseCase.ListByUsuarioAdmin(r.Context(), userAdminID, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	logsResponse := make([]response.LogResponse, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.ToLogResponse(log)
	}

	// Calcular metadados de paginação
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	total := len(logsResponse)
	totalPages := 1
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	pagination := response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	paginatedResponse := response.NewPaginatedResponse(logsResponse, pagination, "Logs do usuário listados com sucesso")

	response.WriteSuccess(w, http.StatusOK, "Logs do usuário listados com sucesso", paginatedResponse)
}

// ListLogsByDateRange lista logs de auditoria filtrados por período
func (h *LogAuditoriaHandler) ListLogsByDateRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
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

	logs, err := h.logAuditoriaUseCase.ListByDateRange(r.Context(), empresaID, startDate, endDate)
	if err != nil {
		if strings.Contains(err.Error(), "formato de data") {
			response.WriteError(w, http.StatusBadRequest, "Formato de data inválido", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	logsResponse := make([]interface{}, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.ToLogResponse(log)
	}

	response.WriteSuccess(w, http.StatusOK, "Logs por período listados com sucesso", logsResponse)
}

// ListLogsByAction lista logs de auditoria filtrados por tipo de ação
func (h *LogAuditoriaHandler) ListLogsByAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	acao := r.URL.Query().Get("acao")
	if strings.TrimSpace(acao) == "" {
		response.WriteError(w, http.StatusBadRequest, "Ação obrigatória", "Parâmetro acao é obrigatório")
		return
	}

	limit, offset := h.getPaginationParams(r)

	logs, err := h.logAuditoriaUseCase.ListByAction(r.Context(), empresaID, acao, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter entidades para DTOs de resposta
	logsResponse := make([]response.LogResponse, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.ToLogResponse(log)
	}

	// Calcular metadados de paginação
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	total := len(logsResponse)
	totalPages := 1
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	pagination := response.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	paginatedResponse := response.NewPaginatedResponse(logsResponse, pagination, "Logs por ação listados com sucesso")

	response.WriteSuccess(w, http.StatusOK, "Logs por ação listados com sucesso", paginatedResponse)
}

// GetAuditSummary retorna resumo estatístico de auditoria por período
func (h *LogAuditoriaHandler) GetAuditSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
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

	summary, err := h.logAuditoriaUseCase.GetAuditSummary(r.Context(), empresaID, startDate, endDate)
	if err != nil {
		if strings.Contains(err.Error(), "formato de data") {
			response.WriteError(w, http.StatusBadRequest, "Formato de data inválido", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Resumo de auditoria obtido com sucesso", summary)
}

// CleanOldLogs remove logs de auditoria antigos baseado em período de retenção
func (h *LogAuditoriaHandler) CleanOldLogs(w http.ResponseWriter, r *http.Request) {
	var req dto.RetentionRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validar período de retenção contra limites do sistema
	if req.RetentionDays < 30 {
		response.WriteError(w, http.StatusBadRequest, "Período inválido", "Período mínimo de retenção é 30 dias")
		return
	}

	if req.RetentionDays > 2555 {
		response.WriteError(w, http.StatusBadRequest, "Período inválido", "Período máximo de retenção é 2555 dias")
		return
	}

	if err := h.logAuditoriaUseCase.CleanOldLogs(r.Context(), req.RetentionDays); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Limpeza de logs antigos realizada com sucesso", nil)
}

// ExportLogs exporta logs de auditoria em formato específico
func (h *LogAuditoriaHandler) ExportLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	format := r.URL.Query().Get("format")

	if strings.TrimSpace(startDate) == "" {
		response.WriteError(w, http.StatusBadRequest, "Data inicial obrigatória", "Parâmetro start_date é obrigatório")
		return
	}

	if strings.TrimSpace(endDate) == "" {
		response.WriteError(w, http.StatusBadRequest, "Data final obrigatória", "Parâmetro end_date é obrigatório")
		return
	}

	if format == "" {
		format = "csv"
	}

	if !h.isValidExportFormat(format) {
		response.WriteError(w, http.StatusBadRequest, "Formato inválido", "Formato deve ser: csv ou excel")
		return
	}

	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	// Buscar logs do período especificado
	logs, err := h.logAuditoriaUseCase.ListByDateRange(r.Context(), empresaID, startDate, endDate)
	if err != nil {
		if strings.Contains(err.Error(), "formato de data") {
			response.WriteError(w, http.StatusBadRequest, "Formato de data inválido", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Preparar dados de exportação com metadados
	exportData := map[string]interface{}{
		"format":      format,
		"start_date":  startDate,
		"end_date":    endDate,
		"total_logs":  len(logs),
		"exported_by": userAdminID,
		"export_ip":   clientIP,
		"logs":        logs,
	}

	response.WriteSuccess(w, http.StatusOK, "Exportação de logs realizada com sucesso", exportData)
}

// validateLogCreateRequest valida campos obrigatórios e regras de negócio para criação
func (h *LogAuditoriaHandler) validateLogCreateRequest(req *dto.LogAuditoriaCreateRequest) error {
	if req.IDUserAdmin <= 0 {
		return fmt.Errorf("ID do usuário administrador é obrigatório")
	}
	if strings.TrimSpace(req.AcaoRealizada) == "" {
		return fmt.Errorf("ação realizada é obrigatória")
	}
	if len(req.AcaoRealizada) < 3 {
		return fmt.Errorf("ação realizada deve ter pelo menos 3 caracteres")
	}
	if len(req.AcaoRealizada) > 255 {
		return fmt.Errorf("ação realizada não pode exceder 255 caracteres")
	}
	return nil
}

// getPaginationParams extrai parâmetros de paginação da query string
func (h *LogAuditoriaHandler) getPaginationParams(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			offset = v
		}
	}

	return limit, offset
}

// isValidExportFormat verifica se formato de exportação é válido
func (h *LogAuditoriaHandler) isValidExportFormat(format string) bool {
	validFormats := []string{"csv", "excel"}
	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

// getUserAdminIDFromContext extrai ID do usuário administrativo do contexto da requisição
func (h *LogAuditoriaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

// getClientIP extrai endereço IP do cliente considerando proxies
func (h *LogAuditoriaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *LogAuditoriaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/logs-auditoria", h.CreateLogAuditoria).Methods("POST")
	router.HandleFunc("/logs-auditoria/{id:[0-9]+}", h.GetLogAuditoria).Methods("GET")
	router.HandleFunc("/logs-auditoria/clean", h.CleanOldLogs).Methods("POST")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/logs-auditoria", h.ListLogsByEmpresa).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/logs-auditoria/by-date", h.ListLogsByDateRange).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/logs-auditoria/by-action", h.ListLogsByAction).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/logs-auditoria/summary", h.GetAuditSummary).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/logs-auditoria/export", h.ExportLogs).Methods("GET")
	router.HandleFunc("/usuarios-administradores/{user_admin_id:[0-9]+}/logs-auditoria", h.ListLogsByUsuario).Methods("GET")
}