package response

type SetorResponse struct {
	ID             int              `json:"id_setor"`
	NomeSetor      string           `json:"nome_setor"`
	Descricao      string           `json:"descricao"`
	Empresa        *EmpresaResponse `json:"empresa,omitempty"`
	TotalPesquisas int              `json:"total_pesquisas,omitempty"`
}