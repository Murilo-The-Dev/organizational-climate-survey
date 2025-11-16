// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para empresas do sistema de pesquisa.
package entity

import "time"

// Empresa representa uma organização cliente do sistema
type Empresa struct {
    ID           int       `json:"id_empresa"`      // Identificador único da empresa
    NomeFantasia string    `json:"nome_fantasia"`   // Nome comercial
    RazaoSocial  string    `json:"razao_social"`    // Nome jurídico registrado
    CNPJ         string    `json:"cnpj"`            // Cadastro Nacional de Pessoa Jurídica
    DataCadastro time.Time `json:"data_cadastro"`   // Data de registro no sistema
    
    // Relacionamentos organizacionais (carregamento opcional)
    Setores                 []Setor                 `json:"setores,omitempty"`                  // Setores da empresa
    UsuariosAdministradores []UsuarioAdministrador  `json:"usuarios_administradores,omitempty"` // Administradores vinculados
}