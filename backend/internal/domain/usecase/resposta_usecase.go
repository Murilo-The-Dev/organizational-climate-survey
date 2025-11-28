// Package usecase implementa os casos de uso para Respostas.
// Fornece funcionalidades de criação e análise de respostas de pesquisas.
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

// RespostaUseCase implementa casos de uso para gerenciamento de respostas
type RespostaUseCase struct {
	repo              repository.RespostaRepository              // Repositório de respostas
	perguntaRepo      repository.PerguntaRepository              // Repositório de perguntas
	pesquisaRepo      repository.PesquisaRepository              // Repositório de pesquisas
	submissaoUseCase  *SubmissaoPesquisaUseCase                  // NOVO: UseCase de submissões
}

// NewRespostaUseCase cria uma nova instância do caso de uso de respostas
func NewRespostaUseCase(
	repo repository.RespostaRepository,
	perguntaRepo repository.PerguntaRepository,
	pesquisaRepo repository.PesquisaRepository,
	submissaoUseCase *SubmissaoPesquisaUseCase, // NOVO
) *RespostaUseCase {
	return &RespostaUseCase{
		repo:             repo,
		perguntaRepo:     perguntaRepo,
		pesquisaRepo:     pesquisaRepo,
		submissaoUseCase: submissaoUseCase, // NOVO
	}
}

// ValidateResposta valida uma resposta individual
func (uc *RespostaUseCase) ValidateResposta(resposta *entity.Resposta) error {
	if resposta.IDPergunta <= 0 {
		return fmt.Errorf("ID da pergunta é obrigatório")
	}

	if resposta.IDSubmissao <= 0 { // NOVO: validar IDSubmissao
		return fmt.Errorf("ID da submissão é obrigatório")
	}

	if strings.TrimSpace(resposta.ValorResposta) == "" {
		return fmt.Errorf("valor da resposta é obrigatório")
	}

	return nil
}

// CreateBatch cria múltiplas respostas vinculadas a uma submissão anônima
// MODIFICADO: Agora recebe tokenAcesso e valida submissão
func (uc *RespostaUseCase) CreateBatch(ctx context.Context, respostas []*entity.Resposta, tokenAcesso string) error {
	// Validações básicas
	if len(respostas) == 0 {
		return fmt.Errorf("lista de respostas não pode estar vazia")
	}

	if len(respostas) > 100 {
		return fmt.Errorf("máximo de 100 respostas por submissão")
	}

	if strings.TrimSpace(tokenAcesso) == "" {
		return fmt.Errorf("token de acesso é obrigatório")
	}

	// CRÍTICO: Validar token e obter submissão
	submissao, err := uc.submissaoUseCase.ValidateToken(ctx, tokenAcesso)
	if err != nil {
		return fmt.Errorf("token inválido: %v", err)
	}

	// Buscar todas as perguntas da pesquisa para validação
	perguntas, err := uc.perguntaRepo.ListByPesquisa(ctx, submissao.IDPesquisa)
	if err != nil {
		return fmt.Errorf("erro ao buscar perguntas: %v", err)
	}

	// Criar mapa de perguntas válidas
	perguntasValidas := make(map[int]bool)
	for _, p := range perguntas {
		perguntasValidas[p.ID] = true
	}

	// Validar todas as respostas e setar IDSubmissao
	now := time.Now()
	for i, resposta := range respostas {
		// Validação básica
		if err := uc.ValidateResposta(resposta); err != nil {
			return fmt.Errorf("resposta %d inválida: %v", i+1, err)
		}

		// CRÍTICO: Validar que pergunta pertence à pesquisa do token
		if !perguntasValidas[resposta.IDPergunta] {
			return fmt.Errorf("resposta %d: pergunta ID %d não pertence à pesquisa", i+1, resposta.IDPergunta)
		}

		// CRÍTICO: Setar IDSubmissao (vincula ao respondente anônimo)
		resposta.IDSubmissao = submissao.ID

		// Define timestamp se não fornecido
		if resposta.DataSubmissao.IsZero() {
			resposta.DataSubmissao = now
		}

		// Valida valor da resposta baseado no tipo da pergunta
		if err := uc.ValidateResponseValue(ctx, resposta.IDPergunta, resposta.ValorResposta); err != nil {
			return fmt.Errorf("resposta %d: %v", i+1, err)
		}
	}

	// Cria as respostas no banco (transação única)
	if err := uc.repo.CreateBatch(ctx, respostas); err != nil {
		return fmt.Errorf("erro ao salvar respostas: %v", err)
	}

	// CRÍTICO: Marcar submissão como completa
	if err := uc.submissaoUseCase.CompleteSubmission(ctx, submissao.ID); err != nil {
		// Log erro mas não falha - respostas já foram salvas
		log.Printf("AVISO: Respostas salvas mas erro ao marcar submissão como completa (ID %d): %v", submissao.ID, err)
	}

	return nil
}

// REMOVIDO: CreateSingleResponse
// Submissões anônimas sempre em lote vinculadas a um token

// REMOVIDO: ValidateSubmissionRules
// Validação agora feita via ValidateToken do SubmissaoUseCase

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
		"pergunta_id":     perguntaID,
		"tipo_pergunta":   pergunta.TipoPergunta,
		"texto_pergunta":  pergunta.TextoPergunta,
		"total_respostas": totalRespostas,
		"dados_agregados": agregados,
		"opcoes_resposta": pergunta.OpcoesResposta,
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