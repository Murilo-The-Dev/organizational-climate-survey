package entity

import "time"

type UsuarioAdministrador struct {
    ID           int       `json:"id_user_admin"`
    IDEmpresa    int       `json:"id_empresa"`
    NomeAdmin    string    `json:"nome_admin"`
    Email        string    `json:"email"`
    SenhaHash    string    `json:"-"` // "-" oculta na serialização JSON
    DataCadastro time.Time `json:"data_cadastro"`
    Status       string    `json:"status"` // Ativo, Inativo, Pendente
    
    // Relacionamento inverso - opcional
    Empresa *Empresa `json:"empresa,omitempty"`
}