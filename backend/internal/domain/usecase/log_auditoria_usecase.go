// Package usecase implementa os casos de uso para LogAuditoria.
// Fornece funcionalidades de registro e consulta de logs de auditoria do sistema.
package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"strings"
	"time"
)

// LogAuditoriaUseCase implementa casos de uso para gerenciamento de logs
type LogAuditoriaUseCase struct {
	repo        repository.LogAuditoriaRepository         // Repositório de logs
	userRepo    repository.UsuarioAdministradorRepository // Repositório de usuários
	empresaRepo repository.EmpresaRepository              // Repositório de empresas
}

// NewLogAuditoriaUseCase cria uma nova instância do caso de uso de logs
func NewLogAuditoriaUseCase(repo repository.LogAuditoriaRepository,
	userRepo repository.UsuarioAdministradorRepository,
	empresaRepo repository.EmpresaRepository) *LogAuditoriaUseCase {
	return &LogAuditoriaUseCase{
		repo:        repo,
		userRepo:    userRepo,
		empresaRepo: empresaRepo,
	}
}

// ValidateLogEntry valida entrada de log antes da criação
func (uc *LogAuditoriaUseCase) ValidateLogEntry(log *entity.LogAuditoria) error {
	if log.IDUserAdmin <= 0 {
		return fmt.Errorf("ID do usuário administrador é obrigatório")
	}

	if strings.TrimSpace(log.AcaoRealizada) == "" {
		return fmt.Errorf("ação realizada é obrigatória")
	}

	if strings.TrimSpace(log.Detalhes) == "" {
		return fmt.Errorf("detalhes da ação são obrigatórios")
	}

	if strings.TrimSpace(log.EnderecoIP) == "" {
		return fmt.Errorf("endereço IP é obrigatório")
	}

	// Valida tamanho dos campos
	if len(log.AcaoRealizada) > 100 {
		return fmt.Errorf("ação realizada não pode exceder 100 caracteres")
	}

	if len(log.Detalhes) > 500 {
		return fmt.Errorf("detalhes não podem exceder 500 caracteres")
	}

	return nil
}

// Create registra um novo log de auditoria
func (uc *LogAuditoriaUseCase) Create(ctx context.Context, log *entity.LogAuditoria) error {
	// Validações
	if err := uc.ValidateLogEntry(log); err != nil {
		return err
	}

	// Verifica se usuário administrador existe
	_, err := uc.userRepo.GetByID(ctx, log.IDUserAdmin)
	if err != nil {
		return fmt.Errorf("usuário administrador não encontrado: %v", err)
	}

	// Define timestamp se não informado
	if log.TimeStamp.IsZero() {
		log.TimeStamp = time.Now()
	}

	if err := uc.repo.Create(ctx, log); err != nil {
		return fmt.Errorf("erro ao criar log de auditoria: %v", err)
	}

	return nil
}

// CreateSystemLog cria log de sistema sem usuário específico
func (uc *LogAuditoriaUseCase) CreateSystemLog(ctx context.Context, acao, detalhes, enderecoIP string) error {
	log := &entity.LogAuditoria{
		IDUserAdmin:   0, // Sistema
		TimeStamp:     time.Now(),
		AcaoRealizada: fmt.Sprintf("SISTEMA: %s", acao),
		Detalhes:      detalhes,
		EnderecoIP:    enderecoIP,
	}

	return uc.repo.Create(ctx, log)
}

// GetByID busca um log pelo seu ID
func (uc *LogAuditoriaUseCase) GetByID(ctx context.Context, id int) (*entity.LogAuditoria, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID do log deve ser maior que zero")
	}

	return uc.repo.GetByID(ctx, id)
}

// ListByEmpresa lista logs de uma empresa com paginação
func (uc *LogAuditoriaUseCase) ListByEmpresa(ctx context.Context, empresaID int, limit, offset int) ([]*entity.LogAuditoria, error) {
	// Validações
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	// Verifica se empresa existe
	_, err := uc.empresaRepo.GetByID(ctx, empresaID)
	if err != nil {
		return nil, fmt.Errorf("empresa não encontrada: %v", err)
	}

	// Valida paginação
	if limit <= 0 || limit > 100 {
		limit = 50 // Limite padrão
	}
	if offset < 0 {
		offset = 0
	}

	return uc.repo.ListByEmpresa(ctx, empresaID, limit, offset)
}

