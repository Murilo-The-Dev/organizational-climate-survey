package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
)

type LogAuditoriaRepository struct {
    db *DB
}

func NewLogAuditoriaRepository(db *DB) *LogAuditoriaRepository {
    return &LogAuditoriaRepository{db: db}
}

var _ repository.LogAuditoriaRepository = (*LogAuditoriaRepository)(nil)

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
        return fmt.Errorf("erro ao criar log de auditoria: %v", err)
    }
    
    return nil
}

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
        return nil, fmt.Errorf("erro ao buscar log de auditoria: %v", err)
    }
    
    return log, nil
}

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
            return nil, fmt.Errorf("erro ao escanear log de auditoria: %v", err)
        }
        logs = append(logs, log)
    }
    
    return logs, nil
}

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
            return nil, fmt.Errorf("erro ao escanear log de auditoria: %v", err)
        }
        logs = append(logs, log)
    }
    
    return logs, nil
}

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
            return nil, fmt.Errorf("erro ao escanear log de auditoria: %v", err)
        }
        logs = append(logs, log)
    }
    
    return logs, nil
}