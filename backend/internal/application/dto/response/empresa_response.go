package response

import (
    "time"
    "organizational-climate-survey/backend/internal/domain/entity"
)

type EmpresaResponse struct {
    ID           int       `json:"id_empresa"`
    NomeFantasia string    `json:"nome_fantasia"`
    RazaoSocial  string    `json:"razao_social"`
    CNPJ         string    `json:"cnpj"`
    DataCadastro time.Time `json:"data_cadastro"`
    TotalSetores int       `json:"total_setores,omitempty"`
    TotalAdmins  int       `json:"total_admins,omitempty"`
    TotalPesquisas int     `json:"total_pesquisas,omitempty"`
}

// ToEmpresaResponse converte entity para response
func ToEmpresaResponse(empresa *entity.Empresa) EmpresaResponse {
    return EmpresaResponse{
        ID:           empresa.ID,
        NomeFantasia: empresa.NomeFantasia,
        RazaoSocial:  empresa.RazaoSocial,
        CNPJ:         empresa.CNPJ,
        DataCadastro: empresa.DataCadastro,
        // Os campos opcionais serão preenchidos quando necessário
    }
}