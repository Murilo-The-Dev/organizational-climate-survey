// Package postgres implementa a camada de acesso a dados usando PostgreSQL
// Fornece estruturas e métodos para conexão e gerenciamento do banco de dados
package postgres

import (
	"database/sql"
	"fmt"
	"organizational-climate-survey/backend/pkg/logger"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// DB encapsula conexão SQL e logger para operações no banco
type DB struct {
	*sql.DB               // Conexão com banco de dados
	logger  logger.Logger // Logger para operações do banco
}

// Config define os parâmetros necessários para conexão com PostgreSQL
type Config struct {
	Host     string // Endereço do servidor
	Port     string // Porta do servidor
	User     string // Usuário do banco
	Password string // Senha do usuário
	DBName   string // Nome do banco de dados
	SSLMode  string // Modo SSL (disable, require, etc)
}

// NewDB cria uma nova conexão com o banco de dados PostgreSQL
// Retorna erro se a conexão falhar ou o banco estiver inacessível
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

// Close fecha a conexão com o banco de dados
func (db *DB) Close() error {
	return db.DB.Close()
}

// Repositories agrupa todos os repositórios da aplicação
// Facilita o acesso centralizado aos repositórios
type Repositories struct {
	Empresa              *EmpresaRepository
	UsuarioAdministrador *UsuarioAdministradorRepository
	Setor                *SetorRepository
	Pesquisa             *PesquisaRepository
	Pergunta             *PerguntaRepository
	Resposta             *RespostaRepository
	SubmissaoPesquisa    *SubmissaoPesquisaRepository // NOVO
	Dashboard            *DashboardRepository
	LogAuditoria         *LogAuditoriaRepository
}

// NewRepositories inicializa todos os repositórios com a conexão fornecida
// Retorna uma estrutura com todos os repositórios prontos para uso
func NewRepositories(db *DB) *Repositories {
	return &Repositories{
		Empresa:              NewEmpresaRepository(db),
		UsuarioAdministrador: NewUsuarioAdministradorRepository(db),
		Setor:                NewSetorRepository(db),
		Pesquisa:             NewPesquisaRepository(db),
		Pergunta:             NewPerguntaRepository(db),
		Resposta:             NewRespostaRepository(db),
		SubmissaoPesquisa:    NewSubmissaoPesquisaRepository(db), // NOVO
		Dashboard:            NewDashboardRepository(db),
		LogAuditoria:         NewLogAuditoriaRepository(db),
	}
}