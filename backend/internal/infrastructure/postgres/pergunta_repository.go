package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
)

type PerguntaRepository struct {
    db *DB
}

func NewPerguntaRepository(db *DB) *PerguntaRepository {
    return &PerguntaRepository{db: db}
}

var _ repository.PerguntaRepository = (*PerguntaRepository)(nil)

func (r *PerguntaRepository) Create(ctx context.Context, pergunta *entity.Pergunta) error {
    query := `
        INSERT INTO pergunta (id_pesquisa, texto_pergunta, tipo_pergunta, ordem_exibicao, opcoes_resposta)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id_pergunta
    `
    
    err := r.db.QueryRowContext(ctx, query,
        pergunta.IDPesquisa,
        pergunta.TextoPergunta,
        pergunta.TipoPergunta,
        pergunta.OrdemExibicao,
        pergunta.OpcoesResposta,
    ).Scan(&pergunta.ID)
    
    if err != nil {
        return fmt.Errorf("erro ao criar pergunta: %v", err)
    }
    
    return nil
}

func (r *PerguntaRepository) CreateBatch(ctx context.Context, perguntas []*entity.Pergunta) error {
    if len(perguntas) == 0 {
        return nil
    }
    
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("erro ao iniciar transação: %v", err)
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO pergunta (id_pesquisa, texto_pergunta, tipo_pergunta, ordem_exibicao, opcoes_resposta)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id_pergunta
    `)
    if err != nil {
        return fmt.Errorf("erro ao preparar statement: %v", err)
    }
    defer stmt.Close()
    
    for _, pergunta := range perguntas {
        err := stmt.QueryRowContext(ctx,
            pergunta.IDPesquisa,
            pergunta.TextoPergunta,
            pergunta.TipoPergunta,
            pergunta.OrdemExibicao,
            pergunta.OpcoesResposta,
        ).Scan(&pergunta.ID)
        
        if err != nil {
            return fmt.Errorf("erro ao inserir pergunta: %v", err)
        }
    }
    
    return tx.Commit()
}

func (r *PerguntaRepository) GetByID(ctx context.Context, id int) (*entity.Pergunta, error) {
    pergunta := &entity.Pergunta{}
    query := `
        SELECT id_pergunta, id_pesquisa, texto_pergunta, tipo_pergunta, ordem_exibicao, opcoes_resposta
        FROM pergunta
        WHERE id_pergunta = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &pergunta.ID,
        &pergunta.IDPesquisa,
        &pergunta.TextoPergunta,
        &pergunta.TipoPergunta,
        &pergunta.OrdemExibicao,
        &pergunta.OpcoesResposta,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("pergunta com ID %d não encontrada", id)
        }
        return nil, fmt.Errorf("erro ao buscar pergunta: %v", err)
    }
    
    return pergunta, nil
}

// GetByPesquisaID - método necessário para o DashboardUseCase
func (r *PerguntaRepository) GetByPesquisaID(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
    query := `
        SELECT id_pergunta, id_pesquisa, texto_pergunta, tipo_pergunta, ordem_exibicao, opcoes_resposta
        FROM pergunta
        WHERE id_pesquisa = $1
        ORDER BY ordem_exibicao
    `
    
    rows, err := r.db.QueryContext(ctx, query, pesquisaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao listar perguntas: %v", err)
    }
    defer rows.Close()
    
    var perguntas []*entity.Pergunta
    
    for rows.Next() {
        pergunta := &entity.Pergunta{}
        err := rows.Scan(
            &pergunta.ID,
            &pergunta.IDPesquisa,
            &pergunta.TextoPergunta,
            &pergunta.TipoPergunta,
            &pergunta.OrdemExibicao,
            &pergunta.OpcoesResposta,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear pergunta: %v", err)
        }
        perguntas = append(perguntas, pergunta)
    }
    
    return perguntas, nil
}

func (r *PerguntaRepository) ListByPesquisa(ctx context.Context, pesquisaID int) ([]*entity.Pergunta, error) {
    // Este método já existe - pode usar GetByPesquisaID ou manter ambos
    return r.GetByPesquisaID(ctx, pesquisaID)
}

func (r *PerguntaRepository) Update(ctx context.Context, pergunta *entity.Pergunta) error {
    query := `
        UPDATE pergunta 
        SET texto_pergunta = $2, tipo_pergunta = $3, ordem_exibicao = $4, opcoes_resposta = $5
        WHERE id_pergunta = $1
    `
    
    result, err := r.db.ExecContext(ctx, query,
        pergunta.ID,
        pergunta.TextoPergunta,
        pergunta.TipoPergunta,
        pergunta.OrdemExibicao,
        pergunta.OpcoesResposta,
    )
    
    if err != nil {
        return fmt.Errorf("erro ao atualizar pergunta: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("pergunta com ID %d não encontrada", pergunta.ID)
    }
    
    return nil
}

func (r *PerguntaRepository) UpdateOrdem(ctx context.Context, perguntaID int, novaOrdem int) error {
    query := `
        UPDATE pergunta 
        SET ordem_exibicao = $2
        WHERE id_pergunta = $1
    `
    
    result, err := r.db.ExecContext(ctx, query, perguntaID, novaOrdem)
    if err != nil {
        return fmt.Errorf("erro ao atualizar ordem da pergunta: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("pergunta com ID %d não encontrada", perguntaID)
    }
    
    return nil
}

func (r *PerguntaRepository) Delete(ctx context.Context, id int) error {
    // Verifica se há respostas vinculadas
    var count int
    checkQuery := `SELECT COUNT(*) FROM resposta WHERE id_pergunta = $1`
    err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&count)
    if err != nil {
        return fmt.Errorf("erro ao verificar dependências: %v", err)
    }
    
    if count > 0 {
        return fmt.Errorf("não é possível deletar pergunta: possui %d respostas vinculadas", count)
    }
    
    query := `DELETE FROM pergunta WHERE id_pergunta = $1`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("erro ao deletar pergunta: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("pergunta com ID %d não encontrada", id)
    }
    
    return nil
}