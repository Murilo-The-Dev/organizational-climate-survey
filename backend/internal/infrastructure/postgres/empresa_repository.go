// empresa_repository.go
package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
    "organizational-climate-survey/backend/pkg/logger"
)

type EmpresaRepository struct {
    db     *DB
    logger logger.Logger
}

func NewEmpresaRepository(db *DB) *EmpresaRepository {
    return &EmpresaRepository{
        db:     db,
        logger: db.logger,
    }
}

var _ repository.EmpresaRepository = (*EmpresaRepository)(nil)

func (r *EmpresaRepository) Create(ctx context.Context, empresa *entity.Empresa) error {
    query := `
        INSERT INTO empresa (nome_fantasia, razao_social, cnpj, data_cadastro)
        VALUES ($1, $2, $3, $4)
        RETURNING id_empresa
    `
    
    err := r.db.QueryRowContext(ctx, query,
        empresa.NomeFantasia,
        empresa.RazaoSocial,
        empresa.CNPJ,
        empresa.DataCadastro,
    ).Scan(&empresa.ID)
    
    if err != nil {
        r.logger.Error("erro ao criar empresa CNPJ=%s: %v", empresa.CNPJ, err)
        return fmt.Errorf("erro ao criar empresa: %v", err)
    }
    
    return nil
}

func (r *EmpresaRepository) GetByID(ctx context.Context, id int) (*entity.Empresa, error) {
    empresa := &entity.Empresa{}
    query := `
        SELECT id_empresa, nome_fantasia, razao_social, cnpj, data_cadastro
        FROM empresa
        WHERE id_empresa = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &empresa.ID,
        &empresa.NomeFantasia,
        &empresa.RazaoSocial,
        &empresa.CNPJ,
        &empresa.DataCadastro,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("empresa com ID %d não encontrada", id)
        }
        r.logger.Error("erro ao buscar empresa ID=%d: %v", id, err)
        return nil, fmt.Errorf("erro ao buscar empresa: %v", err)
    }
    
    return empresa, nil
}

func (r *EmpresaRepository) GetByCNPJ(ctx context.Context, cnpj string) (*entity.Empresa, error) {
    empresa := &entity.Empresa{}
    query := `
        SELECT id_empresa, nome_fantasia, razao_social, cnpj, data_cadastro
        FROM empresa
        WHERE cnpj = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, cnpj).Scan(
        &empresa.ID,
        &empresa.NomeFantasia,
        &empresa.RazaoSocial,
        &empresa.CNPJ,
        &empresa.DataCadastro,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("empresa com CNPJ %s não encontrada", cnpj)
        }
        r.logger.Error("erro ao buscar empresa CNPJ=%s: %v", cnpj, err)
        return nil, fmt.Errorf("erro ao buscar empresa: %v", err)
    }
    
    return empresa, nil
}

func (r *EmpresaRepository) List(ctx context.Context, limit, offset int) ([]*entity.Empresa, error) {
    query := `
        SELECT id_empresa, nome_fantasia, razao_social, cnpj, data_cadastro
        FROM empresa
        ORDER BY data_cadastro DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        r.logger.Error("erro ao listar empresas: %v", err)
        return nil, fmt.Errorf("erro ao listar empresas: %v", err)
    }
    defer rows.Close()
    
    var empresas []*entity.Empresa
    
    for rows.Next() {
        empresa := &entity.Empresa{}
        err := rows.Scan(
            &empresa.ID,
            &empresa.NomeFantasia,
            &empresa.RazaoSocial,
            &empresa.CNPJ,
            &empresa.DataCadastro,
        )
        if err != nil {
            r.logger.Error("erro ao escanear empresa: %v", err)
            return nil, fmt.Errorf("erro ao escanear empresa: %v", err)
        }
        empresas = append(empresas, empresa)
    }
    
    if err = rows.Err(); err != nil {
        r.logger.Error("erro ao iterar empresas: %v", err)
        return nil, fmt.Errorf("erro ao iterar empresas: %v", err)
    }
    
    return empresas, nil
}

func (r *EmpresaRepository) Update(ctx context.Context, empresa *entity.Empresa) error {
    query := `
        UPDATE empresa 
        SET nome_fantasia = $2, razao_social = $3, cnpj = $4
        WHERE id_empresa = $1
    `
    
    result, err := r.db.ExecContext(ctx, query,
        empresa.ID,
        empresa.NomeFantasia,
        empresa.RazaoSocial,
        empresa.CNPJ,
    )
    
    if err != nil {
        r.logger.Error("erro ao atualizar empresa ID=%d: %v", empresa.ID, err)
        return fmt.Errorf("erro ao atualizar empresa: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("empresa com ID %d não encontrada para atualização", empresa.ID)
    }
    
    return nil
}

func (r *EmpresaRepository) Delete(ctx context.Context, id int) error {
    var count int
    checkQuery := `SELECT COUNT(*) FROM usuario_administrador WHERE id_empresa = $1`
    err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&count)
    if err != nil {
        r.logger.Error("erro ao verificar dependências empresa ID=%d: %v", id, err)
        return fmt.Errorf("erro ao verificar dependências: %v", err)
    }
    
    if count > 0 {
        return fmt.Errorf("não é possível deletar empresa: possui %d usuários administradores vinculados", count)
    }
    
    query := `DELETE FROM empresa WHERE id_empresa = $1`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        r.logger.Error("erro ao deletar empresa ID=%d: %v", id, err)
        return fmt.Errorf("erro ao deletar empresa: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("empresa com ID %d não encontrada para deleção", id)
    }
    
    return nil
}