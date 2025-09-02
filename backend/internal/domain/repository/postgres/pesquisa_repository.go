package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
)

type PesquisaRepository struct {
    db *DB
}

// NewPesquisaRepository cria uma nova instância do repositório de pesquisa
func NewPesquisaRepository(db *DB) *PesquisaRepository {
    return &PesquisaRepository{db: db}
}

// Verifica se implementa a interface
var _ repository.PesquisaRepository = (*PesquisaRepository)(nil)

// Create insere uma nova pesquisa
func (r *PesquisaRepository) Create(ctx context.Context, pesquisa *entity.Pesquisa) error {
    query := `
        INSERT INTO pesquisa (id_empresa, id_user_admin, id_setor, titulo, descricao, 
                            data_criacao, data_abertura, data_fechamento, status, 
                            link_acesso, qrcode_path, config_recorrencia, anonimato)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id_pesquisa
    `
    
    err := r.db.QueryRowContext(ctx, query,
        pesquisa.IDEmpresa,
        pesquisa.IDUserAdmin,
        pesquisa.IDSetor,
        pesquisa.Titulo,
        pesquisa.Descricao,
        pesquisa.DataCriacao,
        pesquisa.DataAbertura,
        pesquisa.DataFechamento,
        pesquisa.Status,
        pesquisa.LinkAcesso,
        pesquisa.QRCodePath,
        pesquisa.ConfigRecorrencia,
        pesquisa.Anonimato,
    ).Scan(&pesquisa.ID)
    
    if err != nil {
        return fmt.Errorf("erro ao criar pesquisa: %v", err)
    }
    
    return nil
}

// GetByID busca uma pesquisa por ID
func (r *PesquisaRepository) GetByID(ctx context.Context, id int) (*entity.Pesquisa, error) {
    pesquisa := &entity.Pesquisa{}
    query := `
        SELECT id_pesquisa, id_empresa, id_user_admin, id_setor, titulo, descricao,
               data_criacao, data_abertura, data_fechamento, status, link_acesso,
               qrcode_path, config_recorrencia, anonimato
        FROM pesquisa
        WHERE id_pesquisa = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &pesquisa.ID,
        &pesquisa.IDEmpresa,
        &pesquisa.IDUserAdmin,
        &pesquisa.IDSetor,
        &pesquisa.Titulo,
        &pesquisa.Descricao,
        &pesquisa.DataCriacao,
        &pesquisa.DataAbertura,
        &pesquisa.DataFechamento,
        &pesquisa.Status,
        &pesquisa.LinkAcesso,
        &pesquisa.QRCodePath,
        &pesquisa.ConfigRecorrencia,
        &pesquisa.Anonimato,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("pesquisa com ID %d não encontrada", id)
        }
        return nil, fmt.Errorf("erro ao buscar pesquisa: %v", err)
    }
    
    return pesquisa, nil
}

// GetByLinkAcesso busca uma pesquisa pelo link de acesso (para respondentes)
func (r *PesquisaRepository) GetByLinkAcesso(ctx context.Context, link string) (*entity.Pesquisa, error) {
    pesquisa := &entity.Pesquisa{}
    query := `
        SELECT id_pesquisa, id_empresa, id_user_admin, id_setor, titulo, descricao,
               data_criacao, data_abertura, data_fechamento, status, link_acesso,
               qrcode_path, config_recorrencia, anonimato
        FROM pesquisa
        WHERE link_acesso = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, link).Scan(
        &pesquisa.ID,
        &pesquisa.IDEmpresa,
        &pesquisa.IDUserAdmin,
        &pesquisa.IDSetor,
        &pesquisa.Titulo,
        &pesquisa.Descricao,
        &pesquisa.DataCriacao,
        &pesquisa.DataAbertura,
        &pesquisa.DataFechamento,
        &pesquisa.Status,
        &pesquisa.LinkAcesso,
        &pesquisa.QRCodePath,
        &pesquisa.ConfigRecorrencia,
        &pesquisa.Anonimato,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("pesquisa com link %s não encontrada", link)
        }
        return nil, fmt.Errorf("erro ao buscar pesquisa: %v", err)
    }
    
    return pesquisa, nil
}

