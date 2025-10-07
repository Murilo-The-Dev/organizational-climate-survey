// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

// Setor representa uma divisão ou departamento dentro de uma empresa
type Setor struct {
	ID        int    `json:"id_setor"`   // Identificador único do setor
	IDEmpresa int    `json:"id_empresa"` // ID da empresa à qual pertence
	NomeSetor string `json:"nome_setor"` // Nome do setor/departamento
	Descricao string `json:"descricao"`  // Descrição detalhada do setor

	// Relacionamento com Empresa (opcional, para carregamento sob demanda)
	Empresa *Empresa `json:"empresa,omitempty"` // Dados da empresa associada
}
