// Package response contém structs utilizadas para formatar as respostas da API.
// UsuarioAdministradorResponse representa os dados de um usuário administrador
// retornados pela API, incluindo informações básicas e a empresa associada.
package response

import "time"

// UsuarioAdministradorResponse representa os dados retornados de um usuário administrador.
type UsuarioAdministradorResponse struct {
	ID           int              `json:"id_user_admin" example:"1"`                                                 // ID único do usuário administrador
	NomeAdmin    string           `json:"nome_admin" example:"Joao Silva"`                                           // Nome completo do administrador
	Email        string           `json:"email" example:"joao@empresa.com"`                                          // E-mail de login do administrador
	DataCadastro time.Time        `json:"data_cadastro" swaggertype:"string" example:"2026-01-01T10:00:00Z"`         // Data e hora de criação do registro
	Status       string           `json:"status" example:"Ativo"`                                                    // Status atual do administrador (Ativo, Inativo, Pendente)
	Empresa      *EmpresaResponse `json:"empresa,omitempty" example:"{\"id_empresa\":1,\"nome_fantasia\":\"Acme\"}"` // Empresa associada ao usuário, se aplicável
}
