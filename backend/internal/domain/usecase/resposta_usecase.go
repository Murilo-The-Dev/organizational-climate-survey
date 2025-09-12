package usecase

import (
	"context"
	"fmt"
	"log"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"strconv"
	"strings"
	"time"
)

type RespostaUseCase struct {
	repo         repository.RespostaRepository
	perguntaRepo repository.PerguntaRepository
	pesquisaRepo repository.PesquisaRepository
}

func NewRespostaUseCase(
	repo repository.RespostaRepository,
	perguntaRepo repository.PerguntaRepository,
	pesquisaRepo repository.PesquisaRepository,
) *RespostaUseCase {
	return &RespostaUseCase{
		repo:         repo,
		perguntaRepo: perguntaRepo,
		pesquisaRepo: pesquisaRepo,
	}
}

// ValidateResposta valida uma resposta individual
func (uc *RespostaUseCase) ValidateResposta(resposta *entity.Resposta) error {
	if resposta.IDPergunta <= 0 {
		return fmt.Errorf("ID da pergunta é obrigatório")
	}
	
	if strings.TrimSpace(resposta.ValorResposta) == "" {
		return fmt.Errorf("valor da resposta é obrigatório")
	}
	
	return nil
}

// ValidateSubmissionRules valida regras de submissão de respostas
func (uc *RespostaUseCase) ValidateSubmissionRules(ctx context.Context, perguntaID int) (*entity.Pesquisa, error) {
	// Busca pergunta
	pergunta, err := uc.perguntaRepo.GetByID(ctx, perguntaID)
	if err != nil {
		return nil, fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	// Busca pesquisa
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pergunta.IDPesquisa)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	// Valida se pesquisa está ativa
	if pesquisa.Status != "Ativa" {
		return nil, fmt.Errorf("pesquisa não está ativa para receber respostas")
	}
	
	// Valida período de abertura/fechamento
	now := time.Now()
	
	if pesquisa.DataAbertura != nil && now.Before(*pesquisa.DataAbertura) {
		return nil, fmt.Errorf("pesquisa ainda não foi aberta para respostas")
	}
	
	if pesquisa.DataFechamento != nil && now.After(*pesquisa.DataFechamento) {
		return nil, fmt.Errorf("período de respostas da pesquisa já foi encerrado")
	}
	
	return pesquisa, nil
}

func (uc *RespostaUseCase) CreateBatch(ctx context.Context, respostas []*entity.Resposta) error {
	// Validações básicas
	if len(respostas) == 0 {
		return fmt.Errorf("lista de respostas não pode estar vazia")
	}
	
	if len(respostas) > 100 {
		return fmt.Errorf("máximo de 100 respostas por submissão")
	}
	
	// Valida todas as respostas e coleta IDs de perguntas
	perguntaIDs := make(map[int]bool)
	
	for i, resposta := range respostas {
		if err := uc.ValidateResposta(resposta); err != nil {
			return fmt.Errorf("resposta %d inválida: %v", i+1, err)
		}
		
		perguntaIDs[resposta.IDPergunta] = true
		
		// Define timestamps se não fornecidos
		now := time.Now()
		if resposta.DataSubmissao.IsZero() {
			resposta.DataSubmissao = now
		}
		if resposta.DataResposta.IsZero() {
			resposta.DataResposta = now
		}
		
		// Valida valor da resposta baseado no tipo da pergunta
		if err := uc.ValidateResponseValue(ctx, resposta.IDPergunta, resposta.ValorResposta); err != nil {
			return fmt.Errorf("resposta %d: %v", i+1, err)
		}
	}
	
	// Valida regras de submissão usando primeira pergunta
	_, err := uc.ValidateSubmissionRules(ctx, respostas[0].IDPergunta)
	if err != nil {
		return err
	}
	
	// Verifica se todas as perguntas pertencem à mesma pesquisa
	for perguntaID := range perguntaIDs {
		pergunta, err := uc.perguntaRepo.GetByID(ctx, perguntaID)
		if err != nil {
			return fmt.Errorf("pergunta ID %d não encontrada: %v", perguntaID, err)
		}
		
		// Define IDPesquisa na resposta se não estiver definido
		for _, resposta := range respostas {
			if resposta.IDPergunta == perguntaID && resposta.IDPesquisa == 0 {
				resposta.IDPesquisa = pergunta.IDPesquisa
			}
		}
	}
	
	// Cria as respostas
	if err := uc.repo.CreateBatch(ctx, respostas); err != nil {
		return fmt.Errorf("erro ao salvar respostas: %v", err)
	}
	
	return nil
}

