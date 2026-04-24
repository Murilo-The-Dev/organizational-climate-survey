package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/testutils"
	"organizational-climate-survey/backend/pkg/crypto"
)

func TestEmpresaUseCaseFlows(t *testing.T) {
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	empresaRepo.GetByCNPJFunc = func(ctx context.Context, cnpj string) (*entity.Empresa, error) {
		return nil, context.Canceled
	}
	created := false
	empresaRepo.CreateFunc = func(ctx context.Context, e *entity.Empresa) error {
		created = true
		e.ID = 12
		return nil
	}
	empresaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Empresa, error) {
		return &entity.Empresa{ID: id, NomeFantasia: "Antigo", RazaoSocial: "Antigo LTDA", CNPJ: validCNPJFromBase("123456780001")}, nil
	}
	empresaRepo.UpdateFunc = func(ctx context.Context, e *entity.Empresa) error {
		return nil
	}
	deleted := false
	empresaRepo.DeleteFunc = func(ctx context.Context, id int) error {
		deleted = true
		return nil
	}

	uc := NewEmpresaUseCase(empresaRepo, logRepo)
	empresa := &entity.Empresa{
		NomeFantasia: "Empresa A",
		RazaoSocial:  "Empresa A LTDA",
		CNPJ:         validCNPJFromBase("123456780001"),
	}

	if err := uc.Create(context.Background(), empresa, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected empresa create error: %v", err)
	}
	if !created {
		t.Fatal("expected empresa create repository call")
	}

	empresa.ID = 12
	if err := uc.Update(context.Background(), empresa, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected empresa update error: %v", err)
	}
	if err := uc.Delete(context.Background(), 12, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected empresa delete error: %v", err)
	}
	if !deleted {
		t.Fatal("expected empresa delete repository call")
	}
}

func TestSetorUseCaseFlows(t *testing.T) {
	setorRepo := &testutils.MockSetorRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	empresaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Empresa, error) {
		return &entity.Empresa{ID: id}, nil
	}
	setorRepo.GetByNomeFunc = func(ctx context.Context, empresaID int, nome string) (*entity.Setor, error) {
		return nil, context.Canceled
	}
	created := false
	setorRepo.CreateFunc = func(ctx context.Context, s *entity.Setor) error {
		created = true
		s.ID = 33
		return nil
	}
	setorRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Setor, error) {
		return &entity.Setor{ID: id, IDEmpresa: 1, NomeSetor: "RH"}, nil
	}
	setorRepo.UpdateFunc = func(ctx context.Context, s *entity.Setor) error { return nil }
	setorRepo.DeleteFunc = func(ctx context.Context, id int) error { return nil }

	uc := NewSetorUseCase(setorRepo, empresaRepo, logRepo)
	setor := &entity.Setor{IDEmpresa: 1, NomeSetor: "RH"}

	if err := uc.Create(context.Background(), setor, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected setor create error: %v", err)
	}
	if !created {
		t.Fatal("expected setor create call")
	}
	setor.ID = 33
	if err := uc.Update(context.Background(), setor, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected setor update error: %v", err)
	}
	if err := uc.Delete(context.Background(), 33, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected setor delete error: %v", err)
	}
}

func TestAnalyticsUseCaseFlows(t *testing.T) {
	analyticsRepo := &testutils.MockAnalyticsRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		switch id {
		case 1:
			return &entity.Pesquisa{ID: 1, IDEmpresa: 10, IDSetor: 0, Titulo: "P1", Status: "Concluída"}, nil
		case 2:
			return &entity.Pesquisa{ID: 2, IDEmpresa: 10, IDSetor: 0, Titulo: "P2", Status: "Concluída"}, nil
		default:
			return &entity.Pesquisa{ID: id, IDEmpresa: 99, Titulo: "P3", Status: "Ativa"}, nil
		}
	}
	analyticsRepo.GetPesquisaMetricsFunc = func(ctx context.Context, pesquisaID int) (map[string]interface{}, error) {
		return map[string]interface{}{"total": 12}, nil
	}
	analyticsRepo.GetComparisonDataFunc = func(ctx context.Context, pesquisaIDs []int) (map[string]interface{}, error) {
		return map[string]interface{}{"comp": len(pesquisaIDs)}, nil
	}
	analyticsRepo.GetSetorComparisonFunc = func(ctx context.Context, empresaID int, pesquisaID int) (map[string]interface{}, error) {
		return map[string]interface{}{"empresa_id": empresaID, "pesquisa_id": pesquisaID}, nil
	}

	uc := NewAnalyticsUseCase(analyticsRepo, pesquisaRepo, logRepo)

	if _, err := uc.GetPesquisaMetrics(context.Background(), 1, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected metrics error: %v", err)
	}
	if _, err := uc.GetComparisonData(context.Background(), []int{1, 2}, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected comparison error: %v", err)
	}
	if _, err := uc.GetSetorComparison(context.Background(), 10, 1, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected setor comparison error: %v", err)
	}
	if _, err := uc.GetTrendAnalysis(context.Background(), 10, "30days", 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected trend analysis error: %v", err)
	}

	if _, err := uc.GetComparisonData(context.Background(), []int{1, 3}, 1, "127.0.0.1"); err == nil {
		t.Fatal("expected cross-company comparison error")
	}
}

