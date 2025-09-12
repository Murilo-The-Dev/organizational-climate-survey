package filter

type PaginationRequest struct {
	Page    int    `form:"page" binding:"omitempty,gte=1"`
	Limit   int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	OrderBy string `form:"order_by"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
}

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

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}