// ListByEmpresa retorna todas as pesquisas de uma empresa
func (r *PesquisaRepository) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Pesquisa, error) {
    query := `
        SELECT id_pesquisa, id_empresa, id_user_admin, id_setor, titulo, descricao,
               data_criacao, data_abertura, data_fechamento, status, link_acesso,
               qrcode_path, config_recorrencia, anonimato
        FROM pesquisa
        WHERE id_empresa = $1
        ORDER BY data_criacao DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, empresaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao listar pesquisas: %v", err)
    }
    defer rows.Close()
    
    var pesquisas []*entity.Pesquisa
    
    for rows.Next() {
        pesquisa := &entity.Pesquisa{}
        err := rows.Scan(
            &pesquisa.ID,
            &pesquisa.IDEmpresa,
            &pesquisa.IDUserAdmin,
            &pesquisa.IDSetor,
            &pesquisa.Titulo,
            &pesquisa.Descricao,
            &pesquisa.DataCriacao,
            &pesquisa.DataAbertura,
            &pesquisa.DataFechamento,
            &pesquisa.Status,
            &pesquisa.LinkAcesso,
            &pesquisa.QRCodePath,
            &pesquisa.ConfigRecorrencia,
            &pesquisa.Anonimato,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear pesquisa: %v", err)
        }
        pesquisas = append(pesquisas, pesquisa)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar pesquisas: %v", err)
    }
    
    return pesquisas, nil
}

// ListBySetor retorna pesquisas de um setor específico
func (r *PesquisaRepository) ListBySetor(ctx context.Context, setorID int) ([]*entity.Pesquisa, error) {
    query := `
        SELECT id_pesquisa, id_empresa, id_user_admin, id_setor, titulo, descricao,
               data_criacao, data_abertura, data_fechamento, status, link_acesso,
               qrcode_path, config_recorrencia, anonimato
        FROM pesquisa
        WHERE id_setor = $1
        ORDER BY data_criacao DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, setorID)
    if err != nil {
        return nil, fmt.Errorf("erro ao listar pesquisas por setor: %v", err)
    }
    defer rows.Close()
    
    var pesquisas []*entity.Pesquisa
    
    for rows.Next() {
        pesquisa := &entity.Pesquisa{}
        err := rows.Scan(
            &pesquisa.ID,
            &pesquisa.IDEmpresa,
            &pesquisa.IDUserAdmin,
            &pesquisa.IDSetor,
            &pesquisa.Titulo,
            &pesquisa.Descricao,
            &pesquisa.DataCriacao,
            &pesquisa.DataAbertura,
            &pesquisa.DataFechamento,
            &pesquisa.Status,
            &pesquisa.LinkAcesso,
            &pesquisa.QRCodePath,
            &pesquisa.ConfigRecorrencia,
            &pesquisa.Anonimato,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear pesquisa: %v", err)
        }
        pesquisas = append(pesquisas, pesquisa)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar pesquisas: %v", err)
    }
    
    return pesquisas, nil
}

