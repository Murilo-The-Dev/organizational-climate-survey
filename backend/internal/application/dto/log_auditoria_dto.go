// Package dto contém estruturas de transferência de dados (Data Transfer Objects)
// utilizadas para comunicação entre as camadas externas e o domínio da aplicação.
// Este arquivo define os DTOs relacionados à retenção de dados e logs de auditoria.

package dto

import "organizational-climate-survey/backend/internal/domain/entity"

// RetentionRequest representa a configuração de retenção de dados,
// especificando o período de armazenamento permitido em dias.
type RetentionRequest struct {
	RetentionDays int `json:"retention_days" binding:"required,min=30,max=2555" example:"90"` // Quantidade de dias de retenção (mínimo 30, máximo 2555)
}

// LogAuditoriaCreateRequest representa os dados necessários para registrar
// uma ação administrativa no log de auditoria do sistema.
type LogAuditoriaCreateRequest struct {
	IDUserAdmin   int    `json:"id_user_admin" binding:"required,gt=0" example:"1"`                             // Identificador do administrador responsável pela ação
	AcaoRealizada string `json:"acao_realizada" binding:"required,min=3,max=255" example:"Criacao de pesquisa"` // Descrição resumida da ação executada
	Detalhes      string `json:"detalhes" binding:"max=1000" example:"Pesquisa trimestral criada com sucesso"`  // Detalhamento adicional da ação (opcional)
	EnderecoIP    string `json:"endereco_ip" binding:"omitempty,ip" example:"192.168.1.10"`                     // Endereço IP de origem (opcional e validado)
}

// ToEntity converte a requisição de criação de log em uma entidade de domínio LogAuditoria,
// pronta para persistência.
func (r *LogAuditoriaCreateRequest) ToEntity() *entity.LogAuditoria {
	return &entity.LogAuditoria{
		IDUserAdmin:   r.IDUserAdmin,
		AcaoRealizada: r.AcaoRealizada,
		Detalhes:      r.Detalhes,
		EnderecoIP:    r.EnderecoIP,
	}
}
