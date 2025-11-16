// Package postgres implementa o repositório de Resposta usando PostgreSQL.
// Fornece operações para gerenciamento e análise de respostas das pesquisas.
package postgres

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/logger"
)

// RespostaRepository implementa a interface repository.RespostaRepository
type RespostaRepository struct {
	db     *DB           // Conexão com o banco de dados
	logger logger.Logger // Logger para operações do repositório
}

// NewRespostaRepository cria uma nova instância do repositório
func NewRespostaRepository(db *DB) *RespostaRepository {
	return &RespostaRepository{
		db:     db,
		logger: db.logger,
	}
}

var _ repository.RespostaRepository = (*RespostaRepository)(nil)

// CreateBatch insere múltiplas respostas em uma única transação
// Otimiza a inserção em lote usando prepared statement
func (r *RespostaRepository) CreateBatch(ctx context.Context, respostas []*entity.Resposta) error {
	if len(respostas) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("erro ao iniciar transação batch respostas: %v", err)
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO resposta (id_pergunta, id_pesquisa, valor_resposta, data_resposta)
        VALUES ($1, $2, $3, $4)
    `)
	if err != nil {
		r.logger.Error("erro ao preparar statement batch respostas: %v", err)
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
			r.logger.Error("erro ao inserir resposta batch: %v", err)
			return fmt.Errorf("erro ao inserir resposta: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("erro ao commit batch respostas: %v", err)
		return fmt.Errorf("erro ao commit: %v", err)
	}

	return nil
}

// CountByPesquisa conta o total de respostas de uma pesquisa
func (r *RespostaRepository) CountByPesquisa(ctx context.Context, pesquisaID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM resposta WHERE id_pesquisa = $1`

	err := r.db.QueryRowContext(ctx, query, pesquisaID).Scan(&count)
	if err != nil {
		r.logger.Error("erro ao contar respostas pesquisa ID=%d: %v", pesquisaID, err)
		return 0, fmt.Errorf("erro ao contar respostas da pesquisa: %v", err)
	}

	return count, nil
}

// CountByPergunta conta o total de respostas para uma pergunta específica
func (r *RespostaRepository) CountByPergunta(ctx context.Context, perguntaID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM resposta WHERE id_pergunta = $1`

	err := r.db.QueryRowContext(ctx, query, perguntaID).Scan(&count)
	if err != nil {
		r.logger.Error("erro ao contar respostas pergunta ID=%d: %v", perguntaID, err)
		return 0, fmt.Errorf("erro ao contar respostas da pergunta: %v", err)
	}

	return count, nil
}

// GetAggregatedByPergunta retorna contagem agrupada de respostas por valor
// Útil para análises e dashboards
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
		r.logger.Error("erro ao buscar agregados pergunta ID=%d: %v", perguntaID, err)
		return nil, fmt.Errorf("erro ao buscar dados agregados: %v", err)
	}
	defer rows.Close()

	result := make(map[string]int)

	for rows.Next() {
		var valor string
		var quantidade int

		err := rows.Scan(&valor, &quantidade)
		if err != nil {
			r.logger.Error("erro ao escanear agregado: %v", err)
			return nil, fmt.Errorf("erro ao escanear resultado: %v", err)
		}

		result[valor] = quantidade
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("erro ao iterar agregados: %v", err)
		return nil, fmt.Errorf("erro durante iteração: %v", err)
	}

	return result, nil
}

// GetAggregatedByPesquisa retorna contagem agrupada de todas as respostas da pesquisa
// Agrupadas por pergunta e valor da resposta
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
		r.logger.Error("erro ao buscar agregados pesquisa ID=%d: %v", pesquisaID, err)
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
			r.logger.Error("erro ao escanear agregado pesquisa: %v", err)
			return nil, fmt.Errorf("erro ao escanear resultado: %v", err)
		}

		if result[perguntaID] == nil {
			result[perguntaID] = make(map[string]int)
		}

		result[perguntaID][valor] = quantidade
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("erro ao iterar agregados pesquisa: %v", err)
		return nil, fmt.Errorf("erro durante iteração: %v", err)
	}

	return result, nil
}

// GetResponsesByDateRange busca respostas dentro de um intervalo de datas
// Usado para análises temporais e relatórios
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
		r.logger.Error("erro ao buscar respostas por período pesquisa ID=%d: %v", pesquisaID, err)
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
			r.logger.Error("erro ao escanear resposta: %v", err)
			return nil, fmt.Errorf("erro ao escanear resposta: %v", err)
		}

		respostas = append(respostas, resposta)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("erro ao iterar respostas: %v", err)
		return nil, fmt.Errorf("erro durante iteração: %v", err)
	}

	return respostas, nil
}

// DeleteByPesquisa remove todas as respostas de uma pesquisa
// Útil para limpeza de dados ou remoção de pesquisas
func (r *RespostaRepository) DeleteByPesquisa(ctx context.Context, pesquisaID int) error {
	query := `DELETE FROM resposta WHERE id_pesquisa = $1`

	result, err := r.db.ExecContext(ctx, query, pesquisaID)
	if err != nil {
		r.logger.Error("erro ao deletar respostas pesquisa ID=%d: %v", pesquisaID, err)
		return fmt.Errorf("erro ao deletar respostas: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	r.logger.Info("respostas deletadas pesquisa ID=%d count=%d", pesquisaID, rowsAffected)

	return nil
}
