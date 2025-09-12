package dto

import (
    "organizational-climate-survey/backend/internal/domain/entity"
    "strings"
)

type DashboardCreateRequest struct {
    IDPesquisa    int     `json:"id_pesquisa" binding:"required"`
    Titulo        string  `json:"titulo" binding:"required,min=3,max=255"`
    ConfigFiltros *string `json:"config_filtros,omitempty"` // Ponteiro para match com entity
}

type DashboardUpdateRequest struct {
    Titulo        *string `json:"titulo,omitempty" binding:"omitempty,min=3,max=255"`
    ConfigFiltros *string `json:"config_filtros,omitempty"`
}

func (r *DashboardCreateRequest) ToEntity() *entity.Dashboard {
    return &entity.Dashboard{
        IDPesquisa:    r.IDPesquisa,
        Titulo:        strings.TrimSpace(r.Titulo),
        ConfigFiltros: r.ConfigFiltros, // Já é ponteiro
    }
}

func (r *DashboardUpdateRequest) ApplyToEntity(dashboard *entity.Dashboard) {
    if r.Titulo != nil {
        dashboard.Titulo = strings.TrimSpace(*r.Titulo)
    }
    if r.ConfigFiltros != nil {
        dashboard.ConfigFiltros = r.ConfigFiltros
    }
}