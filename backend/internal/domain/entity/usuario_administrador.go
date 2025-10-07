// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

import (
	"time"
)

// UsuarioAdministrador representa um usuário com privilégios administrativos no sistema
type UsuarioAdministrador struct {
	ID           int       `json:"id_user_admin"` // Identificador único do usuário
	IDEmpresa    int       `json:"id_empresa"`    // ID da empresa à qual pertence
	NomeAdmin    string    `json:"nome_admin"`    // Nome completo do administrador
	Email        string    `json:"email"`         // Email único para login
	SenhaHash    string    `json:"-"`             // Hash da senha (oculto em JSON)
	DataCadastro time.Time `json:"data_cadastro"` // Data de criação do usuário
	Status       string    `json:"status"`        // Estado atual (Ativo, Inativo, Pendente)

	// Relacionamento com Empresa (opcional, para carregamento sob demanda)
	Empresa *Empresa `json:"empresa,omitempty"` // Dados da empresa associada
}