// CreateSingleResponse cria uma resposta individual
func (uc *RespostaUseCase) CreateSingleResponse(ctx context.Context, resposta *entity.Resposta) error {
	if err := uc.ValidateResposta(resposta); err != nil {
		return err
	}
	
	// Valida regras de submissão
	_, err := uc.ValidateSubmissionRules(ctx, resposta.IDPergunta)
	if err != nil {
		return err
	}
	
	// Define IDPesquisa se não estiver definido
	if resposta.IDPesquisa == 0 {
		pergunta, err := uc.perguntaRepo.GetByID(ctx, resposta.IDPergunta)
		if err != nil {
			return fmt.Errorf("erro ao buscar pergunta: %v", err)
		}
		resposta.IDPesquisa = pergunta.IDPesquisa
	}
	
	// Define timestamps se não fornecidos
	now := time.Now()
	if resposta.DataSubmissao.IsZero() {
		resposta.DataSubmissao = now
	}
	if resposta.DataResposta.IsZero() {
		resposta.DataResposta = now
	}
	
	// Valida valor da resposta
	if err := uc.ValidateResponseValue(ctx, resposta.IDPergunta, resposta.ValorResposta); err != nil {
		return err
	}
	
	// Cria array com uma resposta para usar o método batch
	respostas := []*entity.Resposta{resposta}
	
	return uc.repo.CreateBatch(ctx, respostas)
}

func (uc *RespostaUseCase) CountByPesquisa(ctx context.Context, pesquisaID int) (int, error) {
	if pesquisaID <= 0 {
		return 0, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	// Verifica se pesquisa existe
	_, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return 0, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	return uc.repo.CountByPesquisa(ctx, pesquisaID)
}

func (uc *RespostaUseCase) CountByPergunta(ctx context.Context, perguntaID int) (int, error) {
	if perguntaID <= 0 {
		return 0, fmt.Errorf("ID da pergunta deve ser maior que zero")
	}
	
	// Verifica se pergunta existe
	_, err := uc.perguntaRepo.GetByID(ctx, perguntaID)
	if err != nil {
		return 0, fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	return uc.repo.CountByPergunta(ctx, perguntaID)
}

func (uc *RespostaUseCase) GetAggregatedByPergunta(ctx context.Context, perguntaID int) (map[string]int, error) {
	if perguntaID <= 0 {
		return nil, fmt.Errorf("ID da pergunta deve ser maior que zero")
	}
	
	// Verifica se pergunta existe
	pergunta, err := uc.perguntaRepo.GetByID(ctx, perguntaID)
	if err != nil {
		return nil, fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	// Verifica se pesquisa tem dados suficientes para agregação
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pergunta.IDPesquisa)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	// Só permite agregação de pesquisas ativas ou concluídas
	if pesquisa.Status == "Rascunho" {
		return nil, fmt.Errorf("não é possível agregar dados de pesquisa em rascunho")
	}
	
	return uc.repo.GetAggregatedByPergunta(ctx, perguntaID)
}

func (uc *RespostaUseCase) GetAggregatedByPesquisa(ctx context.Context, pesquisaID int) (map[int]map[string]int, error) {
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	// Verifica se pesquisa existe e permite agregação
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	if pesquisa.Status == "Rascunho" {
		return nil, fmt.Errorf("não é possível agregar dados de pesquisa em rascunho")
	}
	
	return uc.repo.GetAggregatedByPesquisa(ctx, pesquisaID)
}

func (uc *RespostaUseCase) GetResponsesByDateRange(ctx context.Context, pesquisaID int, startDate, endDate string) ([]*entity.Resposta, error) {
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	if strings.TrimSpace(startDate) == "" {
		return nil, fmt.Errorf("data inicial é obrigatória")
	}
	
	if strings.TrimSpace(endDate) == "" {
		return nil, fmt.Errorf("data final é obrigatória")
	}
	
	// Valida formato das datas (assumindo formato YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, fmt.Errorf("formato de data inicial inválido (use YYYY-MM-DD): %v", err)
	}
	
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		return nil, fmt.Errorf("formato de data final inválido (use YYYY-MM-DD): %v", err)
	}
	
	// Verifica se data final é posterior à inicial
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	
	if end.Before(start) {
		return nil, fmt.Errorf("data final deve ser posterior à data inicial")
	}
	
	// Verifica se pesquisa existe
	_, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	return uc.repo.GetResponsesByDateRange(ctx, pesquisaID, startDate, endDate)
}

