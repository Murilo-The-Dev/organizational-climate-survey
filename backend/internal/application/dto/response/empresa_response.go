// Package response contém structs usadas para enviar dados da API como respostas.
package response

import (
    "time"
    "organizational-climate-survey/backend/internal/domain/entity"
)

// EmpresaResponse representa os dados de uma empresa para resposta da API.
type EmpresaResponse struct {
    ID             int       `json:"id_empresa"`            // ID da empresa
    NomeFantasia   string    `json:"nome_fantasia"`         // Nome fantasia da empresa
    RazaoSocial    string    `json:"razao_social"`          // Razão social da empresa
    CNPJ           string    `json:"cnpj"`                  // CNPJ da empresa
    DataCadastro   time.Time `json:"data_cadastro"`         // Data de cadastro da empresa
    TotalSetores   int       `json:"total_setores,omitempty"`  // Número de setores, preenchido opcionalmente
    TotalAdmins    int       `json:"total_admins,omitempty"`   // Número de administradores, opcional
    TotalPesquisas int       `json:"total_pesquisas,omitempty"` // Número de pesquisas, opcional
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
