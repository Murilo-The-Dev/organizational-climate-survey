package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
    "time"
)

type RespostaRepository struct {
    db *DB
}

// NewRespostaRepository cria uma nova instância do repositório de resposta
func NewRespostaRepository(db *DB) *RespostaRepository {
    return &RespostaRepository{db: db}
}

// Verifica se implementa a interface
var _ repository.RespostaRepository = (*RespostaRepository)(nil)

// CreateBatch insere múltiplas respostas de uma vez (performance otimizada)
func (r *RespostaRepository) CreateBatch(ctx context.Context, respostas []*entity.Resposta) error {
    if len(respostas) == 0 {
        return nil
    }
    
    // Usar transação para garantir consistência
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("erro ao iniciar transação: %v", err)
    }
    defer tx.Rollback()
    
    // Preparar statement
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO resposta (id_pergunta, data_submissao, valor_resposta)
        VALUES ($1, $2, $3)
        RETURNING id_resposta
    `)
    if err != nil {
        return fmt.Errorf("erro ao preparar statement: %v", err)
    }
    defer stmt.Close()
    
    // Inserir cada resposta
    for _, resposta := range respostas {
        err := stmt.QueryRowContext(ctx,
            resposta.IDPergunta,
            resposta.DataSubmissao,
            resposta.ValorResposta,
        ).Scan(&resposta.ID)
        
        if err != nil {
            return fmt.Errorf("erro ao inserir resposta: %v", err)
        }
    }
    
    // Confirmar transação
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("erro ao confirmar transação: %v", err)
    }
    
    return nil
}

// CountByPesquisa conta o total de respostas de uma pesquisa
func (r *RespostaRepository) CountByPesquisa(ctx context.Context, pesquisaID int) (int, error) {
    var count int
    query := `
        SELECT COUNT(*)
        FROM resposta r
        INNER JOIN pergunta p ON r.id_pergunta = p.id_pergunta
        WHERE p.id_pesquisa = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, pesquisaID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("erro ao contar respostas da pesquisa: %v", err)
    }
    
    return count, nil
}

// CountByPergunta conta o total de respostas de uma pergunta específica
func (r *RespostaRepository) CountByPergunta(ctx context.Context, perguntaID int) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM resposta WHERE id_pergunta = $1`
    
    err := r.db.QueryRowContext(ctx, query, perguntaID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("erro ao contar respostas da pergunta: %v", err)
    }
    
    return count, nil
}

// GetAggregatedByPergunta retorna dados agregados de uma pergunta (para análise)
func (r *RespostaRepository) GetAggregatedByPergunta(ctx context.Context, perguntaID int) (map[string]int, error) {
    query := `
        SELECT valor_resposta, COUNT(*) as total
        FROM resposta
        WHERE id_pergunta = $1
        GROUP BY valor_resposta
        ORDER BY total DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, perguntaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao obter dados agregados: %v", err)
    }
    defer rows.Close()
    
    result := make(map[string]int)
    
    for rows.Next() {
        var valor string
        var total int
        
        err := rows.Scan(&valor, &total)
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear dados agregados: %v", err)
        }
        
        result[valor] = total
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar dados agregados: %v", err)
    }
    
    return result, nil
}

