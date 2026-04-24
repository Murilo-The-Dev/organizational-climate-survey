package usecase

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"
)

type surveyComparisonEntry struct {
	PesquisaID       int       `json:"pesquisa_id"`
	Titulo           string    `json:"titulo"`
	DataCriacao      time.Time `json:"data_criacao"`
	SetorID          int       `json:"setor_id"`
	TotalRespostas   int       `json:"total_respostas"`
	ScoreMedio       float64   `json:"score_medio"`
	VariacaoAnterior float64   `json:"variacao_anterior"`
}

// ListByEmpresaWithSetor lista dashboards da empresa com filtro opcional por setor.
func (uc *DashboardUseCase) ListByEmpresaWithSetor(ctx context.Context, empresaID int, setorID *int) ([]*entity.Dashboard, error) {
	dashboards, err := uc.ListByEmpresa(ctx, empresaID)
	if err != nil {
		return nil, err
	}

	if setorID == nil {
		return dashboards, nil
	}

	filtered := make([]*entity.Dashboard, 0, len(dashboards))
	for _, dashboard := range dashboards {
		pesquisa, err := uc.pesquisaRepo.GetByID(ctx, dashboard.IDPesquisa)
		if err != nil {
			continue
		}
		if pesquisa.IDSetor == *setorID {
			filtered = append(filtered, dashboard)
		}
	}

	return filtered, nil
}

// GetHistoricalComparison retorna série histórica de score médio para a empresa com filtro opcional por setor.
func (uc *DashboardUseCase) GetHistoricalComparison(ctx context.Context, pesquisaID int, setorID *int) (map[string]interface{}, error) {
	if pesquisaID <= 0 {
		return nil, fmt.Errorf("ID da pesquisa inválido")
	}
	if uc.perguntaRepo == nil || uc.respostaRepo == nil {
		return nil, fmt.Errorf("dependências de perguntas/respostas não inicializadas")
	}

	current, err := uc.pesquisaRepo.GetByID(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("pesquisa não encontrada: %v", err)
	}

	candidates, err := uc.pesquisaRepo.ListByStatus(ctx, current.IDEmpresa, "Concluída")
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pesquisas concluídas: %v", err)
	}

	if current.Status != "Concluída" {
		candidates = append(candidates, current)
	}

	entries := make([]surveyComparisonEntry, 0, len(candidates))
	for _, item := range candidates {
		if setorID != nil && item.IDSetor != *setorID {
			continue
		}
		entry, err := uc.buildSurveyComparisonEntry(ctx, item)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return map[string]interface{}{
			"pesquisa_atual_id": pesquisaID,
			"setor_id":          setorID,
			"comparativo":       []surveyComparisonEntry{},
		}, nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].DataCriacao.Before(entries[j].DataCriacao)
	})

	for i := range entries {
		if i == 0 {
			entries[i].VariacaoAnterior = 0
			continue
		}
		entries[i].VariacaoAnterior = entries[i].ScoreMedio - entries[i-1].ScoreMedio
	}

	return map[string]interface{}{
		"pesquisa_atual_id": current.ID,
		"empresa_id":        current.IDEmpresa,
		"setor_id":          setorID,
		"comparativo":       entries,
	}, nil
}

func (uc *DashboardUseCase) buildSurveyComparisonEntry(ctx context.Context, pesquisa *entity.Pesquisa) (surveyComparisonEntry, error) {
	perguntas, err := uc.perguntaRepo.ListByPesquisa(ctx, pesquisa.ID)
	if err != nil {
		return surveyComparisonEntry{}, err
	}

	totalRespostas, err := uc.respostaRepo.CountByPesquisa(ctx, pesquisa.ID)
	if err != nil {
		return surveyComparisonEntry{}, err
	}

	scoreSum := 0.0
	scoreCount := 0

	for _, pergunta := range perguntas {
		agregados, err := uc.respostaRepo.GetAggregatedByPergunta(ctx, pergunta.ID)
		if err != nil {
			continue
		}

		for value, count := range agregados {
			numeric, ok := normalizeScoreValue(pergunta.TipoPergunta, value)
			if !ok {
				continue
			}
			scoreSum += numeric * float64(count)
			scoreCount += count
		}
	}

	score := 0.0
	if scoreCount > 0 {
		score = scoreSum / float64(scoreCount)
	}

	return surveyComparisonEntry{
		PesquisaID:     pesquisa.ID,
		Titulo:         pesquisa.Titulo,
		DataCriacao:    pesquisa.DataCriacao,
		SetorID:        pesquisa.IDSetor,
		TotalRespostas: totalRespostas,
		ScoreMedio:     score,
	}, nil
}

func normalizeScoreValue(tipoPergunta, value string) (float64, bool) {
	switch strings.ToLower(strings.TrimSpace(tipoPergunta)) {
	case "simnao", "sim_nao", "sim-não", "simnão":
		s := strings.ToLower(strings.TrimSpace(value))
		if s == "sim" {
			return 1.0, true
		}
		if s == "não" || s == "nao" {
			return 0.0, true
		}
		return 0, false
	case "escalanumerica", "escala_numerica", "escala numerica", "escala numérica":
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return 0, false
		}
		return float64(n), true
	default:
		return 0, false
	}
}
