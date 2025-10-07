// Package postgres implementa o repositório de UsuarioAdministrador usando PostgreSQL.
// Fornece operações CRUD e consultas específicas para gerenciamento de usuários administradores.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/logger"
)

// UsuarioAdministradorRepository implementa a interface repository.UsuarioAdministradorRepository
type UsuarioAdministradorRepository struct {
	db     *DB           // Conexão com o banco de dados
	logger logger.Logger // Logger para operações do repositório
}

// NewUsuarioAdministradorRepository cria uma nova instância do repositório
func NewUsuarioAdministradorRepository(db *DB) *UsuarioAdministradorRepository {
	return &UsuarioAdministradorRepository{
		db:     db,
		logger: db.logger,
	}
}

var _ repository.UsuarioAdministradorRepository = (*UsuarioAdministradorRepository)(nil)

// Create insere um novo usuário administrador no banco de dados
// Retorna o ID gerado através do RETURNING
func (r *UsuarioAdministradorRepository) Create(ctx context.Context, usuario *entity.UsuarioAdministrador) error {
	query := `
        INSERT INTO usuario_administrador (id_empresa, nome_admin, email, senha_hash, data_cadastro, status)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id_user_admin
    `

	err := r.db.QueryRowContext(ctx, query,
		usuario.IDEmpresa,
		usuario.NomeAdmin,
		usuario.Email,
		usuario.SenhaHash,
		usuario.DataCadastro,
		usuario.Status,
	).Scan(&usuario.ID)

	if err != nil {
		r.logger.Error("erro ao criar usuário admin email=%s: %v", usuario.Email, err)
		return fmt.Errorf("erro ao criar usuário administrador: %v", err)
	}

	return nil
}