// ListByStatus retorna pesquisas por status
func (r *PesquisaRepository) ListByStatus(ctx context.Context, empresaID int, status string) ([]*entity.Pesquisa, error) {
    query := `
        SELECT id_pesquisa, id_empresa, id_user_admin, id_setor, titulo, descricao,
               data_criacao, data_abertura, data_fechamento, status, link_acesso,
               qrcode_path, config_recorrencia, anonimato
        FROM pesquisa
        WHERE id_empresa = $1 AND status = $2
        ORDER BY data_criacao DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, empresaID, status)
    if err != nil {
        return nil, fmt.Errorf("erro ao listar pesquisas por status: %v", err)
    }
    defer rows.Close()
    
    var pesquisas []*entity.Pesquisa
    
    for rows.Next() {
        pesquisa := &entity.Pesquisa{}
        err := rows.Scan(
            &pesquisa.ID,
            &pesquisa.IDEmpresa,
            &pesquisa.IDUserAdmin,
            &pesquisa.IDSetor,
            &pesquisa.Titulo,
            &pesquisa.Descricao,
            &pesquisa.DataCriacao,
            &pesquisa.DataAbertura,
            &pesquisa.DataFechamento,
            &pesquisa.Status,
            &pesquisa.LinkAcesso,
            &pesquisa.QRCodePath,
            &pesquisa.ConfigRecorrencia,
            &pesquisa.Anonimato,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear pesquisa: %v", err)
        }
        pesquisas = append(pesquisas, pesquisa)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar pesquisas: %v", err)
    }
    
    return pesquisas, nil
}

// ListActive retorna pesquisas ativas
func (r *PesquisaRepository) ListActive(ctx context.Context, empresaID int) ([]*entity.Pesquisa, error) {
    query := `
        SELECT id_pesquisa, id_empresa, id_user_admin, id_setor, titulo, descricao,
               data_criacao, data_abertura, data_fechamento, status, link_acesso,
               qrcode_path, config_recorrencia, anonimato
        FROM pesquisa
        WHERE id_empresa = $1 AND status = 'Ativa'
        AND (data_abertura IS NULL OR data_abertura <= NOW())
        AND (data_fechamento IS NULL OR data_fechamento > NOW())
        ORDER BY data_criacao DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, empresaID)
    if err != nil {
        return nil, fmt.Errorf("erro ao listar pesquisas ativas: %v", err)
    }
    defer rows.Close()
    
    var pesquisas []*entity.Pesquisa
    
    for rows.Next() {
        pesquisa := &entity.Pesquisa{}
        err := rows.Scan(
            &pesquisa.ID,
            &pesquisa.IDEmpresa,
            &pesquisa.IDUserAdmin,
            &pesquisa.IDSetor,
            &pesquisa.Titulo,
            &pesquisa.Descricao,
            &pesquisa.DataCriacao,
            &pesquisa.DataAbertura,
            &pesquisa.DataFechamento,
            &pesquisa.Status,
            &pesquisa.LinkAcesso,
            &pesquisa.QRCodePath,
            &pesquisa.ConfigRecorrencia,
            &pesquisa.Anonimato,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao escanear pesquisa: %v", err)
        }
        pesquisas = append(pesquisas, pesquisa)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar pesquisas: %v", err)
    }
    
    return pesquisas, nil
}

// Update atualiza uma pesquisa
func (r *PesquisaRepository) Update(ctx context.Context, pesquisa *entity.Pesquisa) error {
    query := `
        UPDATE pesquisa 
        SET titulo = $2, descricao = $3, data_abertura = $4, data_fechamento = $5,
            status = $6, config_recorrencia = $7
        WHERE id_pesquisa = $1
    `
    
    result, err := r.db.ExecContext(ctx, query,
        pesquisa.ID,
        pesquisa.Titulo,
        pesquisa.Descricao,
        pesquisa.DataAbertura,
        pesquisa.DataFechamento,
        pesquisa.Status,
        pesquisa.ConfigRecorrencia,
    )
    
    if err != nil {
        return fmt.Errorf("erro ao atualizar pesquisa: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("pesquisa com ID %d não encontrada para atualização", pesquisa.ID)
    }
    
    return nil
}

// UpdateStatus atualiza apenas o status de uma pesquisa
func (r *PesquisaRepository) UpdateStatus(ctx context.Context, id int, status string) error {
    query := `
        UPDATE pesquisa 
        SET status = $2
        WHERE id_pesquisa = $1
    `
    
    result, err := r.db.ExecContext(ctx, query, id, status)
    if err != nil {
        return fmt.Errorf("erro ao atualizar status da pesquisa: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("pesquisa com ID %d não encontrada para atualização de status", id)
    }
    
    return nil
}

// Delete remove uma pesquisa
func (r *PesquisaRepository) Delete(ctx context.Context, id int) error {
    // Primeiro verificamos se há respostas vinculadas
    var count int
    checkQuery := `
        SELECT COUNT(*) 
        FROM resposta r
        INNER JOIN pergunta p ON r.id_pergunta = p.id_pergunta
        WHERE p.id_pesquisa = $1
    `
    err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&count)
    if err != nil {
        return fmt.Errorf("erro ao verificar dependências: %v", err)
    }
    
    if count > 0 {
        // Em vez de deletar, arquiva a pesquisa
        return r.UpdateStatus(ctx, id, "Arquivada")
    }
    
    // Se não há respostas, pode deletar
    query := `DELETE FROM pesquisa WHERE id_pesquisa = $1`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("erro ao deletar pesquisa: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("pesquisa com ID %d não encontrada para deleção", id)
    }
    
    return nil
}