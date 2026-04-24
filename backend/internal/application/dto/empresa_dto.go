// Package dto contém estruturas de transferência de dados (Data Transfer Objects)
// responsáveis por mediar a comunicação entre as camadas externas e o domínio.
// Este arquivo define os DTOs para criação e atualização de entidades do tipo Empresa.

package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

// EmpresaCreateRequest representa os dados necessários para o registro
// de uma nova empresa no sistema.
type EmpresaCreateRequest struct {
	NomeFantasia string `json:"nome_fantasia" binding:"required,min=2,max=255" example:"Acme"`               // Nome comercial da empresa (obrigatório)
	RazaoSocial  string `json:"razao_social" binding:"required,min=2,max=255" example:"Acme Industria LTDA"` // Razão social registrada (obrigatória)
	CNPJ         string `json:"cnpj" binding:"required,len=18" example:"12.345.678/0001-90"`                 // CNPJ formatado (obrigatório e com validação de tamanho)
}

// EmpresaUpdateRequest define os campos opcionais permitidos na atualização
// de uma empresa existente.
type EmpresaUpdateRequest struct {
	NomeFantasia *string `json:"nome_fantasia,omitempty" binding:"omitempty,min=2,max=255" example:"Acme Brasil"`     // Novo nome fantasia (opcional)
	RazaoSocial  *string `json:"razao_social,omitempty" binding:"omitempty,min=2,max=255" example:"Acme Brasil LTDA"` // Nova razão social (opcional)
}

// ToEntity converte o DTO de criação em uma instância da entidade de domínio Empresa,
// aplicando limpeza de espaços em campos textuais.
func (r *EmpresaCreateRequest) ToEntity() *entity.Empresa {
	return &entity.Empresa{
		NomeFantasia: strings.TrimSpace(r.NomeFantasia),
		RazaoSocial:  strings.TrimSpace(r.RazaoSocial),
		CNPJ:         r.CNPJ,
	}
}

// ApplyToEntity aplica as modificações do DTO de atualização sobre uma
// instância existente de Empresa, apenas nos campos não nulos.
func (r *EmpresaUpdateRequest) ApplyToEntity(empresa *entity.Empresa) {
	if r.NomeFantasia != nil {
		empresa.NomeFantasia = strings.TrimSpace(*r.NomeFantasia)
	}
	if r.RazaoSocial != nil {
		empresa.RazaoSocial = strings.TrimSpace(*r.RazaoSocial)
	}
}