func TestBootstrapUseCaseInitializeSystem(t *testing.T) {
	empresaRepo := &testutils.MockEmpresaRepository{}
	usuarioRepo := &testutils.MockUsuarioAdministradorRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}
	cryptoSvc := crypto.NewDefaultCryptoService()

	usuarioRepo.CountFunc = func(ctx context.Context) (int, error) { return 0, nil }
	empresaRepo.GetByCNPJFunc = func(ctx context.Context, cnpj string) (*entity.Empresa, error) {
		return nil, context.Canceled
	}
	empresaRepo.CreateFunc = func(ctx context.Context, e *entity.Empresa) error {
		e.ID = 5
		return nil
	}
	usuarioRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.UsuarioAdministrador, error) {
		return nil, context.Canceled
	}
	createdAdmin := false
	usuarioRepo.CreateFunc = func(ctx context.Context, u *entity.UsuarioAdministrador) error {
		createdAdmin = true
		u.ID = 8
		if !strings.HasPrefix(u.SenhaHash, "$2") {
			t.Fatalf("expected hashed password, got %s", u.SenhaHash)
		}
		return nil
	}

	uc := NewBootstrapUseCase(empresaRepo, usuarioRepo, logRepo, cryptoSvc)
	data := &BootstrapData{
		Empresa: &entity.Empresa{NomeFantasia: "Bootstrap Co", RazaoSocial: "Bootstrap Co Ltda", CNPJ: validCNPJFromBase("987654320001")},
		Usuario: &entity.UsuarioAdministrador{NomeAdmin: "Admin Root", Email: "root@example.com", SenhaHash: "SenhaForte123!"},
	}

	if err := uc.InitializeSystem(context.Background(), data); err != nil {
		t.Fatalf("unexpected bootstrap error: %v", err)
	}
	if !createdAdmin {
		t.Fatal("expected bootstrap admin creation")
	}
	if data.Usuario.IDEmpresa != 5 {
		t.Fatalf("expected user company id 5, got %d", data.Usuario.IDEmpresa)
	}

	usuarioRepo.CountFunc = func(ctx context.Context) (int, error) { return 1, nil }
	if err := uc.InitializeSystem(context.Background(), data); err == nil {
		t.Fatal("expected already-initialized error")
	}

	if !uc.isValidEmailFormat("a@b.com") || !uc.isValidCNPJFormat("12.345.678/0001-95") {
		t.Fatal("expected simple format helpers to return true")
	}
}

