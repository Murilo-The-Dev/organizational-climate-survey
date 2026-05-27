package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/testutils"
	"organizational-climate-survey/backend/pkg/logger"

	"github.com/gorilla/mux"
)

func TestGenerateSurveyLink_QRCodeIncludedInResponse(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("QRCODE_OUTPUT_DIR", tmpDir)
	t.Setenv("QRCODE_PUBLIC_BASE_PATH", "/uploads/qrcodes")
	t.Setenv("FRONTEND_URL", "https://frontend.example")

	pesquisaRepo := &testutils.MockPesquisaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	setorRepo := &testutils.MockSetorRepository{}
	dashboardRepo := &testutils.MockDashboardRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{
			ID:          id,
			IDEmpresa:   1,
			IDSetor:     0,
			IDUserAdmin: 1,
			Titulo:      "Pesquisa QR",
			Status:      "Rascunho",
			LinkAcesso:  "abc123",
		}, nil
	}
	updated := false
	pesquisaRepo.UpdateFunc = func(ctx context.Context, p *entity.Pesquisa) error {
		updated = true
		if p.QRCodePath == "" {
			t.Fatal("expected qrcode path to be set")
		}
		return nil
	}

	uc := usecase.NewPesquisaUseCase(pesquisaRepo, empresaRepo, setorRepo, dashboardRepo, logRepo)
	h := NewPesquisaHandler(uc, logger.NoopLogger{})

	router := mux.NewRouter()
	h.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/pesquisas/1/qrcode", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_admin_id", 1))
	req.Host = "api.example"
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", rec.Code, rec.Body.String())
	}
	if !updated {
		t.Fatal("expected pesquisa update to persist qrcode path")
	}

	var apiResp response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &apiResp); err != nil {
		t.Fatalf("failed to decode API response: %v", err)
	}
	if !apiResp.Success {
		t.Fatalf("expected success response, got %+v", apiResp)
	}

	data, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %T", apiResp.Data)
	}

	base64Value, _ := data["qr_code_base64"].(string)
	if strings.TrimSpace(base64Value) == "" {
		t.Fatal("expected non-empty qr_code_base64")
	}

	surveyURL, _ := data["survey_url"].(string)
	if !strings.Contains(surveyURL, "/pesquisas/link/abc123") {
		t.Fatalf("expected survey url to contain link token, got %s", surveyURL)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "abc123.png")); err != nil {
		t.Fatalf("expected qr code file to be generated: %v", err)
	}
}

func TestQRCodeContent_ContainsSurveyURL(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("QRCODE_OUTPUT_DIR", tmpDir)
	t.Setenv("QRCODE_PUBLIC_BASE_PATH", "/uploads/qrcodes")
	t.Setenv("FRONTEND_URL", "https://frontend.example")

	pesquisaRepo := &testutils.MockPesquisaRepository{}
	empresaRepo := &testutils.MockEmpresaRepository{}
	setorRepo := &testutils.MockSetorRepository{}
	dashboardRepo := &testutils.MockDashboardRepository{}
	logRepo := &testutils.MockLogAuditoriaRepository{}

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{
			ID:          id,
			IDEmpresa:   1,
			IDSetor:     0,
			IDUserAdmin: 1,
			Titulo:      "Pesquisa QR",
			Status:      "Rascunho",
			LinkAcesso:  "xyz789",
		}, nil
	}
	pesquisaRepo.UpdateFunc = func(ctx context.Context, p *entity.Pesquisa) error { return nil }

	uc := usecase.NewPesquisaUseCase(pesquisaRepo, empresaRepo, setorRepo, dashboardRepo, logRepo)
	h := NewPesquisaHandler(uc, logger.NoopLogger{})

	router := mux.NewRouter()
	h.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/pesquisas/1/qrcode", nil)
	req.Host = "api.example"
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", rec.Code, rec.Body.String())
	}

	var apiResp response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &apiResp); err != nil {
		t.Fatalf("failed to decode API response: %v", err)
	}
	data, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %T", apiResp.Data)
	}

	qrPayload, _ := data["qr_payload"].(string)
	if !strings.Contains(qrPayload, "https://frontend.example/pesquisas/link/xyz789") {
		t.Fatalf("expected qr payload to contain survey URL, got %s", qrPayload)
	}
}
