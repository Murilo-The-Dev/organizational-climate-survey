package filter

type LogAuditoriaFilter struct {
	PaginationRequest
	IDUserAdmin   *int    `form:"id_user_admin"`
	IDEmpresa     *int    `form:"id_empresa"`
	DataInicio    *string `form:"data_inicio"` // RFC3339
	DataFim       *string `form:"data_fim"`    // RFC3339
	AcaoRealizada *string `form:"acao_realizada"`
	EnderecoIP    *string `form:"endereco_ip"`
}