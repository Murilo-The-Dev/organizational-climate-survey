// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

import (
	"time"
)

// UsuarioAdministrador representa um usuário com privilégios administrativos no sistema
type UsuarioAdministrador struct {
	ID           int       `json:"id_user_admin" example:"1"`                                         // Identificador único do usuário
	IDEmpresa    int       `json:"id_empresa" example:"1"`                                            // ID da empresa à qual pertence
	NomeAdmin    string    `json:"nome_admin" example:"Joao Silva"`                                   // Nome completo do administrador
	Email        string    `json:"email" example:"joao@empresa.com"`                                  // Email único para login
	SenhaHash    string    `json:"-"`                                                                 // Hash da senha (oculto em JSON)
	DataCadastro time.Time `json:"data_cadastro" swaggertype:"string" example:"2026-01-01T10:00:00Z"` // Data de criação do usuário
	Status       string    `json:"status" example:"Ativo"`                                            // Estado atual (Ativo, Inativo, Pendente)

	// Relacionamento com Empresa (opcional, para carregamento sob demanda)
	Empresa *Empresa `json:"empresa,omitempty" example:"{\"id_empresa\":1,\"nome_fantasia\":\"Acme\"}"` // Dados da empresa associada
}
