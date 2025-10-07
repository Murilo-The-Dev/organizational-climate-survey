// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

import "time"

// Pesquisa representa uma pesquisa de clima organizacional
type Pesquisa struct {
	ID                int        `json:"id_pesquisa"`        // Identificador único da pesquisa
	IDEmpresa         int        `json:"id_empresa"`         // ID da empresa responsável
	IDUserAdmin       int        `json:"id_user_admin"`      // ID do administrador criador
	IDSetor           int        `json:"id_setor"`           // ID do setor alvo
	Titulo            string     `json:"titulo"`             // Título da pesquisa
	Descricao         string     `json:"descricao"`          // Descrição detalhada
	DataCriacao       time.Time  `json:"data_criacao"`       // Data de criação
	DataAbertura      *time.Time `json:"data_abertura"`      // Data de início opcional
	DataFechamento    *time.Time `json:"data_fechamento"`    // Data de término opcional
	Status            string     `json:"status"`             // Estado atual (Rascunho, Ativa, etc)
	LinkAcesso        string     `json:"link_acesso"`        // Link único para respostas
	QRCodePath        string     `json:"qrcode_path"`        // Caminho do QR Code gerado
	ConfigRecorrencia *string    `json:"config_recorrencia"` // Configuração de recorrência
	Anonimato         bool       `json:"anonimato"`          // Se respostas são anônimas

	// Relacionamentos (opcional, para carregamento sob demanda)
	Perguntas            []Pergunta            `json:"perguntas,omitempty"`             // Lista de perguntas
	UsuarioAdministrador *UsuarioAdministrador `json:"usuario_administrador,omitempty"` // Admin criador
	Setor                *Setor                `json:"setor,omitempty"`                 // Setor alvo
	Dashboard            *Dashboard            `json:"dashboard,omitempty"`             // Dashboard associado
}