// GetByID busca um usuário administrador pelo seu ID
// Retorna erro específico quando não encontrado
func (r *UsuarioAdministradorRepository) GetByID(ctx context.Context, id int) (*entity.UsuarioAdministrador, error) {
	usuario := &entity.UsuarioAdministrador{}
	query := `
        SELECT id_user_admin, id_empresa, nome_admin, email, senha_hash, data_cadastro, status
        FROM usuario_administrador
        WHERE id_user_admin = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&usuario.ID,
		&usuario.IDEmpresa,
		&usuario.NomeAdmin,
		&usuario.Email,
		&usuario.SenhaHash,
		&usuario.DataCadastro,
		&usuario.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário administrador com ID %d não encontrado", id)
		}
		r.logger.Error("erro ao buscar usuário admin ID=%d: %v", id, err)
		return nil, fmt.Errorf("erro ao buscar usuário administrador: %v", err)
	}

	return usuario, nil
}

// GetByEmail busca um usuário administrador pelo email
// Retorna erro específico quando não encontrado
func (r *UsuarioAdministradorRepository) GetByEmail(ctx context.Context, email string) (*entity.UsuarioAdministrador, error) {
	usuario := &entity.UsuarioAdministrador{}
	query := `
        SELECT id_user_admin, id_empresa, nome_admin, email, senha_hash, data_cadastro, status
        FROM usuario_administrador
        WHERE email = $1
    `

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&usuario.ID,
		&usuario.IDEmpresa,
		&usuario.NomeAdmin,
		&usuario.Email,
		&usuario.SenhaHash,
		&usuario.DataCadastro,
		&usuario.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário administrador com email %s não encontrado", email)
		}
		r.logger.Error("erro ao buscar usuário por email=%s: %v", email, err)
		return nil, fmt.Errorf("erro ao buscar usuário administrador: %v", err)
	}

	return usuario, nil
}

// ListByEmpresa lista todos os usuários administradores de uma empresa
// Ordenados por data de cadastro decrescente
func (r *UsuarioAdministradorRepository) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.UsuarioAdministrador, error) {
	query := `
        SELECT id_user_admin, id_empresa, nome_admin, email, senha_hash, data_cadastro, status
        FROM usuario_administrador
        WHERE id_empresa = $1
        ORDER BY data_cadastro DESC
    `

	rows, err := r.db.QueryContext(ctx, query, empresaID)
	if err != nil {
		r.logger.Error("erro ao listar usuários empresa ID=%d: %v", empresaID, err)
		return nil, fmt.Errorf("erro ao listar usuários administradores: %v", err)
	}
	defer rows.Close()

	var usuarios []*entity.UsuarioAdministrador

	for rows.Next() {
		usuario := &entity.UsuarioAdministrador{}
		err := rows.Scan(
			&usuario.ID,
			&usuario.IDEmpresa,
			&usuario.NomeAdmin,
			&usuario.Email,
			&usuario.SenhaHash,
			&usuario.DataCadastro,
			&usuario.Status,
		)
		if err != nil {
			r.logger.Error("erro ao escanear usuário admin: %v", err)
			return nil, fmt.Errorf("erro ao escanear usuário administrador: %v", err)
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("erro ao iterar usuários admin: %v", err)
		return nil, fmt.Errorf("erro ao iterar usuários administradores: %v", err)
	}

	return usuarios, nil
}

// ListByStatus lista usuários administradores de uma empresa filtrados por status
// Status podem ser: Ativo, Inativo, Pendente
func (r *UsuarioAdministradorRepository) ListByStatus(ctx context.Context, empresaID int, status string) ([]*entity.UsuarioAdministrador, error) {
	query := `
        SELECT id_user_admin, id_empresa, nome_admin, email, senha_hash, data_cadastro, status
        FROM usuario_administrador
        WHERE id_empresa = $1 AND status = $2
        ORDER BY data_cadastro DESC
    `

	rows, err := r.db.QueryContext(ctx, query, empresaID, status)
	if err != nil {
		r.logger.Error("erro ao listar usuários por status empresa ID=%d: %v", empresaID, err)
		return nil, fmt.Errorf("erro ao listar usuários administradores por status: %v", err)
	}
	defer rows.Close()

	var usuarios []*entity.UsuarioAdministrador

	for rows.Next() {
		usuario := &entity.UsuarioAdministrador{}
		err := rows.Scan(
			&usuario.ID,
			&usuario.IDEmpresa,
			&usuario.NomeAdmin,
			&usuario.Email,
			&usuario.SenhaHash,
			&usuario.DataCadastro,
			&usuario.Status,
		)
		if err != nil {
			r.logger.Error("erro ao escanear usuário admin: %v", err)
			return nil, fmt.Errorf("erro ao escanear usuário administrador: %v", err)
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("erro ao iterar usuários admin: %v", err)
		return nil, fmt.Errorf("erro ao iterar usuários administradores: %v", err)
	}

	return usuarios, nil
}

// Update atualiza os dados de um usuário administrador existente
// Retorna erro se o usuário não for encontrado
func (r *UsuarioAdministradorRepository) Update(ctx context.Context, usuario *entity.UsuarioAdministrador) error {
	query := `
        UPDATE usuario_administrador 
        SET nome_admin = $2, email = $3, status = $4
        WHERE id_user_admin = $1
    `

	result, err := r.db.ExecContext(ctx, query,
		usuario.ID,
		usuario.NomeAdmin,
		usuario.Email,
		usuario.Status,
	)

	if err != nil {
		r.logger.Error("erro ao atualizar usuário admin ID=%d: %v", usuario.ID, err)
		return fmt.Errorf("erro ao atualizar usuário administrador: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário administrador com ID %d não encontrado para atualização", usuario.ID)
	}

	return nil
}

// UpdatePassword atualiza apenas a senha hash do usuário
// Usado para operações de redefinição de senha
func (r *UsuarioAdministradorRepository) UpdatePassword(ctx context.Context, id int, senhaHash string) error {
	query := `
        UPDATE usuario_administrador 
        SET senha_hash = $2
        WHERE id_user_admin = $1
    `

	result, err := r.db.ExecContext(ctx, query, id, senhaHash)
	if err != nil {
		r.logger.Error("erro ao atualizar senha usuário ID=%d: %v", id, err)
		return fmt.Errorf("erro ao atualizar senha: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário administrador com ID %d não encontrado para atualização de senha", id)
	}

	return nil
}

// UpdateStatus atualiza apenas o status do usuário
// Usado para ativar/inativar contas
func (r *UsuarioAdministradorRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
        UPDATE usuario_administrador 
        SET status = $2
        WHERE id_user_admin = $1
    `

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		r.logger.Error("erro ao atualizar status usuário ID=%d: %v", id, err)
		return fmt.Errorf("erro ao atualizar status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário administrador com ID %d não encontrado para atualização de status", id)
	}

	return nil
}

// Delete remove um usuário administrador ou inativa se houver dependências
// Verifica pesquisas associadas antes da deleção
func (r *UsuarioAdministradorRepository) Delete(ctx context.Context, id int) error {
	var count int
	checkQuery := `SELECT COUNT(*) FROM pesquisa WHERE id_user_admin = $1`
	err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&count)
	if err != nil {
		r.logger.Error("erro ao verificar dependências usuário ID=%d: %v", id, err)
		return fmt.Errorf("erro ao verificar dependências: %v", err)
	}

	if count > 0 {
		return r.UpdateStatus(ctx, id, "Inativo")
	}

	query := `DELETE FROM usuario_administrador WHERE id_user_admin = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("erro ao deletar usuário admin ID=%d: %v", id, err)
		return fmt.Errorf("erro ao deletar usuário administrador: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário administrador com ID %d não encontrado para deleção", id)
	}

	return nil
}
