// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para empresas do sistema de pesquisa.
package entity

import "time"

// Empresa representa uma organização cliente do sistema
type Empresa struct {
	ID           int       `json:"id_empresa" example:"1"`                                            // Identificador único da empresa
	NomeFantasia string    `json:"nome_fantasia" example:"Acme"`                                      // Nome comercial
	RazaoSocial  string    `json:"razao_social" example:"Acme LTDA"`                                  // Nome jurídico registrado
	CNPJ         string    `json:"cnpj" example:"12.345.678/0001-90"`                                 // Cadastro Nacional de Pessoa Jurídica
	DataCadastro time.Time `json:"data_cadastro" swaggertype:"string" example:"2026-01-01T10:00:00Z"` // Data de registro no sistema

	// Relacionamentos organizacionais (carregamento opcional)
	Setores                 []Setor                `json:"setores,omitempty" example:"[{\"id_setor\":2,\"nome_setor\":\"RH\"}]"`                               // Setores da empresa
	UsuariosAdministradores []UsuarioAdministrador `json:"usuarios_administradores,omitempty" example:"[{\"id_user_admin\":1,\"nome_admin\":\"Joao Silva\"}]"` // Administradores vinculados
}
