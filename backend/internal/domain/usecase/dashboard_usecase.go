package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"strconv"
	"strings"
	"time"
)

type DashboardUseCase struct {
    repo             repository.DashboardRepository
    pesquisaRepo     repository.PesquisaRepository
    perguntaRepo     repository.PerguntaRepository  // Adicionar
    respostaRepo     repository.RespostaRepository  // Adicionar
    empresaRepo      repository.EmpresaRepository
    logAuditoriaRepo repository.LogAuditoriaRepository
}

func NewDashboardUseCase(
	repo repository.DashboardRepository,
	pesquisaRepo repository.PesquisaRepository,
	empresaRepo repository.EmpresaRepository,
	logRepo repository.LogAuditoriaRepository,
) *DashboardUseCase {
	return &DashboardUseCase{
		repo:             repo,
		pesquisaRepo:     pesquisaRepo,
		empresaRepo:      empresaRepo,
		logAuditoriaRepo: logRepo,
	}
}

// ValidateConfigFiltros valida JSON de configuração de filtros
func (uc *DashboardUseCase) ValidateConfigFiltros(configFiltros *string) error {
	if configFiltros != nil && strings.TrimSpace(*configFiltros) != "" {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(*configFiltros), &config); err != nil {
			return fmt.Errorf("configuração de filtros inválida (JSON malformado): %v", err)
		}
	}
	return nil
}

