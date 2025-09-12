package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

type SetorCreateRequest struct {
	IDEmpresa int    `json:"id_empresa" binding:"required,gt=0"`
	NomeSetor string `json:"nome_setor" binding:"required,min=2,max=255"`
	Descricao string `json:"descricao" binding:"max=500"`
}

type SetorUpdateRequest struct {
	NomeSetor *string `json:"nome_setor,omitempty" binding:"omitempty,min=2,max=255"`
	Descricao *string `json:"descricao,omitempty" binding:"omitempty,max=500"`
}

func (r *SetorCreateRequest) ToEntity() *entity.Setor {
	return &entity.Setor{
		IDEmpresa: r.IDEmpresa,
		NomeSetor: strings.TrimSpace(r.NomeSetor),
		Descricao: strings.TrimSpace(r.Descricao),
	}
}

func (r *SetorUpdateRequest) ApplyToEntity(setor *entity.Setor) {
	if r.NomeSetor != nil {
		setor.NomeSetor = strings.TrimSpace(*r.NomeSetor)
	}
	if r.Descricao != nil {
		setor.Descricao = strings.TrimSpace(*r.Descricao)
	}
}