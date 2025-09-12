package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"strings"
	"time"
)

type PerguntaUseCase struct {
	repo             repository.PerguntaRepository
	respostaRepo     repository.RespostaRepository
	pesquisaRepo     repository.PesquisaRepository
	logAuditoriaRepo repository.LogAuditoriaRepository
}

func NewPerguntaUseCase(
	repo repository.PerguntaRepository,
	respostaRepo repository.RespostaRepository,
	pesquisaRepo repository.PesquisaRepository,
	logRepo repository.LogAuditoriaRepository,
) *PerguntaUseCase {
	return &PerguntaUseCase{
		repo:             repo,
		respostaRepo:     respostaRepo,
		pesquisaRepo:     pesquisaRepo,
		logAuditoriaRepo: logRepo,
	}
}

func (uc *PerguntaUseCase) Create(ctx context.Context, pergunta *entity.Pergunta, userAdminID int, enderecoIP string) error {
	// Validações básicas
	if pergunta.IDPesquisa <= 0 {
		return fmt.Errorf("ID da pesquisa é obrigatório")
	}

	if strings.TrimSpace(pergunta.TextoPergunta) == "" {
		return fmt.Errorf("texto da pergunta é obrigatório")
	}
	
	// Verifica se pesquisa existe
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pergunta.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	// Não permite adicionar perguntas em pesquisas ativas
	if pesquisa.Status == "Ativa" || pesquisa.Status == "Concluída" {
		return fmt.Errorf("não é possível adicionar perguntas em pesquisas ativas ou concluídas")
	}
	
	// Valida tipo de pergunta
	validTipos := map[string]bool{
		"MultiplaEscolha": true,
		"RespostaAberta": true,
		"EscalaNumerica": true,
		"SimNao":         true,
	}
	
	if !validTipos[pergunta.TipoPergunta] {
		return fmt.Errorf("tipo de pergunta inválido: %s", pergunta.TipoPergunta)
	}
	
	// Define ordem se não informada
	if pergunta.OrdemExibicao <= 0 {
		// Busca próxima ordem disponível
		perguntas, err := uc.repo.ListByPesquisa(ctx, pergunta.IDPesquisa)
		if err != nil {
			return fmt.Errorf("erro ao determinar ordem: %v", err)
		}
		pergunta.OrdemExibicao = len(perguntas) + 1
	}
	
	if err := uc.repo.Create(ctx, pergunta); err != nil {
		return fmt.Errorf("erro ao criar pergunta: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Pergunta Criada",
			Detalhes:       fmt.Sprintf("Pergunta criada na pesquisa '%s': %s (ID: %d)", pesquisa.Titulo, pergunta.TextoPergunta, pergunta.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *PerguntaUseCase) CreateBatch(ctx context.Context, perguntas []*entity.Pergunta, userAdminID int, enderecoIP string) error {
	if len(perguntas) == 0 {
		return fmt.Errorf("lista de perguntas não pode estar vazia")
	}
	
	// Validações para todas as perguntas
	pesquisaID := perguntas[0].IDPesquisa
	for i, pergunta := range perguntas {
		if pergunta.IDPesquisa != pesquisaID {
			return fmt.Errorf("todas as perguntas devem pertencer à mesma pesquisa")
		}
		
		if strings.TrimSpace(pergunta.TextoPergunta) == "" {
			return fmt.Errorf("pergunta %d: texto é obrigatório", i+1)
		}
	}
	
	// Verifica se pesquisa existe e se permite edição
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Ativa" || pesquisa.Status == "Concluída" {
		return fmt.Errorf("não é possível adicionar perguntas em pesquisas ativas ou concluídas")
	}
	
	if err := uc.repo.CreateBatch(ctx, perguntas); err != nil {
		return fmt.Errorf("erro ao criar perguntas: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Perguntas Criadas em Lote",
			Detalhes:       fmt.Sprintf("%d perguntas criadas na pesquisa '%s'", len(perguntas), pesquisa.Titulo),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *PerguntaUseCase) GetByID(ctx context.Context, id int) (*entity.Pergunta, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID da pergunta deve ser maior que zero")
	}
	
	return uc.repo.GetByID(ctx, id)
}

func (uc *PerguntaUseCase) ListByPesquisa(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	return uc.repo.ListByPesquisa(ctx, pesquisaID)
}

func (uc *PerguntaUseCase) Update(ctx context.Context, pergunta *entity.Pergunta, userAdminID int, enderecoIP string) error {
	// Validações
	if pergunta.ID <= 0 {
		return fmt.Errorf("ID da pergunta inválido")
	}
	
	if strings.TrimSpace(pergunta.TextoPergunta) == "" {
		return fmt.Errorf("texto da pergunta é obrigatório")
	}
	
	// Verifica se pergunta existe
	existing, err := uc.repo.GetByID(ctx, pergunta.ID)
	if err != nil {
		return fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	// Verifica se pesquisa permite edição
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, existing.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Ativa" || pesquisa.Status == "Concluída" {
		return fmt.Errorf("não é possível editar perguntas de pesquisas ativas ou concluídas")
	}
	
	if err := uc.repo.Update(ctx, pergunta); err != nil {
		return fmt.Errorf("erro ao atualizar pergunta: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Pergunta Atualizada",
			Detalhes:       fmt.Sprintf("Pergunta atualizada na pesquisa '%s' (ID: %d)", pesquisa.Titulo, pergunta.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *PerguntaUseCase) Delete(ctx context.Context, id int, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID da pergunta inválido")
	}
	
	// Busca pergunta para validações e log
	pergunta, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	// Verifica se pesquisa permite edição
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pergunta.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Ativa" || pesquisa.Status == "Concluída" {
		return fmt.Errorf("não é possível deletar perguntas de pesquisas ativas ou concluídas")
	}
	
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar pergunta: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Pergunta Deletada",
			Detalhes:       fmt.Sprintf("Pergunta deletada da pesquisa '%s' (ID: %d)", pesquisa.Titulo, pergunta.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *PerguntaUseCase) UpdateOrdem(ctx context.Context, perguntaID int, novaOrdem int, userAdminID int, enderecoIP string) error {
	if perguntaID <= 0 {
		return fmt.Errorf("ID da pergunta inválido")
	}
	
	if novaOrdem <= 0 {
		return fmt.Errorf("ordem deve ser maior que zero")
	}
	
	// Verifica se pergunta existe
	pergunta, err := uc.repo.GetByID(ctx, perguntaID)
	if err != nil {
		return fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	// Verifica se pesquisa permite edição
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pergunta.IDPesquisa)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Ativa" || pesquisa.Status == "Concluída" {
		return fmt.Errorf("não é possível reordenar perguntas de pesquisas ativas ou concluídas")
	}
	
	if err := uc.repo.UpdateOrdem(ctx, perguntaID, novaOrdem); err != nil {
		return fmt.Errorf("erro ao atualizar ordem: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Ordem Pergunta Alterada",
			Detalhes:       fmt.Sprintf("Ordem alterada para %d - Pergunta ID: %d da pesquisa '%s'", novaOrdem, perguntaID, pesquisa.Titulo),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

// ReorderPerguntas reordena todas as perguntas de uma pesquisa de uma só vez
func (uc *PerguntaUseCase) ReorderPerguntas(ctx context.Context, pesquisaID int, perguntaIDs []int, userAdminID int, enderecoIP string) error {
	if pesquisaID <= 0 {
		return fmt.Errorf("ID da pesquisa inválido")
	}
	
	if len(perguntaIDs) == 0 {
		return fmt.Errorf("lista de IDs não pode estar vazia")
	}
	
	// Verifica se pesquisa existe e permite edição
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Ativa" || pesquisa.Status == "Concluída" {
		return fmt.Errorf("não é possível reordenar perguntas de pesquisas ativas ou concluídas")
	}
	
	// Busca todas as perguntas atuais da pesquisa
	perguntasAtuais, err := uc.repo.ListByPesquisa(ctx, pesquisaID)
	if err != nil {
		return fmt.Errorf("erro ao buscar perguntas: %v", err)
	}
	
	// Valida se todos os IDs pertencem à pesquisa
	perguntaMap := make(map[int]bool)
	for _, p := range perguntasAtuais {
		perguntaMap[p.ID] = true
	}
	
	for _, id := range perguntaIDs {
		if !perguntaMap[id] {
			return fmt.Errorf("pergunta ID %d não pertence à pesquisa %d", id, pesquisaID)
		}
	}
	
	// Verifica se todos os IDs estão incluídos
	if len(perguntaIDs) != len(perguntasAtuais) {
		return fmt.Errorf("todos os IDs das perguntas devem estar incluídos na reordenação")
	}
	
	// Atualiza a ordem de cada pergunta
	for i, perguntaID := range perguntaIDs {
		novaOrdem := i + 1
		if err := uc.repo.UpdateOrdem(ctx, perguntaID, novaOrdem); err != nil {
			return fmt.Errorf("erro ao atualizar ordem da pergunta %d: %v", perguntaID, err)
		}
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Perguntas Reordenadas",
			Detalhes:       fmt.Sprintf("Reordenadas %d perguntas da pesquisa '%s'", len(perguntaIDs), pesquisa.Titulo),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

// PerguntaComEstatisticas representa uma pergunta com suas estatísticas
type PerguntaComEstatisticas struct {
	ID               int                    `json:"id_pergunta"`
	TextoPergunta    string                 `json:"texto_pergunta"`
	TipoPergunta     string                 `json:"tipo_pergunta"`
	OrdemExibicao    int                    `json:"ordem_exibicao"`
	OpcoesResposta   *string                `json:"opcoes_resposta"`
	TotalRespostas   int                    `json:"total_respostas"`
	Estatisticas     map[string]interface{} `json:"estatisticas"`
}

// GetPerguntasWithStats retorna perguntas com estatísticas de respostas
func (uc *PerguntaUseCase) GetPerguntasWithStats(ctx context.Context, pesquisaID int) ([]*PerguntaComEstatisticas, error) {
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa inválido")
	}
	
	// Busca perguntas da pesquisa
	perguntas, err := uc.repo.ListByPesquisa(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar perguntas: %v", err)
	}
	
	result := make([]*PerguntaComEstatisticas, len(perguntas))
	
	for i, pergunta := range perguntas {
		// Conta total de respostas para esta pergunta
		totalRespostas, err := uc.respostaRepo.CountByPergunta(ctx, pergunta.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao contar respostas da pergunta %d: %v", pergunta.ID, err)
		}
		
		// Busca estatísticas agregadas
		stats := make(map[string]interface{})
		
		if totalRespostas > 0 {
			aggregated, err := uc.respostaRepo.GetAggregatedByPergunta(ctx, pergunta.ID)
			if err != nil {
				return nil, fmt.Errorf("erro ao buscar estatísticas da pergunta %d: %v", pergunta.ID, err)
			}
			
			// Processa estatísticas baseadas no tipo de pergunta
			switch pergunta.TipoPergunta {
			case "MultiplaEscolha":
				stats["distribuicao_opcoes"] = aggregated
				stats["opcao_mais_escolhida"] = getMostFrequentOption(aggregated)
				
			case "EscalaNumerica":
				stats["distribuicao_valores"] = aggregated
				stats["media"] = calculateAverage(aggregated)
				stats["valor_mais_comum"] = getMostFrequentOption(aggregated)
				
			case "SimNao":
				stats["distribuicao"] = aggregated
				if sim, exists := aggregated["Sim"]; exists {
					total := float64(totalRespostas)
					stats["percentual_sim"] = float64(sim) / total * 100
					stats["percentual_nao"] = (total - float64(sim)) / total * 100
				}
				
			case "RespostaAberta":
				stats["total_respostas_texto"] = totalRespostas
				// Para respostas abertas, não fazemos agregação automática
			}
		}
		
		result[i] = &PerguntaComEstatisticas{
			ID:             pergunta.ID,
			TextoPergunta:  pergunta.TextoPergunta,
			TipoPergunta:   pergunta.TipoPergunta,
			OrdemExibicao:  pergunta.OrdemExibicao,
			OpcoesResposta: pergunta.OpcoesResposta,
			TotalRespostas: totalRespostas,
			Estatisticas:   stats,
		}
	}
	
	return result, nil
}

// Funções auxiliares para processamento de estatísticas

func getMostFrequentOption(aggregated map[string]int) string {
	maxCount := 0
	mostFrequent := ""
	
	for option, count := range aggregated {
		if count > maxCount {
			maxCount = count
			mostFrequent = option
		}
	}
	
	return mostFrequent
}

func calculateAverage(aggregated map[string]int) float64 {
	total := 0.0
	count := 0.0
	
	for valueStr, freq := range aggregated {
		// Tenta converter string para número
		var value float64
		if _, err := fmt.Sscanf(valueStr, "%f", &value); err == nil {
			total += value * float64(freq)
			count += float64(freq)
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return total / count
}