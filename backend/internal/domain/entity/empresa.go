package entity

import "time"

type Empresa struct {
    ID           int       `json:"id_empresa"`
    NomeFantasia string    `json:"nome_fantasia"`
    RazaoSocial  string    `json:"razao_social"`
    CNPJ         string    `json:"cnpj"`
    DataCadastro time.Time `json:"data_cadastro"`
    
    // Relacionamentos - carregados sob demanda
    Setores              []Setor              `json:"setores,omitempty"`
    UsuariosAdministradores []UsuarioAdministrador `json:"usuarios_administradores,omitempty"`
}