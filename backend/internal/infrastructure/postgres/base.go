package postgres

import (
    "database/sql"
    "fmt"
    "organizational-climate-survey/backend/pkg/logger"
    _ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
    logger logger.Logger
}

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func NewDB(host, port, user, password, dbname string) (*DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("erro ao conectar com o banco: %v", err)
    }
    
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("erro ao fazer ping no banco: %v", err)
    }
    
    log := logger.New(nil)
    log.Info("Conectado ao PostgreSQL")
    
    return &DB{DB: db, logger: log}, nil
}

func (db *DB) Close() error {
    return db.DB.Close()
}

type Repositories struct {
    Empresa             *EmpresaRepository
    UsuarioAdministrador *UsuarioAdministradorRepository
    Setor               *SetorRepository
    Pesquisa            *PesquisaRepository
    Pergunta            *PerguntaRepository
    Resposta            *RespostaRepository
    Dashboard           *DashboardRepository
    LogAuditoria        *LogAuditoriaRepository
}

func NewRepositories(db *DB) *Repositories {
    return &Repositories{
        Empresa:             NewEmpresaRepository(db),
        UsuarioAdministrador: NewUsuarioAdministradorRepository(db),
        Setor:               NewSetorRepository(db),
        Pesquisa:            NewPesquisaRepository(db),
        Pergunta:            NewPerguntaRepository(db),
        Resposta:            NewRespostaRepository(db),
        Dashboard:           NewDashboardRepository(db),
        LogAuditoria:        NewLogAuditoriaRepository(db),
    }
}