// Package response contém structs usadas para enviar dados da API como respostas.
// SetorResponse representa a estrutura de resposta para informações de um setor.
package response

// SetorResponse retorna os dados de um setor, incluindo informações da empresa
// e quantidade de pesquisas associadas.
type SetorResponse struct {
	ID             int              `json:"id_setor" example:"2"`                                                      // ID único do setor
	NomeSetor      string           `json:"nome_setor" example:"Recursos Humanos"`                                     // Nome do setor
	Descricao      string           `json:"descricao" example:"Gestão de pessoas"`                                     // Descrição do setor
	Empresa        *EmpresaResponse `json:"empresa,omitempty" example:"{\"id_empresa\":1,\"nome_fantasia\":\"Acme\"}"` // Informações da empresa associada, opcional
	TotalPesquisas int              `json:"total_pesquisas,omitempty" example:"5"`                                     // Quantidade total de pesquisas vinculadas ao setor, opcional
}
