package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/testutils"
)

func TestDashboardCreateSuccessWithDefaultConfig(t *testing.T) {
	dashboardRepo := &testutils.MockDashboardRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Titulo: "Pesquisa A", Status: "Rascunho"}, nil
	}
	dashboardRepo.GetByPesquisaIDFunc = func(ctx context.Context, id int) (*entity.Dashboard, error) {
		return nil, context.DeadlineExceeded
	}
	dashboardRepo.CreateFunc = func(ctx context.Context, d *entity.Dashboard) error {
		d.ID = 10
		return nil
	}

	uc := NewDashboardUseCase(dashboardRepo, pesquisaRepo, perguntaRepo, respostaRepo, empresaRepo, logRepo)
	d := &entity.Dashboard{IDPesquisa: 1, Titulo: "Dashboard Pesquisa A"}

	if err := uc.Create(context.Background(), d, 1, "127.0.0.1"); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if d.ConfigFiltros == nil || !strings.Contains(*d.ConfigFiltros, "filtros_padrao") {
		t.Fatalf("expected default config to be set, got: %v", d.ConfigFiltros)
	}
}

func TestDashboardListByEmpresaWithSetorFilter(t *testing.T) {
	dashboardRepo := &testutils.MockDashboardRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	empresaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Empresa, error) {
		return &entity.Empresa{ID: id}, nil
	}
	dashboardRepo.ListByEmpresaFunc = func(ctx context.Context, id int) ([]*entity.Dashboard, error) {
		return []*entity.Dashboard{
			{ID: 1, IDPesquisa: 100},
			{ID: 2, IDPesquisa: 101},
			{ID: 3, IDPesquisa: 102},
		}, nil
	}
	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		switch id {
		case 100:
			return &entity.Pesquisa{ID: id, IDSetor: 7}, nil
		case 101:
			return &entity.Pesquisa{ID: id, IDSetor: 8}, nil
		default:
			return &entity.Pesquisa{ID: id, IDSetor: 7}, nil
		}
	}

	uc := NewDashboardUseCase(dashboardRepo, pesquisaRepo, perguntaRepo, respostaRepo, empresaRepo, logRepo)

	all, err := uc.ListByEmpresaWithSetor(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error listing all: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 dashboards, got %d", len(all))
	}

	setorID := 7
	filtered, err := uc.ListByEmpresaWithSetor(context.Background(), 1, &setorID)
	if err != nil {
		t.Fatalf("unexpected error filtering dashboards: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 dashboards, got %d", len(filtered))
	}
}

func TestDashboardHistoricalComparisonSuccess(t *testing.T) {
	dashboardRepo := &testutils.MockDashboardRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	now := time.Now()
	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		switch id {
		case 1:
			return &entity.Pesquisa{ID: 1, IDEmpresa: 55, IDSetor: 7, Status: "Concluída", DataCriacao: now.Add(-48 * time.Hour), Titulo: "P1"}, nil
		case 2:
			return &entity.Pesquisa{ID: 2, IDEmpresa: 55, IDSetor: 7, Status: "Concluída", DataCriacao: now.Add(-24 * time.Hour), Titulo: "P2"}, nil
		default:
			return nil, context.Canceled
		}
	}
	pesquisaRepo.ListByStatusFunc = func(ctx context.Context, empresaID int, status string) ([]*entity.Pesquisa, error) {
		return []*entity.Pesquisa{
			{ID: 1, IDEmpresa: 55, IDSetor: 7, Status: "Concluída", DataCriacao: now.Add(-48 * time.Hour), Titulo: "P1"},
			{ID: 2, IDEmpresa: 55, IDSetor: 7, Status: "Concluída", DataCriacao: now.Add(-24 * time.Hour), Titulo: "P2"},
		}, nil
	}
	perguntaRepo.ListByPesquisaFunc = func(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
		return []*entity.Pergunta{{ID: pesquisaID * 10, IDPesquisa: pesquisaID, TipoPergunta: "EscalaNumerica"}}, nil
	}
	respostaRepo.CountByPesquisaFunc = func(ctx context.Context, pesquisaID int) (int, error) {
		return 5, nil
	}
	respostaRepo.GetAggregatedByPerguntaFunc = func(ctx context.Context, perguntaID int) (map[string]int, error) {
		if perguntaID == 10 {
			return map[string]int{"8": 2, "9": 1}, nil
		}
		return map[string]int{"6": 1, "7": 2}, nil
	}

	uc := NewDashboardUseCase(dashboardRepo, pesquisaRepo, perguntaRepo, respostaRepo, empresaRepo, logRepo)
	result, err := uc.GetHistoricalComparison(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error in historical comparison: %v", err)
	}

	comparativo, ok := result["comparativo"].([]surveyComparisonEntry)
	if !ok {
		t.Fatalf("expected comparativo as []surveyComparisonEntry, got %T", result["comparativo"])
	}
	if len(comparativo) != 2 {
		t.Fatalf("expected 2 comparison entries, got %d", len(comparativo))
	}
}

