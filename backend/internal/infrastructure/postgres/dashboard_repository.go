// Package postgres implementa o repositório de Dashboard usando PostgreSQL.
// Fornece operações CRUD e consultas específicas para gerenciamento de dashboards.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/logger"
)

// DashboardRepository implementa a interface repository.DashboardRepository
type DashboardRepository struct {
	db     *DB           // Conexão com o banco de dados
	logger logger.Logger // Logger para operações do repositório
}

// NewDashboardRepository cria uma nova instância do repositório
func NewDashboardRepository(db *DB) *DashboardRepository {
	return &DashboardRepository{
		db:     db,
		logger: db.logger,
	}
}

// Garante que DashboardRepository implementa a interface correta
var _ repository.DashboardRepository = (*DashboardRepository)(nil)

// Create insere um novo dashboard no banco de dados
// Retorna o ID gerado através do RETURNING
func (r *DashboardRepository) Create(ctx context.Context, dashboard *entity.Dashboard) error {
	query := `
        INSERT INTO dashboard (id_pesquisa, titulo, data_criacao, config_filtros)
        VALUES ($1, $2, $3, $4)
        RETURNING id_dashboard
    `

	err := r.db.QueryRowContext(ctx, query,
		dashboard.IDPesquisa,
		dashboard.Titulo,
		dashboard.DataCriacao,
		dashboard.ConfigFiltros,
	).Scan(&dashboard.ID)

	if err != nil {
		r.logger.Error("erro ao criar dashboard pesquisa ID=%d: %v", dashboard.IDPesquisa, err)
		return fmt.Errorf("erro ao criar dashboard: %v", err)
	}

	return nil
}

// GetByID busca um dashboard pelo seu ID
// Retorna erro específico quando não encontrado
func (r *DashboardRepository) GetByID(ctx context.Context, id int) (*entity.Dashboard, error) {
	dashboard := &entity.Dashboard{}
	query := `
        SELECT id_dashboard, id_pesquisa, titulo, data_criacao, config_filtros
        FROM dashboard
        WHERE id_dashboard = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dashboard.ID,
		&dashboard.IDPesquisa,
		&dashboard.Titulo,
		&dashboard.DataCriacao,
		&dashboard.ConfigFiltros,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dashboard com ID %d não encontrado", id)
		}
		r.logger.Error("erro ao buscar dashboard ID=%d: %v", id, err)
		return nil, fmt.Errorf("erro ao buscar dashboard: %v", err)
	}

	return dashboard, nil
}

// GetByPesquisaID busca um dashboard pelo ID da pesquisa associada
// Retorna erro específico quando não encontrado
func (r *DashboardRepository) GetByPesquisaID(ctx context.Context, pesquisaID int) (*entity.Dashboard, error) {
	dashboard := &entity.Dashboard{}
	query := `
        SELECT id_dashboard, id_pesquisa, titulo, data_criacao, config_filtros
        FROM dashboard
        WHERE id_pesquisa = $1
    `

	err := r.db.QueryRowContext(ctx, query, pesquisaID).Scan(
		&dashboard.ID,
		&dashboard.IDPesquisa,
		&dashboard.Titulo,
		&dashboard.DataCriacao,
		&dashboard.ConfigFiltros,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dashboard para pesquisa ID %d não encontrado", pesquisaID)
		}
		r.logger.Error("erro ao buscar dashboard pesquisa ID=%d: %v", pesquisaID, err)
		return nil, fmt.Errorf("erro ao buscar dashboard: %v", err)
	}

	return dashboard, nil
}

// ListByEmpresa lista todos os dashboards de uma empresa
// Ordena por data de criação decrescente
func (r *DashboardRepository) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Dashboard, error) {
	query := `
        SELECT d.id_dashboard, d.id_pesquisa, d.titulo, d.data_criacao, d.config_filtros
        FROM dashboard d
        INNER JOIN pesquisa p ON d.id_pesquisa = p.id_pesquisa
        WHERE p.id_empresa = $1
        ORDER BY d.data_criacao DESC
    `

	rows, err := r.db.QueryContext(ctx, query, empresaID)
	if err != nil {
		r.logger.Error("erro ao listar dashboards empresa ID=%d: %v", empresaID, err)
		return nil, fmt.Errorf("erro ao listar dashboards: %v", err)
	}
	defer rows.Close()

	var dashboards []*entity.Dashboard

	for rows.Next() {
		dashboard := &entity.Dashboard{}
		err := rows.Scan(
			&dashboard.ID,
			&dashboard.IDPesquisa,
			&dashboard.Titulo,
			&dashboard.DataCriacao,
			&dashboard.ConfigFiltros,
		)
		if err != nil {
			r.logger.Error("erro ao escanear dashboard: %v", err)
			return nil, fmt.Errorf("erro ao escanear dashboard: %v", err)
		}
		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}

// Update atualiza os dados de um dashboard existente
// Retorna erro se o dashboard não for encontrado
func (r *DashboardRepository) Update(ctx context.Context, dashboard *entity.Dashboard) error {
	query := `
        UPDATE dashboard 
        SET titulo = $2, config_filtros = $3
        WHERE id_dashboard = $1
    `

	result, err := r.db.ExecContext(ctx, query,
		dashboard.ID,
		dashboard.Titulo,
		dashboard.ConfigFiltros,
	)

	if err != nil {
		r.logger.Error("erro ao atualizar dashboard ID=%d: %v", dashboard.ID, err)
		return fmt.Errorf("erro ao atualizar dashboard: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("dashboard com ID %d não encontrado", dashboard.ID)
	}

	return nil
}

// Delete remove um dashboard do banco de dados
// Retorna erro se o dashboard não for encontrado
func (r *DashboardRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM dashboard WHERE id_dashboard = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("erro ao deletar dashboard ID=%d: %v", id, err)
		return fmt.Errorf("erro ao deletar dashboard: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("dashboard com ID %d não encontrado", id)
	}

	return nil
}
