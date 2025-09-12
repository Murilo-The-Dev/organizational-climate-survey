package dto

import (
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
	"time"
)

type PesquisaCreateRequest struct {
	IDEmpresa         int     `json:"id_empresa" binding:"required,gt=0"`
	IDUserAdmin       int     `json:"id_user_admin" binding:"required,gt=0"`
	IDSetor           int     `json:"id_setor" binding:"required,gt=0"`
	Titulo            string  `json:"titulo" binding:"required,min=3,max=255"`
	Descricao         string  `json:"descricao" binding:"max=1000"`
	Status            string  `json:"status" binding:"required,oneof=Rascunho Ativa Concluída Arquivada"`
	ConfigRecorrencia *string `json:"config_recorrencia,omitempty"`
	Anonimato         bool    `json:"anonimato"`
	DataAbertura      *string `json:"data_abertura,omitempty"`
	DataFechamento    *string `json:"data_fechamento,omitempty"`
}

type PesquisaUpdateRequest struct {
	Titulo            *string `json:"titulo,omitempty" binding:"omitempty,min=3,max=255"`
	Descricao         *string `json:"descricao,omitempty" binding:"omitempty,max=1000"`
	Status            *string `json:"status,omitempty" binding:"omitempty,oneof=Rascunho Ativa Concluída Arquivada"`
	ConfigRecorrencia *string `json:"config_recorrencia,omitempty"`
	DataAbertura      *string `json:"data_abertura,omitempty"`
	DataFechamento    *string `json:"data_fechamento,omitempty"`
}

func (r *PesquisaCreateRequest) ToEntity() (*entity.Pesquisa, error) {
	pesquisa := &entity.Pesquisa{
		IDEmpresa:         r.IDEmpresa,
		IDUserAdmin:       r.IDUserAdmin,
		IDSetor:           r.IDSetor,
		Titulo:            strings.TrimSpace(r.Titulo),
		Descricao:         strings.TrimSpace(r.Descricao),
		Status:            r.Status,
		ConfigRecorrencia: r.ConfigRecorrencia,
		Anonimato:         r.Anonimato,
	}

	if r.DataAbertura != nil {
		if t, err := time.Parse(time.RFC3339, *r.DataAbertura); err == nil {
			pesquisa.DataAbertura = &t
		} else {
			return nil, fmt.Errorf("data_abertura inválida: %v", err)
		}
	}

	if r.DataFechamento != nil {
		if t, err := time.Parse(time.RFC3339, *r.DataFechamento); err == nil {
			pesquisa.DataFechamento = &t
		} else {
			return nil, fmt.Errorf("data_fechamento inválida: %v", err)
		}
	}

	return pesquisa, nil
}

func (r *PesquisaUpdateRequest) ApplyToEntity(pesquisa *entity.Pesquisa) error {
	if r.Titulo != nil {
		pesquisa.Titulo = strings.TrimSpace(*r.Titulo)
	}
	if r.Descricao != nil {
		pesquisa.Descricao = strings.TrimSpace(*r.Descricao)
	}
	if r.Status != nil {
		pesquisa.Status = *r.Status
	}
	if r.ConfigRecorrencia != nil {
		pesquisa.ConfigRecorrencia = r.ConfigRecorrencia
	}

	if r.DataAbertura != nil {
		if t, err := time.Parse(time.RFC3339, *r.DataAbertura); err == nil {
			pesquisa.DataAbertura = &t
		} else {
			return fmt.Errorf("data_abertura inválida: %v", err)
		}
	}

	if r.DataFechamento != nil {
		if t, err := time.Parse(time.RFC3339, *r.DataFechamento); err == nil {
			pesquisa.DataFechamento = &t
		} else {
			return fmt.Errorf("data_fechamento inválida: %v", err)
		}
	}

	return nil
}