package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/testutils"
	"organizational-climate-survey/backend/pkg/crypto"
)

func TestGenerateAccessToken_Success(t *testing.T) {
	repo := &testutils.MockSubmissaoPesquisaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	cryptoSvc := crypto.NewDefaultCryptoService()

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Status: "Ativa"}, nil
	}
	repo.CountByPesquisaAndIPHashFunc = func(ctx context.Context, pesquisaID int, ipHash string, since time.Time) (int, error) {
		return 0, nil
	}
	created := false
	repo.CreateFunc = func(ctx context.Context, s *entity.SubmissaoPesquisa) error {
		created = true
		if s.TokenAcesso == "" {
			t.Fatal("token should not be empty")
		}
		return nil
	}

	uc := NewSubmissaoPesquisaUseCase(repo, pesquisaRepo, cryptoSvc, "salt")
	token, expiresAt, err := uc.GenerateAccessToken(context.Background(), 1, "127.0.0.1", "fp-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if !created {
		t.Fatal("expected repository create call")
	}
	if !expiresAt.After(time.Now()) {
		t.Fatal("expected future expiration")
	}
}

func TestGenerateAccessToken_RateLimitExceeded(t *testing.T) {
	repo := &testutils.MockSubmissaoPesquisaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	cryptoSvc := crypto.NewDefaultCryptoService()

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Status: "Ativa"}, nil
	}
	repo.CountByPesquisaAndIPHashFunc = func(ctx context.Context, pesquisaID int, ipHash string, since time.Time) (int, error) {
		return 3, nil
	}

	uc := NewSubmissaoPesquisaUseCase(repo, pesquisaRepo, cryptoSvc, "salt")
	_, _, err := uc.GenerateAccessToken(context.Background(), 1, "127.0.0.1", "")
	if err == nil {
		t.Fatal("expected rate limit error")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	repo := &testutils.MockSubmissaoPesquisaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	cryptoSvc := crypto.NewDefaultCryptoService()

	repo.GetByTokenFunc = func(ctx context.Context, token string) (*entity.SubmissaoPesquisa, error) {
		return nil, errors.New("token inválido")
	}

	uc := NewSubmissaoPesquisaUseCase(repo, pesquisaRepo, cryptoSvc, "salt")
	_, err := uc.ValidateToken(context.Background(), "bad-token")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGetSubmissionStats_Success(t *testing.T) {
	repo := &testutils.MockSubmissaoPesquisaRepository{}
	pesquisaRepo := &testutils.MockPesquisaRepository{}
	cryptoSvc := crypto.NewDefaultCryptoService()

	pesquisaRepo.GetByIDFunc = func(ctx context.Context, id int) (*entity.Pesquisa, error) {
		return &entity.Pesquisa{ID: id, Status: "Ativa"}, nil
	}
	repo.CountCompleteByPesquisaFunc = func(ctx context.Context, pesquisaID int) (int, error) {
		return 2, nil
	}
	repo.ListByPesquisaFunc = func(ctx context.Context, pesquisaID int) ([]*entity.SubmissaoPesquisa, error) {
		now := time.Now()
		return []*entity.SubmissaoPesquisa{
			{ID: 1, IDPesquisa: pesquisaID, Status: "completa", DataExpiracao: now.Add(time.Hour)},
			{ID: 2, IDPesquisa: pesquisaID, Status: "pendente", DataExpiracao: now.Add(time.Hour)},
			{ID: 3, IDPesquisa: pesquisaID, Status: "pendente", DataExpiracao: now.Add(-time.Hour)},
		}, nil
	}

	uc := NewSubmissaoPesquisaUseCase(repo, pesquisaRepo, cryptoSvc, "salt")
	stats, err := uc.GetSubmissionStats(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats["total_submissoes"].(int) != 3 {
		t.Fatalf("expected total 3, got %v", stats["total_submissoes"])
	}
}
