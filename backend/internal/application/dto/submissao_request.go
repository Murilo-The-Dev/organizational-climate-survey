// Package dto define objetos de transferência de dados para submissões.
package dto

// GenerateTokenRequest representa requisição para gerar token de acesso à pesquisa
type GenerateTokenRequest struct {
	Fingerprint string `json:"fingerprint"` // Fingerprint do browser (opcional)
}

// SubmitRespostasRequest representa requisição de submissão de respostas com token
type SubmitRespostasRequest struct {
	TokenAcesso string                   `json:"token_acesso"` // Token obtido via GenerateToken
	Respostas   []RespostaCreateRequest  `json:"respostas"`    // Lista de respostas da pesquisa
}