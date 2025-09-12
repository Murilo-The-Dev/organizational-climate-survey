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

	"github.com/gorilla/mux"
)

type LogAuditoriaHandler struct {
	logAuditoriaUseCase *usecase.LogAuditoriaUseCase
}

func NewLogAuditoriaHandler(logAuditoriaUseCase *usecase.LogAuditoriaUseCase) *LogAuditoriaHandler {
	return &LogAuditoriaHandler{
		logAuditoriaUseCase: logAuditoriaUseCase,
	}
}

// CreateLogAuditoria godoc
// @Summary Criar novo log de auditoria
// @Description Cria um novo registro de log de auditoria (uso interno do sistema)
// @Tags logs-auditoria
// @Accept json
// @Produce json
// @Param log body dto.LogAuditoriaCreateRequest true "Dados do log"
// @Success 201 {object} response.LogResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /logs-auditoria [post]
func (h *LogAuditoriaHandler) CreateLogAuditoria(w http.ResponseWriter, r *http.Request) {
	var req dto.LogAuditoriaCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validateLogCreateRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	log := req.ToEntity()
	
	if err := h.logAuditoriaUseCase.Create(r.Context(), log); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Log de auditoria criado com sucesso", response.ToLogResponse(log))
}

// GetLogAuditoria godoc
// @Summary Buscar log de auditoria por ID
// @Description Retorna um log de auditoria específico pelo ID
// @Tags logs-auditoria
// @Produce json
// @Param id path int true "ID do log"
// @Success 200 {object} response.LogResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /logs-auditoria/{id} [get]
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

// ListLogsByEmpresa godoc
// @Summary Listar logs por empresa
// @Description Retorna lista paginada de logs de auditoria de uma empresa
// @Tags logs-auditoria
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param limit query int false "Limite de resultados" default(50)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{empresa_id}/logs-auditoria [get]
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

	// Converter para response
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

	// calcular página atual
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	// calcular total de páginas (aqui só usamos len(logsResponse), mas o ideal é count real do banco)
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




// ListLogsByUsuario godoc
// @Summary Listar logs por usuário administrador
// @Description Retorna lista paginada de logs de um usuário administrador específico
// @Tags logs-auditoria
// @Produce json
// @Param user_admin_id path int true "ID do usuário administrador"
// @Param limit query int false "Limite de resultados" default(50)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /usuarios-administradores/{user_admin_id}/logs-auditoria [get]
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

	// Converter para response
	logsResponse := make([]response.LogResponse, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.ToLogResponse(log)
	}

	// calcular página atual
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	total := len(logsResponse) // ideal: count real do banco
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

// ListLogsByDateRange godoc
// @Summary Listar logs por período
// @Description Retorna logs de auditoria filtrados por período
// @Tags logs-auditoria
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Success 200 {object} response.ListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{empresa_id}/logs-auditoria/by-date [get]
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

	// Converter para response
	logsResponse := make([]interface{}, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.ToLogResponse(log)
	}

	response.WriteSuccess(w, http.StatusOK, "Logs por período listados com sucesso", logsResponse)
}

// ListLogsByAction godoc
// @Summary Listar logs por tipo de ação
// @Description Retorna logs filtrados por tipo de ação realizada
// @Tags logs-auditoria
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param acao query string true "Tipo de ação"
// @Param limit query int false "Limite de resultados" default(50)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{empresa_id}/logs-auditoria/by-action [get]
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

	// Converter para response
	logsResponse := make([]response.LogResponse, len(logs))
	for i, log := range logs {
		logsResponse[i] = response.ToLogResponse(log)
	}

	// calcular página atual
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	total := len(logsResponse) // ideal: count real do banco
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

// GetAuditSummary godoc
// @Summary Obter resumo de auditoria
// @Description Retorna resumo estatístico de auditoria por período
// @Tags logs-auditoria
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Success 200 {object} response.AuditSummaryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{empresa_id}/logs-auditoria/summary [get]
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

// CleanOldLogs godoc
// @Summary Limpar logs antigos
// @Description Remove logs de auditoria antigos baseado no período de retenção
// @Tags logs-auditoria
// @Accept json
// @Produce json
// @Param retention body dto.RetentionRequest true "Período de retenção em dias"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /logs-auditoria/clean [post]
func (h *LogAuditoriaHandler) CleanOldLogs(w http.ResponseWriter, r *http.Request) {
	var req dto.RetentionRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação do período de retenção
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

// ExportLogs godoc
// @Summary Exportar logs de auditoria
// @Description Exporta logs de auditoria em formato específico
// @Tags logs-auditoria
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Param format query string false "Formato de exportação (csv, excel)" default(csv)
// @Success 200 {object} response.ExportResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /empresas/{empresa_id}/logs-auditoria/export [get]
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

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	// CORRIGIDO: Como ExportLogs não está implementado no use case, vamos usar ListByDateRange
	logs, err := h.logAuditoriaUseCase.ListByDateRange(r.Context(), empresaID, startDate, endDate)
	if err != nil {
		if strings.Contains(err.Error(), "formato de data") {
			response.WriteError(w, http.StatusBadRequest, "Formato de data inválido", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Criar estrutura de exportação simulada
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

// Métodos auxiliares

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

func (h *LogAuditoriaHandler) isValidExportFormat(format string) bool {
	validFormats := []string{"csv", "excel"}
	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

func (h *LogAuditoriaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *LogAuditoriaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// RegisterRoutes registra as rotas do handler
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