// Package filter contém structs usadas para filtros e paginação em consultas.
package filter

// PaginationRequest define os parâmetros de paginação e ordenação.
type PaginationRequest struct {
	Page    int    `form:"page" binding:"omitempty,gte=1"`           // Número da página, começa em 1
	Limit   int    `form:"limit" binding:"omitempty,gte=1,lte=100"` // Quantidade de itens por página
	OrderBy string `form:"order_by"`                                 // Campo para ordenação
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"` // Ordem ascendente ou descendente
}

// SetDefaults define valores padrão caso os campos não sejam fornecidos.
func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Order == "" {
		p.Order = "desc"
	}
}

// GetOffset calcula o offset usado em consultas paginadas (SQL, por exemplo).
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