func (uc *DashboardUseCase) Create(ctx context.Context, dashboard *entity.Dashboard, userAdminID int, enderecoIP string) error {
	// Validações básicas
	if dashboard.IDPesquisa <= 0 {
		return fmt.Errorf("ID da pesquisa é obrigatório")
	}
	
	if strings.TrimSpace(dashboard.Titulo) == "" {
		return fmt.Errorf("título do dashboard é obrigatório")
	}
	
	// Verifica se pesquisa existe
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, dashboard.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	// Verifica se já existe dashboard para esta pesquisa (relação 1:1)
	existing, err := uc.repo.GetByPesquisaID(ctx, dashboard.IDPesquisa)
	if err == nil && existing != nil {
		return fmt.Errorf("já existe um dashboard para esta pesquisa")
	}
	
	// Valida configuração de filtros se fornecida
	if err := uc.ValidateConfigFiltros(dashboard.ConfigFiltros); err != nil {
		return err
	}
	
	// Define valores padrão
	dashboard.DataCriacao = time.Now()
	
	// Se não há configuração, define padrão
	if dashboard.ConfigFiltros == nil {
		defaultConfig := `{"filtros_padrao": true, "periodo_default": "30days"}`
		dashboard.ConfigFiltros = &defaultConfig
	}
	
	if err := uc.repo.Create(ctx, dashboard); err != nil {
		return fmt.Errorf("erro ao criar dashboard: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Dashboard Criado",
			Detalhes:       fmt.Sprintf("Dashboard criado: %s para pesquisa '%s' (ID: %d)", dashboard.Titulo, pesquisa.Titulo, dashboard.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *DashboardUseCase) GetByID(ctx context.Context, id int) (*entity.Dashboard, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID do dashboard deve ser maior que zero")
	}
	
	return uc.repo.GetByID(ctx, id)
}

func (uc *DashboardUseCase) GetByPesquisaID(ctx context.Context, pesquisaID int) (*entity.Dashboard, error) {
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	// Verifica se pesquisa existe
	_, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	return uc.repo.GetByPesquisaID(ctx, pesquisaID)
}

func (uc *DashboardUseCase) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Dashboard, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	// Verifica se empresa existe
	_, err := uc.empresaRepo.GetByID(ctx, empresaID)
	if err != nil {
		return nil, fmt.Errorf("empresa não encontrada: %v", err)
	}
	
	return uc.repo.ListByEmpresa(ctx, empresaID)
}

func (uc *DashboardUseCase) Update(ctx context.Context, dashboard *entity.Dashboard, userAdminID int, enderecoIP string) error {
	// Validações
	if dashboard.ID <= 0 {
		return fmt.Errorf("ID do dashboard inválido")
	}
	
	if strings.TrimSpace(dashboard.Titulo) == "" {
		return fmt.Errorf("título do dashboard é obrigatório")
	}
	
	// Verifica se dashboard existe
	existing, err := uc.repo.GetByID(ctx, dashboard.ID)
	if err != nil {
		return fmt.Errorf("dashboard não encontrado: %v", err)
	}
	
	// Verifica se pesquisa existe (não permite alterar a pesquisa)
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, existing.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa associada não encontrada: %v", err)
	}
	
	// Não permite alterar a pesquisa associada
	dashboard.IDPesquisa = existing.IDPesquisa
	
	// Valida configuração de filtros se fornecida
	if err := uc.ValidateConfigFiltros(dashboard.ConfigFiltros); err != nil {
		return err
	}
	
	// Preserva data de criação
	dashboard.DataCriacao = existing.DataCriacao
	
	if err := uc.repo.Update(ctx, dashboard); err != nil {
		return fmt.Errorf("erro ao atualizar dashboard: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Dashboard Atualizado",
			Detalhes:       fmt.Sprintf("Dashboard atualizado: %s -> %s da pesquisa '%s' (ID: %d)", existing.Titulo, dashboard.Titulo, pesquisa.Titulo, dashboard.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *DashboardUseCase) UpdateConfig(ctx context.Context, dashboardID int, configFiltros string, userAdminID int, enderecoIP string) error {
	if dashboardID <= 0 {
		return fmt.Errorf("ID do dashboard inválido")
	}
	
	// Verifica se dashboard existe
	dashboard, err := uc.repo.GetByID(ctx, dashboardID)
	if err != nil {
		return fmt.Errorf("dashboard não encontrado: %v", err)
	}
	
	// Valida nova configuração
	if err := uc.ValidateConfigFiltros(&configFiltros); err != nil {
		return err
	}
	
	// Atualiza apenas a configuração
	dashboard.ConfigFiltros = &configFiltros
	
	if err := uc.repo.Update(ctx, dashboard); err != nil {
		return fmt.Errorf("erro ao atualizar configuração: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Configuração Dashboard Atualizada",
			Detalhes:       fmt.Sprintf("Configuração atualizada do dashboard: %s (ID: %d)", dashboard.Titulo, dashboard.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *DashboardUseCase) Delete(ctx context.Context, id int, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID do dashboard inválido")
	}
	
	// Busca dashboard para log e validações
	dashboard, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("dashboard não encontrado: %v", err)
	}
	
	// Busca pesquisa associada
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, dashboard.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa associada não encontrada: %v", err)
	}
	
	// Não permite deletar dashboard de pesquisa ativa
	if pesquisa.Status == "Ativa" {
		return fmt.Errorf("não é possível deletar dashboard de pesquisa ativa")
	}
	
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar dashboard: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Dashboard Deletado",
			Detalhes:       fmt.Sprintf("Dashboard deletado: %s da pesquisa '%s' (ID: %d)", dashboard.Titulo, pesquisa.Titulo, dashboard.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

// GenerateReport gera relatório baseado no dashboard
func (uc *DashboardUseCase) GenerateReport(ctx context.Context, dashboardID int, format string, userAdminID int, enderecoIP string) ([]byte, error) {
	if dashboardID <= 0 {
		return nil, fmt.Errorf("ID do dashboard inválido")
	}
	
	validFormats := map[string]bool{
		"pdf":  true,
		"xlsx": true,
		"csv":  true,
	}
	
	if !validFormats[format] {
		return nil, fmt.Errorf("formato inválido: %s. Formatos válidos: pdf, xlsx, csv", format)
	}
	
	// Verifica se dashboard existe
	dashboard, err := uc.repo.GetByID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("dashboard não encontrado: %v", err)
	}
	
	// Verifica se pesquisa tem dados suficientes
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, dashboard.IDPesquisa)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Rascunho" {
		return nil, fmt.Errorf("não é possível gerar relatório de pesquisa em rascunho")
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Relatório Gerado",
			Detalhes:       fmt.Sprintf("Relatório gerado (%s) do dashboard: %s (ID: %d)", format, dashboard.Titulo, dashboard.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	// Aqui seria implementada a lógica de geração do relatório
	// Por enquanto, retorna dados simulados
	reportContent := fmt.Sprintf("Relatório do Dashboard: %s\nFormato: %s\nPesquisa: %s\nGerado em: %s", 
		dashboard.Titulo, format, pesquisa.Titulo, time.Now().Format("2006-01-02 15:04:05"))
	
	return []byte(reportContent), nil
}

func (uc *DashboardUseCase) GetDashboardData(ctx context.Context, dashboardID int, filters string) (interface{}, error) {
    // Buscar dashboard
    dashboard, err := uc.repo.GetByID(ctx, dashboardID)
    if err != nil {
        return nil, fmt.Errorf("dashboard não encontrado: %v", err)
    }

    // Usar método que existe para buscar dados agregados
    respostasAgregadas, err := uc.respostaRepo.GetAggregatedByPesquisa(ctx, dashboard.IDPesquisa)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar respostas agregadas: %v", err)
    }

    // Buscar perguntas da pesquisa  
    perguntas, err := uc.perguntaRepo.ListByPesquisa(ctx, dashboard.IDPesquisa)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar perguntas: %v", err)
    }

    // Processar dados usando dados agregados
    dadosProcessados := make(map[string]interface{})
    
    for _, pergunta := range perguntas {
        if respostasPergunta, exists := respostasAgregadas[pergunta.ID]; exists {
            dadosProcessados[fmt.Sprintf("pergunta_%d", pergunta.ID)] = processarDadosAgregados(pergunta.TipoPergunta, respostasPergunta)
        }
    }

    // Contar total de respostas usando método que existe
    totalRespostas, err := uc.respostaRepo.CountByPesquisa(ctx, dashboard.IDPesquisa)
    if err != nil {
        return nil, fmt.Errorf("erro ao contar respostas: %v", err)
    }

    return map[string]interface{}{
        "dashboard_id":      dashboardID,
        "total_respostas":   totalRespostas,
        "dados_processados": dadosProcessados,
        "ultima_atualizacao": time.Now(),
    }, nil
}

// Função auxiliar para processar dados agregados
func processarDadosAgregados(tipoPergunta string, dadosAgregados map[string]int) map[string]interface{} {
    switch tipoPergunta {
    case "Multipla Escolha":
        return map[string]interface{}{
            "tipo": "multipla_escolha",
            "distribuicao": dadosAgregados,
        }
    case "Escala Numerica":
        return processarEscalaAgregada(dadosAgregados)
    default:
        return map[string]interface{}{
            "tipo": tipoPergunta,
            "dados": dadosAgregados,
        }
    }
}

func processarEscalaAgregada(dados map[string]int) map[string]interface{} {
    var soma, total int
    for valor, count := range dados {
        if v, err := strconv.Atoi(valor); err == nil {
            soma += v * count
            total += count
        }
    }
    
    media := float64(soma) / float64(total)
    
    return map[string]interface{}{
        "tipo": "escala",
        "distribuicao": dados,
        "media": media,
        "total_respostas": total,
    }
}

// Funções auxiliares para processamento
func filtrarRespostasPorPergunta(respostas []*entity.Resposta, perguntaID int) []*entity.Resposta {
    var filtered []*entity.Resposta
    for _, resposta := range respostas {
        if resposta.IDPergunta == perguntaID {
            filtered = append(filtered, resposta)
        }
    }
    return filtered
}

func processarMultiplaEscolha(respostas []*entity.Resposta) map[string]interface{} {
    contadores := make(map[string]int)
    
    for _, resposta := range respostas {
        contadores[resposta.ValorResposta]++
    }
    
    return map[string]interface{}{
        "tipo": "multipla_escolha",
        "total_respostas": len(respostas),
        "distribuicao": contadores,
    }
}

func processarEscala(respostas []*entity.Resposta) map[string]interface{} {
    var valores []float64
    var soma float64
    
    for _, resposta := range respostas {
        if valor, err := strconv.ParseFloat(resposta.ValorResposta, 64); err == nil {
            valores = append(valores, valor)
            soma += valor
        }
    }
    
    media := soma / float64(len(valores))
    
    return map[string]interface{}{
        "tipo": "escala",
        "total_respostas": len(valores),
        "media": media,
        "valores": valores,
    }
}

// RefreshDashboard com cálculo real
func (uc *DashboardUseCase) RefreshDashboard(ctx context.Context, dashboardID, userAdminID int, clientIP string) error {
    // Validações
    if dashboardID <= 0 {
        return fmt.Errorf("ID do dashboard inválido")
    }
    
    // Buscar dashboard
    dashboard, err := uc.repo.GetByID(ctx, dashboardID)
    if err != nil {
        return fmt.Errorf("dashboard não encontrado: %v", err)
    }

    // Recalcular métricas usando método correto
    totalRespostas, err := uc.respostaRepo.CountByPesquisa(ctx, dashboard.IDPesquisa)
    if err != nil {
        return fmt.Errorf("erro ao contar respostas: %v", err)
    }

    // Calcular taxa de participação (precisa saber quantos foram convidados)
    // Por enquanto, usar um valor padrão ou buscar de outra fonte
    taxaParticipacao := float64(totalRespostas) / 100.0 // Exemplo: assumir 100 convidados

    // Atualizar dashboard
    dashboard.TotalRespostas = totalRespostas
    dashboard.TaxaParticipacao = taxaParticipacao
    dashboard.Metricas = map[string]interface{}{
        "ultima_atualizacao": time.Now(),
        "total_respostas": totalRespostas,
        "taxa_participacao": taxaParticipacao,
    }

    // Atualizar no repository
    if err := uc.repo.Update(ctx, dashboard); err != nil {
        return fmt.Errorf("erro ao atualizar dashboard: %v", err)
    }

    // Log de auditoria
    if userAdminID > 0 {
        log := &entity.LogAuditoria{
            IDUserAdmin:    userAdminID,
            TimeStamp:      time.Now(),
            AcaoRealizada:  "Dashboard Atualizado",
            Detalhes:       fmt.Sprintf("Dashboard atualizado: %s (ID: %d)", dashboard.Titulo, dashboard.ID),
            EnderecoIP:     clientIP,
        }
        uc.logAuditoriaRepo.Create(ctx, log)
    }

    return nil
}

// GetDashboardMetrics com dados reais
func (uc *DashboardUseCase) GetDashboardMetrics(ctx context.Context, dashboardID int) (interface{}, error) {
    // Validações
    if dashboardID <= 0 {
        return nil, fmt.Errorf("ID do dashboard inválido")
    }
    
    // Buscar dashboard
    dashboard, err := uc.repo.GetByID(ctx, dashboardID)
    if err != nil {
        return nil, fmt.Errorf("dashboard não encontrado: %v", err)
    }

    // Buscar dados reais usando métodos corretos
    totalRespostas, err := uc.respostaRepo.CountByPesquisa(ctx, dashboard.IDPesquisa)
    if err != nil {
        return nil, fmt.Errorf("erro ao contar respostas: %v", err)
    }

    // Para última resposta, usar GetResponsesByDateRange ou não incluir por enquanto
    // ultimaResposta := time.Now() // Placeholder

    // Calcular métricas por tipo de pergunta usando método correto
    perguntas, err := uc.perguntaRepo.ListByPesquisa(ctx, dashboard.IDPesquisa)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar perguntas: %v", err)
    }

    tiposPergunta := make(map[string]int)
    for _, pergunta := range perguntas {
        tiposPergunta[pergunta.TipoPergunta]++
    }

    return map[string]interface{}{
        "total_respostas": totalRespostas,
        // "data_ultima_resposta": ultimaResposta, // Remover por enquanto
        "taxa_participacao": dashboard.TaxaParticipacao,
        "resumo_estatistico": map[string]interface{}{
            "total_perguntas": len(perguntas),
            "tipos_pergunta": tiposPergunta,
        },
    }, nil
}