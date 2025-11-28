// Package dto contém estruturas de transferência de dados (Data Transfer Objects)
// empregadas na comunicação entre as camadas de entrada (handlers) e o domínio.
// Este arquivo define os DTOs de criação e atualização de dashboards,
// que representam visões analíticas associadas a uma pesquisa.

package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

// DashboardCreateRequest representa os dados necessários para criar um novo dashboard
// vinculado a uma pesquisa existente.
type DashboardCreateRequest struct {
	IDPesquisa    int     `json:"id_pesquisa" binding:"required"`             // Identificador da pesquisa associada (obrigatório)
	Titulo        string  `json:"titulo" binding:"required,min=3,max=255"`    // Título do dashboard (obrigatório)
	ConfigFiltros *string `json:"config_filtros,omitempty"`                   // Configuração JSON de filtros (opcional)
}

// DashboardUpdateRequest define os campos opcionais disponíveis para
// atualização de um dashboard existente.
type DashboardUpdateRequest struct {
	Titulo        *string `json:"titulo,omitempty" binding:"omitempty,min=3,max=255"` // Novo título (opcional)
	ConfigFiltros *string `json:"config_filtros,omitempty"`                           // Nova configuração de filtros (opcional)
}

// ToEntity converte o DTO de criação em uma instância de entidade Dashboard,
// aplicando limpeza de espaços em campos textuais.
func (r *DashboardCreateRequest) ToEntity() *entity.Dashboard {
	return &entity.Dashboard{
		IDPesquisa:    r.IDPesquisa,
		Titulo:        strings.TrimSpace(r.Titulo),
		ConfigFiltros: r.ConfigFiltros,
	}
}

// ApplyToEntity aplica as modificações contidas no DTO de atualização
// sobre uma entidade Dashboard existente, respeitando os campos não nulos.
func (r *DashboardUpdateRequest) ApplyToEntity(dashboard *entity.Dashboard) {
	if r.Titulo != nil {
		dashboard.Titulo = strings.TrimSpace(*r.Titulo)
	}
	if r.ConfigFiltros != nil {
		dashboard.ConfigFiltros = r.ConfigFiltros
	}
}
