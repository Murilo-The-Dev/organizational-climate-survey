package filter

type EmpresaFilter struct {
	PaginationRequest
	NomeFantasia *string `form:"nome_fantasia"`
	RazaoSocial  *string `form:"razao_social"`
	CNPJ         *string `form:"cnpj"`
	DataInicio   *string `form:"data_inicio"` // RFC3339
	DataFim      *string `form:"data_fim"`    // RFC3339
}