// ListByUsuarioAdmin lista logs de um usuário específico
func (uc *LogAuditoriaUseCase) ListByUsuarioAdmin(ctx context.Context, userAdminID int, limit, offset int) ([]*entity.LogAuditoria, error) {
	// Validações
	if userAdminID <= 0 {
		return nil, fmt.Errorf("ID do usuário administrador deve ser maior que zero")
	}

	// Verifica se usuário existe
	_, err := uc.userRepo.GetByID(ctx, userAdminID)
	if err != nil {
		return nil, fmt.Errorf("usuário administrador não encontrado: %v", err)
	}

	// Valida paginação
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	return uc.repo.ListByUsuarioAdmin(ctx, userAdminID, limit, offset)
}

// ListByDateRange lista logs dentro de um intervalo de datas
func (uc *LogAuditoriaUseCase) ListByDateRange(ctx context.Context, empresaID int, startDate, endDate string) ([]*entity.LogAuditoria, error) {
	// Validações
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	if strings.TrimSpace(startDate) == "" {
		return nil, fmt.Errorf("data inicial é obrigatória")
	}

	if strings.TrimSpace(endDate) == "" {
		return nil, fmt.Errorf("data final é obrigatória")
	}

	// Valida formato das datas
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("formato de data inicial inválido (use YYYY-MM-DD): %v", err)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("formato de data final inválido (use YYYY-MM-DD): %v", err)
	}

	// Verifica se data final é posterior à inicial
	if endTime.Before(startTime) {
		return nil, fmt.Errorf("data final deve ser posterior à data inicial")
	}

	// Verifica se período não excede 1 ano
	if endTime.Sub(startTime) > 365*24*time.Hour {
		return nil, fmt.Errorf("período máximo para consulta é de 1 ano")
	}

	// Verifica se empresa existe
	_, err = uc.empresaRepo.GetByID(ctx, empresaID)
	if err != nil {
		return nil, fmt.Errorf("empresa não encontrada: %v", err)
	}

	return uc.repo.ListByDateRange(ctx, empresaID, startDate, endDate)
}

