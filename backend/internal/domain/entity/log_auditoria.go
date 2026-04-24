// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para auditoria do sistema de pesquisa.
package entity

import "time"

// LogAuditoria registra ações administrativas realizadas no sistema
type LogAuditoria struct {
	ID            int       `json:"id_log" example:"55"`                                           // Identificador único do registro
	IDUserAdmin   int       `json:"id_user_admin" example:"1"`                                     // ID do administrador responsável
	TimeStamp     time.Time `json:"timestamp" swaggertype:"string" example:"2026-01-10T09:00:00Z"` // Momento da ação
	AcaoRealizada string    `json:"acao_realizada" example:"Criacao de pesquisa"`                  // Descrição da operação executada
	Detalhes      string    `json:"detalhes" example:"Pesquisa trimestral criada"`                 // Informações complementares
	EnderecoIP    string    `json:"endereco_ip" example:"192.168.1.10"`                            // Endereço IP de origem

	// Relacionamento com administrador (carregamento opcional)
	UsuarioAdministrador *UsuarioAdministrador `json:"usuario_administrador,omitempty" example:"{\"id_user_admin\":1,\"nome_admin\":\"Joao Silva\"}"`
}
