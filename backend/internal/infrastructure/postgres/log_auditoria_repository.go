// Package postgres implementa o repositório de LogAuditoria usando PostgreSQL.
// Fornece operações para registro e consulta de logs de auditoria do sistema.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/logger"
)

// LogAuditoriaRepository implementa a interface repository.LogAuditoriaRepository
type LogAuditoriaRepository struct {
	db     *DB           // Conexão com o banco de dados
	logger logger.Logger // Logger para operações do repositório
}

// NewLogAuditoriaRepository cria uma nova instância do repositório
func NewLogAuditoriaRepository(db *DB) *LogAuditoriaRepository {
	return &LogAuditoriaRepository{
		db:     db,
		logger: db.logger,
	}
}

// Garante que LogAuditoriaRepository implementa a interface correta
var _ repository.LogAuditoriaRepository = (*LogAuditoriaRepository)(nil)

// Create registra um novo log de auditoria no banco de dados
func (r *LogAuditoriaRepository) Create(ctx context.Context, log *entity.LogAuditoria) error {
	query := `
        INSERT INTO log_auditoria (id_user_admin, timestamp, acao_realizada, detalhes, endereco_ip)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id_log
    `

	err := r.db.QueryRowContext(ctx, query,
		log.IDUserAdmin,
		log.TimeStamp,
		log.AcaoRealizada,
		log.Detalhes,
		log.EnderecoIP,
	).Scan(&log.ID)

	if err != nil {
		r.logger.Error("erro ao criar log auditoria: %v", err)
		return fmt.Errorf("erro ao criar log de auditoria: %v", err)
	}

	return nil
}

// GetByID busca um log de auditoria pelo seu ID
func (r *LogAuditoriaRepository) GetByID(ctx context.Context, id int) (*entity.LogAuditoria, error) {
	log := &entity.LogAuditoria{}
	query := `
        SELECT id_log, id_user_admin, timestamp, acao_realizada, detalhes, endereco_ip
        FROM log_auditoria
        WHERE id_log = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.IDUserAdmin,
		&log.TimeStamp,
		&log.AcaoRealizada,
		&log.Detalhes,
		&log.EnderecoIP,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("log de auditoria com ID %d não encontrado", id)
		}
		r.logger.Error("erro ao buscar log auditoria ID=%d: %v", id, err)
		return nil, fmt.Errorf("erro ao buscar log de auditoria: %v", err)
	}

	return log, nil
}

// ListByEmpresa lista logs de auditoria de uma empresa com paginação
func (r *LogAuditoriaRepository) ListByEmpresa(ctx context.Context, empresaID int, limit, offset int) ([]*entity.LogAuditoria, error) {
	query := `
        SELECT l.id_log, l.id_user_admin, l.timestamp, l.acao_realizada, l.detalhes, l.endereco_ip
        FROM log_auditoria l
        INNER JOIN usuario_administrador ua ON l.id_user_admin = ua.id_user_admin
        WHERE ua.id_empresa = $1
        ORDER BY l.timestamp DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.QueryContext(ctx, query, empresaID, limit, offset)
	if err != nil {
		r.logger.Error("erro ao listar logs empresa ID=%d: %v", empresaID, err)
		return nil, fmt.Errorf("erro ao listar logs de auditoria: %v", err)
	}
	defer rows.Close()

	var logs []*entity.LogAuditoria

	for rows.Next() {
		log := &entity.LogAuditoria{}
		err := rows.Scan(
			&log.ID,
			&log.IDUserAdmin,
			&log.TimeStamp,
			&log.AcaoRealizada,
			&log.Detalhes,
			&log.EnderecoIP,
		)
		if err != nil {
			r.logger.Error("erro ao escanear log auditoria: %v", err)
			return nil, fmt.Errorf("erro ao escanear log de auditoria: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ListByUsuarioAdmin lista logs de auditoria de um usuário específico
func (r *LogAuditoriaRepository) ListByUsuarioAdmin(ctx context.Context, userAdminID int, limit, offset int) ([]*entity.LogAuditoria, error) {
	query := `
        SELECT id_log, id_user_admin, timestamp, acao_realizada, detalhes, endereco_ip
        FROM log_auditoria
        WHERE id_user_admin = $1
        ORDER BY timestamp DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.QueryContext(ctx, query, userAdminID, limit, offset)
	if err != nil {
		r.logger.Error("erro ao listar logs usuário ID=%d: %v", userAdminID, err)
		return nil, fmt.Errorf("erro ao listar logs por usuário: %v", err)
	}
	defer rows.Close()

	var logs []*entity.LogAuditoria

	for rows.Next() {
		log := &entity.LogAuditoria{}
		err := rows.Scan(
			&log.ID,
			&log.IDUserAdmin,
			&log.TimeStamp,
			&log.AcaoRealizada,
			&log.Detalhes,
			&log.EnderecoIP,
		)
		if err != nil {
			r.logger.Error("erro ao escanear log auditoria: %v", err)
			return nil, fmt.Errorf("erro ao escanear log de auditoria: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ListByDateRange lista logs de auditoria dentro de um intervalo de datas
func (r *LogAuditoriaRepository) ListByDateRange(ctx context.Context, empresaID int, startDate, endDate string) ([]*entity.LogAuditoria, error) {
	query := `
        SELECT l.id_log, l.id_user_admin, l.timestamp, l.acao_realizada, l.detalhes, l.endereco_ip
        FROM log_auditoria l
        INNER JOIN usuario_administrador ua ON l.id_user_admin = ua.id_user_admin
        WHERE ua.id_empresa = $1 
        AND l.timestamp >= $2 
        AND l.timestamp <= $3
        ORDER BY l.timestamp DESC
    `

	rows, err := r.db.QueryContext(ctx, query, empresaID, startDate, endDate)
	if err != nil {
		r.logger.Error("erro ao listar logs por período empresa ID=%d: %v", empresaID, err)
		return nil, fmt.Errorf("erro ao listar logs por período: %v", err)
	}
	defer rows.Close()

	var logs []*entity.LogAuditoria

	for rows.Next() {
		log := &entity.LogAuditoria{}
		err := rows.Scan(
			&log.ID,
			&log.IDUserAdmin,
			&log.TimeStamp,
			&log.AcaoRealizada,
			&log.Detalhes,
			&log.EnderecoIP,
		)
		if err != nil {
			r.logger.Error("erro ao escanear log auditoria: %v", err)
			return nil, fmt.Errorf("erro ao escanear log de auditoria: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
