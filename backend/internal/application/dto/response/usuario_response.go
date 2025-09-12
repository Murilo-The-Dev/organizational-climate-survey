package response

import "time"

type UsuarioAdministradorResponse struct {
	ID           int              `json:"id_user_admin"`
	NomeAdmin    string           `json:"nome_admin"`
	Email        string           `json:"email"`
	DataCadastro time.Time        `json:"data_cadastro"`
	Status       string           `json:"status"`
	Empresa      *EmpresaResponse `json:"empresa,omitempty"`
}