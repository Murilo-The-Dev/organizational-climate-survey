package filter

type SetorFilter struct {
	PaginationRequest
	IDEmpresa *int    `form:"id_empresa"`
	Nome      *string `form:"nome"`
}