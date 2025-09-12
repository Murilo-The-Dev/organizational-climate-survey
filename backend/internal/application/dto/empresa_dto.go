package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

type EmpresaCreateRequest struct {
	NomeFantasia string `json:"nome_fantasia" binding:"required,min=2,max=255"`
	RazaoSocial  string `json:"razao_social" binding:"required,min=2,max=255"`
	CNPJ         string `json:"cnpj" binding:"required,len=18"`
}

type EmpresaUpdateRequest struct {
	NomeFantasia *string `json:"nome_fantasia,omitempty" binding:"omitempty,min=2,max=255"`
	RazaoSocial  *string `json:"razao_social,omitempty" binding:"omitempty,min=2,max=255"`
}

func (r *EmpresaCreateRequest) ToEntity() *entity.Empresa {
	return &entity.Empresa{
		NomeFantasia: strings.TrimSpace(r.NomeFantasia),
		RazaoSocial:  strings.TrimSpace(r.RazaoSocial),
		CNPJ:         r.CNPJ,
	}
}

func (r *EmpresaUpdateRequest) ApplyToEntity(empresa *entity.Empresa) {
	if r.NomeFantasia != nil {
		empresa.NomeFantasia = strings.TrimSpace(*r.NomeFantasia)
	}
	if r.RazaoSocial != nil {
		empresa.RazaoSocial = strings.TrimSpace(*r.RazaoSocial)
	}
}