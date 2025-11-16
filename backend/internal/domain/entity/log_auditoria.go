// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para auditoria do sistema de pesquisa.
package entity

import "time"

// LogAuditoria registra ações administrativas realizadas no sistema
type LogAuditoria struct {
    ID            int       `json:"id_log"`            // Identificador único do registro
    IDUserAdmin   int       `json:"id_user_admin"`     // ID do administrador responsável
    TimeStamp     time.Time `json:"timestamp"`         // Momento da ação
    AcaoRealizada string    `json:"acao_realizada"`    // Descrição da operação executada
    Detalhes      string    `json:"detalhes"`          // Informações complementares
    EnderecoIP    string    `json:"endereco_ip"`       // Endereço IP de origem
    
    // Relacionamento com administrador (carregamento opcional)
    UsuarioAdministrador *UsuarioAdministrador `json:"usuario_administrador,omitempty"`
}