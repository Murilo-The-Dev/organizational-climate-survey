package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"time"
)

type AnalyticsUseCase struct {
	repo             repository.AnalyticsRepository
	pesquisaRepo     repository.PesquisaRepository
	logAuditoriaRepo repository.LogAuditoriaRepository
}

func NewAnalyticsUseCase(
	repo repository.AnalyticsRepository,
	pesquisaRepo repository.PesquisaRepository,
	logRepo repository.LogAuditoriaRepository,
) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		repo:             repo,
		pesquisaRepo:     pesquisaRepo,
		logAuditoriaRepo: logRepo,
	}
}

func (uc *AnalyticsUseCase) GetPesquisaMetrics(ctx context.Context, pesquisaID int, userAdminID int, enderecoIP string) (map[string]interface{}, error) {
	// Validações
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	// Verifica se pesquisa existe
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	// Só permite visualizar métricas de pesquisas ativas ou concluídas
	if pesquisa.Status == "Rascunho" {
		return nil, fmt.Errorf("não é possível visualizar métricas de pesquisa em rascunho")
	}
	
	metrics, err := uc.repo.GetPesquisaMetrics(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar métricas: %v", err)
	}
	
	// Log de auditoria para acesso às métricas
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Métricas Acessadas",
			Detalhes:       fmt.Sprintf("Métricas acessadas da pesquisa: %s (ID: %d)", pesquisa.Titulo, pesquisaID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return metrics, nil
}

func (uc *AnalyticsUseCase) GetComparisonData(ctx context.Context, pesquisaIDs []int, userAdminID int, enderecoIP string) (map[string]interface{}, error) {
	// Validações
	if len(pesquisaIDs) == 0 {
		return nil, fmt.Errorf("lista de pesquisas não pode estar vazia")
	}
	
	if len(pesquisaIDs) > 10 {
		return nil, fmt.Errorf("máximo de 10 pesquisas para comparação")
	}
	
	// Valida IDs
	for i, id := range pesquisaIDs {
		if id <= 0 {
			return nil, fmt.Errorf("ID inválido na posição %d: %d", i, id)
		}
	}
	
	// Verifica se todas as pesquisas existem e pertencem à mesma empresa
	var empresaID int
	pesquisaTitulos := make([]string, 0, len(pesquisaIDs))
	
	for i, pesquisaID := range pesquisaIDs {
		pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
		if err != nil {
			return nil, fmt.Errorf("pesquisa ID %d não encontrada: %v", pesquisaID, err)
		}
		
		// Primeira pesquisa define a empresa
		if i == 0 {
			empresaID = pesquisa.IDEmpresa
		} else if pesquisa.IDEmpresa != empresaID {
			return nil, fmt.Errorf("todas as pesquisas devem pertencer à mesma empresa")
		}
		
		// Só permite comparar pesquisas concluídas
		if pesquisa.Status != "Concluída" {
			return nil, fmt.Errorf("só é possível comparar pesquisas concluídas. Pesquisa '%s' está com status: %s", pesquisa.Titulo, pesquisa.Status)
		}
		
		pesquisaTitulos = append(pesquisaTitulos, pesquisa.Titulo)
	}
	
	comparison, err := uc.repo.GetComparisonData(ctx, pesquisaIDs)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar comparação: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Comparação de Pesquisas Gerada",
			Detalhes:       fmt.Sprintf("Comparação gerada entre %d pesquisas da empresa ID: %d", len(pesquisaIDs), empresaID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return comparison, nil
}

func (uc *AnalyticsUseCase) GetSetorComparison(ctx context.Context, empresaID int, pesquisaID int, userAdminID int, enderecoIP string) (map[string]interface{}, error) {
	// Validações
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	// Verifica se pesquisa existe e pertence à empresa
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.IDEmpresa != empresaID {
		return nil, fmt.Errorf("pesquisa não pertence à empresa informada")
	}
	
	// Só permite comparação de setores se pesquisa não tem setor específico
	if pesquisa.IDSetor > 0 {
		return nil, fmt.Errorf("pesquisa é específica de um setor, não é possível fazer comparação entre setores")
	}
	
	// Pesquisa deve estar concluída para comparação entre setores
	if pesquisa.Status != "Concluída" {
		return nil, fmt.Errorf("só é possível comparar setores em pesquisas concluídas")
	}
	
	comparison, err := uc.repo.GetSetorComparison(ctx, empresaID, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar comparação por setor: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Comparação por Setor Gerada",
			Detalhes:       fmt.Sprintf("Comparação por setor gerada para pesquisa: %s (ID: %d)", pesquisa.Titulo, pesquisaID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return comparison, nil
}

// GetTrendAnalysis - análise de tendências ao longo do tempo
func (uc *AnalyticsUseCase) GetTrendAnalysis(ctx context.Context, empresaID int, period string, userAdminID int, enderecoIP string) (map[string]interface{}, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	validPeriods := map[string]bool{
		"7days":   true,
		"30days":  true,
		"90days":  true,
		"1year":   true,
	}
	
	if !validPeriods[period] {
		return nil, fmt.Errorf("período inválido: %s. Valores válidos: 7days, 30days, 90days, 1year", period)
	}
	
	// Este método precisaria ser implementado no repository
	// trends, err := uc.repo.GetTrendAnalysis(ctx, empresaID, period)
	// if err != nil {
	//     return nil, fmt.Errorf("erro ao gerar análise de tendências: %v", err)
	// }
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Análise de Tendências Acessada",
			Detalhes:       fmt.Sprintf("Análise de tendências acessada para empresa ID: %d (período: %s)", empresaID, period),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	// Retorna estrutura básica por enquanto
	return map[string]interface{}{
		"period": period,
		"empresa_id": empresaID,
		"message": "Análise de tendências em desenvolvimento",
	}, nil
}