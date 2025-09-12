package filter

type PesquisaFilter struct {
	PaginationRequest
	IDEmpresa      *int    `form:"id_empresa"`
	IDSetor        *int    `form:"id_setor"`
	IDUserAdmin    *int    `form:"id_user_admin"`
	Status         *string `form:"status" binding:"omitempty,oneof=Rascunho Ativa Concluída Arquivada"`
	DataInicio     *string `form:"data_inicio"` // RFC3339
	DataFim        *string `form:"data_fim"`    // RFC3339
	Titulo         *string `form:"titulo"`
	Anonimato      *bool   `form:"anonimato"`
	ComPerguntas   *bool   `form:"com_perguntas"`
	ComRespostas   *bool   `form:"com_respostas"`
}