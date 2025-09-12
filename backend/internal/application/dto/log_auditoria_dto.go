package dto

import "organizational-climate-survey/backend/internal/domain/entity"

type RetentionRequest struct {
	RetentionDays int `json:"retention_days" binding:"required,min=30,max=2555"`
}

type LogAuditoriaCreateRequest struct {
	IDUserAdmin   int    `json:"id_user_admin" binding:"required,gt=0"`
	AcaoRealizada string `json:"acao_realizada" binding:"required,min=3,max=255"`
	Detalhes      string `json:"detalhes" binding:"max=1000"`
	EnderecoIP    string `json:"endereco_ip" binding:"omitempty,ip"`
}

func (r *LogAuditoriaCreateRequest) ToEntity() *entity.LogAuditoria {
	return &entity.LogAuditoria{
		IDUserAdmin:   r.IDUserAdmin,
		AcaoRealizada: r.AcaoRealizada,
		Detalhes:      r.Detalhes,
		EnderecoIP:    r.EnderecoIP,
	}
}