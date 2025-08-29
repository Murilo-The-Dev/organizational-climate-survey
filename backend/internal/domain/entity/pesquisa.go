package entity

import "time"

type Pesquisa struct {
    ID                int        `json:"id_pesquisa"`
    IDEmpresa         int        `json:"id_empresa"`
    IDUserAdmin       int        `json:"id_user_admin"`
    IDSetor           int        `json:"id_setor"`
    Titulo            string     `json:"titulo"`
    Descricao         string     `json:"descricao"`
    DataCriacao       time.Time  `json:"data_criacao"`
    DataAbertura      *time.Time `json:"data_abertura"`
    DataFechamento    *time.Time `json:"data_fechamento"`
    Status            string     `json:"status"` // Rascunho, Ativa, Concluída, Arquivada
    LinkAcesso        string     `json:"link_acesso"`
    QRCodePath        string     `json:"qrcode_path"`
    ConfigRecorrencia *string    `json:"config_recorrencia"` // JSON string
    Anonimato         bool       `json:"anonimato"`
    
    // Relacionamentos - carregados sob demanda
    Perguntas            []Pergunta           `json:"perguntas,omitempty"`
    UsuarioAdministrador *UsuarioAdministrador `json:"usuario_administrador,omitempty"`
    Setor                *Setor               `json:"setor,omitempty"`
    Dashboard            *Dashboard           `json:"dashboard,omitempty"`
}