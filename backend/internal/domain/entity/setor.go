// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados fundamentais do sistema de pesquisa de clima.
package entity

// Setor representa uma divisão ou departamento dentro de uma empresa
type Setor struct {
	ID        int    `json:"id_setor" example:"2"`                                        // Identificador único do setor
	IDEmpresa int    `json:"id_empresa" example:"1"`                                      // ID da empresa à qual pertence
	NomeSetor string `json:"nome_setor" example:"Recursos Humanos"`                       // Nome do setor/departamento
	Descricao string `json:"descricao" example:"Setor responsável por gestão de pessoas"` // Descrição detalhada do setor

	// Relacionamento com Empresa (opcional, para carregamento sob demanda)
	Empresa *Empresa `json:"empresa,omitempty" example:"{\"id_empresa\":1,\"nome_fantasia\":\"Acme\"}"` // Dados da empresa associada
}
