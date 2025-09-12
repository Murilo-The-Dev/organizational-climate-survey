package postgres

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq" // Driver PostgreSQL
)

// DB representa a conexão com o banco de dados
type DB struct {
    *sql.DB
}

// Config contém as configurações para conexão com o banco de dados
type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// NewDB cria uma nova conexão com PostgreSQL
func NewDB(host, port, user, password, dbname string) (*DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("erro ao conectar com o banco: %v", err)
    }
    
    // Testa a conexão
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("erro ao fazer ping no banco: %v", err)
    }
    
    log.Println("Conectado ao PostgreSQL com sucesso!")
    
    return &DB{db}, nil
}

// Close fecha a conexão com o banco
func (db *DB) Close() error {
    return db.DB.Close()
}

// Repositories struct agrupa todos os repositórios
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

// NewRepositories cria todas as instâncias dos repositórios
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