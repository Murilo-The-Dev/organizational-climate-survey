// Package dto contém estruturas responsáveis pela transferência de dados entre
// as camadas de entrada (handlers) e o domínio da aplicação. Este arquivo define
// os DTOs usados para criação e atualização de setores, garantindo validação e limpeza de dados.

package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

// SetorCreateRequest representa os dados necessários para criar um novo setor
// associado a uma empresa. Inclui validações de integridade e formato.
type SetorCreateRequest struct {
	IDEmpresa int    `json:"id_empresa" binding:"required,gt=0"`         // Identificador da empresa associada
	NomeSetor string `json:"nome_setor" binding:"required,min=2,max=255"` // Nome do setor (obrigatório e limitado)
	Descricao string `json:"descricao" binding:"max=500"`                 // Descrição opcional, com limite de tamanho
}

// SetorUpdateRequest define os campos opcionais para atualização de um setor existente.
// Usa ponteiros para diferenciar entre ausência de campo e valor vazio.
type SetorUpdateRequest struct {
	NomeSetor *string `json:"nome_setor,omitempty" binding:"omitempty,min=2,max=255"` // Novo nome do setor (opcional)
	Descricao *string `json:"descricao,omitempty" binding:"omitempty,max=500"`        // Nova descrição (opcional)
}

// ToEntity converte o DTO de criação em uma entidade de domínio Setor,
// aplicando trim e normalização de espaços em branco.
func (r *SetorCreateRequest) ToEntity() *entity.Setor {
	return &entity.Setor{
		IDEmpresa: r.IDEmpresa,
		NomeSetor: strings.TrimSpace(r.NomeSetor),
		Descricao: strings.TrimSpace(r.Descricao),
	}
}

// ApplyToEntity aplica as alterações do DTO de atualização sobre a entidade Setor existente.
// Apenas campos não nulos são atualizados, preservando valores anteriores.
func (r *SetorUpdateRequest) ApplyToEntity(setor *entity.Setor) {
	if r.NomeSetor != nil {
		setor.NomeSetor = strings.TrimSpace(*r.NomeSetor)
	}
	if r.Descricao != nil {
		setor.Descricao = strings.TrimSpace(*r.Descricao)
	}
}