func TestUsuarioAdministradorUseCaseFlows(t *testing.T) {
	usuarioRepo := &testutils.MockUsuarioAdministradorRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}
	cryptoSvc := crypto.NewDefaultCryptoService()

	hashed, _ := cryptoSvc.HashPassword("SenhaForte123!")
	usuarioRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.UsuarioAdministrador, error) {
		if email == "admin@example.com" {
			return &entity.UsuarioAdministrador{ID: 1, IDEmpresa: 1, NomeAdmin: "Admin", Email: email, SenhaHash: hashed, Status: "Ativo"}, nil
		}
		return nil, context.Canceled
	}
	usuarioRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.UsuarioAdministrador, error) {
		return &entity.UsuarioAdministrador{ID: id, IDEmpresa: 1, NomeAdmin: "Admin", Email: "admin@example.com", SenhaHash: hashed, Status: "Ativo"}, nil
	}
	empresaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Empresa, error) {
		return &entity.Empresa{ID: id}, nil
	}
	usuarioRepo.CreateFunc = func(ctx context.Context, u *entity.UsuarioAdministrador) error {
		u.ID = 3
		return nil
	}
	usuarioRepo.UpdateFunc = func(ctx context.Context, u *entity.UsuarioAdministrador) error { return nil }
	usuarioRepo.UpdatePasswordFunc = func(ctx context.Context, userID int, hash string) error { return nil }
	usuarioRepo.UpdateStatusFunc = func(ctx context.Context, id int, status string) error { return nil }
	usuarioRepo.CountFunc = func(ctx context.Context) (int, error) { return 0, nil }

	uc := NewUsuarioAdministradorUseCase(usuarioRepo, empresaRepo, logRepo, cryptoSvc)

	if _, err := uc.Authenticate(context.Background(), "admin@example.com", "SenhaForte123!", "127.0.0.1"); err != nil {
		t.Fatalf("unexpected authenticate error: %v", err)
	}
	if _, err := uc.Authenticate(context.Background(), "admin@example.com", "senha-errada", "127.0.0.1"); err == nil {
		t.Fatal("expected invalid credentials error")
	}

	newUser := &entity.UsuarioAdministrador{IDEmpresa: 1, NomeAdmin: "Novo", Email: "novo@example.com", SenhaHash: "SenhaForte123!"}
	if err := uc.Create(context.Background(), newUser, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected user create error: %v", err)
	}
	if err := uc.UpdatePassword(context.Background(), 1, "NovaSenha123!", 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected update password error: %v", err)
	}

	if err := uc.Delete(context.Background(), 1, 1, "127.0.0.1"); err != nil {
		t.Fatalf("unexpected soft delete error: %v", err)
	}

	bootstrapUser := &entity.UsuarioAdministrador{IDEmpresa: 1, NomeAdmin: "Bootstrap", Email: "boot@example.com", SenhaHash: "SenhaForte123!"}
	if err := uc.CreateBootstrap(context.Background(), bootstrapUser); err != nil {
		t.Fatalf("unexpected create bootstrap error: %v", err)
	}

	if _, err := uc.Count(context.Background()); err != nil {
		t.Fatalf("unexpected count error: %v", err)
	}
}

func TestLogAuditoriaUseCaseFlows(t *testing.T) {
	logRepo := &testutils.MockLogAuditoriaRepository{}
	userRepo := &testutils.MockUsuarioAdministradorRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}

	userRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.UsuarioAdministrador, error) {
		return &entity.UsuarioAdministrador{ID: id}, nil
	}
	empresaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Empresa, error) {
		return &entity.Empresa{ID: id}, nil
	}
	stored := make([]*entity.LogAuditoria, 0)
	logRepo.CreateFunc = func(ctx context.Context, l *entity.LogAuditoria) error {
		stored = append(stored, l)
		if l.ID == 0 {
			l.ID = len(stored)
		}
		return nil
	}
	logRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.LogAuditoria, error) {
		return &entity.LogAuditoria{ID: id, IDUserAdmin: 1, TimeStamp: time.Now(), AcaoRealizada: "Teste", Detalhes: "det", EnderecoIP: "127.0.0.1"}, nil
	}
	logRepo.ListByDateRangeFunc = func(ctx context.Context, empresaID int, startDate, endDate string) ([]*entity.LogAuditoria, error) {
		return []*entity.LogAuditoria{
			{ID: 1, IDUserAdmin: 1, TimeStamp: time.Now(), AcaoRealizada: "Login", Detalhes: "ok", EnderecoIP: "127.0.0.1"},
			{ID: 2, IDUserAdmin: 2, TimeStamp: time.Now(), AcaoRealizada: "Export", Detalhes: "ok", EnderecoIP: "127.0.0.1"},
		}, nil
	}
	logRepo.ListByEmpresaFunc = func(ctx context.Context, empresaID int, limit, offset int) ([]*entity.LogAuditoria, error) {
		return []*entity.LogAuditoria{
			{AcaoRealizada: "Login"},
			{AcaoRealizada: "Export"},
			{AcaoRealizada: "Login com MFA"},
		}, nil
	}
	logRepo.ListByUsuarioAdminFunc = func(ctx context.Context, userAdminID int, limit, offset int) ([]*entity.LogAuditoria, error) {
		return []*entity.LogAuditoria{{IDUserAdmin: userAdminID}}, nil
	}

	uc := NewLogAuditoriaUseCase(logRepo, userRepo, empresaRepo)

	err := uc.Create(context.Background(), &entity.LogAuditoria{IDUserAdmin: 1, AcaoRealizada: "Ação", Detalhes: "Detalhes", EnderecoIP: "127.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected log create error: %v", err)
	}
	if err := uc.CreateSystemLog(context.Background(), "Job", "Run", "127.0.0.1"); err != nil {
		t.Fatalf("unexpected system log error: %v", err)
	}
	if _, err := uc.GetByID(context.Background(), 1); err != nil {
		t.Fatalf("unexpected get by id error: %v", err)
	}
	if _, err := uc.ListByEmpresa(context.Background(), 1, 10, 0); err != nil {
		t.Fatalf("unexpected list by empresa error: %v", err)
	}
	if _, err := uc.ListByUsuarioAdmin(context.Background(), 1, 10, 0); err != nil {
		t.Fatalf("unexpected list by user error: %v", err)
	}
	if _, err := uc.ListByDateRange(context.Background(), 1, "2024-01-01", "2024-01-31"); err != nil {
		t.Fatalf("unexpected list by date range error: %v", err)
	}
	if _, err := uc.ListByAction(context.Background(), 1, "login", 10, 0); err != nil {
		t.Fatalf("unexpected list by action error: %v", err)
	}
	if _, err := uc.GetAuditSummary(context.Background(), 1, "2024-01-01", "2024-01-31"); err != nil {
		t.Fatalf("unexpected summary error: %v", err)
	}
	if err := uc.CleanOldLogs(context.Background(), 60); err != nil {
		t.Fatalf("unexpected clean old logs error: %v", err)
	}

	exportData, err := uc.ExportLogs(context.Background(), 1, "2024-01-01", "2024-01-31", "json", 1, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected export logs error: %v", err)
	}
	if exportData["total_logs"].(int) != 2 {
		t.Fatalf("expected 2 logs exported, got %v", exportData["total_logs"])
	}

	if _, err := uc.GetLogStatistics(context.Background(), 1); err != nil {
		t.Fatalf("unexpected log statistics error: %v", err)
	}

	if _, err := uc.ExportLogs(context.Background(), 1, "2024-01-01", "2024-01-31", "xml", 1, "127.0.0.1"); err == nil {
		t.Fatal("expected invalid export format error")
	}
}