func TestDashboardGenerateReportAndHelpers(t *testing.T) {
	dashboardRepo := &testutils.MockDashboardRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	dashboardRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Dashboard, error) {
		return &entity.Dashboard{ID: id, IDPesquisa: 22, Titulo: "D22"}, nil
	}
	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Titulo: "Pesquisa 22", Status: "Concluída"}, nil
	}
	perguntaRepo.ListByPesquisaFunc = func(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
		return []*entity.Pergunta{
			{ID: 501, IDPesquisa: pesquisaID, TipoPergunta: "EscalaNumerica", TextoPergunta: "Nota"},
			{ID: 502, IDPesquisa: pesquisaID, TipoPergunta: "SimNao", TextoPergunta: "Recomenda?"},
		}, nil
	}
	respostaRepo.GetAggregatedByPerguntaFunc = func(ctx context.Context, perguntaID int) (map[string]int, error) {
		if perguntaID == 501 {
			return map[string]int{"9": 2, "10": 1}, nil
		}
		return map[string]int{"Sim": 2, "Não": 1}, nil
	}

	uc := NewDashboardUseCase(dashboardRepo, pesquisaRepo, perguntaRepo, respostaRepo, empresaRepo, logRepo)

	csvData, err := uc.GenerateReport(context.Background(), 3, "csv", 10, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected csv export error: %v", err)
	}
	if !strings.Contains(string(csvData), "pergunta_id") {
		t.Fatalf("expected csv header in export, got: %s", string(csvData))
	}

	xlsxData, err := uc.GenerateReport(context.Background(), 3, "xlsx", 10, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected xlsx export error: %v", err)
	}
	if len(xlsxData) == 0 {
		t.Fatal("expected xlsx bytes")
	}

	pdfData, err := uc.GenerateReport(context.Background(), 3, "pdf", 10, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected pdf export error: %v", err)
	}
	if len(pdfData) == 0 {
		t.Fatal("expected pdf bytes")
	}

	if got := computeNumericAverage(map[string]int{"7": 2, "9": 1}); got == "" {
		t.Fatal("expected numeric average string")
	}

	if _, ok := normalizeScoreValue("SimNao", "Sim"); !ok {
		t.Fatal("expected SimNao score normalization to work")
	}
}

func TestDashboardDataMetricsAndRefresh(t *testing.T) {
	dashboardRepo := &testutils.MockDashboardRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	dashboardRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Dashboard, error) {
		return &entity.Dashboard{ID: id, IDPesquisa: 77, Titulo: "D77", TaxaParticipacao: 0.5}, nil
	}
	respostaRepo.GetAggregatedByPesquisaFunc = func(ctx context.Context, pesquisaID int) (map[int]map[string]int, error) {
		return map[int]map[string]int{11: {"7": 2, "9": 1}}, nil
	}
	perguntaRepo.ListByPesquisaFunc = func(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
		return []*entity.Pergunta{{ID: 11, IDPesquisa: pesquisaID, TipoPergunta: "Escala Numerica"}}, nil
	}
	respostaRepo.CountByPesquisaFunc = func(ctx context.Context, pesquisaID int) (int, error) {
		return 3, nil
	}
	updated := false
	dashboardRepo.UpdateFunc = func(ctx context.Context, d *entity.Dashboard) error {
		updated = true
		if d.TotalRespostas != 3 {
			t.Fatalf("expected total respostas=3, got %d", d.TotalRespostas)
		}
		return nil
	}

	uc := NewDashboardUseCase(dashboardRepo, pesquisaRepo, perguntaRepo, respostaRepo, empresaRepo, logRepo)

	data, err := uc.GetDashboardData(context.Background(), 77, "")
	if err != nil {
		t.Fatalf("unexpected GetDashboardData error: %v", err)
	}
	if data == nil {
		t.Fatal("expected dashboard data")
	}

	metrics, err := uc.GetDashboardMetrics(context.Background(), 77)
	if err != nil {
		t.Fatalf("unexpected GetDashboardMetrics error: %v", err)
	}
	if metrics == nil {
		t.Fatal("expected dashboard metrics")
	}

	if err := uc.RefreshDashboard(context.Background(), 77, 9, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected RefreshDashboard error: %v", err)
	}
	if !updated {
		t.Fatal("expected dashboard update call")
	}

	filtered := filtrarRespostasPorPergunta([]*entity.Resposta{{IDPergunta: 1}, {IDPergunta: 2}}, 2)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resposta, got %d", len(filtered))
	}

	multi := processarMultiplaEscolha([]*entity.Resposta{{ValorResposta: "A"}, {ValorResposta: "A"}, {ValorResposta: "B"}})
	if multi["tipo"] != "multipla_escolha" {
		t.Fatalf("unexpected multipla escolha payload: %v", multi)
	}

	escala := processarEscala([]*entity.Resposta{{ValorResposta: "7"}, {ValorResposta: "9"}})
	if escala["tipo"] != "escala" {
		t.Fatalf("unexpected escala payload: %v", escala)
	}
}