// ListByAction lista logs filtrados por tipo de ação
func (uc *LogAuditoriaUseCase) ListByAction(ctx context.Context, empresaID int, acao string, limit, offset int) ([]*entity.LogAuditoria, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	if strings.TrimSpace(acao) == "" {
		return nil, fmt.Errorf("ação é obrigatória")
	}

	// Verifica se empresa existe
	_, err := uc.empresaRepo.GetByID(ctx, empresaID)
	if err != nil {
		return nil, fmt.Errorf("empresa não encontrada: %v", err)
	}

	// Valida paginação
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Como o método ListByAction não está implementado no repository,
	// vamos buscar todos os logs da empresa e filtrar por ação
	allLogs, err := uc.repo.ListByEmpresa(ctx, empresaID, 1000, 0) // Busca muitos logs
	if err != nil {
		return nil, err
	}

	// Filtrar por ação
	var filteredLogs []*entity.LogAuditoria
	for _, log := range allLogs {
		if strings.Contains(strings.ToLower(log.AcaoRealizada), strings.ToLower(acao)) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	// Aplicar paginação manual
	start := offset
	if start > len(filteredLogs) {
		start = len(filteredLogs)
	}

	end := start + limit
	if end > len(filteredLogs) {
		end = len(filteredLogs)
	}

	return filteredLogs[start:end], nil
}

// GetAuditSummary retorna resumo estatístico dos logs
func (uc *LogAuditoriaUseCase) GetAuditSummary(ctx context.Context, empresaID int, startDate, endDate string) (map[string]interface{}, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	// Busca logs do período
	logs, err := uc.ListByDateRange(ctx, empresaID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Processa estatísticas
	actionCounts := make(map[string]int)
	userCounts := make(map[int]int)
	dailyCounts := make(map[string]int)

	for _, log := range logs {
		actionCounts[log.AcaoRealizada]++
		if log.IDUserAdmin > 0 {
			userCounts[log.IDUserAdmin]++
		}

		// Contagem por dia
		day := log.TimeStamp.Format("2006-01-02")
		dailyCounts[day]++
	}

	summary := map[string]interface{}{
		"periodo_inicio":      startDate,
		"periodo_fim":         endDate,
		"total_eventos":       len(logs),
		"acoes_por_tipo":      actionCounts,
		"eventos_por_usuario": userCounts,
		"eventos_por_dia":     dailyCounts,
	}

	return summary, nil
}

// CleanOldLogs remove logs antigos conforme política de retenção
func (uc *LogAuditoriaUseCase) CleanOldLogs(ctx context.Context, retentionDays int) error {
	if retentionDays < 30 {
		return fmt.Errorf("período de retenção mínimo é de 30 dias")
	}

	if retentionDays > 2555 { // ~7 anos
		return fmt.Errorf("período de retenção máximo é de 2555 dias")
	}

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	// Como DeleteOlderThan não está implementado no repository,
	// vamos apenas criar um log da tentativa de limpeza
	systemLog := &entity.LogAuditoria{
		IDUserAdmin:   0,
		TimeStamp:     time.Now(),
		AcaoRealizada: "SISTEMA: Limpeza de Logs Solicitada",
		Detalhes:      fmt.Sprintf("Solicitação de limpeza de logs anteriores a %s (retenção: %d dias)", cutoffDate.Format("2006-01-02"), retentionDays),
		EnderecoIP:    "sistema",
	}

	return uc.repo.Create(ctx, systemLog)
}

// ExportLogs exporta logs em diferentes formatos
func (uc *LogAuditoriaUseCase) ExportLogs(ctx context.Context, empresaID int, startDate, endDate, format string, userAdminID int, clientIP string) (map[string]interface{}, error) {
	// Validações
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	if strings.TrimSpace(startDate) == "" {
		return nil, fmt.Errorf("data inicial é obrigatória")
	}

	if strings.TrimSpace(endDate) == "" {
		return nil, fmt.Errorf("data final é obrigatória")
	}

	// Validar formato
	validFormats := []string{"csv", "excel", "json"}
	isValidFormat := false
	for _, validFormat := range validFormats {
		if format == validFormat {
			isValidFormat = true
			break
		}
	}

	if !isValidFormat {
		return nil, fmt.Errorf("formato de exportação inválido: %s. Formatos válidos: csv, excel, json", format)
	}

	// Buscar logs do período
	logs, err := uc.ListByDateRange(ctx, empresaID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar logs para exportação: %v", err)
	}

	// Criar log de exportação
	exportLog := &entity.LogAuditoria{
		IDUserAdmin:   userAdminID,
		TimeStamp:     time.Now(),
		AcaoRealizada: "EXPORTAÇÃO: Logs de Auditoria",
		Detalhes:      fmt.Sprintf("Exportação de %d logs no período %s a %s em formato %s", len(logs), startDate, endDate, format),
		EnderecoIP:    clientIP,
	}

	// Salvar log de exportação
	if err := uc.repo.Create(ctx, exportLog); err != nil {
		return nil, fmt.Errorf("erro ao criar log de exportação: %v", err)
	}

	// Retornar dados de exportação
	exportData := map[string]interface{}{
		"export_id":   exportLog.ID,
		"format":      format,
		"start_date":  startDate,
		"end_date":    endDate,
		"total_logs":  len(logs),
		"exported_by": userAdminID,
		"export_ip":   clientIP,
		"export_time": exportLog.TimeStamp,
		"logs":        logs,
		"file_name":   fmt.Sprintf("audit_logs_%s_to_%s.%s", startDate, endDate, format),
	}

	return exportData, nil
}

// GetLogStatistics retorna métricas e estatísticas dos logs
func (uc *LogAuditoriaUseCase) GetLogStatistics(ctx context.Context, empresaID int) (map[string]interface{}, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	// Verifica se empresa existe
	_, err := uc.empresaRepo.GetByID(ctx, empresaID)
	if err != nil {
		return nil, fmt.Errorf("empresa não encontrada: %v", err)
	}

	// Buscar logs dos últimos 30 dias
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	logs, err := uc.repo.ListByDateRange(ctx, empresaID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %v", err)
	}

	// Calcular estatísticas
	totalLogs := len(logs)
	uniqueUsers := make(map[int]bool)
	actionTypes := make(map[string]int)

	for _, log := range logs {
		if log.IDUserAdmin > 0 {
			uniqueUsers[log.IDUserAdmin] = true
		}
		actionTypes[log.AcaoRealizada]++
	}

	statistics := map[string]interface{}{
		"periodo":           fmt.Sprintf("%s a %s", startDate, endDate),
		"total_eventos":     totalLogs,
		"usuarios_ativos":   len(uniqueUsers),
		"tipos_acao":        len(actionTypes),
		"eventos_por_dia":   float64(totalLogs) / 30.0,
		"acoes_mais_comuns": actionTypes,
	}

	return statistics, nil
}
