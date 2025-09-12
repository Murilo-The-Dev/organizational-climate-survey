package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

type UsuarioAdministradorCreateRequest struct {
	IDEmpresa int    `json:"id_empresa" binding:"required,gt=0"`
	NomeAdmin string `json:"nome_admin" binding:"required,min=2,max=255"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Senha     string `json:"senha" binding:"required,min=8,max=128"`
	Status    string `json:"status" binding:"required,oneof=Ativo Inativo Pendente"`
}

type UsuarioAdministradorUpdateRequest struct {
	NomeAdmin *string `json:"nome_admin,omitempty" binding:"omitempty,min=2,max=255"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	Status    *string `json:"status,omitempty" binding:"omitempty,oneof=Ativo Inativo Pendente"`
}

type UsuarioAdministradorLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

func (r *UsuarioAdministradorCreateRequest) ToEntity(senhaHash string) *entity.UsuarioAdministrador {
	return &entity.UsuarioAdministrador{
		IDEmpresa: r.IDEmpresa,
		NomeAdmin: strings.TrimSpace(r.NomeAdmin),
		Email:     strings.ToLower(strings.TrimSpace(r.Email)),
		SenhaHash: senhaHash,
		Status:    r.Status,
	}
}

func (r *UsuarioAdministradorUpdateRequest) ApplyToEntity(user *entity.UsuarioAdministrador) {
	if r.NomeAdmin != nil {
		user.NomeAdmin = strings.TrimSpace(*r.NomeAdmin)
	}
	if r.Email != nil {
		user.Email = strings.ToLower(strings.TrimSpace(*r.Email))
	}
	if r.Status != nil {
		user.Status = *r.Status
	}
}