package postgres

import (
    "context"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
)

type RespostaRepository struct {
    db *DB
}

func NewRespostaRepository(db *DB) *RespostaRepository {
    return &RespostaRepository{db: db}
}

var _ repository.RespostaRepository = (*RespostaRepository)(nil)

func (r *RespostaRepository) CreateBatch(ctx context.Context, respostas []*entity.Resposta) error {
    if len(respostas) == 0 {
        return nil
    }
    
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("erro ao iniciar transação: %v", err)
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO resposta (id_pergunta, id_pesquisa, valor_resposta, data_resposta)
        VALUES ($1, $2, $3, $4)
    `)
    if err != nil {
        return fmt.Errorf("erro ao preparar statement: %v", err)
    }
    defer stmt.Close()
    
    for _, resposta := range respostas {
        _, err := stmt.ExecContext(ctx,
            resposta.IDPergunta,
            resposta.IDPesquisa,
            resposta.ValorResposta,
            resposta.DataResposta,
        )
        if err != nil {
            return fmt.Errorf("erro ao inserir resposta: %v", err)
        }
    }
    
    return tx.Commit()
}

func (r *RespostaRepository) CountByPesquisa(ctx context.Context, pesquisaID int) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM resposta WHERE id_pesquisa = $1`
    
    err := r.db.QueryRowContext(ctx, query, pesquisaID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("erro ao contar respostas da pesquisa: %v", err)
    }
    
    return count, nil
}

func (r *RespostaRepository) CountByPergunta(ctx context.Context, perguntaID int) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM resposta WHERE id_pergunta = $1`
    
    err := r.db.QueryRowContext(ctx, query, perguntaID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("erro ao contar respostas da pergunta: %v", err)
    }
    
    return count, nil
}

func (r *RespostaRepository) GetAggregatedByPergunta(ctx context.Context, perguntaID int) (map[string]int, error) {
    query := `
        SELECT valor_resposta, COUNT(*) as quantidade
        FROM resposta 
        WHERE id_pergunta = $1 
        GROUP BY valor_resposta
        ORDER BY quantidade DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, perguntaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar dados agregados: %v", err)
    }
    defer rows.Close()
    
    result := make(map[string]int)
    
    for rows.Next() {
        var valor string
        var quantidade int
        
        err := rows.Scan(&valor, &quantidade)
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear resultado: %v", err)
        }
        
        result[valor] = quantidade
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erro durante iteração: %v", err)
    }
    
    return result, nil
}

func (r *RespostaRepository) GetAggregatedByPesquisa(ctx context.Context, pesquisaID int) (map[int]map[string]int, error) {
    query := `
        SELECT id_pergunta, valor_resposta, COUNT(*) as quantidade
        FROM resposta 
        WHERE id_pesquisa = $1 
        GROUP BY id_pergunta, valor_resposta
        ORDER BY id_pergunta, quantidade DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, pesquisaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar dados agregados por pesquisa: %v", err)
    }
    defer rows.Close()
    
    result := make(map[int]map[string]int)
    
    for rows.Next() {
        var perguntaID int
        var valor string
        var quantidade int
        
        err := rows.Scan(&perguntaID, &valor, &quantidade)
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear resultado: %v", err)
        }
        
        if result[perguntaID] == nil {
            result[perguntaID] = make(map[string]int)
        }
        
        result[perguntaID][valor] = quantidade
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erro durante iteração: %v", err)
    }
    
    return result, nil
}

func (r *RespostaRepository) GetResponsesByDateRange(ctx context.Context, pesquisaID int, startDate, endDate string) ([]*entity.Resposta, error) {
    query := `
        SELECT id_resposta, id_pergunta, id_pesquisa, valor_resposta, data_resposta, data_submissao
        FROM resposta 
        WHERE id_pesquisa = $1 
          AND data_submissao BETWEEN $2 AND $3
        ORDER BY data_submissao
    `
    
    rows, err := r.db.QueryContext(ctx, query, pesquisaID, startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar respostas por período: %v", err)
    }
    defer rows.Close()
    
    var respostas []*entity.Resposta
    
    for rows.Next() {
        resposta := &entity.Resposta{}
        err := rows.Scan(
            &resposta.ID,
            &resposta.IDPergunta,
            &resposta.IDPesquisa,
            &resposta.ValorResposta,
            &resposta.DataResposta,
            &resposta.DataSubmissao,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear resposta: %v", err)
        }
        
        respostas = append(respostas, resposta)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erro durante iteração: %v", err)
    }
    
    return respostas, nil
}

func (r *RespostaRepository) DeleteByPesquisa(ctx context.Context, pesquisaID int) error {
    query := `DELETE FROM resposta WHERE id_pesquisa = $1`
    
    result, err := r.db.ExecContext(ctx, query, pesquisaID)
    if err != nil {
        return fmt.Errorf("erro ao deletar respostas: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    // Log para auditoria - pode adicionar se necessário
    _ = rowsAffected
    
    return nil
}