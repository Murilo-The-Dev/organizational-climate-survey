package filter

type UsuarioAdministradorFilter struct {
	PaginationRequest
	IDEmpresa *int    `form:"id_empresa"`
	Status    *string `form:"status" binding:"omitempty,oneof=Ativo Inativo Pendente"`
	Email     *string `form:"email"`
	Nome      *string `form:"nome"`
}