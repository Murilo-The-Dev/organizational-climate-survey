// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// UsuarioAdministradorFilter define filtros para listagem de administradores de usuário.
type UsuarioAdministradorFilter struct {
	PaginationRequest           // Campos de paginação (herança)
	IDEmpresa *int    `form:"id_empresa"`                     // Filtra pelo ID da empresa
	Status    *string `form:"status" binding:"omitempty,oneof=Ativo Inativo Pendente"` // Filtra pelo status do usuário
	Email     *string `form:"email"`                           // Filtra pelo email
	Nome      *string `form:"nome"`                            // Filtra pelo nome
}
