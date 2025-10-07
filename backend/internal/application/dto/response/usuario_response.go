// Package response contém structs utilizadas para formatar as respostas da API.
// UsuarioAdministradorResponse representa os dados de um usuário administrador 
// retornados pela API, incluindo informações básicas e a empresa associada.
package response

import "time"

// UsuarioAdministradorResponse representa os dados retornados de um usuário administrador.
type UsuarioAdministradorResponse struct {
	ID           int              `json:"id_user_admin"`      // ID único do usuário administrador
	NomeAdmin    string           `json:"nome_admin"`         // Nome completo do administrador
	Email        string           `json:"email"`              // E-mail de login do administrador
	DataCadastro time.Time        `json:"data_cadastro"`      // Data e hora de criação do registro
	Status       string           `json:"status"`             // Status atual do administrador (Ativo, Inativo, Pendente)
	Empresa      *EmpresaResponse `json:"empresa,omitempty"`  // Empresa associada ao usuário, se aplicável
}