func (uc *RespostaUseCase) DeleteByPesquisa(ctx context.Context, pesquisaID int, userAdminID int, motivo string) error {
	if pesquisaID <= 0 {
		return fmt.Errorf("ID da pesquisa deve ser maior que zero")
	}
	
	if strings.TrimSpace(motivo) == "" {
		return fmt.Errorf("motivo da exclusão é obrigatório")
	}
	
	// Verifica se pesquisa existe
	pesquisa, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return fmt.Errorf("pesquisa não encontrada: %v", err)
	}
	
	// Verifica quantidade de respostas antes da exclusão (para log)
	count, err := uc.repo.CountByPesquisa(ctx, pesquisaID)
	if err != nil {
		count = 0 // Se falhar na contagem, continua com 0
	}
	
	// Só permite exclusão de respostas se pesquisa não estiver ativa
	if pesquisa.Status == "Ativa" {
		return fmt.Errorf("não é possível excluir respostas de pesquisa ativa")
	}
	
	if err := uc.repo.DeleteByPesquisa(ctx, pesquisaID); err != nil {
		return fmt.Errorf("erro ao excluir respostas: %v", err)
	}
	
	// Log da operação de exclusão
	log.Printf("Respostas excluídas da pesquisa %d: %d respostas removidas. Motivo: %s. Admin ID: %d", 
		pesquisaID, count, motivo, userAdminID)
	
	return nil
}

// GetStatisticsByPergunta retorna estatísticas específicas de uma pergunta
func (uc *RespostaUseCase) GetStatisticsByPergunta(ctx context.Context, perguntaID int) (map[string]interface{}, error) {
	if perguntaID <= 0 {
		return nil, fmt.Errorf("ID da pergunta deve ser maior que zero")
	}
	
	// Verifica se pergunta existe
	pergunta, err := uc.perguntaRepo.GetByID(ctx, perguntaID)
	if err != nil {
		return nil, fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	// Busca contagem total
	totalRespostas, err := uc.repo.CountByPergunta(ctx, perguntaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar respostas: %v", err)
	}
	
	// Busca dados agregados
	agregados, err := uc.repo.GetAggregatedByPergunta(ctx, perguntaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados agregados: %v", err)
	}
	
	stats := map[string]interface{}{
		"pergunta_id":      perguntaID,
		"tipo_pergunta":    pergunta.TipoPergunta,
		"texto_pergunta":   pergunta.TextoPergunta,
		"total_respostas":  totalRespostas,
		"dados_agregados":  agregados,
		"opcoes_resposta":  pergunta.OpcoesResposta,
	}
	
	return stats, nil
}

// ValidateResponseValue valida valor da resposta baseado no tipo da pergunta
func (uc *RespostaUseCase) ValidateResponseValue(ctx context.Context, perguntaID int, valorResposta string) error {
	pergunta, err := uc.perguntaRepo.GetByID(ctx, perguntaID)
	if err != nil {
		return fmt.Errorf("pergunta não encontrada: %v", err)
	}
	
	valorResposta = strings.TrimSpace(valorResposta)
	
	switch pergunta.TipoPergunta {
	case "SimNao":
		if valorResposta != "Sim" && valorResposta != "Não" {
			return fmt.Errorf("resposta deve ser 'Sim' ou 'Não'")
		}
		
	case "EscalaNumerica":
		// Valida se é um número de 1 a 10 (ou outro range)
		if num, err := strconv.Atoi(valorResposta); err != nil || num < 1 || num > 10 {
			return fmt.Errorf("resposta deve ser um valor numérico entre 1 e 10")
		}
		
	case "MultiplaEscolha":
		// Para múltipla escolha, seria necessário validar contra as opções disponíveis
		// Seria necessário parsear o JSON de OpcoesResposta
		if valorResposta == "" {
			return fmt.Errorf("uma opção deve ser selecionada")
		}
		// TODO: Validar contra opções específicas no JSON
		
	case "RespostaAberta":
		if len(valorResposta) > 1000 {
			return fmt.Errorf("resposta de texto livre não pode exceder 1000 caracteres")
		}
		if len(valorResposta) < 1 {
			return fmt.Errorf("resposta de texto livre não pode estar vazia")
		}
		
	default:
		return fmt.Errorf("tipo de pergunta não reconhecido: %s", pergunta.TipoPergunta)
	}
	
	return nil
}