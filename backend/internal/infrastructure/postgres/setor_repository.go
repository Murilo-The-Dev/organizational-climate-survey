// setor_repository.go
package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
    "organizational-climate-survey/backend/pkg/logger"
)

type SetorRepository struct {
    db     *DB
    logger logger.Logger
}

func NewSetorRepository(db *DB) *SetorRepository {
    return &SetorRepository{
        db:     db,
        logger: db.logger,
    }
}

var _ repository.SetorRepository = (*SetorRepository)(nil)

func (r *SetorRepository) Create(ctx context.Context, setor *entity.Setor) error {
    query := `
        INSERT INTO setor (id_empresa, nome_setor, descricao)
        VALUES ($1, $2, $3)
        RETURNING id_setor
    `
    
    err := r.db.QueryRowContext(ctx, query,
        setor.IDEmpresa,
        setor.NomeSetor,
        setor.Descricao,
    ).Scan(&setor.ID)
    
    if err != nil {
        r.logger.Error("erro ao criar setor nome=%s: %v", setor.NomeSetor, err)
        return fmt.Errorf("erro ao criar setor: %v", err)
    }
    
    return nil
}

func (r *SetorRepository) GetByID(ctx context.Context, id int) (*entity.Setor, error) {
    setor := &entity.Setor{}
    query := `
        SELECT id_setor, id_empresa, nome_setor, descricao
        FROM setor
        WHERE id_setor = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &setor.ID,
        &setor.IDEmpresa,
        &setor.NomeSetor,
        &setor.Descricao,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("setor com ID %d não encontrado", id)
        }
        r.logger.Error("erro ao buscar setor ID=%d: %v", id, err)
        return nil, fmt.Errorf("erro ao buscar setor: %v", err)
    }
    
    return setor, nil
}

func (r *SetorRepository) GetByNome(ctx context.Context, empresaID int, nome string) (*entity.Setor, error) {
    setor := &entity.Setor{}
    query := `
        SELECT id_setor, id_empresa, nome_setor, descricao
        FROM setor
        WHERE id_empresa = $1 AND nome_setor = $2
    `
    
    err := r.db.QueryRowContext(ctx, query, empresaID, nome).Scan(
        &setor.ID,
        &setor.IDEmpresa,
        &setor.NomeSetor,
        &setor.Descricao,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("setor %s não encontrado na empresa ID %d", nome, empresaID)
        }
        r.logger.Error("erro ao buscar setor nome=%s empresa ID=%d: %v", nome, empresaID, err)
        return nil, fmt.Errorf("erro ao buscar setor: %v", err)
    }
    
    return setor, nil
}

func (r *SetorRepository) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Setor, error) {
    query := `
        SELECT id_setor, id_empresa, nome_setor, descricao
        FROM setor
        WHERE id_empresa = $1
        ORDER BY nome_setor
    `
    
    rows, err := r.db.QueryContext(ctx, query, empresaID)
    if err != nil {
        r.logger.Error("erro ao listar setores empresa ID=%d: %v", empresaID, err)
        return nil, fmt.Errorf("erro ao listar setores: %v", err)
    }
    defer rows.Close()
    
    var setores []*entity.Setor
    
    for rows.Next() {
        setor := &entity.Setor{}
        err := rows.Scan(
            &setor.ID,
            &setor.IDEmpresa,
            &setor.NomeSetor,
            &setor.Descricao,
        )
        if err != nil {
            r.logger.Error("erro ao escanear setor: %v", err)
            return nil, fmt.Errorf("erro ao escanear setor: %v", err)
        }
        setores = append(setores, setor)
    }
    
    return setores, nil
}

func (r *SetorRepository) Update(ctx context.Context, setor *entity.Setor) error {
    query := `
        UPDATE setor 
        SET nome_setor = $2, descricao = $3
        WHERE id_setor = $1
    `
    
    result, err := r.db.ExecContext(ctx, query,
        setor.ID,
        setor.NomeSetor,
        setor.Descricao,
    )
    
    if err != nil {
        r.logger.Error("erro ao atualizar setor ID=%d: %v", setor.ID, err)
        return fmt.Errorf("erro ao atualizar setor: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("setor com ID %d não encontrado", setor.ID)
    }
    
    return nil
}

func (r *SetorRepository) Delete(ctx context.Context, id int) error {
    var count int
    checkQuery := `SELECT COUNT(*) FROM pesquisa WHERE id_setor = $1`
    err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&count)
    if err != nil {
        r.logger.Error("erro ao verificar dependências setor ID=%d: %v", id, err)
        return fmt.Errorf("erro ao verificar dependências: %v", err)
    }
    
    if count > 0 {
        return fmt.Errorf("não é possível deletar setor: possui %d pesquisas vinculadas", count)
    }
    
    query := `DELETE FROM setor WHERE id_setor = $1`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        r.logger.Error("erro ao deletar setor ID=%d: %v", id, err)
        return fmt.Errorf("erro ao deletar setor: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("setor com ID %d não encontrado", id)
    }
    
    return nil
}