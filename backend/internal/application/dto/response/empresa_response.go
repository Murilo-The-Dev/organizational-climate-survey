// Package response contém structs usadas para enviar dados da API como respostas.
package response

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"time"
)

// EmpresaResponse representa os dados de uma empresa para resposta da API.
type EmpresaResponse struct {
	ID             int       `json:"id_empresa" example:"1"`                                            // ID da empresa
	NomeFantasia   string    `json:"nome_fantasia" example:"Acme"`                                      // Nome fantasia da empresa
	RazaoSocial    string    `json:"razao_social" example:"Acme LTDA"`                                  // Razão social da empresa
	CNPJ           string    `json:"cnpj" example:"12.345.678/0001-90"`                                 // CNPJ da empresa
	DataCadastro   time.Time `json:"data_cadastro" swaggertype:"string" example:"2026-01-01T10:00:00Z"` // Data de cadastro da empresa
	TotalSetores   int       `json:"total_setores,omitempty" example:"4"`                               // Número de setores, preenchido opcionalmente
	TotalAdmins    int       `json:"total_admins,omitempty" example:"2"`                                // Número de administradores, opcional
	TotalPesquisas int       `json:"total_pesquisas,omitempty" example:"8"`                             // Número de pesquisas, opcional
}

// ToEmpresaResponse converte a entidade Empresa para a struct de resposta
func ToEmpresaResponse(empresa *entity.Empresa) EmpresaResponse {
	return EmpresaResponse{
		ID:           empresa.ID,
		NomeFantasia: empresa.NomeFantasia,
		RazaoSocial:  empresa.RazaoSocial,
		CNPJ:         empresa.CNPJ,
		DataCadastro: empresa.DataCadastro,
		// Campos opcionais podem ser preenchidos posteriormente
	}
}
