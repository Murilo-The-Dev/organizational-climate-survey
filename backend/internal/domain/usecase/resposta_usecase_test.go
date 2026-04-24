package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/testutils"
	"organizational-climate-survey/backend/pkg/crypto"
)

func buildUseCasesForRespostaTests(subRepo *testutils.MockSubmissaoPesquisaRepository, perguntaRepo *testutils.MockPerguntaRepository, pesquisaRepo *testutils.MockPesquisaRepository, respostaRepo *testutils.MockRespostaRepository) (*SubmissaoPesquisaUseCase, *RespostaUseCase) {
	cryptoSvc := crypto.NewDefaultCryptoService()
	subUC := NewSubmissaoPesquisaUseCase(subRepo, pesquisaRepo, cryptoSvc, "salt")
	respUC := NewRespostaUseCase(respostaRepo, perguntaRepo, pesquisaRepo, subUC)
	return subUC, respUC
}

func TestCreateBatch_SuccessAnonymousFlow(t *testing.T) {
	subRepo := &testutils.MockSubmissaoPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}

	subRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.SubmissaoPesquisa, error) {
		return &entity.SubmissaoPesquisa{ID: 42, IDPesquisa: 7, Status: "pendente", DataExpiracao: time.Now().Add(time.Hour)}, nil
	}
	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Status: "Ativa"}, nil
	}
	perguntaRepo.ListByPesquisaFunc = func(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
		return []*entity.Pergunta{{ID: 5, IDPesquisa: pesquisaID}}, nil
	}
	perguntaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pergunta, error) {
		return &entity.Pergunta{ID: id, IDPesquisa: 7, TipoPergunta: "SimNao", TextoPergunta: "Pergunta"}, nil
	}

	created := false
	respostaRepo.CreateBatchFunc = func(ctx context.Context, respostas []*entity.Resposta) error {
		created = true
		if len(respostas) != 1 {
			t.Fatalf("expected 1 resposta, got %d", len(respostas))
		}
		if respostas[0].IDSubmissao != 42 {
			t.Fatalf("expected IDSubmissao=42, got %d", respostas[0].IDSubmissao)
		}
		return nil
	}

	completed := false
	subRepo.MarkAsCompletedFunc = func(ctx context.Context, id int) error {
		completed = true
		if id != 42 {
			t.Fatalf("expected completion for submission 42, got %d", id)
		}
		return nil
	}

	_, respUC := buildUseCasesForRespostaTests(subRepo, perguntaRepo, pesquisaRepo, respostaRepo)
	err := respUC.CreateBatch(context.Background(), []*entity.Resposta{{IDPergunta: 5, ValorResposta: "Sim"}}, "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected batch creation")
	}
	if !completed {
		t.Fatal("expected submission completion")
	}
}

func TestCreateBatch_RejectDuplicateSubmission(t *testing.T) {
	subRepo := &testutils.MockSubmissaoPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}

	subRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.SubmissaoPesquisa, error) {
		return nil, errors.New("token já utilizado")
	}

	_, respUC := buildUseCasesForRespostaTests(subRepo, perguntaRepo, pesquisaRepo, respostaRepo)
	err := respUC.CreateBatch(context.Background(), []*entity.Resposta{{IDPergunta: 5, ValorResposta: "Sim"}}, "used-token")
	if err == nil {
		t.Fatal("expected duplicate submission error")
	}
	if !strings.Contains(err.Error(), "token inválido") {
		t.Fatalf("expected token error, got: %v", err)
	}
}

func TestCreateBatch_RejectExpiredOrInvalidToken(t *testing.T) {
	subRepo := &testutils.MockSubmissaoPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}

	subRepo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.SubmissaoPesquisa, error) {
		return nil, errors.New("token expirado")
	}

	_, respUC := buildUseCasesForRespostaTests(subRepo, perguntaRepo, pesquisaRepo, respostaRepo)
	err := respUC.CreateBatch(context.Background(), []*entity.Resposta{{IDPergunta: 5, ValorResposta: "Sim"}}, "expired-token")
	if err == nil {
		t.Fatal("expected expired token error")
	}
}

func TestDeleteByPesquisa_LGPDFlow(t *testing.T) {
	subRepo := &testutils.MockSubmissaoPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}

	deleted := false
	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Status: "Concluída"}, nil
	}
	respostaRepo.CountByPesquisaFunc = func(ctx context.Context, pesquisaID int) (int, error) {
		return 12, nil
	}
	respostaRepo.DeleteByPesquisaFunc = func(ctx context.Context, pesquisaID int) error {
		deleted = true
		return nil
	}

	_, respUC := buildUseCasesForRespostaTests(subRepo, perguntaRepo, pesquisaRepo, respostaRepo)
	err := respUC.DeleteByPesquisa(context.Background(), 7, 1, "LGPD cleanup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Fatal("expected delete by pesquisa to be called")
	}
}

func TestDeleteByPesquisa_ActiveSurveyRejected(t *testing.T) {
	subRepo := &testutils.MockSubmissaoPesquisaRepository{}
	perguntaRepo := &testutils.MockPerguntaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	respostaRepo := &testutils.MockRespostaRepository{}

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Status: "Ativa"}, nil
	}

	_, respUC := buildUseCasesForRespostaTests(subRepo, perguntaRepo, pesquisaRepo, respostaRepo)
	err := respUC.DeleteByPesquisa(context.Background(), 7, 1, "LGPD cleanup")
	if err == nil {
		t.Fatal("expected error for active survey")
	}
}