func TestUseCaseValidationBranches(t *testing.T) {
	analyticsRepo := &testutils.MockAnalyticsRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}
	analyticsUC := NewAnalyticsUseCase(analyticsRepo, pesquisaRepo, logRepo)
	if _, err := analyticsUC.GetTrendAnalysis(context.Background(), 1, "2days", 1, "127.0.0.1"); err == nil {
		t.Fatal("expected invalid trend period error")
	}

	dashboardUC := NewDashboardUseCase(&testutils.MockDashboardRepository{}, &testutils.MockPesquisaRepository{}, &testutils.MockPerguntaRepository{}, &testutils.MockRespostaRepository{}, &testutils.MockEmpresaRepository{}, &testutils.MockLogAuditoriaRepository{})
	invalidJSON := "{invalid-json"
	if err := dashboardUC.ValidateConfigFiltros(&invalidJSON); err == nil {
		t.Fatal("expected invalid dashboard config JSON error")
	}

	logUC := NewLogAuditoriaUseCase(&testutils.MockLogAuditoriaRepository{}, &testutils.MockUsuarioAdministradorRepository{}, &testutils.MockEmpresaRepository{})
	if err := logUC.ValidateLogEntry(&entity.LogAuditoria{IDUserAdmin: 0, AcaoRealizada: "", Detalhes: "", EnderecoIP: ""}); err == nil {
		t.Fatal("expected log validation error")
	}
	if err := logUC.CleanOldLogs(context.Background(), 10); err == nil {
		t.Fatal("expected retention lower bound error")
	}

	pesquisaUC := NewPesquisaUseCase(&testutils.MockPesquisaRepository{}, &testutils.MockEmpresaRepository{}, &testutils.MockSetorRepository{}, &testutils.MockDashboardRepository{}, &testutils.MockLogAuditoriaRepository{})
	if _, err := pesquisaUC.ListByStatus(context.Background(), 1, "INVALIDO"); err == nil {
		t.Fatal("expected invalid status filter error")
	}
}

func validCNPJFromBase(base12 string) string {
	if len(base12) != 12 {
		panic(fmt.Sprintf("base must have 12 digits, got %d", len(base12)))
	}
	digits := make([]int, 0, 14)
	for _, c := range base12 {
		digits = append(digits, int(c-'0'))
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += digits[i] * weights1[i]
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	digits = append(digits, d1)

	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += digits[i] * weights2[i]
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}

	return fmt.Sprintf("%s%d%d", base12, d1, d2)
}