// GetAggregatedByPesquisa retorna dados agregados de todas as perguntas de uma pesquisa
func (r *RespostaRepository) GetAggregatedByPesquisa(ctx context.Context, pesquisaID int) (map[int]map[string]int, error) {
    query := `
        SELECT p.id_pergunta, r.valor_resposta, COUNT(*) as total
        FROM resposta r
        INNER JOIN pergunta p ON r.id_pergunta = p.id_pergunta
        WHERE p.id_pesquisa = $1
        GROUP BY p.id_pergunta, r.valor_resposta
        ORDER BY p.ordem_exibicao, total DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, pesquisaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao obter dados agregados da pesquisa: %v", err)
    }
    defer rows.Close()
    
    result := make(map[int]map[string]int)
    
    for rows.Next() {
        var perguntaID int
        var valor string
        var total int
        
        err := rows.Scan(&perguntaID, &valor, &total)
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear dados agregados: %v", err)
        }
        
        if result[perguntaID] == nil {
            result[perguntaID] = make(map[string]int)
        }
        
        result[perguntaID][valor] = total
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar dados agregados: %v", err)
    }
    
    return result, nil
}

// GetResponsesByDateRange retorna respostas dentro de um período (para análise temporal)
func (r *RespostaRepository) GetResponsesByDateRange(ctx context.Context, pesquisaID int, startDate, endDate string) ([]*entity.Resposta, error) {
    query := `
        SELECT r.id_resposta, r.id_pergunta, r.data_submissao, r.valor_resposta
        FROM resposta r
        INNER JOIN pergunta p ON r.id_pergunta = p.id_pergunta
        WHERE p.id_pesquisa = $1 
        AND r.data_submissao >= $2 
        AND r.data_submissao <= $3
        ORDER BY r.data_submissao DESC
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
            &resposta.DataSubmissao,
            &resposta.ValorResposta,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear resposta: %v", err)
        }
        respostas = append(respostas, resposta)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar respostas: %v", err)
    }
    
    return respostas, nil
}

// DeleteByPesquisa remove todas as respostas de uma pesquisa (para limpeza)
func (r *RespostaRepository) DeleteByPesquisa(ctx context.Context, pesquisaID int) error {
    query := `
        DELETE FROM resposta 
        WHERE id_pergunta IN (
            SELECT id_pergunta FROM pergunta WHERE id_pesquisa = $1
        )
    `
    
    result, err := r.db.ExecContext(ctx, query, pesquisaID)
    if err != nil {
        return fmt.Errorf("erro ao deletar respostas da pesquisa: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    // Log informativo sobre quantas respostas foram deletadas
    if rowsAffected > 0 {
        fmt.Printf("Deletadas %d respostas da pesquisa ID %d\n", rowsAffected, pesquisaID)
    }
    
    return nil
}

// GetStatsByPesquisa retorna estatísticas básicas de uma pesquisa
func (r *RespostaRepository) GetStatsByPesquisa(ctx context.Context, pesquisaID int) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Total de respostas
    totalRespostas, err := r.CountByPesquisa(ctx, pesquisaID)
    if err != nil {
        return nil, err
    }
    stats["total_respostas"] = totalRespostas
    
    // Primeira e última resposta
    query := `
        SELECT MIN(r.data_submissao) as primeira_resposta,
               MAX(r.data_submissao) as ultima_resposta
        FROM resposta r
        INNER JOIN pergunta p ON r.id_pergunta = p.id_pergunta
        WHERE p.id_pesquisa = $1
    `
    
    var primeiraResposta, ultimaResposta sql.NullTime
    err = r.db.QueryRowContext(ctx, query, pesquisaID).Scan(&primeiraResposta, &ultimaResposta)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("erro ao buscar estatísticas: %v", err)
    }
    
    if primeiraResposta.Valid {
        stats["primeira_resposta"] = primeiraResposta.Time
    }
    if ultimaResposta.Valid {
        stats["ultima_resposta"] = ultimaResposta.Time
    }
    
    // Respostas por dia (útil para análise temporal)
    queryDiaria := `
        SELECT DATE(r.data_submissao) as data, COUNT(*) as total
        FROM resposta r
        INNER JOIN pergunta p ON r.id_pergunta = p.id_pergunta
        WHERE p.id_pesquisa = $1
        GROUP BY DATE(r.data_submissao)
        ORDER BY data DESC
        LIMIT 7
    `
    
    rows, err := r.db.QueryContext(ctx, queryDiaria, pesquisaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar dados diários: %v", err)
    }
    defer rows.Close()
    
    respostasPorDia := make(map[string]int)
    for rows.Next() {
        var data time.Time
        var total int
        err := rows.Scan(&data, &total)
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear dados diários: %v", err)
        }
        respostasPorDia[data.Format("2006-01-02")] = total
    }
    
    stats["respostas_por_dia"] = respostasPorDia
    
    return stats, nil
}